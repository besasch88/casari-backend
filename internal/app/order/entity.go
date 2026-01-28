package order

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
)

type orderEntity ceng_pubsub.OrderEventEntity
type courseEntity ceng_pubsub.CourseEventEntity
type courseSelectionEntity ceng_pubsub.CourseSelectionEventEntity

type courseEntityWithChilds struct {
	courseEntity
	Items []courseSelectionEntity `json:"items"`
}

type orderEntityWithChilds struct {
	orderEntity
	Courses []courseEntityWithChilds `json:"courses"`
}
