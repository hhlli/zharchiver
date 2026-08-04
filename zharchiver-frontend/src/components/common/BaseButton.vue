<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center px-4 py-2 text-sm font-medium rounded-lg transition-colors shadow-sm focus:outline-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed',
      variantClasses[variant],
      customClass
    ]"
    @click="$emit('click')"
  >
    <svg 
      v-if="loading" 
      class="animate-spin -ml-1 mr-2 h-4 w-4" 
      :class="variant === 'outline' ? 'text-gray-500' : 'text-current'" 
      fill="none" 
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot></slot>
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'button'
  },
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'danger', 'outline', 'soft', 'danger-soft'].includes(value)
  },
  loading: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  customClass: {
    type: String,
    default: 'w-full sm:w-auto'
  }
})

defineEmits(['click'])

const variantClasses = {
  'primary': 'bg-blue-600 hover:bg-blue-700 text-white',
  'secondary': 'bg-gray-800 hover:bg-gray-900 text-white',
  'danger': 'bg-red-600 hover:bg-red-700 text-white',
  'outline': 'bg-white hover:bg-gray-50 text-gray-700 border border-gray-300',
  'soft': 'bg-blue-50 hover:bg-blue-100 text-blue-700 border-transparent',
  'danger-soft': 'bg-red-50 hover:bg-red-100 text-red-700 border-transparent'
}
</script>
