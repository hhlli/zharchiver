<template>
  <!-- 铺满整个屏幕 -->
  <div class="h-screen w-screen bg-white flex flex-col overflow-hidden font-sans antialiased select-none">
    
    <AppHeader 
      v-model:searchQuery="searchQuery" 
      @add-archive="showAddModal = true" 
      @toggle-sidebar="toggleSidebar"
    />

    <!-- 2. 主体区 (侧边栏 + 内容视图) -->
    <div class="flex-1 flex overflow-hidden">
      
      <AppSidebar 
        v-model:currentView="currentView"
        v-model:activeCategory="activeCategory"
        v-model:activeSetting="activeSetting"
        :tags="tags"
        :hasCurrentAnswer="!!currentAnswer"
        :isOpen="isMobileSidebarOpen"
        :isDesktopOpen="isDesktopSidebarOpen"
        @close="isMobileSidebarOpen = false"
        @collapse-desktop="isDesktopSidebarOpen = false"
        @clear-answer="currentAnswer = null"
      />

      <!-- 迷你侧边栏 (仅展开按钮) -->
      <div 
        v-if="!isDesktopSidebarOpen" 
        class="hidden md:flex flex-col justify-end border-r border-gray-200 bg-[#f5f5f7] flex-shrink-0"
      >
        <div class="p-3">
          <button @click="isDesktopSidebarOpen = true" class="p-1 rounded text-gray-400 hover:text-gray-700 hover:bg-gray-200 transition cursor-pointer flex items-center justify-center" title="展开边栏">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 19l7-7-7-7M5 19l7-7-7-7"></path></svg>
          </button>
        </div>
      </div>

      <!-- 右侧主内容展示区 -->
      <main class="flex-1 bg-gray-50/30 md:bg-white relative flex flex-col min-w-0">
        
        <div class="flex-1 overflow-y-auto p-4 md:p-6 pb-[calc(6rem+env(safe-area-inset-bottom))] md:pb-6 select-text">
        <!-- 删除确认弹窗 -->
        <DeleteConfirmModal
          :show="itemToDelete !== null"
          :title="itemToDelete?.title"
          @cancel="itemToDelete = null"
          @confirm="confirmDelete"
        />

        <!-- 主页视图 -->
        <div v-if="currentView === 'home'" class="h-full">
          <ArticleDetail 
            v-if="currentAnswer" 
            :answer="currentAnswer" 
            @back="currentAnswer = null" 
            @refresh="fetchAnswersList"
          />
          <HomeView 
            v-else 
            :answers="filteredAnswers" 
            :viewMode="viewMode" 
            @view="selectAnswer" 
            @delete="itemToDelete = $event" 
          />
        </div>

        <!-- 设置界面 -->
        <div v-else-if="currentView === 'settings'" class="h-full">
          <PlatformAuth v-if="activeSetting === 'auth'" />
          <DataManagement v-else-if="activeSetting === 'data'" />
          <AccountSecurity v-else-if="activeSetting === 'security'" />
          <AISettings v-else-if="activeSetting === 'ai'" />
          <LogCenter v-else-if="activeSetting === 'logs'" />
        </div>

        </div>

        <!-- 悬浮的视图切换控件 (Floating View Toggle) -->
        <div v-if="currentView === 'home' && !currentAnswer" class="absolute bottom-[calc(2rem+env(safe-area-inset-bottom))] md:bottom-6 left-1/2 transform -translate-x-1/2 bg-white/80 backdrop-blur-md shadow-[0_2px_12px_rgba(0,0,0,0.1)] border border-gray-200/60 rounded-full p-1 flex items-center text-xs font-medium text-gray-500 z-10 transition-all hover:shadow-[0_4px_16px_rgba(0,0,0,0.12)]">
          <button 
            @click="viewMode = 'grid'"
            :class="['px-5 py-1.5 rounded-full transition-all cursor-pointer flex items-center space-x-1.5', viewMode === 'grid' ? 'bg-blue-50 text-blue-600 shadow-sm' : 'hover:text-gray-900']"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
            <span>网格</span>
          </button>
          <button 
            @click="viewMode = 'list'"
            :class="['px-5 py-1.5 rounded-full transition-all cursor-pointer flex items-center space-x-1.5', viewMode === 'list' ? 'bg-blue-50 text-blue-600 shadow-sm' : 'hover:text-gray-900']"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            <span>列表</span>
          </button>
        </div>

      </main>
    </div>

    <!-- 底部状态栏 (PC) -->
    <footer class="hidden md:flex h-8 bg-[#f5f5f7] border-t border-gray-200 items-center justify-between px-6 flex-shrink-0 text-[11px] text-gray-500">
      <span>知乎内容本地归档库</span>
      <span>已加载 {{ answers.length }} 个项目</span>
    </footer>

    <!-- 弹窗集群 -->
    <AddArchiveModal v-model:show="showAddModal" :tags="tags" @success="onArchiveSuccess" />
    <LoginModal v-model:show="showLoginModal" @success="fetchAnswersList" />

    <!-- 全局消息弹窗 -->
    <BaseModal 
      :show="showGlobalAlert" 
      @close="showGlobalAlert = false"
      closeOnOutside
      maxWidthClass="max-w-lg"
    >
      <div class="p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-3">{{ globalAlertTitle }}</h3>
        <div class="text-sm text-gray-600 whitespace-pre-wrap max-h-[60vh] overflow-y-auto bg-gray-50 p-3 rounded-lg border border-gray-100 font-mono">{{ globalAlertMessage }}</div>
        <div class="mt-6 flex justify-end">
          <BaseButton variant="primary" @click="showGlobalAlert = false">
            确定
          </BaseButton>
        </div>
      </div>
    </BaseModal>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, provide } from 'vue'
