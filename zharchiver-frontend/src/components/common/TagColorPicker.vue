<template>
  <div class="space-y-2 relative">
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

      <!-- 展开自定义调色盘按钮 -->
      <button
        @click="openCustomPanel"
        class="w-5 h-5 rounded-full flex-shrink-0 transition-transform hover:scale-110 border border-black/10 shadow-sm"
        :class="isCustomSelected ? 'ring-2 ring-offset-2 ring-blue-400 scale-110' : ''"
        style="background: conic-gradient(red, yellow, lime, aqua, blue, magenta, red)"
        title="自定义颜色"
      ></button>
    </div>

    <!-- 自定义颜色设置面板 (绝对定位浮层) -->
    <div v-if="showCustomPanel" class="absolute z-50 top-8 left-0 rounded-xl border border-gray-200 bg-white p-3 flex flex-col gap-3 shadow-xl w-64">
      
      <!-- 模式切换 -->
      <div class="flex bg-gray-100 p-1 rounded-lg">
        <button 
          class="flex-1 text-xs py-1.5 rounded-md transition"
          :class="customMode === 'solid' ? 'bg-white shadow-sm font-medium text-gray-800' : 'text-gray-500 hover:text-gray-700'"
          @click="customMode = 'solid'"
        >纯色</button>
        <button 
          class="flex-1 text-xs py-1.5 rounded-md transition"
          :class="customMode === 'gradient' ? 'bg-white shadow-sm font-medium text-gray-800' : 'text-gray-500 hover:text-gray-700'"
          @click="customMode = 'gradient'"
        >双色渐变</button>
      </div>

      <!-- 纯色模式内容 -->
      <div v-if="customMode === 'solid'" class="flex items-center gap-3 py-2 px-1">
        <div class="relative w-8 h-8 rounded-md overflow-hidden shadow-sm border border-gray-200">
          <input type="color" v-model="solidColor" class="absolute -inset-4 w-16 h-16 cursor-pointer" title="点击选择颜色" />
        </div>
        <span class="text-xs text-gray-600 uppercase font-mono tracking-wider">{{ solidColor }}</span>
      </div>

      <!-- 渐变模式内容 -->
      <div v-if="customMode === 'gradient'" class="flex flex-col gap-3 pt-1">
        <div class="flex items-center justify-between gap-3">
          
          <!-- 颜色 1 -->
          <div class="flex flex-col items-center gap-1">
            <span class="text-[10px] text-gray-400">起点</span>
            <div class="relative w-7 h-7 rounded-md overflow-hidden shadow-sm border border-gray-200">
              <input type="color" v-model="gradColor1" class="absolute -inset-4 w-16 h-16 cursor-pointer" title="点击选择起点颜色" />
            </div>
          </div>
          
          <!-- 角度控制 -->
          <div class="flex-1 flex flex-col items-center gap-1">
            <span class="text-[10px] text-gray-400">角度: {{ gradAngle }}°</span>
            <input type="range" v-model="gradAngle" min="0" max="360" class="w-full h-1 bg-gray-200 rounded-lg appearance-none cursor-pointer" />
          </div>

          <!-- 颜色 2 -->
          <div class="flex flex-col items-center gap-1">
            <span class="text-[10px] text-gray-400">终点</span>
            <div class="relative w-7 h-7 rounded-md overflow-hidden shadow-sm border border-gray-200">
              <input type="color" v-model="gradColor2" class="absolute -inset-4 w-16 h-16 cursor-pointer" title="点击选择终点颜色" />
            </div>
          </div>
          
        </div>
        
        <!-- 渐变预览块 -->
        <div class="h-10 w-full rounded-lg shadow-inner border border-black/10 transition-all" :style="{ background: generatedGradient }"></div>
      </div>

      <!-- 底部操作按钮 -->
      <div class="flex justify-end gap-2 mt-1">
        <button @click="showCustomPanel = false" class="text-xs text-gray-500 hover:text-gray-700 px-3 py-1.5 rounded-md cursor-pointer transition">取消</button>
        <button @click="applyCustomColor" class="text-xs bg-brand text-white hover:bg-brand-hover px-3 py-1.5 rounded-md cursor-pointer transition font-medium shadow-sm">应用</button>
      </div>
      
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

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
  showCustomPanel.value = false // 选中预设时关闭自定义面板
}

// 判断当前是否选中了“非预设”的自定义颜色
const isCustomSelected = computed(() => {
  if (!props.modelValue) return false
  const isPreset = SOLID_PRESETS.some(p => p.value === props.modelValue) || 
                   GRADIENT_PRESETS.some(p => p.value === props.modelValue)
  return !isPreset
})

// ── 自定义面板状态 ────────────────────────────────────────────────────────────
const showCustomPanel = ref(false)
const customMode = ref('solid') // 'solid' | 'gradient'
const solidColor = ref('#3b82f6')
const gradColor1 = ref('#4facfe')
const gradColor2 = ref('#00f2fe')
const gradAngle = ref(135)

const generatedGradient = computed(() => {
  return `linear-gradient(${gradAngle.value}deg, ${gradColor1.value} 0%, ${gradColor2.value} 100%)`
})

// 点击彩虹圆点：打开面板并初始化状态
const openCustomPanel = () => {
  if (showCustomPanel.value) {
    showCustomPanel.value = false
    return
  }

  // 根据当前绑定的颜色初始化面板状态
  const current = props.modelValue
  if (current && current.startsWith('linear-gradient')) {
    customMode.value = 'gradient'
    // 尝试正则解析现有的渐变角度和两个颜色
    const match = current.match(/linear-gradient\((\d+)deg,\s*(#[0-9a-fA-F]{6})\s*0%,\s*(#[0-9a-fA-F]{6})\s*100%\)/)
    if (match) {
      gradAngle.value = parseInt(match[1], 10)
      gradColor1.value = match[2]
      gradColor2.value = match[3]
    } else {
      // 无法解析时给默认值
      gradAngle.value = 135
      gradColor1.value = '#4facfe'
      gradColor2.value = '#00f2fe'
    }
  } else if (current && current.startsWith('#') && current.length === 7) {
    customMode.value = 'solid'
    solidColor.value = current
  } else {
    // 兜底默认值
    customMode.value = 'solid'
    solidColor.value = '#3b82f6'
  }
  
  showCustomPanel.value = true
}

const applyCustomColor = () => {
  if (customMode.value === 'solid') {
    emit('update:modelValue', solidColor.value)
  } else {
    emit('update:modelValue', generatedGradient.value)
  }
  showCustomPanel.value = false
}
</script>
