package response

type APIResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}
