<template>
  <div class="editor-container">
    <div class="pane-header">Markdownエディター</div>
    <div class="editor-wrapper">
      <textarea
        ref="textareaRef"
        v-model="content"
        class="editor-textarea"
        placeholder="Markdownを入力してください..."
        @input="handleInput"
        @selectionchange="handleSelectionChange"
        @click="handleSelectionChange"
        @keyup="handleSelectionChange"
      />
      <CursorOverlay
        :cursors="otherCursors"
        :editor-content="content"
        :textarea-ref="textareaRef"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import CursorOverlay from './CursorOverlay.vue'
import type { CursorInfo } from '../composables/useCursors'

type Props = {
  modelValue: string
  otherCursors: CursorInfo[]
}

type Emits = {
  (e: 'update:modelValue', value: string): void
  (e: 'cursorMove', position: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const textareaRef = ref<HTMLTextAreaElement>()
const content = ref('')
let debounceTimer: number | null = null

// プロップスの変更を監視
watch(
  () => props.modelValue,
  newValue => {
    if (content.value !== newValue) {
      const cursorPos = textareaRef.value?.selectionStart || 0
      content.value = newValue

      // DOM更新後にカーソル位置を復元
      setTimeout(() => {
        if (textareaRef.value) {
          textareaRef.value.setSelectionRange(cursorPos, cursorPos)
        }
      }, 0)
    }
  }
)

// 初期値を設定
onMounted(() => {
  content.value = props.modelValue
})

const handleInput = () => {
  emit('update:modelValue', content.value)

  // デバウンス処理でカーソル位置を送信
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    handleSelectionChange()
  }, 100)
}

const handleSelectionChange = () => {
  if (textareaRef.value) {
    const position = textareaRef.value.selectionStart
    emit('cursorMove', position)
  }
}

// 外部からテキストエリアにフォーカスするための関数
const focus = () => {
  textareaRef.value?.focus()
}

defineExpose({
  focus,
  textareaRef,
})
</script>

<style scoped>
.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  border-right: 2px solid #ddd;
}

.pane-header {
  padding: 8px 12px;
  background-color: #e9ecef;
  border-bottom: 1px solid #ddd;
  font-weight: bold;
  font-size: 14px;
}

.editor-wrapper {
  flex: 1;
  position: relative;
}

.editor-textarea {
  width: 100%;
  height: 100%;
  padding: 10px;
  border: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  outline: none;
  background-color: #fff;
  box-sizing: border-box;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .editor-container {
    border-right: none;
    border-bottom: 2px solid #ddd;
  }
}
</style>
