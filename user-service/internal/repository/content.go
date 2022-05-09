package repository

import (
	"database/sql"
	"user-service/internal/model"

	"github.com/lib/pq"
)

type Content struct {
	*sql.DB
}

func NewContent(db *sql.DB) *Content {
	return &Content{db}
}

func (db *Content) GetUserSubs(userID string) ([]string, error) {
	rows, err := db.Query("SELECT subbed_to_user_id FROM user_subs WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSubs []string

	for rows.Next() {
		var sub string
		if err := rows.Scan(&sub); err != nil {
			return nil, err
		}
		userSubs = append(userSubs, sub)
	}

	return userSubs, nil
}

func (db *Content) GetUserFeed(userID string, amount int) ([]model.Post, error) {
	subs, err := db.GetUserSubs(userID)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT post_id, poster_id, body, paywall_locked, paywall_tier, image_ref FROM posts WHERE poster_id=ANY($1)", pq.Array(subs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.PostID, &post.PosterID, &post.Body, &post.PaywallLocked, &post.PaywallTier, &post.ImageRef); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
