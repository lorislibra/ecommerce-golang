package dtos

import (
	"time"

	"github.com/donnjedarko/paninaro/src/entities"
)

type ProductResp struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    int       `json:"quantity"`
	Sku         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
}

func ProductRespFromEntity(p *entities.Product) *ProductResp {
	return &ProductResp{
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		Sku:         p.Sku,
		CreatedAt:   p.CreatedAt,
	}
}
