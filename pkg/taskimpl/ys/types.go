package ys

import "nebula/third_party/http"

type ProductDetailsResponse struct {
	ID                     string             `json:"id"`
	Name                   string             `json:"name"`
	ModelNumber            string             `json:"model_number"`
	ProductType            string             `json:"product_type"`
	MetaData               MetaData           `json:"meta_data"`
	YeezyPDPCallout        []string           `json:"yeezyPDPCallout"`
	ViewList               []ViewList         `json:"view_list"`
	PricingInformation     PricingInformation `json:"pricing_information"`
	AttributeList          AttributeList      `json:"attribute_list"`
	ProductDescription     ProductDescription `json:"product_description"`
	RecommendationsEnabled bool               `json:"recommendationsEnabled"`
	ProductLinkList        []interface{}      `json:"product_link_list"`
}

type AttributeList struct {
	IsWaitingRoomProduct     bool       `json:"isWaitingRoomProduct"`
	BadgeText                string     `json:"badge_text"`
	BadgeStyle               string     `json:"badge_style"`
	Brand                    string     `json:"brand"`
	Collection               []string   `json:"collection"`
	Category                 string     `json:"category"`
	Color                    string     `json:"color"`
	ReturnType               string     `json:"return_type"`
	Gender                   string     `json:"gender"`
	Personalizable           bool       `json:"personalizable"`
	MandatoryPersonalization bool       `json:"mandatory_personalization"`
	Customizable             bool       `json:"customizable"`
	Pricebook                string     `json:"pricebook"`
	Sale                     bool       `json:"sale"`
	Outlet                   bool       `json:"outlet"`
	IsCNCRestricted          bool       `json:"isCnCRestricted"`
	SizeChartLink            string     `json:"size_chart_link"`
	Sport                    []string   `json:"sport"`
	SizeFitBar               SizeFitBar `json:"size_fit_bar"`
	PreviewTo                string     `json:"preview_to"`
	ComingSoonSignup         bool       `json:"coming_soon_signup"`
	MaxOrderQuantity         int64      `json:"max_order_quantity"`
	ProductType              []string   `json:"productType"`
	SearchColor              string     `json:"search_color"`
	SpecialLaunch            bool       `json:"specialLaunch"`
	SpecialLaunchType        string     `json:"specialLaunchType"`
	SearchColorRaw           string     `json:"search_color_raw"`
}

type SizeFitBar struct {
	Value               string `json:"value"`
	SelectedMarkerIndex int64  `json:"selectedMarkerIndex"`
	MarkerCount         int64  `json:"markerCount"`
}

