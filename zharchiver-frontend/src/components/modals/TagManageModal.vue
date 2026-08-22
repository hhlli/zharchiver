<template>
  <BaseModal 
    :show="store.showTagManageModal" 
    :closeOnOutside="true"
    @close="closeModal"
  >
    <div class="px-5 pt-5 pb-3 flex items-center justify-between border-b border-line-light">
      <h3 class="text-base font-semibold text-primary">管理标签</h3>
      <button @click="closeModal" class="text-gray-400 hover:text-secondary text-lg cursor-pointer">&times;</button>
    </div>
    <div class="p-0 max-h-[70vh] overflow-y-auto">
      <div v-if="store.allTags && store.allTags.length > 0" class="divide-y divide-line-light">
        <div v-for="tag in store.allTags" :key="tag.name" class="px-5 py-3 hover:bg-surface-hover transition-colors">

          <!-- 展示态 -->
          <div 
            v-if="editingTag !== tag.name" 
            class="flex items-center justify-between"
            draggable="true"
            @dragstart="onDragStart($event, tag.name)"
            @dragover.prevent
            @dragenter.prevent="onDragEnter(tag.name)"
            @drop="onDrop(tag.name)"
            :class="{'opacity-50': draggingTag === tag.name}"
          >
            <div class="flex items-center space-x-3 pointer-events-none">
              <span
                class="w-2.5 h-2.5 rounded-full flex-shrink-0 shadow-[0_0_2px_rgba(0,0,0,0.12)]"
                :style="{ background: resolveTagBackground(tag.color) }"
              ></span>
              <span class="text-sm text-secondary font-medium">{{ tag.name }}</span>
            </div>
            <div class="flex items-center space-x-1.5 text-gray-400">
              <!-- 拖拽把手 -->
              <div class="p-1.5 cursor-grab active:cursor-grabbing hover:text-secondary transition" title="拖动排序">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"></path></svg>
              </div>
              <!-- 编辑 -->
              <button @click="startEdit(tag)" class="p-1.5 hover:text-brand transition rounded cursor-pointer" title="编辑">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
              </button>
              <!-- 删除 -->
              <button @click="deleteTag(tag.name)" :disabled="isProcessing" class="p-1.5 hover:text-red-500 transition rounded cursor-pointer" title="删除">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
              </button>
            </div>
          </div>

          <!-- 编辑态 -->
          <div v-else class="space-y-3">
            <div>
              <label class="text-xs text-muted block mb-1">标签名称</label>
              <input
                v-model="editForm.name"
                type="text"
                class="w-full text-sm border border-line rounded-md px-2.5 py-1.5 focus:outline-none focus:ring-2 focus:ring-blue-500 bg-surface"
              />
            </div>
            <div>
              <label class="text-xs text-muted block mb-2">标签颜色</label>
              <TagColorPicker v-model="currentColor" />
            </div>
            <div class="flex justify-end space-x-2 pt-1">
              <button @click="editingTag = null" class="px-3 py-1 border border-line text-secondary rounded-md text-xs hover:bg-surface-hover transition cursor-pointer">取消</button>
              <button @click="saveTag(tag.name)" :disabled="isProcessing || !editForm.name.trim()" class="px-3 py-1 bg-brand text-white rounded-md text-xs font-medium hover:bg-brand-hover transition cursor-pointer disabled:opacity-50">保存</button>
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
import TagColorPicker from '../common/TagColorPicker.vue'
import { resolveTagBackground } from '../../utils/tagColor'

const store = useArchiveStore()
const apiFetch = store.apiFetch

const editingTag = ref(null)
const editForm = ref({ name: '' })
const isProcessing = ref(false)
const currentColor = ref('#3b82f6')

const closeModal = () => {
  editingTag.value = null
  store.showTagManageModal = false
}

const startEdit = (tag) => {
  editingTag.value = tag.name
  editForm.value = { name: tag.name }
  // resolveTagBackground 兼容旧关键词，确保 currentColor 始终是合法的 CSS 颜色值
  currentColor.value = resolveTagBackground(tag.color)
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
        color: currentColor.value,
      })
    })

    if (res.ok) {
      editingTag.value = null
      store.showToast('标签已更新')
      await store.fetchAnswersList()
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
      await store.fetchAnswersList()
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

const draggingTag = ref(null)

const onDragStart = (e, tagName) => {
  draggingTag.value = tagName
  e.dataTransfer.effectAllowed = 'move'
  // 必须设置 data，不然 firefox 不支持拖拽
  e.dataTransfer.setData('text/plain', tagName)
}

const onDragEnter = (targetTagName) => {
  if (!draggingTag.value || draggingTag.value === targetTagName) return
  
  const tags = [...store.allTags]
  const fromIdx = tags.findIndex(t => t.name === draggingTag.value)
  const toIdx = tags.findIndex(t => t.name === targetTagName)
  
  if (fromIdx === -1 || toIdx === -1) return
  
  // 交换数组元素
  const temp = tags[fromIdx]
  tags.splice(fromIdx, 1)
  tags.splice(toIdx, 0, temp)
  
  store.allTags = tags
}

const onDrop = async (targetTagName) => {
  draggingTag.value = null
  
  // 保存到后端
  const sortOrder = store.allTags.map(t => t.name)
  try {
    await apiFetch('/api/settings/preferences', {
      method: 'POST',
      body: JSON.stringify({
        tag_sort_order: JSON.stringify(sortOrder)
      })
    })
  } catch (err) {
    console.error('保存标签排序失败:', err)
  }
}
</script>
