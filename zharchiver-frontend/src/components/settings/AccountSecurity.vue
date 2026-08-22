<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-2xl font-semibold text-primary mb-8">账户安全</h2>

    <div class="bg-surface rounded-xl border border-line overflow-hidden shadow-sm">
      <div class="px-5 py-4 border-b border-line-light bg-surface-hover/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
        <h3 class="text-sm font-semibold text-secondary">访问密码控制</h3>
        <span class="text-xs text-gray-400">防止未授权用户访问您的应用数据</span>
      </div>
      
      <div class="p-4 md:p-5 space-y-4">
        <div class="flex items-center justify-between h-8">
          <label class="block text-sm font-medium text-secondary">启用访问密码</label>
          <label class="relative inline-flex items-center cursor-pointer flex-shrink-0">
            <input type="checkbox" v-model="isPasswordEnabled" class="sr-only peer">
            <div class="w-11 h-6 bg-line peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-200 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-brand"></div>
          </label>
        </div>

        <!-- 原密码验证 (如果系统当前已开启密码，不管是修改还是关闭，都必须输入原密码) -->
        <div v-if="initialPasswordEnabled" class="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-1 sm:space-y-0 pt-4 border-t border-line-light">
          <label class="block text-sm font-medium text-secondary sm:w-32">原密码</label>
          <input
            type="password"
            v-model="oldPassword"
            placeholder="验证当前密码"
            class="flex-1 w-full sm:max-w-md text-sm border border-line rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
        </div>

        <!-- 密码设置表单 (开启时才显示新密码框) -->
        <div v-if="isPasswordEnabled" class="space-y-4 pt-4 border-t border-line-light">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-1 sm:space-y-0">
            <label class="block text-sm font-medium text-secondary sm:w-32">新密码</label>
            <input
              type="password"
              v-model="newPassword"
              placeholder="输入新密码"
              class="flex-1 w-full sm:max-w-md text-sm border border-line rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-1 sm:space-y-0">
            <label class="block text-sm font-medium text-secondary sm:w-32">确认密码</label>
            <input
              type="password"
              v-model="confirmPassword"
              placeholder="再次输入新密码"
              class="flex-1 w-full sm:max-w-md text-sm border border-line rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
        </div>
        
        <div class="flex justify-end pt-4">
          <BaseButton
            variant="primary"
            @click="updatePassword"
            :disabled="(initialPasswordEnabled && !oldPassword) || (isPasswordEnabled && (!newPassword || !confirmPassword))"
          >
            保存密码设置
          </BaseButton>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseButton from '../common/BaseButton.vue'

const store = useArchiveStore()
const apiFetch = store.apiFetch
const showAlert = store.showAlert

const isPasswordEnabled = ref(false)
const initialPasswordEnabled = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

// 获取后端基础地址，根据你的实际环境调整
const API_BASE = ''

// 组件挂载时获取当前安全状态
onMounted(async () => {
  try {
    const res = await apiFetch(`${API_BASE}/api/auth/status`)
    if (res.ok) {
      const data = await res.json()
      isPasswordEnabled.value = data.enabled
      initialPasswordEnabled.value = data.enabled
    }
  } catch (error) {
    console.error('获取状态失败:', error)
  }
})

const updatePassword = async () => {
  if (isPasswordEnabled.value) {
    if (newPassword.value || confirmPassword.value) {
      if (!oldPassword.value) {
        store.showToast('请填写原密码', 'error')
        return
      }
      if (newPassword.value !== confirmPassword.value) {
        store.showToast('两次输入的密码不一致', 'error')
        return
      }
      if (newPassword.value.length < 6) {
        store.showToast('密码长度至少为 6 位', 'error')
        return
      }
    }
  }

  try {
    const res = await apiFetch(`${API_BASE}/api/auth/update`, {
      method: 'POST',
      body: JSON.stringify({
        enabled: isPasswordEnabled.value,
        password: newPassword.value,
        oldPassword: oldPassword.value
      })
    })

    if (res.ok) {
      store.showToast('账户安全设置已保存')
      oldPassword.value = ''
      newPassword.value = ''
      confirmPassword.value = ''
      
      initialPasswordEnabled.value = isPasswordEnabled.value
      
      if (initialPasswordEnabled.value) {
         window.location.reload()
      }
    } else {
      store.showToast('保存失败', 'error')
    }
  } catch (error) {
    console.error(error)
    store.showToast('网络请求失败', 'error')
  }
}
</script>