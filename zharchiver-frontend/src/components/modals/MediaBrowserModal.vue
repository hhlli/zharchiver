<template>
  <div class="absolute inset-0 z-50 bg-surface-hover flex flex-col items-center justify-center">
    <!-- Header -->
    <div class="absolute top-0 w-full h-16 flex items-center justify-between px-6 bg-surface/90 backdrop-blur border-b border-line z-10">
      <div class="text-primary font-medium text-lg tracking-wide flex items-center space-x-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
        <span>视图数据浏览</span>
        <span class="text-sm text-muted ml-2 font-normal" v-if="!loading">{{ mediaList.length }} 个资源 ({{ mediaGroups.length }} 个归档)</span>
      </div>
      <button 
        @click="$emit('close')" 
        class="w-8 h-8 rounded-full hover:bg-line flex items-center justify-center text-secondary transition-colors cursor-pointer"
        title="返回"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>

    <!-- Content -->
    <div class="w-full h-full overflow-y-auto px-4 sm:px-8 pt-20 pb-12 custom-scrollbar">
      
      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center justify-center h-full text-muted space-y-4">
        <svg class="animate-spin w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>正在加载资源库...</span>
      </div>

      <!-- Empty -->
      <div v-else-if="mediaGroups.length === 0" class="flex flex-col items-center justify-center h-full text-gray-400 space-y-3">
        <svg class="w-12 h-12 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
        <p>暂无任何归档的图片或视频</p>
      </div>

      <!-- Masonry Grid -->
      <div v-else class="columns-2 md:columns-3 lg:columns-4 xl:columns-5 gap-4 space-y-4">
        <div 
          v-for="(group, idx) in displayedMediaGroups" 
          :key="idx"
          @click="openGroup(group)"
          class="relative group break-inside-avoid overflow-hidden rounded-xl bg-surface shadow-sm border border-line transition-transform duration-300 hover:shadow-md hover:-translate-y-1 cursor-pointer"
        >
          <!-- Badge for multiple items -->
          <div v-if="group.items.length > 1" class="absolute top-2 right-2 bg-black/60 backdrop-blur-md text-white text-xs px-2 py-1 rounded-md flex items-center gap-1 z-10 font-medium shadow-sm">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
            {{ group.items.length }}
          </div>

          <!-- Cover Image -->
          <img 
            v-if="group.cover.type === 'image'" 
            :src="getMediaUrl(group.cover.url)" 
            loading="lazy" 
            class="w-full h-auto object-cover pointer-events-none"
          />
          <!-- Cover Video -->
          <div v-else-if="group.cover.type === 'video'" class="relative">
            <video 
              preload="metadata"
              class="w-full h-auto object-cover pointer-events-none"
            >
              <source :src="getMediaUrl(group.cover.url)" />
            </video>
            <!-- Play Icon Overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/10 group-hover:bg-black/20 transition-colors">
              <svg class="w-12 h-12 text-white/80" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
            </div>
          </div>

          <!-- Overlay Info -->
          <div class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/90 via-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none">
            <div class="text-white/90 text-sm font-medium truncate">归档于 {{ formatDate(group.saved_at) }}</div>
            <div class="text-white/50 text-xs mt-0.5 uppercase">{{ group.cover.type }}</div>
          </div>
        </div>
      </div>
      
      <!-- Load More Sentinel -->
      <div ref="loadMoreSentinel" class="h-10 mt-4"></div>

    </div>

    <!-- Lightbox Overlay -->
    <div 
      v-if="selectedGroup" 
      class="fixed inset-0 z-[200] bg-black/95 flex items-center justify-center backdrop-blur-md"
      @click="closeLightbox"
    >
      <div class="absolute top-6 right-6 flex items-center space-x-3 z-20">
        <!-- Index Indicator -->
        <div v-if="selectedGroup.items.length > 1" class="h-12 px-4 rounded-full bg-surface/10 text-white flex items-center font-medium text-sm">
          {{ selectedIndex + 1 }} / {{ selectedGroup.items.length }}
        </div>
        <!-- Jump to Article Button -->
        <button 
          @click.stop="jumpToArticle"
          class="h-12 px-5 rounded-full bg-brand hover:bg-brand-hover text-white flex items-center space-x-2 font-medium transition-colors shadow-lg cursor-pointer"
          title="跳转至原文"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
          <span>查看原文</span>
        </button>
        <!-- Close Button -->
        <button 
          class="w-12 h-12 rounded-full bg-surface/10 hover:bg-surface/20 flex items-center justify-center text-white transition-colors cursor-pointer"
          @click.stop="closeLightbox"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <!-- Prev Button -->
      <button 
        v-if="selectedGroup.items.length > 1"
        @click.stop="prevMedia"
        class="absolute left-6 top-1/2 -translate-y-1/2 w-12 h-12 rounded-full bg-surface/10 hover:bg-surface/20 flex items-center justify-center text-white transition-colors z-20 cursor-pointer"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
      </button>

      <!-- Next Button -->
      <button 
        v-if="selectedGroup.items.length > 1"
        @click.stop="nextMedia"
        class="absolute right-6 top-1/2 -translate-y-1/2 w-12 h-12 rounded-full bg-surface/10 hover:bg-surface/20 flex items-center justify-center text-white transition-colors z-20 cursor-pointer"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
      </button>

      <!-- Current Media -->
      <img 
        v-if="currentMedia.type === 'image'" 
        :key="currentMedia.url"
        :src="getMediaUrl(currentMedia.url)" 
        class="max-w-full max-h-full object-contain select-none transition-opacity"
        @click.stop
      />
      <video 
        v-else-if="currentMedia.type === 'video'" 
        :key="currentMedia.url"
        controls 
        autoplay
        class="max-w-[90%] max-h-[90%] object-contain rounded-lg shadow-2xl bg-black"
        @click.stop
      >
        <source :src="getMediaUrl(currentMedia.url)" />
      </video>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()
