package models

type Category struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title" valid:"required~Title is required"`
}
