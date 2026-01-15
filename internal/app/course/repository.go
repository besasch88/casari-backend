package course

import (
	"fmt"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type courseRepositoryInterface interface {
	checkTableExists(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (bool, error)
	listCourses(tx *gorm.DB, tableID uuid.UUID, forUpdate bool) ([]courseEntity, int64, error)
	getCourseByID(tx *gorm.DB, courseID uuid.UUID, forUpdate bool) (courseEntity, error)
	saveCourse(tx *gorm.DB, course courseEntity, operation ceng_db.SaveOperation) (courseEntity, error)
	deleteCourse(tx *gorm.DB, course courseEntity) (courseEntity, error)
}

type courseRepository struct {
	relevanceThresholdConfig float64
}

func newCourseRepository(relevanceThresholdConfig float64) courseRepository {
	return courseRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r courseRepository) checkTableExists(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (bool, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID)
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 || ceng_utils.IsEmpty(model) {
		return false, nil
	}
	return true, nil
}

func (r courseRepository) listCourses(tx *gorm.DB, tableID uuid.UUID, forUpdate bool) ([]courseEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*courseModel
	query := tx.Model(courseModel{}).Where("table_id = ?", tableID)
	queryCount := tx.Model(courseModel{}).Where("table_id = ?", tableID)

	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "created_at", ceng_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []courseEntity{}, 0, result.Error
	}
	var entities []courseEntity = []courseEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r courseRepository) getCourseByID(tx *gorm.DB, courseID uuid.UUID, forUpdate bool) (courseEntity, error) {
	var model *courseModel
	query := tx.Where("id = ?", courseID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return courseEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return courseEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r courseRepository) saveCourse(tx *gorm.DB, course courseEntity, operation ceng_db.SaveOperation) (courseEntity, error) {
	var model = courseModel(course)
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
		return courseEntity{}, err
	}
	return course, nil
}

func (r courseRepository) deleteCourse(tx *gorm.DB, course courseEntity) (courseEntity, error) {
	var model = courseModel(course)
	err := tx.Delete(model).Error
	if err != nil {
		return courseEntity{}, err
	}
	return course, nil
}
