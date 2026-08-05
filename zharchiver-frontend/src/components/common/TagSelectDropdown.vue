<template>
  <div class="relative inline-block text-left" ref="dropdownRef">
    <!-- 触发按钮 -->
    <button 
      @click="isOpen = !isOpen"
      class="md:hidden flex items-center justify-center w-8 h-8 bg-gray-100/80 hover:bg-gray-200/80 text-gray-600 rounded-lg border-none focus:outline-none transition-colors"
      title="选择分类"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
    </button>

    <!-- 展开菜单 -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 translate-y-1 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-1 scale-95"
    >
      <div v-if="isOpen" class="absolute right-0 mt-2 w-32 rounded-2xl shadow-xl bg-white/95 backdrop-blur-md border border-gray-100 focus:outline-none z-50 overflow-hidden py-2">
        <button
          @click="selectOption('all')"
          :class="['block w-full text-left px-4 py-2 text-sm font-medium transition-colors', modelValue === 'all' ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >
          所有
        </button>
        
        <div v-if="tags && tags.length > 0" class="border-t border-gray-100 my-1 mx-2"></div>

        <button
          v-for="tag in tags"
          :key="tag.name"
          @click="selectOption(tag.name)"
          :class="['block w-full text-left px-4 py-2 text-sm font-medium transition-colors truncate', modelValue === tag.name ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50']"
        >
          {{ tag.name }}
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  tags: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const dropdownRef = ref(null)

const selectOption = (val) => {
  emit('update:modelValue', val)
  isOpen.value = false
}

const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
  document.addEventListener('touchstart', handleClickOutside)
})
onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
  document.removeEventListener('touchstart', handleClickOutside)
})
</script>
