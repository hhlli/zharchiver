<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-xl md:text-2xl font-semibold text-gray-800 mb-8">工具配置</h2>
    
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-gray-700">视觉大模型 (Vision LLM) 接口</h3>
        <span class="text-xs text-gray-400">接口需兼容 OpenAI 的 /v1/chat/completions 规范</span>
      </div>
      <div class="p-4 md:p-5">
        <p class="text-sm text-gray-500 mb-6">
          配置大模型 API 接口，即可在 Telegram 机器人中发送包含正文和标题的截图，AI 将自动为您提炼文字并归档。
        </p>

        <!-- Dropdown for profiles -->
        <div class="mb-6 flex items-center justify-between bg-gray-50 p-3 rounded-lg border border-gray-200">
          <div class="flex-1 flex items-center space-x-3">
            <label class="text-sm font-medium text-gray-700">当前使用配置:</label>
            <select v-model="config.ai_active_profile_id" class="flex-1 max-w-[200px] border border-gray-300 bg-white text-gray-700 rounded-md px-3 py-1.5 text-sm focus:outline-none focus:border-blue-500 cursor-pointer">
              <option v-for="profile in config.ai_profiles" :key="profile.id" :value="profile.id">
                {{ profile.name }}
              </option>
            </select>
          </div>
          <div class="flex items-center space-x-2">
            <button @click="addNewProfile" class="text-brand hover:text-brand-hover text-sm font-medium px-2 py-1 transition cursor-pointer">
              + 新建配置
            </button>
            <button @click="deleteActiveProfile" :disabled="config.ai_profiles.length <= 1" :class="['text-sm font-medium px-2 py-1 transition', config.ai_profiles.length <= 1 ? 'text-gray-400 cursor-not-allowed' : 'text-red-500 hover:text-red-600 cursor-pointer']" title="删除当前配置">
              删除
            </button>
          </div>
        </div>

        <div v-if="activeProfile" class="space-y-4 pt-2 border-t border-gray-100 mt-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">配置名称</label>
            <input 
              v-model="activeProfile.name" 
              type="text" 
              placeholder="例如: 通义千问 Qwen-VL-Max"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">API Base URL</label>
            <input 
              v-model="activeProfile.base_url" 
              type="text" 
              placeholder="例如: https://api.openai.com/v1/chat/completions"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">API Key</label>
            <input 
              v-model="activeProfile.api_key" 
              type="password" 
              placeholder="sk-..."
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">模型名称 (Model Name)</label>
            <input 
              v-model="activeProfile.model_name" 
              type="text" 
              placeholder="例如: gpt-4o, qwen-vl-max"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
        
        <div class="pt-4 flex justify-end space-x-3">
          <BaseButton 
            variant="outline"
            :loading="isTesting"
            @click="testConnection"
          >
            测试连通性
          </BaseButton>
          
          <BaseButton 
            variant="primary"
            :loading="isSaving"
            @click="saveConfig('ai')"
          >
            保存配置
          </BaseButton>
        </div>
      </div>
    </div>
      
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-gray-700">Telegram 归档机器人</h3>
        <span class="text-xs text-gray-400">负责监听您的消息并自动抓取</span>
      </div>
      <div class="p-4 md:p-5">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Telegram API Endpoint</label>
            <input 
              v-model="tgConfig.telegram_api_endpoint" 
              type="text" 
              placeholder="默认: https://api.telegram.org" 
              class="w-full h-10 px-3 rounded-lg border border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition"
            />
            <p class="text-xs text-gray-500 mt-1">
              如需突破 50MB 官方限制，请填写您的 Local Telegram Bot API 服务器地址 (例如: http://127.0.0.1:8081)
            </p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Bot Token</label>
            <input 
              v-model="tgConfig.telegram_bot_token" 
              type="text" 
              placeholder="例如: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">您的 Chat ID</label>
            <input 
              v-model="tgConfig.telegram_chat_id" 
              type="text" 
              placeholder="例如: 123456789"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="pt-4 flex justify-end space-x-3">
          <BaseButton 
            variant="outline"
            :loading="isTestingTgArchive"
            @click="testConnectionTg('archive')"
          >
            测试归档机器人
          </BaseButton>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-gray-700">Telegram 推送机器人</h3>
        <span class="text-xs text-gray-400">负责一键推送文章给您或频道</span>
      </div>
      <div class="p-4 md:p-5">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Push Bot Token</label>
            <input 
              v-model="tgConfig.telegram_push_bot_token" 
              type="text" 
              placeholder="推送机器人的 Token"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">目标频道或 Chat ID</label>
            <input 
              v-model="tgConfig.telegram_push_chat_id" 
              type="text" 
              placeholder="例如: @my_channel_id 或 123456789"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="pt-4 flex justify-end space-x-3">
          <BaseButton 
            variant="outline"
            :loading="isTestingTgPush"
            @click="testConnectionTg('push')"
          >
            测试推送机器人
          </BaseButton>
          <BaseButton 
            variant="primary"
            :loading="isSaving"
            @click="saveConfig('tg')"
          >
            保存所有 TG 配置
          </BaseButton>
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
const apiFetch = store.apiFetch
const showAlert = store.showAlert

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
      const errorText = await res.text()
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
      const errorText = await res.text()
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
