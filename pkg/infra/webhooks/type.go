package webhooks

import "encoding/json"

func UnmarshalWebhookData(data []byte) (WebhookData, error) {
	var webhook WebhookData
	err := json.Unmarshal(data, &webhook)
	return webhook, err
}

type OrderStatus string
type WebhookField string

const (
	OrderStatusPlaced OrderStatus = "Order placed"
	OrderStatusFailed OrderStatus = "Order failed"
	ColorStatusPlaced OrderStatus = "placed"
	ColorStatusFailed OrderStatus = "failed"

	TaskField     WebhookField = "Task"
	ProductField  WebhookField = "Product"
	StoreField    WebhookField = "Store"
	ProfileField  WebhookField = "Profile"
	ProxyField    WebhookField = "Proxy"
	SizeField     WebhookField = "Size"
	QuantityField WebhookField = "Quantity"
	OrderField    WebhookField = "Order"
	EmailField    WebhookField = "Email"
)

type Product struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	Image    string `json:"image"`
	Size     string `json:"size"`
	Url      string `json:"url"`
	Quantity string `json:"quantity"`
}

type Store struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProfileEmbed struct {
	ID        string `json:"profile"`
	GroupName string `json:"groupName"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

type GroupedObject struct {
	GroupName string `json:"groupName"`
	Name      string `json:"name"`
}

type Order struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type Task struct {
	ID        string `json:"id"`
	Mode      string `json:"mode"`
	GroupName string `json:"groupName"`
}

type CheckoutData struct {
	Task    Task          `json:"task"`
	Product Product       `json:"product"`
	Store   Store         `json:"store"`
	Profile ProfileEmbed  `json:"profile"`
	Proxy   GroupedObject `json:"proxy"`
	Order   Order         `json:"order"`
	Status  OrderStatus   `json:"status"`
}

type Field struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type WebhookData struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Active      bool           `json:"active"`
	Declines    bool           `json:"declines"`
	Sensitivity bool           `json:"sensitivity"`
	Profiles    []Profile      `json:"profiles"`
	Fields      []WebhookField `json:"fields"`
	Url         string         `json:"url"`
	Type        string         `json:"type"`
}
