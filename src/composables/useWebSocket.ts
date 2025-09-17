import { ref, onUnmounted } from 'vue'

export type WebSocketMessage = {
  type: 'init' | 'update' | 'cursor' | 'cursor_disconnect'
  content?: string
  renderedHtml?: string
  userId?: string
  position?: number
}

export type ConnectionStatus = {
  isConnected: boolean
  status: 'connecting' | 'connected' | 'disconnected' | 'error'
  statusText: string
}

export function useWebSocket(url: string) {
  const ws = ref<WebSocket | null>(null)
  const connectionStatus = ref<ConnectionStatus>({
    isConnected: false,
    status: 'connecting',
    statusText: '接続中...',
  })

  const messageHandlers = ref<Map<string, (message: WebSocketMessage) => void>>(
    new Map()
  )
  let reconnectTimer: number | null = null

  const connect = () => {
    try {
      ws.value = new WebSocket(url)

      ws.value.onopen = () => {
        connectionStatus.value = {
          isConnected: true,
          status: 'connected',
          statusText: '接続済み',
        }
      }

      ws.value.onmessage = event => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          const handler = messageHandlers.value.get(message.type)
          if (handler) {
            handler(message)
          }
        } catch (error) {
          console.error('WebSocket message parsing error:', error)
        }
      }

      ws.value.onclose = () => {
        connectionStatus.value = {
          isConnected: false,
          status: 'disconnected',
          statusText: '切断されました - 再接続中...',
        }

        // 1秒後に再接続を試行
        if (reconnectTimer) clearTimeout(reconnectTimer)
        reconnectTimer = setTimeout(connect, 1000)
      }

      ws.value.onerror = () => {
        connectionStatus.value = {
          isConnected: false,
          status: 'error',
          statusText: 'エラーが発生しました - 再接続中...',
        }
      }
    } catch (error) {
      console.error('WebSocket connection error:', error)
      connectionStatus.value = {
        isConnected: false,
        status: 'error',
        statusText: '接続エラー',
      }
    }
  }

  const sendMessage = (message: WebSocketMessage) => {
    if (ws.value && connectionStatus.value.isConnected) {
      try {
        ws.value.send(JSON.stringify(message))
      } catch (error) {
        console.error('WebSocket send error:', error)
      }
    }
  }

  const registerHandler = (
    messageType: string,
    handler: (message: WebSocketMessage) => void
  ) => {
    messageHandlers.value.set(messageType, handler)
  }

  const unregisterHandler = (messageType: string) => {
    messageHandlers.value.delete(messageType)
  }

  const disconnect = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }

    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connectionStatus,
    connect,
    sendMessage,
    registerHandler,
    unregisterHandler,
    disconnect,
  }
}
