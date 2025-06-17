package model

import "time"

type Comment struct {
	PostID    int       `bson:"post_id" json:"post_id"`
	UserID    int       `bson:"user_id" json:"user_id"`
	Username  string    `bson:"username" json:"username"`
	Text      string    `bson:"text" json:"text"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
