package foots

type SessionResponse struct {
	Data SessionData `json:"data"`
}

type SessionData struct {
	CsrfToken string `json:"csrfToken"`
}

type StockResponse struct {
	Name              string         `json:"name"`
	SellableUnits     []SellableUnit `json:"sellableUnits"`
	VariantAttributes []struct {
		Code          string `json:"code"`
		Sku           string `json:"sku"`
		SkuLaunchDate string `json:"skuLaunchDate,omitempty"`
	} `json:"variantAttributes"`
}

type SellableUnit struct {
	Attributes []SellableUnitAttribute `json:"attributes"`
	Code       string                  `json:"code"`
}

type SellableUnitAttribute struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AddToCartResponse struct {
	Code       string `json:"code"`
	Guid       string `json:"guid"`
	TotalPrice struct {
		FormattedValue string  `json:"formattedValue"`
		Value          float64 `json:"value"`
	} `json:"totalPrice"`
}

type SuccessfulCheckout struct {
	Order struct {
		Code       string `json:"code"`
		TotalPrice struct {
			Value float64 `json:"value"`
			FormattedValue string `json:"formattedValue"`
		} `json:"totalPrice"`
	} `json:"order"`
}

type FailedCheckout struct {
	Errors []struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"errors"`
}
