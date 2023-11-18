package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

type ResetToken struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type ResetTokenService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (service *ResetTokenService) Create(email string) (*ResetToken, error) {
	return nil, fmt.Errorf("TODO: implement this func")
}

func (service *ResetTokenService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: implement this func")
}

func (service *ResetTokenService) User(token string) (*User, error) {
	tokenHash := service.generateTokenHash(token)

	service.DB.Stats()
	res := service.DB.QueryRow(`
		SELECT u.id user_id, u.email, u.pw_hash FROM reset_tokens rt
		LEFT JOIN users u on rt.user_id = u.id
		WHERE rt.token_hash = $1`,
		tokenHash,
	)

	user := User{}

	err := res.Scan(&user.ID, &user.Email, &user.PwHash)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("provided reset token is expired or a fake")
	}

	return &user, nil
}

func (service *ResetTokenService) generateTokenHash(token string) string {
	bytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(bytes[:])
}
