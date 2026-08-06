<template>
  <div class="select-none md:h-full">
    <!-- 状态提示：暂无数据或无搜索结果 -->
    <div v-if="store.filteredAnswers.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-400 space-y-2">
      <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path></svg>
      <span class="text-sm">暂无匹配的归档内容</span>
    </div>

    <!-- 网格视图 (Grid) -->
    <div v-if="store.viewMode === 'grid'" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div 
        v-for="item in store.filteredAnswers" 
        :key="item.answer_id"
        @click="store.selectAnswer(item.answer_id)"
        :class="['group bg-white border border-gray-200/80 rounded-xl px-4 pt-4 pb-1 hover:shadow-md transition cursor-pointer flex flex-col justify-between h-full relative overflow-hidden gap-2', item.question_idß === 'twitter' ? 'hover:border-gray-800' : 'hover:border-blue-400']"
      >
        <div class="space-y-1 flex-1">
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center space-x-2">
              <PlatformIcon :questionId="item.question_id" svgClass="w-5 h-5 transition-transform group-hover:scale-105" />
              <span class="text-[10px] font-semibold tracking-wider uppercase text-gray-400">{{ item.tag || 'ANSWER' }}</span>
            </div>
            <button 
              @click.stop="store.itemToDelete = item"
              class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50"
              title="删除归档"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
          <h3 :class="['font-medium text-gray-800 text-xs line-clamp-3 transition leading-5 h-[60px]', item.question_id === 'twitter' ? 'group-hover:text-gray-900' : 'group-hover:text-blue-600']">
            {{ item.title }}
          </h3>
        </div>

        <div class="mt-auto pt-1 border-t border-gray-100 flex items-center justify-end text-[11px] text-gray-400">
          <span>归档：{{ formatDate(item.saved_at) }}</span>
        </div>
      </div>
    </div>

    <!-- 列表视图 (List) -->
    <div v-else class="divide-y divide-gray-100 border border-gray-200/80 rounded-xl overflow-hidden">
      <div 
        v-for="item in store.filteredAnswers" 
        :key="item.answer_id"
        @click="store.selectAnswer(item.answer_id)"
        :class="['px-3 sm:px-4 py-3 transition cursor-pointer flex items-center justify-between group', item.question_id === 'twitter' ? 'hover:bg-gray-100' : 'hover:bg-blue-50/50']"
      >
        <div class="flex items-center space-x-3 pr-2 sm:pr-4 flex-1 overflow-hidden">
          <PlatformIcon :questionId="item.question_id" svgClass="w-4 h-4 hidden sm:block grayscale group-hover:grayscale-0 opacity-70 group-hover:opacity-100 transition-all" />
          <div class="flex flex-col flex-1 overflow-hidden">
            <span :class="['text-xs font-medium text-gray-800 line-clamp-2 sm:truncate sm:mb-0 mb-1', item.question_id === 'twitter' ? 'group-hover:text-gray-900' : 'group-hover:text-blue-600']">{{ item.title }}</span>
            <div class="flex sm:hidden items-center text-[10px] text-gray-400">
              <span>归档：{{ formatDate(item.saved_at) }}</span>
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-2 sm:space-x-6 text-[11px] text-gray-400 flex-shrink-0">
          <span class="hidden sm:inline">归档：{{ formatDate(item.saved_at) }}</span>
          <button 
            @click.stop="store.itemToDelete = item"
            class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50"
            title="删除归档"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'
import PlatformIcon from '../common/PlatformIcon.vue'

const store = useArchiveStore()

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}

const getHexColor = (colorKey) => {
  const hexColors = {
    blue: '#3b82f6', red: '#ef4444', green: '#10b981', yellow: '#f59e0b',
    purple: '#8b5cf6', pink: '#ec4899', indigo: '#6366f1', teal: '#14b8a6',
    orange: '#f97316', cyan: '#06b6d4', slate: '#64748b'
  }
  return hexColors[colorKey] || hexColors.blue;
}
</script>
