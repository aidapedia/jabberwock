package policy

import "time"

type ServiceType string

const (
	HTTPServiceType ServiceType = "http"
)

type Role struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Resource struct {
	ID        int64
	Type      ServiceType
	Method    string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetRoleByUserIDResp []Role

type LoadPolicyResponse []Policy

type Policy struct {
	Role   string
	Path   string
	Type   string
	Method string
}
