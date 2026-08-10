<template>
  <div class="select-none md:h-full">
    <!-- 空状态 -->
    <div v-if="store.groupedAnswers.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-400 space-y-2">
      <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path></svg>
      <span class="text-sm">暂无匹配的归档内容</span>
    </div>

    <!-- 网格视图 (Grid) -->
    <div v-if="store.viewMode === 'grid'" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div
        v-for="group in store.groupedAnswers"
        :key="group.groupKey"
        @click="handleGroupClick(group)"
        :class="['group bg-white border border-gray-200/80 rounded-xl px-4 pt-3 pb-1 hover:shadow-md transition cursor-pointer flex flex-col justify-between h-full relative overflow-hidden gap-2', isTwitter(group) ? 'hover:border-black' : 'hover:border-blue-400']"
      >
        <div class="space-y-1 flex-1">
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center space-x-2">
              <PlatformIcon :questionId="group.answers[0].question_id" svgClass="w-5 h-5 transition-transform group-hover:scale-105" />
              <AppBadge size="xs">{{ group.answers[0].tag || 'ANSWER' }}</AppBadge>
            </div>
            <button
              v-if="group.count === 1"
              @click.stop="store.itemToDelete = group.answers[0]"
              class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50"
              title="删除归档"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
            </button>
          </div>
          <h3 :class="['font-medium text-gray-800 text-sm line-clamp-3 transition leading-snug h-[56px]', isTwitter(group) ? 'group-hover:text-gray-900' : 'group-hover:text-brand']">
            {{ group.title || '（无标题）' }}
          </h3>
        </div>
        <div class="mt-auto pt-1 border-t border-gray-100 flex items-center justify-between text-[11px] text-gray-400">
          <!-- 多条回答角标 -->
          <span v-if="group.count > 1" class="inline-flex items-center space-x-1">
            <span class="w-3.5 h-3.5 rounded-full bg-blue-500 text-white text-[9px] font-bold flex items-center justify-center leading-none">{{ group.count }}</span>
          </span>
          <span v-else></span>
          <span>归档：{{ formatDate(group.answers[0].saved_at) }}</span>
        </div>
      </div>
    </div>

    <!-- 列表视图 (List) -->
    <div v-else class="divide-y divide-gray-100 border border-gray-200/80 rounded-xl overflow-hidden">
      <div
        v-for="group in store.groupedAnswers"
        :key="group.groupKey"
        @click="handleGroupClick(group)"
        :class="['px-3 sm:px-4 py-3 transition cursor-pointer flex items-center justify-between group', isTwitter(group) ? 'hover:bg-gray-100' : 'hover:bg-blue-50/50']"
      >
        <div class="flex items-center space-x-3 pr-2 sm:pr-4 flex-1 overflow-hidden">
          <PlatformIcon :questionId="group.answers[0].question_id" svgClass="w-4 h-4 hidden sm:block grayscale group-hover:grayscale-0 opacity-70 group-hover:opacity-100 transition-all" />
          <div class="flex flex-col flex-1 overflow-hidden">
            <span :class="['text-sm font-medium text-gray-800 line-clamp-2 sm:truncate sm:mb-0 mb-1', isTwitter(group) ? 'group-hover:text-gray-900' : 'group-hover:text-brand']">
              {{ group.title || '（无标题）' }}
            </span>
            <div class="flex sm:hidden items-center text-xs text-gray-400">
              <span>归档：{{ formatDate(group.answers[0].saved_at) }}</span>
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-2 sm:space-x-6 text-xs text-gray-400 flex-shrink-0">
          <span class="hidden sm:inline">归档：{{ formatDate(group.answers[0].saved_at) }}</span>
          <button
            v-if="group.count === 1"
            @click.stop="store.itemToDelete = group.answers[0]"
            class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50"
            title="删除归档"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 无限滚动触发哨兵 -->
    <div ref="sentinel" class="h-8 flex items-center justify-center mt-2">
      <span v-if="store.isLoadingMore" class="text-xs text-gray-400 flex items-center space-x-1.5">
        <svg class="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        <span>加载中...</span>
      </span>
      <span v-else-if="!store.hasMore && store.totalCount > 0" class="text-xs text-gray-300">
        已加载全部 {{ store.totalCount }} 条
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import PlatformIcon from '../common/PlatformIcon.vue'
import AppBadge from '../common/AppBadge.vue'

const store = useArchiveStore()
const sentinel = ref(null)
let observer = null

onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && store.hasMore && !store.isLoadingMore) {
      store.loadMoreAnswers()
    }
  }, { threshold: 0.1 })
  if (sentinel.value) observer.observe(sentinel.value)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})

const isTwitter = (group) => group.answers[0]?.question_id === 'twitter'

const handleGroupClick = (group) => {
  if (group.count === 1) {
    store.selectAnswer(group.answers[0].answer_id)
  } else {
    store.selectGroup(group)
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString()
}
</script>
