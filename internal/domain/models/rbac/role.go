package rbac

type RoleName string

const (
	Admin = "Admin"
	User  = "User"
)

type Role struct {
	Name        RoleName
	Description string
	Permissions []Permission
}
