<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 flex items-center justify-center p-4" :class="zIndexClass">
      <div 
        class="absolute inset-0 bg-black/40 backdrop-blur-sm" 
        @click="closeOnOutside ? $emit('close') : null"
      ></div>
      <div :class="['relative bg-surface rounded-2xl shadow-2xl w-full', maxWidthClass, contentClass]">
        <slot></slot>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
defineProps({
  show: { 
    type: Boolean, 
    required: true 
  },
  zIndexClass: { 
    type: String, 
    default: 'z-[80]' 
  },
  maxWidthClass: { 
    type: String, 
    default: 'max-w-sm' 
  },
  contentClass: { 
    type: String, 
    default: 'overflow-hidden transform transition-all' 
  },
  closeOnOutside: { 
    type: Boolean, 
    default: false 
  }
})

defineEmits(['close'])
</script>
