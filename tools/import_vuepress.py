#!/usr/bin/env python3
"""Import the legacy VuePress (vdoing) site into the blog through its admin API.

Rules are documented in docs/MIGRATION.md. The source tree is READ-ONLY here.

    python3 tools/import_vuepress.py            # dry-run: parse + report only (default)
    python3 tools/import_vuepress.py --apply    # write through the API (idempotent, by slug)
"""
import argparse, hashlib, json, mimetypes, os, re, sys, urllib.error, urllib.request, uuid
from collections import Counter, defaultdict

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DEFAULT_SRC = os.path.join(ROOT, "Vilivala")
REPORT_DIR = os.path.join(ROOT, "tools", "reports")
EARLIEST = "2020-01-12T11:51:53+08:00"

# ---- rule tables (MIGRATION.md §2–§4) ---------------------------------------------------
PAGES = {  # source file (relative to docs/) -> page slug
    "05.about/01.关于.md": "about",
    "04.more/1.日志.md": "changelog",
    "04.more/2.赞助&支持.md": "sponsor",
    "04.more/3.友情链接.md": "links",
}
SLUG_WORDS = {  # readable slugs for Chinese category names
    "随笔": "essay", "人设": "characters", "素材": "material", "诗": "poetry", "摘抄": "excerpts",
    "笔记": "notes", "世界观体系": "worldview", "书单": "booklist", "区块链投资": "blockchain",
    "教程": "tutorial", "英语四六级": "cet", "技术": "tech", "版本控制与管理": "vcs",
    "博客开发指南": "blog-dev", "测评": "review",
}


def map_category(cats, rel):
    """Normalise a legacy categories list to a canonical path of up to 3 names."""
    c = list(cats)
    if not c:
        return None
    one = c[0] if len(c) == 1 else None
    if one == "随笔":
        return ["fun", "随笔"]
    if one in ("《笔记》", "笔记"):
        return ["ToRead", "笔记"]
    if one == "《人物设定集》":
        return ["ToRead", "人设"]
    if one == "世界体系":
        return ["ToRead", "世界观体系"]
    if one == "categories":
        return ["ToRead", "书单"]
    if one in ("git", "licence", "theme"):
        base = ["code", "版本控制与管理"]
        return base + ["博客开发指南"] if "博客开发指南" in rel else base
    if c == ["fun", "测评"]:
        return ["code", "测评"]
    return c[:3]


def ascii_slug(s):
    return re.sub(r"[^a-z0-9]+", "-", s.lower()).strip("-")


def cat_slug(path):
    parts = []
    for name in path:
        a = ascii_slug(name)
        parts.append(a or SLUG_WORDS.get(name) or "c" + hashlib.md5(name.encode()).hexdigest()[:5])
    return "-".join(parts)


def tag_slug(name):
    return ascii_slug(name) or "t" + hashlib.md5(name.encode()).hexdigest()[:6]


# ---- parsing -----------------------------------------------------------------------------
def parse_front(text):
    m = re.match(r"^---\r?\n(.*?)\r?\n---\r?\n?(.*)$", text, re.S)
    if not m:
        return {}, text
    meta, key = {}, None
    for line in m.group(1).splitlines():
        km = re.match(r"^([A-Za-z_]\w*):\s*(.*)$", line)
        if km:
            key, val = km.group(1), km.group(2).strip()
            if val.startswith("[") and val.endswith("]"):
                meta[key] = [x.strip().strip("'\"") for x in val[1:-1].split(",") if x.strip()]
            elif val == "":
                meta[key] = []  # maybe a block list; items appended below
            else:
                meta[key] = val.strip("'\"")
        elif key and re.match(r"^\s*-", line) and isinstance(meta.get(key), list):
            meta[key].append(line.strip()[1:].strip().strip("'\""))
    return meta, m.group(2)


def parse_date(raw):
    if not raw:
        return None
    m = re.match(r"^(\d{4})-(\d{1,2})-(\d{1,2})(?:[ T](\d{1,2}):(\d{1,2})(?::(\d{1,2}))?)?", raw)
    if not m:
        return None
    y, mo, d, h, mi, s = (int(x) if x else 0 for x in m.groups())
    return f"{y:04d}-{mo:02d}-{d:02d}T{h:02d}:{mi:02d}:{s:02d}+08:00"


