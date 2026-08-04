<template>
  <div class="max-w-6xl mx-auto pt-4 pb-8 px-4 md:px-8">
    <h2 class="text-2xl font-semibold text-gray-800 mb-8">数据管理</h2>
    
    <div class="space-y-6">
      <!-- 备份与恢复 -->
      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
        <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50">
          <h3 class="text-sm font-semibold text-gray-700">本地备份与恢复</h3>
        </div>
        <div class="divide-y divide-gray-100">
          <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
            <div class="flex-col pr-0 sm:pr-4">
              <span class="block text-sm font-medium text-gray-700">导出完整备份</span>
              <span class="block text-xs text-gray-400 mt-1">将数据库和所有本地图片打包为 ZIP 下载到本地，确保数据安全。</span>
            </div>
            <BaseButton variant="soft" @click="downloadBackup">
              下载备份包
            </BaseButton>
          </div>
          <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
            <div class="flex-col pr-0 sm:pr-4">
              <span class="block text-sm font-medium text-gray-700">从备份恢复</span>
              <span class="block text-xs text-red-500 mt-1">上传 ZIP 备份包进行恢复。警告：此操作将覆盖当前所有数据！</span>
            </div>
            <div class="relative flex-shrink-0 w-full sm:w-auto">
              <input 
                type="file" 
                accept=".zip" 
                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                @change="restoreBackup"
                :disabled="restoring"
              />
              <BaseButton variant="danger-soft" :loading="restoring" customClass="w-full sm:w-auto pointer-events-none">
                {{ restoring ? '恢复中...' : '上传并恢复' }}
              </BaseButton>
            </div>
          </div>
        </div>
      </div>

      <!-- Telegram 云备份 -->
      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
        <div class="px-5 py-4 border-b border-gray-100 bg-gray-50/50 flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-1 sm:space-y-0">
          <h3 class="text-sm font-semibold text-gray-700">Telegram 自动备份</h3>
          <span class="text-xs text-gray-400">每天自动将备份文件发送至您的 Bot 对话</span>
        </div>
        <div class="p-4 md:p-5 space-y-4">
          <div class="flex items-center justify-between h-8">
            <div class="flex items-center space-x-3">
              <label class="block text-sm font-medium text-gray-700">自动备份时间</label>
              <select 
                v-if="tgConfig.telegram_backup_enabled" 
                v-model="tgConfig.telegram_backup_time" 
                @change="saveTgConfig"
                class="border border-gray-300 text-gray-600 rounded-lg px-2 py-1 text-xs focus:outline-none focus:border-blue-500 cursor-pointer"
              >
                <option v-for="h in 24" :key="h" :value="String(h-1).padStart(2, '0') + ':00'">{{ String(h-1).padStart(2, '0') }}:00</option>
              </select>
              <BaseButton 
                v-if="tgConfig.telegram_backup_enabled"
                variant="outline" 
                customClass="!px-3 !py-0.5 !h-[26px] text-xs leading-none" 
                :loading="sendingBackup" 
                @click="sendManualBackup"
              >
                发送备份
              </BaseButton>
            </div>
            <div class="flex items-center space-x-3">
              <div 
                @click="toggleBackup"
                :class="['w-11 h-6 rounded-full cursor-pointer transition-colors relative flex-shrink-0', tgConfig.telegram_backup_enabled ? 'bg-blue-600' : 'bg-gray-300']"
              >
                <div :class="['w-5 h-5 bg-white rounded-full absolute top-0.5 shadow transition-transform', tgConfig.telegram_backup_enabled ? 'translate-x-5.5' : 'translate-x-0.5']"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 其他 -->
      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm divide-y divide-gray-100">
        <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
          <div class="flex-col">
            <span class="block text-sm font-medium text-gray-700">清理本地图片</span>
            <span class="block text-xs text-gray-400 mt-1">扫描并删除已移除归档遗留的孤立图片文件，释放磁盘空间。</span>
          </div>
          <BaseButton variant="outline">
            开始清理
          </BaseButton>
        </div>
        <div class="p-4 md:p-5 flex flex-col sm:flex-row sm:items-center justify-between space-y-3 sm:space-y-0">
          <div class="flex-col">
            <span class="block text-sm font-medium text-gray-700">重建数据库索引</span>
            <span class="block text-xs text-gray-400 mt-1">重新生成 SQLite 虚拟表索引，优化检索性能。</span>
          </div>
          <BaseButton variant="outline">
            重建索引
          </BaseButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import BaseButton from '../common/BaseButton.vue'

const apiFetch = inject('apiFetch')
const API_BASE = ''

const tgConfig = ref({
  telegram_backup_enabled: false,
  telegram_backup_time: ''
})

const saving = ref(false)
const restoring = ref(false)
const sendingBackup = ref(false)

const sendManualBackup = async () => {
  sendingBackup.value = true
  try {
    const res = await apiFetch(`${API_BASE}/api/backup/telegram/send`, { method: 'POST' })
    if (res.ok) {
      alert('备份已成功发送至 Telegram！')
    } else {
      const err = await res.text()
      alert('发送失败: ' + err)
    }
  } catch (e) {
    console.error(e)
    alert('网络错误')
  } finally {
    sendingBackup.value = false
  }
}

const loadTgConfig = async () => {
  try {
    const res = await apiFetch(`${API_BASE}/api/settings/backup`)
    if (res.ok) {
      const data = await res.json()
      tgConfig.value = {
        telegram_backup_enabled: data.telegram_backup_enabled === 'true',
        telegram_backup_time: data.telegram_backup_time || ''
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const toggleBackup = () => {
  tgConfig.value.telegram_backup_enabled = !tgConfig.value.telegram_backup_enabled
  if (tgConfig.value.telegram_backup_enabled && !tgConfig.value.telegram_backup_time) {
    tgConfig.value.telegram_backup_time = '00:00'
  }
  saveTgConfig()
}

const saveTgConfig = async () => {
  saving.value = true
  try {
    const res = await apiFetch(`${API_BASE}/api/settings/backup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        telegram_backup_enabled: tgConfig.value.telegram_backup_enabled ? 'true' : 'false',
        telegram_backup_time: tgConfig.value.telegram_backup_time
      })
    })
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

const downloadBackup = () => {
  const token = localStorage.getItem('token') || ''
  window.open(`${API_BASE}/api/backup/download?token=${encodeURIComponent(token)}`, '_blank')
}

const restoreBackup = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  if (!confirm('警告：从备份恢复将抹除现有所有归档数据和图片！是否确认继续？')) {
    event.target.value = ''
    return
  }

  restoring.value = true
  const formData = new FormData()
  formData.append('backup', file)

  try {
    const token = localStorage.getItem('token')
    // 直接用原生 fetch 因为 formData 会自己带上带 boundary 的 content-type
    const headers = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    const res = await fetch(`${API_BASE}/api/backup/restore`, {
      method: 'POST',
      headers,
      body: formData
    })
    
    if (res.ok) {
      alert('数据恢复成功！请刷新页面。')
      window.location.reload()
    } else {
      const err = await res.text()
      alert('恢复失败: ' + err)
    }
  } catch (e) {
    console.error(e)
    alert('网络错误')
  } finally {
    restoring.value = false
    event.target.value = ''
  }
}

onMounted(() => {
  loadTgConfig()
})
</script>