<template>
  <div class="space-y-2">
    <!-- 预设色点区域 -->
    <div class="flex items-center gap-2 flex-wrap">
      <!-- 7 个纯色 -->
      <button
        v-for="c in SOLID_PRESETS"
        :key="c.value"
        @click="select(c.value)"
        class="w-5 h-5 rounded-full flex-shrink-0 transition-transform hover:scale-110 border border-black/10 shadow-sm"
        :class="isSelected(c.value) ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : ''"
        :style="{ background: c.value }"
        :title="c.label"
      ></button>

      <!-- 分隔线 -->
      <div class="w-px h-4 bg-gray-200 mx-0.5 flex-shrink-0"></div>

      <!-- 2 个渐变色 -->
      <button
        v-for="g in GRADIENT_PRESETS"
        :key="g.value"
        @click="select(g.value)"
        class="w-5 h-5 rounded-full flex-shrink-0 transition-transform hover:scale-110 border border-black/10 shadow-sm"
        :class="isSelected(g.value) ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : ''"
        :style="{ background: g.value }"
        :title="g.label"
      ></button>

      <!-- 分隔线 -->
      <div class="w-px h-4 bg-gray-200 mx-0.5 flex-shrink-0"></div>

      <!-- 展开自定义调色盘按钮：彩虹圆圈 -->
      <button
        @click="togglePicker"
        class="w-5 h-5 rounded-full flex-shrink-0 transition-transform hover:scale-110 border border-black/10 shadow-sm"
        :class="showPicker ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : ''"
        style="background: conic-gradient(red, yellow, lime, aqua, blue, magenta, red)"
        title="自定义颜色"
      ></button>
    </div>

    <!-- 展开的调色板面板 -->
    <div v-if="showPicker" class="rounded-xl border border-gray-100 bg-gray-50 p-2 flex flex-col items-center gap-1.5 shadow-inner">
      <ColorPicker
        isWidget
        :pureColor="pickerPureColor"
        :gradientColor="pickerGradientColor"
        :activeKey="pickerActiveKey"
        @update:pureColor="onPickerPureColor"
        @update:gradientColor="onPickerGradientColor"
        @update:activeKey="val => pickerActiveKey = val"
      />
      <button
        @click="showPicker = false"
        class="self-end text-xs text-gray-400 hover:text-brand px-2 py-0.5 rounded cursor-pointer transition"
      >收起 ↑</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ColorPicker } from 'vue3-colorpicker'
import 'vue3-colorpicker/style.css'
import { colorToPickerState, pickerStateToColor } from '../../utils/tagColor'

// ── 预设颜色 ─────────────────────────────────────────────────────────────────
const SOLID_PRESETS = [
  { label: '蓝色',   value: '#3b82f6' },
  { label: '红色',   value: '#ef4444' },
  { label: '绿色',   value: '#10b981' },
  { label: '紫色',   value: '#8b5cf6' },
  { label: '粉色',   value: '#ec4899' },
  { label: '橙色',   value: '#f97316' },
  { label: '青绿',   value: '#14b8a6' },
]

const GRADIENT_PRESETS = [
  { label: '落日',   value: 'linear-gradient(135deg, #f97316 0%, #ec4899 100%)' },
  { label: '深空',   value: 'linear-gradient(135deg, #667eea 0%, #4facfe 100%)' },
]

// ── Props / Emits ────────────────────────────────────────────────────────────
const props = defineProps({
  modelValue: { type: String, default: '#3b82f6' }
})
const emit = defineEmits(['update:modelValue'])

// ── 预设点选 ─────────────────────────────────────────────────────────────────
const isSelected = (val) => props.modelValue === val

const select = (val) => {
  emit('update:modelValue', val)
  // 同步 picker 状态，保持一致
  syncPickerState(val)
}

// ── 调色盘状态 ───────────────────────────────────────────────────────────────
const showPicker = ref(false)
const pickerPureColor = ref('#3b82f6')
const pickerGradientColor = ref('linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)')
const pickerActiveKey = ref('pure')

const syncPickerState = (colorVal) => {
  const state = colorToPickerState(colorVal)
  pickerPureColor.value = state.pureColor
  pickerGradientColor.value = state.gradientColor
  pickerActiveKey.value = state.activeKey
}

// 外部 modelValue 变化时同步 picker（比如选了现有标签）
watch(() => props.modelValue, (val) => {
  syncPickerState(val)
}, { immediate: true })

const togglePicker = () => {
  showPicker.value = !showPicker.value
}

const onPickerPureColor = (val) => {
  pickerPureColor.value = val
  if (pickerActiveKey.value === 'pure') {
    emit('update:modelValue', val)
  }
}

const onPickerGradientColor = (val) => {
  pickerGradientColor.value = val
  if (pickerActiveKey.value === 'gradient') {
    emit('update:modelValue', val)
  }
}

// activeKey 切换时，同步当前值给父组件
watch(pickerActiveKey, (key) => {
  const val = pickerStateToColor(key, pickerPureColor.value, pickerGradientColor.value)
  emit('update:modelValue', val)
})
</script>
