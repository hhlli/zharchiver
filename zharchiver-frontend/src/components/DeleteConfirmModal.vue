<template>
  <BaseModal :show="show" @close="$emit('cancel')">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between bg-red-50/50">
      <div class="flex items-center space-x-2 text-red-600">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <h3 class="font-bold">确认删除</h3>
      </div>
      <button @click="$emit('cancel')" class="text-gray-400 hover:text-gray-600 transition cursor-pointer">&times;</button>
    </div>
    
    <!-- Body -->
    <div class="px-6 py-5 text-sm text-gray-600 leading-relaxed">
      您确定要永久删除归档 <span class="font-bold text-gray-800">"{{ title }}"</span> 吗？
      <p class="mt-2 text-xs text-red-500 font-medium">此操作不可逆，将同时清除相关的本地图片文件和评论数据。</p>
    </div>
    
    <!-- Footer -->
    <div class="px-6 py-4 bg-gray-50 flex justify-end space-x-3 border-t border-gray-100">
      <button 
        @click="$emit('cancel')"
        class="px-4 py-1.5 rounded-lg text-sm text-gray-600 border border-gray-300 hover:bg-gray-100 transition cursor-pointer"
      >
        取消
      </button>
      <button 
        @click="$emit('confirm')"
        class="px-4 py-1.5 rounded-lg text-sm text-white bg-red-600 hover:bg-red-700 transition cursor-pointer font-medium shadow-sm"
      >
        确认删除
      </button>
    </div>
  </BaseModal>
</template>

<script setup>
import BaseModal from './common/BaseModal.vue'
defineProps({
  show: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: ''
  }
})

defineEmits(['confirm', 'cancel'])
</script>
