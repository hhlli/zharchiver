<template>
  <div class="mt-8 pt-8 border-t border-gray-200">
    <h3 class="text-lg font-bold text-gray-800 mb-4 flex items-center space-x-2">
      <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z"></path></svg>
      <span>评论与笔记 ({{ comments.length }})</span>
    </h3>

    <div v-if="loading" class="text-sm text-gray-500 py-4">
      正在加载评论...
    </div>
    
    <div v-else-if="comments.length === 0" class="text-sm text-gray-400 py-4 italic bg-gray-50 rounded-lg text-center">
      暂无评论或笔记，点击上方按钮添加。
    </div>

    <div v-else class="space-y-4">
      <div 
        v-for="comment in comments" 
        :key="comment.id"
        class="bg-gray-50 border border-gray-100 rounded-xl p-4 shadow-sm"
      >
        <div class="text-[11px] text-gray-400 mb-2 border-b border-gray-200/60 pb-2">
          添加于 {{ new Date(comment.created_at).toLocaleString() }}
        </div>
        <div class="text-sm text-gray-700 whitespace-pre-wrap leading-relaxed">
          {{ comment.content }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, inject } from 'vue'

const props = defineProps({
  answerId: {
    type: String,
    required: true
  },
  refreshKey: {
    type: Number,
    default: 0
  }
})

const comments = ref([])
const loading = ref(false)

const API_BASE = ''
const apiFetch = inject('apiFetch') // 会在 App.vue 提供

const fetchComments = async () => {
  if (!props.answerId) return
  loading.value = true
  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${props.answerId}/comments`)
    if (res.ok) {
      comments.value = await res.json()
    }
  } catch (err) {
    console.error('获取评论失败', err)
  } finally {
    loading.value = false
  }
}

watch(() => props.answerId, fetchComments, { immediate: true })
watch(() => props.refreshKey, fetchComments)
</script>
