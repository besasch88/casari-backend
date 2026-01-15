package table

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type tableRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type tableRouter struct {
	service tableServiceInterface
}

func newTableRouter(service tableServiceInterface) tableRouter {
	return tableRouter{
		service: service,
	}
}

// Implementation
func (r tableRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/tables",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_OTHER_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request ListTablesInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			items, totalCount, err := r.service.listTables(ctx, request)

			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.GET(
		"/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_OTHER_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getTableByID(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_OTHER_TABLES, ceng_auth.WRITE_OTHER_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateTable(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errTableSameNameAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_OTHER_TABLES, ceng_auth.WRITE_OTHER_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteTable(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})

	router.GET(
		"/my/tables",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MY_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request ListTablesInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			authUser := ceng_auth.GetAuthenticatedUserFromSession(ctx)
			authUserId := uuid.MustParse(authUser.ID)
			request.UserId = &authUserId
			items, totalCount, err := r.service.listTables(ctx, request)

			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.GET(
		"/my/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MY_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			authUser := ceng_auth.GetAuthenticatedUserFromSession(ctx)
			authUserId := uuid.MustParse(authUser.ID)
			request.UserId = &authUserId
			item, err := r.service.getTableByID(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/my/tables",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MY_TABLES, ceng_auth.WRITE_MY_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			authUser := ceng_auth.GetAuthenticatedUserFromSession(ctx)
			authUserId := uuid.MustParse(authUser.ID)
			request.UserId = &authUserId
			item, err := r.service.createTable(ctx, request)
			if err == errTableSameNameAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/my/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MY_TABLES, ceng_auth.WRITE_MY_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			authUser := ceng_auth.GetAuthenticatedUserFromSession(ctx)
			authUserId := uuid.MustParse(authUser.ID)
			request.UserId = &authUserId
			item, err := r.service.updateTable(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errTableSameNameAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/my/tables/:tableId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_MY_TABLES, ceng_auth.WRITE_MY_TABLES}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteTableInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			authUser := ceng_auth.GetAuthenticatedUserFromSession(ctx)
			authUserId := uuid.MustParse(authUser.ID)
			request.UserId = &authUserId
			_, err := r.service.deleteTable(ctx, request)
			if err == errTableNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})
}
