package model

// AggregateID for Sale aggregate.
const AggregateID = 3

// Sale represents the Sale aggregate.
type Sale struct {
	SaleID    string     `bson:"saleID,omitempty" json:"saleID,omitempty"`
	Items     []SoldItem `bson:"items,omitempty" json:"items,omitempty"`
	Timestamp int64      `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

// SoldItem defines an item sold in sale.
type SoldItem struct {
	ItemID      string  `bson:"itemID,omitempty" json:"itemID,omitempty"`
	DateArrived int64   `bson:"dateArrived,omitempty" json:"dateArrived,omitempty"`
	Lot         string  `bson:"lot,omitempty" json:"lot,omitempty"`
	Name        string  `bson:"name,omitempty" json:"name,omitempty"`
	Price       float64 `bson:"price,omitempty" json:"price,omitempty"`
	SKU         string  `bson:"sku,omitempty" json:"sku,omitempty"`
	SoldWeight  float64 `bson:"soldWeight,omitempty" json:"soldWeight,omitempty"`
	TotalWeight float64 `bson:"totalWeight,omitempty" json:"totalWeight,omitempty"`
	UPC         string  `bson:"upc,omitempty" json:"upc,omitempty"`
}
