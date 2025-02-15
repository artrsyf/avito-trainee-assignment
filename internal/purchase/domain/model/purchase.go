package model

type Purchase struct {
	ID             uint `db:"id"`
	PurchaserID    uint `db:"purchaser_id"`
	PurchaseTypeID uint `db:"purchase_type_id"`
}
