<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-2xl font-semibold text-gray-800 mb-8">偏好设置</h2>
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex-col pr-0 sm:pr-4">
          <span class="block text-sm font-medium text-gray-700">列表显示模式</span>
          <span class="block text-xs text-gray-400 mt-1">选择首页回答列表的展现方式，网格模式更紧凑，列表模式展示更多内容。</span>
        </div>
        
        <div class="flex bg-gray-100 p-1 rounded-lg">
          <button 
            @click="updateViewMode('grid')"
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'grid' ? 'bg-white text-brand shadow-sm' : 'text-gray-500 hover:text-gray-900']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
            <span>网格</span>
          </button>
          <button 
            @click="updateViewMode('list')"
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'list' ? 'bg-white text-brand shadow-sm' : 'text-gray-500 hover:text-gray-900']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            <span>列表</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm mt-6">
      <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-gray-700">自动化动作</h3>
        <span class="text-xs text-gray-400">归档完成后的后续处理</span>
      </div>
      <div class="p-4 md:p-5 space-y-4">
        <div class="flex items-center justify-between h-8">
          <label class="block text-sm font-medium text-gray-700">归档成功后自动推送至 Telegram</label>
          <div class="flex items-center space-x-3">
            <div 
              @click="toggleAutoPush"
              :class="['w-11 h-6 rounded-full cursor-pointer transition-colors relative flex-shrink-0', autoPushEnabled ? 'bg-brand' : 'bg-gray-300']"
            >
              <div 
                :class="['w-5 h-5 bg-white rounded-full absolute top-0.5 shadow transition-transform', autoPushEnabled ? 'translate-x-5.5' : 'translate-x-0.5']"
              ></div>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-400">
          * 开启后，无论是网页手动归档还是发给 Telegram 归档机器人，归档成功后都会调用推送机器人自动发送图文。
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()
const autoPushEnabled = ref(false)

const updateViewMode = (mode) => {
  store.setViewMode(mode)
}

onMounted(async () => {
  try {
    const res = await store.apiFetch('/api/settings/preferences')
    if (res.ok) {
      const data = await res.json()
      autoPushEnabled.value = data.auto_push_enabled === 'true'
    }
  } catch (err) {
    console.error('获取偏好设置失败:', err)
  }
})

const toggleAutoPush = async () => {
  autoPushEnabled.value = !autoPushEnabled.value
  try {
    const res = await store.apiFetch('/api/settings/preferences', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auto_push_enabled: autoPushEnabled.value ? 'true' : 'false' })
    })
    if (res.ok) {
      store.showToast('自动推送设置已更新')
    } else {
      store.showToast('设置更新失败，请重试', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('网络请求失败', 'error')
  }
}
</script>
