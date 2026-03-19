package user

type User struct {
	ID           int    `db:"id" json:"id"`
	Login        string `db:"login" json:"login"`
	Phone        string `db:"phone" json:"phone"`
	PasswordHash string `db:"password_hash" json:"-"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}
