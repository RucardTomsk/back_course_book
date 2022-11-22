package model

type User struct {
	FIO      string `json:"fio" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Guid     string `json:"guid"`
}
