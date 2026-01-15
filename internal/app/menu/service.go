package menu

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_err"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type menuServiceInterface interface {
	getMenu(ctx *gin.Context) (menu, error)
}

type menuService struct {
	storage     *gorm.DB
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  menuRepositoryInterface
}

func newMenuService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository menuRepositoryInterface) menuService {
	return menuService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuService) getMenu(ctx *gin.Context) (menu, error) {
	categories, _, err := s.repository.listMenuCategories(s.storage, false)
	if err != nil || categories == nil {
		return menu{}, ceng_err.ErrGeneric
	}
	items, _, err := s.repository.listMenuItems(s.storage, false)
	if err != nil || items == nil {
		return menu{}, ceng_err.ErrGeneric
	}
	options, _, err := s.repository.listMenuOptions(s.storage, false)
	if err != nil || options == nil {
		return menu{}, ceng_err.ErrGeneric
	}

	menu := menu{
		Categories: []menuCategory{},
	}
	for _, category := range categories {
		if !*category.Active {
			continue
		}
		menuCategory := menuCategory{
			menuCategoryEntity: category,
			Items:              []menuItem{},
		}
		for _, item := range items {
			menuItem := menuItem{
				menuItemEntity: item,
				Options:        []menuOption{},
			}
			if !*item.Active {
				continue
			}
			if item.MenuCategoryID.String() != category.ID.String() {
				continue
			}
			for _, option := range options {
				if !*option.Active {
					continue
				}
				if option.MenuItemID.String() != item.ID.String() {
					continue
				}
				menuOption := menuOption{
					menuOptionEntity: option,
				}
				menuItem.Options = append(menuItem.Options, menuOption)
			}
			menuCategory.Items = append(menuCategory.Items, menuItem)
		}
		menu.Categories = append(menu.Categories, menuCategory)

	}
	return menu, nil
}
