package course

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

type courseServiceInterface interface {
	listCourses(ctx *gin.Context, input ListCoursesInputDto) ([]courseEntity, int64, error)
	getCourseByID(ctx *gin.Context, input getCourseInputDto) (courseEntity, error)
	createCourse(ctx *gin.Context, input createCourseInputDto) (courseEntity, error)
	updateCourse(ctx *gin.Context, input updateCourseInputDto) (courseEntity, error)
	deleteCourse(ctx *gin.Context, input deleteCourseInputDto) (courseEntity, error)
}

type courseService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  courseRepositoryInterface
}

func newCourseService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository courseRepositoryInterface) courseService {
	return courseService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s courseService) listCourses(ctx *gin.Context, input ListCoursesInputDto) ([]courseEntity, int64, error) {
	tableID := uuid.MustParse(input.TableId)
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
		userID = nil
	}
	if exists, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
		return []courseEntity{}, 0, ceng_err.ErrGeneric
	} else if !exists {
		return []courseEntity{}, 0, errTableNotFound
	}
	items, totalCount, err := s.repository.listCourses(s.storage, tableID, false)
	if err != nil || items == nil {
		return []courseEntity{}, 0, ceng_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s courseService) getCourseByID(ctx *gin.Context, input getCourseInputDto) (courseEntity, error) {
	courseID := uuid.MustParse(input.ID)
	item, err := s.repository.getCourseByID(s.storage, courseID, false)
	if err != nil {
		return courseEntity{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(item) {
		return courseEntity{}, errCourseNotFound
	}
	// Check if the course is part of a table available by the user
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	// Skip check if the user can access other tables
	if !slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
		userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
		if exists, err := s.repository.checkTableExists(s.storage, userID, item.TableID); err != nil {
			return courseEntity{}, ceng_err.ErrGeneric
		} else if !exists {
			return courseEntity{}, errCourseNotFound
		}
	}
	return item, nil
}

func (s courseService) createCourse(ctx *gin.Context, input createCourseInputDto) (courseEntity, error) {
	tableID := uuid.MustParse(input.TableId)
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.WRITE_OTHER_TABLES) {
		userID = nil
	}
	if exists, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
		return courseEntity{}, ceng_err.ErrGeneric
	} else if !exists {
		return courseEntity{}, errTableNotFound
	}
	now := time.Now()
	newCourse := courseEntity{
		ID:        uuid.New(),
		TableID:   tableID,
		UserID:    uuid.MustParse(requester.ID),
		Close:     ceng_utils.BoolPtr(false),
		CreatedAt: now,
		UpdatedAt: now,
	}
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		if _, err := s.repository.saveCourse(tx, newCourse, ceng_db.Create); err != nil {
			return ceng_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicCourseV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.CourseCreatedEvent,
				EventEntity: &ceng_pubsub.CourseEventEntity{
					ID:        newCourse.ID,
					UserID:    newCourse.UserID,
					TableID:   newCourse.TableID,
					Close:     newCourse.Close,
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
		return courseEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newCourse, nil
}

func (s courseService) updateCourse(ctx *gin.Context, input updateCourseInputDto) (courseEntity, error) {
	now := time.Now()
	var updatedCourse courseEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		courseId := uuid.MustParse(input.ID)
		currentCourse, err := s.repository.getCourseByID(tx, courseId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentCourse) {
			return errCourseNotFound
		}
		// Check if the course is part of a table available by the user
		requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
		// Skip check if the user can access other tables
		if !slices.Contains(requester.Permissions, ceng_auth.WRITE_OTHER_TABLES) {
			userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
			if exists, err := s.repository.checkTableExists(s.storage, userID, currentCourse.TableID); err != nil {
				return ceng_err.ErrGeneric
			} else if !exists {
				return errCourseNotFound
			}
		}
		updatedCourse = currentCourse
		updatedCourse.UpdatedAt = now
		if input.Close != nil {
			updatedCourse.Close = input.Close
		}
		if _, err = s.repository.saveCourse(tx, updatedCourse, ceng_db.Update); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of course updated
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicCourseV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.CourseUpdatedEvent,
				EventEntity: &ceng_pubsub.CourseEventEntity{
					ID:        updatedCourse.ID,
					TableID:   updatedCourse.TableID,
					UserID:    updatedCourse.UserID,
					Close:     updatedCourse.Close,
					CreatedAt: updatedCourse.CreatedAt,
					UpdatedAt: updatedCourse.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentCourse, updatedCourse),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return courseEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedCourse, nil
}

func (s courseService) deleteCourse(ctx *gin.Context, input deleteCourseInputDto) (courseEntity, error) {
	var currentCourse courseEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		// Check if the Menu Item exists
		courseId := uuid.MustParse(input.ID)
		currentCourse, err := s.repository.getCourseByID(tx, courseId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentCourse) {
			return errCourseNotFound
		}
		// Check if the course is part of a table available by the user
		requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
		// Skip check if the user can access other tables
		if !slices.Contains(requester.Permissions, ceng_auth.WRITE_OTHER_TABLES) {
			userID := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
			if exists, err := s.repository.checkTableExists(s.storage, userID, currentCourse.TableID); err != nil {
				return ceng_err.ErrGeneric
			} else if !exists {
				return errCourseNotFound
			}
		}
		s.repository.deleteCourse(tx, currentCourse)

		// Send an event of course deleted
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicCourseV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.CourseDeletedEvent,
				EventEntity: &ceng_pubsub.CourseEventEntity{
					ID:        currentCourse.ID,
					TableID:   currentCourse.TableID,
					UserID:    currentCourse.UserID,
					Close:     currentCourse.Close,
					CreatedAt: currentCourse.CreatedAt,
					UpdatedAt: currentCourse.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentCourse, courseEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return courseEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentCourse, nil
}
