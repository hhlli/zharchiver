<template>
  <div class="h-full flex flex-col pt-4 pb-8 px-4 md:px-8 max-w-6xl mx-auto">
    <!-- 顶部控制栏 -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-4 flex-shrink-0 space-y-3 sm:space-y-0">
      <h2 class="text-xl md:text-2xl font-semibold text-gray-800">日志中心</h2>
      <div class="flex space-x-3 w-full sm:w-auto">
        <select v-model="levelFilter" class="flex-1 sm:flex-none border border-gray-300 bg-white rounded-md px-3 py-1.5 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500">
          <option value="ALL">全部日志</option>
          <option value="INFO">INFO</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
        </select>
        <BaseButton variant="outline" @click="clearLogs" customClass="flex-1 sm:flex-none">
          清空屏幕
        </BaseButton>
      </div>
    </div>

    <!-- 日志输出终端 -->
    <div class="flex-1 bg-white rounded-xl shadow-sm overflow-hidden flex flex-col border border-gray-200">
      <div class="flex-1 overflow-y-auto p-4 font-mono text-[12px] leading-relaxed select-text" ref="logContainer">
        <div v-for="(log, index) in filteredLogs" :key="index" class="mb-1 flex space-x-3 hover:bg-gray-50 px-2 py-0.5 rounded transition-colors">
          <span class="text-gray-400 flex-shrink-0 w-20">{{ log.time }}</span>
          <span :class="levelColor(log.level)" class="flex-shrink-0 w-12 font-semibold">{{ log.level }}</span>
          <span class="text-gray-700 break-all whitespace-pre-wrap">{{ log.message }}</span>
        </div>
        <div v-if="filteredLogs.length === 0" class="text-gray-400 text-center mt-10">等待接收日志数据...</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import BaseButton from '../common/BaseButton.vue'

const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const API_WS = `${protocol}//${window.location.host}/api/logs/ws`;

const logs = ref([])
const levelFilter = ref('ALL')
const logContainer = ref(null)
let ws = null

// 根据等级过滤日志
const filteredLogs = computed(() => {
  if (levelFilter.value === 'ALL') return logs.value
  return logs.value.filter(log => log.level === levelFilter.value)
})

// 日志等级着色
const levelColor = (level) => {
  switch (level) {
    case 'INFO': return 'text-blue-600'
    case 'WARN': return 'text-yellow-600'
    case 'ERROR': return 'text-red-600'
    default: return 'text-gray-500'
  }
}

const clearLogs = () => {
  logs.value = []
}

// 自动滚动到最底部
const scrollToBottom = () => {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

// 建立 WebSocket 连接
const initWebSocket = () => {
  const token = localStorage.getItem('token') || ''
  const wsUrl = token ? `${API_WS}?token=${token}` : API_WS
  ws = new WebSocket(wsUrl)
  
  ws.onmessage = (event) => {
    try {
      const logEntry = JSON.parse(event.data)
      logs.value.push(logEntry)
      
      // 限制最大保留行数防止内存溢出
      if (logs.value.length > 2000) {
        logs.value.shift()
      }
      
      scrollToBottom()
    } catch (e) {
      console.error('日志解析失败:', e)
    }
  }

  ws.onerror = (error) => {
    console.error('WebSocket 错误:', error)
  }

  ws.onclose = () => {
    // 断线后 3 秒自动重连
    setTimeout(initWebSocket, 3000)
  }
}

onMounted(() => {
  initWebSocket()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})
</script>