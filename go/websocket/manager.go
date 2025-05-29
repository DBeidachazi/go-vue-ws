// websocket/manager.go
package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const (
	TextMessage          = websocket.TextMessage
	CloseGoingAway       = websocket.CloseGoingAway
	CloseAbnormalClosure = websocket.CloseAbnormalClosure
)

type Upgrader = websocket.Upgrader

func IsUnexpectedCloseError(err error, codes ...int) bool {
	return websocket.IsUnexpectedCloseError(err, codes...)
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// NewManager 创建WebSocket管理器
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) RegisterClient(client *Client) {
	m.register <- client
}

func (m *Manager) UnregisterClient(client *Client) {
	m.unregister <- client
}

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
			log.Println("新客户端连接")
		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				m.mutex.Lock()
				delete(m.clients, client)
				close(client.Send)
				m.mutex.Unlock()
				log.Println("客户端断开连接")
			}
		case message := <-m.broadcast:
			m.mutex.Lock()
			for client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, client)
				}
			}
			m.mutex.Unlock()
		}
	}
}

// BroadcastPollUpdate 广播投票统计更新
func (m *Manager) BroadcastPollUpdate(message []byte) {
	m.broadcast <- message
}
