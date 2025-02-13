package model

type Purchase struct {
	ID             uint `db:"id"`
	PurchaserId    uint `db:"purchaser_id"`
	PurchaseTypeId uint `db:"purchase_type_id"`
}
