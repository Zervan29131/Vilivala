<template>
  <div class="max-w-4xl space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <div>
        <h1 class="text-xl font-bold text-[#213547] dark:text-white">站点设置</h1>
        <p class="text-xs text-gray-400 mt-1">动态调整前台文案、公告、打赏/版权开关与备案信息，保存后立即生效。</p>
      </div>
      <button type="button" class="btn-primary" :disabled="saving" @click="save">
        {{ saving ? '保存中…' : '保存全部设置' }}
      </button>
    </div>

    <!-- 1. Branding & Hero -->
    <section class="card p-5 space-y-4">
      <h2 class="text-sm font-bold text-[#213547] dark:text-white">一、站点与首页 Hero</h2>
      <div class="grid sm:grid-cols-2 gap-4">
        <div>
          <label class="lbl">站点名称</label>
          <input v-model="form.site_title" class="input-field" placeholder="Zervan的小站" />
        </div>
        <div>
          <label class="lbl">站点简述</label>
          <input v-model="form.site_description" class="input-field" placeholder="Zervanのblog" />
        </div>
      </div>
      <div class="grid sm:grid-cols-2 gap-4">
        <div>
          <label class="lbl">首页 Hero 主标题</label>
          <input v-model="form.hero_title" class="input-field" />
        </div>
        <div>
          <label class="lbl">默认副标题</label>
          <input v-model="form.hero_tagline" class="input-field" />
        </div>
      </div>
      <div>
        <label class="lbl">打字机轮播副标题（每行一句，仅留一行则静止显示）</label>
        <textarea v-model="subtitlesText" rows="3" class="input-field font-mono text-xs" />
      </div>
    </section>

    <!-- 2. Blogger Card & Socials -->
    <section class="card p-5 space-y-4">
      <h2 class="text-sm font-bold text-[#213547] dark:text-white">二、博主名片与社交外链</h2>
      <div class="grid sm:grid-cols-3 gap-4">
        <div>
          <label class="lbl">博主名称</label>
          <input v-model="form.blogger_name" class="input-field" />
        </div>
        <div>
          <label class="lbl">个性签名</label>
          <input v-model="form.blogger_slogan" class="input-field" />
        </div>
        <div>
          <label class="lbl">头像地址</label>
          <div class="flex gap-2">
            <input v-model="form.blogger_avatar" class="input-field font-mono text-xs" placeholder="/img/console.jpg" />
            <label class="btn-secondary !px-2.5 cursor-pointer shrink-0 text-xs">
              上传
              <input type="file" accept="image/*" class="hidden" @change="uploadAvatar" />
            </label>
          </div>
        </div>
      </div>

      <div>
        <div class="flex items-center justify-between mb-2">
          <label class="lbl !mb-0">社交外链列表</label>
          <button type="button" class="text-xs text-[#42b883] hover:underline" @click="addSocial">
            + 添加外链
          </button>
        </div>
        <div class="space-y-2">
          <div
            v-for="(s, idx) in form.socials"
            :key="idx"
            class="grid grid-cols-[130px_1fr_auto] gap-2 items-center"
          >
            <input v-model="s.title" class="input-field" placeholder="名称（如 知乎）" />
            <input v-model="s.link" class="input-field font-mono text-xs" placeholder="https://... 或 mailto:..." />
            <button type="button" class="text-xs text-red-500 hover:underline px-2" @click="form.socials.splice(idx, 1)">
              删除
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- 3. Broadcast, Reward & Copyright -->
    <section class="card p-5 space-y-4">
      <h2 class="text-sm font-bold text-[#213547] dark:text-white">三、公告栏、赞赏与版权声明</h2>

      <div class="flex flex-wrap items-center gap-6 text-sm">
        <label class="inline-flex items-center gap-2 cursor-pointer">
          <input v-model="form.notice_enabled" type="checkbox" class="rounded accent-[#42b883]" />
          <span>启用侧边栏公告</span>
        </label>
        <label class="inline-flex items-center gap-2 cursor-pointer">
          <input v-model="form.reward_enabled" type="checkbox" class="rounded accent-[#42b883]" />
          <span>启用文末赞赏按钮</span>
        </label>
        <label class="inline-flex items-center gap-2 cursor-pointer">
          <input v-model="form.copyright_enabled" type="checkbox" class="rounded accent-[#42b883]" />
          <span>启用文末版权声明框</span>
        </label>
      </div>

      <div class="grid sm:grid-cols-[180px_1fr] gap-4">
        <div>
          <label class="lbl">公告模式</label>
          <select v-model="form.notice_mode" class="input-field">
            <option value="custom">固定自定义公告</option>
            <option value="hitokoto">每日一言（自动拉取）</option>
          </select>
        </div>
        <div>
          <label class="lbl">公告文案（一言模式下作为离线兜底）</label>
          <input v-model="form.notice_text" class="input-field" />
        </div>
      </div>

      <div class="grid sm:grid-cols-2 gap-4">
        <div>
          <label class="lbl">赞赏弹窗提示语</label>
          <input v-model="form.reward_wording" class="input-field" />
        </div>
        <div>
          <label class="lbl">原创版权声明文案</label>
          <input v-model="form.copyright_notice" class="input-field" />
        </div>
      </div>
    </section>

    <!-- 4. Footer & Filing -->
    <section class="card p-5 space-y-4">
      <h2 class="text-sm font-bold text-[#213547] dark:text-white">四、页脚年份与备案信息（留空则不显示备案行）</h2>
      <div class="grid sm:grid-cols-3 gap-4">
        <div>
          <label class="lbl">建站起始年份</label>
          <input v-model.number="form.since_year" type="number" class="input-field" placeholder="2021" />
        </div>
        <div>
          <label class="lbl">ICP 备案号</label>
          <input v-model="form.icp_number" class="input-field" placeholder="例如：鲁ICP备2021036977号-1" />
        </div>
        <div>
          <label class="lbl">ICP 跳转链接</label>
          <input v-model="form.icp_url" class="input-field font-mono text-xs" placeholder="https://beian.miit.gov.cn/" />
        </div>
      </div>
      <div class="grid sm:grid-cols-2 gap-4">
        <div>
          <label class="lbl">公安备案号</label>
          <input v-model="form.gongan_number" class="input-field" placeholder="可选" />
        </div>
        <div>
          <label class="lbl">公安备案链接</label>
          <input v-model="form.gongan_url" class="input-field font-mono text-xs" placeholder="https://www.beian.gov.cn/" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { adminUpdateSiteConfig, adminUploadImage, getSiteConfig } from '@/api'
