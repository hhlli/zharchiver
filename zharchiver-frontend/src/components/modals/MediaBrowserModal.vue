<template>
  <div class="absolute inset-0 z-50 bg-gray-50 flex flex-col items-center justify-center">
    <!-- Header -->
    <div class="absolute top-0 w-full h-16 flex items-center justify-between px-6 bg-white/90 backdrop-blur border-b border-gray-200 z-10">
      <div class="text-gray-800 font-medium text-lg tracking-wide flex items-center space-x-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
        <span>视图数据浏览</span>
        <span class="text-sm text-gray-500 ml-2 font-normal" v-if="!loading">{{ mediaList.length }} 个资源</span>
      </div>
      <button 
        @click="$emit('close')" 
        class="w-8 h-8 rounded-full hover:bg-gray-200 flex items-center justify-center text-gray-600 transition-colors cursor-pointer"
        title="返回"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
      </button>
    </div>

    <!-- Content -->
    <div class="w-full h-full overflow-y-auto px-4 sm:px-8 pt-20 pb-12 custom-scrollbar">
      
      <!-- Loading -->
      <div v-if="loading" class="flex flex-col items-center justify-center h-full text-gray-500 space-y-4">
        <svg class="animate-spin w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>正在加载资源库...</span>
      </div>

      <!-- Empty -->
      <div v-else-if="mediaList.length === 0" class="flex flex-col items-center justify-center h-full text-gray-400 space-y-3">
        <svg class="w-12 h-12 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
        <p>暂无任何归档的图片或视频</p>
      </div>

      <!-- Masonry Grid -->
      <div v-else class="columns-2 md:columns-3 lg:columns-4 xl:columns-5 gap-4 space-y-4">
        <div 
          v-for="(item, idx) in displayedMediaList" 
          :key="idx"
          @click="selectedMedia = item"
          class="relative group break-inside-avoid overflow-hidden rounded-xl bg-white shadow-sm border border-gray-200 transition-transform duration-300 hover:shadow-md hover:-translate-y-1 cursor-pointer"
        >
          <!-- Image -->
          <img 
            v-if="item.type === 'image'" 
            :src="getMediaUrl(item.url)" 
            loading="lazy" 
            class="w-full h-auto object-cover pointer-events-none"
          />
          <!-- Video -->
          <div v-else-if="item.type === 'video'" class="relative">
            <video 
              preload="metadata"
              class="w-full h-auto object-cover pointer-events-none"
            >
              <source :src="getMediaUrl(item.url)" />
            </video>
            <!-- Play Icon Overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/10 group-hover:bg-black/20 transition-colors">
              <svg class="w-12 h-12 text-white/80" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
            </div>
          </div>

          <!-- Overlay Info -->
          <div class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/90 via-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none">
            <div class="text-white/90 text-xs font-medium truncate">归档于 {{ formatDate(item.saved_at) }}</div>
            <div class="text-white/50 text-[10px] mt-0.5 uppercase">{{ item.type }}</div>
          </div>
        </div>
      </div>
      
      <!-- Load More Sentinel -->
      <div ref="loadMoreSentinel" class="h-10 mt-4"></div>

    </div>

    <!-- Lightbox Overlay -->
    <div 
      v-if="selectedMedia" 
      class="fixed inset-0 z-[200] bg-black/95 flex items-center justify-center backdrop-blur-md"
      @click="selectedMedia = null"
    >
      <div class="absolute top-6 right-6 flex items-center space-x-3 z-10">
        <!-- Jump to Article Button -->
        <button 
          @click.stop="jumpToArticle"
          class="h-12 px-5 rounded-full bg-blue-600 hover:bg-blue-700 text-white flex items-center space-x-2 font-medium transition-colors shadow-lg"
          title="跳转至原文"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
          <span>查看原文</span>
        </button>
        <!-- Close Button -->
        <button 
          class="w-12 h-12 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center text-white transition-colors"
          @click.stop="selectedMedia = null"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <img 
        v-if="selectedMedia.type === 'image'" 
        :src="getMediaUrl(selectedMedia.url)" 
        class="max-w-full max-h-full object-contain"
        @click.stop
      />
      <video 
        v-else-if="selectedMedia.type === 'video'" 
        controls 
        autoplay
        class="max-w-[90%] max-h-[90%] object-contain rounded-lg shadow-2xl bg-black"
        @click.stop
      >
        <source :src="getMediaUrl(selectedMedia.url)" />
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
const selectedMedia = ref(null)

const page = ref(1)
const pageSize = 30
const displayedMediaList = computed(() => {
  return mediaList.value.slice(0, page.value * pageSize)
})

const loadMoreSentinel = ref(null)
let observer = null

const jumpToArticle = async () => {
  if (selectedMedia.value && selectedMedia.value.answer_id) {
    store.currentView = 'home'
    emit('close')
    // Give time for view transition before selecting answer
    setTimeout(() => {
      store.selectAnswer(selectedMedia.value.answer_id)
    }, 100)
  }
}

onMounted(async () => {
  document.body.style.overflow = 'hidden' // Lock scroll
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
      if (page.value * pageSize < mediaList.value.length) {
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
