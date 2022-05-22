package repository

import (
	"database/sql"
	"user-service/internal/model"

	"github.com/google/uuid"
)

type User struct {
	*sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db}
}

func (db *User) GetUserIDByName(name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.QueryRow("SELECT id FROM users WHERE username =$1", name).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (db *User) GetUserPublicData(userID uuid.UUID) (model.User, error) {
	var u model.User
	err := db.QueryRow("SELECT id, username, avatar_ref FROM users WHERE id =$1", userID).Scan(&u.ID, &u.Name, &u.AvatarRef)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (db *User) ChangeUserAvatar(userID uuid.UUID, avatarRef string) error {
	var u model.User
	err := db.QueryRow("SELECT id, username, avatar_ref FROM users WHERE id =$1", userID).Scan(&u.ID, &u.Name, &u.AvatarRef)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE users SET avatar_ref=$1 WHERE id=$2", avatarRef, userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *User) GetFullUserPublicData(userID uuid.UUID) (model.User, error) {
	var u model.User
	sqlQuery := `SELECT u.id, u.username, u.avatar_ref, c.cnt
	FROM users u
	INNER JOIN (select subbed_to_user_id,count(subbed_to_user_id) as cnt
				FROM user_subs
				group by subbed_to_user_id) c on u.id=c.subbed_to_user_id
	where u.id=$1`
	rows, err := db.Query(sqlQuery, userID)
	if err != nil {
		return u, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.Name, &u.AvatarRef, &u.SubscriberCount); err != nil {
			return u, err
		}
	}

	return u, nil
}
