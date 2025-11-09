package model

import (
	"database/sql"
	"fmt"
	"strings"
	"synk/gateway/app/util"
)

type Users struct {
	db *sql.DB
}

type UsersList struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	CreatedAt string `json:"created_at"`
}

func NewUsers(db *sql.DB) *Users {
	users := Users{db: db}

	return &users
}

func (u *Users) List(id string) ([]UsersList, error) {
	var users []UsersList

	whereList := []string{}
	whereValues := []any{}

	if id != "" {
		whereList = append(whereList, "user_id = ?")
		whereValues = append(whereValues, id)
	}

	where := ""

	if len(whereList) > 0 {
		where = " AND " + strings.Join(whereList, " AND ")
	}

	rows, rowsErr := u.db.Query(
		`SELECT user.user_id, user.user_name, user.user_email, user.created_at
        FROM user
        WHERE user.deleted_at IS NULL `+where, whereValues...,
	)

	if rowsErr != nil {
		return nil, fmt.Errorf("models.users.list: %s", rowsErr.Error())
	}

	defer rows.Close()

	rowsErr = rows.Err()

	if rowsErr != nil {
		return nil, fmt.Errorf("models.users.list: %s", rowsErr.Error())
	}

	for rows.Next() {
		var user UsersList

		exception := rows.Scan(
			&user.UserId,
			&user.UserName,
			&user.UserEmail,
			&user.CreatedAt,
		)

		user.CreatedAt = util.ToTimeBR(user.CreatedAt)

		if exception != nil {
			return nil, fmt.Errorf("models.users.list: %s", exception.Error())
		}

		users = append(users, user)
	}

	return users, nil
}
