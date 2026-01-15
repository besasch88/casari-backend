package auth

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_env"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_scheduler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
Init the module by registering new APIs
*/
func Init(envs *ceng_env.Envs, dbStorage *gorm.DB, cron *ceng_scheduler.Scheduler, routerGroup *gin.RouterGroup) {
	zap.L().Info("Initialize Auth package...")
	var repository authRepositoryInterface
	var userRepository authUserRepositoryInterface
	var util authUtilInterface
	var service authServiceInterface
	var scheduler authSchedulerInterface
	var router authRouterInterface

	repository = newAuthRepository()
	userRepository = newAuthUserRepository()
	util = newAuthUtil(envs.AuthJwtSecret, envs.AuthJwtAccessTokenDuration, envs.AuthJwtRefreshTokenDuration)
	service = newAuthService(dbStorage, repository, userRepository, util)
	scheduler = newAuthScheduler(dbStorage, cron, repository)
	scheduler.init()

	router = newAuthRouter(service)
	router.register(routerGroup)
	zap.L().Info("Auth package initialized")
}
