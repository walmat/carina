package shopify

type CaptchaChannelResult struct {
	Token string
	Error error
}

type GetConfigResponse struct {
	PaymentInstruments struct {
		AccessToken     string      `json:"accessToken"`
		AmazonPayConfig interface{} `json:"amazonPayConfig"`
		ApplePayConfig  struct {
			ShopID                        int      `json:"shopId"`
			CountryCode                   string   `json:"countryCode"`
			CurrencyCode                  string   `json:"currencyCode"`
			MerchantCapabilities          []string `json:"merchantCapabilities"`
			MerchantID                    string   `json:"merchantId"`
			MerchantName                  string   `json:"merchantName"`
			RequiredBillingContactFields  []string `json:"requiredBillingContactFields"`
			RequiredShippingContactFields []string `json:"requiredShippingContactFields"`
			ShippingType                  string   `json:"shippingType"`
			SupportedNetworks             []string `json:"supportedNetworks"`
			Total                         struct {
				Type   string `json:"type"`
				Label  string `json:"label"`
				Amount string `json:"amount"`
			} `json:"total"`
		} `json:"applePayConfig"`
		CheckoutConfig struct {
			Domain string `json:"domain"`
			ShopID int    `json:"shopId"`
		} `json:"checkoutConfig"`
		ShopifyPayConfig struct {
			Domain                string `json:"domain"`
			ShopID                int    `json:"shopId"`
			Accelerated           bool   `json:"accelerated"`
			SupportsLogin         bool   `json:"supportsLogin"`
			ExperimentTestGroup   bool   `json:"experimentTestGroup"`
			MerchantID            string `json:"merchantId"`
			SupportsSubscriptions bool   `json:"supportsSubscriptions"`
		} `json:"shopifyPayConfig"`
		Currency        string      `json:"currency"`
		GooglePayConfig interface{} `json:"googlePayConfig"`
		Locale          string      `json:"locale"`
		PaypalConfig    struct {
			Domain         string `json:"domain"`
			Environment    string `json:"environment"`
			MerchantID     string `json:"merchantId"`
			ButtonVersion  string `json:"buttonVersion"`
			VenmoSupported bool   `json:"venmoSupported"`
			Locale         string `json:"locale"`
			ShopID         int    `json:"shopId"`
		} `json:"paypalConfig"`
		OffsiteConfigs    interface{} `json:"offsiteConfigs"`
		SupportsDiscounts bool        `json:"supportsDiscounts"`
		SupportsGiftCards bool        `json:"supportsGiftCards"`
		CheckoutDisabled  bool        `json:"checkoutDisabled"`
	} `json:"paymentInstruments"`
}

