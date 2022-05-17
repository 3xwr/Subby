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

type post struct {
	postID   uuid.UUID
	posterID uuid.UUID
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

func (db *Content) RemoveUserSubscription(currentUser string, userToUnsubscribeFrom string) error {
	err := db.checkIfUserExists(userToUnsubscribeFrom)
	if err != nil {
		return err
	}
	subs, err := db.GetUserSubs(currentUser)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		if userToUnsubscribeFrom == sub {
			_, err = db.Exec("DELETE FROM user_subs WHERE user_id=$1 AND subbed_to_user_id=$2", currentUser, userToUnsubscribeFrom)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *Content) CheckSubscribe(subbingUser uuid.UUID, checkUser uuid.UUID) (bool, error) {
	var sUser uuid.UUID
	var cUser uuid.UUID
	err := db.QueryRow("SELECT user_id, subbed_to_user_id FROM user_subs WHERE user_id = $1 AND subbed_to_user_id = $2", subbingUser, checkUser).Scan(&sUser, &cUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *Content) SaveNewPost(post model.Post) error {
	_, err := db.Exec("INSERT INTO posts (post_id, posted_at, poster_id, body, membership_locked, membership_tier, image_ref) VALUES ($1, $2, $3, $4, $5, $6, $7)", post.PostID, post.PostedAt, post.PosterID, post.Body, post.MembershipLocked, post.MembershipTier, post.ImageRef)
	if err != nil {
		return err
	}
	return nil
}

func (db *Content) DeletePostFromDB(userID string, postID string) error {
	var p post
	err := db.QueryRow("SELECT poster_id, post_id FROM posts WHERE poster_id = $1 AND post_id = $2", userID, postID).Scan(&p.posterID, &p.postID)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM posts WHERE post_id=$1", postID)
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
	sqlQuery := `SELECT post_id, posted_at, poster_id, body, membership_locked, membership_tier, image_ref, username, avatar_ref 
	FROM posts
	INNER JOIN users ON posts.poster_id = users.id
	WHERE poster_id=ANY($1) LIMIT $2`
	rows, err := db.Query(sqlQuery, pq.Array(subs), amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.PostID, &post.PostedAt, &post.PosterID, &post.Body, &post.MembershipLocked, &post.MembershipTier, &post.ImageRef, &post.PosterUsername, &post.PosterAvatarRef); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Content) GetUserPosts(posterID uuid.UUID, amount int) ([]model.Post, error) {
	sqlQuery := `SELECT post_id, posted_at, poster_id, body, membership_locked, membership_tier, image_ref, username, avatar_ref 
	FROM posts
	INNER JOIN users ON posts.poster_id = users.id
	WHERE poster_id=$1 LIMIT $2`
	rows, err := db.Query(sqlQuery, posterID, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.PostID, &post.PostedAt, &post.PosterID, &post.Body, &post.MembershipLocked, &post.MembershipTier, &post.ImageRef, &post.PosterUsername, &post.PosterAvatarRef); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
