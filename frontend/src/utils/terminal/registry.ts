export interface CommandDoc {
  name: string
  aliases?: string[]
  category: 'space' | 'blog' | 'web' | 'tools' | 'system'
  usage: string
  summary: string
  options?: string[]
  examples?: string[]
}

export const COMMAND_DOCS: CommandDoc[] = [
  // 1. Space / Virtual Directory
  {
    name: 'ls',
    aliases: ['list', 'dir'],
    category: 'space',
    usage: 'ls [路径|posts|categories|pages] [-r]',
    summary: '列出当前目录内容（在分类目录下自动列出该类文章）',
    options: ['-r, --recursive  递归展示子目录'],
    examples: ['ls', 'ls -r', 'ls /code', 'ls posts'],
  },
  {
    name: 'cd',
    category: 'space',
    usage: 'cd [目录路径 | .. | ~]',
    summary: '进入博客分类目录或自定义书签文件夹',
    examples: ['cd code', 'cd /ToRead', 'cd ..', 'cd ~'],
  },
  {
    name: 'pwd',
    category: 'space',
    usage: 'pwd',
    summary: '显示当前所在的虚拟目录路径',
  },
  {
    name: 'mkdir',
    category: 'space',
    usage: 'mkdir <目录名>',
    summary: '在当前目录新建自定义书签文件夹',
    examples: ['mkdir AI工具', 'mkdir /dev/前端'],
  },
  {
    name: 'touch',
    aliases: ['add'],
    category: 'space',
    usage: 'touch <书签名> <网址>',
    summary: '在当前目录新增一个网站书签',
    examples: ['touch V2EX https://v2ex.com', 'touch /dev/Vite https://cn.vitejs.dev'],
  },
  {
    name: 'mv',
    aliases: ['move'],
    category: 'space',
    usage: 'mv <源书签> <目标路径>',
    summary: '移动或重命名自定义书签',
    examples: ['mv V2EX /dev/V2EX'],
  },
  {
    name: 'cp',
    aliases: ['copy'],
    category: 'space',
    usage: 'cp <源书签> <目标路径>',
    summary: '复制自定义书签到新位置',
  },
  {
    name: 'rm',
    aliases: ['del'],
    category: 'space',
    usage: 'rm <路径> [-r]',
    summary: '删除自定义书签或书签文件夹（博客分类只读）',
    options: ['-r, --recursive  递归删除非空书签文件夹'],
    examples: ['rm V2EX', 'rm -r AI工具'],
  },

  // 2. Blog Content
  {
    name: 's',
    aliases: ['find'],
    category: 'blog',
    usage: 's <关键词>',
    summary: '站内全文搜索博客文章（支持中英文）',
    examples: ['s MySQL', 's 博客'],
  },
  {
    name: 'hot',
    category: 'blog',
    usage: 'hot',
    summary: '弹出本站阅读量 Top 10 热门文章排行榜卡片',
  },
  {
    name: 'cat',
    category: 'blog',
    usage: 'cat <文章slug>',
    summary: '在终端内预览指定文章摘要与信息',
    examples: ['cat c6d7d3'],
  },
  {
    name: 'open',
    category: 'blog',
    usage: 'open <文章slug>',
    summary: '直接跳转阅读指定文章',
    examples: ['open c6d7d3'],
  },
  {
    name: 'whoami',
    aliases: ['about', 'info'],
    category: 'blog',
    usage: 'whoami',
    summary: '查看博主名片与全站实时统计',
  },

  // 3. Web Search & Navigation
  {
    name: 'search',
    category: 'web',
    usage: 'search <关键词> [-f 搜索源] [-s]',
    summary: '多引擎网页检索（默认 Google，可切知乎/必应/豆瓣/MDN 等）',
    options: [
      '-f, --from <源>  可选 google|baidu|bing|github|bili|zhihu|douban|mdn|blog',
      '-s, --self       在当前标签页直接打开',
    ],
    examples: ['search Vue3 -f mdn', 'search 机器学习 -f zhihu', 'search Go并发 -f bing -s'],
  },
  {
    name: 'google',
    aliases: ['g'],
    category: 'web',
    usage: 'g <关键词> [-s]',
    summary: 'Google 搜索',
  },
  {
    name: 'baidu',
    aliases: ['bd'],
    category: 'web',
    usage: 'bd <关键词> [-s]',
    summary: '百度搜索',
  },
  {
    name: 'bing',
    category: 'web',
    usage: 'bing <关键词> [-s]',
    summary: '必应 Bing 搜索',
  },
  {
    name: 'github',
    aliases: ['gh'],
    category: 'web',
    usage: 'gh <关键词> [-s]',
    summary: 'GitHub 开源仓库搜索',
  },
  {
    name: 'bilibili',
    aliases: ['bili'],
    category: 'web',
    usage: 'bili <关键词> [-s]',
    summary: 'Bilibili 视频搜索',
  },
  {
    name: 'zhihu',
    category: 'web',
    usage: 'zhihu <关键词> [-s]',
    summary: '知乎问答/文章搜索',
  },
  {
    name: 'douban',
    category: 'web',
    usage: 'douban <关键词> [-s]',
    summary: '豆瓣书影音搜索',
  },
  {
    name: 'mdn',
    category: 'web',
    usage: 'mdn <关键词> [-s]',
    summary: 'MDN Web 开发者文档搜索',
  },
  {
    name: 'goto',
    aliases: ['visit'],
    category: 'web',
    usage: 'goto <书签名或网址> [-s]',
    summary: '直达已保存书签或任意网址',
    examples: ['goto GitHub', 'goto vuejs.org -s'],
  },
  {
    name: 'bm',
    aliases: ['bookmarks'],
    category: 'web',
    usage: 'bm [add <名> <网址> | rm <名>]',
    summary: '查看或快捷增删书签',
  },

  // 4. Developer Micro-tools
  {
    name: 'timing',
    aliases: ['timer', 'pomodoro'],
    category: 'tools',
    usage: 'timing <时长，如 25m / 90s / 60> [标题]',
    summary: '启动行内交互式番茄钟 / 倒计时卡片（支持暂停与重置）',
    examples: ['timing 25m 专注编码', 'timing 60s'],
  },
  {
    name: 'varbook',
    aliases: ['var'],
    category: 'tools',
    usage: 'varbook <中文词汇或英文词组>',
    summary: '变量命名转换器（输出小驼峰/大驼峰/下划线等 5 种写法，点按即复制）',
    examples: ['varbook 用户配置', 'varbook 文章列表', 'varbook post detail'],
  },
  {
    name: 'todo',
    category: 'tools',
    usage: 'todo [add <事项> | done <序号> | rm <序号>]',
    summary: '轻量本地待办事项管理',
    examples: ['todo add 阅读完 Go 并发章节', 'todo done 1'],
  },
  {
    name: 'shortcut',
    aliases: ['keys'],
    category: 'tools',
    usage: 'shortcut',
    summary: '弹出终端快捷键一览卡片',
  },

  // 5. System & Appearance
  {
    name: 'bg',
    aliases: ['background'],
    category: 'system',
    usage: 'bg [grid|emerald|midnight|slate|random|<图片URL>]',
    summary: '切换终端背景预设、随机壁纸或自定义图片网址',
    examples: ['bg emerald', 'bg random', 'bg /img/net.png'],
  },
  {
    name: 'theme',
    category: 'system',
    usage: 'theme [system|light|dark]',
    summary: '切换全站明暗主题',
  },
  {
    name: 'hint',
    category: 'system',
    usage: 'hint [on|off]',
    summary: '开启或关闭输入框下方的实时命令语法提示',
  },
  {
    name: 'reset',
    category: 'system',
    usage: 'reset [--all]',
    summary: '重置壁纸、书签目录与提示开关为初始默认状态（加 --all 同时清空待办）',
  },
  {
    name: 'help',
    aliases: ['man'],
    category: 'system',
    usage: 'help [命令名]',
    summary: '查看全部命令总览，或查询某条命令的详细参数与示例',
    examples: ['help', 'help search', 'help cd', 'help timing'],
  },
  {
    name: 'history',
    aliases: ['h'],
    category: 'system',
    usage: 'history',
    summary: '查看最近执行的命令记录（支持 !! 或 !序号 重跑）',
  },
  {
    name: 'date',
    category: 'system',
    usage: 'date',
    summary: '显示当前系统日期与时间',
  },
  {
    name: 'clear',
    aliases: ['cl'],
    category: 'system',
    usage: 'clear',
    summary: '清空终端屏幕（快捷键 Ctrl + L）',
  },
  {
    name: 'home',
    aliases: ['exit'],
    category: 'system',
    usage: 'home',
    summary: '返回博客主站首页',
  },
]

export function findCommandDoc(nameOrAlias: string): CommandDoc | undefined {
  const q = nameOrAlias.trim().toLowerCase()
  if (!q) return undefined
  return COMMAND_DOCS.find((d) => d.name === q || d.aliases?.includes(q))
}

export function allCommandKeywords(): string[] {
  const set = new Set<string>()
  for (const d of COMMAND_DOCS) {
    set.add(d.name)
    d.aliases?.forEach((a) => set.add(a))
  }
  return Array.from(set)
}
