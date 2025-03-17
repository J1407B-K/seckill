package model

type User struct {
	UserId   string  `json:"userId"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}
