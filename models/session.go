package models

import (
	"database/sql"
	"fmt"

	"github.com/vityayka/go-zero/rand"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (service *SessionService) Create(userId int) (*Session, error) {
	token, error := rand.SessionToken()
	if error != nil {
		return nil, error
	}

	tokenHash, error := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if error != nil {
		return nil, fmt.Errorf("generating hash failed: %v", error)
	}

	session := Session{
		UserID:    userId,
		Token:     token,
		TokenHash: string(tokenHash),
	}

	fmt.Printf("userId: %d, hash: %s \n", userId, session.TokenHash)
	res := service.DB.QueryRow("INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2) RETURNING id;", userId, session.TokenHash)

	error = res.Scan(&session.ID)
	if error != nil {
		return nil, fmt.Errorf("scan failed: %v", error)
	}

	return &session, nil
}

func (service *SessionService) User(token string) (*User, error) {
	if !rand.IsSessionToken(token) {
		return nil, fmt.Errorf("provided string is not a session token")
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("generating hash failed: %v", err)
	}

	res := service.DB.QueryRow(`
		SELECT u.id user_id, u.email. u.pw_hash FROM sessions s
		LEFT JOIN users u on s.user_ud = u.id
		WHERE s.token_hash = $1`,
		tokenHash,
	)

	user := User{}

	err = res.Scan(&user.ID, &user.Email, &user.PwHash)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("provided session token is expired or a fake")
	}

	return &user, nil
}
