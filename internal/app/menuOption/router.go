package menuOption

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuOptionRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuOptionRouter struct {
	service menuOptionServiceInterface
}

func newMenuOptionRouter(service menuOptionServiceInterface) menuOptionRouter {
	return menuOptionRouter{
		service: service,
	}
}

// Implementation
func (r menuOptionRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu/items/:menuItemId/options",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request listMenuOptionsInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			items, totalCount, err := r.service.listMenuOptions(ctx, request)
			if err == errMenuItemNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/menu/items/:menuItemId/options",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createMenuOptionInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createMenuOption(ctx, request)
			if err == errMenuItemNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuOptionSameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/menu/options/:menuOptionId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuOptionInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenuOptionByID(ctx, request)
			if err == errMenuOptionNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/menu/options/:menuOptionId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateMenuOptionInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateMenuOption(ctx, request)
			if err == errMenuOptionNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuOptionSameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/menu/options/:menuOptionId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteMenuOptionInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteMenuOption(ctx, request)
			if err == errMenuOptionNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})
}
