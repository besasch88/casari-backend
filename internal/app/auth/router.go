package auth

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type authRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type authRouter struct {
	service authServiceInterface
}

func newAuthRouter(service authServiceInterface) authRouter {
	return authRouter{
		service: service,
	}
}

// Implementation
func (r authRouter) register(router *gin.RouterGroup) {
	router.POST(
		"/auth/login",
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input
			var request loginUserInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.login(ctx, request.Username, request.Password)
			if err == errInvalidCredentials {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/auth/refresh",
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			var request refreshTokenInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.refreshToken(ctx, request.RefreshToken)
			if err == errExpiredRefreshToken {
				ceng_router.ReturnUnauthorizedError(ctx)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/auth/logout",
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input
			var request refreshTokenInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			err := r.service.revokeRefreshToken(ctx, request.RefreshToken)
			if err == errExpiredRefreshToken {
				ceng_router.ReturnUnauthorizedError(ctx)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})
}
