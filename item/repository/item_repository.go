package repository

import (
	"context"
	"go_boilerplate/item"
	"go_boilerplate/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type itemRepository struct {
	sb squirrel.StatementBuilderType
	db *sqlx.DB
}

func NewItemRepository(sb squirrel.StatementBuilderType, db *sqlx.DB) item.Repository {

	return &itemRepository{sb, db}
}

func (repo itemRepository) scanResult(ctx context.Context, rows *sqlx.Rows) (data []*models.Item, err error) {

	result := []*models.Item{}

	for rows.Next() {
		var itemModel models.Item
		err := rows.StructScan(&itemModel)
		if err != nil {
			return nil, err
		}
		result = append(result, &itemModel)
	}
	return result, nil
}
func (repo *itemRepository) GetAll(ctx context.Context, num uint) ([]*models.Item, error) {
	q, args, err := repo.sb.
		Select("item.id",
			"item.title",
			"item.description",
			"item.price",
			"category.title as category").
		From("item").
		Join("category", squirrel.Eq{"item.category_id": "category.id"}).
		Join("item_images", squirrel.Eq{"item.id": "item_images.item_id"}).
		OrderBy("item.title").
		Limit(10).
		Offset(uint64(num)).
		ToSql()
	rows, err := repo.db.QueryxContext(ctx, q, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := repo.scanResult(ctx, rows)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (repo *itemRepository) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	q := repo.sb.Select("id, title, description, price, created_at, updated_at").
		From("item").Where(squirrel.Eq{"id": id})
	itemModel := &models.Item{}
	err := q.ScanContext(ctx, itemModel)
	if err != nil {
		return nil, err
	}
	return itemModel, nil
}

func (repo *itemRepository) GetByTitle(ctx context.Context, title string) (*models.Item, error) {
	q := repo.sb.Select("id, title, description, price, created_at, updated_at").
		From("item").Where(squirrel.Eq{"title": title})
	itemModel := &models.Item{}
	err := q.ScanContext(ctx, itemModel)
	if err != nil {
		return nil, err
	}
	return itemModel, nil
}

func (repo *itemRepository) Store(ctx context.Context, i *models.Item, fileNames []string) (uint, error) {

	q, args, err := repo.sb.Insert("item").
		Columns("title", "description", "price", "user_id", "category_id", "created_at").
		Values(i.Title, i.Description, i.Price, i.UserID, i.CategoryID, time.Now()).
		ToSql()

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	res, err := tx.ExecContext(ctx, q, args)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	iq := repo.sb.Insert("item_images").
		Columns("item_id", "hash", "created_at")

	var iqs string
	var iargs []interface{}
	for fileName := range fileNames {
		iqs, iargs, err = iq.Values(id, fileName, time.Now()).ToSql()
		_, err = tx.ExecContext(ctx, iqs, iargs)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return uint(id), nil
}

func (repo *itemRepository) Delete(ctx context.Context, id uint) error {
	q := repo.sb.Delete("").
		From("item").
		Where("id = ?", id)
	res, err := q.ExecContext(ctx)
	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil || c == 0 {
		return err
	}
	return nil
}

func (repo *itemRepository) Update(ctx context.Context, i *models.Item) error {
	data := map[string]interface{}{
		"title":       i.Title,
		"description": i.Description,
		"price":       i.Price,
	}
	q := repo.sb.Update("item").
		SetMap(data).
		Where("id = ?", i.ID)
	res, err := q.ExecContext(ctx)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil || c == 0 {

		return err
	}

	return nil
}
