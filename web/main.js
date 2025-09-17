const editor = document.getElementById('editor')
const status = document.getElementById('status')
const cursorOverlay = document.getElementById('cursor-overlay')
const preview = document.getElementById('preview')
let ws
let isConnected = false
let myUserId = 'user_' + Math.random().toString(36).substring(2, 9)
let cursors = {}
let debounceTimer = null
let lastContent = ''

console.log('My User ID:', myUserId)

function calculateCursorPosition(text, position) {
  const lines = text.substring(0, position).split('\n')
  const line = lines.length - 1
  const column = lines[lines.length - 1].length
  return { line, column }
}

function updatePreview(html) {
  if (html) {
    preview.innerHTML = html
  } else {
    preview.innerHTML =
      '<p style="color: #6a737d; font-style: italic;">プレビューがここに表示されます...</p>'
  }
}

function updateCursorDisplay() {
  cursorOverlay.innerHTML = ''
  console.log('Updating cursor display, cursors:', cursors)

  for (const userId in cursors) {
    if (userId === myUserId) continue

    const cursor = cursors[userId]
    const pos = calculateCursorPosition(editor.value, cursor.position)
    console.log('Cursor position for', userId, ':', pos)

    const cursorElement = document.createElement('div')
    cursorElement.className = 'cursor'
    cursorElement.style.backgroundColor = '#FF6B6B'
    cursorElement.style.left = pos.column * 8.4 + 'px'
    cursorElement.style.top = pos.line * 21 + 'px'

    const label = document.createElement('div')
    label.className = 'cursor-label'
    label.style.backgroundColor = '#FF6B6B'
    label.textContent = cursor.userId

    cursorElement.appendChild(label)
    cursorOverlay.appendChild(cursorElement)
  }
}

function sendCursorPosition() {
  if (isConnected) {
    const message = JSON.stringify({
      type: 'cursor',
      userId: myUserId,
      position: editor.selectionStart,
      content: '#FF6B6B',
    })
    console.log('Sending cursor position:', message)
    ws.send(message)
  }
}

function sendUpdateMessage() {
  if (isConnected && editor.value !== lastContent) {
    const message = JSON.stringify({
      type: 'update',
      content: editor.value,
    })
    ws.send(message)
    lastContent = editor.value
    sendCursorPosition()
  }
}

function debouncedSendUpdate() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(sendUpdateMessage, 300) // 300msのデバウンス
}

function connect() {
  ws = new WebSocket('ws://localhost:8888/ws')

  ws.onopen = function () {
    isConnected = true
    status.textContent = '接続済み'
    status.className = 'status connected'
  }

  ws.onmessage = function (event) {
    const message = JSON.parse(event.data)

    if (message.type === 'init') {
      if (editor.value !== message.content) {
        const cursorPos = editor.selectionStart
        editor.value = message.content
        editor.setSelectionRange(cursorPos, cursorPos)
      }
      // レンダリング結果をプレビューに表示
      updatePreview(message.renderedHtml)
    } else if (message.type === 'update') {
      if (editor.value !== message.content) {
        const cursorPos = editor.selectionStart
        editor.value = message.content
        editor.setSelectionRange(cursorPos, cursorPos)
      }
      // レンダリング結果をプレビューに表示
      updatePreview(message.renderedHtml)
    } else if (message.type === 'cursor') {
      console.log('Received cursor message:', message)
      if (message.userId !== myUserId) {
        cursors[message.userId] = {
          userId: message.userId,
          position: message.position,
        }
        console.log('Updated cursors:', cursors)
        updateCursorDisplay()
      }
    } else if (message.type === 'cursor_disconnect') {
      delete cursors[message.userId]
      updateCursorDisplay()
    }
  }

  ws.onclose = function () {
    isConnected = false
    status.textContent = '切断されました - 再接続中...'
    status.className = 'status disconnected'
    setTimeout(connect, 1000)
  }

  ws.onerror = function () {
    isConnected = false
    status.textContent = 'エラーが発生しました - 再接続中...'
    status.className = 'status disconnected'
  }
}

editor.addEventListener('input', debouncedSendUpdate)

editor.addEventListener('selectionchange', sendCursorPosition)
editor.addEventListener('click', sendCursorPosition)
editor.addEventListener('keyup', sendCursorPosition)

// 初期状態でプレビューを設定
updatePreview('')

connect()
