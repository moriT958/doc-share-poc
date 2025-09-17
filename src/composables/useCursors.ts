import { ref, computed } from 'vue'

export type CursorInfo = {
  userId: string
  position: number
  color: string
}

export type CursorPosition = {
  line: number
  column: number
}

export function useCursors(myUserId: string) {
  const cursors = ref<Map<string, CursorInfo>>(new Map())

  // ユーザーごとに固有の色を生成
  const generateUserColor = (userId: string): string => {
    // ユーザーIDからハッシュ値を生成
    let hash = 0
    for (let i = 0; i < userId.length; i++) {
      const char = userId.charCodeAt(i)
      hash = (hash << 5) - hash + char
      hash = hash & hash // 32bit整数に変換
    }

    // HSLカラーで生成（彩度と明度を固定し、色相をバリエーション）
    const hue = Math.abs(hash) % 360
    const saturation = 70 + (Math.abs(hash) % 20) // 70-90%
    const lightness = 45 + (Math.abs(hash) % 20) // 45-65%

    return `hsl(${hue}, ${saturation}%, ${lightness}%)`
  }

  // テキスト内の位置から行と列を計算
  const calculateCursorPosition = (
    text: string,
    position: number
  ): CursorPosition => {
    const lines = text.substring(0, position).split('\n')
    const line = lines.length - 1
    const column = lines[lines.length - 1].length
    return { line, column }
  }

  // 自分以外のカーソル情報を取得
  const otherCursors = computed(() => {
    const result: CursorInfo[] = []
    for (const [userId, cursor] of cursors.value) {
      if (userId !== myUserId) {
        result.push(cursor)
      }
    }
    return result
  })

  // カーソル情報を更新
  const updateCursor = (userId: string, position: number) => {
    const color = generateUserColor(userId)
    cursors.value.set(userId, {
      userId,
      position,
      color,
    })
  }

  // カーソルを削除
  const removeCursor = (userId: string) => {
    cursors.value.delete(userId)
  }

  // すべてのカーソルをクリア
  const clearCursors = () => {
    cursors.value.clear()
  }

  // デバッグ用：全カーソル情報を取得
  const getAllCursors = computed(() => {
    return Array.from(cursors.value.entries()).map(([id, cursor]) => ({
      id,
      ...cursor,
    }))
  })

  return {
    cursors,
    otherCursors,
    updateCursor,
    removeCursor,
    clearCursors,
    calculateCursorPosition,
    generateUserColor,
    getAllCursors,
  }
}
