<template>
  <BaseModal 
    :show="store.showTagManageModal" 
    :closeOnOutside="true"
    @close="closeModal"
  >
    <div class="px-5 pt-5 pb-3 flex items-center justify-between border-b border-gray-100">
      <h3 class="text-base font-semibold text-gray-800">管理标签</h3>
      <button @click="closeModal" class="text-gray-400 hover:text-gray-600 text-lg cursor-pointer">&times;</button>
    </div>
    <div class="p-0 max-h-[70vh] overflow-y-auto">
      <div v-if="store.tags && store.tags.length > 0" class="divide-y divide-gray-100">
        <div v-for="tag in store.tags" :key="tag.name" class="px-5 py-3 hover:bg-gray-50 transition-colors group">
          
          <div v-if="editingTag !== tag.name" class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <span class="w-2.5 h-2.5 rounded-full flex-shrink-0" :style="{ backgroundColor: hexColors[tag.color] || hexColors.blue }"></span>
              <span class="text-sm text-gray-700">{{ tag.name }}</span>
            </div>
            <div class="flex items-center space-x-1">
              <button @click="startEdit(tag)" class="text-[11px] text-gray-400 hover:text-brand transition px-2 py-1 rounded">
                编辑
              </button>
              <button @click="deleteTag(tag.name)" :disabled="isProcessing" class="text-[11px] text-gray-400 hover:text-red-500 transition px-2 py-1 rounded">
                删除
              </button>
            </div>
          </div>

          <div v-else class="space-y-3">
            <div>
              <label class="text-xs text-gray-500 block mb-1">标签名称</label>
              <input 
                v-model="editForm.name" 
                type="text" 
                class="w-full text-sm border border-gray-300 rounded-md px-2.5 py-1.5 focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
              />
            </div>
            <div>
              <label class="text-xs text-gray-500 block mb-1">标签颜色</label>
              <div class="flex flex-wrap gap-2 pt-0.5">
                <button v-for="(hex, c) in hexColors" :key="c" @click="editForm.color = c" :class="['w-3.5 h-3.5 rounded-full cursor-pointer transition flex-shrink-0', editForm.color === c ? 'ring-2 ring-offset-1 ring-blue-400 scale-110' : 'hover:scale-110']" :style="{ backgroundColor: hex }"></button>
              </div>
            </div>
            <div class="flex justify-end space-x-2 pt-1">
              <button @click="editingTag = null" class="px-3 py-1 border border-gray-300 text-gray-600 rounded-md text-xs hover:bg-gray-100 transition cursor-pointer">
                取消
              </button>
              <button @click="saveTag(tag.name)" :disabled="isProcessing || !editForm.name.trim()" class="px-3 py-1 bg-brand text-white rounded-md text-xs font-medium hover:bg-brand-hover transition cursor-pointer disabled:opacity-50">
                保存
              </button>
            </div>
          </div>

        </div>
      </div>
      <div v-else class="text-center py-8 text-gray-400 text-sm">
        暂无标签数据
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseModal from '../common/BaseModal.vue'

const store = useArchiveStore()
const apiFetch = store.apiFetch

const hexColors = {
  blue: '#3b82f6',
  red: '#ef4444',
  green: '#10b981',
  yellow: '#eab308',
  purple: '#8b5cf6',
  pink: '#ec4899',
  orange: '#f97316',
  teal: '#14b8a6',
  cyan: '#06b6d4',
  gray: '#6b7280'
}

const editingTag = ref(null)
const editForm = ref({ name: '', color: 'blue' })
const isProcessing = ref(false)

const closeModal = () => {
  editingTag.value = null
  store.showTagManageModal = false
}

const startEdit = (tag) => {
  editingTag.value = tag.name
  editForm.value = { name: tag.name, color: tag.color }
}

const saveTag = async (oldName) => {
  if (!editForm.value.name.trim()) return
  isProcessing.value = true
  
  try {
    const res = await apiFetch('/api/tags', {
      method: 'PUT',
      body: JSON.stringify({
        old_tag: oldName,
        new_tag: editForm.value.name.trim(),
        color: editForm.value.color
      })
    })

    if (res.ok) {
      editingTag.value = null
      await store.fetchAnswersList() // 刷新列表数据以更新标签
      store.showToast('标签已更新')
      setTimeout(() => store.fetchAnswersList(), 100)
    } else {
      const errData = await res.json()
      store.showToast(errData.message || '更新标签失败', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('更新标签失败，请检查网络', 'error')
  } finally {
    isProcessing.value = false
  }
}

const deleteTag = async (name) => {
  const confirmed = window.confirm(`确定要删除标签 "${name}" 吗？此操作将移除所有归档上的该标签。`)
  if (!confirmed) return
  
  isProcessing.value = true
  try {
    const res = await apiFetch('/api/tags', {
      method: 'DELETE',
      body: JSON.stringify({ tag: name })
    })

    if (res.ok) {
      store.showToast('标签已删除')
      setTimeout(() => store.fetchAnswersList(), 100)
    } else {
      const errData = await res.json()
      store.showToast(errData.message || '删除标签失败', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('删除标签失败，请检查网络', 'error')
  } finally {
    isProcessing.value = false
  }
}
</script>
