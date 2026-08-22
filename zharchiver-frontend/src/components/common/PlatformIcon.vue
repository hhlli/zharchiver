<template>
  <div class="flex-shrink-0 flex items-center justify-center">
    <!-- Twitter / X -->
    <svg v-if="platform === 'twitter'" :class="[svgClass, 'text-primary']" viewBox="0 0 24 24" fill="currentColor" stroke="none">
      <path stroke="none" d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 24.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
    </svg>
    <!-- Zhihu -->
    <svg v-else :class="[svgClass, 'text-blue-500']" viewBox="0 0 24 24">
      <rect width="24" height="24" rx="5" fill="currentColor"/>
      <text x="12" y="17" fill="white" font-size="15" font-weight="bold" font-family="sans-serif" text-anchor="middle">知</text>
    </svg>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  questionId: {
    type: String,
    required: true
  },
  svgClass: {
    type: String,
    default: 'w-5 h-5'
  }
})

const platform = computed(() => {
  if (props.questionId === 'twitter') return 'twitter'
  return 'zhihu'
})
</script>
