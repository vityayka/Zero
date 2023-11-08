package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID     int
	Email  string
	PwHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("pw hash generating error: %w", err)
	}

	user := User{
		Email:  email,
		PwHash: string(hash),
	}

	fmt.Printf("Email: %s", email)
	res := us.DB.QueryRow("INSERT INTO users (email, pw_hash) VALUES ($1, $2) returning id", email, user.PwHash)

	err = res.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting user error: %w", err)
	}

	return &user, nil
}

func (us *UserService) Auth(email, password string) (*User, error) {
	user := User{
		Email: strings.ToLower(email),
	}

	row := us.DB.QueryRow("SELECT id, pw_hash FROM users where email = $1;", email)
	err := row.Scan(&user.ID, &user.PwHash)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not found")
	}

	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	fmt.Printf("User: %v", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, fmt.Errorf("password is incorrect")
	}

	return &user, nil
}
