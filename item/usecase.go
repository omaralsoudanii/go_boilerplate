package item

import (
	"context"

	"go_boilerplate/models"
)

// UseCase represent the item usecases
type UseCase interface {
	GetAll(ctx context.Context, num uint) ([]*models.Item, error)
	GetByID(ctx context.Context, id uint) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	GetByTitle(ctx context.Context, title string) (*models.Item, error)
	Store(ctx context.Context, item *models.Item, images []File) (uint, error)
	Delete(ctx context.Context, id uint) error
}
