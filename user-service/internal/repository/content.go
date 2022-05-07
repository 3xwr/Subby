package repository

import (
	"database/sql"
	"user-service/internal/model"
)

type Content struct {
	*sql.DB
}

func NewContent(db *sql.DB) *Content {
	return &Content{db}
}

func (db *Content) GetUserPosts(user_id string) ([]model.Post, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=$1", user_id).Scan(&username)
	if err != nil {
		return []model.Post{}, err
	}
	return []model.Post{
		{
			Body: username,
		},
	}, nil
}
