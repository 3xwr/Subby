package repository

import (
	"database/sql"
)

type Subscriptions struct {
	*sql.DB
}

func NewSubscriptions(db *sql.DB) *Subscriptions {
	return &Subscriptions{db}
}

func (db *Subscriptions) GetUserSubs(user_id string) ([]string, error) {
	rows, err := db.Query("SELECT subbed_to_user_id FROM user_subs WHERE user_id=$1", user_id)
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
