package item

import (
	"context"

	"go_boilerplate/models"
)

// Repository represent the item's repository contract
type Repository interface {
	GetAll(ctx context.Context, num uint) (res []*models.Item, err error)
	GetByID(ctx context.Context, id uint) (*models.Item, error)
	GetByTitle(ctx context.Context, title string) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Store(ctx context.Context, item *models.Item, fileNames []string) (uint, error)
	Delete(ctx context.Context, id uint) error
}
