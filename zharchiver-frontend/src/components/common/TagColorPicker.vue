<template>
  <div class="space-y-2">
    <!-- 预设色点区域 -->
    <div class="flex items-center gap-2 flex-wrap">
      <!-- 7 个纯色预设 -->
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
      <div class="w-px h-4 bg-line mx-0.5 flex-shrink-0"></div>

      <!-- 展开原生自定义调色盘按钮 -->
      <label
        class="relative w-5 h-5 rounded-full flex-shrink-0 transition-transform hover:scale-110 border border-black/10 shadow-sm cursor-pointer block overflow-hidden"
        :class="isCustomSelected ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : ''"
        style="background: conic-gradient(red, yellow, lime, aqua, blue, magenta, red)"
        title="自定义颜色"
      >
        <input 
          type="color" 
          :value="customColor"
          @input="onCustomColorChange"
          class="absolute -inset-4 w-12 h-12 opacity-0 cursor-pointer"
        />
      </label>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

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

// ── Props / Emits ────────────────────────────────────────────────────────────
const props = defineProps({
  modelValue: { type: String, default: '#3b82f6' }
})
const emit = defineEmits(['update:modelValue'])

// ── 预设点选 ─────────────────────────────────────────────────────────────────
const isSelected = (val) => props.modelValue === val

const select = (val) => {
  emit('update:modelValue', val)
}

// 判断当前是否选中了“非预设”的自定义颜色
const isCustomSelected = computed(() => {
  if (!props.modelValue) return false
  const isPreset = SOLID_PRESETS.some(p => p.value === props.modelValue)
  return !isPreset
})

// 原生 input type="color" 只能接受 #RRGGBB 格式。
// 兜底一个默认颜色防止原生 input 报错。
const customColor = computed(() => {
  if (props.modelValue && props.modelValue.startsWith('#') && props.modelValue.length === 7) {
    return props.modelValue
  }
  return '#3b82f6' // 回退到默认蓝色
})

const onCustomColorChange = (e) => {
  emit('update:modelValue', e.target.value)
}
</script>
