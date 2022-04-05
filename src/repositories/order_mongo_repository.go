package repositories

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepository struct {
	mongo   db.MongoInstance
	colName string
}

func NewOrderRepository(mongo db.MongoInstance) domains.OrderRepository {
	repo := &orderRepository{
		mongo:   mongo,
		colName: entities.OrderCollectionName,
	}

	repo.InitCollection()

	return repo
}

func (r *orderRepository) Collection() *mongo.Collection {
	return r.mongo.Db().Collection(r.colName)
}

func (r *orderRepository) InitCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := r.mongo.Db().CreateCollection(ctx, r.colName)
	if err != nil && err.(mongo.CommandError).Name != "NamespaceExists" {
		panic(err)
	}

	indexes := []mongo.IndexModel{}

	if len(indexes) > 0 {
		_, err = r.Collection().Indexes().CreateMany(ctx, indexes)
		if err != nil {
			panic(err)
		}
	}
}

func (r *orderRepository) Find(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) (*entities.Order, error) {
	foundOrder := new(entities.Order)
	err := r.Collection().FindOne(ctx, bson.M{"_id": orderOid, "user_id": userOid}).Decode(foundOrder)
	if err != nil {
		return foundOrder, err
	}

	return foundOrder, nil
}

func (r *orderRepository) FindAll(ctx context.Context) ([]*entities.Order, error) {
	curr, err := r.Collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var orders []*entities.Order
	if err = curr.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindAllByUser(ctx context.Context, userOid primitive.ObjectID) ([]*entities.Order, error) {
	curr, err := r.Collection().Find(ctx, bson.M{"user_id": userOid})
	if err != nil {
		return nil, err
	}

	var orders []*entities.Order
	if err = curr.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindFull(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) (*entities.Order, error) {
	foundOrder := new(entities.Order)
	err := r.Collection().FindOne(ctx, bson.M{"_id": orderOid, "user_id": orderOid}).Decode(foundOrder)
	if err != nil {
		return foundOrder, err
	}

	return foundOrder, nil
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
	res, err := r.Collection().InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.Oid = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *orderRepository) EditStatus(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID, status string) (bool, error) {
	res, err := r.Collection().UpdateOne(ctx, bson.M{"_id": orderOid, "user_id": userOid}, bson.M{"$set": bson.M{
		"status": status,
	}})
	if err != nil {
		return false, err
	}

	if res.ModifiedCount == 0 {
		return false, nil
	}

	return true, nil
}

func (r *orderRepository) Cancel(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) error {
	ok, err := r.EditStatus(ctx, orderOid, userOid, "cancelled")
	if err != nil {
		return err
	}

	if !ok {
		return mongo.ErrNoDocuments
	}

	return nil
}
