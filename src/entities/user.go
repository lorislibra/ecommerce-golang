package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UserCollectionName = "Users"

type User struct {
	Oid       primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	Username  string             `bson:"username"`
	FirstName string             `bson:"firstname"`
	LastName  string             `bson:"lastname"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      Role               `bson:"role"`
}
