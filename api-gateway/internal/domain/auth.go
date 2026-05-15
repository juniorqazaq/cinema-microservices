package domain

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type AuthUser struct {
	ID    string
	Email string
	Role  Role
}

func (u AuthUser) IsAdmin() bool {
	return u.Role == RoleAdmin
}
