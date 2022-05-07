package repository

import (
	"database/sql"
	"fmt"
	"user-service/internal/model"
)

type Content struct {
	*sql.DB
}

func NewContent(db *sql.DB) *Content {
	return &Content{db}
}

func (db *Content) GetUserPosts(user_id string) []model.Post {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=$1", user_id).Scan(&username)
	if err != nil {
		fmt.Println("query 2", err)
		return []model.Post{}
	}
	return []model.Post{
		{
			Body: username,
		},
	}
}
