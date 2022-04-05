package dtos

import "github.com/donnjedarko/paninaro/src/entities"

type OrderCreateBodyItem struct {
	Sku      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type OrderCreateBody struct {
	Items []OrderCreateBodyItem `json:"items"`
}

func (o *OrderCreateBody) ToEntity(products []*entities.Product) *entities.Order {
	items := make([]entities.OrderItem, 0, len(o.Items))

	for _, item := range o.Items {
		var prodFound *entities.Product
		for _, prod := range products {
			if prod.Sku == item.Sku {
				prodFound = prod
				break
			}
		}

		if prodFound != nil {
			items = append(items, entities.OrderItem{
				Quantity: item.Quantity,
				Sku:      prodFound.Sku,
				Price:    prodFound.Price,
				Procuct:  prodFound,
			})
		}
	}

	return &entities.Order{
		Items:  items,
		Status: "created",
	}
}

func (o *OrderCreateBody) Skus() []string {
	skus := make([]string, 0, len(o.Items))
	for _, item := range o.Items {
		skus = append(skus, item.Sku)
	}
	return skus
}
