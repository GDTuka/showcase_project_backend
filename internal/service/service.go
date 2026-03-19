package service

import (
	"fmt"

	"showcase_project/config"
	"showcase_project/data/dto/user"
	"showcase_project/data/request/auth"
	"showcase_project/data/request/utils"
	e "showcase_project/internal/error_service"
	"showcase_project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(req auth.RegisterRequest) (*int, *TokenDetails, *TokenDetails, e.IAppError)
	Login(req auth.LoginRequest) (*user.User, *TokenDetails, *TokenDetails, e.IAppError)
	RefreshToken(refreshToken string) (*TokenDetails, *TokenDetails, e.IAppError)
}

type Utils interface {
	CheckLoginUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
	CheckPhoneUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
}

type Service struct {
	Auth
	Utils
	JWT
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	jwtService := NewJWTService(cfg)
	return &Service{
		Auth:  NewAuthService(repo.Auth, jwtService),
		Utils: NewUtilsService(repo.Utils),
		JWT:   jwtService,
	}
}

// Auth Service
type AuthService struct {
	repo repository.Auth
	jwt  JWT
}

func NewAuthService(repo repository.Auth, jwt JWT) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(req auth.RegisterRequest) (*int, *TokenDetails, *TokenDetails, e.IAppError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, nil, e.NewAppError(fmt.Errorf("failed to hash password"), 500)
	}

	userId, appErr := s.repo.Register(req, string(hashedPassword))
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	at, rt, appErr := s.jwt.GenerateTokens(*userId)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	return userId, at, rt, nil
}

func (s *AuthService) Login(req auth.LoginRequest) (*user.User, *TokenDetails, *TokenDetails, e.IAppError) {
	u, err := s.repo.GetUserByLogin(req.Login)
	if err != nil {
		if err.Code() == 404 {
			return nil, nil, nil, e.NewAppError(fmt.Errorf("invalid login or password"), 401)
		}
		return nil, nil, nil, err
	}

	cryptErr := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if cryptErr != nil {
		return nil, nil, nil, e.NewAppError(fmt.Errorf("invalid login or password"), 401)
	}

	at, rt, appErr := s.jwt.GenerateTokens(u.ID)
	if appErr != nil {
		return nil, nil, nil, appErr
	}

	return u, at, rt, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenDetails, *TokenDetails, e.IAppError) {
	userId, err := s.jwt.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, nil, err
	}

	return s.jwt.GenerateTokens(userId)
}

// Utils Service
type UtilsService struct {
	repo repository.Utils
}

func NewUtilsService(repo repository.Utils) *UtilsService {
	return &UtilsService{repo: repo}
}

func (s *UtilsService) CheckLoginUnique(req utils.CheckUniqueRequest) (bool, e.IAppError) {
	return s.repo.CheckLoginUnique(req)
}

func (s *UtilsService) CheckPhoneUnique(req utils.CheckUniqueRequest) (bool, e.IAppError) {
	return s.repo.CheckPhoneUnique(req)
}
