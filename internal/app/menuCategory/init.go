package menuCategory

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
	zap.L().Info("Initialize MenuCategory package...")
	var repository menuCategoryRepositoryInterface
	var service menuCategoryServiceInterface
	var router menuCategoryRouterInterface

	repository = newMenuCategoryRepository(envs.SearchRelevanceThreshold)
	service = newMenuCategoryService(dbStorage, pubSubAgent, repository)
	router = newMenuCategoryRouter(service)
	router.register(routerGroup)
	zap.L().Info("MenuCategory package initialized")
}