const emit = defineEmits(['close'])

const loading = ref(true)
const mediaList = ref([])

const selectedGroup = ref(null)
const selectedIndex = ref(0)
const currentMedia = computed(() => {
  if (!selectedGroup.value) return null
  return selectedGroup.value.items[selectedIndex.value]
})

const mediaGroups = computed(() => {
  const groups = []
  const map = {}
  for (const item of mediaList.value) {
    if (!map[item.answer_id]) {
      const newGroup = {
        answer_id: item.answer_id,
        saved_at: item.saved_at,
        cover: item,
        items: []
      }
      map[item.answer_id] = newGroup
      groups.push(newGroup)
    }
    map[item.answer_id].items.push(item)
  }
  return groups
})

const page = ref(1)
const pageSize = 30
const displayedMediaGroups = computed(() => {
  return mediaGroups.value.slice(0, page.value * pageSize)
})

const loadMoreSentinel = ref(null)
let observer = null

const openGroup = (group) => {
  selectedGroup.value = group
  selectedIndex.value = 0
}

const closeLightbox = () => {
  selectedGroup.value = null
  selectedIndex.value = 0
}

const prevMedia = () => {
  if (!selectedGroup.value) return
  if (selectedIndex.value > 0) {
    selectedIndex.value--
  } else {
    selectedIndex.value = selectedGroup.value.items.length - 1
  }
}

const nextMedia = () => {
  if (!selectedGroup.value) return
  if (selectedIndex.value < selectedGroup.value.items.length - 1) {
    selectedIndex.value++
  } else {
    selectedIndex.value = 0
  }
}

// Keyboard navigation
const handleKeydown = (e) => {
  if (!selectedGroup.value) return
  if (e.key === 'ArrowLeft') {
    prevMedia()
  } else if (e.key === 'ArrowRight') {
    nextMedia()
  } else if (e.key === 'Escape') {
    closeLightbox()
  }
}

const jumpToArticle = async () => {
  if (selectedGroup.value && selectedGroup.value.answer_id) {
    store.currentView = 'home'
    emit('close')
    // Give time for view transition before selecting answer
    setTimeout(() => {
      store.selectAnswer(selectedGroup.value.answer_id)
    }, 100)
  }
}

onMounted(async () => {
  document.body.style.overflow = 'hidden' // Lock scroll
  window.addEventListener('keydown', handleKeydown)
  try {
    const res = await store.apiFetch('/api/media')
    if (res.ok) {
      mediaList.value = await res.json()
    } else {
      store.showToast('加载媒体库失败', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('网络请求失败', 'error')
  } finally {
    loading.value = false
  }

  // Setup infinite scroll
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
      if (page.value * pageSize < mediaGroups.value.length) {
        page.value++
      }
    }
  }, { rootMargin: '300px' })

  setTimeout(() => {
    if (loadMoreSentinel.value) observer.observe(loadMoreSentinel.value)
  }, 500)
})

onUnmounted(() => {
  document.body.style.overflow = '' // Restore scroll
  window.removeEventListener('keydown', handleKeydown)
  if (observer) observer.disconnect()
})

const formatDate = (dateStr) => {
  if (!dateStr) return '未知时间'
  const d = new Date(dateStr)
  return `${d.getFullYear()}/${d.getMonth()+1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const getMediaUrl = (url) => {
  const token = localStorage.getItem('token')
  const baseUrl = store.API_BASE + url
  if (token) {
    return `${baseUrl}?token=${token}`
  }
  return baseUrl
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.2);
}
</style>
