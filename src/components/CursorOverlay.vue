<template>
  <div class="cursor-overlay" ref="overlayRef">
    <div
      v-for="cursor in cursors"
      :key="cursor.userId"
      :style="getCursorStyle(cursor)"
      class="cursor"
    >
      <div class="cursor-label" :style="{ backgroundColor: cursor.color }">
        {{ cursor.userId }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import type { CursorInfo } from '../composables/useCursors'

type Props = {
  cursors: CursorInfo[]
  editorContent: string
  textareaRef?: HTMLTextAreaElement
}

const props = defineProps<Props>()
const overlayRef = ref<HTMLElement>()

// テキスト内の位置から行と列を計算
const calculateCursorPosition = (text: string, position: number) => {
  const lines = text.substring(0, position).split('\n')
  const line = lines.length - 1
  const column = lines[lines.length - 1].length
  return { line, column }
}

// カーソルのスタイルを計算
const getCursorStyle = (
  cursor: CursorInfo
): Record<string, string | number> => {
  if (!props.editorContent) return { display: 'none' }

  const pos = calculateCursorPosition(props.editorContent, cursor.position)

  // フォントサイズとライン高さはCSSと同じ値を使用
  const fontSize = 14
  const lineHeight = 1.5
  const charWidth = fontSize * 0.6 // モノスペースフォントの概算文字幅
  const lineHeightPx = fontSize * lineHeight

  return {
    position: 'absolute',
    left: `${pos.column * charWidth}px`,
    top: `${pos.line * lineHeightPx}px`,
    width: '2px',
    height: `${lineHeightPx}px`,
    backgroundColor: cursor.color,
    zIndex: 10,
    pointerEvents: 'none' as const,
    animation: 'blink 1s infinite',
  }
}

// カーソル位置の更新を監視
watch(
  [() => props.cursors, () => props.editorContent],
  () => {
    nextTick(() => {
      // DOM更新後に再計算
    })
  },
  { deep: true }
)
</script>

<style scoped>
.cursor-overlay {
  position: absolute;
  top: 10px;
  left: 10px;
  pointer-events: none;
  z-index: 2;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  white-space: pre;
  line-height: 1.5;
  font-size: 14px;
  width: calc(100% - 20px);
  height: calc(100% - 20px);
}

.cursor {
  position: absolute;
  width: 2px;
  height: 20px;
}

.cursor-label {
  position: absolute;
  top: -25px;
  left: -5px;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  color: white;
  white-space: nowrap;
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  z-index: 11;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}
</style>
