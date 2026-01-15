package menu

import (
	"fmt"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuRepositoryInterface interface {
	listMenuCategories(tx *gorm.DB, forUpdate bool) ([]menuCategoryEntity, int64, error)
	listMenuItems(tx *gorm.DB, forUpdate bool) ([]menuItemEntity, int64, error)
	listMenuOptions(tx *gorm.DB, forUpdate bool) ([]menuOptionEntity, int64, error)
}

type menuRepository struct {
	relevanceThresholdConfig float64
}

func newMenuRepository(relevanceThresholdConfig float64) menuRepository {
	return menuRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r menuRepository) listMenuCategories(tx *gorm.DB, forUpdate bool) ([]menuCategoryEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuCategoryModel
	query := tx.Model(menuCategoryModel{})
	queryCount := tx.Model(menuCategoryModel{})
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", ceng_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuCategoryEntity{}, 0, result.Error
	}
	var entities []menuCategoryEntity = []menuCategoryEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r menuRepository) listMenuItems(tx *gorm.DB, forUpdate bool) ([]menuItemEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuItemModel
	query := tx.Model(menuItemModel{})
	queryCount := tx.Model(menuItemModel{})
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", ceng_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuItemEntity{}, 0, result.Error
	}
	var entities []menuItemEntity = []menuItemEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r menuRepository) listMenuOptions(tx *gorm.DB, forUpdate bool) ([]menuOptionEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuOptionModel
	query := tx.Model(menuOptionModel{})
	queryCount := tx.Model(menuOptionModel{})
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", ceng_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuOptionEntity{}, 0, result.Error
	}
	var entities []menuOptionEntity = []menuOptionEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}
