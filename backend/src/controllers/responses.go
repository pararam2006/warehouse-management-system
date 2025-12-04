package controllers

// ErrorResponse — единый формат ошибки для всех HTTP-контроллеров.
type ErrorResponse struct {
	Error string `json:"error"`
}


