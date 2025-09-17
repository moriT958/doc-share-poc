package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	document   string
	cursors    map[string]*CursorInfo
	mutex      sync.RWMutex
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	color  string
}

type Message struct {
	Type         string `json:"type"`
	Content      string `json:"content"`
	RenderedHTML string `json:"renderedHtml,omitempty"`
	UserID       string `json:"userId,omitempty"`
	Position     int    `json:"position,omitempty"`
}

type CursorInfo struct {
	UserID   string `json:"userId"`
	Position int    `json:"position"`
	Color    string `json:"color"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	mdParser      goldmark.Markdown
	htmlSanitizer *bluemonday.Policy
)

func init() {
	// Goldmarkパーサーの初期化
	mdParser = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	// HTMLサニタイザーの初期化
	htmlSanitizer = bluemonday.UGCPolicy()
}

// MarkdownをHTMLにレンダリングしてサニタイズする
func renderMarkdown(source string) string {
	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(source), &buf); err != nil {
		log.Printf("Markdown rendering error: %v", err)
		return ""
	}
	return htmlSanitizer.Sanitize(buf.String())
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		document:   "",
		cursors:    make(map[string]*CursorInfo),
	}
}

func generateUserID() string {
	rand.Seed(time.Now().UnixNano())
	return "user_" + string(rune(rand.Intn(26)+65)) + string(rune(rand.Intn(26)+65)) + string(rune(rand.Intn(26)+65))
}

func generateColor() string {
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", "#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E9"}
	rand.Seed(time.Now().UnixNano())
	return colors[rand.Intn(len(colors))]
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.mutex.RLock()
			content := h.document
			cursors := make(map[string]*CursorInfo)
			for k, v := range h.cursors {
				cursors[k] = v
			}
			h.mutex.RUnlock()

			initMsg := Message{
				Type:         "init",
				Content:      content,
				RenderedHTML: renderMarkdown(content),
			}
			initData, _ := json.Marshal(initMsg)

			select {
			case client.send <- initData:
				for _, cursor := range cursors {
					if cursor.UserID != client.userID {
						cursorData, _ := json.Marshal(Message{
							Type:     "cursor",
							UserID:   cursor.UserID,
							Position: cursor.Position,
							Content:  cursor.Color,
						})
						client.send <- cursorData
					}
				}
			default:
				close(client.send)
				delete(h.clients, client)
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.mutex.Lock()
				delete(h.cursors, client.userID)
				h.mutex.Unlock()

				disconnectMsg := Message{
					Type:   "cursor_disconnect",
					UserID: client.userID,
				}
				disconnectData, _ := json.Marshal(disconnectMsg)
				for otherClient := range h.clients {
					select {
					case otherClient.send <- disconnectData:
					default:
						close(otherClient.send)
						delete(h.clients, otherClient)
					}
				}
			}

		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			h.mutex.Lock()
			switch msg.Type {
			case "update":
				h.document = msg.Content
				// Markdownレンダリングを実行してメッセージに追加
				msg.RenderedHTML = renderMarkdown(msg.Content)
				message, _ = json.Marshal(msg) // レンダリング結果を含めたメッセージで更新
			case "cursor":
				h.cursors[msg.UserID] = &CursorInfo{
					UserID:   msg.UserID,
					Position: msg.Position,
					Color:    msg.Content,
				}
			}
			h.mutex.Unlock()

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		msg, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: generateUserID(),
		color:  generateColor(),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	log.Println("サーバーを http://localhost:8888 で起動しています...")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
