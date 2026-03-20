package service

import (
	"showcase_project/data/dto/profile"
	"showcase_project/data/dto/user"
	reqUser "showcase_project/data/request/user"
	e "showcase_project/internal/error_service"
	"showcase_project/internal/repository"
)

type User interface {
	SendSmsCode(phone string) e.IAppError
	CheckSmsCode(phone string, code string) (bool, e.IAppError)
	GetCurrentUser(userId int) (*profile.UserWithProfile, e.IAppError)
	GetUserById(userId int) (*profile.UserWithProfile, e.IAppError)
	SearchUsers(req reqUser.SearchRequest, currentUserId int) ([]user.User, e.IAppError)
}

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SendSmsCode(phone string) e.IAppError {
	return s.repo.SendSmsCode(phone)
}

func (s *UserService) CheckSmsCode(phone string, code string) (bool, e.IAppError) {
	return s.repo.CheckSmsCode(phone, code)
}

func (s *UserService) GetCurrentUser(userId int) (*profile.UserWithProfile, e.IAppError) {
	u, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.GetUserProfile(userId)
	if err != nil {
		return nil, err
	}

	return &profile.UserWithProfile{
		User:    *u,
		Profile: p,
	}, nil
}

func (s *UserService) GetUserById(userId int) (*profile.UserWithProfile, e.IAppError) {
	u, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.GetUserProfile(userId)
	if err != nil {
		return nil, err
	}

	return &profile.UserWithProfile{
		User:    *u,
		Profile: p,
	}, nil
}

func (s *UserService) SearchUsers(req reqUser.SearchRequest, currentUserId int) ([]user.User, e.IAppError) {
	return s.repo.SearchUsers(req, currentUserId)
}
