package user

type SearchRequest struct {
	Login        string `form:"login"`
	RelationType string `form:"relation_type"`
	Limit        int    `form:"limit,default=10"`
	Offset       int    `form:"offset,default=0"`
}

type CheckSmsCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

type SendSmsCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}
