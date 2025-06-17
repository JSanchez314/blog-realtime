package api

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/JSanchez314/blog-realtime/internal/model"
)

type PostHandler struct {
	db *sqlx.DB
}

func NewPostHandler(db *sqlx.DB) *PostHandler { return &PostHandler{db}}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	var posts []model.Post
	h.db.Select(&posts, "SELECT * FROM posts ORDER BY created_at DESC")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p model.Post
	json.NewDecoder(r.Body).Decode(&p)
	err := h.db.QueryRowx
	("INSERT INTO posts(author_id, category_id, title, body) VALUES($1,$2,$3,$4) RETURNING id,created_at", 
	p.AuthorID, p.CategoryID, p.Title, p.Body,).Scan(&p.ID, &p.CreatedAt)
	if err != nil { http.Error(w, err.Error(),500); return}
	json.NewEncoder(w).Encode(p)
}


func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p model.Post
	err := h.db.Get(&p, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil { http.Error(w, "Not found",404); return}
	json.NewEncoder(w).Encode(p)
}