package storage

import (
	"avito-task/internal/model"
	"database/sql"
	"log"
)

type UsersPostgresStorage struct {
	Db *sql.DB
}

func (u *UsersPostgresStorage) GetAllUsers(conn *sql.DB) ([]model.User, error) {
	rows, err := conn.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var allUsers []model.User
	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			log.Fatal(err)
		}
		allUsers = append(allUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (u *UsersPostgresStorage) StoreUser(userName string, conn *sql.DB) error {
	_, err := conn.Exec(`INSERT INTO users (name) VALUES ($1)`, userName)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersPostgresStorage) UpdateUser(user model.User, conn *sql.DB) error {
	_, err := conn.Exec(`UPDATE users SET name = $2 WHERE id = $1`, user.Id, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersPostgresStorage) DeleteUser(userId int64, conn *sql.DB) error {
	_, err := conn.Exec(`DELETE FROM users WHERE id = $1`, userId)
	if err != nil {
		return err
	}
	return nil
}
