package course

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
	zap.L().Info("Initialize Course package...")
	var repository courseRepositoryInterface
	var service courseServiceInterface
	var router courseRouterInterface

	repository = newCourseRepository(envs.SearchRelevanceThreshold)
	service = newCourseService(dbStorage, pubSubAgent, repository)
	router = newCourseRouter(service)
	router.register(routerGroup)
	zap.L().Info("Course package initialized")
}