def clean_content(body):
    out = []
    for line in body.replace("\r\n", "\n").split("\n"):
        if re.fullmatch(r"\s*(</?br\s*/?>\s*)+", line):
            out.append("")  # a line holding only <br> tags becomes a blank line
        else:
            out.append(re.sub(r"</?br\s*/?>", "  \n", line))
    return "\n".join(out).strip() + "\n"


def make_summary(body, n=120):
    t = re.sub(r"```.*?```", " ", body, flags=re.S)
    t = re.sub(r"!\[[^\]]*\]\([^)]*\)", " ", t)
    t = re.sub(r"\[([^\]]*)\]\([^)]*\)", r"\1", t)
    t = re.sub(r"<[^>]+>", " ", t)
    t = re.sub(r"^[#>\-\*\s]+", "", t, flags=re.M)
    t = re.sub(r"[*_`~|]", "", t)
    t = re.sub(r"\s+", " ", t).strip()
    return t[:n]


IMG_RE = re.compile(r"""(?:<img[^>]+src=["']([^"']+)["'])|(?:!\[[^\]]*\]\(([^)\s]+))""")


def local_images(body):
    return [a or b for a, b in IMG_RE.findall(body) if (a or b) and not re.match(r"^(https?:)?//|^data:", a or b)]


def resolve_image(ref, src_root, file_dir):
    ref = ref.split("?")[0].split("#")[0]
    cands = [os.path.join(file_dir, ref),
             os.path.join(src_root, "docs", ".vuepress", "public", ref.lstrip("/")),
             os.path.join(src_root, ref.lstrip("/")),
             os.path.join(src_root, "diagram", os.path.basename(ref))]
    return next((c for c in cands if os.path.isfile(c)), None)


# ---- scan --------------------------------------------------------------------------------
def scan(src):
    docs = os.path.join(src, "docs")
    posts, pages, skipped = [], [], []
    files = []
    for dp, dn, fn in os.walk(docs):
        dn[:] = sorted(d for d in dn if d != ".vuepress")
        for f in sorted(fn):
            if f.endswith(".md"):
                files.append(os.path.join(dp, f))
    for p in files:
        rel = os.path.relpath(p, docs)
        meta, body = parse_front(open(p, encoding="utf-8", errors="replace").read())
        title = meta.get("title", "")
        if rel in PAGES:
            pages.append(dict(slug=PAGES[rel], title=title or PAGES[rel], content=clean_content(body), rel=rel))
            continue
        if rel == "index.md" or meta.get("home") or "pageComponent" in meta or meta.get("article") == "false" \
                or rel.startswith("@pages") or rel.startswith("00.首页"):
            skipped.append((rel, "非文章页（首页/板块入口/站内索引），由新站取代"))
            continue
        cats = [c for c in (meta.get("categories") or []) if c]
        path = map_category(cats, rel)
        if not path:
            skipped.append((rel, "没有分类"))
            continue
        if len(body.strip()) < 30:
            skipped.append((rel, "正文为空"))
            continue
        tags = []
        for t in meta.get("tags") or []:
            if t and t not in tags:
                tags.append(t)
        posts.append(dict(rel=rel, dir=os.path.dirname(p), title=title or os.path.splitext(os.path.basename(rel))[0],
                          permalink=meta.get("permalink", ""), date_raw=meta.get("date", ""),
                          legacy_cats=cats, cat_path=path, tags=tags, is_top="sticky" in meta,
                          body=body))
    return posts, pages, skipped


def finalize(posts, src):
    """Slugs, dates, summaries, images: everything that needs cross-post context."""
    notes = defaultdict(list)
    seen = Counter()
    for p in posts:
        base = re.sub(r"^/?pages/", "", p["permalink"]).strip("/") or "p" + hashlib.md5(p["rel"].encode()).hexdigest()[:6]
        base = re.sub(r"[^A-Za-z0-9_-]+", "-", base).strip("-") or "post"
        seen[base] += 1
        p["slug"] = base if seen[base] == 1 else f"{base}-{seen[base]}"
        if seen[base] > 1:
            notes["slug_conflicts"].append(f"{p['rel']}  →  {p['slug']}（与 {base} 冲突）")
    by_dir = defaultdict(list)
    for p in posts:
        by_dir[p["dir"]].append(p)
    for p in posts:
        d = parse_date(p["date_raw"])
        if d and p["date_raw"] != d[:19].replace("T", " "):
            notes["date_normalised"].append(f"{p['rel']}: '{p['date_raw']}' → {d}")
        if not d:
            sib = [q for q in by_dir[p["dir"]] if q is not p and parse_date(q["date_raw"])]
            before = [q for q in sib if q["rel"] < p["rel"]]
            pick = (before[-1] if before else (sib[0] if sib else None))
            d = parse_date(pick["date_raw"]) if pick else EARLIEST
            notes["date_missing"].append(f"{p['rel']} → {d}（{'取同目录相邻文章' if pick else '站点最早日期'}）")
        p["created_at"] = d
        p["content"] = clean_content(p["body"])
        p["summary"] = make_summary(p["body"])
        p["images"] = local_images(p["body"])
    return notes


