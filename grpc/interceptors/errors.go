package interceptors

import (
	"context"
	"errors"
	"fmt"
	apperror "github.com/ALexfonSchneider/blog-common/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"strconv"
)

type ErrorsInterceptor struct {
	appCodeToGrpcStatusMappings map[int]codes.Code
	domain                      string
	logger                      *slog.Logger
}

type Config struct {
	Domain                      string
	AppCodeToGrpcStatusMappings map[int]codes.Code
	Logger                      *slog.Logger
}

func NewErrorsInterceptor(config Config) *ErrorsInterceptor {
	logger := config.Logger
	if logger == nil {
		logger = slog.Default()
	}

	interceptor := &ErrorsInterceptor{
		appCodeToGrpcStatusMappings: config.AppCodeToGrpcStatusMappings,
		domain:                      config.Domain,
		logger:                      logger,
	}

	return interceptor
}

func (i *ErrorsInterceptor) Interceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var appErr *apperror.ApplicationError
		if errors.As(err, &appErr) {
			return nil, i.buildStatusWithDetails(appErr, info.FullMethod)
		}

		unknownErr := apperror.InternalServerError.Wrap(err)

		return nil, i.buildStatusWithDetails(unknownErr, info.FullMethod)
	}
}

func (i *ErrorsInterceptor) buildStatusWithDetails(
	appErr *apperror.ApplicationError,
	method string,
) error {
	// Определяем gRPC код
	grpcCode, ok := i.appCodeToGrpcStatusMappings[appErr.Code()]
	if !ok {
		grpcCode = httpToGRPCCode(appErr.HttpCode())
	}

	// Формируем сообщение
	msg := appErr.Message()
	if appErr.Detail() != "" {
		msg = fmt.Sprintf("%s: %s", appErr.Message(), appErr.Detail())
	}

	// Создаём статус
	st := status.New(grpcCode, msg)

	// Добавляем ErrorInfo
	withInfo, err := st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: strconv.Itoa(appErr.Code()),
			Domain: i.domain,
		},
	)

	if err != nil {
		// Если не удалось добавить детали — возвращаем хотя бы статус
		return st.Err()
	}

	// Логируем
	i.logger.Error("gRPC handler error",
		"method", method,
		"message", msg,
		"cause", appErr.Cause(),
		"app_code", appErr.Code(),
		"grpc_code", grpcCode,
	)

	return withInfo.Err()
}

func httpToGRPCCode(httpCode int) codes.Code {
	switch httpCode {
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 409:
		return codes.AlreadyExists
	case 429:
		return codes.ResourceExhausted
	case 500:
		return codes.Internal
	case 503:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
