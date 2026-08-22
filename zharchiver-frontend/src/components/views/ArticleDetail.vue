<template>
  <div class="max-w-6xl mx-auto space-y-6 md:h-full pb-8">
    <button 
      @click="store.goBackToList()" 
      class="inline-flex items-center space-x-1 text-sm md:text-xs font-medium text-muted hover:text-brand mb-2 md:mb-2 cursor-pointer transition select-none py-2 px-1 -ml-1 md:py-0 md:px-0 md:ml-0"
    >
      <svg class="w-5 h-5 md:w-4 md:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
      <span>返回列表</span>
    </button>

    <header class="border-b pb-4 space-y-2">
      <h1 v-if="!isEditing" class="text-xl md:text-2xl font-bold text-primary leading-snug">{{ store.currentAnswer?.title }}</h1>
      <textarea 
        v-else 
        v-model="editTitle" 
        ref="titleTextarea"
        @input="autoResizeTitle"
        rows="1"
        class="w-full text-xl md:text-2xl font-bold text-primary leading-snug bg-transparent resize-none focus:outline-none overflow-hidden" 
        placeholder="文章标题" 
      ></textarea>
      <div class="flex flex-col sm:flex-row sm:items-center text-xs text-muted space-y-2 sm:space-y-0">
        <div class="flex items-center">
          <span class="font-medium text-secondary">归档：{{ formatDate(store.currentAnswer?.saved_at) }}</span>
        </div>
        
        <div class="sm:ml-auto flex items-center space-x-3">
          <div class="flex items-center text-xs">
            <span 
              @click="startEditTag" 
              class="border px-2 py-0.5 rounded-md font-medium tracking-wider cursor-pointer select-none hover:opacity-80 transition"
              :style="resolveTagStyle(store.currentAnswer?.tag_color)"
              title="点击编辑标签"
            >
              {{ store.currentAnswer?.tag || 'ANSWER' }}
            </span>
          </div>
          <div v-if="isEditing" class="flex items-center space-x-2 ml-4 border-l pl-4 border-line">
            <button @click="cancelEdit" class="text-xs text-muted hover:text-secondary transition px-2 py-1 cursor-pointer">取消</button>
            <button @click="saveEdit" class="text-xs bg-brand text-white px-3 py-1 rounded-md hover:bg-brand-hover transition flex items-center space-x-1 shadow-sm cursor-pointer" :disabled="saving">
              <span v-if="saving">保存中...</span>
              <span v-else>保存</span>
            </button>
          </div>

          <AnswerActions v-else :answerId="store.currentAnswer?.answer_id" @commentAdded="onCommentAdded" @edit="enterEditMode" />
        </div>
      </div>
    </header>

    <!-- 同问题多回答的作者切换标签 -->
    <div v-if="store.currentGroup && store.currentGroup.count > 1" class="flex flex-wrap gap-2 my-2">
      <button
        v-for="ans in store.currentGroup.answers"
        :key="ans.answer_id"
        @click="store.selectAnswer(ans.answer_id, true)"
        :class="['px-3 py-1.5 rounded-full text-xs font-medium transition-colors cursor-pointer border', store.currentAnswer?.answer_id === ans.answer_id ? 'bg-blue-50 dark:bg-blue-900/30 border-blue-200 dark:border-blue-800 text-blue-600 dark:text-blue-400 shadow-sm' : 'bg-surface-hover border-line text-muted hover:bg-surface-hover']"
      >
        {{ ans.author_name || '匿名用户' }}
      </button>
    </div>

    <article v-if="!isEditing"
      class="prose prose-sm md:prose-base prose-slate max-w-none text-primary leading-relaxed prose-img:max-w-full prose-img:rounded-lg prose-img:shadow-sm"
      v-html="processHtmlContent(store.currentAnswer?.content_html)"
    ></article>
    <div v-else class="py-2 animate-in fade-in slide-in-from-bottom-2 duration-300">
      <RichEditor v-model="editContent" />
    </div>

    <div class="mt-8 pt-3 border-t border-line-light flex flex-col sm:flex-row sm:items-center justify-between text-xs text-muted mb-3 gap-3">
      <div class="flex items-center space-x-4">
        <span class="font-medium text-secondary">作者：{{ store.currentAnswer?.author_name }}</span>
        <span>发布于 {{ formatTimestamp(store.currentAnswer?.created_time) }}</span>
      </div>
      <div class="flex items-center space-x-3">
        <!-- 添加评论按钮 -->
        <button 
          @click="showCommentModal = true"
          class="inline-flex items-center space-x-1 text-muted hover:text-green-600 transition-colors cursor-pointer"
          title="添加评论"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"></path></svg>
          <span>添加评论</span>
        </button>
        <div class="w-px h-3 bg-line"></div>

        <a :href="getOriginalUrl(store.currentAnswer)" target="_blank" rel="noopener noreferrer" class="inline-flex items-center space-x-1 text-muted hover:text-brand transition-colors cursor-pointer">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
          <span></span>
        </a>
        <div class="w-px h-3 bg-line"></div>
        <button 
          @click="store.itemToDelete = store.currentAnswer; store.itemToDeleteType = 'answer'" 
          class="inline-flex items-center space-x-1 text-muted hover:text-red-500 transition-colors cursor-pointer"
          title="删除此回答"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
          <span></span>
        </button>
      </div>
    </div>

    <AnswerComments :answerId="store.currentAnswer?.answer_id" :refreshKey="commentRefreshKey" />

    <!-- 编辑标签弹窗 -->
    <BaseModal 
      :show="isEditingTag" 
      @close="isEditingTag = false"
      maxWidthClass="max-w-sm"
      contentClass="p-6 space-y-4"
      zIndexClass="z-[60]"
    >
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-primary">编辑标签</h3>
        <button @click="isEditingTag = false" class="text-gray-400 hover:text-secondary text-lg cursor-pointer">&times;</button>
      </div>
      
      <div class="space-y-4 pt-2">
        <div class="space-y-2">
          <label class="text-xs text-muted">标签名称</label>
          <div class="relative">
            <input 
              type="text" 
              v-model="editTagValue" 
              class="w-full border border-line rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 pr-8" 
              @keyup.enter="saveTag"
              @focus="showTagDropdown = true"
              @blur="hideTagDropdown"
            >
            <div class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
            </div>
            
            <div v-if="showTagDropdown && store.allTags.length > 0" class="absolute z-50 w-full mt-1 bg-surface border border-line rounded-lg shadow-lg max-h-48 overflow-y-auto">
              <div 
                v-for="t in store.allTags" 
                :key="t.name" 
                @mousedown.prevent="selectExistingTag(t)"
                class="px-3 py-2 text-sm text-secondary hover:bg-blue-50 dark:hover:bg-blue-900/30 cursor-pointer flex items-center space-x-2"
              >
                <span class="w-2 h-2 rounded-full flex-shrink-0 shadow-[0_0_2px_rgba(0,0,0,0.12)]" :style="{ background: resolveTagBackground(t.color) }"></span>
                <span class="truncate">{{ t.name }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-4">
          <label class="text-xs text-muted block mb-2">标签颜色</label>
          <TagColorPicker v-model="tagColor" />
        </div>
      </div>

      <div class="flex justify-end space-x-2 pt-4">
        <button 
          @click="isEditingTag = false"
          class="px-4 py-1.5 border border-line text-secondary rounded-lg text-xs hover:bg-surface-hover transition cursor-pointer"
        >取消</button>
        <button 
          @click="saveTag"
          class="px-4 py-1.5 bg-brand text-white rounded-lg text-xs font-medium hover:bg-brand-hover transition cursor-pointer"
        >保存</button>
      </div>
    </BaseModal>

    <!-- 添加评论弹窗 -->
    <BaseModal 
      :show="showCommentModal" 
      @close="showCommentModal = false"
      maxWidthClass="max-w-lg"
      contentClass="p-6 space-y-4"
      zIndexClass="z-[70]"
    >
      <div class="flex items-center justify-between border-b pb-3">
        <h3 class="text-sm font-semibold text-primary">添加评论或笔记</h3>
        <button @click="showCommentModal = false" class="text-gray-400 hover:text-secondary text-lg cursor-pointer">&times;</button>
      </div>
      
      <div class="space-y-2">
        <textarea 
          v-model="commentContent"
          rows="4"
          placeholder="在这里粘贴您的评论或学习笔记..."
          class="w-full border border-line rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
          :disabled="isSubmittingComment"
        ></textarea>
      </div>

      <div class="flex justify-end space-x-2 pt-2">
        <button 
          @click="showCommentModal = false"
          class="px-4 py-1.5 border border-line text-secondary rounded-lg text-xs hover:bg-surface-hover transition cursor-pointer"
          :disabled="isSubmittingComment"
        >
          取消
        </button>
        <button 
          @click="submitComment"
          class="px-4 py-1.5 bg-brand text-white rounded-lg text-xs font-medium hover:bg-brand-hover transition disabled:bg-gray-400 cursor-pointer"
          :disabled="isSubmittingComment || !commentContent.trim()"
        >
          {{ isSubmittingComment ? '提交中...' : '提交评论' }}
        </button>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useArchiveStore } from '../../stores/archive'
