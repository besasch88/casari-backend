package order

import (
	"slices"
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_auth"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_err"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderServiceInterface interface {
	createOrderFromEvent(event ceng_pubsub.TableEventEntity) error
	getOrder(ctx *gin.Context, input getOrderInputDto) (orderEntityWithChilds, error)
	updateOrder(ctx *gin.Context, input updateOrderInputDto) (orderEntityWithChilds, error)
}

type orderService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  orderRepositoryInterface
}

func newOrderService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository orderRepositoryInterface) orderService {
	return orderService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s orderService) createOrderFromEvent(event ceng_pubsub.TableEventEntity) error {
	now := time.Now()
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		orderID := uuid.New()
		newOrder := orderEntity{
			ID:        orderID,
			TableID:   event.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.repository.saveOrder(s.storage, newOrder, ceng_db.Create)
		// Send an event of order created
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicOrderV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.OrderCreatedEvent,
				EventEntity: &ceng_pubsub.OrderEventEntity{
					ID:        newOrder.ID,
					TableID:   newOrder.TableID,
					CreatedAt: newOrder.CreatedAt,
					UpdatedAt: newOrder.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(orderEntity{}, newOrder),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		newCourse := courseEntity{
			ID:        uuid.New(),
			OrderID:   orderID,
			Position:  1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.repository.saveCourse(s.storage, newCourse, ceng_db.Create)
		// Send an event of course created
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicCourseV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.CourseCreatedEvent,
				EventEntity: &ceng_pubsub.CourseEventEntity{
					ID:        newCourse.ID,
					OrderID:   newCourse.OrderID,
					Position:  newCourse.Position,
					CreatedAt: newCourse.CreatedAt,
					UpdatedAt: newCourse.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(courseEntity{}, newCourse),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return nil
}

func (s orderService) getOrder(ctx *gin.Context, input getOrderInputDto) (orderEntityWithChilds, error) {
	tableID := uuid.MustParse(input.TableID)
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
		userID = nil
	}
	// Check if the table exists for that user
	if item, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
		return orderEntityWithChilds{}, ceng_err.ErrGeneric
	} else if ceng_utils.IsEmpty(item) {
		return orderEntityWithChilds{}, errTableNotFound
	}
	// Retrieve the order
	order, err := s.repository.getOrderByTableID(s.storage, tableID, false)
	if err != nil {
		return orderEntityWithChilds{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(order) {
		return orderEntityWithChilds{}, errOrderNotFound
	}
	// Retrieve all courses
	courses, _, err := s.repository.listCoursesByOrderID(s.storage, order.ID, false)
	if err != nil {
		return orderEntityWithChilds{}, ceng_err.ErrGeneric
	}
	// For each course, retrieve its selections
	coursesWithChilds := []courseEntityWithChilds{}
	for _, course := range courses {
		selections, _, err := s.repository.listCourseSelectionsByCourseID(s.storage, course.ID, false)
		if err != nil {
			return orderEntityWithChilds{}, ceng_err.ErrGeneric
		}
		coursesWithChilds = append(coursesWithChilds, courseEntityWithChilds{
			courseEntity: course,
			Items:        selections,
		})
	}
	return orderEntityWithChilds{
		orderEntity: order,
		Courses:     coursesWithChilds,
	}, nil
}

func (s orderService) updateOrder(ctx *gin.Context, input updateOrderInputDto) (orderEntityWithChilds, error) {

	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		tableID := uuid.MustParse(input.TableID)
		requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
		userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
		// Check if the user has a permission
		if slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
			userID = nil
		}
		// Check if the table exists for that user
		if item, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
			return ceng_err.ErrGeneric
		} else if ceng_utils.IsEmpty(item) {
			return errTableNotFound
		}
		// Check if the Order exists
		order, err := s.repository.getOrderByTableID(s.storage, tableID, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		// If it does not exists, return an error
		if ceng_utils.IsEmpty(order) {
			return errOrderNotFound
		}

		// Get existing courses
		courses, total, err := s.repository.listCoursesByOrderID(s.storage, order.ID, false)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		// Now for each course received from input, check if already exists or needed to be created a new one
		var lastPosition int64 = total
		for index, inputCourse := range input.Courses {
			var course courseEntity
			if len(courses) >= index+1 {
				course = courses[index]
				// Check if the order of items are respected
				if course.ID.String() != inputCourse.ID {
					return errCourseMismatch
				}
			} else {
				// If there are new course input, let's create them
				lastPosition++
				course = courseEntity{
					ID:        ceng_utils.GetUUIDFromString(inputCourse.ID),
					OrderID:   order.ID,
					Position:  lastPosition,
					CreatedAt: now,
					UpdatedAt: now,
				}
				if _, err := s.repository.saveCourse(s.storage, course, ceng_db.Create); err != nil {
					return ceng_err.ErrGeneric
				}
				// Send an event of order created
				if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicCourseV1, ceng_pubsub.PubSubMessage{
					Message: ceng_pubsub.PubSubEvent{
						EventID:   uuid.New(),
						EventTime: time.Now(),
						EventType: ceng_pubsub.OrderCreatedEvent,
						EventEntity: &ceng_pubsub.OrderEventEntity{
							ID:        order.ID,
							TableID:   order.TableID,
							CreatedAt: order.CreatedAt,
							UpdatedAt: order.UpdatedAt,
						},
						EventChangedFields: ceng_utils.DiffStructs(courseEntity{}, course),
					},
				}); err != nil {
					return ceng_err.ErrGeneric
				} else {
					eventsToPublish = append(eventsToPublish, event)
				}
			}
			if err := s.repository.deleteSelectionsByCourseID(s.storage, course.ID); err != nil {
				return ceng_err.ErrGeneric
			}
			for _, inputSelection := range inputCourse.Items {
				selection := courseSelectionEntity{
					ID:           uuid.New(),
					CourseID:     course.ID,
					MenuItemID:   ceng_utils.GetUUIDFromString(inputSelection.MenuItemID),
					MenuOptionID: ceng_utils.GetOptionalUUIDFromString(inputSelection.MenuOptionID),
					Quantity:     inputSelection.Quantity,
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				if _, err := s.repository.saveSelection(s.storage, selection, ceng_db.Create); err != nil {
					return ceng_err.ErrGeneric
				}
			}
		}
		return nil
	})
	if errTransaction != nil {
		return orderEntityWithChilds{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return s.getOrder(ctx, getOrderInputDto{TableID: input.TableID})
}
