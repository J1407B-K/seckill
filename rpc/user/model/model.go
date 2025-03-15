package model

type User struct {
	Userid   int     `json:"userid" gorm:"primary_key;auto_increment"`
	Username string  `json:"username" gorm:"unique;not null"`
	Password string  `json:"password" gorm:"not null;size:255"`
	Email    *string `json:"email" gorm:"unique"`
}
