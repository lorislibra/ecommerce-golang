package dtos

import "github.com/donnjedarko/paninaro/src/entities"

type OrderRespItem struct {
	Sku      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

type OrderResp struct {
	Items  []OrderRespItem `json:"items"`
	Status string          `json:"status"`
	Id     string          `json:"id"`
}

func OrderRespFromEntity(o *entities.Order) *OrderResp {
	items := make([]OrderRespItem, 0, len(o.Items))

	for _, i := range o.Items {
		items = append(items, OrderRespItem{
			Sku:      i.Sku,
			Quantity: i.Quantity,
			Price:    i.Price,
		})
	}

	return &OrderResp{
		Items:  items,
		Status: o.Status,
		Id:     o.Oid.Hex(),
	}
}
