# リアルタイム Markdown 共同編集アプリケーション

複数のユーザーが同時にMarkdownドキュメントを編集できるリアルタイム共同編集アプリです。
WebSocketを使用してリアルタイムでの文書同期と、他のユーザーのカーソル位置表示を実現しています。

## 🚀 機能

- **リアルタイム共同編集**: 複数のユーザーが同時にMarkdownを編集
- **ライブプレビュー**: 編集内容をリアルタイムでHTML表示
- **カーソル表示**: 他のユーザーのカーソル位置をリアルタイム表示
- **接続状態表示**: WebSocket接続の状態を視覚的に確認
- **Markdownサポート**: 標準的なMarkdown記法に対応

## 🛠 技術スタック

### フロントエンド

- **Vue 3** - プログレッシブJavaScriptフレームワーク
- **TypeScript** - 型安全性を提供
- **Vite** - 高速ビルドツール
- **Composition API** - Vue 3の最新API

### バックエンド

- **Go** - 高性能なバックエンドサーバー
- **Gorilla WebSocket** - WebSocket通信ライブラリ
- **goldmark** - Markdown → HTML変換
- **bluemonday** - HTMLサニタイズ（XSS対策）

## 📁 プロジェクト構造

```text
doc-share-poc-2/
├── 📁 src/                      # Vue.js フロントエンド
│   ├── 📁 components/
│   │   ├── EditorApp.vue        # メインアプリケーション
│   │   ├── MarkdownEditor.vue   # Markdownエディター
│   │   ├── PreviewPane.vue      # HTMLプレビューペイン
│   │   ├── CursorOverlay.vue    # 他ユーザーのカーソル表示
│   │   └── StatusIndicator.vue  # 接続状態インジケーター
│   ├── 📁 composables/
│   │   ├── useWebSocket.ts      # WebSocket接続管理
│   │   └── useCursors.ts        # カーソル位置管理
│   ├── App.vue                  # ルートコンポーネント
│   └── main.ts                  # エントリーポイント
├── 📁 internal/                 # Go バックエンド
│   ├── websocket.go             # WebSocket接続ハンドリング
│   └── render.go                # Markdownレンダリング
├── main.go                      # Goサーバーエントリーポイント
├── package.json                 # Node.js依存関係
├── go.mod                       # Go モジュール定義
└── Makefile                     # ビルドスクリプト
```

## 🏗 セットアップ

### 必要な環境

- **Node.js** 18以上
- **Go** 1.22以上
- **npm** または **yarn**

### インストール手順

1. **リポジトリをクローン**

```bash
git clone <repository-url>
cd doc-share-poc-2
```

2. **フロントエンド依存関係をインストール**

```bash
npm install
```

3. **Go依存関係をインストール**

```bash
go mod tidy
```

## 🚀 起動方法

### 開発モード（推奨）

**ターミナル1: Goサーバーを起動**

```bash
make dev
# または
go run main.go
```

サーバーは `http://localhost:8000` で起動します

**ターミナル2: Vue開発サーバーを起動**

```bash
npm run dev
```

フロントエンドは `http://localhost:5173` で起動します

### プロダクションビルド

```bash
# フロントエンドをビルド
npm run build

# Goサーバーを起動（静的ファイルも配信）
go run main.go
```

## 📋 開発コマンド

| コマンド          | 説明                   |
| ----------------- | ---------------------- |
| `npm run dev`     | Vue開発サーバー起動    |
| `npm run build`   | プロダクションビルド   |
| `npm run preview` | ビルド結果をプレビュー |
| `npm run fmt`     | コードフォーマット     |
| `make dev`        | Goサーバー起動         |
| `go mod tidy`     | Go依存関係を整理       |

## 🏗 アーキテクチャ

```
┌─────────────────┐    WebSocket     ┌─────────────────┐
│   Vue.js        │ ◄─────────────► │   Go Server     │
│   Frontend      │    Port 8000     │   Backend       │
│   Port 5173     │                  │                 │
└─────────────────┘                  └─────────────────┘
        │                                      │
        ▼                                      ▼
┌─────────────────┐                  ┌─────────────────┐
│  - エディター    │                  │  - WebSocket    │
│  - プレビュー    │                  │    Hub管理      │
│  - カーソル表示  │                  │  - Markdown     │
│  - 状態管理     │                  │    レンダリング  │
└─────────────────┘                  └─────────────────┘
```

### 通信フロー

1. ユーザーがエディターでテキストを編集
2. WebSocket経由でサーバーに変更を送信
3. サーバーが全接続ユーザーに変更を配信
4. 各クライアントでエディターとプレビューを更新

## 🤝 開発に参加する

1. フィーチャーブランチを作成
2. 変更を実装
3. `npm run fmt` でコードをフォーマット
4. プルリクエストを作成

## 📝 今後の改善予定

- [ ] ユーザー登録機能
- [ ] 文書の保存/読み込み機能
- [ ] より豊富なMarkdown記法サポート(Pukiwikiの記法に対応)
