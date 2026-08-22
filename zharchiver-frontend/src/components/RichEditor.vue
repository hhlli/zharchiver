<template>
  <div class="border border-line rounded-xl overflow-hidden bg-surface flex flex-col shadow-sm">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 p-2 border-b border-line-light bg-surface-hover/80" v-if="editor">
      <button @click="editor.chain().focus().toggleBold().run()" :class="['p-1.5 rounded-md hover:bg-line text-secondary transition', editor.isActive('bold') ? 'bg-line text-primary' : '']" title="粗体">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z"></path></svg>
      </button>
      <button @click="editor.chain().focus().toggleItalic().run()" :class="['p-1.5 rounded-md hover:bg-line text-secondary transition', editor.isActive('italic') ? 'bg-line text-primary' : '']" title="斜体">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 4h8 M14 4l-4 16 M6 20h8"></path></svg>
      </button>
      <button @click="editor.chain().focus().toggleStrike().run()" :class="['p-1.5 rounded-md hover:bg-line text-secondary transition', editor.isActive('strike') ? 'bg-line text-primary' : '']" title="删除线">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14 M10 4a4 4 0 00-4 4 4 4 0 004 4 M14 20a4 4 0 004-4 4 4 0 00-4-4"></path></svg>
      </button>
      
      <button @click="editor.chain().focus().unsetAllMarks().clearNodes().run()" class="p-1.5 rounded-md hover:bg-line text-secondary transition" title="清除格式">
        <!-- SVG for clear formatting -->
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 2l-7 7 7 7 7-7-7-7z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 22h18"></path></svg>
      </button>

      <div class="w-px h-4 bg-line mx-1"></div>
      
      <div class="flex items-center space-x-1 ml-1">
        <select 
          :value="editor.getAttributes('textStyle').fontSize || '16px'"
          @change="e => e.target.value ? editor.chain().focus().setFontSize(e.target.value).run() : editor.chain().focus().unsetFontSize().run()" 
          class="text-xs border border-line rounded px-1.5 py-1 bg-surface focus:outline-none text-secondary cursor-pointer font-sans"
          title="字号"
        >
          <option value="12px">12</option>
          <option value="14px">14</option>
          <option value="16px">16</option>
          <option value="18px">18</option>
          <option value="20px">20</option>
          <option value="24px">24</option>
          <option value="36px">36</option>
        </select>
      </div>

      <div class="flex items-center relative mx-1">
        <div class="relative group flex items-center justify-center p-1 rounded-md hover:bg-line transition" title="字体颜色">
          <!-- Text Color Icon ('A') -->
          <svg class="w-4 h-4 text-secondary" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 19l5-14 5 14M9.5 13h5"></path></svg>
          <input 
            type="color" 
            class="absolute inset-0 opacity-0 cursor-pointer w-full h-full"
            :value="editor.getAttributes('textStyle').color || '#000000'"
            @input="e => editor.chain().focus().setColor(e.target.value).run()"
          />
          <div class="absolute bottom-0.5 left-1.5 right-1.5 h-0.5 rounded-sm" :style="{ backgroundColor: editor.getAttributes('textStyle').color || '#000000' }"></div>
        </div>
      </div>
      
      <div class="w-px h-4 bg-line mx-1"></div>
      
      <button @click="editor.chain().focus().undo().run()" :disabled="!editor.can().undo()" class="p-1.5 rounded-md hover:bg-line text-secondary transition disabled:opacity-30 disabled:hover:bg-transparent" title="撤销">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2 M3 10l6 6 M3 10l6-6"></path></svg>
      </button>
      <button @click="editor.chain().focus().redo().run()" :disabled="!editor.can().redo()" class="p-1.5 rounded-md hover:bg-line text-secondary transition disabled:opacity-30 disabled:hover:bg-transparent" title="重做">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10h-10a8 8 0 00-8 8v2 M21 10l-6 6 M21 10l-6-6"></path></svg>
      </button>
    </div>

    <!-- Editor Content -->
    <div class="flex-1 overflow-y-auto p-4 md:p-6 cursor-text min-h-[400px] max-h-[70vh] bg-surface" @click="editor?.commands.focus()">
      <editor-content :editor="editor" class="prose prose-sm md:prose-base prose-slate max-w-none text-primary leading-relaxed prose-img:max-w-full prose-img:rounded-lg prose-img:shadow-sm focus:outline-none" />
    </div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, watch } from 'vue'
import { Extension } from '@tiptap/core'
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import { TextStyle } from '@tiptap/extension-text-style'
import { Color } from '@tiptap/extension-color'

const FontSize = Extension.create({
  name: 'fontSize',
  addOptions() { return { types: ['textStyle'] } },
  addGlobalAttributes() {
    return [
      {
        types: this.options.types,
        attributes: {
          fontSize: {
            default: null,
            parseHTML: element => element.style.fontSize?.replace(/['"]+/g, ''),
            renderHTML: attributes => {
              if (!attributes.fontSize) return {}
              return { style: `font-size: ${attributes.fontSize}` }
            },
          },
        },
      },
    ]
  },
  addCommands() {
    return {
      setFontSize: fontSize => ({ chain }) => chain().setMark('textStyle', { fontSize }).run(),
      unsetFontSize: () => ({ chain }) => chain().setMark('textStyle', { fontSize: null }).removeEmptyTextStyle().run(),
    }
  },
})

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const editor = new Editor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Image.configure({
      inline: true,
      allowBase64: true,
    }),
    Link.configure({
      openOnClick: false,
    }),
    TextStyle,
    FontSize,
    Color
  ],
  editorProps: {
    attributes: {
      class: 'focus:outline-none h-full',
    },
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  },
})

watch(() => props.modelValue, (value) => {
  const isSame = editor.getHTML() === value
  if (isSame) {
    return
  }
  editor.commands.setContent(value, false)
})

onBeforeUnmount(() => {
  editor.destroy()
})
</script>

<style>
/* 针对 tiptap 默认焦点的蓝框进行清理 */
.ProseMirror:focus {
  outline: none;
}
/* 给图片添加选中样式，方便用户知道选中了可以删除 */
.ProseMirror img.ProseMirror-selectednode {
  outline: 3px solid #3b82f6;
  border-radius: 8px;
}
</style>
