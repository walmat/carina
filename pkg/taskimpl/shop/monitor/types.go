package monitor


type Result struct {
	Product *StoreProduct
	Error error
}

type StoreProduct struct {
	ProductName string
	ImageUrl string
	Sku string
	Variants    interface{}
}

type StoresPool struct {
	ProductsByStore map[string]map[string]*StoreProduct
}

type ProductResponseVariant struct {
	ID                  int64       `json:"id"`
	Title               string      `json:"title"`
	Option1             string      `json:"option1"`
	Option2             interface{} `json:"option2"`
	Option3             interface{} `json:"option3"`
	Sku                 string      `json:"sku"`
	RequiresShipping    bool        `json:"requires_shipping"`
	Taxable             bool        `json:"taxable"`
	FeaturedImage       interface{} `json:"featured_image"`
	Available           bool        `json:"available"`
	Name                string      `json:"name"`
	PublicTitle         string      `json:"public_title"`
	Options             []string    `json:"options"`
	Price               int         `json:"price"`
	Weight              int         `json:"weight"`
	CompareAtPrice      interface{} `json:"compare_at_price"`
	InventoryManagement string      `json:"inventory_management"`
	Barcode             interface{} `json:"barcode"`
}

type ProductResponseJs struct {
	ID                   int64       `json:"id"`
	Title                string      `json:"title"`
	Handle               string      `json:"handle"`
	Description          string      `json:"description"`
	PublishedAt          string      `json:"published_at"`
	CreatedAt            string      `json:"created_at"`
	Vendor               string      `json:"vendor"`
	Type                 string      `json:"type"`
	Tags                 []string    `json:"tags"`
	Price                int         `json:"price"`
	PriceMin             int         `json:"price_min"`
	PriceMax             int         `json:"price_max"`
	Available            bool        `json:"available"`
	PriceVaries          bool        `json:"price_varies"`
	CompareAtPrice       interface{} `json:"compare_at_price"`
	CompareAtPriceMin    int         `json:"compare_at_price_min"`
	CompareAtPriceMax    int         `json:"compare_at_price_max"`
	CompareAtPriceVaries bool        `json:"compare_at_price_varies"`
	Variants             []ProductResponseVariant `json:"variants"`
	Images        []string `json:"images"`
	FeaturedImage string   `json:"featured_image"`
	Options       []struct {
		Name     string   `json:"name"`
		Position int      `json:"position"`
		Values   []string `json:"values"`
	} `json:"options"`
	URL   string `json:"url"`
	Media []struct {
		Alt          interface{} `json:"alt"`
		ID           int64       `json:"id"`
		Position     int         `json:"position"`
		PreviewImage struct {
			AspectRatio float64 `json:"aspect_ratio"`
			Height      int     `json:"height"`
			Width       int     `json:"width"`
			Src         string  `json:"src"`
		} `json:"preview_image"`
		AspectRatio float64 `json:"aspect_ratio"`
		Height      int     `json:"height"`
		MediaType   string  `json:"media_type"`
		Src         string  `json:"src"`
		Width       int     `json:"width"`
	} `json:"media"`
}

type GetJsonProductsVariant struct {
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

type GetJsonProductsResponse struct {
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
		Variants    []GetJsonProductsVariant `json:"variants"`
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