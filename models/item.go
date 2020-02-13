package models

type Item struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title" valid:"required~Title is required"`
	Description string   `json:"description" valid:"required~Description is required"`
	Price       float64  `json:"price" valid:"required~Price is required"`
	CategoryID  int      `json:"category_id,omitempty" valid:"required~category_id is required"`
	UserID      int      `json:"user_id,omitempty" db:"user_id" valid:"required~user_id is required"`
	Hash        string   `json:"hash"`
	CreatedAt   NullTime `json:"created_at" db:"created_at"`
	UpdatedAt   NullTime `json:"updated_at" db:"updated_at"`
}
