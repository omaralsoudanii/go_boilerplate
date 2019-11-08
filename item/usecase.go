package item

import (
	"context"

	"go_boilerplate/models"
)

// UseCase represent the item usecases
type UseCase interface {
	Fetch(ctx context.Context, num int64) ([]*models.Item, error)
	GetByID(ctx context.Context, id int64) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	GetByTitle(ctx context.Context, name string) (*models.Item, error)
	Store(ctx context.Context, item *models.Item, images []File) (uint, error)
	Delete(ctx context.Context, id int64) error
}
