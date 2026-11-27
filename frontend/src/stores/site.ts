import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'
import { getSiteConfig } from '@/api'
import type { SiteConfigData } from '@/api/types'
import { siteConfig } from '@/site.config'

export function defaultSiteConfig(): SiteConfigData {
  return {
    site_title: siteConfig.title,
    site_description: siteConfig.description,
    hero_title: siteConfig.hero.title,
    hero_tagline: siteConfig.hero.tagline,
    subtitles: [
      siteConfig.hero.tagline,
      '记录代码、阅读与思考的角落',
      'Stay hungry, stay foolish',
    ],
    blogger_name: siteConfig.blogger.name,
    blogger_slogan: siteConfig.blogger.slogan,
    blogger_avatar: siteConfig.blogger.avatar,
    notice_enabled: true,
    notice_mode: 'custom',
    notice_text: 'Stay hungry, stay foolish · 欢迎来到 Zervan 的小站',
    reward_enabled: true,
    reward_wording: '觉得文章有帮助？请我喝杯咖啡吧 ~',
    copyright_enabled: true,
    copyright_notice: '本博客所有文章除特别声明外，著作权归作者所有。转载请注明出处！',
    since_year: 2021,
    icp_number: '',
    icp_url: 'https://beian.miit.gov.cn/',
    gongan_number: '',
    gongan_url: 'https://www.beian.gov.cn/',
    socials: siteConfig.socials.map((s) => ({ title: s.title, link: s.link })),
    stats: {
      published_posts: 0,
      total_words: 0,
      total_views: 0,
    },
  }
}

export const useSiteStore = defineStore('site', () => {
  const config = reactive<SiteConfigData>(defaultSiteConfig())
  const loaded = ref(false)
  let inflight: Promise<void> | null = null

  function applyData(data: SiteConfigData) {
    Object.assign(config, data)
    if (!Array.isArray(data.subtitles) || !data.subtitles.length) {
      config.subtitles = [data.hero_tagline || siteConfig.hero.tagline]
    }
    if (!Array.isArray(data.socials)) {
      config.socials = []
    }
  }

  async function ensureLoaded(force = false) {
    if (loaded.value && !force) return
    if (inflight && !force) return inflight
    inflight = (async () => {
      try {
        const res = await getSiteConfig()
        applyData(res.data)
        loaded.value = true
      } catch {
        /* keep local defaults gracefully */
      } finally {
        inflight = null
      }
    })()
    return inflight
  }

  return { config, loaded, applyData, ensureLoaded }
})
