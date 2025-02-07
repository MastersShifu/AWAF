package dto

type UserDTO struct {
	UsernameOrEmail string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}
