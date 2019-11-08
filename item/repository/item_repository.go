package repository

import (
	"context"
	"fmt"
	"go_boilerplate/item"
	"go_boilerplate/models"
	"time"

	_lib "go_boilerplate/lib"

	"github.com/jmoiron/sqlx"
)

var log = _lib.GetLogger()

type itemRepository struct {
	Conn *sqlx.DB
}

func NewItemRepository(Conn *sqlx.DB) item.Repository {

	return &itemRepository{Conn}
}

func (repo itemRepository) fetch(ctx context.Context, query string, args ...interface{}) (data []*models.Item, err error) {
	rows, err := repo.Conn.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.Item, 0)
	for rows.Next() {
		itemModel := new(models.Item)
		err := rows.Scan(
			&itemModel.ID,
			&itemModel.Title,
			&itemModel.Description,
			&itemModel.Price,
			&itemModel.Category,
			&itemModel.Hash,
			&itemModel.CreatedAt.Time,
			&itemModel.UpdatedAt.Time,
		)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		result = append(result, itemModel)
	}
	return result, nil
}
func (repo *itemRepository) Fetch(ctx context.Context, num int64) ([]*models.Item, error) {
	fmt.Println(num)
	query := `SELECT tbl_item.id,tbl_item.title,tbl_item.description,tbl_item.price,tbl_category.title as category , tbl_item_images.hash ,
	tbl_item.created_at,tbl_item.updated_at
			FROM tbl_item 
			join 
			tbl_category 
			on 
			tbl_item.category_id = tbl_category.id
			join 
			tbl_item_images 
			on
			tbl_item.id = tbl_item_images.item_id
			ORDER BY title LIMIT 12 OFFSET $1 `
	res, err := repo.fetch(ctx, query, num)
	if err != nil {
		log.Error(err)
		return nil, _lib.ErrInternalServerError
	}
	return res, err
}
func (repo *itemRepository) GetByID(ctx context.Context, id int64) (*models.Item, error) {
	query := `SELECT id,title,description,price,created_at,updated_at
  						FROM tbl_item WHERE id = $1`
	itemModel := &models.Item{}
	err := repo.Conn.GetContext(ctx, itemModel, query, id)
	if err != nil {
		log.Error(err)
		return nil, _lib.ErrNotFound
	}

	return itemModel, nil
}

func (repo *itemRepository) GetByTitle(ctx context.Context, title string) (*models.Item, error) {
	query := `SELECT id,title,description,price,created_at,updated_at
  						FROM tbl_item WHERE title = $1`
	itemModel := &models.Item{}
	err := repo.Conn.GetContext(ctx, itemModel, query, title)
	if err != nil {
		log.Error(err)
		return nil, _lib.ErrNotFound
	}

	return itemModel, nil
}

func (repo *itemRepository) Store(ctx context.Context, i *models.Item, fileNames []string) (uint, error) {

	query := "INSERT INTO tbl_item ( title, description , price , user_id , category_id, created_at , updated_at) " +
		"VALUES ($1, $2 , $3 , $4 , $5 ,$6 , $7) RETURNING id"
	fileQuery := "INSERT INTO tbl_item_images ( item_id, hash , created_at , updated_at) " +
		"VALUES ($1, $2 , $3 , $4)"

	tx, err := repo.Conn.Beginx()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return 0, _lib.ErrBadParamInput
	}
	result := tx.QueryRowxContext(ctx, query, i.Title, i.Description, i.Price, i.UserID, i.CategoryID, time.Now(), time.Now())
	var id uint
	err = result.Scan(&id)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return 0, _lib.ErrBadParamInput
	}
	for _, fileName := range fileNames {
		_, err = tx.ExecContext(ctx, fileQuery, id, fileName, time.Now(), time.Now())
		if err != nil {
			tx.Rollback()
			log.Error(err)
			return 0, _lib.ErrBadParamInput
		}
	}
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return 0, _lib.ErrBadParamInput
	}
	err = tx.Commit()
	return id, nil
}

func (repo *itemRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM tbl_item WHERE id = $1"

	stmt, err := repo.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error(err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return _lib.ErrNotFound
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAffected)
		return _lib.ErrNotFound
	}

	return nil
}
func (repo *itemRepository) Update(ctx context.Context, i *models.Item) error {
	query := `UPDATE tbl_item set title=?, description=? , price=? , updated_at
		WHERE ID = $1`

	stmt, err := repo.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, i.Title, i.Description, i.Price, time.Now(), i.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return _lib.ErrNotFound
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}
