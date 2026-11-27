<template>
  <div class="card p-5 text-center">
    <img
      :src="siteConfig.blogger.avatar"
      :alt="siteConfig.blogger.name"
      class="w-20 h-20 rounded-full mx-auto object-cover ring-2 ring-[#42b883]/40 shadow-sm hover:rotate-6 transition-transform duration-300"
    />
    <h3 class="mt-3 text-base font-bold text-[#213547] dark:text-white">{{ siteConfig.blogger.name }}</h3>
    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ siteConfig.blogger.slogan }}</p>

    <div class="mt-4 flex items-center justify-center flex-wrap gap-2">
      <a
        v-for="s in siteConfig.socials"
        :key="s.title"
        :href="s.link"
        target="_blank"
        rel="noopener"
        :title="s.title"
        class="px-2.5 py-1 rounded-md text-xs bg-gray-100 dark:bg-white/10 text-gray-600 dark:text-gray-300 hover:bg-[#42b883] hover:text-white transition-colors"
      >
        {{ s.title }}
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { siteConfig } from '@/site.config'
</script>
