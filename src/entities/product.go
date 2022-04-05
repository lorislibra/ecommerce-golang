package entities

import (
	"time"
)

const ProductsCollectionName = "Products"

type Product struct {
	Sku         string  `bson:"_id"`
	Title       string  `bson:"title"`
	Description string  `bson:"description"`
	Price       float32 `bson:"price"`
	Quantity    int     `bson:"quantity"`
	Hidden      bool    `bson:"hidden"`
	Image       string  `bson:"image"`

	CreatedAt time.Time `bson:"created_at"`
}