# ---- report ------------------------------------------------------------------------------
def build_report(src, posts, pages, skipped, notes, applied=None):
    L = []
    w = L.append
    w(f"# 导入{'结果' if applied else '预演'}报告（{'已写入' if applied else 'dry-run，未写入任何数据'}）\n")
    w(f"- 源：`{src}`")
    w(f"- 将导入文章：**{len(posts)}**；独立页面：**{len(pages)}**；跳过：**{len(skipped)}**\n")
    tree = defaultdict(Counter)
    for p in posts:
        for i in range(1, len(p["cat_path"]) + 1):
            tree[tuple(p["cat_path"][:i])]["total"] += 1
        tree[tuple(p["cat_path"])]["direct"] += 1
    w("## 分类树（直属 / 含子类）\n")
    w("```")
    for k in sorted(tree):
        w(f"{'  ' * (len(k) - 1)}{k[-1]}  ({tree[k]['direct']} / {tree[k]['total']})   slug={cat_slug(list(k))}")
    w("```\n")
    changed = [p for p in posts if p["legacy_cats"] != p["cat_path"]]
    w(f"## 分类被归并的文章（{len(changed)}）\n")
    grp = defaultdict(list)
    for p in changed:
        grp[(' / '.join(p['legacy_cats']), ' / '.join(p['cat_path']))].append(p["title"])
    for (a, b), ts in sorted(grp.items()):
        w(f"- `{a}` → `{b}`：{len(ts)} 篇（如《{ts[0]}》）")
    w("")
    for key, title in [("slug_conflicts", "slug 冲突已改名"), ("date_missing", "缺失日期已补全"), ("date_normalised", "日期格式已规范")]:
        w(f"## {title}（{len(notes[key])}）\n")
        for x in notes[key][:30]:
            w(f"- {x}")
        w("")
    tags = Counter(t for p in posts for t in p["tags"])
    w(f"## 标签：{len(tags)} 个（已丢弃空标签），前 15：{', '.join(f'{t}×{n}' for t, n in tags.most_common(15))}\n")
    w(f"## 置顶文章（{sum(p['is_top'] for p in posts)}）\n")
    for p in posts:
        if p["is_top"]:
            w(f"- 《{p['title']}》")
    imgs = [(p["rel"], i, resolve_image(i, src, p["dir"])) for p in posts for i in p["images"]]
    w(f"\n## 本地图片（{len(imgs)}）\n")
    for rel, ref, found in imgs:
        w(f"- {rel}: `{ref}` → {'找到 ' + os.path.relpath(found, src) if found else '**未找到，保留原文**'}")
    w(f"\n## 独立页面（{len(pages)}）\n")
    for pg in pages:
        w(f"- `{pg['slug']}`：{pg['title']}（{len(pg['content'])} 字符）")
    w(f"\n## 跳过（{len(skipped)}）\n")
    for rel, why in skipped:
        w(f"- `{rel}` — {why}")
    if applied:
        w("\n## 写入结果\n")
        for k, v in applied.items():
            w(f"- {k}: {v}")
    return "\n".join(L) + "\n"


# ---- API client --------------------------------------------------------------------------
class Api:
    def __init__(self, base, user, pw):
        self.base, self.token = base.rstrip("/"), None
        self.token = self.call("POST", "/api/auth/login", {"username": user, "password": pw})["token"]

    def call(self, method, path, body=None, raw=None, headers=None):
        h = {"Authorization": f"Bearer {self.token}"} if self.token else {}
        data = raw
        if body is not None:
            data, h["Content-Type"] = json.dumps(body).encode(), "application/json"
        h.update(headers or {})
        req = urllib.request.Request(self.base + path, data=data, method=method, headers=h)
        try:
            resp = json.load(urllib.request.urlopen(req, timeout=30))
        except urllib.error.HTTPError as e:
            resp = json.load(e)
        if resp.get("code") != 0:
            raise RuntimeError(f"{method} {path}: {resp.get('message')}")
        return resp["data"]

    def upload(self, filepath):
        b = uuid.uuid4().hex
        ct = mimetypes.guess_type(filepath)[0] or "application/octet-stream"
        body = (f"--{b}\r\nContent-Disposition: form-data; name=\"file\"; filename=\"{os.path.basename(filepath)}\"\r\n"
                f"Content-Type: {ct}\r\n\r\n").encode() + open(filepath, "rb").read() + f"\r\n--{b}--\r\n".encode()
        return self.call("POST", "/api/admin/upload", raw=body, headers={"Content-Type": f"multipart/form-data; boundary={b}"})["url"]


