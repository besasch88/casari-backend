package course

import (
	"time"

	"github.com/google/uuid"
)

type tableModel struct {
	ID uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
}

func (m tableModel) TableName() string {
	return "ceng_table"
}

type courseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID    uuid.UUID `gorm:"column:user_id;type:varchar(36)"`
	TableID   uuid.UUID `gorm:"column:table_id;type:varchar(36)"`
	Close     *bool     `gorm:"column:close;type:boolean"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m courseModel) TableName() string {
	return "ceng_course"
}

func (m courseModel) toEntity() courseEntity {
	return courseEntity(m)
}
