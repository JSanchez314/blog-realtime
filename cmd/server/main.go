// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/tuusuario/blog-realtime/internal/api"
	"github.com/tuusuario/blog-realtime/internal/chat"
	"github.com/tuusuario/blog-realtime/internal/config"
	"github.com/tuusuario/blog-realtime/internal/db"
)

func main() {
	cfg := config.Load()

	pgDB, err := db.NewPostgres(cfg.PgURL)
	if err != nil {
		log.Fatalf("error conectando a Postgres: %v", err)
	}
	defer pgDB.Close()

	mongoCol, err := db.NewMongo(cfg.MongoURI, cfg.MongoCAPath)
	if err != nil {
		log.Fatalf("error conectando a MongoDB: %v", err)
	}

	hub := chat.NewHub(mongoCol)
	go hub.Run()

	r := mux.NewRouter()

	postsH := api.NewPostsHandler(pgDB)
	r.HandleFunc("/api/posts", postsH.List).Methods(http.MethodGet)
	r.HandleFunc("/api/posts/{id}", postsH.Get).Methods(http.MethodGet)
	r.Handle("/api/posts",
		api.JWTMiddleware(cfg.JwtSecret)(
			http.HandlerFunc(postsH.Create),
		),
	).Methods(http.MethodPost)

	r.HandleFunc("/ws/comments", api.CommentSocket(hub))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor escuchando en %s …", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
