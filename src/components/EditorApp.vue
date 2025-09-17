<template>
  <div class="editor-app">
    <div class="header">
      <h1 class="app-title">リアルタイムMarkdownエディター</h1>
      <StatusIndicator :status="connectionStatus" />
    </div>
    <div class="main-container">
      <MarkdownEditor
        v-model="content"
        :other-cursors="otherCursors"
        @cursor-move="handleCursorMove"
      />
      <PreviewPane :rendered-html="renderedHtml" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MarkdownEditor from './MarkdownEditor.vue'
import PreviewPane from './PreviewPane.vue'
import StatusIndicator from './StatusIndicator.vue'
import { useWebSocket } from '../composables/useWebSocket'
import { useCursors } from '../composables/useCursors'

const content = ref('')
const renderedHtml = ref('')
const myUserId = `user_${Math.random().toString(36).substring(2, 9)}`

// WebSocket接続
const { connectionStatus, connect, sendMessage, registerHandler } =
  useWebSocket('ws://localhost:8000/ws')

// カーソル管理
const { otherCursors, updateCursor, removeCursor } = useCursors(myUserId)

let debounceTimer: number | null = null
let lastContent = ''

// デバウンス付きでコンテンツ更新を送信
const sendContentUpdate = () => {
  if (connectionStatus.value.isConnected && content.value !== lastContent) {
    sendMessage({
      type: 'update',
      content: content.value,
    })
    lastContent = content.value
    sendCursorUpdate()
  }
}

const debouncedSendUpdate = () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(sendContentUpdate, 300)
}

// カーソル位置を送信
const sendCursorUpdate = () => {
  if (connectionStatus.value.isConnected) {
    const position = 0 // エディターコンポーネントから取得する必要がある
    sendMessage({
      type: 'cursor',
      userId: myUserId,
      position,
      content: '#FF6B6B', // 後でuseCursorsから取得する色に変更
    })
  }
}

// エディターからのカーソル移動イベント
const handleCursorMove = (position: number) => {
  if (connectionStatus.value.isConnected) {
    sendMessage({
      type: 'cursor',
      userId: myUserId,
      position,
      content: '#FF6B6B',
    })
  }
}

// コンテンツ変更の監視
const handleContentChange = () => {
  debouncedSendUpdate()
}

// WebSocketメッセージハンドラー
const setupMessageHandlers = () => {
  registerHandler('init', message => {
    if (content.value !== message.content) {
      content.value = message.content || ''
    }
    if (message.renderedHtml) {
      renderedHtml.value = message.renderedHtml
    }
  })

  registerHandler('update', message => {
    if (content.value !== message.content) {
      content.value = message.content || ''
    }
    if (message.renderedHtml) {
      renderedHtml.value = message.renderedHtml
    }
  })

  registerHandler('cursor', message => {
    if (
      message.userId &&
      message.userId !== myUserId &&
      message.position !== undefined
    ) {
      updateCursor(message.userId, message.position)
    }
  })

  registerHandler('cursor_disconnect', message => {
    if (message.userId) {
      removeCursor(message.userId)
    }
  })
}

// コンテンツ変更の監視
const watchContent = () => {
  let previousContent = content.value
  setInterval(() => {
    if (content.value !== previousContent) {
      previousContent = content.value
      handleContentChange()
    }
  }, 100)
}

onMounted(() => {
  console.log('My User ID:', myUserId)
  setupMessageHandlers()
  connect()
  watchContent()
})
</script>

<style scoped>
.editor-app {
  font-family:
    -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  margin: 0;
  padding: 0;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.header {
  padding: 10px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #ddd;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-title {
  margin: 0;
  font-size: 18px;
}

.main-container {
  display: flex;
  flex: 1;
  height: calc(100vh - 60px);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .main-container {
    flex-direction: column;
  }
}
</style>
