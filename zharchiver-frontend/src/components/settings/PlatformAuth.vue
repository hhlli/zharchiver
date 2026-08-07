<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-2xl font-semibold text-gray-800 mb-8">平台鉴权</h2>
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex-col pr-0 sm:pr-4">
          <div class="flex items-center gap-3">
            <span class="block text-sm font-medium text-gray-700">知乎账号授权</span>
            <span v-if="isZhihuConfigured" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
              已鉴权
            </span>
          </div>
          <span class="block text-xs text-gray-400 mt-1">使用知乎 App 扫码登录，获取新鲜的凭证。凭证过期或抓取失败时可重新授权。</span>
        </div>
        <BaseButton variant="primary" @click="openQRModal">
          {{ isZhihuConfigured ? '重新扫描授权' : '扫描二维码登录' }}
        </BaseButton>
      </div>
    </div>

    <!-- 扫码弹窗 -->
    <BaseModal
      :show="showModal"
      @close="closeModal"
    >
      <div class="p-6 flex flex-col items-center justify-center min-h-[250px]">
        <div v-if="wsStatus === 'loading' || wsStatus === 'connecting'" class="flex flex-col items-center">
          <div class="w-8 h-8 border-4 border-blue-200 border-t-brand rounded-full animate-spin mb-4"></div>
          <p class="text-sm text-gray-600">{{ wsMessage || '正在连接安全环境...' }}</p>
        </div>
        
        <div v-else-if="wsStatus === 'qrcode' || wsStatus === 'waiting'" class="flex flex-col items-center">
          <img :src="qrImage" alt="Zhihu QR Code" class="w-48 h-48 border border-gray-200 rounded-lg shadow-sm mb-4" />
          <p class="text-sm font-medium text-gray-800">{{ wsMessage }}</p>
          <p class="text-xs text-gray-400 mt-2">请使用知乎 App 扫描上方二维码</p>
        </div>

        <div v-else-if="wsStatus === 'success' || wsStatus === 'done'" class="flex flex-col items-center">
          <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
          </div>
          <p class="text-sm font-medium text-gray-800">{{ wsMessage }}</p>
        </div>

        <div v-else-if="wsStatus === 'error'" class="flex flex-col items-center">
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </div>
          <p class="text-sm font-medium text-red-600">{{ wsMessage }}</p>
          <BaseButton variant="outline" @click="retryConnection" customClass="mt-4">
            重试
          </BaseButton>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseButton from '../common/BaseButton.vue'
import BaseModal from '../common/BaseModal.vue'

const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const API_WS_QR = `${protocol}//${window.location.host}/api/auth/zhihu/qrcode/ws`
const API_AUTH_STATUS = '/api/auth/status'

const store = useArchiveStore()
const apiFetch = store.apiFetch

const showModal = ref(false)
const wsStatus = ref('connecting') // connecting, loading, qrcode, waiting, success, done, error
const wsMessage = ref('')
const qrImage = ref('')
const isZhihuConfigured = ref(false)

let ws = null

onMounted(() => {
  checkStatus()
})

const checkStatus = async () => {
  try {
    const res = await apiFetch(API_AUTH_STATUS)
    if (res.ok) {
      const data = await res.json()
      isZhihuConfigured.value = data.zhihu_configured || false
    }
  } catch (err) {
    console.error('Failed to check auth status', err)
  }
}

const openQRModal = () => {
  showModal.value = true
  connectWebSocket()
}

const closeModal = () => {
  showModal.value = false
  if (ws) {
    ws.close()
    ws = null
  }
}

const retryConnection = () => {
  if (ws) ws.close()
  connectWebSocket()
}

const connectWebSocket = () => {
  wsStatus.value = 'connecting'
  wsMessage.value = '正在连接...'
  qrImage.value = ''
  
  const token = localStorage.getItem('token') || ''
  const wsUrl = token ? `${API_WS_QR}?token=${token}` : API_WS_QR
  
  ws = new WebSocket(wsUrl)
  
  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      wsStatus.value = data.status
      wsMessage.value = data.message
      
      if (data.image) {
        qrImage.value = data.image
      }
      
      if (data.status === 'done') {
        isZhihuConfigured.value = true
        setTimeout(() => {
          closeModal()
        }, 2000)
      }
    } catch (e) {
      console.error('WS 解析失败', e)
    }
  }
  
  ws.onerror = () => {
    wsStatus.value = 'error'
    wsMessage.value = 'WebSocket 连接失败'
  }
  
  ws.onclose = () => {
    if (wsStatus.value !== 'done' && wsStatus.value !== 'error') {
      wsStatus.value = 'error'
      wsMessage.value = '连接已意外断开'
    }
  }
}
</script>