package repositories

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepository struct {
	mongo   db.MongoInstance
	colName string
}

func NewProductRepository(mongo db.MongoInstance) domains.ProductRepository {
	repo := &productRepository{
		mongo:   mongo,
		colName: entities.ProductsCollectionName,
	}

	repo.InitCollection()

	return repo
}

func (r *productRepository) Collection() *mongo.Collection {
	return r.mongo.Db().Collection(r.colName)
}

func (r *productRepository) InitCollection() {
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

func (r *productRepository) Find(ctx context.Context, sku string) (*entities.Product, error) {
	foundProduct := new(entities.Product)

	err := r.Collection().FindOne(ctx, bson.M{"_id": sku}).Decode(foundProduct)
	if err != nil {
		return foundProduct, err
	}

	return foundProduct, nil
}

func (r *productRepository) FindMany(ctx context.Context, skus []string) ([]*entities.Product, error) {
	curr, err := r.Collection().Find(ctx, bson.M{"_id": bson.M{"$in": skus}})
	if err != nil {
		return nil, err
	}

	var foundProducts []*entities.Product
	if err = curr.All(ctx, &foundProducts); err != nil {
		return nil, err
	}

	return foundProducts, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	curr, err := r.Collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var products []*entities.Product
	if err = curr.All(ctx, &products); err != nil {
		return products, err
	}

	return products, nil
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	res, err := r.Collection().InsertOne(ctx, product)
	if err != nil {
		return err
	}
	_ = res
	return nil
}

func (r *productRepository) SetHidden(ctx context.Context, sku string, value bool) (bool, error) {
	res, err := r.Collection().UpdateOne(ctx, bson.M{"_id": sku}, bson.M{"$set": bson.M{"hidden": value}})
	if err != nil {
		return false, err
	}

	if res.ModifiedCount == 0 {
		return false, nil
	}

	return true, nil
}

func (r *productRepository) Edit(ctx context.Context, sku string, product *dtos.ProductUpdateBody) (bool, error) {
	query := bson.M{}

	if product.Description != "" {
		query["description"] = product.Description
	}

	if product.Price > 0 {
		query["price"] = product.Price
	}

	if product.Quantity >= 0 {
		query["quantity"] = product.Quantity
	}

	res, err := r.Collection().UpdateOne(ctx, bson.M{"_id": sku}, bson.M{"$set": query})
	if err != nil {
		return false, err
	}

	if res.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}
