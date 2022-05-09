package repository

import (
	"database/sql"
	"user-service/internal/model"

	"github.com/google/uuid"
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

func (db *Content) checkIfUserExists(userID string) error {
	var u user
	err := db.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&u.id)
	if err != nil {
		return err
	}
	return nil
}

func (db *Content) AddUserSubscription(currentUser string, userToSubsrcibeTo string) error {
	err := db.checkIfUserExists(userToSubsrcibeTo)
	if err != nil {
		return err
	}
	subs, err := db.GetUserSubs(currentUser)
	if err != nil {
		return err
	}
	//If user is already subscribed just exit
	for _, sub := range subs {
		if userToSubsrcibeTo == sub {
			return nil
		}
	}
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO user_subs (sub_id, user_id, subbed_to_user_id) VALUES ($1, $2, $3)", tokenID, currentUser, userToSubsrcibeTo)
	if err != nil {
		return err
	}
	return nil
}

func (db *Content) GetUserFeed(userID string, amount int) ([]model.Post, error) {
	subs, err := db.GetUserSubs(userID)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT post_id, posted_at, poster_id, body, paywall_locked, paywall_tier, image_ref FROM posts WHERE poster_id=ANY($1) LIMIT $2", pq.Array(subs), amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.PostID, &post.PostedAt, &post.PosterID, &post.Body, &post.PaywallLocked, &post.PaywallTier, &post.ImageRef); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
