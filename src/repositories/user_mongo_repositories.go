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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userMongoRepository struct {
	mongo   db.MongoInstance
	colName string
}

func NewUserMongoRepository(mongo db.MongoInstance) domains.UserRepository {
	repo := &userMongoRepository{
		mongo:   mongo,
		colName: entities.UserCollectionName,
	}
	repo.InitCollection()

	return repo
}

func (r *userMongoRepository) InitCollection() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := r.mongo.Db().CreateCollection(ctx, r.colName)
	if err != nil && err.(mongo.CommandError).Name != "NamespaceExists" {
		panic(err)
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
	}

	if len(indexes) > 0 {
		_, err = r.Collection().Indexes().CreateMany(ctx, indexes)
		if err != nil {
			panic(err)
		}
	}

}

func (r *userMongoRepository) Collection() *mongo.Collection {
	return r.mongo.Db().Collection(r.colName)
}

func (r *userMongoRepository) FindAll(ctx context.Context) ([]*entities.User, error) {
	curr, err := r.Collection().Find(ctx, bson.M{
		"role": bson.M{
			"$in": bson.A{entities.Student, entities.Admin},
		},
	})
	if err != nil {
		return nil, err
	}

	var users []*entities.User
	if err = curr.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userMongoRepository) FindByUserOrEmail(ctx context.Context, username string, email string) (*entities.User, error) {
	foundUser := new(entities.User)
	err := r.Collection().FindOne(ctx, bson.M{
		"$or": bson.A{
			bson.M{"username": username},
			bson.M{"email": email},
		},
	}).Decode(foundUser)
	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

func (r *userMongoRepository) Find(ctx context.Context, oid primitive.ObjectID) (*entities.User, error) {
	foundUser := new(entities.User)
	err := r.Collection().FindOne(ctx, bson.M{"_id": oid}).Decode(foundUser)
	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

func (r *userMongoRepository) Create(ctx context.Context, user *entities.User) error {
	res, err := r.Collection().InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.Oid = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *userMongoRepository) Delete(ctx context.Context, oid primitive.ObjectID) error {
	res, err := r.Collection().DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
