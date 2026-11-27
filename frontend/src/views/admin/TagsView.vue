<template>
  <div>
    <h1 class="text-xl font-bold text-[#213547] dark:text-white mb-6">标签管理（{{ items.length }}）</h1>

    <form class="card p-4 mb-6 flex flex-wrap gap-2 items-center" @submit.prevent="save">
      <input v-model="form.name" class="input-field !w-48" placeholder="标签名称" />
      <input v-model="form.slug" class="input-field !w-48 font-mono" placeholder="slug（留空自动生成）" />
      <button class="btn-primary" :disabled="!form.name.trim()">{{ editingId ? '保存修改' : '+ 新增' }}</button>
      <button v-if="editingId" type="button" class="btn-secondary" @click="reset">取消</button>
    </form>

    <div class="card overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-white/5 text-gray-500 dark:text-gray-400 text-left">
          <tr>
            <th class="px-4 py-3 font-medium">名称</th>
            <th class="px-4 py-3 font-medium">Slug</th>
            <th class="px-4 py-3 font-medium">文章数</th>
            <th class="px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="i in items" :key="i.id" class="border-t border-gray-100 dark:border-white/10">
            <td class="px-4 py-3 font-medium text-[#213547] dark:text-gray-100">{{ i.name }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">{{ i.slug }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400">{{ i.post_count }}</td>
            <td class="px-4 py-3 text-right space-x-3">
              <button class="text-[#42b883] hover:underline" @click="edit(i)">编辑</button>
              <button class="text-red-500 hover:underline" @click="remove(i)">删除</button>
            </td>
          </tr>
          <tr v-if="!items.length">
            <td colspan="4" class="px-4 py-12 text-center text-gray-400">暂无标签</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { adminCreateTag, adminDeleteTag, adminUpdateTag, getTags } from '@/api'
import type { Tag } from '@/api/types'
import { makeSlug } from '@/utils/format'
import { toast } from '@/utils/toast'

const items = ref<Tag[]>([])
const editingId = ref<number | null>(null)
const form = reactive({ name: '', slug: '' })

async function load() {
  try {
    items.value = (await getTags()).data
  } catch {
    /* handled by interceptor */
  }
}
function reset() {
  editingId.value = null
  form.name = ''
  form.slug = ''
}
function edit(i: Tag) {
  editingId.value = i.id
  form.name = i.name
  form.slug = i.slug
}
async function save() {
  const payload = { name: form.name.trim(), slug: form.slug.trim() || makeSlug(form.name, 'tag') }
  if (editingId.value) await adminUpdateTag(editingId.value, payload)
  else await adminCreateTag(payload)
  toast.success('保存成功')
  reset()
  await load()
}
async function remove(i: Tag) {
  if (!confirm(`确定删除标签「${i.name}」吗？`)) return
  await adminDeleteTag(i.id)
  toast.success('已删除')
  await load()
}
onMounted(load)
</script>
