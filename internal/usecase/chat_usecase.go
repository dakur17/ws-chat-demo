package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nickyrolly/ws-chat-demo/internal/repository"
)

type userConnMap map[string][]*websocket.Conn

type ChatBox struct {
	clients userConnMap
	mu      sync.Mutex
}

func NewChatBox() *ChatBox {
	return &ChatBox{
		clients: make(userConnMap),
	}
}

func (cb *ChatBox) AddClient(chatboxID string, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[chatboxID]; !ok {
		cb.clients[chatboxID] = []*websocket.Conn{}
	}
	cb.clients[chatboxID] = append(cb.clients[chatboxID], conn)
	log.Printf("Add client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *ChatBox) RemoveClient(chatboxID string, conn *websocket.Conn) {
	cb.mu.Lock()
	if _, ok := cb.clients[chatboxID]; ok {
		// Find conn index
		idx := cb.findConn(chatboxID, conn)
		// Remove conn from slice if conn found
		if idx != -1 {
			cb.clients[chatboxID] = append(cb.clients[chatboxID][:idx], cb.clients[chatboxID][idx+1:]...)
		}
	}
	log.Printf("Remove client : %+v\n", cb.clients)
	cb.mu.Unlock()
}

func (cb *ChatBox) Broadcast(chatboxID string, curConn *websocket.Conn, message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	log.Println("chatboxID :", chatboxID)

	// Send to chatbox user pools except for current connection
	for _, conn := range cb.clients[chatboxID] {
		if conn == curConn {
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error broadcasting message to user :", err)
		}
	}

	log.Printf("Broadcast clients : %+v\n", cb.clients)
}

func (cb *ChatBox) findConn(chatboxID string, conn *websocket.Conn) int {
	for i, c := range cb.clients[chatboxID] {
		if c == conn {
			return i
		}
	}
	return -1
}

func (cb *ChatBox) PublishSaveChatHistory(params repository.ChatHistoryData) error {
	// messageBody, err := json.Marshal(params)
	// if err != nil {
	// 	log.Println("Error Marshal:", err)
	// 	return err
	// }

	// Exercise 3.1.2
	// Please complete this block to publish personal chat message via NSQ

	return nil
}

func (cb *ChatBox) GetChatHistory(ctx context.Context, params repository.ChatHistoryData) ([]map[string]interface{}, error) {
	var (
		chatHistoryData = []map[string]interface{}{}
	)
	// Sort user id for A and B to be from lowest to highest
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()
	log.Printf("%+v", params)

	// Exercise 3.3.3
	// Please complete this block to get chat history data from database

	// Uncomment code below after get chat history from DB has implemented (Exercise 3.3.3)
	// for rows.Next() {
	// 	var (
	// 		data         map[string]interface{}
	// 		senderUserID int
	// 		message      string
	// 		replyTime    time.Time
	// 	)

	// 	err := rows.Scan(&senderUserID, &message, &replyTime)
	// 	if err != nil {
	// 		log.Println("[GetChatHistory] Error Scan: ", err)
	// 		return chatHistoryData, err
	// 	}

	// 	data = map[string]interface{}{
	// 		"sender_user_id": senderUserID,
	// 		"message":        message,
	// 		"reply_time":     replyTime,
	// 	}

	// 	chatHistoryData = append(chatHistoryData, data)
	// }

	return chatHistoryData, nil
}
