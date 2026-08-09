<template>
  <div class="min-h-screen md:h-screen w-full bg-white md:flex md:flex-col font-sans antialiased select-none md:overflow-hidden">
    
    <div class="sticky md:static top-0 z-40">
      <AppHeader />
    </div>

    <div class="md:flex-1 md:flex md:overflow-hidden">
      
      <div v-show="store.isDesktopSidebarOpen" class="hidden md:block flex-shrink-0 z-30 h-full">
        <AppSidebar />
      </div>

      <!-- 迷你侧边栏 (仅展开按钮) -->
      <div 
        v-if="!store.isDesktopSidebarOpen" 
        class="hidden md:flex flex-col justify-end border-r border-gray-200 bg-[#f5f5f7] flex-shrink-0 h-full z-30"
      >
        <div class="p-3">
          <button @click="store.isDesktopSidebarOpen = true" class="p-1 rounded text-gray-400 hover:text-gray-700 hover:bg-gray-200 transition cursor-pointer flex items-center justify-center" title="展开边栏">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 19l7-7-7-7M5 19l7-7-7-7"></path></svg>
          </button>
        </div>
      </div>

      <!-- 右侧主内容展示区 -->
      <main class="md:flex-1 bg-gray-50/30 md:bg-white relative md:flex md:flex-col min-w-0 md:h-full">
        
        <div class="md:flex-1 p-4 md:p-6 pb-[calc(6rem+env(safe-area-inset-bottom))] md:pb-6 select-text md:overflow-y-auto">
        <!-- 删除确认弹窗 -->
        <DeleteConfirmModal />

        <!-- 主页视图 -->
        <div v-if="store.currentView === 'home'" class="md:h-full">
          <ArticleDetail v-if="store.currentAnswer" />
          <GroupDetail v-else-if="store.currentGroup" />
          <HomeView v-else />
        </div>

        <!-- 设置界面 -->
        <div v-else-if="store.currentView === 'settings'" class="md:h-full">
          <PlatformAuth v-if="store.activeSetting === 'auth'" />
          <DataManagement v-else-if="store.activeSetting === 'data'" />
          <AccountSecurity v-else-if="store.activeSetting === 'security'" />
          <AISettings v-else-if="store.activeSetting === 'ai'" />
          <LogCenter v-else-if="store.activeSetting === 'logs'" />
          <DisplaySettings v-else-if="store.activeSetting === 'display'" />
        </div>

        </div>

      </main>
    </div>

    <!-- 底部状态栏 (PC) -->
    <footer class="hidden md:flex h-8 bg-[#f5f5f7] border-t border-gray-200 items-center justify-between px-6 flex-shrink-0 text-[11px] text-gray-500">
      <span>知乎内容本地归档库</span>
      <span>已加载 {{ store.answers.length }} 个项目</span>
    </footer>

    <!-- 移动端悬浮菜单 (FAB) -->
    <MobileFloatingMenu />

    <!-- 底部通知 Toast -->
    <ToastNotification />

    <!-- 弹窗集群 -->
    <AddArchiveModal />
    <LoginModal />
    <TagManageModal />

    <!-- 全局消息弹窗 -->
    <BaseModal 
      :show="store.showGlobalAlert" 
      @close="store.showGlobalAlert = false"
      closeOnOutside
      maxWidthClass="max-w-lg"
    >
      <div class="p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-3">{{ store.globalAlertTitle }}</h3>
        <div class="text-sm text-gray-600 whitespace-pre-wrap max-h-[60vh] overflow-y-auto bg-gray-50 p-3 rounded-lg border border-gray-100 font-mono">{{ store.globalAlertMessage }}</div>
        <div class="mt-6 flex justify-end">
          <BaseButton variant="primary" @click="store.showGlobalAlert = false">
            确定
          </BaseButton>
        </div>
      </div>
    </BaseModal>

    <!-- 悬浮的后台进度条组件 -->
    <GlobalProgress />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useArchiveStore } from './stores/archive'
import AppHeader from './components/layout/AppHeader.vue'
import AppSidebar from './components/layout/AppSidebar.vue'
import HomeView from './components/views/HomeView.vue'
import GroupDetail from './components/views/GroupDetail.vue'
import ArticleDetail from './components/views/ArticleDetail.vue'
import AddArchiveModal from './components/modals/AddArchiveModal.vue'
import LoginModal from './components/modals/LoginModal.vue'
import GlobalProgress from './components/common/GlobalProgress.vue'
import TagManageModal from './components/modals/TagManageModal.vue'

import PlatformAuth from './components/settings/PlatformAuth.vue'
import DataManagement from './components/settings/DataManagement.vue'
import AccountSecurity from './components/settings/AccountSecurity.vue'
import AISettings from './components/settings/AISettings.vue'
import LogCenter from './components/settings/LogCenter.vue'
import DisplaySettings from './components/settings/DisplaySettings.vue'
import DeleteConfirmModal from './components/DeleteConfirmModal.vue'
import BaseModal from './components/common/BaseModal.vue'
import BaseButton from './components/common/BaseButton.vue'
import ToastNotification from './components/common/ToastNotification.vue'
import MobileFloatingMenu from './components/layout/MobileFloatingMenu.vue'

const store = useArchiveStore()

onMounted(() => {
  store.fetchAnswersList()
})
</script>