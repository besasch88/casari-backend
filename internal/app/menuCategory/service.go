package menuCategory

import (
	"math"
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_err"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type menuCategoryServiceInterface interface {
	listMenuCategories(ctx *gin.Context) ([]menuCategoryEntity, int64, error)
	getMenuCategoryByID(ctx *gin.Context, input getMenuCategoryInputDto) (menuCategoryEntity, error)
	createMenuCategory(ctx *gin.Context, input createMenuCategoryInputDto) (menuCategoryEntity, error)
	updateMenuCategory(ctx *gin.Context, input updateMenuCategoryInputDto) (menuCategoryEntity, error)
	deleteMenuCategory(ctx *gin.Context, input deleteMenuCategoryInputDto) (menuCategoryEntity, error)
}

type menuCategoryService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  menuCategoryRepositoryInterface
}

func newMenuCategoryService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository menuCategoryRepositoryInterface) menuCategoryService {
	return menuCategoryService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuCategoryService) listMenuCategories(ctx *gin.Context) ([]menuCategoryEntity, int64, error) {
	items, totalCount, err := s.repository.listMenuCategories(s.storage, false)
	if err != nil || items == nil {
		return []menuCategoryEntity{}, 0, ceng_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s menuCategoryService) getMenuCategoryByID(ctx *gin.Context, input getMenuCategoryInputDto) (menuCategoryEntity, error) {
	menuCategoryID := uuid.MustParse(input.ID)
	item, err := s.repository.getMenuCategoryByID(s.storage, menuCategoryID, false)
	if err != nil {
		return menuCategoryEntity{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(item) {
		return menuCategoryEntity{}, errMenuCategoryNotFound
	}
	return item, nil
}

func (s menuCategoryService) createMenuCategory(ctx *gin.Context, input createMenuCategoryInputDto) (menuCategoryEntity, error) {
	now := time.Now()
	maxValue := int64(math.MaxInt64)
	newMenuCategory := menuCategoryEntity{
		ID:        uuid.New(),
		Title:     input.Title,
		Position:  maxValue,
		Active:    ceng_utils.BoolPtr(false),
		CreatedAt: now,
		UpdatedAt: now,
	}
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuCategoryEntity
		if item, err := s.repository.getMenuCategoryByTitle(tx, input.Title, false); err != nil {
			return ceng_err.ErrGeneric
		} else if !ceng_utils.IsEmpty(item) {
			return errMenuCategorySameTitleAlreadyExists
		} else if _, err = s.repository.saveMenuCategory(tx, newMenuCategory, ceng_db.Create); err != nil {
			return ceng_err.ErrGeneric
		} else if updatedEntities, err = s.repository.recalculateMenuCategorysPosition(tx); err != nil {
			return ceng_err.ErrGeneric
		} else if newMenuCategory, err = s.repository.getMenuCategoryByID(tx, newMenuCategory.ID, false); err != nil {
			return ceng_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuCategoryCreatedEvent,
				EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
					ID:        newMenuCategory.ID,
					Title:     newMenuCategory.Title,
					Position:  newMenuCategory.Position,
					Active:    newMenuCategory.Active,
					CreatedAt: newMenuCategory.CreatedAt,
					UpdatedAt: newMenuCategory.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(menuCategoryEntity{}, newMenuCategory),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == newMenuCategory.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
				Message: ceng_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: ceng_pubsub.MenuCategoryUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
						ID:        updatedEntity.ID,
						Title:     updatedEntity.Title,
						Position:  updatedEntity.Position,
						Active:    updatedEntity.Active,
						CreatedAt: updatedEntity.CreatedAt,
						UpdatedAt: updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuCategoryEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newMenuCategory, nil
}

func (s menuCategoryService) updateMenuCategory(ctx *gin.Context, input updateMenuCategoryInputDto) (menuCategoryEntity, error) {
	now := time.Now()
	var updatedMenuCategory menuCategoryEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuCategoryEntity
		menuCategoryId := uuid.MustParse(input.ID)
		currentMenuCategory, err := s.repository.getMenuCategoryByID(tx, menuCategoryId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentMenuCategory) {
			return errMenuCategoryNotFound
		}
		updatedMenuCategory = currentMenuCategory
		updatedMenuCategory.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			menuCategorySameTitle, err := s.repository.getMenuCategoryByTitle(tx, *input.Title, false)
			if err != nil {
				return ceng_err.ErrGeneric
			}
			if !ceng_utils.IsEmpty(menuCategorySameTitle) && menuCategorySameTitle.ID.String() != input.ID {
				return errMenuCategorySameTitleAlreadyExists
			}
			updatedMenuCategory.Title = *input.Title
		}
		if input.Active != nil {
			updatedMenuCategory.Active = input.Active
		}
		if input.Position != nil {
			// If the step is moving in a lower position (e.g. from 10 to 3),
			// we need to move it one step more, so that, the algorith to re-sort all steps correctly
			if updatedMenuCategory.Position < *input.Position {
				*input.Position++
			}
			updatedMenuCategory.Position = *input.Position
		}
		if _, err = s.repository.saveMenuCategory(tx, updatedMenuCategory, ceng_db.Update); err != nil {
			return ceng_err.ErrGeneric
		}
		if updatedEntities, err = s.repository.recalculateMenuCategorysPosition(tx); err != nil {
			return ceng_err.ErrGeneric
		}
		if updatedMenuCategory, err = s.repository.getMenuCategoryByID(tx, updatedMenuCategory.ID, false); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of menuCategory updated
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuCategoryUpdatedEvent,
				EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
					ID:        updatedMenuCategory.ID,
					Title:     updatedMenuCategory.Title,
					Position:  updatedMenuCategory.Position,
					Active:    updatedMenuCategory.Active,
					CreatedAt: updatedMenuCategory.CreatedAt,
					UpdatedAt: updatedMenuCategory.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentMenuCategory, updatedMenuCategory),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == updatedMenuCategory.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
				Message: ceng_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: ceng_pubsub.MenuCategoryUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
						ID:        updatedEntity.ID,
						Title:     updatedEntity.Title,
						Position:  updatedEntity.Position,
						Active:    updatedEntity.Active,
						CreatedAt: updatedEntity.CreatedAt,
						UpdatedAt: updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuCategoryEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedMenuCategory, nil
}

func (s menuCategoryService) deleteMenuCategory(ctx *gin.Context, input deleteMenuCategoryInputDto) (menuCategoryEntity, error) {
	var currentMenuCategory menuCategoryEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuCategoryEntity
		// Check if exists
		menuCategoryId := uuid.MustParse(input.ID)
		currentMenuCategory, err := s.repository.getMenuCategoryByID(tx, menuCategoryId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentMenuCategory) {
			return errMenuCategoryNotFound
		}
		s.repository.deleteMenuCategory(tx, currentMenuCategory)

		if updatedEntities, err = s.repository.recalculateMenuCategorysPosition(tx); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of menuCategory deleted
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuCategoryDeletedEvent,
				EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
					ID:        currentMenuCategory.ID,
					Title:     currentMenuCategory.Title,
					Position:  currentMenuCategory.Position,
					Active:    currentMenuCategory.Active,
					CreatedAt: currentMenuCategory.CreatedAt,
					UpdatedAt: currentMenuCategory.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentMenuCategory, menuCategoryEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuCategoryV1, ceng_pubsub.PubSubMessage{
				Message: ceng_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: ceng_pubsub.MenuCategoryUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuCategoryEventEntity{
						ID:        updatedEntity.ID,
						Title:     updatedEntity.Title,
						Position:  updatedEntity.Position,
						Active:    updatedEntity.Active,
						CreatedAt: updatedEntity.CreatedAt,
						UpdatedAt: updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuCategoryEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentMenuCategory, nil
}