type VariantOptionResponse struct {
	ID               int64       `json:"id"`
	Title            string      `json:"title"`
	Option1          string      `json:"option1"`
	Option2          interface{} `json:"option2"`
	Option3          interface{} `json:"option3"`
	Sku              string      `json:"sku"`
	RequiresShipping bool        `json:"requires_shipping"`
	Taxable          bool        `json:"taxable"`
	FeaturedImage    interface{} `json:"featured_image"`
	Available        bool        `json:"available"`
	Price            string      `json:"price"`
	Grams            int         `json:"grams"`
	CompareAtPrice   interface{} `json:"compare_at_price"`
	Position         int         `json:"position"`
	ProductID        int64       `json:"product_id"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
}

type GetProductResponse struct {
	Products []struct {
		ID          int64                   `json:"id"`
		Title       string                  `json:"title"`
		Handle      string                  `json:"handle"`
		BodyHTML    string                  `json:"body_html"`
		PublishedAt string                  `json:"published_at"`
		CreatedAt   string                  `json:"created_at"`
		UpdatedAt   string                  `json:"updated_at"`
		Vendor      string                  `json:"vendor"`
		ProductType string                  `json:"product_type"`
		Tags        []string                `json:"tags"`
		Variants    []VariantOptionResponse `json:"variants"`
		Images      []struct {
			ID         int64         `json:"id"`
			CreatedAt  string        `json:"created_at"`
			Position   int           `json:"position"`
			UpdatedAt  string        `json:"updated_at"`
			ProductID  int64         `json:"product_id"`
			VariantIds []interface{} `json:"variant_ids"`
			Src        string        `json:"src"`
			Width      int           `json:"width"`
			Height     int           `json:"height"`
		} `json:"images"`
		Options []struct {
			Name     string   `json:"name"`
			Position int      `json:"position"`
			Values   []string `json:"values"`
		} `json:"options"`
	} `json:"products"`
}

type SubmitCartResponse struct {
	ID                           int64         `json:"id"`
	Properties                   interface{}   `json:"properties"`
	Quantity                     int           `json:"quantity"`
	VariantID                    int64         `json:"variant_id"`
	Key                          string        `json:"key"`
	Title                        string        `json:"title"`
	Price                        int           `json:"price"`
	OriginalPrice                int           `json:"original_price"`
	DiscountedPrice              int           `json:"discounted_price"`
	LinePrice                    int           `json:"line_price"`
	OriginalLinePrice            int           `json:"original_line_price"`
	TotalDiscount                int           `json:"total_discount"`
	Discounts                    []interface{} `json:"discounts"`
	Sku                          string        `json:"sku"`
	Grams                        int           `json:"grams"`
	Vendor                       string        `json:"vendor"`
	Taxable                      bool          `json:"taxable"`
	ProductID                    int64         `json:"product_id"`
	ProductHasOnlyDefaultVariant bool          `json:"product_has_only_default_variant"`
	GiftCard                     bool          `json:"gift_card"`
	FinalPrice                   int           `json:"final_price"`
	FinalLinePrice               int           `json:"final_line_price"`
	URL                          string        `json:"url"`
	FeaturedImage                struct {
		AspectRatio float64 `json:"aspect_ratio"`
		Alt         string  `json:"alt"`
		Height      int     `json:"height"`
		URL         string  `json:"url"`
		Width       int     `json:"width"`
	} `json:"featured_image"`
	Image              string   `json:"image"`
	Handle             string   `json:"handle"`
	RequiresShipping   bool     `json:"requires_shipping"`
	ProductType        string   `json:"product_type"`
	ProductTitle       string   `json:"product_title"`
	ProductDescription string   `json:"product_description"`
	VariantTitle       string   `json:"variant_title"`
	VariantOptions     []string `json:"variant_options"`
	OptionsWithValues  []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"options_with_values"`
	LineLevelDiscountAllocations []interface{} `json:"line_level_discount_allocations"`
	LineLevelTotalDiscount       int           `json:"line_level_total_discount"`
}

type CustomerDataApiAddress struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	CountryCode  string `json:"country_code"`
	ProvinceCode string `json:"province_code"`
	City         string `json:"city"`
	Zip          string `json:"zip"`
	Phone        string `json:"phone"`
}

type LineItemsApiForm struct {
	VariantID  interface{}       `json:"variant_id"`
	Quantity   int64             `json:"quantity"`
	Properties map[string]string `json:"properties"`
}

type CustomerDataApiCheckoutForm struct {
	Email           string                 `json:"email"`
	Secret          bool                   `json:"secret"`
	ShippingAddress CustomerDataApiAddress `json:"shipping_address"`
	BillingAddress  CustomerDataApiAddress `json:"billing_address"`
	LineItems       []LineItemsApiForm     `json:"line_items"`
}

type CustomerDataApiForm struct {
	Checkout CustomerDataApiCheckoutForm `json:"checkout"`
}

