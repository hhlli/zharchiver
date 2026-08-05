<template>
  <!-- 铺满整个屏幕，移动端使用 min-h-screen 允许原生滚动，电脑端使用 h-screen 和 overflow-hidden 锁定高度 -->
  <div class="min-h-screen md:h-screen w-full bg-white flex flex-col font-sans antialiased select-none md:overflow-hidden">
    
    <div class="sticky md:static top-0 z-40">
      <AppHeader 
        v-model:searchQuery="searchQuery" 
        v-model:activeCategory="activeCategory"
        :tags="tags"
        :currentView="currentView"
        @add-archive="showAddModal = true" 
      />
    </div>

    <!-- 2. 主体区 (侧边栏 + 内容视图) -->
    <div class="flex-1 flex md:overflow-hidden">
      
      <div class="hidden md:block flex-shrink-0 z-30 h-full">
        <AppSidebar 
          v-model:currentView="currentView"
          v-model:activeCategory="activeCategory"
          v-model:activeSetting="activeSetting"
          :tags="tags"
          :hasCurrentAnswer="!!currentAnswer"
          :isDesktopOpen="isDesktopSidebarOpen"
          @collapse-desktop="isDesktopSidebarOpen = false"
          @clear-answer="currentAnswer = null"
        />
      </div>

      <!-- 迷你侧边栏 (仅展开按钮) -->
      <div 
        v-if="!isDesktopSidebarOpen" 
        class="hidden md:flex flex-col justify-end border-r border-gray-200 bg-[#f5f5f7] flex-shrink-0 h-full z-30"
      >
        <div class="p-3">
          <button @click="isDesktopSidebarOpen = true" class="p-1 rounded text-gray-400 hover:text-gray-700 hover:bg-gray-200 transition cursor-pointer flex items-center justify-center" title="展开边栏">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 19l7-7-7-7M5 19l7-7-7-7"></path></svg>
          </button>
        </div>
      </div>

      <!-- 右侧主内容展示区 -->
      <main class="flex-1 bg-gray-50/30 md:bg-white relative flex flex-col min-w-0 md:h-full">
        
        <div class="flex-1 p-4 md:p-6 pb-[calc(6rem+env(safe-area-inset-bottom))] md:pb-6 select-text md:overflow-y-auto">
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
          <DisplaySettings v-else-if="activeSetting === 'display'" />
        </div>

        </div>

      </main>
    </div>

    <!-- 底部状态栏 (PC) -->
    <footer class="hidden md:flex h-8 bg-[#f5f5f7] border-t border-gray-200 items-center justify-between px-6 flex-shrink-0 text-[11px] text-gray-500">
      <span>知乎内容本地归档库</span>
      <span>已加载 {{ answers.length }} 个项目</span>
    </footer>

    <!-- 移动端悬浮菜单 (FAB) -->
    <MobileFloatingMenu 
      v-model:currentView="currentView"
      v-model:activeSetting="activeSetting"
    />

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
import DisplaySettings from './components/settings/DisplaySettings.vue'
import DeleteConfirmModal from './components/DeleteConfirmModal.vue'
import BaseModal from './components/common/BaseModal.vue'
import BaseButton from './components/common/BaseButton.vue'
import MobileFloatingMenu from './components/layout/MobileFloatingMenu.vue'

const API_BASE = ''

const isDesktopSidebarOpen = ref(true)
const showAddModal = ref(false)
const showLoginModal = ref(false)

const currentView = ref('home') 
const activeSetting = ref('auth') 
const viewMode = ref(localStorage.getItem('viewMode') || 'grid') 
provide('viewMode', viewMode)
provide('setViewMode', (mode) => {
  viewMode.value = mode
  localStorage.setItem('viewMode', mode)
})

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
</script>