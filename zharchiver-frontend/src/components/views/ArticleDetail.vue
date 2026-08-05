<template>
  <div class="max-w-6xl mx-auto space-y-6 md:h-full pb-8">
    <button 
      @click="store.currentAnswer = null" 
      class="inline-flex items-center space-x-1 text-sm md:text-xs font-medium text-gray-500 hover:text-blue-600 mb-2 md:mb-2 cursor-pointer transition select-none py-2 px-1 -ml-1 md:py-0 md:px-0 md:ml-0"
    >
      <svg class="w-5 h-5 md:w-4 md:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
      <span>返回列表</span>
    </button>

    <header class="border-b pb-4 space-y-2">
      <h1 class="text-xl md:text-2xl font-bold text-gray-900 leading-snug">{{ store.currentAnswer?.title }}</h1>
      <div class="flex flex-col sm:flex-row sm:items-center text-xs text-gray-500 space-y-2 sm:space-y-0">
        <div class="flex items-center">
          <span class="font-medium text-gray-700">归档：{{ formatDate(store.currentAnswer?.saved_at) }}</span>
        </div>
        
        <div class="sm:ml-auto flex items-center space-x-3">
          <div class="flex items-center text-xs">
            <span 
              @dblclick="startEditTag" 
              class="bg-blue-50 text-blue-600 border border-blue-100 px-2 py-0.5 rounded-md font-medium tracking-wider cursor-pointer select-none hover:bg-blue-100 transition"
              title="双击编辑标签"
            >
              {{ store.currentAnswer?.tag || 'ANSWER' }}
            </span>
          </div>
          <AnswerActions :answerId="store.currentAnswer?.answer_id" @commentAdded="onCommentAdded" />
        </div>
      </div>
    </header>

    <article 
      class="prose prose-sm md:prose-base prose-slate max-w-none text-gray-800 leading-relaxed space-y-4 prose-img:max-w-full prose-img:rounded-lg prose-img:shadow-sm"
      v-html="processHtmlContent(store.currentAnswer?.content_html)"
    ></article>

    <div class="mt-8 pt-6 border-t border-gray-100 flex flex-col items-center justify-center text-sm text-gray-500 space-y-1 bg-gray-50/50 rounded-xl pb-6">
      <span class="font-medium text-gray-700 mt-4">原文作者：{{ store.currentAnswer?.author_name }}</span>
      <span class="text-xs">原文发布：{{ formatTimestamp(store.currentAnswer?.created_time) }}</span>
    </div>

    <AnswerComments :answerId="store.currentAnswer?.answer_id" :refreshKey="commentRefreshKey" />

    <!-- 编辑标签弹窗 -->
    <BaseModal 
      :show="isEditingTag" 
      @close="isEditingTag = false"
      maxWidthClass="max-w-sm"
      contentClass="p-6 space-y-4"
      zIndexClass="z-[60]"
    >
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-gray-800">编辑标签</h3>
        <button @click="isEditingTag = false" class="text-gray-400 hover:text-gray-600 text-lg cursor-pointer">&times;</button>
      </div>
      
      <div class="space-y-4 pt-2">
        <div class="space-y-2">
          <label class="text-xs text-gray-500">标签名称</label>
          <input type="text" v-model="editTagValue" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" @keyup.enter="saveTag">
        </div>
        <div class="space-y-2">
          <label class="text-xs text-gray-500">标签颜色</label>
          <div class="flex space-x-2">
            <button v-for="(hex, c) in hexColors" :key="c" @click="editTagColor = c" :class="['w-6 h-6 rounded-full cursor-pointer transition', editTagColor === c ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : 'hover:scale-110']" :style="{ backgroundColor: hex }"></button>
          </div>
        </div>
      </div>

      <div class="flex justify-end space-x-2 pt-4">
        <button 
          @click="isEditingTag = false"
          class="px-4 py-1.5 border border-gray-300 text-gray-600 rounded-lg text-xs hover:bg-gray-50 transition cursor-pointer"
        >取消</button>
        <button 
          @click="saveTag"
          class="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-xs font-medium hover:bg-blue-700 transition cursor-pointer"
        >保存</button>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import AnswerActions from '../AnswerActions.vue'
import AnswerComments from '../AnswerComments.vue'
import BaseModal from '../common/BaseModal.vue'
import DOMPurify from 'dompurify'

const store = useArchiveStore()

const commentRefreshKey = ref(0)
const onCommentAdded = () => {
  commentRefreshKey.value++
}

const isEditingTag = ref(false)
const editTagValue = ref('')
const editTagColor = ref('blue')

const hexColors = {
  blue: '#3b82f6',
  red: '#ef4444',
  green: '#10b981',
  yellow: '#f59e0b',
  purple: '#8b5cf6'
}

const startEditTag = () => {
  editTagValue.value = store.currentAnswer?.tag || ''
  editTagColor.value = store.currentAnswer?.tag_color || 'blue'
  isEditingTag.value = true
}

const saveTag = async () => {
  if (!store.currentAnswer) return

  try {
    const res = await store.apiFetch(`/api/answers/${store.currentAnswer.answer_id}/tag`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ tag: editTagValue.value.trim(), color: editTagColor.value })
    })
    if (!res.ok) throw new Error("标签更新失败")
    
    // 乐观更新
    store.currentAnswer.tag = editTagValue.value.trim()
    store.currentAnswer.tag_color = editTagColor.value
    isEditingTag.value = false
    store.fetchAnswersList()
  } catch (err) {
    store.showAlert('错误', err.message)
  }
}

const processHtmlContent = (html) => {
  if (!html) return ''
  const token = localStorage.getItem('token') || ''
  
  const processed = html.replace(/src="\/storage\/([^"]+)"/g, (match, p1) => {
    const url = `${store.API_BASE}/storage/${p1}`
    return token ? `src="${url}?token=${token}"` : `src="${url}"`
  })
  
  return DOMPurify.sanitize(processed)
}

const formatTimestamp = (ts) => {
  if (!ts) return ''
  return new Date(ts * 1000).toLocaleString()
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}
</script>
