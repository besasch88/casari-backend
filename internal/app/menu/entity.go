package menu

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
)

type menuCategoryEntity ceng_pubsub.MenuCategoryEventEntity
type menuItemEntity ceng_pubsub.MenuItemEventEntity
type menuOptionEntity ceng_pubsub.MenuOptionEventEntity

type menuOption struct {
	menuOptionEntity
}
type menuItem struct {
	menuItemEntity
	Options []menuOption `json:"options"`
}
type menuCategory struct {
	menuCategoryEntity
	Items []menuItem `json:"items"`
}

type menu struct {
	Categories []menuCategory `json:"categories"`
}
