<template>
  <div v-if="store.showAddModal" class="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 space-y-4 border border-gray-100">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-gray-800">添加知乎回答归档</h3>
        <button @click="store.showAddModal = false" class="text-gray-400 hover:text-gray-600 text-lg cursor-pointer">&times;</button>
      </div>
      
      <div class="space-y-2">
        <label class="text-xs text-gray-500">知乎回答链接</label>
        <input 
          v-model="inputUrl"
          type="text" 
          placeholder="https://www.zhihu.com/question/xxx/answer/xxx"
          class="w-full border border-gray-300 rounded-lg px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
          :disabled="loading"
        />
        <p v-if="errorMsg" class="text-red-500 text-xs mt-1">{{ errorMsg }}</p>
      </div>

      <div class="space-y-2">
        <label class="text-xs text-gray-500">标签 (可选)</label>
        <div class="flex items-center space-x-2">
          <input 
            v-model="inputTag"
            type="text" 
            list="existing-tags"
            placeholder="输入或选择标签"
            class="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-0"
            :disabled="loading"
          />
          <datalist id="existing-tags">
            <option v-for="t in store.tags" :key="t.name" :value="t.name"></option>
          </datalist>
          <div class="flex items-center space-x-1.5 px-1 flex-shrink-0">
            <button v-for="(hex, c) in hexColors" :key="c" @click="inputTagColor = c" :class="['w-4 h-4 rounded-full cursor-pointer transition', inputTagColor === c ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : 'hover:scale-110']" :style="{ backgroundColor: hex }"></button>
          </div>
        </div>
      </div>

      <div class="flex justify-end space-x-2 pt-2">
        <button 
          @click="store.showAddModal = false"
          class="px-4 py-1.5 border border-gray-300 text-gray-600 rounded-lg text-xs hover:bg-gray-50 transition cursor-pointer"
          :disabled="loading"
        >
          取消
        </button>
        <button 
          @click="submitArchive"
          class="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-xs font-medium hover:bg-blue-700 transition disabled:bg-gray-400 cursor-pointer"
          :disabled="loading"
        >
          {{ loading ? '抓取归档中...' : '开始归档' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()

const inputUrl = ref('')
const inputTag = ref('')
const inputTagColor = ref('blue')
const loading = ref(false)
const errorMsg = ref('')

const hexColors = {
  blue: '#3b82f6',
  red: '#ef4444',
  green: '#10b981',
  yellow: '#f59e0b',
  purple: '#8b5cf6'
}

const submitArchive = async () => {
  if (!inputUrl.value) {
    errorMsg.value = '请输入链接'
    return
  }
  
  loading.value = true
  errorMsg.value = ''
  
  try {
    const res = await store.apiFetch(`${store.API_BASE}/api/archive`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url: inputUrl.value.trim(), tag: inputTag.value.trim(), color: inputTagColor.value })
    })

    if (!res.ok) {
      const errText = await res.text()
      throw new Error(errText || '归档请求失败')
    }

    const result = await res.json()
    inputUrl.value = ''
    inputTag.value = ''
    store.showAddModal = false
    store.onArchiveSuccess(result.data?.answer_id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}
</script>
