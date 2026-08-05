<template>
  <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-3 flex-shrink-0 select-none">
    <!-- 左侧：Logo 与 基础控制 -->
    <div class="flex items-center space-x-3 md:space-x-6">
      <div class="flex items-center space-x-2">
        <!-- Logo: 使用一个更精致的SVG图标代表归档/箱子 -->
        <div class="w-7 h-7 bg-blue-600 rounded-lg flex items-center justify-center shadow-sm">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
          </svg>
        </div>
        <span class="text-[15px] font-bold text-gray-800 tracking-tight">ZHArchiver</span>
      </div>
    </div>

    <!-- 右侧：搜索框 与 新建归档 -->
    <div class="flex items-center space-x-2 md:space-x-4">
      <div class="w-24 sm:w-32 md:w-52 relative">
        <input 
          :value="searchQuery"
          @input="$emit('update:searchQuery', $event.target.value)"
          type="text" 
          placeholder="搜索回答..." 
          class="w-full bg-gray-200/60 text-xs text-gray-800 pl-7 pr-3 py-1.5 rounded-lg border border-transparent focus:bg-white focus:border-blue-500 focus:outline-none transition"
        />
        <svg class="w-3.5 h-3.5 text-gray-400 absolute left-2 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
        </svg>
      </div>

      <!-- 分类下拉框 (仅移动端) -->
      <select 
        v-if="currentView === 'home'"
        :value="activeCategory"
        @change="$emit('update:activeCategory', $event.target.value)"
        class="md:hidden w-20 sm:w-24 bg-gray-100 text-xs text-gray-700 py-1.5 px-2 rounded-lg border-none focus:ring-0 appearance-none text-center cursor-pointer"
      >
        <option value="all">所有</option>
        <option v-for="tag in tags" :key="tag.name" :value="tag.name">{{ tag.name }}</option>
      </select>
      
      <!-- 新增归档按钮 -->
      <button 
        @click="$emit('add-archive')"
        class="px-2 md:px-3 py-1.5 bg-blue-600 text-white rounded-lg text-xs font-medium hover:bg-blue-700 transition flex items-center md:space-x-1 cursor-pointer shadow-sm"
      >
        <svg class="w-4 h-4 md:w-3.5 md:h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        <span class="hidden md:inline">新建归档</span>
      </button>
    </div>
  </header>
</template>

<script setup>
defineProps({
  searchQuery: {
    type: String,
    default: ''
  },
  activeCategory: {
    type: String,
    default: 'all'
  },
  tags: {
    type: Array,
    default: () => []
  },
  currentView: {
    type: String,
    default: 'home'
  }
})
defineEmits(['update:searchQuery', 'update:activeCategory', 'add-archive'])
</script>
