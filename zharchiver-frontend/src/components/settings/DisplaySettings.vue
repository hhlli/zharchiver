<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-2xl font-semibold text-primary mb-8">偏好设置</h2>
    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm">
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex-col pr-0 sm:pr-4">
          <span class="block text-sm font-medium text-secondary">列表显示模式</span>
          <span class="block text-xs text-gray-400 mt-1">选择首页回答列表的展现方式，网格模式更紧凑，列表模式展示更多内容。</span>
        </div>
        
        <div class="flex bg-surface-hover p-1 rounded-lg">
          <button 
            @click="updateViewMode('grid')"
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'grid' ? 'bg-surface text-brand shadow-sm' : 'text-muted hover:text-primary']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
            <span>网格</span>
          </button>
          <button 
            @click="updateViewMode('list')"
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'list' ? 'bg-surface text-brand shadow-sm' : 'text-muted hover:text-primary']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            <span>列表</span>
          </button>
        </div>
      </div>

      <div class="border-t border-line-light p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex-col pr-0 sm:pr-4">
          <span class="block text-sm font-medium text-secondary">夜间模式</span>
          <span class="block text-xs text-gray-400 mt-1">选择页面主题颜色。注意：目前仅提供功能入口，尚未完全适配各组件的深色样式。</span>
        </div>
        
        <div class="flex bg-surface-hover p-1 rounded-lg">
          <button 
            @click="updateThemeMode('light')"
            :class="['px-3 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1', store.themeMode === 'light' ? 'bg-surface text-brand shadow-sm' : 'text-muted hover:text-primary']"
          >
            <span>浅色</span>
          </button>
          <button 
            @click="updateThemeMode('dark')"
            :class="['px-3 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1', store.themeMode === 'dark' ? 'bg-surface text-brand shadow-sm' : 'text-muted hover:text-primary']"
          >
            <span>深色</span>
          </button>
          <button 
            @click="updateThemeMode('auto')"
            :class="['px-3 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1', store.themeMode === 'auto' ? 'bg-surface text-brand shadow-sm' : 'text-muted hover:text-primary']"
          >
            <span>跟随系统</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm mt-6">
      <div class="px-5 py-4 border-b border-line-light bg-surface-hover/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-secondary">自动化动作</h3>
        <span class="text-xs text-gray-400">归档完成后的后续处理</span>
      </div>
      <div class="p-4 md:p-5 space-y-4">
        <!-- 自动推送 TG -->
        <div class="flex items-center justify-between h-8">
          <label class="block text-sm font-medium text-secondary">归档成功后自动推送至 Telegram</label>
          <div class="flex items-center space-x-3">
            <div 
              @click="toggleAutoPush"
              :class="['w-11 h-6 rounded-full cursor-pointer transition-colors relative flex-shrink-0', autoPushEnabled ? 'bg-brand' : 'bg-line']"
            >
              <div 
                :class="['w-5 h-5 bg-white rounded-full absolute top-0.5 shadow transition-transform', autoPushEnabled ? 'translate-x-5.5' : 'translate-x-0.5']"
              ></div>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-400 border-b border-line-light pb-4">
          * 开启后，无论是网页手动归档还是发给 Telegram 归档机器人，归档成功后都会调用推送机器人自动发送图文。
        </p>

        <!-- AI 自动归类 -->
        <div class="flex items-center justify-between h-8 pt-2">
          <label class="block text-sm font-medium text-secondary">归档时由 AI 自动生成分类标签</label>
          <div class="flex items-center space-x-3">
            <div 
              @click="toggleAutoCategorize"
              :class="['w-11 h-6 rounded-full cursor-pointer transition-colors relative flex-shrink-0', autoCategorizeEnabled ? 'bg-brand' : 'bg-line']"
            >
              <div 
                :class="['w-5 h-5 bg-white rounded-full absolute top-0.5 shadow transition-transform', autoCategorizeEnabled ? 'translate-x-5.5' : 'translate-x-0.5']"
              ></div>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-400">
          * 开启后，系统在归档网页、推特或截图时，若未指定标签，大模型将根据内容智能匹配您的已有标签或创建新标签。
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
const autoCategorizeEnabled = ref(false)

const updateViewMode = (mode) => {
  store.setViewMode(mode)
}

const updateThemeMode = async (mode) => {
  store.applyTheme(mode)
  await updatePreferences()
}

onMounted(async () => {
  try {
    const res = await store.apiFetch('/api/settings/preferences')
    if (res.ok) {
      const data = await res.json()
      autoPushEnabled.value = data.auto_push_enabled === 'true'
      autoCategorizeEnabled.value = data.auto_categorization_enabled === 'true'
      
      // 同步 themeMode 到 store
      if (data.theme_mode) {
        store.applyTheme(data.theme_mode)
      }
    }
  } catch (err) {
    console.error('获取偏好设置失败:', err)
  }
})

const toggleAutoPush = async () => {
  autoPushEnabled.value = !autoPushEnabled.value
  await updatePreferences()
}

const toggleAutoCategorize = async () => {
  autoCategorizeEnabled.value = !autoCategorizeEnabled.value
  await updatePreferences()
}

const updatePreferences = async () => {
  try {
    const res = await store.apiFetch('/api/settings/preferences', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        auto_push_enabled: autoPushEnabled.value ? 'true' : 'false',
        auto_categorization_enabled: autoCategorizeEnabled.value ? 'true' : 'false',
        theme_mode: store.themeMode
      })
    })
    if (res.ok) {
      store.showToast('设置已更新')
    } else {
      store.showToast('设置更新失败，请重试', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('网络请求失败', 'error')
  }
}
</script>