import AnswerActions from '../AnswerActions.vue'
import AnswerComments from '../AnswerComments.vue'
import BaseModal from '../common/BaseModal.vue'
import RichEditor from '../RichEditor.vue'
import DOMPurify from 'dompurify'
import TagColorPicker from '../common/TagColorPicker.vue'
import { resolveTagBackground, resolveTagStyle } from '../../utils/tagColor'

const store = useArchiveStore()
// ── 文章编辑 ──────────────────────────────────────────────────────────────────
const isEditing = ref(false)
const editTitle = ref('')
const editContent = ref('')
const saving = ref(false)
const titleTextarea = ref(null)

const autoResizeTitle = () => {
  if (titleTextarea.value) {
    titleTextarea.value.style.height = 'auto'
    titleTextarea.value.style.height = titleTextarea.value.scrollHeight + 'px'
  }
}

const enterEditMode = async () => {
  const { nextTick } = await import('vue')
  editTitle.value = store.currentAnswer?.title || ''
  editContent.value = processHtmlContent(store.currentAnswer?.content_html)
  isEditing.value = true
  await nextTick()
  autoResizeTitle()
}

const cancelEdit = () => { isEditing.value = false }

const saveEdit = async () => {
  if (!store.currentAnswer) return
  saving.value = true
  try {
    const cleanContent = editContent.value.replace(
      new RegExp(`\\?token=${localStorage.getItem('token') || ''}`, 'g'), ''
    )
    const res = await store.apiFetch(`/api/answers/${store.currentAnswer.answer_id}/content`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: editTitle.value.trim(), content_html: cleanContent })
    })
    if (!res.ok) throw new Error('保存失败')
    store.currentAnswer.title = editTitle.value.trim()
    store.currentAnswer.content_html = cleanContent
    isEditing.value = false
    store.fetchAnswersList()
  } catch (err) {
    store.showToast(err.message, 'error')
  } finally {
    saving.value = false
  }
}

