package domain

type RefreshToken struct {
	UserID uint   `gorm:"primarykey"`
	Token  string `gorm:"not null"`
}

type AuthRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
