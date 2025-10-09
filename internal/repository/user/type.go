package user

type User struct {
	ID              int64
	Name            string
	Password        string
	Email           string
	ImageURL        string
	Phone           string
	IsPhoneVerified bool
	Status          Status
	Type            Type
}

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
