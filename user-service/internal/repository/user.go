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

func (db *User) GetUserPublicData(userID uuid.UUID) (model.User, error) {
	var u model.User
	err := db.QueryRow("SELECT id, username, avatar_ref FROM users WHERE id =$1", userID).Scan(&u.ID, &u.Name, &u.AvatarRef)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}
