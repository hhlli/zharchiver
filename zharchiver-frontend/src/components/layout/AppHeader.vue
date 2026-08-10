<template>
  <header class="bg-white border-b border-gray-200 flex-shrink-0 select-none pt-[env(safe-area-inset-top)]">
    <div class="h-14 flex items-center justify-between px-3">
      <!-- 左侧：Logo 与 基础控制 -->
      <div class="flex items-center space-x-3 md:space-x-6">
      <!-- Logo + 品牌名 -->
      <div class="flex items-center space-x-2">
        <div class="w-7 h-7 bg-brand rounded-lg flex items-center justify-center shadow-sm">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
          </svg>
        </div>
        <span class="text-tweet font-bold text-gray-800 tracking-tight">ZHArchiver</span>
      </div>
    </div>

    <!-- 右侧：搜索框 与 新建归档 -->
    <div class="flex items-center space-x-2 md:space-x-4">
      <div v-if="store.currentView === 'home'" class="w-24 sm:w-32 md:w-52 relative">
        <input 
          v-model="store.searchQuery"
          type="text" 
          placeholder="搜索回答..." 
          class="w-full h-8 bg-gray-200/60 text-xs text-gray-800 pl-7 pr-3 rounded-[10px] border border-transparent focus:bg-white focus:border-blue-500 focus:outline-none transition appearance-none"
        />
        <svg class="w-3.5 h-3.5 text-gray-400 absolute left-2 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
        </svg>
      </div>

      <!-- 分类下拉框 (仅移动端) -->
      <TagSelectDropdown
        v-if="store.currentView === 'home'"
        v-model="store.activeCategory"
        :tags="store.tags"
      />
      
      <!-- 新增归档按钮 -->
      <button 
        @click="store.showAddModal = true"
        class="w-8 h-8 md:w-auto md:px-3 bg-brand text-white rounded-[10px] text-xs font-medium hover:bg-brand-hover transition flex items-center justify-center md:space-x-1 cursor-pointer shadow-sm flex-shrink-0"
      >
        <svg class="w-4 h-4 md:w-3.5 md:h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        <span class="hidden md:inline">新建归档</span>
      </button>
    </div>
    </div>
  </header>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'
import TagSelectDropdown from '../common/TagSelectDropdown.vue'

const store = useArchiveStore()
</script>
