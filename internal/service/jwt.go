package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"showcase_project/config"
	e "showcase_project/internal/error_service"
)

type TokenDetails struct {
	Token     string
	ExpiresAt int64
}

type JWT interface {
	GenerateTokens(userId int) (*TokenDetails, *TokenDetails, e.IAppError)
	ValidateToken(tokenString string, expectedType string) (int, e.IAppError)
}

type JWTService struct {
	secret string
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secret: cfg.JWT.Secret,
	}
}

func (s *JWTService) GenerateTokens(userId int) (*TokenDetails, *TokenDetails, e.IAppError) {
	// Access token - 1 week (as requested)
	atExpiration := time.Now().Add(time.Hour * 24 * 7).Unix()
	atClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     atExpiration,
		"type":    "access",
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.secret))
	if err != nil {
		return nil, nil, e.NewAppError(fmt.Errorf("failed to generate access token"), 500)
	}

	// Refresh token - 4 weeks
	rtExpiration := time.Now().Add(time.Hour * 24 * 28).Unix()
	rtClaims := jwt.MapClaims{
		"user_id": userId,
		"exp":     rtExpiration,
		"type":    "refresh",
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.secret))
	if err != nil {
		return nil, nil, e.NewAppError(fmt.Errorf("failed to generate refresh token"), 500)
	}

	return &TokenDetails{
			Token:     accessToken,
			ExpiresAt: atExpiration,
		}, &TokenDetails{
			Token:     refreshToken,
			ExpiresAt: rtExpiration,
		}, nil
}

func (s *JWTService) ValidateToken(tokenString string, expectedType string) (int, e.IAppError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return 0, e.NewAppError(fmt.Errorf("invalid token"), 401)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["type"] != expectedType {
			return 0, e.NewAppError(fmt.Errorf("invalid token type"), 401)
		}
		userId := int(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, e.NewAppError(fmt.Errorf("invalid token claims"), 401)
}
