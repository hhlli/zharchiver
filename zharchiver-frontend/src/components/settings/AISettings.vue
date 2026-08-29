<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8 transition-opacity duration-300" :class="isLoading ? 'opacity-0' : 'opacity-100'">
    <h2 class="text-xl md:text-2xl font-semibold text-primary mb-8">工具配置</h2>
    
    <!-- AI -->
    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-line-light bg-surface-hover/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-secondary">视觉大模型 (Vision LLM) 接口</h3>
        <span class="text-xs text-gray-400">接口需兼容 OpenAI 的 /v1/chat/completions 规范</span>
      </div>
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex flex-col space-y-1.5">
          <div class="flex items-center space-x-2">
            <span v-if="config.ai_profiles && config.ai_profiles.length > 0" class="px-2.5 py-0.5 bg-green-50 dark:bg-green-900/30 text-green-600 dark:text-green-400 text-xs font-medium rounded-full border border-green-100 dark:border-green-900/50">已配置</span>
            <span v-else class="px-2.5 py-0.5 bg-surface-hover text-muted text-xs font-medium rounded-full border border-line">未配置</span>
          </div>
          <p class="text-xs text-muted">配置大模型 API，让机器人在归档时智能提取截图正文和标题，并自动为您分类。</p>
        </div>
        <button @click="openModal('ai')" class="w-full sm:w-auto px-4 py-2 bg-surface hover:bg-surface-hover text-brand text-sm font-medium rounded-lg transition-colors border border-line shadow-sm flex items-center justify-center space-x-1.5 whitespace-nowrap">
          <span>⚙️ 配置接口</span>
        </button>
      </div>
    </div>

    <!-- TG Archive -->
    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-line-light bg-surface-hover/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-secondary">Telegram 归档机器人</h3>
        <span class="text-xs text-gray-400">负责监听您的消息并自动抓取</span>
      </div>
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex flex-col space-y-1.5">
          <div class="flex items-center space-x-2">
            <span v-if="tgConfig.telegram_bot_token && tgConfig.telegram_chat_id" class="px-2.5 py-0.5 bg-green-50 dark:bg-green-900/30 text-green-600 dark:text-green-400 text-xs font-medium rounded-full border border-green-100 dark:border-green-900/50">已配置</span>
            <span v-else class="px-2.5 py-0.5 bg-surface-hover text-muted text-xs font-medium rounded-full border border-line">未配置</span>
          </div>
          <p class="text-xs text-muted">配置专属 Bot，监听您的私聊消息或群消息，将链接、文本、图片自动抓取到归档系统。</p>
        </div>
        <button @click="openModal('tg_archive')" class="w-full sm:w-auto px-4 py-2 bg-surface hover:bg-surface-hover text-brand text-sm font-medium rounded-lg transition-colors border border-line shadow-sm flex items-center justify-center space-x-1.5 whitespace-nowrap">
          <span>⚙️ 配置机器人</span>
        </button>
      </div>
    </div>

    <!-- TG Push -->
    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-line-light bg-surface-hover/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-secondary">Telegram 推送机器人</h3>
        <span class="text-xs text-gray-400">负责一键推送文章给您或频道</span>
      </div>
      <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
        <div class="flex flex-col space-y-1.5">
          <div class="flex items-center space-x-2">
            <span v-if="tgConfig.telegram_push_bot_token && tgConfig.telegram_push_chat_id" class="px-2.5 py-0.5 bg-green-50 dark:bg-green-900/30 text-green-600 dark:text-green-400 text-xs font-medium rounded-full border border-green-100 dark:border-green-900/50">已配置</span>
            <span v-else class="px-2.5 py-0.5 bg-surface-hover text-muted text-xs font-medium rounded-full border border-line">未配置</span>
          </div>
          <p class="text-xs text-muted">配置推送频道，在网页端一键将排版文章推送至您的 Telegram。</p>
        </div>
        <button @click="openModal('tg_push')" class="w-full sm:w-auto px-4 py-2 bg-surface hover:bg-surface-hover text-brand text-sm font-medium rounded-lg transition-colors border border-line shadow-sm flex items-center justify-center space-x-1.5 whitespace-nowrap">
          <span>⚙️ 配置推送</span>
        </button>
      </div>
    </div>

    <!-- 全局统一样式的弹窗层 -->
    <div v-if="activeModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm transition-opacity">
      <div class="bg-surface rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden transform transition-all flex flex-col max-h-[90vh]">
        <!-- 弹窗 Header -->
        <div class="px-6 py-4 border-b border-line-light flex justify-between items-center bg-surface-hover/50">
          <h3 class="text-base font-bold text-primary flex items-center gap-2">
            <span v-if="activeModal === 'ai'">🤖 视觉大模型配置</span>
            <span v-if="activeModal === 'tg_archive'">📥 归档机器人配置</span>
            <span v-if="activeModal === 'tg_push'">📤 推送机器人配置</span>
          </h3>
          <button @click="closeModal" class="text-gray-400 hover:text-secondary transition-colors cursor-pointer">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>

        <!-- 弹窗 Body -->
        <div class="p-6 overflow-y-auto space-y-5">
          <!-- 1. AI 模型配置内容 -->
          <template v-if="activeModal === 'ai'">
            <!-- 配置选择器 -->
            <div class="bg-blue-50/50 dark:bg-blue-900/20 p-3.5 rounded-xl border border-blue-100 dark:border-blue-900/50">
              <label class="block text-xs font-semibold text-blue-900 dark:text-blue-400 mb-2">切换配置方案:</label>
              <div class="flex items-center gap-2.5">
                <select v-model="config.ai_active_profile_id" class="flex-1 min-w-0 bg-surface border border-blue-200 dark:border-blue-800 text-blue-900 dark:text-blue-400 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-brand cursor-pointer shadow-sm">
                  <option v-for="profile in config.ai_profiles" :key="profile.id" :value="profile.id">
                    {{ profile.name }}
                  </option>
                </select>
                <div class="flex items-center space-x-2 flex-shrink-0">
                  <button @click="addNewProfile" class="w-9 h-9 bg-surface text-blue-600 dark:text-blue-400 rounded-lg border border-blue-200 dark:border-blue-800 hover:bg-blue-50 dark:hover:bg-blue-900/30 dark:hover:bg-blue-900/40 transition flex items-center justify-center cursor-pointer shadow-sm" title="新建方案">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
                  </button>
                  <button @click="deleteActiveProfile" :disabled="config.ai_profiles.length <= 1" :class="['w-9 h-9 rounded-lg border transition flex items-center justify-center shadow-sm', config.ai_profiles.length <= 1 ? 'bg-surface-hover text-gray-400 border-line cursor-not-allowed' : 'bg-surface text-red-500 border-red-200 hover:bg-red-50 dark:hover:bg-red-900/30 cursor-pointer']" title="删除当前方案">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </div>
              </div>
            </div>

            <!-- 当前配置编辑区 -->
            <div v-if="activeProfile" class="space-y-4">
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">配置名称</label>
                <input v-model="activeProfile.name" type="text" placeholder="例如: OpenAI GPT-4o" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">API Base URL</label>
                <input v-model="activeProfile.base_url" type="text" placeholder="https://api.openai.com/v1/chat/completions" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">API Key</label>
                <input v-model="activeProfile.api_key" type="password" placeholder="sk-..." class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">Model Name (模型标识)</label>
                <input v-model="activeProfile.model_name" type="text" placeholder="例如: gpt-4o, qwen-vl-max" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
            </div>
          </template>

          <!-- 2. TG 归档机器人配置内容 -->
          <template v-if="activeModal === 'tg_archive'">
            <div class="space-y-4">
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">Telegram API Endpoint (可选)</label>
                <input v-model="tgConfig.telegram_api_endpoint" type="text" placeholder="默认: https://api.telegram.org" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
                <p class="text-[11px] text-gray-400 mt-1.5 leading-tight">如需抓取大于 20MB 的文件或国内直连，请填写您自建的 Local Telegram Bot API 服务器地址 (例如: http://127.0.0.1:8081)。留空则使用官方 API。</p>
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">Bot Token</label>
                <input v-model="tgConfig.telegram_bot_token" type="password" placeholder="例如: 123456:ABC-DEF1234ghIkl..." class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">您的 Chat ID (用于接收抓取的对话)</label>
                <input v-model="tgConfig.telegram_chat_id" type="text" placeholder="例如: 123456789" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
            </div>
          </template>

          <!-- 3. TG 推送机器人配置内容 -->
          <template v-if="activeModal === 'tg_push'">
            <div class="space-y-4">
              <div class="bg-blue-50/50 dark:bg-blue-900/20 p-3 rounded-lg border border-blue-100 dark:border-blue-900/50 mb-2 text-xs text-blue-800 dark:text-blue-300 leading-relaxed">
                推送机器人用于在网页上点击“分享到 TG”时，将格式化排版好的文章推送到指定频道或个人对话中。您可以与归档机器人使用同一个 Bot。
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">Push Bot Token</label>
                <input v-model="tgConfig.telegram_push_bot_token" type="password" placeholder="推送机器人的 Token" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-secondary mb-1.5">目标频道或 Chat ID</label>
                <input v-model="tgConfig.telegram_push_chat_id" type="text" placeholder="例如: @my_channel_id 或 123456789" class="w-full h-10 px-3 bg-surface-hover border border-line rounded-lg text-sm font-mono focus:outline-none focus:bg-surface focus:border-brand focus:ring-1 focus:ring-brand transition" />
              </div>
            </div>
          </template>
        </div>
        
        <!-- 弹窗 Footer -->
        <div class="px-6 py-4 bg-surface-hover/50 border-t border-line-light flex justify-end space-x-3">
          <template v-if="activeModal === 'ai'">
            <BaseButton variant="outline" :loading="isTesting" @click="testConnection">测试连通性</BaseButton>
            <BaseButton variant="primary" :loading="isSaving" @click="saveConfig('ai')">保存 AI 配置</BaseButton>
          </template>
          
          <template v-if="activeModal === 'tg_archive'">
            <BaseButton variant="outline" :loading="isTestingTgArchive" @click="testConnectionTg('archive')">测试归档机器人</BaseButton>
            <BaseButton variant="primary" :loading="isSaving" @click="saveConfig('tg')">保存 TG 配置</BaseButton>
          </template>

          <template v-if="activeModal === 'tg_push'">
            <BaseButton variant="outline" :loading="isTestingTgPush" @click="testConnectionTg('push')">测试推送</BaseButton>
            <BaseButton variant="primary" :loading="isSaving" @click="saveConfig('tg')">保存 TG 配置</BaseButton>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseButton from '../common/BaseButton.vue'

const store = useArchiveStore()
const isLoading = ref(true)
const apiFetch = store.apiFetch

// 控制当前打开的弹窗类型: '' | 'ai' | 'tg_archive' | 'tg_push'
const activeModal = ref('')

const isSaving = ref(false)
const isTesting = ref(false)
const isTestingTgArchive = ref(false)
const isTestingTgPush = ref(false)

const config = ref({
  ai_profiles: [],
  ai_active_profile_id: ''
})

const tgConfig = ref({
  telegram_bot_token: '',
  telegram_chat_id: '',
  telegram_push_bot_token: '',
  telegram_push_chat_id: '',
  telegram_api_endpoint: ''
})

const activeProfile = computed(() => {
  return config.value.ai_profiles.find(p => p.id === config.value.ai_active_profile_id) || null
})

const openModal = (type) => {
  activeModal.value = type
}

const closeModal = () => {
  activeModal.value = ''
  // 可选：在这里可以重新拉取一次数据，丢弃未保存的修改
}

const addNewProfile = () => {
  const newId = 'profile_' + Date.now()
  config.value.ai_profiles.push({
    id: newId,
    name: '新建配置',
    base_url: '',
    api_key: '',
    model_name: ''
  })
  config.value.ai_active_profile_id = newId
}

const deleteActiveProfile = () => {
  if (config.value.ai_profiles.length <= 1) return
  const idToDelete = config.value.ai_active_profile_id
  config.value.ai_profiles = config.value.ai_profiles.filter(p => p.id !== idToDelete)
  config.value.ai_active_profile_id = config.value.ai_profiles[0].id
}

onMounted(async () => {
  try {
    const res = await apiFetch('/api/settings/ai')
    if (res.ok) {
      const data = await res.json()
      config.value.ai_profiles = data.ai_profiles || []
      config.value.ai_active_profile_id = data.ai_active_profile_id || ''
    }
  } catch (err) {
    console.error('获取 AI 配置失败:', err)
  }

  try {
    const resTg = await apiFetch('/api/settings/telegram')
    if (resTg.ok) {
      const dataTg = await resTg.json()
      tgConfig.value = dataTg
    }
  } catch (err) {
    console.error('获取 Telegram 配置失败:', err)
  }
  isLoading.value = false
})

const saveConfig = async (type) => {
  isSaving.value = true
  try {
    let url = ''
    let payload = {}

    if (type === 'ai') {
      url = '/api/settings/ai'
      payload = config.value
    } else {
      url = '/api/settings/telegram'
      payload = tgConfig.value
    }

    const res = await apiFetch(url, {
      method: 'POST',
      body: JSON.stringify(payload)
    })
    
    if (res.ok) {
      if (type === 'ai') {
        store.showToast('AI 配置已成功更新！')
      } else {
        store.showToast('Telegram 配置已成功更新！')
      }
      closeModal()
    } else {
      store.showToast('保存失败', 'error')
    }
  } catch (err) {
    console.error('保存失败:', err)
    store.showToast('请求发生错误', 'error')
  } finally {
    isSaving.value = false
  }
}

const testConnection = async () => {
  if (!activeProfile.value) {
    store.showToast('请先选择或创建一个模型', 'error')
    return
  }
  if (!activeProfile.value.base_url || !activeProfile.value.api_key || !activeProfile.value.model_name) {
    store.showToast('请先填写完整所有配置', 'error')
    return
  }

  isTesting.value = true
  try {
    const res = await apiFetch('/api/settings/ai/test', {
      method: 'POST',
      body: JSON.stringify(activeProfile.value)
    })
    
    if (res.ok) {
      store.showToast('AI 配置连通性测试成功')
    } else {
      store.showToast('测试失败', 'error')
    }
  } catch (err) {
    console.error('测试失败:', err)
    store.showToast('请求错误，请检查网络', 'error')
  } finally {
    isTesting.value = false
  }
}

const testConnectionTg = async (botType) => {
  if (botType === 'archive' && (!tgConfig.value.telegram_bot_token || !tgConfig.value.telegram_chat_id)) {
    store.showToast('请先填写完整的归档机器人配置', 'error')
    return
  }
  if (botType === 'push' && (!tgConfig.value.telegram_push_bot_token || !tgConfig.value.telegram_push_chat_id)) {
    store.showToast('请先填写完整的推送机器人配置', 'error')
    return
  }

  if (botType === 'archive') {
    isTestingTgArchive.value = true
  } else {
    isTestingTgPush.value = true
  }
  
  try {
    const res = await apiFetch('/api/settings/telegram/test', {
      method: 'POST',
      body: JSON.stringify({
        ...tgConfig.value,
        bot_type: botType
      })
    })
    
    if (res.ok) {
      store.showToast('测试消息已发送')
    } else {
      store.showToast('发送失败', 'error')
    }
  } catch (err) {
    console.error('发送失败:', err)
    store.showToast('请求错误', 'error')
  } finally {
    if (botType === 'archive') {
      isTestingTgArchive.value = false
    } else {
      isTestingTgPush.value = false
    }
  }
}
</script>
