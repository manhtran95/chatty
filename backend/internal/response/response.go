package response

type StatusType string

const (
	StatusSuccess StatusType = "success"
	StatusError   StatusType = "error"
)

type ErrorType string

const (
	ErrorBadRequest          ErrorType = "badRequest"
	ErrorInternalServerError ErrorType = "internalServerError"
	ErrorInvalidCredentials  ErrorType = "invalidCredentials"
	ErrorAccessTokenExpired  ErrorType = "accessTokenExpired"
	ErrorRefreshTokenExpired ErrorType = "refreshTokenExpired"
	ErrorUnauthorized        ErrorType = "unauthorized"
	ErrorForbidden           ErrorType = "forbidden"
)

type APIResponse[T any] struct {
	Data    T          `json:"data,omitempty"`
	Message string     `json:"message,omitempty"`
	Status  StatusType `json:"status,omitempty"`
	Error   ErrorType  `json:"error,omitempty"`
}
