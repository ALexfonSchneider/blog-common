package errors

var (
	// InternalServerError — ошибка 500
	InternalServerError = New(0, 500, "Internal Server Error", "", nil)

	// UserNotFound — пользователь не найден
	UserNotFound = New(1000, 404, "User not found", "", nil)

	// InvalidCredentials — неверные учётные данные
	InvalidCredentials = New(1001, 400, "Invalid credentials", "", nil)

	// UserAlreadyExists — пользователь уже существует
	UserAlreadyExists = New(1002, 400, "User already exists", "", nil)

	// AccessTokenExpired — access token просрочен
	AccessTokenExpired = New(1003, 401, "Access token expired", "", nil)

	// RefreshTokenExpired — refresh token просрочен
	RefreshTokenExpired = New(1004, 401, "Refresh token expired", "", nil)

	// Unauthorized — доступ запрещён
	Unauthorized = New(1005, 401, "Unauthorized", "", nil)

	// ValidationFailed — ошибка валидации
	ValidationFailed = New(1006, 400, "Validation failed", "", nil)
)
