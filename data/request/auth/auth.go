package auth

type RegisterRequest struct {
	Login string `json:"login" binding:"required,alpha,min=3,max=30"`
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
}

type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
}
