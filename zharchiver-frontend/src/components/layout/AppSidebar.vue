<template>
  <aside class="flex inset-y-0 left-0 z-50 w-64 md:w-52 bg-[#f5f5f7] border-r border-gray-200 p-3 flex-col justify-between flex-shrink-0 h-full overflow-y-auto">
    
    <!-- 主界面导航模式 -->
    <div v-if="store.currentView === 'home'" class="space-y-6">
      <div>
        <nav class="space-y-1">
          <button 
            @click="store.goHome()"
            :class="['w-full flex items-center space-x-2.5 px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeCategory === 'all' && !store.currentAnswer ? 'bg-brand text-white shadow-sm' : 'text-gray-700 hover:bg-gray-200/60']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path></svg>
            <span>所有回答</span>
          </button>
        </nav>
      </div>

      <div v-if="store.allTags && store.allTags.length > 0">
        <div class="flex items-center justify-between px-3 mb-2">
          <div class="text-xs font-semibold text-gray-400 tracking-wider">标签</div>
          <button 
            @click="store.showTagManageModal = true"
            class="text-gray-400 hover:text-blue-500 transition-colors"
            title="管理标签"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
          </button>
        </div>
        <div class="space-y-1">
          <button
            v-for="tag in store.allTags"
            :key="tag.name"
            @click="selectCategory(tag.name)"
            :class="['w-full flex items-center space-x-2 px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeCategory === tag.name && !store.currentAnswer ? 'bg-blue-100 text-brand-hover shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >
            <span class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: getHexColor(tag.color) }"></span>
            <span class="truncate">{{ tag.name }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 设置界面导航模式 -->
    <div v-else-if="store.currentView === 'settings'" class="space-y-6">
      <div>
        <button 
          @click="store.currentView = 'home'" 
          class="w-full flex items-center space-x-2 px-3 py-2 text-sm font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-200/60 rounded-lg transition mb-4 cursor-pointer"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          <span>返回主界面</span>
        </button>

        <div class="text-xs font-semibold text-gray-400 px-3 mb-2 tracking-wider">系统设置</div>
        <nav class="space-y-1">
          <button 
            @click="store.activeSetting = 'security'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'security' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >账户安全</button>
          <button 
            @click="store.activeSetting = 'auth'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'auth' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >平台鉴权</button>
          <button 
            @click="store.activeSetting = 'data'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'data' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >数据管理</button>
          <button 
            @click="store.activeSetting = 'display'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'display' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >偏好设置</button>
          <button 
            @click="store.activeSetting = 'ai'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'ai' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >工具</button>
          <button 
            @click="store.activeSetting = 'logs'"
            :class="['w-full text-left px-3 py-1.5 rounded-lg text-sm font-medium transition cursor-pointer', store.activeSetting === 'logs' ? 'bg-gray-200/80 text-gray-900 shadow-sm' : 'text-gray-600 hover:bg-gray-200/50']"
          >日志中心</button>
        </nav>
      </div>
    </div>

    <div class="flex items-center justify-between mt-auto">
      <div class="flex items-center space-x-1">
        <button @click="store.isDesktopSidebarOpen = false" class="hidden md:flex p-1 rounded text-gray-400 hover:text-gray-700 hover:bg-gray-200 transition cursor-pointer items-center justify-center" title="收起边栏">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"></path></svg>
        </button>
        <div class="text-[11px] text-gray-400 font-mono hidden md:block">ZHArchiver v1.1.7</div>
      </div>
      <button @click="openSettings" :class="['p-1 rounded transition cursor-pointer', store.currentView === 'settings' ? 'text-gray-800 bg-gray-200' : 'text-gray-400 hover:text-gray-700 hover:bg-gray-200']" title="设置">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()

const selectCategory = (category) => {
  store.activeCategory = category
  store.currentAnswer = null
  store.currentGroup = null
  // 切换标签后滚回顶部
  setTimeout(() => {
    const el = document.getElementById('main-scroll-container')
    if (el) el.scrollTop = 0
  }, 30)
}

const openSettings = () => {
  store.currentView = 'settings'
  store.currentAnswer = null
}

const getHexColor = (colorKey) => {
  if (!colorKey) return '#3b82f6';
  if (colorKey.startsWith('#')) return colorKey;

  const hexColors = {
    blue: '#3b82f6', red: '#ef4444', green: '#10b981', yellow: '#f59e0b',
    purple: '#8b5cf6', pink: '#ec4899', indigo: '#6366f1', teal: '#14b8a6',
    orange: '#f97316', cyan: '#06b6d4', slate: '#64748b'
  }
  return hexColors[colorKey] || hexColors.blue;
}
</script>
