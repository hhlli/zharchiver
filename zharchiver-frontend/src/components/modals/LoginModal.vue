<template>
  <BaseModal 
    :show="store.showLoginModal"
    contentClass="p-6 space-y-6"
    zIndexClass="z-[90]"
  >
    <div class="text-center space-y-2">
      <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-6 h-6 text-brand" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
      </div>
      <h3 class="text-xl font-bold text-gray-900">需要验证身份</h3>
      <p class="text-sm text-gray-500">此项目已开启安全保护，请输入密码</p>
    </div>

    <form @submit.prevent="submitLogin" class="space-y-4">
      <div>
        <input 
          type="password" 
          v-model="loginPassword"
          placeholder="请输入密码"
          class="w-full border border-gray-300 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition"
          :disabled="loginLoading"
          autofocus
        />
        <p v-if="loginErrorMsg" class="text-red-500 text-xs mt-1">{{ loginErrorMsg }}</p>
      </div>
      <AppButton
        type="submit"
        class="w-full py-3 text-sm rounded-xl"
        :loading="loginLoading"
        :disabled="!loginPassword.trim()"
      >
        {{ loginLoading ? '验证中...' : '确认进入' }}
      </AppButton>
    </form>
  </BaseModal>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import BaseModal from '../common/BaseModal.vue'
import AppButton from '../common/AppButton.vue'

const store = useArchiveStore()

const loginPassword = ref('')
const loginLoading = ref(false)
const loginErrorMsg = ref('')

const submitLogin = async () => {
  if (!loginPassword.value.trim()) return
  loginLoading.value = true
  loginErrorMsg.value = ''

  try {
    const res = await fetch(`${store.API_BASE}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: loginPassword.value.trim() })
    })

    if (!res.ok) {
      throw new Error('密码错误')
    }

    const data = await res.json()
    localStorage.setItem('token', data.token)
    store.showLoginModal = false
    loginPassword.value = ''
    store.fetchAnswersList()
  } catch (err) {
    loginErrorMsg.value = err.message
  } finally {
    loginLoading.value = false
  }
}
</script>
