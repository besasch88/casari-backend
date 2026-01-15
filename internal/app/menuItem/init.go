package menuItem

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
	zap.L().Info("Initialize MenuItem package...")
	var repository menuItemRepositoryInterface
	var service menuItemServiceInterface
	var router menuItemRouterInterface

	repository = newMenuItemRepository(envs.SearchRelevanceThreshold)
	service = newMenuItemService(dbStorage, pubSubAgent, repository)
	router = newMenuItemRouter(service)
	router.register(routerGroup)
	zap.L().Info("MenuItem package initialized")
}
