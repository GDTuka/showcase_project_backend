package repository

import (
	"database/sql"
	"fmt"

	"showcase_project/data/dto/profile"
	"showcase_project/data/dto/user"
	reqUser "showcase_project/data/request/user"
	e "showcase_project/internal/error_service"

	"github.com/jmoiron/sqlx"
)

type User interface {
	SendSmsCode(phone string) e.IAppError
	CheckSmsCode(phone string, code string) (bool, e.IAppError)
	GetUserById(userId int) (*user.User, e.IAppError)
	GetUserProfile(userId int) (*profile.UserProfile, e.IAppError)
	SearchUsers(req reqUser.SearchRequest, currentUserId int) ([]user.User, e.IAppError)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SendSmsCode(phone string) e.IAppError {
	tx, err := r.db.Beginx()
	if err != nil {
		return e.NewAppError(err, 500)
	}
	defer tx.Rollback()

	// Clean up expired codes (older than 1 hour)
	_, err = tx.Exec("DELETE FROM sms_code WHERE created_at < datetime('now', '-1 hour')")
	if err != nil {
		return e.NewAppError(err, 500)
	}

	// Delete existing code for this phone
	_, err = tx.Exec("DELETE FROM sms_code WHERE phone = ?", phone)
	if err != nil {
		return e.NewAppError(err, 500)
	}

	// Insert new code "000000"
	_, err = tx.Exec("INSERT INTO sms_code (phone, code) VALUES (?, '000000')", phone)
	if err != nil {
		return e.NewAppError(err, 500)
	}

	if err = tx.Commit(); err != nil {
		return e.NewAppError(err, 500)
	}

	return nil
}

func (r *UserRepository) CheckSmsCode(phone string, code string) (bool, e.IAppError) {
	// First clean up expired codes
	_, err := r.db.Exec("DELETE FROM sms_code WHERE created_at < datetime('now', '-1 hour')")
	if err != nil {
		return false, e.NewAppError(err, 500)
	}

	var count int
	err = r.db.Get(&count, "SELECT COUNT(*) FROM sms_code WHERE phone = ? AND code = ?", phone, code)
	if err != nil && err != sql.ErrNoRows {
		return false, e.NewAppError(err, 500)
	}

	if count > 0 {
		// Code is valid, delete it so it can't be used again
		_, _ = r.db.Exec("DELETE FROM sms_code WHERE phone = ?", phone)
		return true, nil
	}

	return false, nil
}

func (r *UserRepository) GetUserById(userId int) (*user.User, e.IAppError) {
	var u user.User
	err := r.db.Get(&u, "SELECT id, login, phone, avatar, created_at FROM user WHERE id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NewAppError(fmt.Errorf("user not found"), 404)
		}
		return nil, e.NewAppError(err, 500)
	}
	return &u, nil
}

func (r *UserRepository) GetUserProfile(userId int) (*profile.UserProfile, e.IAppError) {
	var p profile.UserProfile
	err := r.db.Get(&p, "SELECT * FROM user_profile WHERE user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No profile yet
		}
		return nil, e.NewAppError(err, 500)
	}
	return &p, nil
}

func (r *UserRepository) SearchUsers(req reqUser.SearchRequest, currentUserId int) ([]user.User, e.IAppError) {
	query := "SELECT u.id, u.login, u.phone, u.password_hash, IFNULL(u.avatar, '') as avatar, u.created_at FROM user u "
	args := []interface{}{}

	if req.RelationType != "" {
		query += "JOIN user_relation ur ON u.id = ur.related_user_id AND ur.user_id = ? AND ur.relation_type = ? "
		args = append(args, currentUserId, req.RelationType)
	}

	query += "WHERE u.id != ? "
	args = append(args, currentUserId)

	if req.Login != "" {
		query += "AND u.login LIKE ? "
		args = append(args, "%"+req.Login+"%")
	}

	query += "LIMIT ? OFFSET ?"
	args = append(args, req.Limit, req.Offset)

	var users []user.User
	err := r.db.Select(&users, query, args...)
	if err != nil {
		return nil, e.NewAppError(err, 500)
	}
	
	if users == nil {
		users = []user.User{}
	}

	return users, nil
}