def apply(api, src, posts, pages):
    stats = Counter()
    # categories (reuse by parent + name)
    existing = {(c["parent_id"], c["name"]): c["id"] for c in api.call("GET", "/api/categories")}
    slugs_used = {c["slug"] for c in api.call("GET", "/api/categories")}

    def ensure_cat(path):
        parent = 0
        for i, name in enumerate(path):
            key = (parent, name)
            if key not in existing:
                slug = cat_slug(path[: i + 1])
                n = 2
                while slug in slugs_used:
                    slug, n = f"{cat_slug(path[: i + 1])}-{n}", n + 1
                existing[key] = api.call("POST", "/api/admin/categories", {"name": name, "slug": slug, "parent_id": parent})["id"]
                slugs_used.add(slug)
                stats["分类新建"] += 1
            parent = existing[key]
        return parent

    tag_ids = {t["name"]: t["id"] for t in api.call("GET", "/api/tags")}

    def ensure_tag(name):
        if name not in tag_ids:
            tag_ids[name] = api.call("POST", "/api/admin/tags", {"name": name, "slug": tag_slug(name)})["id"]
            stats["标签新建"] += 1
        return tag_ids[name]

    by_slug, page = {}, 1
    while True:
        d = api.call("GET", f"/api/admin/posts?page={page}&page_size=50")
        by_slug.update({p["slug"]: p["id"] for p in d["list"]})
        if page * 50 >= d["total"]:
            break
        page += 1

    img_cache = {}
    for p in posts:
        try:
            content = p["content"]
            for ref in dict.fromkeys(p["images"]):
                f = resolve_image(ref, src, p["dir"])
                if not f:
                    stats["图片未找到"] += 1
                    continue
                if f not in img_cache:
                    img_cache[f] = api.upload(f)
                    stats["图片上传"] += 1
                content = content.replace(ref, img_cache[f])
            payload = dict(title=p["title"], slug=p["slug"], summary=p["summary"], content=content, status="published",
                           category_id=ensure_cat(p["cat_path"]), tag_ids=[ensure_tag(t) for t in p["tags"]],
                           is_top=p["is_top"], created_at=p["created_at"])
            if p["slug"] in by_slug:
                api.call("PUT", f"/api/admin/posts/{by_slug[p['slug']]}", payload)
                stats["文章更新"] += 1
            else:
                api.call("POST", "/api/admin/posts", payload)
                stats["文章新建"] += 1
        except Exception as e:  # keep going; report at the end
            stats["文章失败"] += 1
            print(f"  ! {p['rel']}: {e}", file=sys.stderr)
    for pg in pages:
        api.call("PUT", f"/api/admin/pages/{pg['slug']}", {"title": pg["title"][:64], "content": pg["content"]})
        stats["页面写入"] += 1
    return dict(stats)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--src", default=DEFAULT_SRC)
    ap.add_argument("--api", default="http://localhost:8080")
    ap.add_argument("--user", default="admin")
    ap.add_argument("--password", default=os.environ.get("BLOG_ADMIN_PASSWORD", "admin123"))
    ap.add_argument("--apply", action="store_true", help="write through the API (default: dry-run)")
    a = ap.parse_args()

    posts, pages, skipped = scan(a.src)
    notes = finalize(posts, a.src)
    os.makedirs(REPORT_DIR, exist_ok=True)
    result = None
    if a.apply:
        result = apply(Api(a.api, a.user, a.password), a.src, posts, pages)
    report = build_report(a.src, posts, pages, skipped, notes, result)
    out = os.path.join(REPORT_DIR, "import-result.md" if a.apply else "dry-run.md")
    open(out, "w", encoding="utf-8").write(report)
    print(report if len(report) < 6000 else report[:6000] + "\n…（完整报告见文件）")
    print(f"\n报告已写入: {out}")


if __name__ == "__main__":
    main()
