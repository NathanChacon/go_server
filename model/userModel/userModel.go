package userModel

import (
	"fmt"

	"github.com/google/uuid"
	"test.com/events/database"
)

type User struct {
	ID       string `json:"-"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func EmailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	err := database.DB.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func CreateUser(user User) (string, error) {
	id := uuid.New().String()
	_, err := database.DB.Exec(
		"INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)",
		id, user.Name, user.Email, user.Password,
	)

	if err != nil {
		fmt.Print("create user error", err)
		return "", err
	}

	return id, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT id, email, password FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		fmt.Print("get user error", err)
		return User{}, err
	}

	return user, nil
}
