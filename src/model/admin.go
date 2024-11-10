package model

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
 
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Username string `json:"username"`
	Message string `json:"message"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