type MetaData struct {
	PageTitle   string `json:"page_title"`
	SiteName    string `json:"site_name"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Canonical   string `json:"canonical"`
}

type PricingInformation struct {
	StandardPrice      int64 `json:"standard_price"`
	StandardPriceNoVat int64 `json:"standard_price_no_vat"`
	CurrentPrice       int64 `json:"currentPrice"`
}

type ProductDescription struct {
	Title             string            `json:"title"`
	Usps              []string          `json:"usps"`
	DescriptionAssets DescriptionAssets `json:"description_assets"`
}

type DescriptionAssets struct{}

type ViewList struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	Source   string `json:"source"`
}

type ProductStockResponse struct {
	ID                 string      `json:"id"`
	AvailabilityStatus string      `json:"availability_status"`
	VariationList      []Variation `json:"variation_list"`
}

type Variation struct {
	Sku                string `json:"sku"`
	Size               string `json:"size"`
	Availability       int64  `json:"availability"`
	AvailabilityStatus string `json:"availability_status"`
}

type AddToCartResponse struct {
	BasketID              string           `json:"basketId"`
	Currency              string           `json:"currency"`
	ModifiedDate          string           `json:"modifiedDate"`
	Pricing               AddToCartPricing `json:"pricing"`
	ResourceState         string           `json:"resourceState"`
	TaxationPolicy        string           `json:"taxationPolicy"`
	TotalProductCount     int64            `json:"totalProductCount"`
	MessageList           []MessageList    `json:"messageList"`
	ShipmentList          []ShipmentList   `json:"shipmentList"`
	Customer              Customer         `json:"customer"`
	FreeShippingThreshold int64            `json:"freeShippingThreshold"`
}

type Customer struct {
	CustomerID string `json:"customerId"`
	IsLoggedIn bool   `json:"isLoggedIn"`
}

type MessageList struct {
	Type    string   `json:"type"`
	Details *Details `json:"details,omitempty"`
}

type Details struct {
	ShipmentID string `json:"shipmentId"`
}

type AddToCartPricing struct {
	Total                            int64 `json:"total"`
	BaseTotal                        int64 `json:"baseTotal"`
	TotalTax                         int64 `json:"totalTax"`
	ProductTotal                     int64 `json:"productTotal"`
	ProductTotalBeforeDiscounts      int64 `json:"productTotalBeforeDiscounts"`
	ProductTotalBeforeOrderDiscounts int64 `json:"productTotalBeforeOrderDiscounts"`
	ShippingTotal                    int64 `json:"shippingTotal"`
	ShippingBaseTotal                int64 `json:"shippingBaseTotal"`
}

type ShipmentList struct {
	ShipmentID          string                `json:"shipmentId"`
	ShipmentType        string                `json:"shipmentType"`
	ProductLineItemList []ProductLineItemList `json:"productLineItemList"`
	ShippingLineItem    ShippingLineItem      `json:"shippingLineItem"`
	ShippingOnDate      string                `json:"shippingOnDate"`
}

type ProductLineItemList struct {
	ItemID               string                     `json:"itemId"`
	ProductID            string                     `json:"productId"`
	ProductName          string                     `json:"productName"`
	Category             string                     `json:"category"`
	CanonicalProductName string                     `json:"canonicalProductName"`
	ProductImage         string                     `json:"productImage"`
	Quantity             int64                      `json:"quantity"`
	Pricing              ProductLineItemListPricing `json:"pricing"`
	Gender               string                     `json:"gender"`
	Color                string                     `json:"color"`
	Size                 string                     `json:"size"`
	AllowedActions       AllowedActions             `json:"allowedActions"`
	MaxQuantityAllowed   int64                      `json:"maxQuantityAllowed"`
	IsBonusProduct       bool                       `json:"isBonusProduct"`
	ProductType          string                     `json:"productType"`
	AvailableStock       int64                      `json:"availableStock"`
	LastAdded            bool                       `json:"lastAdded"`
	IsFlashProduct       bool                       `json:"isFlashProduct"`
	SpecialLaunchProduct bool                       `json:"specialLaunchProduct"`
}

type AllowedActions struct {
	Delete         bool `json:"delete"`
	Edit           bool `json:"edit"`
	MoveToWishlist bool `json:"moveToWishlist"`
}

type ProductLineItemListPricing struct {
	BaseUnitPrice          int64 `json:"baseUnitPrice"`
	UnitPrice              int64 `json:"unitPrice"`
	BasePrice              int64 `json:"basePrice"`
	Price                  int64 `json:"price"`
	PriceAfterAllDiscounts int64 `json:"priceAfterAllDiscounts"`
	UnitPriceWithoutTax    int64 `json:"unitPriceWithoutTax"`
}

type ShippingLineItem struct {
	Name                     string                  `json:"name"`
	Description              string                  `json:"description"`
	ID                       string                  `json:"id"`
	Pricing                  ShippingLineItemPricing `json:"pricing"`
	CarrierServiceName       string                  `json:"carrierServiceName"`
	FreeShippingThresholdMin int64                   `json:"freeShippingThresholdMin"`
	FreeShippingThresholdMax int64                   `json:"freeShippingThresholdMax"`
}

type ShippingLineItemPricing struct {
	BasePrice int64 `json:"basePrice"`
	Price     int64 `json:"price"`
}

type AddToCartForm struct {
	ProductId1          string `json:"product_id"`
	ProductVariationSku string `json:"product_variation_sku"`
	ProductId2          string `json:"productId"`
	Quantity            int    `json:"quantity"`
	Size                string `json:"size"`
	DisplaySize         string `json:"displaySize"`
}

type MethodList struct {
	ID               string `json:"id"`
	ShipmentId       string `json:"shipmentId"`
	CollectionPeriod string `json:"collectionPeriod"`
	DeliveryPeriod   string `json:"deliveryPeriod"`
}

type Pricing struct {
	Total                            int `json:"total"`
	BaseTotal                        int `json:"baseTotal"`
	TotalTax                         int `json:"totalTax"`
	ProductTotal                     int `json:"productTotal"`
	ProductTotalBeforeDiscounts      int `json:"productTotalBeforeDiscounts"`
	ProductTotalBeforeOrderDiscounts int `json:"productTotalBeforeOrderDiscounts"`
	ShippingTotal                    int `json:"shippingTotal"`
	ShippingBaseTotal                int `json:"shippingBaseTotal"`
}

type SubmitOrderResponse struct {
	CreationDate  string  `json:"creationDate"`
	Currency      string  `json:"currency"`
	Exported      bool    `json:"exported"`
	OrderID       string  `json:"orderId"`
	PaymentStatus string  `json:"paymentStatus"`
	Pricing       Pricing `json:"pricing"`
	ResourceState string  `json:"resourceState"`
	Status        string  `json:"status"`
}

type FormFields struct {
	PaReq       string `json:"PaReq"`
	EncodedData string `json:"EncodedData"`
	MD          string `json:"MD"`
}

type PaRedirectForm struct {
	FormMethod string     `json:"formMethod"`
	FormAction string     `json:"formAction"`
	FormFields FormFields `json:"formFields"`
}

type ThreeDSForm struct {
	FormMethod string     `json:"formMethod"`
	FormAction string     `json:"formAction"`
	FormFields FormFields `json:"formFields"`
	OrderId    string     `json:"orderId"`
}

type ThreeDSResponse struct {
	OrderID           string         `json:"orderId"`
	ResourceState     string         `json:"resourceState"`
	PaymentStatus     string         `json:"paymentStatus"`
	Status            string         `json:"status"`
	AuthorizationType string         `json:"authorizationType"`
	PaRedirectForm    PaRedirectForm `json:"paRedirectForm"`
}

type ThreeDSData struct {
	Data       map[string]interface{} `json:"data"`
	OrderID    string                 `json:"orderId"`
	PaymentUrl string                 `json:"paymentUrl"`
	TermUrl    string                 `json:"termUrl"`
	Cookies    []http.Cookie          `json:"cookies"`
}
