package models

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       int
	Email    string
	Username sql.NullString
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(email string) (*User, error) {
	row := m.DB.QueryRow(
		"INSERT INTO users (email) VALUES ($1) RETURNING id, email, username",
		email,
	)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username)
	return &user, err
}

func (m *UserModel) GetListWithUsername() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, email, username FROM users WHERE username IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Email, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	row := m.DB.QueryRow("SELECT id, email, username FROM users WHERE username = $1", username)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	row := m.DB.QueryRow("SELECT id, email, username FROM users WHERE email = $1", email)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	row := m.DB.QueryRow("SELECT id, email, username FROM users WHERE id = $1", id)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Update(id int, username string) (*User, error) {
	row := m.DB.QueryRow(
		"UPDATE users SET username = $1 WHERE id = $2 RETURNING id, email, username",
		username, id,
	)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
