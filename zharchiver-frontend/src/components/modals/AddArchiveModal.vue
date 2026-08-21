<template>
  <div v-if="store.showAddModal" class="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md p-6 space-y-4 border border-gray-100">
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-gray-800">添加归档</h3>
        <button @click="store.showAddModal = false" class="text-gray-400 hover:text-gray-600 text-lg cursor-pointer">&times;</button>
      </div>

      <div class="space-y-2">
        <label class="text-xs text-gray-500">粘贴链接</label>
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

        <!-- 标签名输入 + 下拉 -->
        <div class="relative">
          <input
            v-model="inputTag"
            type="text"
            placeholder="输入或选择标签"
            class="w-full border border-gray-300 rounded-lg px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500 pr-7"
            :disabled="loading"
            @focus="showTagDropdown = true"
            @blur="hideTagDropdown"
          />
          <div class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
          </div>
          <div v-if="showTagDropdown && store.allTags.length > 0" class="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-40 overflow-y-auto">
            <div
              v-for="t in store.allTags"
              :key="t.name"
              @mousedown.prevent="selectExistingTag(t)"
              class="px-3 py-2 text-xs text-gray-700 hover:bg-blue-50 cursor-pointer flex items-center space-x-2"
            >
              <span class="w-2 h-2 rounded-full flex-shrink-0 shadow-[0_0_2px_rgba(0,0,0,0.12)]" :style="{ background: resolveTagBackground(t.color) }"></span>
              <span class="truncate">{{ t.name }}</span>
            </div>
          </div>
        </div>

        <!-- 颜色选择：实时预览圆点 + 调色盘弹出 -->
        <div v-if="inputTag.trim()" class="flex items-center space-x-2 pt-0.5">
          <span class="text-xs text-gray-400 flex-shrink-0">标签颜色</span>
          <div class="relative">
            <div
              class="w-4 h-4 rounded-full cursor-pointer border border-black/10 shadow-sm transition hover:scale-110"
              :style="{ background: currentPreviewColor }"
              @click="showColorPicker = !showColorPicker"
              title="点击选择颜色"
            ></div>
            <div
              v-if="showColorPicker"
              class="absolute z-50 left-0 top-6 bg-white rounded-xl shadow-xl border border-gray-100 p-2"
            >
              <ColorPicker
                isWidget
                :pureColor="pureColor"
                :gradientColor="gradientColor"
                :activeKey="activeKey"
                @update:pureColor="onPureColorChange"
                @update:gradientColor="onGradientColorChange"
                @update:activeKey="val => activeKey = val"
              />
              <div class="flex justify-end mt-1.5">
                <button @click="showColorPicker = false" class="text-xs text-gray-500 hover:text-brand px-2 py-0.5 rounded cursor-pointer">完成</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end space-x-2 pt-2">
        <button
          @click="store.showAddModal = false"
          class="px-4 py-1.5 border border-gray-300 text-gray-600 rounded-lg text-xs hover:bg-gray-50 transition cursor-pointer"
          :disabled="loading"
        >取消</button>
        <button
          @click="submitArchive"
          class="px-4 py-1.5 bg-brand text-white rounded-lg text-xs font-medium hover:bg-brand-hover transition disabled:bg-gray-400 cursor-pointer"
          :disabled="loading"
        >{{ loading ? '抓取归档中...' : '开始归档' }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import { ColorPicker } from 'vue3-colorpicker'
import 'vue3-colorpicker/style.css'
import { resolveTagBackground, colorToPickerState, pickerStateToColor } from '../../utils/tagColor'

const store = useArchiveStore()

const inputUrl = ref('')
const inputTag = ref('')
const loading = ref(false)
const errorMsg = ref('')
const showTagDropdown = ref(false)
const showColorPicker = ref(false)

// ColorPicker 三状态
const pureColor = ref('#3b82f6')
const gradientColor = ref('linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)')
const activeKey = ref('pure')

const currentPreviewColor = computed(() =>
  pickerStateToColor(activeKey.value, pureColor.value, gradientColor.value)
)

const onPureColorChange = (val) => { pureColor.value = val }
const onGradientColorChange = (val) => { gradientColor.value = val }

const hideTagDropdown = () => { showTagDropdown.value = false }

const selectExistingTag = (t) => {
  inputTag.value = t.name
  const state = colorToPickerState(t.color)
  pureColor.value = state.pureColor
  gradientColor.value = state.gradientColor
  activeKey.value = state.activeKey
  showTagDropdown.value = false
}

// 手动输入标签名时，如果匹配已有标签，自动同步颜色
watch(inputTag, (newVal) => {
  const existing = store.allTags.find(t => t.name === newVal)
  if (existing) {
    const state = colorToPickerState(existing.color)
    pureColor.value = state.pureColor
    gradientColor.value = state.gradientColor
    activeKey.value = state.activeKey
  }
})

const submitArchive = async () => {
  if (!inputUrl.value) {
    errorMsg.value = '请输入链接'
    return
  }

  loading.value = true
  errorMsg.value = ''

  const finalColor = inputTag.value.trim()
    ? pickerStateToColor(activeKey.value, pureColor.value, gradientColor.value)
    : ''

  try {
    const res = await store.apiFetch(`${store.API_BASE}/api/archive`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        url: inputUrl.value.trim(),
        tag: inputTag.value.trim(),
        color: finalColor,
      })
    })

    if (!res.ok) {
      const errText = await res.text()
      throw new Error(errText || '归档请求失败')
    }

    inputUrl.value = ''
    inputTag.value = ''
    showColorPicker.value = false
    store.showAddModal = false
    store.showToast('已加入后台解析任务，请留意页面底部进度条')
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}
</script>
