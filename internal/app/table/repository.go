package table

import (
	"fmt"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tableRepositoryInterface interface {
	listTables(tx *gorm.DB, userId *uuid.UUID, includeClosed *bool, forUpdate bool) ([]tableEntity, int64, error)
	getTableByID(tx *gorm.DB, tableID uuid.UUID, userId *uuid.UUID, forUpdate bool) (tableEntity, error)
	getOpenTableByName(tx *gorm.DB, tableName string, forUpdate bool) (tableEntity, error)
	saveTable(tx *gorm.DB, table tableEntity, operation ceng_db.SaveOperation) (tableEntity, error)
	deleteTable(tx *gorm.DB, table tableEntity) (tableEntity, error)
}

type tableRepository struct {
	relevanceThresholdConfig float64
}

func newTableRepository(relevanceThresholdConfig float64) tableRepository {
	return tableRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r tableRepository) listTables(tx *gorm.DB, userId *uuid.UUID, includeClosed *bool, forUpdate bool) ([]tableEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*tableModel
	query := tx.Model(tableModel{})
	queryCount := tx.Model(tableModel{})

	if userId != nil {
		query = query.Where("user_id = ?", userId)
		queryCount = queryCount.Where("user_id = ?", userId)
	}
	if includeClosed == nil || !*includeClosed {
		query = query.Where("is_open = ?", true)
		queryCount = queryCount.Where("is_open = ?", true)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "created_at", ceng_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []tableEntity{}, 0, result.Error
	}
	var entities []tableEntity = []tableEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r tableRepository) getTableByID(tx *gorm.DB, tableID uuid.UUID, userId *uuid.UUID, forUpdate bool) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID)
	if userId != nil {
		query = query.Where("user_id = ?", userId)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return tableEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return tableEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r tableRepository) getOpenTableByName(tx *gorm.DB, tableName string, forUpdate bool) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("name = ?", tableName).Where("is_open = ?", true)

	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return tableEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return tableEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r tableRepository) saveTable(tx *gorm.DB, table tableEntity, operation ceng_db.SaveOperation) (tableEntity, error) {
	var model = tableModel(table)
	var err error
	switch operation {
	case ceng_db.Create:
		err = tx.Create(model).Error
	case ceng_db.Update:
		err = tx.Updates(model).Error
	case ceng_db.Upsert:
		err = tx.Save(model).Error
	}
	if err != nil {
		return tableEntity{}, err
	}
	return table, nil
}

func (r tableRepository) deleteTable(tx *gorm.DB, table tableEntity) (tableEntity, error) {
	var model = tableModel(table)
	err := tx.Delete(model).Error
	if err != nil {
		return tableEntity{}, err
	}
	return table, nil
}
