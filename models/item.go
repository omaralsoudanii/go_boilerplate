package models

import (
	pg "github.com/lib/pq"
)

type Item struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title" valid:"required~Title is required"`
	Description string      `json:"description" valid:"required~Description is required"`
	Hash        string      `json:"hash"`
	Category    string      `json:"category"`
	Price       float64     `json:"price" valid:"required~Price is required"`
	CategoryID  int         `json:"category_id" valid:"required~category is required"`
	UserID      int         `json:"user_id" db:"user_id"`
	CreatedAt   pg.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt   pg.NullTime `json:"updated_at" db:"updated_at"`
}
