package statistics

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_err"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type statisticsServiceInterface interface {
	getStatistics(ctx *gin.Context) (statisticsEntity, error)
}

type statisticsService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  statisticsRepositoryInterface
}

func newStatisticsService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository statisticsRepositoryInterface) statisticsService {
	return statisticsService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s statisticsService) getStatistics(ctx *gin.Context) (statisticsEntity, error) {
	avgTableDuration, err := s.repository.getAverageTableDuration(s.storage)
	if err != nil {
		return statisticsEntity{}, ceng_err.ErrGeneric
	}
	paymentMethodsTakins, err := s.repository.getPaymentMethodsTakins(s.storage)
	if err != nil {
		return statisticsEntity{}, ceng_err.ErrGeneric
	}
	menuItemStats, err := s.repository.getMenuItemStats(s.storage)
	if err != nil {
		return statisticsEntity{}, ceng_err.ErrGeneric
	}
	return statisticsEntity{
		AvgTableDuration: avgTableDuration,
		PaymentsTakins:   paymentMethodsTakins,
		MenuItemStats:    menuItemStats,
	}, nil
}
