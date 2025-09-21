package errors

var (
	// ErrInternalServerError — ошибка 500
	ErrInternalServerError = New(0, "Internal Server Error", "")

	// ErrUserNotFound — пользователь не найден
	ErrUserNotFound = New(1000, "User not found", "")

	// ErrInvalidCredentials — неверные учётные данные
	ErrInvalidCredentials = New(1001, "Invalid credentials", "")

	// ErrUserAlreadyExists — пользователь уже существует
	ErrUserAlreadyExists = New(1002, "User already exists", "")

	// ErrAccessTokenExpired — access token просрочен
	ErrAccessTokenExpired = New(1003, "Access token expired", "")

	// ErrRefreshTokenExpired — refresh token просрочен
	ErrRefreshTokenExpired = New(1004, "Refresh token expired", "")

	// ErrUnauthorized — доступ запрещён
	ErrUnauthorized = New(1005, "Unauthorized", "")

	// ErrValidationFailed — ошибка валидации
	ErrValidationFailed = New(1006, "Validation failed", "")

	// ErrUserDeleted - Пользователь был удален
	ErrUserDeleted = New(1007, "User deleted", "")
)