import type { SiteConfigData } from '@/api/types'
import { defaultSiteConfig, useSiteStore } from '@/stores/site'
import { toast } from '@/utils/toast'

const siteStore = useSiteStore()
const saving = ref(false)
const form = reactive<SiteConfigData>(defaultSiteConfig())
const subtitlesText = ref('')

function syncFrom(data: SiteConfigData) {
  Object.assign(form, {
    ...data,
    socials: (data.socials || []).map((s) => ({ ...s })),
  })
  subtitlesText.value = (data.subtitles || []).join('\n')
}

function addSocial() {
  form.socials.push({ title: '', link: '' })
}

async function uploadAvatar(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const res = await adminUploadImage(file)
    form.blogger_avatar = res.data.url
    toast.success('头像上传成功')
  } catch {
    /* handled by interceptor */
  } finally {
    input.value = ''
  }
}

async function save() {
  saving.value = true
  try {
    const payload: Omit<SiteConfigData, 'stats'> = {
      site_title: form.site_title,
      site_description: form.site_description,
      hero_title: form.hero_title,
      hero_tagline: form.hero_tagline,
      subtitles: subtitlesText.value
        .split(/\r?\n/)
        .map((s) => s.trim())
        .filter(Boolean),
      blogger_name: form.blogger_name,
      blogger_slogan: form.blogger_slogan,
      blogger_avatar: form.blogger_avatar,
      notice_enabled: form.notice_enabled,
      notice_mode: form.notice_mode,
      notice_text: form.notice_text,
      reward_enabled: form.reward_enabled,
      reward_wording: form.reward_wording,
      copyright_enabled: form.copyright_enabled,
      copyright_notice: form.copyright_notice,
      since_year: Number(form.since_year) || 2021,
      icp_number: form.icp_number,
      icp_url: form.icp_url,
      gongan_number: form.gongan_number,
      gongan_url: form.gongan_url,
      socials: form.socials,
    }
    const res = await adminUpdateSiteConfig(payload)
    siteStore.applyData(res.data)
    syncFrom(res.data)
    toast.success('站点设置已保存')
  } catch {
    /* handled by interceptor */
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const res = await getSiteConfig()
    siteStore.applyData(res.data)
    syncFrom(res.data)
  } catch {
    syncFrom(siteStore.config)
  }
})
</script>

<style scoped>
.lbl {
  @apply block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1;
}
</style>
