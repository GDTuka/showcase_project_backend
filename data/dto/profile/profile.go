package profile

import "showcase_project/data/dto/user"

type UserProfile struct {
	UserID         int    `db:"user_id" json:"user_id"`
	FirstName      *string `db:"first_name" json:"first_name"`
	LastName       *string `db:"last_name" json:"last_name"`
	MiddleName     *string `db:"middle_name" json:"middle_name"`
	Status         *string `db:"status" json:"status"`
	PrivateProfile bool   `db:"private_profile" json:"private_profile"`
	BirthDate      *string `db:"birth_date" json:"birth_date"`
	Gender         *string `db:"gender" json:"gender"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}

type UserWithProfile struct {
	user.User
	Profile *UserProfile `json:"profile"`
}