import AppHeader from './components/layout/AppHeader.vue'
import AppSidebar from './components/layout/AppSidebar.vue'
import HomeView from './components/views/HomeView.vue'
import ArticleDetail from './components/views/ArticleDetail.vue'
import AddArchiveModal from './components/modals/AddArchiveModal.vue'
import LoginModal from './components/modals/LoginModal.vue'

import PlatformAuth from './components/settings/PlatformAuth.vue'
import DataManagement from './components/settings/DataManagement.vue'
import AccountSecurity from './components/settings/AccountSecurity.vue'
import AISettings from './components/settings/AISettings.vue'
import LogCenter from './components/settings/LogCenter.vue'
import DeleteConfirmModal from './components/DeleteConfirmModal.vue'
import BaseModal from './components/common/BaseModal.vue'
import BaseButton from './components/common/BaseButton.vue'

const API_BASE = ''

const isDesktopSidebarOpen = ref(true)
const isMobileSidebarOpen = ref(false)
const showAddModal = ref(false)
const showLoginModal = ref(false)

const currentView = ref('home') 
const activeSetting = ref('auth') 
const viewMode = ref('grid') 
const searchQuery = ref('')
const activeCategory = ref('all') // 'all', 'marked', 'images'

const answers = ref([])
const currentAnswer = ref(null)
const itemToDelete = ref(null)

const showGlobalAlert = ref(false)
const globalAlertTitle = ref('')
const globalAlertMessage = ref('')

const showAlert = (title, message) => {
  globalAlertTitle.value = title
  globalAlertMessage.value = message
  showGlobalAlert.value = true
}
provide('showAlert', showAlert)

const apiFetch = async (url, options = {}) => {
  const token = localStorage.getItem('token')
  const headers = { ...options.headers }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  const res = await fetch(url, { ...options, headers })
  
  if (res.status === 401) {
    showLoginModal.value = true
    throw new Error('未授权，请重新登录')
  }
  
  return res
}
provide('apiFetch', apiFetch)

const tags = computed(() => {
  const map = new Map()
  answers.value.forEach(a => {
    if (a.tag && !map.has(a.tag)) {
      map.set(a.tag, a.tag_color || 'blue')
    }
  })
  return Array.from(map.entries())
    .map(([name, color]) => ({ name, color }))
    .sort((a, b) => a.name.localeCompare(b.name))
})

const filteredAnswers = computed(() => {
  let list = answers.value

  if (activeCategory.value !== 'all') {
    list = list.filter(item => item.tag === activeCategory.value)
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(item => 
      item.title.toLowerCase().includes(q) || 
      item.author_name.toLowerCase().includes(q)
    )
  }
  
  return list
})

const fetchAnswersList = async () => {
  try {
    const res = await apiFetch(`${API_BASE}/api/answers`)
    if (res.ok) {
      answers.value = await res.json()
    }
  } catch (err) {
    console.error('获取列表失败:', err)
  }
}

const selectAnswer = async (id) => {
  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${id}`)
    if (res.ok) {
      currentAnswer.value = await res.json()
    }
  } catch (err) {
    console.error('获取详情失败:', err)
  }
}

const onArchiveSuccess = async (id) => {
  await fetchAnswersList()
  if (id) {
    selectAnswer(id)
  }
}

const confirmDelete = async () => {
  if (!itemToDelete.value) return
  const id = itemToDelete.value.answer_id
  itemToDelete.value = null
  
  try {
    const res = await apiFetch(`${API_BASE}/api/answers/${id}`, { method: 'DELETE' })
    if (res.ok) {
      if (currentAnswer.value && currentAnswer.value.answer_id === id) {
        currentAnswer.value = null
      }
      await fetchAnswersList()
    } else {
      const errData = await res.json()
      showAlert('删除失败', errData.message)
    }
  } catch (e) {
    console.error(e)
    showAlert('错误', '网络请求失败')
  }
}

onMounted(() => {
  fetchAnswersList()
})

const toggleSidebar = () => {
  if (window.innerWidth >= 768) {
    isDesktopSidebarOpen.value = !isDesktopSidebarOpen.value
  } else {
    isMobileSidebarOpen.value = !isMobileSidebarOpen.value
  }
}
</script>