type CustomerDataResponse struct {
	Checkout struct {
		CompletedAt         interface{} `json:"completed_at"`
		CreatedAt           string      `json:"created_at"`
		Currency            string      `json:"currency"`
		PresentmentCurrency string      `json:"presentment_currency"`
		CustomerID          int64       `json:"customer_id"`
		CustomerLocale      string      `json:"customer_locale"`
		DeviceID            interface{} `json:"device_id"`
		DiscountCode        interface{} `json:"discount_code"`
		Email               string      `json:"email"`
		LegalNoticeURL      interface{} `json:"legal_notice_url"`
		LocationID          interface{} `json:"location_id"`
		Name                string      `json:"name"`
		Note                string      `json:"note"`
		NoteAttributes      struct {
			CheckoutClicked string `json:"checkout_clicked"`
		} `json:"note_attributes"`
		OrderID                  interface{}   `json:"order_id"`
		OrderStatusURL           interface{}   `json:"order_status_url"`
		Order                    interface{}   `json:"order"`
		PaymentDue               string        `json:"payment_due"`
		PaymentURL               string        `json:"payment_url"`
		Payments                 []interface{} `json:"payments"`
		Phone                    interface{}   `json:"phone"`
		ShopifyPaymentsAccountID interface{}   `json:"shopify_payments_account_id"`
		PrivacyPolicyURL         interface{}   `json:"privacy_policy_url"`
		RefundPolicyURL          string        `json:"refund_policy_url"`
		RequiresShipping         bool          `json:"requires_shipping"`
		ReservationTimeLeft      int           `json:"reservation_time_left"`
		ReservationTime          interface{}   `json:"reservation_time"`
		SourceIdentifier         interface{}   `json:"source_identifier"`
		SourceName               string        `json:"source_name"`
		SourceURL                interface{}   `json:"source_url"`
		SubscriptionPolicyURL    interface{}   `json:"subscription_policy_url"`
		SubtotalPrice            string        `json:"subtotal_price"`
		ShippingPolicyURL        interface{}   `json:"shipping_policy_url"`
		TaxExempt                bool          `json:"tax_exempt"`
		TaxesIncluded            bool          `json:"taxes_included"`
		TermsOfSaleURL           interface{}   `json:"terms_of_sale_url"`
		TermsOfServiceURL        string        `json:"terms_of_service_url"`
		Token                    string        `json:"token"`
		TotalPrice               string        `json:"total_price"`
		TotalTax                 string        `json:"total_tax"`
		TotalTipReceived         string        `json:"total_tip_received"`
		TotalLineItemsPrice      string        `json:"total_line_items_price"`
		UpdatedAt                string        `json:"updated_at"`
		UserID                   interface{}   `json:"user_id"`
		WebURL                   string        `json:"web_url"`
		LineItems                []struct {
			ID               string      `json:"id"`
			Key              string      `json:"key"`
			ProductID        int64       `json:"product_id"`
			VariantID        int64       `json:"variant_id"`
			Sku              string      `json:"sku"`
			Vendor           string      `json:"vendor"`
			Title            string      `json:"title"`
			VariantTitle     string      `json:"variant_title"`
			ImageURL         string      `json:"image_url"`
			Taxable          bool        `json:"taxable"`
			RequiresShipping bool        `json:"requires_shipping"`
			GiftCard         bool        `json:"gift_card"`
			Price            string      `json:"price"`
			CompareAtPrice   interface{} `json:"compare_at_price"`
			LinePrice        string      `json:"line_price"`
			Properties       struct {
			} `json:"properties"`
			Quantity            int           `json:"quantity"`
			Grams               int           `json:"grams"`
			FulfillmentService  string        `json:"fulfillment_service"`
			AppliedDiscounts    []interface{} `json:"applied_discounts"`
			DiscountAllocations []interface{} `json:"discount_allocations"`
			TaxLines            []interface{} `json:"tax_lines"`
		} `json:"line_items"`
		GiftCards        []interface{} `json:"gift_cards"`
		TaxLines         []interface{} `json:"tax_lines"`
		TaxManipulations []interface{} `json:"tax_manipulations"`
		ShippingLine     interface{}   `json:"shipping_line"`
		ShippingRate     interface{}   `json:"shipping_rate"`
		ShippingAddress  struct {
			ID           int64       `json:"id"`
			FirstName    string      `json:"first_name"`
			LastName     string      `json:"last_name"`
			Phone        string      `json:"phone"`
			Company      interface{} `json:"company"`
			Address1     string      `json:"address1"`
			Address2     string      `json:"address2"`
			City         string      `json:"city"`
			Province     string      `json:"province"`
			ProvinceCode string      `json:"province_code"`
			Country      string      `json:"country"`
			CountryCode  string      `json:"country_code"`
			Zip          string      `json:"zip"`
		} `json:"shipping_address"`
		CreditCard     interface{} `json:"credit_card"`
		BillingAddress struct {
			ID           int64       `json:"id"`
			FirstName    string      `json:"first_name"`
			LastName     string      `json:"last_name"`
			Phone        string      `json:"phone"`
			Company      interface{} `json:"company"`
			Address1     string      `json:"address1"`
			Address2     string      `json:"address2"`
			City         string      `json:"city"`
			Province     string      `json:"province"`
			ProvinceCode string      `json:"province_code"`
			Country      string      `json:"country"`
			CountryCode  string      `json:"country_code"`
			Zip          string      `json:"zip"`
		} `json:"billing_address"`
		AppliedDiscount interface{} `json:"applied_discount"`
	} `json:"checkout"`
}

