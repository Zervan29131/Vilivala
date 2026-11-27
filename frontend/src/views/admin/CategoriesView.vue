<template>
  <div>
    <h1 class="text-xl font-bold text-[#213547] dark:text-white mb-6">分类管理（最多三级）</h1>

    <form class="card p-4 mb-6 grid sm:grid-cols-[1fr_1fr_1fr_90px_auto] gap-2 items-center" @submit.prevent="save">
      <input v-model="form.name" class="input-field" placeholder="分类名称" />
      <input v-model="form.slug" class="input-field font-mono" placeholder="slug（留空自动生成）" />
      <select v-model="form.parent_id" class="input-field">
        <option :value="0">顶级分类（无父级）</option>
        <option v-for="p in parentOptions" :key="p.id" :value="p.id">{{ p.path }}</option>
      </select>
      <input v-model.number="form.sort" type="number" class="input-field" placeholder="排序" title="同级排序（小在前）" />
      <div class="flex gap-2">
        <button class="btn-primary shrink-0" :disabled="!form.name.trim()">{{ editingId ? '保存' : '+ 新增' }}</button>
        <button v-if="editingId" type="button" class="btn-secondary shrink-0" @click="reset">取消</button>
      </div>
    </form>

    <div class="card overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-white/5 text-gray-500 dark:text-gray-400 text-left">
          <tr>
            <th class="px-4 py-3 font-medium">分类层级</th>
            <th class="px-4 py-3 font-medium">Slug</th>
            <th class="px-4 py-3 font-medium">排序</th>
            <th class="px-4 py-3 font-medium">文章（直属 / 含子类）</th>
            <th class="px-4 py-3 font-medium text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in flatTree" :key="c.id" class="border-t border-gray-100 dark:border-white/10">
            <td class="px-4 py-3 font-medium text-[#213547] dark:text-gray-100">
              <span :style="{ paddingLeft: `${(c.depth - 1) * 20}px` }" class="inline-flex items-center gap-1.5">
                <span v-if="c.depth > 1" class="text-gray-400">└</span>
                <span>{{ c.name }}</span>
                <span class="text-[11px] text-gray-400 font-normal">L{{ c.depth }}</span>
              </span>
            </td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">{{ c.slug }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 text-xs">{{ c.sort }}</td>
            <td class="px-4 py-3 text-gray-500 dark:text-gray-400 text-xs">{{ c.post_count }} / {{ c.total_count }}</td>
            <td class="px-4 py-3 text-right space-x-3 whitespace-nowrap">
              <button class="text-[#42b883] hover:underline" @click="edit(c)">编辑</button>
              <button class="text-red-500 hover:underline" @click="remove(c)">删除</button>
            </td>
          </tr>
          <tr v-if="!flatTree.length">
            <td colspan="5" class="px-4 py-12 text-center text-gray-400">暂无分类</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { adminCreateCategory, adminDeleteCategory, adminUpdateCategory, getCategories } from '@/api'
import type { Category } from '@/api/types'
import { makeSlug } from '@/utils/format'
import { toast } from '@/utils/toast'

const tree = ref<Category[]>([])
const editingId = ref<number | null>(null)
const form = reactive({ name: '', slug: '', parent_id: 0, sort: 0 })

function flatten(nodes: Category[], out: Category[] = []): Category[] {
  for (const n of nodes) {
    out.push(n)
    if (n.children?.length) flatten(n.children, out)
  }
  return out
}

const flatTree = computed(() => flatten(tree.value))

// Only depth 1 or 2 can be a parent (max depth = 3), and exclude self when editing
const parentOptions = computed(() =>
  flatTree.value.filter((c) => c.depth < 3 && c.id !== editingId.value)
)

async function load() {
  try {
    tree.value = (await getCategories(true)).data
  } catch {
    /* handled by interceptor */
  }
}

function reset() {
  editingId.value = null
  form.name = ''
  form.slug = ''
  form.parent_id = 0
  form.sort = 0
}

function edit(c: Category) {
  editingId.value = c.id
  form.name = c.name
  form.slug = c.slug
  form.parent_id = c.parent_id
  form.sort = c.sort
}

async function save() {
  const payload = {
    name: form.name.trim(),
    slug: form.slug.trim() || makeSlug(form.name, 'cat'),
    parent_id: form.parent_id,
    sort: form.sort || 0,
  }
  if (editingId.value) await adminUpdateCategory(editingId.value, payload)
  else await adminCreateCategory(payload)
  toast.success('保存成功')
  reset()
  await load()
}

async function remove(c: Category) {
  if (!confirm(`确定删除分类「${c.path}」吗？`)) return
  await adminDeleteCategory(c.id)
  toast.success('已删除')
  await load()
}

onMounted(load)
</script>
