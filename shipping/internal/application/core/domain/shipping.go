package domain

type Item struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

type Shipping struct {
	ID         int64  `json:"id"`
	PurchaseID int64  `json:"purchase_id"`
	Items      []Item `json:"items"`
}

func NewShipping(purchaseId int64, items []Item) Shipping {
	return Shipping{
		PurchaseID: purchaseId,
		Items:      items,
	}
}
