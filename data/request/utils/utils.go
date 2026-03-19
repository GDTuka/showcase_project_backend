package utils

type CheckUniqueRequest struct {
	Value string `json:"value" binding:"required"`
}
