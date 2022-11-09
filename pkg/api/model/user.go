package model

import (
	"github.com/google/uuid"
	"nebula/pkg/security"
)

type User struct {
	Uid          uuid.UUID
	Email        string
	PasswordHash string
	Key          string
}

type Session struct {
	Uid      uuid.UUID
	Sid      uuid.UUID
	Hwid     security.HWID
	Verified bool
}
