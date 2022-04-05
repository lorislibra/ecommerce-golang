package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OrderCollectionName = "Orders"

type OrderItem struct {
	Quantity int      `bson:"quantity"`
	Price    float32  `bson:"price"`
	Sku      string   `bson:"sku"`
	Procuct  *Product `bson:"-"`
}

type Order struct {
	Oid primitive.ObjectID `bson:"_id,omitempty"`

	UserOid primitive.ObjectID `bson:"user_id"`
	User    *User              `bson:"-"`

	Items     []OrderItem `bson:"items"`
	CreatedAt time.Time   `bson:"created_at"`
	Status    string      `bson:"status"`
}
