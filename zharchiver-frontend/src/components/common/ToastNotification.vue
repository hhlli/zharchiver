<template>
  <Transition name="toast">
    <div 
      v-if="store.toastVisible" 
      class="fixed bottom-10 left-1/2 transform -translate-x-1/2 z-50 flex items-center px-4 py-2.5 bg-white border border-gray-100 rounded-lg shadow-sm font-medium text-sm text-gray-700 space-x-2"
    >
      <!-- Success Icon -->
      <svg v-if="store.toastType === 'success'" class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"></path>
      </svg>
      <!-- Error Icon -->
      <svg v-else-if="store.toastType === 'error'" class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12"></path>
      </svg>
      <!-- Info Icon -->
      <svg v-else class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <span>{{ store.toastMessage }}</span>
    </div>
  </Transition>
</template>

<script setup>
import { useArchiveStore } from '../../stores/archive'

const store = useArchiveStore()
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}
</style>
