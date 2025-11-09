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
	Type       Type

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) IsEmpty() bool {
	return u.ID == 0
}

type Verified uint64

// Binary representation of verification status
// VerifiedNone  = 0000
// VerifiedPhone = 0001
// VerifiedEmail = 0010
const (
	VerifiedNone  Verified = 0
	VerifiedPhone Verified = 1
	VerifiedEmail Verified = 2
)

type Status int8

const (
	StatusActive Status = iota
	StatusBlocked
)

type Type int8

const (
	TypeAdmin Type = iota
	TypeUser
)

func (t Type) String() string {
	switch t {
	case TypeAdmin:
		return "admin"
	case TypeUser:
		return "user"
	default:
		return "unknown"
	}
}
