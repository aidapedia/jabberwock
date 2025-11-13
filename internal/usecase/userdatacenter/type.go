package userdatacenter

import (
	userRepo "github.com/aidapedia/jabberwock/internal/repository/user"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (u *User) Transform(user userRepo.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.Phone = user.Phone
	u.Email = user.Email
}
