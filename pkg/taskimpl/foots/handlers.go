package foots

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/pkg/util"
	"nebula/third_party/adyen"
	"nebula/third_party/http"
	"strings"
	"time"
)

type SiteHeaderInfo struct {
	UserAgent     string
	ApiKey        string
	ApiIdentifier string
}

type FootsiteType string

const (
	FootsiteKFTL       FootsiteType = "kftl"
	FootsiteFTL        FootsiteType = "ftl"
	FootsiteEastBay    FootsiteType = "eastbay"
	FootsiteChamps     FootsiteType = "champs"
	FootsiteFootaction FootsiteType = "fta"
)

var (
	siteMap = map[FootsiteType]SiteHeaderInfo{
		FootsiteKFTL: {
			UserAgent:     "Kidsfootlocker/CFNetwork/Darwin",
			ApiKey:        "uiZXf5lQ4BCnOVMk3PBC4UmVh8Uxjcm4",
			ApiIdentifier: "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614",
		},
		FootsiteFTL: {
			UserAgent:     "FootLocker/CFNetwork/Darwin",
			ApiKey:        "m38t5V0ZmfTsRpKIiQlszub1Tx4FbnGG",
			ApiIdentifier: "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614",
		},
		FootsiteEastBay: {
			UserAgent:     "Eastbay/CFNetwork/Darwin",
			ApiKey:        "5IVJBiZI99JVraO5PAmsbxFCZnf7mmBx",
			ApiIdentifier: "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614",
		},
		FootsiteChamps: {
			UserAgent:     "ChampsSports/CFNetwork/Darwin",
			ApiKey:        "tFbOwo6SafnkT4BKbbbNh3sIlULmwSKn",
			ApiIdentifier: "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614",
		},
		FootsiteFootaction: {
			UserAgent:     "Footaction/CFNetwork/Darwin",
			ApiKey:        "Q0ZAbdoZoVsMuGGdjn6UO6DETP1ESgpk",
			ApiIdentifier: "921B2b33cAfba5WWcb0bc32d5ix89c6b0f613",
		},
	}
)

func handleSetup(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	getFootsiteType(runningCtx)
	setStore(runningCtx.FootsiteType, runningCtx)
	configureClient(taskCtx, runningCtx)

	return GetSession, nil
}

func handleGetSession(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Creating session")

	runningCtx.ProductId = staticCtx.Default.Sku

	getFootsiteType(runningCtx)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("https://bwi.pops.fastly-analytics.com/apigate/v3/session?timestamp=%d", time.Now().UnixNano()/int64(time.Millisecond))).
		Header("accept", "application/json").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("cache-control", "no-cache").
		Header("pragma", "no-cache").
		Header("cookie", makeCookies(runningCtx)).
		Header("referer", runningCtx.StoreUrl.String()).
		Header("user-agent", siteMap[runningCtx.FootsiteType].UserAgent).
		Header("x-api-key", siteMap[runningCtx.FootsiteType].ApiKey).
		Header("x-flapi-api-identifier", siteMap[runningCtx.FootsiteType].ApiIdentifier).
		Header("origin", runningCtx.StoreUrl.String()).
		Header("host", runningCtx.StoreUrl.Hostname()).
		Header("x-host", runningCtx.StoreUrl.Hostname()).
		Header("spoof-host", runningCtx.StoreUrl.Hostname()).
		Header("sec-fetch-dest", "empty").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-site", "same-origin").
		Header("x-fl-request-id", uuid.New().String()).
		SendAndClose()
	if err != nil {
		fmt.Println(err)
		taskCtx.SendStatusColored("Error creating session", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetSession, nil
	}

	extractCookies(resp.Header, runningCtx)

	// NOTE: MAX_CONN Fastly Error
	if resp.StatusCode == 503 {
		taskCtx.SendStatusColored("Max connections", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetSession, nil
	}

	if resp.StatusCode == 529 {
		return handleQueue(taskCtx, GetSession), nil
	}

	var body SessionResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Nil session, retrying..", util.ColorFailed)
		time.Sleep(time.Second * 1)
		return GetSession, nil
	}

	runningCtx.CsrfToken = body.Data.CsrfToken

	return GetStock, nil
}

