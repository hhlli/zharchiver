<template>
  <header class="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-3 flex-shrink-0 select-none">
    <!-- 左侧：Logo 与 基础控制 -->
    <div class="flex items-center space-x-3 md:space-x-6">
      <div class="flex items-center space-x-2">
        <!-- 移动端汉堡菜单按钮 -->
        <button 
          @click="$emit('toggle-sidebar')"
          class="md:hidden p-1 -ml-2 text-gray-500 hover:text-gray-800 transition rounded-md hover:bg-gray-100"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
        </button>
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
      <div class="w-36 md:w-52 relative">
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
  }
})
defineEmits(['update:searchQuery', 'add-archive', 'toggle-sidebar'])
</script>
