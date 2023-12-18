package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/vityayka/go-zero/rand"
)

const (
	MinResetTokenSize    = 32
	DefaultTokenDuration = time.Hour
)

type ResetToken struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
	UserService   *UserService
}

func (service *PasswordResetService) Create(email string) (*ResetToken, error) {
	user, err := service.UserService.Find(email)
	if err != nil {
		return nil, err
	}
	token, err := service.ResetToken()
	if err != nil {
		return nil, err
	}

	duration := service.Duration
	if duration == 0 {
		duration = DefaultTokenDuration
	}

	resetToken := ResetToken{
		UserID:    user.ID,
		Token:     token,
		TokenHash: service.generateTokenHash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	res := service.DB.QueryRow(`
		INSERT INTO reset_tokens (user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3)  
		ON CONFLICT (user_id) DO UPDATE SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, user.ID, resetToken.TokenHash, resetToken.ExpiresAt)

	err = res.Scan(&resetToken.ID)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %v", err)
	}

	return &resetToken, nil
}

func (service *PasswordResetService) Consume(token ResetToken) error {
	token.TokenHash = service.generateTokenHash(token.Token)
	_, err := service.DB.Exec("DELETE FROM reset_tokens WHERE token_hash = $1", token.TokenHash)
	return err
}

func (service *PasswordResetService) User(token string) (*User, error) {
	tokenHash := service.generateTokenHash(token)

	res := service.DB.QueryRow(`
		SELECT u.id user_id, u.email, u.pw_hash FROM reset_tokens rt
		LEFT JOIN users u on rt.user_id = u.id
		WHERE rt.token_hash = $1 and rt.expires_at >= $2`,
		tokenHash, time.Now(),
	)

	user := User{}

	err := res.Scan(&user.ID, &user.Email, &user.PwHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("provided reset token is expired or a fake")
	}

	return &user, nil
}

func (service *PasswordResetService) ResetToken() (string, error) {
	size := service.BytesPerToken
	if size < MinResetTokenSize {
		size = MinResetTokenSize
	}
	return rand.String(size)
}

func (service *PasswordResetService) generateTokenHash(token string) string {
	bytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(bytes[:])
}
