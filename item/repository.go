package item

import (
	"context"

	"go_boilerplate/models"
)

// Repository represent the item's repository contract
type Repository interface {
	GetAll(ctx context.Context, num uint) (res []ItemMapper, err error)
	GetByID(ctx context.Context, id uint) (*ItemMapper, error)
	GetByTitle(ctx context.Context, title string) (*ItemMapper, error)
	Update(ctx context.Context, item *models.Item) error
	Store(ctx context.Context, item *models.Item, fileNames []string) (int64, error)
	Delete(ctx context.Context, id uint) error
}
