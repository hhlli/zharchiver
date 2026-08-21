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
        <div v-for="tag in store.allTags" :key="tag.name" class="px-5 py-3 hover:bg-gray-50 transition-colors">

          <!-- 展示态 -->
          <div v-if="editingTag !== tag.name" class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <span
                class="w-2.5 h-2.5 rounded-full flex-shrink-0 shadow-[0_0_2px_rgba(0,0,0,0.12)]"
                :style="{ background: resolveTagBackground(tag.color) }"
              ></span>
              <span class="text-sm text-gray-700">{{ tag.name }}</span>
            </div>
            <div class="flex items-center space-x-1">
              <button @click="startEdit(tag)" class="text-[11px] text-gray-400 hover:text-brand transition px-2 py-1 rounded cursor-pointer">编辑</button>
              <button @click="deleteTag(tag.name)" :disabled="isProcessing" class="text-[11px] text-gray-400 hover:text-red-500 transition px-2 py-1 rounded cursor-pointer">删除</button>
            </div>
          </div>

          <!-- 编辑态 -->
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
              <label class="text-xs text-gray-500 block mb-2">标签颜色</label>
              <!-- 实时预览 + 触发调色盘 -->
              <div class="flex items-center space-x-3">
                <div
                  class="w-5 h-5 rounded-full flex-shrink-0 shadow-sm border border-black/10 cursor-pointer transition hover:scale-110"
                  :style="{ background: currentPreviewColor }"
                  title="当前颜色预览"
                ></div>
                <div class="flex-1 bg-gray-50 rounded-lg border border-gray-100 p-2 flex justify-center">
                  <ColorPicker
                    isWidget
                    :pureColor="pureColor"
                    :gradientColor="gradientColor"
                    :activeKey="activeKey"
                    @update:pureColor="onPureColorChange"
                    @update:gradientColor="onGradientColorChange"
                    @update:activeKey="val => activeKey = val"
                  />
                </div>
              </div>
            </div>
            <div class="flex justify-end space-x-2 pt-1">
              <button @click="editingTag = null" class="px-3 py-1 border border-gray-300 text-gray-600 rounded-md text-xs hover:bg-gray-100 transition cursor-pointer">取消</button>
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
import { ref, computed } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseModal from '../common/BaseModal.vue'
import { ColorPicker } from 'vue3-colorpicker'
import 'vue3-colorpicker/style.css'
import { resolveTagBackground, colorToPickerState, pickerStateToColor } from '../../utils/tagColor'

const store = useArchiveStore()
const apiFetch = store.apiFetch

const editingTag = ref(null)
const editForm = ref({ name: '' })
const isProcessing = ref(false)

// ColorPicker 三状态
const pureColor = ref('#3b82f6')
const gradientColor = ref('linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)')
const activeKey = ref('pure')

// 实时预览色
const currentPreviewColor = computed(() =>
  pickerStateToColor(activeKey.value, pureColor.value, gradientColor.value)
)

const onPureColorChange = (val) => { pureColor.value = val }
const onGradientColorChange = (val) => { gradientColor.value = val }

const closeModal = () => {
  editingTag.value = null
  store.showTagManageModal = false
}

const startEdit = (tag) => {
  editingTag.value = tag.name
  editForm.value = { name: tag.name }
  const state = colorToPickerState(tag.color)
  pureColor.value = state.pureColor
  gradientColor.value = state.gradientColor
  activeKey.value = state.activeKey
}

const saveTag = async (oldName) => {
  if (!editForm.value.name.trim()) return
  isProcessing.value = true

  const finalColor = pickerStateToColor(activeKey.value, pureColor.value, gradientColor.value)

  try {
    const res = await apiFetch('/api/tags', {
      method: 'PUT',
      body: JSON.stringify({
        old_tag: oldName,
        new_tag: editForm.value.name.trim(),
        color: finalColor,
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
</script>
