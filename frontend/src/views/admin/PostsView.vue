<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-bold text-[#213547] dark:text-white">文章管理（{{ total }}）</h1>
      <router-link to="/admin/posts/new" class="btn-primary">+ 新建文章</router-link>
    </div>

    <div class="flex flex-wrap gap-2 mb-4">
      <select v-model="status" class="input-field !w-28 sm:!w-32" @change="reload(1)">
        <option value="">全部状态</option>
        <option value="published">已发布</option>
        <option value="draft">草稿</option>
      </select>
      <input v-model="keyword" class="input-field flex-1 min-w-[160px] sm:!w-56 sm:flex-none" placeholder="搜索标题、摘要或正文…" @keyup.enter="reload(1)" />
      <button class="btn-secondary" @click="reload(1)">搜索</button>
    </div>

    <div class="card overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-white/5 text-gray-500 dark:text-gray-400 text-left">
          <tr>
            <th class="px-3.5 sm:px-4 py-3 font-medium">标题</th>
            <th class="hidden lg:table-cell px-4 py-3 font-medium">分类</th>
            <th class="px-3 sm:px-4 py-3 font-medium whitespace-nowrap">状态</th>
            <th class="hidden md:table-cell px-4 py-3 font-medium">字数 / 阅读</th>
            <th class="hidden sm:table-cell px-4 py-3 font-medium">创建时间</th>
            <th class="px-3.5 sm:px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in posts" :key="p.id" class="border-t border-gray-100 dark:border-white/10 hover:bg-gray-50/60 dark:hover:bg-white/[0.03]">
            <td class="px-3.5 sm:px-4 py-3 font-medium text-[#213547] dark:text-gray-100 max-w-[180px] sm:max-w-xs">
              <div class="truncate">
                <span v-if="p.is_top" class="badge bg-[#42b883]/15 text-[#42b883] mr-1.5">顶</span>
                {{ p.title }}
              </div>
              <div class="sm:hidden mt-1 text-[11px] font-normal text-gray-400 truncate">
                {{ p.category?.name || '未分类' }} · {{ formatDate(p.created_at) }}
              </div>
            </td>
            <td class="hidden lg:table-cell px-4 py-3 text-gray-500 dark:text-gray-400 text-xs">
              {{ p.category?.path?.map((n) => n.name).join(' / ') || '—' }}
            </td>
            <td class="px-3 sm:px-4 py-3 whitespace-nowrap">
              <span
                class="badge"
                :class="p.status === 'published' ? 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400' : 'bg-amber-500/15 text-amber-600 dark:text-amber-400'"
              >
                {{ p.status === 'published' ? '已发' : '草稿' }}
              </span>
            </td>
            <td class="hidden md:table-cell px-4 py-3 text-gray-500 dark:text-gray-400 text-xs whitespace-nowrap">
              {{ p.word_count }} 字 · {{ p.view_count }}
            </td>
            <td class="hidden sm:table-cell px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs whitespace-nowrap">{{ formatDate(p.created_at) }}</td>
            <td class="px-3.5 sm:px-4 py-3 text-right whitespace-nowrap space-x-2 sm:space-x-3 text-xs sm:text-sm">
              <router-link :to="`/posts/${p.slug}`" target="_blank" class="text-gray-500 dark:text-gray-400 hover:text-[#42b883]">
                预览
              </router-link>
              <button class="text-gray-600 dark:text-gray-300 hover:text-[#42b883]" @click="toggle(p)">
                {{ p.status === 'published' ? '转草稿' : '发布' }}
              </button>
              <router-link :to="`/admin/posts/${p.id}/edit`" class="text-[#42b883] hover:underline">编辑</router-link>
              <button class="text-red-500 hover:underline" @click="remove(p)">删除</button>
            </td>
          </tr>
          <tr v-if="!posts.length">
            <td colspan="6" class="px-4 py-12 text-center text-gray-400">暂无文章</td>
          </tr>
        </tbody>
      </table>
    </div>
    <Pagination :page="page" :page-size="pageSize" :total="total" @change="reload" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { adminDeletePost, adminGetPosts, adminUpdatePostStatus } from '@/api'
import type { Post } from '@/api/types'
import Pagination from '@/components/Pagination.vue'
import { formatDate } from '@/utils/format'
import { toast } from '@/utils/toast'

const posts = ref<Post[]>([])
const page = ref(1)
const pageSize = 15
const total = ref(0)
const status = ref('')
const keyword = ref('')

async function reload(p = page.value) {
  page.value = p
  try {
    const res = await adminGetPosts({
      page: p,
      page_size: pageSize,
      status: status.value || undefined,
      keyword: keyword.value.trim() || undefined,
    })
    posts.value = res.data.list
    total.value = res.data.total
  } catch {
    /* handled by request interceptor */
  }
}

async function toggle(p: Post) {
  const next = p.status === 'published' ? 'draft' : 'published'
  try {
    await adminUpdatePostStatus(p.id, next)
    toast.success(next === 'published' ? '已发布' : '已转为草稿')
    await reload()
  } catch {
    /* handled by interceptor */
  }
}

async function remove(p: Post) {
  if (!confirm(`确定删除《${p.title}》吗？此操作不可恢复。`)) return
  try {
    await adminDeletePost(p.id)
    toast.success('已删除')
    if (posts.value.length === 1 && page.value > 1) page.value--
    await reload()
  } catch {
    /* handled by interceptor */
  }
}

onMounted(() => reload(1))
</script>
