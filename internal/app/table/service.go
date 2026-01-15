package table

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

type tableServiceInterface interface {
	listTables(ctx *gin.Context, input ListTablesInputDto) ([]tableEntity, int64, error)
	getTableByID(ctx *gin.Context, input getTableInputDto) (tableEntity, error)
	createTable(ctx *gin.Context, input createTableInputDto) (tableEntity, error)
	updateTable(ctx *gin.Context, input updateTableInputDto) (tableEntity, error)
	deleteTable(ctx *gin.Context, input deleteTableInputDto) (tableEntity, error)
}

type tableService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  tableRepositoryInterface
}

func newTableService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository tableRepositoryInterface) tableService {
	return tableService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s tableService) listTables(ctx *gin.Context, input ListTablesInputDto) ([]tableEntity, int64, error) {
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userId := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
		userId = nil
	}
	items, totalCount, err := s.repository.listTables(s.storage, userId, input.IncludeClosed, false)
	if err != nil || items == nil {
		return []tableEntity{}, 0, ceng_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s tableService) getTableByID(ctx *gin.Context, input getTableInputDto) (tableEntity, error) {
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userId := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.READ_OTHER_TABLES) {
		userId = nil
	}
	tableID := uuid.MustParse(input.ID)
	item, err := s.repository.getTableByID(s.storage, tableID, userId, false)
	if err != nil {
		return tableEntity{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(item) {
		return tableEntity{}, errTableNotFound
	}
	return item, nil
}

func (s tableService) createTable(ctx *gin.Context, input createTableInputDto) (tableEntity, error) {
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userId := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	now := time.Now()
	newTable := tableEntity{
		ID:            uuid.New(),
		UserID:        *userId,
		Name:          input.Name,
		Close:         ceng_utils.BoolPtr(false),
		PaymentMethod: nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		if item, err := s.repository.getOpenTableByName(tx, input.Name, false); err != nil {
			return ceng_err.ErrGeneric
		} else if !ceng_utils.IsEmpty(item) {
			return errTableSameNameAlreadyExists
		} else if _, err = s.repository.saveTable(tx, newTable, ceng_db.Create); err != nil {
			return ceng_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicTableV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.TableCreatedEvent,
				EventEntity: &ceng_pubsub.TableEventEntity{
					ID:            newTable.ID,
					UserID:        newTable.UserID,
					Name:          newTable.Name,
					Close:         newTable.Close,
					PaymentMethod: newTable.PaymentMethod,
					CreatedAt:     newTable.CreatedAt,
					UpdatedAt:     newTable.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(tableEntity{}, newTable),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newTable, nil
}

func (s tableService) updateTable(ctx *gin.Context, input updateTableInputDto) (tableEntity, error) {
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userId := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.WRITE_OTHER_TABLES) {
		userId = nil
	}
	now := time.Now()
	var updatedTable tableEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		tableId := uuid.MustParse(input.ID)
		currentTable, err := s.repository.getTableByID(tx, tableId, userId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentTable) {
			return errTableNotFound
		}
		updatedTable = currentTable
		updatedTable.UpdatedAt = now
		// If the input contains a new name, check for collision
		if input.Name != nil {
			tableSameName, err := s.repository.getOpenTableByName(tx, *input.Name, false)
			if err != nil {
				return ceng_err.ErrGeneric
			}
			if !ceng_utils.IsEmpty(tableSameName) && tableSameName.ID.String() != input.ID {
				return errTableSameNameAlreadyExists
			}
			updatedTable.Name = *input.Name
		}
		if input.Close != nil {
			updatedTable.Close = input.Close
		}
		if input.PaymentMethod != nil {
			updatedTable.PaymentMethod = input.PaymentMethod
		}
		if _, err = s.repository.saveTable(tx, updatedTable, ceng_db.Update); err != nil {
			return ceng_err.ErrGeneric
		}
		if updatedTable, err = s.repository.getTableByID(tx, updatedTable.ID, userId, false); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of table updated
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicTableV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.TableUpdatedEvent,
				EventEntity: &ceng_pubsub.TableEventEntity{
					ID:            updatedTable.ID,
					UserID:        updatedTable.UserID,
					Name:          updatedTable.Name,
					Close:         updatedTable.Close,
					PaymentMethod: updatedTable.PaymentMethod,
					CreatedAt:     updatedTable.CreatedAt,
					UpdatedAt:     updatedTable.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentTable, updatedTable),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedTable, nil
}

func (s tableService) deleteTable(ctx *gin.Context, input deleteTableInputDto) (tableEntity, error) {
	requester := ceng_auth.GetAuthenticatedUserFromSession(ctx)
	userId := ceng_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, ceng_auth.WRITE_OTHER_TABLES) {
		userId = nil
	}
	var currentTable tableEntity
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		// Check if the Menu Item exists
		tableId := uuid.MustParse(input.ID)
		currentTable, err := s.repository.getTableByID(tx, tableId, userId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentTable) {
			return errTableNotFound
		}
		s.repository.deleteTable(tx, currentTable)
		// Send an event of table deleted
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicTableV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.TableDeletedEvent,
				EventEntity: &ceng_pubsub.TableEventEntity{
					ID:            currentTable.ID,
					UserID:        currentTable.UserID,
					Name:          currentTable.Name,
					Close:         currentTable.Close,
					PaymentMethod: currentTable.PaymentMethod,
					CreatedAt:     currentTable.CreatedAt,
					UpdatedAt:     currentTable.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentTable, tableEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentTable, nil
}
