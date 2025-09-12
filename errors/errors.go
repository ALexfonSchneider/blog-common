package errors

var (
	// ErrInternalServerError — ошибка 500
	ErrInternalServerError = New(0, 500, "Internal Server Error", "", nil)

	// ErrUserNotFound — пользователь не найден
	ErrUserNotFound = New(1000, 404, "User not found", "", nil)

	// ErrInvalidCredentials — неверные учётные данные
	ErrInvalidCredentials = New(1001, 400, "Invalid credentials", "", nil)

	// ErrUserAlreadyExists — пользователь уже существует
	ErrUserAlreadyExists = New(1002, 400, "User already exists", "", nil)

	// ErrAccessTokenExpired — access token просрочен
	ErrAccessTokenExpired = New(1003, 401, "Access token expired", "", nil)

	// ErrRefreshTokenExpired — refresh token просрочен
	ErrRefreshTokenExpired = New(1004, 401, "Refresh token expired", "", nil)

	// ErrUnauthorized — доступ запрещён
	ErrUnauthorized = New(1005, 401, "Unauthorized", "", nil)

	// ErrValidationFailed — ошибка валидации
	ErrValidationFailed = New(1006, 400, "Validation failed", "", nil)

	// ErrUserDeleted - Пользователь был удален
	ErrUserDeleted = New(1007, 400, "User deleted", "", nil)
)
