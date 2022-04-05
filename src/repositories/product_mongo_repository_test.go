package repositories_test

import (
	"context"
	"testing"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/repositories"
)

func TestCreateProduct(t *testing.T) {
	_ = config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoInstance := db.NewMongoInstance()
	repo := repositories.NewProductRepository(mongoInstance)

	product1 := dtos.ProductCreateBody{
		Title:       "big mac",
		Description: "big mac gordo buono infinito",
		Price:       4.5,
		Quantity:    100,
		Sku:         "BM01",
	}

	productEntity := product1.ToEntity()

	err := repo.Create(ctx, productEntity)
	if err != nil {
		t.Error(err)
	}
}
