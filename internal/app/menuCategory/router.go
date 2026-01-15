package menuCategory

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuCategoryRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuCategoryRouter struct {
	service menuCategoryServiceInterface
}

func newMenuCategoryRouter(service menuCategoryServiceInterface) menuCategoryRouter {
	return menuCategoryRouter{
		service: service,
	}
}

// Implementation
func (r menuCategoryRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu/categories",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			items, totalCount, err := r.service.listMenuCategories(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/menu/categories",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createMenuCategoryInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createMenuCategory(ctx, request)
			if err == errMenuCategorySameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/menu/categories/:menuCategoryId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuCategoryInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenuCategoryByID(ctx, request)
			if err == errMenuCategoryNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/menu/categories/:menuCategoryId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateMenuCategoryInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateMenuCategory(ctx, request)
			if err == errMenuCategoryNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuCategorySameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/menu/categories/:menuCategoryId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MENU, ceng_auth.WRITE_MENU}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteMenuCategoryInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteMenuCategory(ctx, request)
			if err == errMenuCategoryNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})
}
