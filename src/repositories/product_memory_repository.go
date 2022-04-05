package repositories

import (
	"context"
	"sync"

	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

type memoryRepository struct {
	mu sync.RWMutex
	db map[string]*entities.Product
}

func NewMemoryRepository() domains.ProductRepository {
	return &memoryRepository{
		db: make(map[string]*entities.Product),
	}
}

func (r *memoryRepository) Find(ctx context.Context, sku string) (*entities.Product, error) {
	r.mu.RLock()
	defer r.mu.Unlock()

	product, found := r.db[sku]
	if !found {
		return nil, mongo.ErrNoDocuments
	}
	return product, nil
}

func (r *memoryRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entities.Product, 0, len(r.db))
	for _, p := range r.db {
		products = append(products, p)
	}

	return products, nil
}

func (r *memoryRepository) Create(ctx context.Context, product *entities.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, found := r.db[product.Sku]
	if found {
		return mongo.ErrNoDocuments
	}

	r.db[product.Sku] = product
	return nil
}

func (r *memoryRepository) SetHidden(ctx context.Context, sku string, value bool) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, found := r.db[sku]
	if found {
		product.Hidden = value
		return true, nil
	}

	return false, mongo.ErrNoDocuments
}

func (r *memoryRepository) Edit(ctx context.Context, sku string, product *dtos.ProductUpdateBody) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, found := r.db[sku]
	if !found {
		return false, mongo.ErrNoDocuments
	}

	if product.Description != "" {
		p.Description = product.Description
	}

	if product.Price > 0 {
		p.Price = product.Price
	}

	if product.Quantity >= 0 {
		p.Quantity = product.Quantity
	}

	return true, nil
}

func (r *memoryRepository) FindMany(ctx context.Context, skus []string) ([]*entities.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entities.Product, 0, len(skus))
	for _, sku := range skus {
		p, found := r.db[sku]
		if found {
			products = append(products, p)
		}
	}

	if len(products) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return products, nil

}
