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

      <!-- 自定义颜色：使用原生 input type="color" -->
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
}

// ── 自定义颜色 ───────────────────────────────────────────────────────────────
// 判断当前是否选中了“非预设”的颜色（即自定义颜色）
const isCustomSelected = computed(() => {
  if (!props.modelValue) return false
  const isPreset = SOLID_PRESETS.some(p => p.value === props.modelValue) || 
                   GRADIENT_PRESETS.some(p => p.value === props.modelValue)
  return !isPreset
})

// 原生 input type="color" 只能接受 #RRGGBB 格式。
// 如果当前颜色是渐变色或无效格式，我们给一个默认的黑色兜底，防止原生 input 报错。
const customColor = computed(() => {
  if (props.modelValue && props.modelValue.startsWith('#') && props.modelValue.length === 7) {
    return props.modelValue
  }
  return '#000000'
})

const onCustomColorChange = (e) => {
  emit('update:modelValue', e.target.value)
}
</script>
