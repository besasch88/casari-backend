package printer

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_env"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
Init the module by registering new APIs and PubSub consumers.
*/
func Init(envs *ceng_env.Envs, dbStorage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, routerGroup *gin.RouterGroup) {
	zap.L().Info("Initialize Printer package...")
	var repository printerRepositoryInterface
	var service printerServiceInterface
	var router printerRouterInterface

	repository = newPrinterRepository(envs.SearchRelevanceThreshold)
	service = newPrinterService(dbStorage, pubSubAgent, repository)
	router = newPrinterRouter(service)
	router.register(routerGroup)
	zap.L().Info("Printer package initialized")
}
