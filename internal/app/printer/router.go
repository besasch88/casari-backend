package printer

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type printerRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type printerRouter struct {
	service printerServiceInterface
}

func newPrinterRouter(service printerServiceInterface) printerRouter {
	return printerRouter{
		service: service,
	}
}

// Implementation
func (r printerRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/printers",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_PRINTER}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			items, totalCount, err := r.service.listPrinters(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/printers",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_PRINTER, ceng_auth.WRITE_PRINTER}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createPrinterInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createPrinter(ctx, request)
			if err == errPrinterSameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/printers/:printerId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_PRINTER}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getPrinterInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getPrinterByID(ctx, request)
			if err == errPrinterNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/printers/:printerId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_PRINTER, ceng_auth.WRITE_PRINTER}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updatePrinterInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updatePrinter(ctx, request)
			if err == errPrinterNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errPrinterSameTitleAlreadyExists {
				ceng_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/printers/:printerId",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_PRINTER, ceng_auth.WRITE_PRINTER}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deletePrinterInputDto
			if err := ceng_router.BindParameters(ctx, &request); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				ceng_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deletePrinter(ctx, request)
			if err == errPrinterNotFound {
				ceng_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnNoContent(ctx)
		})
}