type CartApiResponse struct {
	Token      string      `json:"token"`
	Note       interface{} `json:"note"`
	Attributes struct {
	} `json:"attributes"`
	OriginalTotalPrice int     `json:"original_total_price"`
	TotalPrice         int     `json:"total_price"`
	TotalDiscount      int     `json:"total_discount"`
	TotalWeight        float64 `json:"total_weight"`
	ItemCount          int     `json:"item_count"`
	Items              []struct {
		ID         int64 `json:"id"`
		Properties struct {
			Upsell string `json:"upsell"`
		} `json:"properties"`
		Quantity                     int           `json:"quantity"`
		VariantID                    int64         `json:"variant_id"`
		Key                          string        `json:"key"`
		Title                        string        `json:"title"`
		Price                        int           `json:"price"`
		OriginalPrice                int           `json:"original_price"`
		DiscountedPrice              int           `json:"discounted_price"`
		LinePrice                    int           `json:"line_price"`
		OriginalLinePrice            int           `json:"original_line_price"`
		TotalDiscount                int           `json:"total_discount"`
		Discounts                    []interface{} `json:"discounts"`
		Sku                          string        `json:"sku"`
		Grams                        int           `json:"grams"`
		Vendor                       string        `json:"vendor"`
		Taxable                      bool          `json:"taxable"`
		ProductID                    int64         `json:"product_id"`
		ProductHasOnlyDefaultVariant bool          `json:"product_has_only_default_variant"`
		GiftCard                     bool          `json:"gift_card"`
		FinalPrice                   int           `json:"final_price"`
		FinalLinePrice               int           `json:"final_line_price"`
		URL                          string        `json:"url"`
		FeaturedImage                struct {
			AspectRatio float64 `json:"aspect_ratio"`
			Alt         string  `json:"alt"`
			Height      int     `json:"height"`
			URL         string  `json:"url"`
			Width       int     `json:"width"`
		} `json:"featured_image"`
		Image              string   `json:"image"`
		Handle             string   `json:"handle"`
		RequiresShipping   bool     `json:"requires_shipping"`
		ProductType        string   `json:"product_type"`
		ProductTitle       string   `json:"product_title"`
		ProductDescription string   `json:"product_description"`
		VariantTitle       string   `json:"variant_title"`
		VariantOptions     []string `json:"variant_options"`
		OptionsWithValues  []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"options_with_values"`
		LineLevelDiscountAllocations []interface{} `json:"line_level_discount_allocations"`
		LineLevelTotalDiscount       int           `json:"line_level_total_discount"`
	} `json:"items"`
	RequiresShipping              bool          `json:"requires_shipping"`
	Currency                      string        `json:"currency"`
	ItemsSubtotalPrice            int           `json:"items_subtotal_price"`
	CartLevelDiscountApplications []interface{} `json:"cart_level_discount_applications"`
}

type ShippingApiResponse struct {
	ShippingRates []struct {
		ID       string `json:"id"`
		Price    string `json:"price"`
		Title    string `json:"title"`
		Checkout struct {
			TotalTax      string `json:"total_tax"`
			TotalPrice    string `json:"total_price"`
			SubtotalPrice string `json:"subtotal_price"`
		} `json:"checkout"`
		PhoneRequired          bool        `json:"phone_required"`
		DeliveryRange          interface{} `json:"delivery_range"`
		EstimatedTimeInTransit interface{} `json:"estimated_time_in_transit"`
	} `json:"shipping_rates"`
}

type PaymentSessionCardData struct {
	Number            string `json:"number"`
	Name              string `json:"name"`
	Month             string `json:"month"`
	Year              string `json:"year"`
	VerificationValue string `json:"verification_value"`
}

type PaymentSessionPayload struct {
	CreditCard          PaymentSessionCardData `json:"credit_card"`
	PaymentSessionScope string                 `json:"payment_session_scope"`
}

type PaymentSessionResponse struct {
	Id string `json:"id"`
}

type CustomerCheckoutPayload struct {
	FirstName  string `json:"checkout[shipping_address][first_name]"`
	LastName   string `json:"checkout[shipping_address][last_name}"`
	Line1      string `json:"checkout[shipping_address][address1]"`
	Line2      string `json:"checkout[shipping_address][address2]"`
	City       string `json:"checkout[shipping_address][city]"`
	Country    string `json:"checkout[shipping_address][country]"`
	PostalCode string `json:"checkout[shipping_address][zip]"`
	Phone      string `json:"checkout[shipping_address][phone]"`
}

type QueueAvailability struct {
	Available bool `json:"available"`
}

type QPoll struct {
	__Typename                 string              `json:"__typename"`
	Token                      string              `json:"token"`
	QueueEtaSeconds            int                 `json:"queueEtaSeconds"`
	ProductVariantAvailability []QueueAvailability `json:"productVariantAvailability"`
}

type QData struct {
	Poll QPoll `json:"poll"`
}

type GetNextQueueResponse struct {
	Data QData `json:"data"`
}
