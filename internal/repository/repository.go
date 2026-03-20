package repository

import (
	"database/sql"
	"fmt"

	"showcase_project/data/dto/user"
	"showcase_project/data/request/auth"
	"showcase_project/data/request/utils"
	e "showcase_project/internal/error_service"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	Register(req auth.RegisterRequest) (*int, e.IAppError)
	GetUserByLogin(login string) (*user.User, e.IAppError)
	GetUserByPhone(phone string) (*user.User, e.IAppError)
}

type Utils interface {
	CheckLoginUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
	CheckPhoneUnique(req utils.CheckUniqueRequest) (bool, e.IAppError)
}

type Repository struct {
	Auth
	Utils
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:  NewAuthRepository(db),
		Utils: NewUtilsRepository(db),
		User:  NewUserRepository(db),
	}
}

// Auth Repository
type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Register(req auth.RegisterRequest) (*int, e.IAppError) {
	query := `INSERT INTO user (login, phone) VALUES (?, ?)`
	result, err := r.db.Exec(query, req.Login, req.Phone)
	if err != nil {
		return nil, e.NewAppError(fmt.Errorf("failed to insert user: %w", err), 409) // Conflict
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, e.NewAppError(err, 500)
	}

	intId := int(id)
	return &intId, nil
}

func (r *AuthRepository) GetUserByLogin(login string) (*user.User, e.IAppError) {
	var u user.User
	err := r.db.Get(&u, "SELECT id, login, phone, IFNULL(avatar, '') as avatar, created_at FROM user WHERE login = ?", login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NewAppError(fmt.Errorf("user not found"), 404)
		}
		return nil, e.NewAppError(err, 500)
	}
	return &u, nil
}

func (r *AuthRepository) GetUserByPhone(phone string) (*user.User, e.IAppError) {
	var u user.User
	err := r.db.Get(&u, "SELECT id, login, phone, IFNULL(avatar, '') as avatar, created_at FROM user WHERE phone = ?", phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NewAppError(fmt.Errorf("user not found"), 404)
		}
		return nil, e.NewAppError(err, 500)
	}
	return &u, nil
}

// Utils Repository
type UtilsRepository struct {
	db *sqlx.DB
}

func NewUtilsRepository(db *sqlx.DB) *UtilsRepository {
	return &UtilsRepository{db: db}
}

func (r *UtilsRepository) CheckLoginUnique(req utils.CheckUniqueRequest) (bool, e.IAppError) {
	var count int
	err := r.db.Get(&count, "SELECT count(id) FROM user WHERE login = ?", req.Value)
	if err != nil && err != sql.ErrNoRows {
		return false, e.NewAppError(err, 500)
	}
	return count == 0, nil
}

func (r *UtilsRepository) CheckPhoneUnique(req utils.CheckUniqueRequest) (bool, e.IAppError) {
	var count int
	err := r.db.Get(&count, "SELECT count(id) FROM user WHERE phone = ?", req.Value)
	if err != nil && err != sql.ErrNoRows {
		return false, e.NewAppError(err, 500)
	}
	return count == 0, nil
}
