<template>
  <div class="relative inline-block text-left" ref="dropdownRef">
    <!-- 触发按钮 -->
    <button 
      @click="isOpen = !isOpen"
      class="md:hidden flex items-center justify-between w-20 sm:w-24 bg-gray-100/80 hover:bg-gray-200/80 text-xs text-gray-700 py-1.5 px-2.5 rounded-lg border-none focus:outline-none transition-colors"
    >
      <span class="truncate font-medium">{{ currentLabel }}</span>
      <svg class="w-3 h-3 text-gray-500 flex-shrink-0 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
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
      <div v-if="isOpen" class="absolute right-0 mt-2 w-36 rounded-2xl shadow-xl bg-white/95 backdrop-blur-md border border-gray-100 focus:outline-none z-50 overflow-hidden py-1.5">
        <button
          @click="selectOption('all')"
          :class="['block w-full text-left px-4 py-2 text-sm transition-colors', modelValue === 'all' ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-50']"
        >
          <div class="flex items-center justify-between">
            <span>所有</span>
            <svg v-if="modelValue === 'all'" class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
          </div>
        </button>
        
        <div v-if="tags && tags.length > 0" class="border-t border-gray-100 my-1 mx-2"></div>

        <button
          v-for="tag in tags"
          :key="tag.name"
          @click="selectOption(tag.name)"
          :class="['block w-full text-left px-4 py-2 text-sm transition-colors truncate', modelValue === tag.name ? 'bg-blue-50 text-blue-600 font-medium' : 'text-gray-700 hover:bg-gray-50']"
        >
          <div class="flex items-center justify-between">
            <span class="truncate mr-2">{{ tag.name }}</span>
            <svg v-if="modelValue === tag.name" class="w-4 h-4 text-blue-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
          </div>
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

const currentLabel = computed(() => {
  if (props.modelValue === 'all') return '所有'
  const tag = props.tags.find(t => t.name === props.modelValue)
  return tag ? tag.name : '所有'
})

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
