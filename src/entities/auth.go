package entities

import "time"

type CreateAuth struct {
	Identifier string `json:"identifier" validate:"email"`
	Password   string `json:"password"`
}

type AuthEntity struct {
	Id         int        `json:"id"`
	Identifier string     `json:"identifier"`
	Password   string     `json:"-"`
	IsDisabled bool       `json:"isDisabled"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	SignedInAt *time.Time `json:"-"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
