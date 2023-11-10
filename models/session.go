package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/vityayka/go-zero/rand"
)

const (
	MinSessionTokenSize = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB         *sql.DB
	TokenBytes int
}

func (service *SessionService) Create(userId int) (*Session, error) {
	token, err := service.SessionToken()
	if err != nil {
		return nil, err

	}

	session := Session{
		UserID:    userId,
		Token:     token,
		TokenHash: service.generateTokenHash(token),
	}

	res := service.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash) 
		VALUES ($1, $2)  
		ON CONFLICT (user_id) DO UPDATE SET token_hash = $2
		RETURNING id;
	`, userId, session.TokenHash)

	err = res.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %v", err)
	}

	return &session, nil
}

func (service *SessionService) User(token string) (*User, error) {
	tokenHash := service.generateTokenHash(token)

	service.DB.Stats()
	res := service.DB.QueryRow(`
		SELECT u.id user_id, u.email, u.pw_hash FROM sessions s
		LEFT JOIN users u on s.user_id = u.id
		WHERE s.token_hash = $1`,
		tokenHash,
	)

	user := User{}

	err := res.Scan(&user.ID, &user.Email, &user.PwHash)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("provided session token is expired or a fake")
	}

	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	hash := service.generateTokenHash(token)

	_, err := service.DB.Exec("DELETE FROM sessions WHERE token_hash = $1", hash)

	return err
}

func (service *SessionService) SessionToken() (string, error) {
	size := service.TokenBytes
	if size < MinSessionTokenSize {
		size = MinSessionTokenSize
	}
	return rand.String(size)
}

func (service *SessionService) generateTokenHash(token string) string {
	bytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(bytes[:])
}
