package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type hub struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	document   string
	cursors    map[string]*CursorInfo
	mutex      sync.RWMutex
}

type client struct {
	hub    *hub
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub() *hub {
	return &hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		document:   "",
		cursors:    make(map[string]*CursorInfo),
	}
}

func (h *hub) Start() {
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

func (c *client) readPump() {
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

func (c *client) writePump() {
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

func ServeWS(hub *hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{
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
