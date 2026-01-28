package order

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type orderRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type orderRouter struct {
	service      orderServiceInterface
	printService printServiceInterface
}

func newOrderRouter(service orderServiceInterface, printService printServiceInterface) orderRouter {
	return orderRouter{
		service:      service,
		printService: printService,
	}
}

// Implementation
func (r orderRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/tables/:tableId/order",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getOrderInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getOrder(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errOrderNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "order-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/tables/:tableId/order",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateOrderInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateOrder(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errOrderNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errCourseMismatch {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}

			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "order-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/tables/:tableId/order/print",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request printOrderInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"success": true})
		})

}
