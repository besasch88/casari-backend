package statistics

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_router"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type statisticsRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type statisticsRouter struct {
	service statisticsServiceInterface
}

func newStatisticsRouter(service statisticsServiceInterface) statisticsRouter {
	return statisticsRouter{
		service: service,
	}
}

// Implementation
func (r statisticsRouter) register(router *gin.RouterGroup) {

	router.GET(
		"/statistics",
		ceng_auth.AuthMiddleware([]string{ceng_auth.READ_STATISTICS}),
		ceng_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			item, err := r.service.getStatistics(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "statistics-router"), zap.Error(err))
				ceng_router.ReturnGenericError(ctx)
				return
			}
			ceng_router.ReturnOk(ctx, &gin.H{"item": item})
		})

}
