<template>
  <BaseModal :show="store.itemToDelete !== null" @close="store.itemToDelete = null">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-line-light flex items-center justify-between bg-red-50/50 dark:bg-red-900/20">
      <div class="flex items-center space-x-2 text-red-600">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <h3 class="font-bold">确认删除</h3>
      </div>
      <button @click="store.itemToDelete = null" class="text-gray-400 hover:text-secondary transition cursor-pointer">&times;</button>
    </div>
    
    <!-- Body -->
    <div class="px-6 py-5 text-sm text-secondary leading-relaxed">
      <template v-if="store.itemToDeleteType === 'group'">
        您确定要永久删除该问题下的 <span class="font-bold text-primary">所有 {{ store.itemToDelete?.count }} 条</span> 归档回答吗？
        <p class="mt-2 text-xs text-red-500 font-medium">此操作不可逆，将同时清除相关的本地图片文件和评论数据。</p>
      </template>
      <template v-else-if="store.itemToDeleteType === 'answer'">
        您确定要永久删除归档 <span class="font-bold text-primary">"{{ store.itemToDelete?.title }}"</span> 吗？
        <p class="mt-2 text-xs text-red-500 font-medium">此操作不可逆，将同时清除相关的本地图片文件和评论数据。</p>
      </template>
      <template v-else-if="store.itemToDeleteType === 'comment'">
        您确定要永久删除这条评论/笔记吗？
        <p class="mt-2 text-xs text-red-500 font-medium">此操作不可逆，删除后无法恢复。</p>
      </template>
    </div>
    
    <!-- Footer -->
    <div class="px-6 py-4 bg-surface-hover flex justify-end space-x-3 border-t border-line-light">
      <button 
        @click="store.itemToDelete = null"
        class="px-4 py-1.5 rounded-lg text-sm text-secondary border border-line hover:bg-surface-hover transition cursor-pointer"
      >
        取消
      </button>
      <button 
        @click="store.confirmDelete"
        class="px-4 py-1.5 rounded-lg text-sm text-white bg-red-600 hover:bg-red-700 transition cursor-pointer font-medium shadow-sm"
      >
        确认删除
      </button>
    </div>
  </BaseModal>
</template>

<script setup>
import { useArchiveStore } from '../stores/archive'
import BaseModal from './common/BaseModal.vue'

const store = useArchiveStore()
</script>
