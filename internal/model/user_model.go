package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"emai"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"emai"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email    string `json:"emai"`
	Password string `json:"password"`
}
