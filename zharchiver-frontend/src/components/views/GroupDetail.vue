<template>
  <div class="md:h-full flex flex-col">
    <!-- 顶部返回栏 -->
    <div class="flex items-center space-x-3 mb-4">
      <button
        @click="store.currentGroup = null"
        class="flex items-center space-x-1.5 text-sm text-gray-500 hover:text-blue-600 transition cursor-pointer group"
      >
        <svg class="w-4 h-4 transition-transform group-hover:-translate-x-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        <span>返回</span>
      </button>
    </div>

    <!-- 问题标题 -->
    <h2 class="text-base font-semibold text-gray-800 mb-4 leading-snug">
      {{ store.currentGroup.title }}
    </h2>

    <!-- 回答列表 -->
    <div class="divide-y divide-gray-100 border border-gray-200/80 rounded-xl overflow-hidden">
      <div
        v-for="item in store.currentGroup.answers"
        :key="item.answer_id"
        @click="store.selectAnswer(item.answer_id)"
        class="px-4 py-3 flex items-center justify-between hover:bg-blue-50/50 cursor-pointer transition group"
      >
        <div class="flex flex-col flex-1 overflow-hidden pr-4">
          <span class="text-sm font-medium text-gray-800 group-hover:text-brand transition truncate">
            {{ item.author_name || '匿名用户' }}
          </span>
          <span class="text-xs text-gray-400 mt-0.5">归档于 {{ formatDate(item.saved_at) }}</span>
        </div>
        <div class="flex items-center space-x-3 flex-shrink-0">
          <button
            @click.stop="store.itemToDelete = item"
            class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50"
            title="删除此回答"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
          <svg class="w-4 h-4 text-gray-300 group-hover:text-blue-400 transition" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}
</script>
