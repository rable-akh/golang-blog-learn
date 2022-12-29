package responses

import "akh/blog/models"

type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type SignInResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type SignDataResponse struct {
	Token string `json:"token"`
	models.User
}
