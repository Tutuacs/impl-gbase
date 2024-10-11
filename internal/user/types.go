package user

import (
	"time"

	"github.com/Tutuacs/pkg/enums"
)

// TODO: Create types and dtos for user

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Role      enums.Role `json:"role"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
}

type NewUserDto struct {
	Name     string     `json:"name" validation:"required"`
	Role     enums.Role `json:"role" validation:"required,min=0,max=1"`
	Email    string     `json:"email" validation:"required,email"`
	Password string     `json:"-" validation:"required"`
}

type UpdateUserDto struct {
	Name  string     `json:"name"`
	Role  enums.Role `json:"role" validation:"min=0,max=1"`
	Email string     `json:"email" validation:"required,email"`
}
