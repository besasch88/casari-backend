package menuItem

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

type menuItemServiceInterface interface {
	listMenuItems(ctx *gin.Context, input listMenuItemsInputDto) ([]menuItemEntity, int64, error)
	getMenuItemByID(ctx *gin.Context, input getMenuItemInputDto) (menuItemEntity, error)
	createMenuItem(ctx *gin.Context, input createMenuItemInputDto) (menuItemEntity, error)
	updateMenuItem(ctx *gin.Context, input updateMenuItemInputDto) (menuItemEntity, error)
	deleteMenuItem(ctx *gin.Context, input deleteMenuItemInputDto) (menuItemEntity, error)
}

type menuItemService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  menuItemRepositoryInterface
}

func newMenuItemService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository menuItemRepositoryInterface) menuItemService {
	return menuItemService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuItemService) listMenuItems(ctx *gin.Context, input listMenuItemsInputDto) ([]menuItemEntity, int64, error) {
	menuCategoryID := uuid.MustParse(input.MenuCategoryId)
	if exists, err := s.repository.checkMenuCategoryExists(s.storage, menuCategoryID); err != nil {
		return []menuItemEntity{}, 0, ceng_err.ErrGeneric
	} else if !exists {
		return []menuItemEntity{}, 0, errMenuCategoryNotFound
	}
	items, totalCount, err := s.repository.listMenuItems(s.storage, menuCategoryID, false)
	if err != nil || items == nil {
		return []menuItemEntity{}, 0, ceng_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s menuItemService) getMenuItemByID(ctx *gin.Context, input getMenuItemInputDto) (menuItemEntity, error) {
	menuItemID := uuid.MustParse(input.ID)
	item, err := s.repository.getMenuItemByID(s.storage, menuItemID, false)
	if err != nil {
		return menuItemEntity{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(item) {
		return menuItemEntity{}, errMenuItemNotFound
	}
	return item, nil
}

func (s menuItemService) createMenuItem(ctx *gin.Context, input createMenuItemInputDto) (menuItemEntity, error) {
	menuCategoryId := uuid.MustParse(input.MenuCategoryId)
	menuCategoryID := uuid.MustParse(input.MenuCategoryId)
	if exists, err := s.repository.checkMenuCategoryExists(s.storage, menuCategoryID); err != nil {
		return menuItemEntity{}, ceng_err.ErrGeneric
	} else if !exists {
		return menuItemEntity{}, errMenuCategoryNotFound
	}
	now := time.Now()
	maxValue := int64(math.MaxInt64)
	newMenuItem := menuItemEntity{
		ID:             uuid.New(),
		MenuCategoryID: menuCategoryId,
		Position:       maxValue,
		Title:          input.Title,
		Active:         ceng_utils.BoolPtr(false),
		Price:          input.Price,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		if item, err := s.repository.getMenuItemByTitle(tx, input.Title, false); err != nil {
			return ceng_err.ErrGeneric
		} else if !ceng_utils.IsEmpty(item) {
			return errMenuItemSameTitleAlreadyExists
		} else if _, err = s.repository.saveMenuItem(tx, newMenuItem, ceng_db.Create); err != nil {
			return ceng_err.ErrGeneric
		} else if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, menuCategoryId); err != nil {
			return ceng_err.ErrGeneric
		} else if newMenuItem, err = s.repository.getMenuItemByID(tx, newMenuItem.ID, false); err != nil {
			return ceng_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuItemV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuItemCreatedEvent,
				EventEntity: &ceng_pubsub.MenuItemEventEntity{
					ID:        newMenuItem.ID,
					Title:     newMenuItem.Title,
					Position:  newMenuItem.Position,
					Active:    newMenuItem.Active,
					Price:     newMenuItem.Price,
					CreatedAt: newMenuItem.CreatedAt,
					UpdatedAt: newMenuItem.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(menuItemEntity{}, newMenuItem),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == newMenuItem.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuItemV1, ceng_pubsub.PubSubMessage{
				Message: ceng_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: ceng_pubsub.MenuItemUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuItemEventEntity{
						ID:             updatedEntity.ID,
						MenuCategoryID: updatedEntity.MenuCategoryID,
						Title:          updatedEntity.Title,
						Position:       updatedEntity.Position,
						Active:         updatedEntity.Active,
						Price:          updatedEntity.Price,
						CreatedAt:      updatedEntity.CreatedAt,
						UpdatedAt:      updatedEntity.UpdatedAt,
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
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newMenuItem, nil
}

func (s menuItemService) updateMenuItem(ctx *gin.Context, input updateMenuItemInputDto) (menuItemEntity, error) {
	now := time.Now()
	var updatedMenuItem menuItemEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		menuItemId := uuid.MustParse(input.ID)
		currentMenuItem, err := s.repository.getMenuItemByID(tx, menuItemId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentMenuItem) {
			return errMenuItemNotFound
		}
		updatedMenuItem = currentMenuItem
		updatedMenuItem.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			menuItemSameTitle, err := s.repository.getMenuItemByTitle(tx, *input.Title, false)
			if err != nil {
				return ceng_err.ErrGeneric
			}
			if !ceng_utils.IsEmpty(menuItemSameTitle) && menuItemSameTitle.ID.String() != input.ID {
				return errMenuItemSameTitleAlreadyExists
			}
			updatedMenuItem.Title = *input.Title
		}
		if input.Active != nil {
			updatedMenuItem.Active = input.Active
		}
		if input.Price != nil {
			updatedMenuItem.Price = *input.Price
		}
		if input.Position != nil {
			// If the step is moving in a lower position (e.g. from 10 to 3),
			// we need to move it one step more, so that, the algorith to re-sort all steps correctly
			if updatedMenuItem.Position < *input.Position {
				*input.Position++
			}
			updatedMenuItem.Position = *input.Position
		}
		if _, err = s.repository.saveMenuItem(tx, updatedMenuItem, ceng_db.Update); err != nil {
			return ceng_err.ErrGeneric
		}
		if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, updatedMenuItem.MenuCategoryID); err != nil {
			return ceng_err.ErrGeneric
		}
		if updatedMenuItem, err = s.repository.getMenuItemByID(tx, updatedMenuItem.ID, false); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of menuItem updated
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuItemV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuItemUpdatedEvent,
				EventEntity: &ceng_pubsub.MenuItemEventEntity{
					ID:             updatedMenuItem.ID,
					MenuCategoryID: updatedMenuItem.MenuCategoryID,
					Title:          updatedMenuItem.Title,
					Position:       updatedMenuItem.Position,
					Active:         updatedMenuItem.Active,
					Price:          updatedMenuItem.Price,
					CreatedAt:      updatedMenuItem.CreatedAt,
					UpdatedAt:      updatedMenuItem.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentMenuItem, updatedMenuItem),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == updatedMenuItem.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuItemV1, ceng_pubsub.PubSubMessage{
				Message: ceng_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: ceng_pubsub.MenuItemUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuItemEventEntity{
						ID:             updatedEntity.ID,
						MenuCategoryID: updatedEntity.MenuCategoryID,
						Title:          updatedEntity.Title,
						Position:       updatedEntity.Position,
						Active:         updatedEntity.Active,
						Price:          updatedEntity.Price,
						CreatedAt:      updatedEntity.CreatedAt,
						UpdatedAt:      updatedEntity.UpdatedAt,
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
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedMenuItem, nil
}

func (s menuItemService) deleteMenuItem(ctx *gin.Context, input deleteMenuItemInputDto) (menuItemEntity, error) {
	var currentMenuItem menuItemEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		// Check if the Menu Item exists
		menuItemId := uuid.MustParse(input.ID)
		currentMenuItem, err := s.repository.getMenuItemByID(tx, menuItemId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentMenuItem) {
			return errMenuItemNotFound
		}
		s.repository.deleteMenuItem(tx, currentMenuItem)
		if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, currentMenuItem.MenuCategoryID); err != nil {
			return ceng_err.ErrGeneric
		}
		// Send an event of menuItem deleted
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicMenuItemV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.MenuItemDeletedEvent,
				EventEntity: &ceng_pubsub.MenuItemEventEntity{
					ID:             currentMenuItem.ID,
					MenuCategoryID: currentMenuItem.MenuCategoryID,
					Title:          currentMenuItem.Title,
					Position:       currentMenuItem.Position,
					Active:         currentMenuItem.Active,
					Price:          currentMenuItem.Price,
					CreatedAt:      currentMenuItem.CreatedAt,
					UpdatedAt:      currentMenuItem.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentMenuItem, menuItemEntity{}),
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
					EventType: ceng_pubsub.MenuItemUpdatedEvent,
					EventEntity: &ceng_pubsub.MenuItemEventEntity{
						ID:             updatedEntity.ID,
						MenuCategoryID: updatedEntity.MenuCategoryID,
						Title:          updatedEntity.Title,
						Position:       updatedEntity.Position,
						Active:         updatedEntity.Active,
						CreatedAt:      updatedEntity.CreatedAt,
						UpdatedAt:      updatedEntity.UpdatedAt,
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
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentMenuItem, nil
}
