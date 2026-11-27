<template>
  <div class="mt-6 flex flex-wrap items-center justify-center gap-3">
    <!-- Reward Button -->
    <button
      v-if="site.config.reward_enabled"
      type="button"
      class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-amber-500/15 hover:bg-amber-500 text-amber-600 dark:text-amber-400 hover:text-white text-xs sm:text-sm font-medium transition-colors"
      @click="rewardOpen = true"
    >
      <span>☕</span>
      <span>赞赏支持</span>
    </button>

    <!-- Share / QR Button -->
    <button
      type="button"
      class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-[#42b883]/12 hover:bg-[#42b883] text-[#42b883] dark:text-[#4fc08d] hover:text-white text-xs sm:text-sm font-medium transition-colors"
      @click="openShare"
    >
      <span>🔗</span>
      <span>分享 / 扫码阅读</span>
    </button>

    <!-- Reward Modal -->
    <teleport to="body">
      <div
        v-if="rewardOpen"
        class="fixed inset-0 z-[65] bg-black/55 backdrop-blur-sm flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        aria-label="赞赏支持"
        @click.self="rewardOpen = false"
      >
        <div class="card max-w-md w-full p-6 text-center space-y-4 shadow-2xl">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#213547] dark:text-white">☕ 赞赏支持</h3>
            <button type="button" class="text-gray-400 hover:text-gray-600 dark:hover:text-white" @click="rewardOpen = false">✕</button>
          </div>
          <p class="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
            {{ site.config.reward_wording }}
          </p>
          <div class="grid grid-cols-2 gap-4 pt-1">
            <div class="p-3 rounded-xl bg-gray-50 dark:bg-white/5 border border-gray-200/70 dark:border-white/10 space-y-2">
              <img src="/img/weichat.png" alt="微信赞赏码" class="w-36 h-36 mx-auto object-contain rounded" />
              <div class="text-xs font-medium text-[#213547] dark:text-gray-200">微信扫一扫</div>
            </div>
            <div class="p-3 rounded-xl bg-gray-50 dark:bg-white/5 border border-gray-200/70 dark:border-white/10 space-y-2">
              <img src="/img/alipay.png" alt="支付宝赞赏码" class="w-36 h-36 mx-auto object-contain rounded" />
              <div class="text-xs font-medium text-[#213547] dark:text-gray-200">支付宝扫一扫</div>
            </div>
          </div>
        </div>
      </div>
    </teleport>

    <!-- Share & QR Modal -->
    <teleport to="body">
      <div
        v-if="shareOpen"
        class="fixed inset-0 z-[65] bg-black/55 backdrop-blur-sm flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        aria-label="分享文章"
        @click.self="shareOpen = false"
      >
        <div class="card max-w-sm w-full p-6 text-center space-y-4 shadow-2xl">
          <div class="flex items-center justify-between">
            <h3 class="text-base font-bold text-[#213547] dark:text-white">分享文章</h3>
            <button type="button" class="text-gray-400 hover:text-gray-600 dark:hover:text-white" @click="shareOpen = false">✕</button>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ title }}</p>

          <!-- Locally generated SVG QR code (zero third-party tracking) -->
          <div
            v-if="qrSvg"
            class="w-40 h-40 mx-auto p-2 bg-white rounded-xl border border-gray-200 shadow-sm flex items-center justify-center"
            v-html="qrSvg"
          />
          <p class="text-[11px] text-gray-400">手机微信 / 浏览器扫一扫阅读</p>

          <div class="flex items-center justify-center gap-2 pt-1">
            <button type="button" class="btn-primary !text-xs !px-3.5 !py-1.5" @click="copyShareLink">
              复制链接
            </button>
            <a
              :href="weiboHref"
              target="_blank"
              rel="noopener noreferrer"
              class="btn-secondary !text-xs !px-3.5 !py-1.5"
            >
              分享到微博
            </a>
          </div>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import QRCode from 'qrcode'
import { useSiteStore } from '@/stores/site'
import { toast } from '@/utils/toast'

const props = defineProps<{
  title: string
  slug: string
}>()

const site = useSiteStore()
const rewardOpen = ref(false)
const shareOpen = ref(false)
const qrSvg = ref('')

const fullUrl = computed(() => `${window.location.origin}/posts/${props.slug}`)
const weiboHref = computed(
  () =>
    `https://service.weibo.com/share/share.php?url=${encodeURIComponent(fullUrl.value)}&title=${encodeURIComponent(props.title)}`
)

async function openShare() {
  shareOpen.value = true
  try {
    qrSvg.value = await QRCode.toString(fullUrl.value, {
      type: 'svg',
      margin: 1,
      width: 144,
      color: { dark: '#213547', light: '#ffffff' },
    })
  } catch {
    qrSvg.value = ''
  }
}

async function copyShareLink() {
  try {
    await navigator.clipboard.writeText(`${props.title} ${fullUrl.value}`)
    toast.success('分享链接已复制')
  } catch {
    toast.error('复制失败')
  }
}
</script>
