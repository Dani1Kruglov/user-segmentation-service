package storage

import (
	"avito-task/internal/model"
	"database/sql"
	"encoding/json"
	"log"
)

type UsersPostgresStorage struct {
	Db *sql.DB
}

type UserJSON struct {
	Id   int64  `json:"user_id"`
	Name string `json:"user_name"`
}

func (u *UsersPostgresStorage) GetAllUsers(conn *sql.DB) ([]byte, error) {
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

	allUsersJSON, err := json.Marshal(allUsers)
	if err != nil {
		return nil, err
	}
	return allUsersJSON, nil
}

func (u *UsersPostgresStorage) StoreUser(userJSON []byte, conn *sql.DB) error {
	var user UserJSON
	err := json.Unmarshal(userJSON, &user)
	if err != nil {
		return err
	}

	_, err = conn.Exec(`INSERT INTO users (name) VALUES ($1)`, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersPostgresStorage) UpdateUser(userJSON []byte, conn *sql.DB) error {
	var user UserJSON
	err := json.Unmarshal(userJSON, &user)
	if err != nil {
		return err
	}

	_, err = conn.Exec(`UPDATE users SET name = $2 WHERE id = $1`, user.Id, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersPostgresStorage) DeleteUser(userJSON []byte, conn *sql.DB) error {
	var user UserJSON
	err := json.Unmarshal(userJSON, &user)
	if err != nil {
		return err
	}

	_, err = conn.Exec(`DELETE FROM users WHERE id = $1`, user.Id)
	if err != nil {
		return err
	}
	return nil
}
