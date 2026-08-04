<template>
  <div class="h-full bg-gray-50 overflow-y-auto">
    <div class="max-w-3xl mx-auto py-10 px-8">
      <h2 class="text-2xl font-semibold text-gray-800 mb-8">设置</h2>

      <!-- 模块 1: 账号与鉴权 -->
      <div class="mb-10">
        <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 px-1">账号与鉴权</h3>
        <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm">
          <div class="p-4 flex items-center justify-between">
            <div class="flex-col pr-4">
              <span class="block text-sm font-medium text-gray-700">知乎 Cookie (d_c0)</span>
              <span class="block text-xs text-gray-400 mt-1">配置此参数用于绕过未登录限制，抓取高清图片。修改后下一次归档生效。</span>
            </div>
            <div class="w-72 shrink-0">
              <input 
                type="text" 
                v-model="settings.cookie" 
                class="w-full text-sm border border-gray-300 rounded-md px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-shadow" 
                placeholder="输入 d_c0 值"
              >
            </div>
          </div>
        </div>
      </div>

      <!-- 模块 2: 数据与存储 -->
      <div class="mb-10">
        <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3 px-1">数据与存储</h3>
        <div class="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm divide-y divide-gray-100">
          <div class="p-4 flex items-center justify-between">
            <div class="flex-col">
              <span class="block text-sm font-medium text-gray-700">清理本地图片</span>
              <span class="block text-xs text-gray-400 mt-1">扫描并删除已移除归档遗留的孤立图片文件，释放磁盘空间。</span>
            </div>
            <button class="px-4 py-1.5 bg-gray-50 hover:bg-gray-100 text-gray-700 text-sm font-medium rounded-md transition-colors border border-gray-200 shadow-sm">
              开始清理
            </button>
          </div>
          <div class="p-4 flex items-center justify-between">
            <div class="flex-col">
              <span class="block text-sm font-medium text-gray-700">重建数据库索引</span>
              <span class="block text-xs text-gray-400 mt-1">重新生成 SQLite 虚拟表索引，优化检索性能。</span>
            </div>
            <button class="px-4 py-1.5 bg-gray-50 hover:bg-gray-100 text-gray-700 text-sm font-medium rounded-md transition-colors border border-gray-200 shadow-sm">
              重建索引
            </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const settings = reactive({
  cookie: 'kdOYbtmYmhyPTuJmX8VHlc5ERZxJtktiX88=|1784191170'
})
</script>