func handleGetStock(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Retrieving stock")

	storeUrl := *runningCtx.StoreUrl
	storeUrl.Path = fmt.Sprintf("/product/~/%s.html", runningCtx.ProductId)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("https://bwi.pops.fastly-analytics.com/apigate/products/pdp/%s?timestamp=%d", runningCtx.ProductId, time.Now().UnixNano()/int64(time.Millisecond))).
		Header("accept", "application/json").
		Header("x-fl-request-id", uuid.New().String()).
		Header("x-flapi-session-id", runningCtx.SessionId).
		Header("cookie", makeCookies(runningCtx)).
		Header("user-agent", siteMap[runningCtx.FootsiteType].UserAgent).
		Header("x-api-key", siteMap[runningCtx.FootsiteType].ApiKey).
		Header("x-flapi-api-identifier", siteMap[runningCtx.FootsiteType].ApiIdentifier).
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("host", runningCtx.StoreUrl.Hostname()).
		Header("x-host", runningCtx.StoreUrl.Hostname()).
		Header("referer", storeUrl.String()).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("cache-control", "no-cache").
		Header("pragma", "no-cache").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error retrieving stock", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetStock, nil
	}

	extractCookies(resp.Header, runningCtx)

	// NOTE: MAX_CONN Fastly Error
	if resp.StatusCode == 503 {
		taskCtx.SendStatusColored("Max connections", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetStock, nil
	}

	if resp.StatusCode == 400 {
		taskCtx.SendStatusColored("Stock not live", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetStock, nil
	}

	if resp.StatusCode == 529 {
		return handleQueue(taskCtx, GetStock), nil
	}

	var stockResp StockResponse
	if err = resp.JsonBody(&stockResp); err != nil {
		taskCtx.SendStatusColored("Error parsing stock", util.ColorFailed)
		time.Sleep(time.Second * 1)
		return GetSession, nil
	}

	runningCtx.ProductName = stockResp.Name
	productUrl := *runningCtx.StoreUrl
	productUrl.Path = fmt.Sprintf("/product/~/%s.html", staticCtx.Default.Sku)
	runningCtx.ProductUrl = productUrl.String()

	style := getSkuCode(staticCtx, stockResp)
	if style == "" {
		taskCtx.SendStatusColored("No style matched, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetStock, nil
	}

	variants := extractStyleCodes(style, stockResp.SellableUnits)
	if len(variants) == 0 {
		taskCtx.SendStatusColored("No variations found, retrying...", util.ColorFailed)
		time.Sleep(time.Millisecond * 1500)
		return GetStock, nil
	}

	var vars []taskutil.Variant
	for _, variant := range variants {
		trimmed := taskutil.Variant{
			ID:      variant.Code,
			Size:    variant.Attributes[0].Value,
			InStock: false,
		}

		vars = append(vars, trimmed)
	}

	matches, err := taskutil.MatchVariants(staticCtx.Default.Sizes, vars, predicate)
	if err != nil {
		taskCtx.SendStatusColored("No variations matched, retrying...", util.ColorFailed)
		time.Sleep(time.Millisecond * 1500)
		return GetStock, nil
	}

	rand.Seed(time.Now().UnixNano())
	chosenVariant := matches[rand.Intn(len(matches))]
	runningCtx.Sku = chosenVariant.ID
	runningCtx.Size = chosenVariant.Size
	runningCtx.ProductImage = fmt.Sprintf("https://images.footlocker.com/pi/%s/small/%s.jpeg", staticCtx.Default.Sku, staticCtx.Default.Sku)

	return AddToCart, nil
}

func handleAddToCart(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Adding to cart")

	storeUrl := *runningCtx.StoreUrl
	storeUrl.Path = fmt.Sprintf("/product/~/%s.html", staticCtx.Default.Sku)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("https://acc.pops.fastly-analytics.com/apigate/users/carts/current/entries?timestamp=%d", time.Now().UnixNano()/int64(time.Millisecond))).
		Method(http.MethodPost).
		Header("accept", "application/json").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("cache-control", "no-cache").
		Header("cookie", makeCookies(runningCtx)).
		Header("content-type", "application/json").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("host", runningCtx.StoreUrl.Hostname()).
		Header("x-host", runningCtx.StoreUrl.Hostname()).
		Header("fastly-ff", "!!cache-bwi5152-BWI").
		Header("spoof-authority", fmt.Sprintf(":%s", runningCtx.StoreUrl.Hostname())).
		Header("referer", storeUrl.String()).
		Header("x-csrf-token", runningCtx.CsrfToken).
		Header("x-fl-productid", fmt.Sprintf("%s", runningCtx.Sku)).
		Header("x-fl-request-id", uuid.New().String()).
		Header("user-agent", siteMap[runningCtx.FootsiteType].UserAgent).
		Header("x-api-key", siteMap[runningCtx.FootsiteType].ApiKey).
		Header("x-flapi-api-identifier", siteMap[runningCtx.FootsiteType].ApiIdentifier).
		Header("x-flapi-session-id", runningCtx.SessionId).
		Header("pragma", "no-cache").
		JsonBody(map[string]interface{}{
			"productId":       runningCtx.Sku,
			"productQuantity": 1,
		}).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error adding to cart", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return AddToCart, nil
	}

	extractCookies(resp.Header, runningCtx)

	if resp.StatusCode == 503 { // NOTE: MAX_CONN Fastly Error
		taskCtx.SendStatusColored("Max connections", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return AddToCart, nil
	}

	if resp.StatusCode == 429 {
		taskCtx.SendStatusColored("Rate limited", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return AddToCart, nil
	}

	if resp.StatusCode == 403 {
		taskCtx.SendStatusColored("DataDome blocked", util.ColorFailed)
		time.Sleep(time.Second * 5)
		return AddToCart, nil
	}

	if resp.StatusCode == 400 || resp.StatusCode == 531 || resp.StatusCode == 550 {
		taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
		time.Sleep(time.Second * 1)
		return AddToCart, nil
	}

	if resp.StatusCode == 529 {
		return handleQueue(taskCtx, AddToCart), nil
	}

	var body AddToCartResponse
	if err = resp.JsonBody(&body); err != nil {
		fmt.Println(err)
		taskCtx.SendStatusColored("Error parsing cart", util.ColorWarning)
		time.Sleep(time.Second * 1)
		return AddToCart, nil
	}

	runningCtx.CartId = body.Guid
	runningCtx.ProductPrice = int64(body.TotalPrice.Value)

	return SubmitInformation, nil
}

func handleSubmitInformation(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting information")

	storeUrl := *runningCtx.StoreUrl
	storeUrl.Path = "/cart"

	// TODO: come back and fix this. This is a hacky way of splitting the name.
	nameParts := strings.Split(runningCtx.ProfileData.Billing.Name, " ")

	var shippingState string
	if runningCtx.ProfileData.Shipping.State != nil {
		shippingState = runningCtx.ProfileData.Shipping.State.Code
	}

	var billingState string
	if runningCtx.ProfileData.Billing.State != nil {
		billingState = runningCtx.ProfileData.Billing.State.Code
	}

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("https://acc.pops.fastly-analytics.com/apigate/users/carts/current/paypal?timestamp=%d", time.Now().UnixNano()/int64(time.Millisecond))).
		Method(http.MethodPost).
		Header("accept", "application/json").
		Header("x-csrf-token", runningCtx.CsrfToken).
		Header("x-fl-request-id", uuid.New().String()).
		Header("x-flapi-session-id", runningCtx.SessionId).
		Header("cookie", makeCookies(runningCtx)).
		Header("user-agent", siteMap[runningCtx.FootsiteType].UserAgent).
		Header("x-api-key", siteMap[runningCtx.FootsiteType].ApiKey).
		Header("x-flapi-api-identifier", siteMap[runningCtx.FootsiteType].ApiIdentifier).
		Header("x-flapi-cart-guid", runningCtx.CartId).
		Header("content-type", "application/json").
		Header("host", runningCtx.StoreUrl.Hostname()).
		Header("x-host", runningCtx.StoreUrl.Hostname()).
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", storeUrl.String()).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("cache-control", "no-cache").
		Header("pragma", "no-cache").
		JsonBody(map[string]interface{}{
			"checkoutType": "CREDIT_CARD",
			"type":         "CreditCard",
			"nonce":        uuid.New().String(),
			"details": map[string]interface{}{
				"email":     runningCtx.ProfileData.Payment.Email,
				"firstName": nameParts[0],
				"lastName":  nameParts[1],
				"payerId":   uuid.New().String(),
				"shippingAddress": map[string]interface{}{
					"recipientName":     runningCtx.ProfileData.Shipping.Name,
					"line1":             runningCtx.ProfileData.Shipping.Line1,
					"line2":             runningCtx.ProfileData.Shipping.Line2,
					"extendedAddress":   runningCtx.ProfileData.Shipping.Line2,
					"streetAddress":     runningCtx.ProfileData.Shipping.Line2,
					"city":              runningCtx.ProfileData.Shipping.City,
					"state":             shippingState,
					"postalCode":        runningCtx.ProfileData.Shipping.PostCode,
					"countryCode":       runningCtx.ProfileData.Shipping.Country.Code,
					"countryCodeAlpha2": runningCtx.ProfileData.Shipping.Country.Code,
					"locality":          runningCtx.ProfileData.Shipping.City,
					"region":            shippingState,
				},
				"phone":       runningCtx.ProfileData.Payment.Phone,
				"countryCode": runningCtx.ProfileData.Shipping.Country.Code,
				"billingAddress": map[string]interface{}{
					"recipientName":     runningCtx.ProfileData.Billing.Name,
					"line1":             runningCtx.ProfileData.Billing.Line1,
					"line2":             runningCtx.ProfileData.Billing.Line2,
					"extendedAddress":   runningCtx.ProfileData.Billing.Line2,
					"streetAddress":     runningCtx.ProfileData.Billing.Line2,
					"city":              runningCtx.ProfileData.Billing.City,
					"state":             billingState,
					"postalCode":        runningCtx.ProfileData.Billing.PostCode,
					"countryCode":       runningCtx.ProfileData.Billing.Country.Code,
					"countryCodeAlpha2": runningCtx.ProfileData.Billing.Country.Code,
					"locality":          runningCtx.ProfileData.Billing.City,
					"region":            billingState,
				},
			},
		}).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error submitting information", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitInformation, nil
	}

	extractCookies(resp.Header, runningCtx)

	if resp.StatusCode == 503 {
		taskCtx.SendStatusColored("Max connections", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitInformation, nil
	}

	if resp.StatusCode == 403 {
		taskCtx.SendStatusColored("DataDome blocked", util.ColorFailed)
		time.Sleep(time.Second * 5)
		return SubmitInformation, nil
	}

	if resp.StatusCode == 529 {
		return handleQueue(taskCtx, SubmitInformation), nil
	}

	return SubmitOrder, nil
}

func handleSubmitOrder(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting order")

	storeUrl := *runningCtx.StoreUrl
	storeUrl.Path = "/checkout"

	a := adyen.NewAdyen("adyenjs_0_1_18$", "A237060180D24CDEF3E4E27D828BDB6A13E12C6959820770D7F2C1671DD0AEF4729670C20C6C5967C664D18955058B69549FBE8BF3609EF64832D7C033008A818700A9B0458641C5824F5FCBB9FF83D5A83EBDF079E73B81ACA9CA52FDBCAD7CD9D6A337A4511759FA21E34CD166B9BABD512DB7B2293C0FE48B97CAB3DE8F6F1A8E49C08D23A98E986B8A995A8F382220F06338622631435736FA064AEAC5BD223BAF42AF2B66F1FEA34EF3C297F09C10B364B994EA287A5602ACF153D0B4B09A604B987397684D19DBC5E6FE7E4FFE72390D28D6E21CA3391FA3CAADAD80A729FEF4823F6BE9711D4D51BF4DFCB6A3607686B34ACCE18329D415350FD0654D")
	info, err := a.EncryptCreditcardDetails(runningCtx.ProfileData.Payment.Number, runningCtx.ProfileData.Payment.ExpMonth, runningCtx.ProfileData.Payment.ExpYear, runningCtx.ProfileData.Payment.Cvv)
	if err != nil {
		return "", err
	}

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("https://acc.pops.fastly-analytics.com/apigate/v2/users/orders?timestamp=%d", time.Now().UnixNano()/int64(time.Millisecond))).
		Method(http.MethodPost).
		Header("accept", "application/json").
		Header("x-csrf-token", runningCtx.CsrfToken).
		Header("x-fl-request-id", uuid.New().String()).
		Header("x-flapi-cart-guid", uuid.New().String()).
		Header("x-flapi-session-id", runningCtx.SessionId).
		Header("cookie", makeCookies(runningCtx)).
		Header("user-agent", siteMap[runningCtx.FootsiteType].UserAgent).
		Header("x-api-key", siteMap[runningCtx.FootsiteType].ApiKey).
		Header("x-flapi-api-identifier", siteMap[runningCtx.FootsiteType].ApiIdentifier).
		Header("content-type", "application/json").
		Header("origin", storeUrl.String()).
		Header("host", runningCtx.StoreUrl.Hostname()).
		Header("x-host", runningCtx.StoreUrl.Hostname()).
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", storeUrl.String()).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("cache-control", "no-cache").
		Header("pragma", "no-cache").
		JsonBody(map[string]interface{}{
			"paymentMethod":         "CREDITCARD",
			"sid":                   "",
			"deviceId":              "0500KYdAroY0whrC8CDqltcUm6G/OehQcmsBicaUaq5smOXApMumK1A4sp9IXhSIOhMWjdv1ut3EqNWozlTHv/VXZHerdMnBmN7FDxhtgBGnUGnwCNfUnz83BQZy2VEREPPtB2tdqu9QgXctIbpYIWRHeNXLvZyqmHpQEIh47Xlfmjq7JroGyb2hi9Hk0CazCu+I+18utShuIOzjkoAwX9CVbiPZK21t+6Xqa0g9NiKg9p9eSinro3l4EVTc2ULntQro4HXsBEgPqLA6JcVeGyREQCrWPnyiT8GngRQ5qG2EsZ+C+vRwHXBY2D29Ow0HVk+WefLZ5t0a74RqJFQweDAygvJty85J5zrR1l1Ue3dKFQBuqIW6GZXhealXeZ8YG/vOg/phdtdt48VTI+/GqNjdICGX8jm9a+ivKrxXNEWRDbK85dG3B2yRW+0Ucb3tJ6QbpVYuAE/dX9/drvCj0k31Gps1k6MFzSLA84KdBHUZnitnqJzvUp51MCr/CcHZslSAjd0oCqHsgZ3QjfmAv0cFwSoBfTSMV1H8DlHGg8IdmH3+5aRbcPY4PysSMhw6FbSUO7XQQPeydJU1vw7/JNIducR+aDpfrVjGpFHd7qaDffmceQTLQZUhR1krSG/GTwP6VGBDUOUglz3PrcBBq9tVF5a50I5j26NCgwX/i3mOCBBiKsekTD5BL4gIusbuZShopBvaQkomm/xi5lMsoj98x/lHBL6KJXSz3f9B+1x3Mv2lzSd8vm88b1+iz3EPsRqALjOfcn++yQsgk9tUZdX+jVQGvCkFcFi2aht1s3lWxdD7H2vEZQH85tyZm0t1ztxpqlqXCPsLcOBKndH+/FDaO1m301CqnwwH0lo/42h+u/JvSyYWtkMpLjTVSw7ObY7taiORDka3KFX1h/MYnGlRtWKxO0vz5730jb0p2H9vBtIwcggZU2+JB83fdLZFaBtVRIfqw8y/uxIo0zVWQGiv+zPlfiYTiYE+QPN5VMzukkqUOu+8pcJsPt1mc+nHzrQO1EVm+4Fi3G73RjsTxTWchPGIW184yEItl4kXnufWcZqpE3GYGg6XJI/27Q5fEHzTnVA/kiLCz7Bm+1VVXAtMMv4axfYWz/nVqBvqpOlsGbd+KfJeP6ChUm+SSAY37/sLGUFeZL4kbXv5m4q4+Sa7hFcXxmZggYLCLkZmH9bgNRfObl2nwqvFZhIEq93jpUm6GPOxrl5QujJtjkNx0lDFpaJzzrazBerj0zvKDJvnYESggcWFWOx2f3lqm0fHlUdiTfTVE0aYHVc/DExrH6wsYUVAUkxA1OGUjN5EbxwhL0mMWZF/pIbbqpUSUM6AA1+QeYBI+cNU9qgFc8ZhEmGlWjwX5kEDn/8AWxK/4HZt9KkQLhMWoe0dPgkZvLPdhuolp/uHwMyX0AxqTHMqIoratsdNmmBhjVM9rdodKMtpNEppTkYWScjTEvgMjjKbWXE1BRxiYq9sUNcLKqYLcs2bA7pfcpoq60dhgHnvwyP9KsLsF5VlYNCoHAxaFEFf8xzJkTGgDJSTBT7SBYa2LO4kWUHPuByrIPFJQGj3oTUwF8LqY9KmXCPFjT8rxdhWLogtfwG/s8lW+uFfmOEeWAx4eIg0aDl0T0N/HJMQBM5V715iP25Jj9HwFQdBrcWlkm0ZiSRAYOqizv5ZXkmrmA6qRMEJ6XLqB5lZJKIuoFoQ2tl5SgML4c1VHB7Kxb40LS25zOP6QxAV39AIb2i2I7a2p+iv7QdjVztcmlUgHvx4xg6jCr3hXFrSQNL8mBc1rKQ42Uqj4q1mW3DKLcQ5EV0fs7sJvoJ/pPrrGyMmdOWkDfChdYnG6xpnpdKqHR9XuhBnzlFb+qT5uLyUC8qfI0xOppJjJ6v555ytySTeh4PLCqagrEMPfMMX2FQ0A7KKLhPrpPfYMUJMlwVC8QRAiqA4bFmTvs2EOhBrE/FhOVTfXp17wLA/yVJnx0/7Qf/qfQv7UhKCjnHKR+OV1bRPeBqPvKAGjwJyg/MflRtyHECCda77jgZoMuqbBBmiQzLYAWM4AYXs2bNXxxLW2DaaQQ4KN10zW2b3dzDDENFfdhUeMDsrrMj9EHjrjLHRF2GXiAwtCZJmUD+bX0LDHc4pwys0lzBuOARqFsQTwMjdfoVjcvRcQK50HYiFre2u9hn9FVzrb1zyl+6KBb5CfR5Mv8vUgLf8eqTwtYBXSyg6mHAzM2BmNNkzA3dQ7Zmdq0WX2Bp2nJ31HCWGqmkd3bswxyTX9q4cqwZrJpPDW8Ek0hjjJ5jk1pOTpE3TO0aB/r2KA+HaBKTl/f1keFRDcXzg5CbcyhOyJTjpNyNMI+9OzZzMHWVX2Bm9ScTiwn4eDx3L3KChZjM8W+uhnTY1ujNageoOJZHF9sbMJjodPjFym+l/4Fx5LGia28Ro1SEc1ipybNbfDUMO84bWqFutP2EYHh5DcUp302xvkpGgcngDu2btzr+rrsCshJWMrnewcOeFDPZ/f62tLkuBMqVGbg0z7ibAXMi1KcyxNnI/ZNxwLgIfvjJ8sqFWY7U9tbGr3k3ZUjxgmgsvmJo73PCLifWo0LsOjaalWaF03DXe9RR0T2r6fOeoKs5393fb+DKzNLGKu2Zs3j26rvk6jCo4Ek1rQGIera5tE1C6wFnFGczEiZ6FrhpLA4zbbnu+njkBmrUjRwv+aRoepB2t8ZYIBmxe8Oh61Gf44tzNWeoUV/YukZbeZKa1j4ygI8TPVzGAOxPYT4jC1ltRivhZaRl9pYn2KCeprtKM8ww+Zruxcz6dnzNpqPyYiyjZjPYYCGK2KESSBsgSwJqiSNOVzAEPG+hUyc3FpRg2+PKFnJ4iqN9d5IOf+itCRVxRHNwcMrWjAWilvrI7AfXTvYJkJpXPAgRpmqY7BhGcfMRmiPWqDH7bj3XrYd/lOgZQ8aFaHmnHhBnB7rYD5NhQo1KkBHLI3lOg7aV18xg6Wi8GiBcyjgXaax08ZAJBEpptH0/o9A21bOcjiFLCh9+MchS4ZuOj26twFNxm5NRYYCT3277WYavMnNCBPoftQcqmm2HGLkxn6IdHN/1V2Q/whIF2ujGaL9DI4x6y1RJBSVj9Xggo4qYcp/OmkM0JmjHLm4c9I3461SvIZNjh422tq1qUJhh1dlmVfdj4tVs7RfngQm243+nefQIp90HkeeVb5DQ/aReNloEoIOKLEOb/kzvCJHKlfW7b8dDtEEv/R/RchWM2clOJSJxD8iZ1Cfyx2L91D6pfGOPkRiA2",
			"cartId":                runningCtx.CartId,
			"encryptedCardNumber":   info.EncryptedCardNumber,
			"encryptedExpiryMonth":  info.EncryptedExpMonth,
			"encryptedExpiryYear":   info.EncryptedExpYear,
			"encryptedSecurityCode": info.EncryptedCvc,
		}).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error submitting order", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitOrder, nil
	}

	extractCookies(resp.Header, runningCtx)


	if resp.StatusCode == 503 {
		taskCtx.SendStatusColored("Max connections", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitOrder, nil
	}

	if resp.StatusCode == 403 {
		taskCtx.SendStatusColored("DataDome blocked", util.ColorFailed)
		time.Sleep(time.Second * 5)
		return SubmitOrder, nil
	}

	if resp.StatusCode == 529 {
		return handleQueue(taskCtx, SubmitOrder), nil
	}

	if resp.StatusCode == 400 {
		addWebhook(taskCtx, runningCtx,  "failed")
		taskCtx.SendStatusColored("Order failed", util.ColorFailed)
		return "", tasks.CompletedErr
	}

	var body SuccessfulCheckout
	if err = resp.JsonBody(&body); err == nil {
		runningCtx.OrderId = body.Order.Code
		runningCtx.ProductPrice = int64(body.Order.TotalPrice.Value)
	}

	addWebhook(taskCtx, runningCtx,  "placed")
	taskCtx.SendStatusColored("Order placed", util.ColorSuccess)
	return "", tasks.CompletedErr
}
