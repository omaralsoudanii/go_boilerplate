package models

type ItemImages struct {
	ID     uint   `json:"id"`
	Hash   string `json:"hash"`
	ItemID uint   `json:"item_id,omitempty"`
	Size   int64  `json:"size" `
	Type   string `json:"type"`
}
