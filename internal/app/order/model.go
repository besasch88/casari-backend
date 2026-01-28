package order

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

type orderModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	TableID   uuid.UUID `gorm:"column:table_id;type:varchar(36)"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m orderModel) TableName() string {
	return "ceng_order"
}

func (m orderModel) toEntity() orderEntity {
	return orderEntity(m)
}

type courseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	OrderID   uuid.UUID `gorm:"column:order_id;type:varchar(36)"`
	Position  int64     `gorm:"column:position;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m courseModel) TableName() string {
	return "ceng_course"
}

func (m courseModel) toEntity() courseEntity {
	return courseEntity(m)
}

type courseSelectionModel struct {
	ID           uuid.UUID  `gorm:"primaryKey;column:id;type:varchar(36)"`
	CourseID     uuid.UUID  `gorm:"column:course_id;type:varchar(36)"`
	MenuItemID   uuid.UUID  `gorm:"column:menu_item_id;type:varchar(36)"`
	MenuOptionID *uuid.UUID `gorm:"column:menu_option_id;type:varchar(36)"`
	Quantity     int64      `gorm:"column:quantity;type:bigint"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m courseSelectionModel) TableName() string {
	return "ceng_course_selection"
}

func (m courseSelectionModel) toEntity() courseSelectionEntity {
	return courseSelectionEntity(m)
}
