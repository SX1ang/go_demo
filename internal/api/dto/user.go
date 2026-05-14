package dto

type SignUpReq struct {
	Username   string `json:"username" binding:"required,min=3,max=32"`
	Password   string `json:"password" binding:"required,min=8,max=64"`
	RePassword string `json:"re_password" binding:"required,min=8,max=64,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	Gender     int32  `json:"gender" binding:"required,oneof=1 2"` // Gender填0值，required validation失败
}

type LoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type LoginRes struct {
	UserID        int64  `json:"user_id,string"`
	Username      string `json:"username"`
	SessionId     string `json:"session_id"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type LogoutReq struct {
	UserId int64
}
