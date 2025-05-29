package router

import (
	"github.com/gin-gonic/gin"
	"goproject/websocket"
	"log"
	"net/http"
)

var (
	wsManager *websocket.Manager
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func init() {
	wsManager = websocket.NewManager()
	go wsManager.Start()
}

func InitWebSocketRouter(routerGroup *gin.RouterGroup) {
	wsRouter := routerGroup.Group("/ws")
	{
		wsRouter.GET("/poll", handleWebSocket)
	}
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级WebSocket连接失败: %v", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	wsManager.RegisterClient(client)

	go writeMessages(client)
	go readMessages(client)
}

// 向客户端发送消息
func writeMessages(client *websocket.Client) {
	defer func() {
		client.Conn.Close()
	}()

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		}
	}
}

func readMessages(client *websocket.Client) {
	defer func() {
		wsManager.UnregisterClient(client)
		client.Conn.Close()
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break
		}
		// do something
	}
}
