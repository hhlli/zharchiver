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
      <div v-if="isOpen" class="absolute bottom-16 right-0 w-48 bg-white/90 backdrop-blur-md border border-gray-200 rounded-2xl shadow-xl overflow-hidden flex flex-col py-2">
        <div class="px-4 py-2 text-xs font-semibold text-gray-400 tracking-wider">
          系统设置
        </div>
        <button 
          @click="openSettings('security')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'security' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >账户安全</button>
        <button 
          @click="openSettings('auth')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'auth' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >平台鉴权</button>
        <button 
          @click="openSettings('data')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'data' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >数据管理</button>
        <button 
          @click="openSettings('display')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'display' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >偏好设置</button>
        <button 
          @click="openSettings('ai')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'ai' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >工具</button>
        <button 
          @click="openSettings('logs')"
          :class="['w-full text-left px-4 py-2 text-sm font-medium transition cursor-pointer', activeSetting === 'logs' && currentView === 'settings' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >日志中心</button>

        <div v-if="currentView === 'settings'" class="border-t border-gray-100 mt-1 pt-1">
          <button 
            @click="goHome"
            class="w-full text-left px-4 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50 transition cursor-pointer"
          >返回主界面</button>
        </div>
      </div>
    </transition>

    <!-- 悬浮按钮 (FAB) -->
    <button 
      @click="toggleMenu"
      :class="['w-12 h-12 rounded-full flex items-center justify-center text-white shadow-lg transition-transform duration-300', isOpen ? 'bg-gray-800 rotate-45' : 'bg-blue-600 hover:bg-blue-700']"
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

const props = defineProps({
  currentView: {
    type: String,
    required: true
  },
  activeSetting: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:currentView', 'update:activeSetting'])

const isOpen = ref(false)

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const openSettings = (setting) => {
  emit('update:currentView', 'settings')
  emit('update:activeSetting', setting)
  isOpen.value = false
}

const goHome = () => {
  emit('update:currentView', 'home')
  isOpen.value = false
}
</script>
