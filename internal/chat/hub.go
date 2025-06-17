package chat

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tuusuario/blog-realtime/internal/model"
)

type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
}

type client struct {
	conn *websocket.Conn
	send chan model.Comment
}

type subscription struct {
	postID int
	client *client
}

type broadcastMsg struct {
	postID  int
	comment model.Comment
}

type Hub struct {
	register   chan subscription
	unregister chan subscription
	broadcast  chan broadcastMsg
	rooms      map[int]map[*client]bool
	mongoCol   MongoCollection
}

func NewHub(col MongoCollection) *Hub {
	return &Hub{
		register:   make(chan subscription),
		unregister: make(chan subscription),
		broadcast:  make(chan broadcastMsg),
		rooms:      make(map[int]map[*client]bool),
		mongoCol:   col,
	}
}

func (h *Hub) Run() {
	for {
		select {

		// Nuevo cliente en sala postID
		case sub := <-h.register:
			clients, ok := h.rooms[sub.postID]
			if !ok {
				clients = make(map[*client]bool)
				h.rooms[sub.postID] = clients
			}
			h.rooms[sub.postID][sub.client] = true

		// Cliente sale de sala
		case sub := <-h.unregister:
			if clients, ok := h.rooms[sub.postID]; ok {
				if _, exists := clients[sub.client]; exists {
					delete(clients, sub.client)
					close(sub.client.send)
				}
				if len(clients) == 0 {
					delete(h.rooms, sub.postID)
				}
			}

		case bm := <-h.broadcast:

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := h.mongoCol.InsertOne(ctx, bm.comment)
			cancel()
			if err != nil {
				log.Println("hub: error guardando comentario:", err)
			}

			if clients, ok := h.rooms[bm.postID]; ok {
				for c := range clients {
					select {
					case c.send <- bm.comment:
					default:

						close(c.send)
						delete(clients, c)
					}
				}
			}
		}
	}
}

func (h *Hub) Register(postID int, c *client) {
	h.register <- subscription{postID, c}
}

func (h *Hub) Unregister(postID int, c *client) {
	h.unregister <- subscription{postID, c}
}

func (h *Hub) Broadcast(postID int, comment model.Comment) {
	h.broadcast <- broadcastMsg{postID, comment}
}

func (c *client) Writer() {
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
