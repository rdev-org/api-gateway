package notification

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rdev-org/api-gateway/internal/config"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*websocket.Conn
}

type NotificationHandler struct {
	kafkaReader *kafka.Reader
	connManager *ConnectionManager
}

func NewNotificationHandler(c *config.KafkaConfig) *NotificationHandler {
	h := &NotificationHandler{
		kafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: c.Brokers,
			Topic:   c.Topic,
			GroupID: c.GroupID,
		}),
		connManager: &ConnectionManager{
			connections: make(map[string]*websocket.Conn),
		},
	}

	go h.consumeKafkaEvents(context.Background())
	return h
}

func (h *NotificationHandler) InAppNotification(w http.ResponseWriter, r *http.Request) {
	// userID := r.Header.Get("X-User-Id")
	// if userID == "" {
	// 	http.Error(w, "Unauthorized: missing user_id", http.StatusBadRequest)
	// 	return
	// }

	userID := "1"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Failed to upgrade websocket: %v", err)
		return
	}

	h.connManager.mu.Lock()
	h.connManager.connections[userID] = conn
	h.connManager.mu.Unlock()

	defer func() {
		h.connManager.mu.Lock()
		delete(h.connManager.connections, userID)
		h.connManager.mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *NotificationHandler) consumeKafkaEvents(ctx context.Context) {
	for {
		m, err := h.kafkaReader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[API-Gateway WS Kafka] Error reading message: %v", err)
			continue
		}

		userID := "1"
		fmt.Println()
		fmt.Println("KAFKA READ: ", m.Topic, string(m.Value), string(m.Key))
		fmt.Println()

		h.connManager.mu.RLock()
		conn, exists := h.connManager.connections[userID]
		h.connManager.mu.RUnlock()

		if exists && conn != nil {
			if err := conn.WriteMessage(websocket.TextMessage, m.Value); err != nil {
				log.Printf("[API-Gateway WS] Push error for user %s: %v", userID, err)
			}
		}
	}
}

func (h *NotificationHandler) Close() error {
	if h.kafkaReader != nil {
		return h.kafkaReader.Close()
	}
	return nil
}
