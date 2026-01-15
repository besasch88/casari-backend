package auth

import (
	"gorm.io/gorm"
)

type authUserRepositoryInterface interface {
	findAuthUserByUserID(tx *gorm.DB, user_id string) (authUserEntity, error)
	findAuthUserByUsername(tx *gorm.DB, username string) (authUserEntity, error)
}

type authUserRepository struct {
}

func newAuthUserRepository() authUserRepository {
	return authUserRepository{}
}

func (r authUserRepository) findAuthUserByUsername(tx *gorm.DB, username string) (authUserEntity, error) {
	var model *authUserModel
	query := tx.Where("username = ?", username)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return authUserEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return authUserEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r authUserRepository) findAuthUserByUserID(tx *gorm.DB, user_id string) (authUserEntity, error) {
	var model *authUserModel
	query := tx.Where("id = ?", user_id)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return authUserEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return authUserEntity{}, nil
	}
	return model.toEntity(), nil
}
