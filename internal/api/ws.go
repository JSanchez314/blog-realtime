package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JSanchez314/blog-realtime/internal/chat"
	"github.com/JSanchez314/blog-realtime/internal/model"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func CommentSocket(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
		conn, _ := upgrader.Upgrade(w, r, nil)
		client := &chat.Client{conn: conn, send: make(chan model.Comment)}
		hub.Register(postID, client)

		defer func() {
			hub.Unregister(postID, client)
			conn.Close()
		}()

		go client.Writer()
		for {
			var msg model.Comment
			if err := conn.ReadJSON(&msg); err != nil {
				break
			}
			msg.Timestamp = time.Now()
			hub.Broadcast(postID, msg)
		}
	}
}
