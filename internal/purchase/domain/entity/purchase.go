package entity

type Purchase struct {
	PurchaserID      uint
	PurchaseTypeName string
}

type PurchaseGroup struct {
	PurchaseTypeName string `json:"type"`
	Quantity         uint   `json:"quantity"`
}

type Inventory []PurchaseGroup
