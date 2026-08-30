<template>
  <div class="select-none md:h-full w-full h-full pb-4">
    <!-- 空状态 -->
    <div v-if="store.groupedAnswers.length === 0" class="h-64 flex flex-col items-center justify-center text-gray-400 space-y-2">
      <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path></svg>
      <span class="text-sm">暂无匹配的归档内容</span>
    </div>

    <!-- 虚拟滚动容器 -->
    <div
      v-else
      class="relative w-full"
      :style="{ height: `${virtualizer.getTotalSize()}px` }"
    >
      <!-- 网格视图 (Grid) -->
      <template v-if="store.viewMode === 'grid'">
        <div
          v-for="virtualRow in virtualizer.getVirtualItems()"
          :key="virtualRow.key"
          :style="{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: `${virtualRow.size}px`,
            transform: `translateY(${virtualRow.start}px)`,
            display: 'grid',
            gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))`
          }"
          class="gap-4 pb-4"
        >
          <div
            v-for="group in gridRows[virtualRow.index]"
            :key="group.groupKey"
            @click="handleGroupClick(group)"
            :class="['group bg-surface border border-line/80 rounded-xl px-4 pt-3 pb-1 hover:shadow-md transition cursor-pointer flex flex-col justify-between h-full relative overflow-hidden gap-2', isTwitter(group) ? 'hover:border-primary' : 'hover:border-blue-400']"
          >
            <div class="space-y-1 flex-1">
              <div class="flex items-center justify-between w-full">
                <div class="flex items-center space-x-2">
                  <PlatformIcon :questionId="group.answers[0].question_id" svgClass="w-5 h-5 transition-transform group-hover:scale-105" />
                  <AppBadge size="xs">{{ group.answers[0].tag || 'ANSWER' }}</AppBadge>
                </div>
                <div class="flex items-center space-x-1 -mr-1">
                  <button
                    @click.stop="store.toggleFavorite(group.count === 1 ? group.answers[0] : group, group.count > 1)"
                    :class="['transition cursor-pointer p-1 rounded-full hover:bg-surface-hover', (group.count === 1 ? group.answers[0].is_favorite : group.answers[0].is_favorite) ? 'text-yellow-400 hover:text-yellow-500' : 'text-gray-300 hover:text-gray-400']"
                    title="收藏/取消收藏"
                  >
                    <svg class="w-4 h-4" :fill="(group.count === 1 ? group.answers[0].is_favorite : group.answers[0].is_favorite) ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
                    </svg>
                  </button>
                  <button
                    @click.stop="store.itemToDelete = group.count === 1 ? group.answers[0] : group; store.itemToDeleteType = group.count === 1 ? 'answer' : 'group'"
                    class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50 dark:hover:bg-red-900/30"
                    title="删除归档"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </div>
              </div>
              <h3 :class="['font-medium text-primary text-sm line-clamp-3 transition leading-snug h-[56px]', isTwitter(group) ? 'group-hover:text-primary' : 'group-hover:text-brand']">
                {{ group.title || '（无标题）' }}
              </h3>
            </div>
            <div class="mt-auto pt-1 border-t border-line-light flex items-center justify-between text-[11px] text-gray-400">
              <span v-if="group.count > 1" class="inline-flex items-center space-x-1">
                <span class="w-3.5 h-3.5 rounded-full bg-blue-500 text-white text-[9px] font-bold flex items-center justify-center leading-none">{{ group.count }}</span>
              </span>
              <span v-else></span>
              <span>归档：{{ formatDate(group.answers[0].saved_at) }}</span>
            </div>
          </div>
        </div>
      </template>

      <!-- 列表视图 (List) -->
      <template v-else>
        <div class="w-full absolute top-0 left-0 bg-surface border border-line/80 rounded-xl overflow-hidden" :style="{ height: `${virtualizer.getTotalSize()}px` }">
          <div
            v-for="virtualRow in virtualizer.getVirtualItems()"
            :key="virtualRow.key"
            :data-index="virtualRow.index"
            :ref="virtualizer.measureElement"
            :style="{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              transform: `translateY(${virtualRow.start}px)`
            }"
            @click="handleGroupClick(store.groupedAnswers[virtualRow.index])"
            :class="['px-3 sm:px-4 py-3 transition cursor-pointer flex items-center justify-between group border-b border-line-light last:border-b-0', isTwitter(store.groupedAnswers[virtualRow.index]) ? 'hover:bg-surface-hover' : 'hover:bg-blue-50 dark:hover:bg-blue-900/30/50 dark:hover:bg-blue-900/20']"
          >
            <div class="flex items-center space-x-3 pr-2 sm:pr-4 flex-1 overflow-hidden">
              <PlatformIcon :questionId="store.groupedAnswers[virtualRow.index].answers[0].question_id" svgClass="w-4 h-4 hidden sm:block" />
              <div class="flex flex-col flex-1 overflow-hidden">
                <span :class="['text-sm font-medium text-primary line-clamp-2 sm:truncate', isTwitter(store.groupedAnswers[virtualRow.index]) ? 'group-hover:text-primary' : 'group-hover:text-brand']">
                  {{ store.groupedAnswers[virtualRow.index].title || '（无标题）' }}
                </span>
              </div>
            </div>
              <div class="flex items-center space-x-2 sm:space-x-4 ml-2 sm:ml-6 -mr-1">
                <button
                  @click.stop="store.toggleFavorite(store.groupedAnswers[virtualRow.index].count === 1 ? store.groupedAnswers[virtualRow.index].answers[0] : store.groupedAnswers[virtualRow.index], store.groupedAnswers[virtualRow.index].count > 1)"
                  :class="['transition cursor-pointer p-1 rounded-full hover:bg-surface-hover', (store.groupedAnswers[virtualRow.index].count === 1 ? store.groupedAnswers[virtualRow.index].answers[0].is_favorite : store.groupedAnswers[virtualRow.index].answers[0].is_favorite) ? 'text-yellow-400 hover:text-yellow-500' : 'text-gray-300 hover:text-gray-400']"
                  title="收藏/取消收藏"
                >
                  <svg class="w-4 h-4" :fill="(store.groupedAnswers[virtualRow.index].count === 1 ? store.groupedAnswers[virtualRow.index].answers[0].is_favorite : store.groupedAnswers[virtualRow.index].answers[0].is_favorite) ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path>
                  </svg>
                </button>
                <button
                  @click.stop="store.itemToDelete = store.groupedAnswers[virtualRow.index].count === 1 ? store.groupedAnswers[virtualRow.index].answers[0] : store.groupedAnswers[virtualRow.index]; store.itemToDeleteType = store.groupedAnswers[virtualRow.index].count === 1 ? 'answer' : 'group'"
                  class="text-gray-300 hover:text-red-500 transition cursor-pointer p-1 rounded-full hover:bg-red-50 dark:hover:bg-red-900/30"
                  title="删除归档"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                </button>
              </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 无限滚动触发哨兵 -->
    <div ref="sentinel" class="h-8 flex items-center justify-center mt-6">
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useVirtualizer } from '@tanstack/vue-virtual'
import { useArchiveStore } from '../../stores/archive'
import PlatformIcon from '../common/PlatformIcon.vue'
import AppBadge from '../common/AppBadge.vue'

const store = useArchiveStore()
const sentinel = ref(null)
let observer = null

// --- 网格列数计算 ---
const columns = ref(1)
const scrollElement = ref(null)
let resizeObserver = null

const updateColumns = () => {
  const width = scrollElement.value?.clientWidth || window.innerWidth
  // Based on a min-width of roughly 300px per card + 16px gap
  columns.value = Math.max(1, Math.floor((width + 16) / 316))
}

// 分块逻辑：将一维数组切分为网格的“行”
const gridRows = computed(() => {
  if (store.viewMode !== 'grid') return []
  const rows = []
  const items = store.groupedAnswers
  const cols = columns.value
  for (let i = 0; i < items.length; i += cols) {
    rows.push(items.slice(i, i + cols))
  }
  return rows
})

// --- 虚拟列表配置 ---
const virtualItemsCount = computed(() => {
  return store.viewMode === 'grid' ? gridRows.value.length : store.groupedAnswers.length
})

const virtualizerOptions = computed(() => ({
  count: virtualItemsCount.value,
  getScrollElement: () => scrollElement.value,
  estimateSize: () => store.viewMode === 'grid' ? 166 : 70, // 网格行高 150+16gap = 166, 列表 70
  overscan: 5,
}))

const virtualizer = useVirtualizer(virtualizerOptions)

onMounted(() => {
  scrollElement.value = document.getElementById('main-scroll-container')
  
  if (scrollElement.value) {
    resizeObserver = new ResizeObserver(() => {
      updateColumns()
      if (store.viewMode === 'list' && virtualizer.value) {
         virtualizer.value.measure()
      }
    })
    resizeObserver.observe(scrollElement.value)
    updateColumns()
  }

  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && store.hasMore && !store.isLoadingMore) {
      store.loadMoreAnswers()
    }
  }, { threshold: 0.1, rootMargin: '200px' })
  
  if (sentinel.value) observer.observe(sentinel.value)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
  if (resizeObserver && scrollElement.value) {
    resizeObserver.disconnect()
  }
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
