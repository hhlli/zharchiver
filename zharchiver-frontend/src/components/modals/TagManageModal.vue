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
      <div v-if="store.allTags && store.allTags.length > 0" class="divide-y divide-gray-100">
        <div v-for="tag in store.allTags" :key="tag.name" class="px-5 py-3 hover:bg-gray-50 transition-colors group">
          
          <div v-if="editingTag !== tag.name" class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <span class="w-2.5 h-2.5 rounded-full flex-shrink-0 shadow-[0_0_2px_rgba(0,0,0,0.1)]" :style="{ background: getTagBackground(tag.color) }"></span>
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
              <label class="text-xs text-gray-500 block mb-1.5">纯色</label>
              <div class="flex flex-wrap gap-2 items-center mb-3">
                <button v-for="(hex, c) in hexColors" :key="c" @click="editForm.color = c" :class="['w-4 h-4 rounded-full cursor-pointer transition flex-shrink-0 shadow-sm border border-black/5', editForm.color === c ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : 'hover:scale-110']" :style="{ background: hex }"></button>
                <div class="w-px h-4 bg-gray-200 mx-1"></div>
                <div class="relative w-4 h-4 rounded-full cursor-pointer transition flex-shrink-0 shadow-sm border border-black/5"
                     :class="[editForm.color.startsWith('#') && !Object.values(hexColors).includes(editForm.color) ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : 'hover:scale-110']"
                     style="background: conic-gradient(red, yellow, lime, aqua, blue, magenta, red);"
                     title="自定义纯色">
                  <div v-if="editForm.color.startsWith('#') && !Object.values(hexColors).includes(editForm.color)" class="absolute inset-0 rounded-full" :style="{ background: editForm.color }"></div>
                  <input type="color" :value="editForm.color.startsWith('#') ? editForm.color : '#ffffff'" @input="editForm.color = $event.target.value" class="absolute -inset-4 w-12 h-12 opacity-0 cursor-pointer" />
                </div>
              </div>

              <label class="text-xs text-gray-500 block mb-1.5">自定义渐变色</label>
              <div class="flex items-center space-x-2">
                <div class="relative w-6 h-6 rounded-md shadow-sm border border-black/10 flex-shrink-0" :style="{ background: gradientStart }">
                  <input type="color" v-model="gradientStart" @change="updateCustomGradient" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                </div>
                <span class="text-gray-300 text-xs">➔</span>
                <div class="relative w-6 h-6 rounded-md shadow-sm border border-black/10 flex-shrink-0" :style="{ background: gradientEnd }">
                  <input type="color" v-model="gradientEnd" @change="updateCustomGradient" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                </div>
                <div class="ml-3 w-16 h-6 rounded-md shadow-sm border border-black/10" :style="{ background: `linear-gradient(135deg, ${gradientStart} 0%, ${gradientEnd} 100%)` }"></div>
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

const gradients = {
  sunset: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
  ocean: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  purplegrad: 'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)',
  forest: 'linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%)',
  night: 'linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)',
  warm: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  cool: 'linear-gradient(135deg, #e0c3fc 0%, #8ec5fc 100%)',
  emerald: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  deep: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
}

const getTagBackground = (colorVal) => {
  if (!colorVal) return hexColors.blue;
  if (colorVal.startsWith('#') || colorVal.startsWith('linear-gradient')) return colorVal;
  return hexColors[colorVal] || gradients[colorVal] || hexColors.blue;
}

const editingTag = ref(null)
const editForm = ref({ name: '', color: 'blue' })
const isProcessing = ref(false)

const gradientStart = ref('#4facfe')
const gradientEnd = ref('#00f2fe')

const updateCustomGradient = () => {
  editForm.value.color = `linear-gradient(135deg, ${gradientStart.value} 0%, ${gradientEnd.value} 100%)`
}

const closeModal = () => {
  editingTag.value = null
  store.showTagManageModal = false
}

const startEdit = (tag) => {
  editingTag.value = tag.name
  editForm.value = { name: tag.name, color: tag.color }
  
  if (tag.color && tag.color.startsWith('linear-gradient')) {
    const match = tag.color.match(/#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})/g)
    if (match && match.length >= 2) {
      gradientStart.value = match[0]
      gradientEnd.value = match[1]
    }
  } else {
    gradientStart.value = '#4facfe'
    gradientEnd.value = '#00f2fe'
  }
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
