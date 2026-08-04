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

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">API Base URL</label>
            <input 
              v-model="config.ai_base_url" 
              type="text" 
              placeholder="例如: https://api.openai.com/v1/chat/completions"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">API Key</label>
            <input 
              v-model="config.ai_api_key" 
              type="password" 
              placeholder="sk-..."
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">模型名称 (Model Name)</label>
            <input 
              v-model="config.ai_model_name" 
              type="text" 
              placeholder="例如: gpt-4o, gemini-1.5-pro"
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
            @click="saveConfig"
          >
            保存配置
          </BaseButton>
        </div>
      </div>
    </div>
      
    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm mb-6">
      <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-gray-700">Telegram 机器人配置</h3>
        <span class="text-xs text-gray-400">实现多设备互通与归档备份</span>
      </div>
      <div class="p-4 md:p-5">
        <p class="text-sm text-gray-500 mb-6">
          配置 Telegram Bot Token 和 Chat ID 后，机器人可以将归档数据直接发送给您。
        </p>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Bot Token</label>
            <input 
              v-model="config.telegram_bot_token" 
              type="text" 
              placeholder="例如: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Chat ID</label>
            <input 
              v-model="config.telegram_chat_id" 
              type="text" 
              placeholder="您的 Telegram Chat ID"
              class="w-full text-sm border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="pt-4 flex justify-end">
          <BaseButton 
            variant="primary"
            :loading="isSaving"
            @click="saveConfig"
          >
            保存配置
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- 消息弹窗 -->
    <BaseModal 
      :show="showModal" 
      @close="showModal = false"
      closeOnOutside
      maxWidthClass="max-w-lg"
    >
      <div class="p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-3">{{ modalTitle }}</h3>
        <div class="text-sm text-gray-600 whitespace-pre-wrap max-h-[60vh] overflow-y-auto bg-gray-50 p-3 rounded-lg border border-gray-100 font-mono">{{ modalMessage }}</div>
        <div class="mt-6 flex justify-end">
          <BaseButton variant="primary" @click="showModal = false">
            确定
          </BaseButton>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import BaseModal from '../common/BaseModal.vue'
import BaseButton from '../common/BaseButton.vue'

const apiFetch = inject('apiFetch')
const isSaving = ref(false)
const isTesting = ref(false)

const showModal = ref(false)
const modalTitle = ref('')
const modalMessage = ref('')

const showAlert = (title, message) => {
  modalTitle.value = title
  modalMessage.value = message
  showModal.value = true
}

const config = ref({
  ai_base_url: '',
  ai_api_key: '',
  ai_model_name: '',
  telegram_bot_token: '',
  telegram_chat_id: ''
})

onMounted(async () => {
  try {
    const res = await apiFetch('/api/settings/ai')
    if (res.ok) {
      const data = await res.json()
      config.value = data
    }
  } catch (err) {
    console.error('获取 AI 配置失败:', err)
  }
})

const saveConfig = async () => {
  isSaving.value = true
  try {
    const res = await apiFetch('/api/settings/ai', {
      method: 'POST',
      body: JSON.stringify(config.value)
    })
    
    if (res.ok) {
      showAlert('保存成功', 'AI 配置已成功更新！')
    } else {
      showAlert('保存失败', '请检查网络或后端日志')
    }
  } catch (err) {
    console.error('保存失败:', err)
    showAlert('错误', '请求发生错误')
  } finally {
    isSaving.value = false
  }
}

const testConnection = async () => {
  if (!config.value.ai_base_url || !config.value.ai_api_key || !config.value.ai_model_name) {
    showAlert('提示', '请先填写完整所有配置')
    return
  }

  isTesting.value = true
  try {
    const res = await apiFetch('/api/settings/ai/test', {
      method: 'POST',
      body: JSON.stringify(config.value)
    })
    
    if (res.ok) {
      showAlert('🎉 连通性测试成功', '您的 AI 配置一切正常，随时可以开始归档！')
    } else {
      const errorText = await res.text()
      showAlert('测试失败', errorText)
    }
  } catch (err) {
    console.error('测试失败:', err)
    showAlert('请求错误', '请检查网络或后端日志')
  } finally {
    isTesting.value = false
  }
}
</script>
