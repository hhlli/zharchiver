<template>
  <div class="flex items-center space-x-3 ml-4 border-l pl-4 border-gray-200">
    <!-- 编辑按钮 (占位) -->
    <button 
      class="text-xs text-gray-500 hover:text-blue-600 transition flex items-center space-x-1 cursor-pointer"
      title="编辑当前回答 (占位)"
    >
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
      <span>编辑</span>
    </button>

    <!-- 添加评论按钮 -->
    <button 
      @click="showModal = true"
      class="text-xs text-gray-500 hover:text-green-600 transition flex items-center space-x-1 cursor-pointer"
    >
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"></path></svg>
      <span>添加评论</span>
    </button>

    <!-- 评论弹窗 -->
    <BaseModal 
      :show="showModal" 
      @close="showModal = false"
      maxWidthClass="max-w-lg"
      contentClass="p-6 space-y-4"
      zIndexClass="z-[70]"
    >
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-gray-800">添加评论或笔记</h3>
        <button @click="showModal = false" class="text-gray-400 hover:text-gray-600 text-lg cursor-pointer">&times;</button>
      </div>
      
      <div class="space-y-2">
        <textarea 
          v-model="commentContent"
          rows="4"
          placeholder="在这里粘贴您的评论或学习笔记..."
          class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
          :disabled="loading"
        ></textarea>
      </div>

      <div class="flex justify-end space-x-2 pt-2">
        <button 
          @click="showModal = false"
          class="px-4 py-1.5 border border-gray-300 text-gray-600 rounded-lg text-xs hover:bg-gray-50 transition cursor-pointer"
          :disabled="loading"
        >
          取消
        </button>
        <button 
          @click="submitComment"
          class="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-xs font-medium hover:bg-blue-700 transition disabled:bg-gray-400 cursor-pointer"
          :disabled="loading || !commentContent.trim()"
        >
          {{ loading ? '提交中...' : '提交评论' }}
        </button>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, inject } from 'vue'
import BaseModal from './common/BaseModal.vue'

const props = defineProps({
  answerId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['commentAdded'])

const showModal = ref(false)
const commentContent = ref('')
const loading = ref(false)

const API_BASE = ''
const apiFetch = inject('apiFetch') // 会在 App.vue 提供

const submitComment = async () => {
  if (!commentContent.value.trim()) return
  loading.value = true

  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${props.answerId}/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: commentContent.value.trim() })
    })

    if (res.ok) {
      commentContent.value = ''
      showModal.value = false
      emit('commentAdded')
    } else {
      alert('评论提交失败')
    }
  } catch (err) {
    console.error(err)
    alert('网络请求失败')
  } finally {
    loading.value = false
  }
}
</script>
