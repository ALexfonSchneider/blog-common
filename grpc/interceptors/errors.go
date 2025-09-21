package interceptors

import (
	"context"
	"fmt"
	apperror "github.com/ALexfonSchneider/blog-common/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
	"log/slog"
	"strconv"
)

type ErrorsInterceptor struct {
	appCodeToGrpcStatus map[int]codes.Code
	domain              string
	logger              *slog.Logger
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
		appCodeToGrpcStatus: config.AppCodeToGrpcStatusMappings,
		domain:              config.Domain,
		logger:              logger,
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

		appErr, ok := err.(*apperror.ApplicationError)
		if !ok {
			appErr = apperror.New(5000, "Internal server error", "").Wrap(err)
		}

		var (
			errCode    = appErr.Code()
			errMessage = appErr.Message()
			errDetail  = appErr.Detail()
			errCause   = appErr.Cause()

			errCodeStr = strconv.Itoa(errCode)
			grpcCode   = i.getGRPCCode(errCode)
		)

		i.logger.Error("gRPC handler error",
			"method", info.FullMethod,
			"app_code", errCode,
			"grpc_code", grpcCode,
			"message", errMessage,
			"detail", errDetail,
			"cause", errCause,
		)

		md := metadata.Pairs(
			"app-code", errCodeStr,
			"app-message", errMessage,
			"app-detail", errDetail,
		)
		if err = grpc.SendHeader(ctx, md); err != nil {
			return nil, err
		}

		st := status.New(grpcCode, errMessage)

		details := []protoadapt.MessageV1{
			&errdetails.ErrorInfo{
				Reason: fmt.Sprintf("APP_ERROR_%s", errCodeStr),
				Domain: i.domain,
				Metadata: map[string]string{
					"app-code":    errCodeStr,
					"app-message": errMessage,
				},
			},
		}

		if errDetail != "" || errCause != nil {
			debug := &errdetails.DebugInfo{
				Detail: appErr.Detail(),
			}
			if errCause != nil {
				debug.StackEntries = []string{errCause.Error()}
			}

			details = append(details, debug)
		}

		stWithDetails, sErr := st.WithDetails(details...)
		if sErr != nil {
			return nil, st.Err()
		}

		return nil, stWithDetails.Err()
	}
}

func (i *ErrorsInterceptor) getGRPCCode(appCode int) codes.Code {
	if code, ok := i.appCodeToGrpcStatus[appCode]; ok {
		return code
	}
	return codes.Unknown
}

func first(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func ExtractAppError(ctx context.Context, err error) *apperror.ApplicationError {
	st, ok := status.FromError(err)
	if !ok {
		return apperror.New(5000, err.Error(), "")
	}

	// --- 1. Попытка достать metadata ---
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		code, _ := strconv.Atoi(first(md["app-code"]))
		message := first(md["app-message"])
		detail := first(md["app-detail"])
		return apperror.New(code, message, detail)
	}

	var (
		code    int
		message string
		detail  string
	)

	// --- 2. Попытка достать из status details ---
	for _, d := range st.Details() {
		switch info := d.(type) {
		case *errdetails.ErrorInfo:
			code, _ = strconv.Atoi(info.Metadata["app-code"])
			message = st.Message()
		case *errdetails.DebugInfo:
			detail = info.Detail
		}
	}

	return apperror.New(code, message, detail)
}
