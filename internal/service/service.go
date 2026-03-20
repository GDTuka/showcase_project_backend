package service

import (
	"fmt"

	"showcase_project/config"
	"showcase_project/data/dto/user"
	"showcase_project/data/request/auth"
	"showcase_project/data/request/utils"
	e "showcase_project/internal/error_service"
	"showcase_project/internal/repository"
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
	User
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	jwtService := NewJWTService(cfg)
	userService := NewUserService(repo.User)
	return &Service{
		Auth:  NewAuthService(repo.Auth, jwtService, repo.User),
		Utils: NewUtilsService(repo.Utils),
		JWT:   jwtService,
		User:  userService,
	}
}

// Auth Service
type AuthService struct {
	repo     repository.Auth
	jwt      JWT
	userRepo repository.User
}

func NewAuthService(repo repository.Auth, jwt JWT, userRepo repository.User) *AuthService {
	return &AuthService{repo: repo, jwt: jwt, userRepo: userRepo}
}

func (s *AuthService) Register(req auth.RegisterRequest) (*int, *TokenDetails, *TokenDetails, e.IAppError) {
	// Verify SMS code
	isValid, appErr := s.userRepo.CheckSmsCode(req.Phone, req.Code)
	if appErr != nil {
		return nil, nil, nil, appErr
	}
	if !isValid {
		return nil, nil, nil, e.NewAppError(fmt.Errorf("invalid or expired SMS code"), 400)
	}

	userId, appErr := s.repo.Register(req)
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
	// Verify SMS code
	isValid, appErr := s.userRepo.CheckSmsCode(req.Phone, req.Code)
	if appErr != nil {
		return nil, nil, nil, appErr
	}
	if !isValid {
		return nil, nil, nil, e.NewAppError(fmt.Errorf("invalid or expired SMS code"), 400)
	}

	u, err := s.repo.GetUserByPhone(req.Phone)
	if err != nil {
		if err.Code() == 404 {
			return nil, nil, nil, e.NewAppError(fmt.Errorf("user not found"), 404)
		}
		return nil, nil, nil, err
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
