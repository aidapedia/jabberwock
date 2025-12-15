package user

import "time"

type User struct {
	ID         int64
	Name       string
	Password   string
	Email      string
	AvatarURL  string
	Phone      string
	IsVerified Verified
	Status     Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

type Verified uint64

// Binary representation of verification status
const (
	VerifiedNone  Verified = 0 // 0000
	VerifiedPhone Verified = 1 // 0001
	VerifiedEmail Verified = 2 // 0010
)

type Status int8

const (
	StatusActive Status = iota
	StatusBlocked
)
