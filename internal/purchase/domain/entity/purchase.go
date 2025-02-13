package entity

type Purchase struct {
	PurchaserId      uint
	PurchaseTypeName string
}

type PurchaseGroup struct {
	PurchaseTypeName string
	Quantity         uint
}

type Inventory []PurchaseGroup
