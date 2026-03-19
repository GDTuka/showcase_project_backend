package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"showcase_project/data/dto/user"
	"showcase_project/data/request/auth"
	"showcase_project/data/request/utils"
	e "showcase_project/internal/error_service"
	"showcase_project/internal/repository"
)

type Auth interface {
	Register(req auth.RegisterRequest) (*int, e.IAppError)
	Login(req auth.LoginRequest) (*user.User, e.IAppError)
}

type Utils interface {
	CheckLoginUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
	CheckPhoneUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
}

type Service struct {
	Auth
	Utils
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:  NewAuthService(repo.Auth),
		Utils: NewUtilsService(repo.Utils),
	}
}

// Auth Service
type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req auth.RegisterRequest) (*int, e.IAppError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, e.NewAppError(fmt.Errorf("failed to hash password"), 500)
	}

	return s.repo.Register(req, string(hashedPassword))
}

func (s *AuthService) Login(req auth.LoginRequest) (*user.User, e.IAppError) {
	u, err := s.repo.GetUserByLogin(req.Login)
	if err != nil {
		if err.Code() == 404 {
			return nil, e.NewAppError(fmt.Errorf("invalid login or password"), 401)
		}
		return nil, err
	}

	cryptErr := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if cryptErr != nil {
		return nil, e.NewAppError(fmt.Errorf("invalid login or password"), 401)
	}

	return u, nil
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
