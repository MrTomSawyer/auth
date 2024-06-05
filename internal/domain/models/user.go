package models

import "github.com/MrTomSawyer/sso/internal/domain/models/rbac"

type User struct {
	ID           string
	Email        string
	PasswordHash []byte
	Role         rbac.Role
}
