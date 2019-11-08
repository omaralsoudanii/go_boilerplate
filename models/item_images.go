package models

type ItemImages struct {
	ID     uint   `json:"id"`
	Hash   string `json:"hash"`
	ItemId uint   `json:"item_id"`
	Size   int64  `json:"size" `
	Type   string `json:"type"`
}
