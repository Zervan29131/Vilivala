// Site-wide branding & hero configuration (ported from legacy Vilivala)

export interface FeatureItem {
  title: string
  details: string
  slug: string
  imgUrl: string
}

export interface SocialItem {
  title: string
  link: string
  icon: 'weibo' | 'bilibili' | 'zhihu' | 'douyin' | 'mail' | 'github'
}

export const siteConfig = {
  title: 'Zervan的小站',
  description: 'Zervanのblog',
  logo: '/img/logo.png',
  hero: {
    image: '/img/net.png',
    title: '这是存在于某处的幻想世界',
    tagline: 'Be the change you want to see in the World',
    actionText: '支持一下',
    actionLink: '/pages/sponsor',
  },
  features: [
    { title: 'code', details: '技术', slug: 'code', imgUrl: '/img/code.png' },
    { title: 'ToRead', details: '思想', slug: 'toread', imgUrl: '/img/ToRead.png' },
    { title: 'fun', details: '未来', slug: 'fun', imgUrl: '/img/fun.png' },
  ] as FeatureItem[],
  blogger: {
    avatar: '/img/console.jpg',
    name: 'Zervan',
    slogan: '天下最普通的人之一',
  },
  pandaWidget: '/img/panda-waving.png',
  githubUrl: 'https://github.com/Zervan29131',
  socials: [
    { title: '微博', link: 'https://weibo.com/zervan', icon: 'weibo' },
    { title: 'bilibili', link: 'https://space.bilibili.com/1650366974', icon: 'bilibili' },
    { title: '知乎', link: 'https://www.zhihu.com/people/zervan', icon: 'zhihu' },
    { title: '抖音', link: 'https://v.douyin.com/L7jT1xp/', icon: 'douyin' },
    { title: '发邮件', link: 'mailto:nickmechlin01@outlook.com', icon: 'mail' },
  ] as SocialItem[],
}
