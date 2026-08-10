<template>
  <transition name="fade-slide">
    <div v-if="show" class="fixed bottom-6 left-1/2 transform -translate-x-1/2 bg-white text-gray-800 rounded-xl shadow-xl border border-gray-200 overflow-hidden z-50 min-w-[280px] pointer-events-auto">
      <div class="px-4 py-3">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center space-x-2">
            <svg v-if="!isSuccess && !isFailed" class="animate-spin h-4 w-4 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else-if="isSuccess" class="h-4 w-4 text-green-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <svg v-else class="h-4 w-4 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            <span class="text-sm font-medium" :class="isSuccess ? 'text-green-600' : isFailed ? 'text-red-600' : 'text-gray-800'">
              {{ title }}
            </span>
          </div>
          <span v-if="!isSuccess && !isFailed && currentFileName" class="text-xs text-gray-400 truncate max-w-[120px]" :title="currentFileName">
            {{ currentFileName }}
          </span>
          <span v-if="!isSuccess && !isFailed" class="text-xs font-semibold text-blue-500">{{ percent }}%</span>
        </div>
        
        <div class="w-full bg-gray-100 rounded-full h-1 overflow-hidden">
          <div class="h-1 rounded-full transition-all duration-300 ease-out" :class="isSuccess ? 'bg-green-500' : isFailed ? 'bg-red-400' : 'bg-blue-500'" :style="{ width: percent + '%' }"></div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()
const show = ref(false)
const title = ref('后台解析任务进行中...')
const percent = ref(0)
const currentFileName = ref('')
const isSuccess = ref(false)
const isFailed = ref(false)
let hideTimeout = null

const handleLogEvent = (e) => {
  try {
    const data = JSON.parse(e.data)
    
    if (data.level === 'TASK_START') {
      show.value = true
      isSuccess.value = false
      isFailed.value = false
      percent.value = 0
      title.value = '后台解析任务进行中...'
      currentFileName.value = ''
      clearTimeout(hideTimeout)
    } 
    else if (data.level === 'PROGRESS') {
      show.value = true
      isSuccess.value = false
      clearTimeout(hideTimeout)
      
      const parts = data.message.split('|')
      if (parts.length >= 2) {
        percent.value = parseInt(parts[0], 10)
        currentFileName.value = parts[1]
        title.value = '正在下载媒体文件...'
      }
    }
    else if (data.level === 'TASK_SUCCESS') {
      show.value = true
      percent.value = 100
      isSuccess.value = true
      isFailed.value = false
      title.value = '归档成功！'
      currentFileName.value = ''
      
      // 触发更新列表
      store.fetchAnswersList()
      
      clearTimeout(hideTimeout)
      hideTimeout = setTimeout(() => {
        show.value = false
      }, 2000)
    }
    else if (data.level === 'TASK_FAILED') {
      show.value = true
      percent.value = 100
      isSuccess.value = false
      isFailed.value = true
      title.value = '解析失败：' + data.message
      
      clearTimeout(hideTimeout)
      hideTimeout = setTimeout(() => {
        show.value = false
      }, 4000)
    }
  } catch (err) {
    // 解析失败不处理
  }
}

let ws = null

onMounted(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // store.API_BASE 可能自带 http:// 或者为空（同域）
  let host = window.location.host
  if (store.API_BASE) {
    try {
      const url = new URL(store.API_BASE)
      host = url.host
    } catch(e) {}
  }
  const wsUrl = `${protocol}//${host}/api/logs/ws`
  
  ws = new WebSocket(wsUrl)
  ws.onmessage = handleLogEvent
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})
</script>

<style scoped>
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}
</style>
