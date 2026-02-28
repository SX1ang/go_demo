package dto

type SignUpReq struct {
	Username   string `json:"username" binding:"required,min=3,max=32"`
	Password   string `json:"password" binding:"required,min=8,max=64"`
	RePassword string `json:"re_password" binding:"required,min=8,max=64,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	Gender     int32  `json:"gender" binding:"required,oneof=1 2"` // Gender填0值，required validation失败
}
