package model

type CheckoutRequest struct {
	Product CheckoutProduct `json:"product"`
}

type CheckoutProduct struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
	Url   *string `json:"url"`
	Size  *string `json:"size"`
	Price *string `json:"price"`
}
