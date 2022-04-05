package dtos

import "github.com/donnjedarko/paninaro/src/entities"

type ProductCreateBody struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
	Sku         string  `json:"sku"`
}

func (p *ProductCreateBody) ToEntity() *entities.Product {
	return &entities.Product{
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
		Sku:         p.Sku,
	}
}

type ProductUpdateBody struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
	Hidden      bool    `json:"hidden"`
}

func (p *ProductUpdateBody) ToEntity() *entities.Product {
	return &entities.Product{
		Hidden:      p.Hidden,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
	}
}
