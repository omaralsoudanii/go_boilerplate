package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_boilerplate/item"
	_lib "go_boilerplate/lib"
	"go_boilerplate/models"
	"time"
)

type itemRepository struct {
	Conn *sqlx.DB
}

func NewItemRepository(Conn *sqlx.DB) item.Repository {

	return &itemRepository{Conn}
}

var selectQuery = `SELECT item.id, item.hash, item.title, item.description, item.price,
						category.title, category.id,
						item_images.id, item_images.hash, item_images.type, item_images.size,
						user.id, user.username, user.email,
						item.created_at,item.updated_at
						FROM
						item
						JOIN
						category
						ON
						item.category_id = category.id
						JOIN
						item_images
						ON
						item.id = item_images.item_id
						JOIN
						user
						ON
						item.user_id = user.id`

func (repo itemRepository) fetch(ctx context.Context, query string, args ...interface{}) (data []*item.ItemMapper, err error) {
	rows, err := repo.Conn.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	records := make(map[uint]*item.ItemMapper, 0)
	images := make(map[uint][]*models.ItemImages, 0)
	for rows.Next() {
		model := new(item.ItemScanner)
		err := rows.StructScan(&model)
		if err != nil {
			return nil, err
		}
		images[model.Item.ID] = append(images[model.Item.ID], model.ItemImages)
		data := item.ItemMapper{
			Item:       model.Item,
			Category:   model.Category,
			User:       model.User,
			ItemImages: images[model.Item.ID],
		}
		records[model.Item.ID] = &data
	}
	result := make([]*item.ItemMapper, 0)
	for _, r := range records {
		result = append(result, r)
	}
	return result, nil
}

func (repo *itemRepository) GetAll(ctx context.Context, num uint) ([]*item.ItemMapper, error) {
	query := selectQuery + " ORDER BY item.title LIMIT 10 OFFSET ? "
	res, err := repo.fetch(ctx, query, num)
	fmt.Print(err)
	if err != nil {
		return nil, _lib.ErrInternalServerError
	}
	return res, nil
}
func (repo *itemRepository) GetByID(ctx context.Context, id uint) (*item.ItemMapper, error) {
	query := selectQuery + " WHERE id = ?"
	itemModel := &item.ItemMapper{}
	err := repo.Conn.GetContext(ctx, itemModel, query, id)
	if err != nil {
		return nil, _lib.ErrNotFound
	}

	return itemModel, nil
}

func (repo *itemRepository) GetByTitle(ctx context.Context, title string) (*item.ItemMapper, error) {
	query := selectQuery + " WHERE title = ?"
	itemModel := &item.ItemMapper{}
	err := repo.Conn.GetContext(ctx, itemModel, query, title)
	if err != nil {
		return nil, _lib.ErrNotFound
	}

	return itemModel, nil
}

func (repo *itemRepository) Store(ctx context.Context, i *models.Item, fileNames []string) (int64, error) {

	query := "INSERT INTO item ( title, description , price , user_id , category_id, created_at , updated_at) " +
		"VALUES (?, ? , ? , ? , ? ,? , ?)"
	fileQuery := "INSERT INTO item_images ( item_id, hash , created_at , updated_at) " +
		"VALUES (?, ? , ? , ?)"

	tx, err := repo.Conn.BeginTxx(ctx, nil)

	if tx == nil {
		return 0, _lib.ErrInternalServerError
	}

	if err != nil {
		_ = tx.Rollback()
		return 0, _lib.ErrInternalServerError
	}

	result, err := tx.ExecContext(ctx, query, i.Title, i.Description, i.Price, i.UserID, i.CategoryID, time.Now(), time.Now())
	if err != nil {
		_ = tx.Rollback()
		return 0, _lib.ErrBadParamInput
	}

	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, _lib.ErrInternalServerError
	}

	for _, fileName := range fileNames {
		_, err = tx.ExecContext(ctx, fileQuery, id, fileName, time.Now(), time.Now())
		if err != nil {
			_ = tx.Rollback()
			return 0, _lib.ErrBadParamInput
		}
	}

	err = tx.Commit()
	return id, nil
}

func (repo *itemRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM item WHERE id = ?"

	stmt, err := repo.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return _lib.ErrNotFound
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAffected)
		return _lib.ErrNotFound
	}

	return nil
}

func (repo *itemRepository) Update(ctx context.Context, i *models.Item) error {
	query := `UPDATE item set title=?, description=? , price=? , updated_at = ? , category_id = ? , user_id = ?
		WHERE ID = ?`

	stmt, err := repo.Conn.PreparexContext(ctx, query)
	if err != nil {
		return _lib.ErrInternalServerError
	}

	res, err := stmt.ExecContext(ctx, i.Title, i.Description, i.Price, time.Now(), i.CategoryID, i.UserID, i.ID)
	if err != nil {
		return _lib.ErrInternalServerError
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return _lib.ErrNotFound
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		return _lib.ErrInternalServerError
	}

	return nil
}
