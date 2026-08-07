<template>
  <span 
    :class="[
      'inline-flex items-center justify-center font-semibold tracking-wider uppercase',
      sizeClasses[size] || sizeClasses.sm,
      variantClasses[variant] || variantClasses.default
    ]"
  >
    <slot></slot>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'primary', 'success', 'warning', 'danger'].includes(value)
  },
  size: {
    type: String,
    default: 'sm',
    validator: (value) => ['xs', 'sm', 'md'].includes(value)
  }
})

const sizeClasses = {
  xs: 'text-[9px] px-1.5 py-0.5 rounded',
  sm: 'text-mini px-2 py-0.5 rounded',
  md: 'text-xs px-2.5 py-1 rounded-md'
}

const variantClasses = {
  default: 'bg-transparent text-gray-400',
  primary: 'bg-brand/10 text-brand',
  success: 'bg-green-100 text-green-700',
  warning: 'bg-yellow-100 text-yellow-700',
  danger: 'bg-red-100 text-red-700'
}
</script>
