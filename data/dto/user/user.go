package user

type User struct {
	ID        int     `db:"id" json:"id"`
	Login     string  `db:"login" json:"login"`
	Phone     string  `db:"phone" json:"phone"`
	Avatar    *string `db:"avatar" json:"avatar"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}
