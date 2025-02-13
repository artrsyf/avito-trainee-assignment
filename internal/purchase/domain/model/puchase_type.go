package model

type PurchaseType struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
	Cost uint   `db:"cost"`
}
