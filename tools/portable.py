#!/usr/bin/env python3
"""
Cross-platform migration & multi-OS build helper for Vilivala (macOS / Linux / Windows).

Commands:
  python3 tools/portable.py pack                  # Export MySQL DB + uploads/ into a portable .tar.gz bundle
  python3 tools/portable.py unpack <bundle.tar.gz># Restore MySQL DB + uploads/ from a bundle
  python3 tools/portable.py build-all             # Build frontend dist/ + Linux/macOS/Windows Go binaries
"""

from __future__ import annotations

import argparse
import datetime as dt
import os
import shutil
import subprocess
import sys
import tarfile
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BACKEND_DIR = ROOT / "backend"
FRONTEND_DIR = ROOT / "frontend"
BACKUPS_DIR = ROOT / "backups"
DIST_RELEASE_DIR = ROOT / "release"


def find_bin(name: str, fallbacks: list[str]) -> str:
    found = shutil.which(name)
    if found:
        return found
    for fb in fallbacks:
        if Path(fb).exists():
            return fb
    raise FileNotFoundError(f"Required executable '{name}' not found in PATH.")


def db_env() -> dict[str, str]:
    return {
        "host": os.environ.get("MYSQL_HOST", "127.0.0.1"),
        "port": os.environ.get("MYSQL_PORT", "3306"),
        "user": os.environ.get("MYSQL_USER", "bloguser"),
        "password": os.environ.get("MYSQL_PASSWORD", "blogpass123"),
        "database": os.environ.get("MYSQL_DATABASE", "personal_blog"),
    }


def cmd_pack(_args: argparse.Namespace) -> None:
    BACKUPS_DIR.mkdir(parents=True, exist_ok=True)
    mysqldump = find_bin(
        "mysqldump",
        ["/opt/homebrew/opt/mysql@8.0/bin/mysqldump", "/usr/local/mysql/bin/mysqldump"],
    )
    cfg = db_env()
    stamp = dt.datetime.now().strftime("%Y%m%d-%H%M%S")
    bundle_path = BACKUPS_DIR / f"vilivala-migration-{stamp}.tar.gz"

    with tempfile.TemporaryDirectory() as tmp:
        sql_path = Path(tmp) / "database.sql"
        env = os.environ.copy()
        env["MYSQL_PWD"] = cfg["password"]
        dump_cmd = [
            mysqldump,
            "-h",
            cfg["host"],
            "-P",
            cfg["port"],
            "-u",
            cfg["user"],
            "--single-transaction",
            "--no-tablespaces",
            "--routines",
            "--triggers",
            "--default-character-set=utf8mb4",
            cfg["database"],
        ]
        with sql_path.open("wb") as f:
            subprocess.run(dump_cmd, stdout=f, env=env, check=True)

        uploads_dir = BACKEND_DIR / "uploads"
        with tarfile.open(bundle_path, "w:gz") as tar:
            tar.add(sql_path, arcname="database.sql")
            if uploads_dir.exists():
                tar.add(uploads_dir, arcname="uploads")

    size_mb = bundle_path.stat().st_size / (1024 * 1024)
    print(f"[OK] Portable migration bundle created: {bundle_path} ({size_mb:.2f} MB)")


def cmd_unpack(args: argparse.Namespace) -> None:
    bundle = Path(args.bundle).resolve()
    if not bundle.exists():
        sys.exit(f"[ERROR] Bundle not found: {bundle}")

    mysql_bin = find_bin(
        "mysql",
        ["/opt/homebrew/opt/mysql@8.0/bin/mysql", "/usr/local/mysql/bin/mysql"],
    )
    cfg = db_env()

    with tempfile.TemporaryDirectory() as tmp:
        tmp_dir = Path(tmp)
        with tarfile.open(bundle, "r:gz") as tar:
            # Safe member extraction (prevent path traversal)
            for member in tar.getmembers():
                target = (tmp_dir / member.name).resolve()
                if not str(target).startswith(str(tmp_dir.resolve())):
                    sys.exit(f"[ERROR] Unsafe archive path: {member.name}")
            tar.extractall(tmp_dir)

        sql_path = tmp_dir / "database.sql"
        if sql_path.exists():
            env = os.environ.copy()
            env["MYSQL_PWD"] = cfg["password"]
            with sql_path.open("rb") as f:
                subprocess.run(
                    [
                        mysql_bin,
                        "-h",
                        cfg["host"],
                        "-P",
                        cfg["port"],
                        "-u",
                        cfg["user"],
                        "--default-character-set=utf8mb4",
                        cfg["database"],
                    ],
                    stdin=f,
                    env=env,
                    check=True,
                )
            print("[OK] Database restored successfully.")

        extracted_uploads = tmp_dir / "uploads"
        if extracted_uploads.exists():
            target_uploads = BACKEND_DIR / "uploads"
            target_uploads.mkdir(parents=True, exist_ok=True)
            shutil.copytree(extracted_uploads, target_uploads, dirs_exist_ok=True)
            print(f"[OK] Uploads restored to {target_uploads}")


def cmd_build_all(_args: argparse.Namespace) -> None:
    npm_cmd = "npm.cmd" if os.name == "nt" else "npm"
    print("[1/2] Building frontend SPA...")
    subprocess.run([npm_cmd, "run", "build"], cwd=FRONTEND_DIR, check=True)

    DIST_RELEASE_DIR.mkdir(parents=True, exist_ok=True)
    targets = [
        ("linux", "amd64", "vilivala-linux-amd64"),
        ("linux", "arm64", "vilivala-linux-arm64"),
        ("darwin", "arm64", "vilivala-darwin-arm64"),
        ("darwin", "amd64", "vilivala-darwin-amd64"),
        ("windows", "amd64", "vilivala-windows-amd64.exe"),
    ]
    print("[2/2] Cross-compiling Go backend for Linux / macOS / Windows...")
    for goos, goarch, out_name in targets:
        env = os.environ.copy()
        env["CGO_ENABLED"] = "0"
        env["GOOS"] = goos
        env["GOARCH"] = goarch
        out_path = DIST_RELEASE_DIR / out_name
        subprocess.run(
            ["go", "build", "-ldflags=-s -w", "-o", str(out_path), "./cmd/server"],
            cwd=BACKEND_DIR,
            env=env,
            check=True,
        )
        print(f"  -> {out_path.relative_to(ROOT)}")

    fe_target = DIST_RELEASE_DIR / "frontend-dist"
    if fe_target.exists():
        shutil.rmtree(fe_target)
    shutil.copytree(FRONTEND_DIR / "dist", fe_target)
    print(f"[OK] Release artifacts ready in {DIST_RELEASE_DIR}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Vilivala cross-platform migration & build tool")
    sub = parser.add_subparsers(dest="command", required=True)
    sub.add_parser("pack", help="Export DB + uploads into a portable .tar.gz bundle")
    up = sub.add_parser("unpack", help="Restore DB + uploads from a portable .tar.gz bundle")
    up.add_argument("bundle", help="Path to vilivala-migration-*.tar.gz")
    sub.add_parser("build-all", help="Build frontend + Linux/macOS/Windows backend binaries")

    args = parser.parse_args()
    if args.command == "pack":
        cmd_pack(args)
    elif args.command == "unpack":
        cmd_unpack(args)
    elif args.command == "build-all":
        cmd_build_all(args)


if __name__ == "__main__":
    main()
