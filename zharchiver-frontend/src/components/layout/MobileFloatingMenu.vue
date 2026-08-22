<template>
  <div class="md:hidden fixed bottom-[calc(1.5rem+env(safe-area-inset-bottom))] right-4 z-[60]">
    <!-- 展开菜单 -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 translate-y-4 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-4 scale-95"
    >
      <div v-if="isOpen" class="absolute bottom-16 right-0 w-48 bg-surface/90 backdrop-blur-md border border-line rounded-2xl shadow-xl overflow-hidden flex flex-col py-2">
        <div class="px-4 py-2 text-xs font-semibold text-gray-400 tracking-wider">
          系统设置
        </div>
        <button 
          @click="openSettings('security')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'security' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >账户安全</button>
        <button 
          @click="openSettings('auth')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'auth' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >平台鉴权</button>
        <button 
          @click="openSettings('data')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'data' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >数据管理</button>
        <button 
          @click="openSettings('display')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'display' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >偏好设置</button>
        <button 
          @click="openSettings('ai')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'ai' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >工具</button>
        <button 
          @click="openSettings('logs')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', store.activeSetting === 'logs' && store.currentView === 'settings' ? 'bg-surface-hover text-primary' : 'text-secondary hover:bg-surface-hover']"
        >日志中心</button>

        <div v-if="store.currentView === 'settings'" class="border-t border-line-light mt-1 pt-1">
          <button 
            @click="goHome"
            class="w-full text-left px-4 py-2 text-sm font-medium text-brand hover:bg-blue-50 dark:hover:bg-blue-900/30 transition cursor-pointer"
          >返回主界面</button>
        </div>
      </div>
    </transition>

    <!-- 悬浮按钮 (FAB) -->
    <button 
      @click="toggleMenu"
      :class="['w-12 h-12 rounded-full flex items-center justify-center text-white shadow-lg transition-transform duration-300', isOpen ? 'bg-gray-800 rotate-45' : 'bg-brand hover:bg-brand-hover']"
    >
      <svg v-if="!isOpen" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
      <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
    </button>

    <!-- 点击空白处关闭 -->
    <div v-if="isOpen" @click="isOpen = false" class="fixed inset-0 z-[-1]"></div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()
const isOpen = ref(false)

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const openSettings = (setting) => {
  store.currentView = 'settings'
  store.activeSetting = setting
  isOpen.value = false
}

const goHome = () => {
  store.currentView = 'home'
  isOpen.value = false
}
</script>