// ── 评论 ──────────────────────────────────────────────────────────────────────
const commentRefreshKey = ref(0)
const onCommentAdded = () => { commentRefreshKey.value++ }

const showCommentModal = ref(false)
const commentContent = ref('')
const isSubmittingComment = ref(false)

const submitComment = async () => {
  if (!commentContent.value.trim() || !store.currentAnswer) return
  isSubmittingComment.value = true

  try {
    const res = await store.apiFetch(`/api/answers/${store.currentAnswer.answer_id}/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: commentContent.value.trim() })
    })

    if (res.ok) {
      showCommentModal.value = false
      commentContent.value = ''
      onCommentAdded()
    } else {
      store.showToast('评论提交失败', 'error')
    }
  } catch (err) {
    console.error(err)
    store.showToast('网络请求失败', 'error')
  } finally {
    isSubmittingComment.value = false
  }
}

// ── 标签编辑 ──────────────────────────────────────────────────────────────────
const isEditingTag = ref(false)
const editTagValue = ref('')
const showTagDropdown = ref(false)
const tagColor = ref('#3b82f6')

const hideTagDropdown = () => { showTagDropdown.value = false }

const selectExistingTag = (t) => {
  editTagValue.value = t.name
  tagColor.value = resolveTagBackground(t.color)
  showTagDropdown.value = false
}

watch(editTagValue, (newVal) => {
  const existing = store.allTags.find(t => t.name === newVal)
  if (existing) tagColor.value = resolveTagBackground(existing.color)
})

const startEditTag = () => {
  editTagValue.value = store.currentAnswer?.tag || ''
  tagColor.value = resolveTagBackground(store.currentAnswer?.tag_color)
  isEditingTag.value = true
}

const saveTag = async () => {
  if (!store.currentAnswer) return
  try {
    const res = await store.apiFetch(`/api/answers/${store.currentAnswer.answer_id}/tag`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ tag: editTagValue.value.trim(), color: tagColor.value })
    })
    if (!res.ok) throw new Error('标签更新失败')
    store.currentAnswer.tag = editTagValue.value.trim()
    store.currentAnswer.tag_color = tagColor.value
    isEditingTag.value = false
    store.fetchAnswersList()
  } catch (err) {
    store.showToast(err.message, 'error')
  }
}

// ── 内容处理 ──────────────────────────────────────────────────────────────────
const processHtmlContent = (html) => {
  if (!html) return ''
  const token = localStorage.getItem('token') || ''
  const processed = html.replace(/src="\/storage\/([^"]+)"/g, (match, p1) => {
    const url = `${store.API_BASE}/storage/${p1}`
    return token ? `src="${url}?token=${token}"` : `src="${url}"`
  })
  return DOMPurify.sanitize(processed, { ADD_TAGS: ['video', 'source'], ADD_ATTR: ['controls'] })
}

const formatTimestamp = (ts) => ts ? new Date(ts * 1000).toLocaleString() : ''
const formatDate = (dateStr) => dateStr ? new Date(dateStr).toLocaleDateString() : ''

const getOriginalUrl = (answer) => {
  if (!answer) return '#'
  if (answer.question_id === 'twitter') return `https://x.com/i/status/${answer.answer_id}`
  return `https://www.zhihu.com/question/${answer.question_id}/answer/${answer.answer_id}`
}
</script>
