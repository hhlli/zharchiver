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
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'grid' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-900']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
            <span>网格</span>
          </button>
          <button 
            @click="updateViewMode('list')"
            :class="['px-4 py-1.5 rounded-md transition-all cursor-pointer text-xs font-medium flex items-center space-x-1.5', store.viewMode === 'list' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-500 hover:text-gray-900']"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            <span>列表</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()

const updateViewMode = (mode) => {
  store.setViewMode(mode)
}
</script>
