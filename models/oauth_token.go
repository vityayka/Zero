package models

import (
	"database/sql"
	"fmt"
	"time"
)

type OAuthToken struct {
	ID        uint
	UserID    int
	Token     string
	ExpiresAt time.Time
	Provider  string
}

type OAuthService struct {
	DB *sql.DB
}

func (service *OAuthService) Create(userId int, token string, expiresAt time.Time, provider string) (*OAuthToken, error) {
	row := service.DB.QueryRow(`
		INSERT INTO oauth_tokens (token, user_id, expires_at, provider) VALUES ($1, $2, $3, $4) RETURNING id`,
		token, userId, expiresAt, provider,
	)

	var oAuthToken OAuthToken
	err := row.Scan(&oAuthToken.ID)

	if err != nil {
		return nil, fmt.Errorf("inserting oauth_token: %v", err)
	}

	oAuthToken.UserID = userId
	oAuthToken.Token = token
	oAuthToken.Provider = provider
	oAuthToken.ExpiresAt = expiresAt

	return &oAuthToken, nil
}
