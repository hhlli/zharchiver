<template>
  <div class="pt-6 border-t border-line">
    <h3 class="text-lg font-bold text-primary mb-4 flex items-center space-x-2">
      <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z"></path></svg>
      <span>评论与笔记 ({{ comments.length }})</span>
    </h3>

    <div v-if="loading" class="text-sm text-muted py-4">
      正在加载评论...
    </div>
    
    <div v-else-if="comments.length === 0" class="text-sm text-gray-400 py-4 italic bg-surface-hover rounded-lg text-center">
      暂无评论或笔记，点击上方按钮添加。
    </div>

    <div v-else class="space-y-4">
      <div 
        v-for="comment in comments" 
        :key="comment.id"
        class="bg-surface-hover border border-line-light rounded-xl p-4 shadow-sm group"
      >
        <div class="flex items-center justify-between text-[11px] text-gray-400 mb-2 border-b border-line/60 pb-2">
          <span>添加于 {{ new Date(comment.created_at).toLocaleString() }}</span>
          <div class="flex items-center space-x-3">
            <button @click="startEdit(comment)" class="text-gray-400 hover:text-blue-500 transition cursor-pointer" title="编辑">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
            </button>
            <button @click="deleteComment(comment)" class="text-gray-400 hover:text-red-500 transition cursor-pointer" title="删除">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
        </div>
        
        <div v-if="editingCommentId === comment.id" class="space-y-2">
          <textarea 
            v-model="editContent"
            rows="3"
            class="w-full text-sm text-secondary border border-line rounded-lg p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          ></textarea>
          <div class="flex justify-end space-x-2">
            <button @click="editingCommentId = null" class="px-3 py-1 text-xs text-secondary bg-line hover:bg-line rounded-md transition cursor-pointer">取消</button>
            <button @click="saveEdit(comment)" class="px-3 py-1 text-xs text-white bg-brand hover:bg-brand-hover rounded-md transition cursor-pointer" :disabled="savingEdit">保存</button>
          </div>
        </div>
        <div v-else class="text-sm text-secondary whitespace-pre-wrap leading-relaxed">
          {{ comment.content }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useArchiveStore } from '../stores/archive'

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
const store = useArchiveStore()
const apiFetch = store.apiFetch

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

const editingCommentId = ref(null)
const editContent = ref('')
const savingEdit = ref(false)

const startEdit = (comment) => {
  editingCommentId.value = comment.id
  editContent.value = comment.content
}

const saveEdit = async (comment) => {
  if (!editContent.value.trim()) return
  savingEdit.value = true
  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${props.answerId}/comments/${comment.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: editContent.value.trim() })
    })
    if (res.ok) {
      editingCommentId.value = null
      await fetchComments()
    } else {
      store.showToast('更新失败', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('网络请求失败', 'error')
  } finally {
    savingEdit.value = false
  }
}

const deleteComment = (comment) => {
  store.itemToDeleteType = 'comment'
  store.itemToDeleteCallback = async () => {
    try {
      const res = await apiFetch(`${API_BASE}/api/answers/${props.answerId}/comments/${comment.id}`, {
        method: 'DELETE'
      })
      if (res.ok) {
        await fetchComments()
      } else {
        store.showToast('删除失败', 'error')
      }
    } catch (err) {
      console.error(err)
      store.showToast('网络请求失败', 'error')
    }
  }
  store.itemToDelete = comment
}

watch(() => props.answerId, fetchComments, { immediate: true })
watch(() => props.refreshKey, fetchComments)
</script>
