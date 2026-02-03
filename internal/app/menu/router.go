package menu

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuRouter struct {
	service menuServiceInterface
}

func newMenuRouter(service menuServiceInterface) menuRouter {
	return menuRouter{
		service: service,
	}
}

// Implementation
func (r menuRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenu(ctx, request)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

}
