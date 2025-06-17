package model

import "time"

type Post struct {
	ID         int       `db:"id" json:"id"`
	AuthorID   int       `db:"author_id" json:"author_id"`
	CategoryID int       `db:"category_id" json:"category_id"`
	Title      string    `db:"title" json:"title"`
	Body       string    `db:"body" json:"body"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
