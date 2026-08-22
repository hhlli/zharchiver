<template>
  <div class="flex items-center space-x-3 ml-4 border-l pl-4 border-line">
    <!-- 编辑按钮 -->
    <button 
      @click="$emit('edit')"
      class="text-xs text-muted hover:text-brand transition flex items-center space-x-1 cursor-pointer"
      title="编辑当前回答"
    >
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
      <span>编辑</span>
    </button>

    <!-- 发送至 Telegram 按钮 -->
    <button 
      @click="shareToTelegram"
      :disabled="sharing"
      class="text-xs text-muted hover:text-blue-500 transition flex items-center space-x-1 cursor-pointer disabled:opacity-50"
      title="一键推送到 Telegram"
    >
      <svg v-if="sharing" class="animate-spin w-3.5 h-3.5 text-blue-500" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path></svg>
      <span>{{ sharing ? '发送中...' : 'TG' }}</span>
    </button>

    <!-- 收藏按钮 -->
    <button 
      v-if="store.currentAnswer"
      @click="store.toggleFavorite(store.currentAnswer, false)"
      :class="['text-xs transition flex items-center cursor-pointer', store.currentAnswer.is_favorite ? 'text-yellow-400 hover:text-yellow-500' : 'text-muted hover:text-yellow-400']"
      title="收藏/取消收藏"
    >
      <svg class="w-3.5 h-3.5" :fill="store.currentAnswer.is_favorite ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../stores/archive'

const props = defineProps({
  answerId: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['edit'])

const sharing = ref(false)

const API_BASE = ''
const store = useArchiveStore()
const apiFetch = store.apiFetch
const showAlert = store.showAlert

const shareToTelegram = async () => {
  if (sharing.value) return
  sharing.value = true

  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${props.answerId}/share/telegram`, {
      method: 'POST'
    })

    if (res.ok) {
      store.showToast('已成功推送到您的 Telegram！')
    } else {
      const errText = await res.text()
      store.showToast('发送失败', 'error')
    }
  } catch (err) {
    console.error('Telegram share failed:', err)
    store.showToast('网络请求失败', 'error')
  } finally {
    sharing.value = false
  }
}
</script>
