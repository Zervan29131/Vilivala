export interface NamingVariant {
  style: string
  label: string
  value: string
}

const CN_DICT: [string, string][] = [
  ['用户配置', 'user config'],
  ['站点设置', 'site settings'],
  ['文章列表', 'post list'],
  ['文章详情', 'post detail'],
  ['分类目录', 'category tree'],
  ['个人博客', 'personal blog'],
  ['阅读时间', 'reading time'],
  ['字数统计', 'word count'],
  ['创建时间', 'created at'],
  ['更新时间', 'updated at'],
  ['发布状态', 'publish status'],
  ['访问次数', 'view count'],
  ['阅读量', 'view count'],
  ['验证码', 'captcha code'],
  ['管理员', 'admin'],
  ['仪表盘', 'dashboard'],
  ['数据库', 'database'],
  ['中间件', 'middleware'],
  ['服务器', 'server'],
  ['客户端', 'client'],
  ['上下文', 'context'],
  ['初始化', 'initialize'],
  ['自定义', 'custom'],
  ['暗黑', 'dark'],
  ['浅色', 'light'],
  ['主题', 'theme'],
  ['用户', 'user'],
  ['账号', 'account'],
  ['密码', 'password'],
  ['令牌', 'token'],
  ['登录', 'login'],
  ['注册', 'register'],
  ['退出', 'logout'],
  ['权限', 'permission'],
  ['角色', 'role'],
  ['文章', 'post'],
  ['标题', 'title'],
  ['摘要', 'summary'],
  ['正文', 'content'],
  ['内容', 'content'],
  ['封面', 'cover'],
  ['分类', 'category'],
  ['标签', 'tag'],
  ['归档', 'archive'],
  ['页面', 'page'],
  ['评论', 'comment'],
  ['回复', 'reply'],
  ['通知', 'notification'],
  ['消息', 'message'],
  ['订单', 'order'],
  ['商品', 'product'],
  ['价格', 'price'],
  ['状态', 'status'],
  ['类型', 'type'],
  ['列表', 'list'],
  ['详情', 'detail'],
  ['总数', 'total'],
  ['数量', 'count'],
  ['分页', 'pagination'],
  ['当前', 'current'],
  ['历史', 'history'],
  ['配置', 'config'],
  ['设置', 'settings'],
  ['参数', 'params'],
  ['选项', 'options'],
  ['结果', 'result'],
  ['响应', 'response'],
  ['请求', 'request'],
  ['接口', 'api'],
  ['错误', 'error'],
  ['成功', 'success'],
  ['缓存', 'cache'],
  ['队列', 'queue'],
  ['任务', 'task'],
  ['定时', 'schedule'],
  ['书签', 'bookmark'],
  ['目录', 'directory'],
  ['路径', 'path'],
  ['链接', 'link'],
  ['地址', 'url'],
  ['文件', 'file'],
  ['图片', 'image'],
  ['头像', 'avatar'],
  ['上传', 'upload'],
  ['下载', 'download'],
  ['创建', 'create'],
  ['新增', 'add'],
  ['更新', 'update'],
  ['修改', 'edit'],
  ['删除', 'delete'],
  ['移除', 'remove'],
  ['查询', 'query'],
  ['搜索', 'search'],
  ['获取', 'get'],
  ['保存', 'save'],
  ['校验', 'validate'],
  ['名称', 'name'],
  ['编号', 'id'],
  ['时间', 'time'],
  ['日期', 'date'],
]

function extractWords(input: string): string[] {
  let working = input.trim()
  if (!working) return []

  // Replace known Chinese terms greedily (longest phrases are listed first)
  for (const [cn, en] of CN_DICT) {
    if (working.includes(cn)) {
      working = working.split(cn).join(` ${en} `)
    }
  }

  // Split camelCase / PascalCase boundaries and non-alphanumeric separators
  const spaced = working
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/([A-Z]+)([A-Z][a-z])/g, '$1 $2')
    .replace(/[^a-zA-Z0-9]+/g, ' ')
    .trim()

  return spaced
    .split(/\s+/)
    .map((w) => w.toLowerCase())
    .filter(Boolean)
}

export function buildNamingVariants(input: string): { words: string[]; variants: NamingVariant[] } {
  const words = extractWords(input)
  if (!words.length) return { words: [], variants: [] }

  const cap = (w: string) => w.charAt(0).toUpperCase() + w.slice(1)
  const camel = words[0] + words.slice(1).map(cap).join('')
  const pascal = words.map(cap).join('')
  const snake = words.join('_')
  const constant = words.join('_').toUpperCase()
  const kebab = words.join('-')

  return {
    words,
    variants: [
      { style: 'camelCase', label: '小驼峰（变量 / 函数）', value: camel },
      { style: 'PascalCase', label: '大驼峰（组件 / 类型）', value: pascal },
      { style: 'snake_case', label: '下划线（字段 / Go / Python）', value: snake },
      { style: 'CONSTANT_CASE', label: '全大写（常量 / 环境变量）', value: constant },
      { style: 'kebab-case', label: '短横线（路由 / CSS 类名）', value: kebab },
    ],
  }
}
