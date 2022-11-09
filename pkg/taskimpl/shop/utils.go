package shopify

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/avct/uasurfer"
	"github.com/elliotchance/orderedmap"
	"io"
	"log"
	"nebula/pkg/infra/profiles"
	"nebula/third_party/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	checkoutHashNotFound    = errors.New("checkout hash not found")
	paymentGatewaysByShopId = map[string]string{
		"6269065":     "26102467",
		"499112":      "39676444746",
		"4356735045":  "24945229893",
		"517042":      "5021317",
		"14644082":    "118964878",
		"26792329325": "37105762413",
		"942252":      "128707719",
		"6163517":     "8308339",
		"6240605":     "8410639",
		"12109358":    "61989958",
		"2085268":     "3723727",
		"2681785":     "3953145",
		"8523376":     "39578599470",
		"2192362":     "6735901",
		"2825850":     "5730877",
		"10116760":    "35379780",
		"1460732":     "19494852",
		"18048548964": "38994542692",
		"16141961":    "113981520",
	}
)

/*
	Builds the `sec-ch-ua` header based on the used User-Agent

	should return => `sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`
*/
func buildSecurityHeader(userAgent string) string {
	ua := uasurfer.Parse(userAgent)

	return fmt.Sprintf("\"Chromium\";v=\"%d\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"%d\"", ua.Browser.Version.Major, ua.Browser.Version.Major)
}

// toggleRedirect Makes the client follow or not follow a redirect
func toggleRedirect(followRedirect bool, ctx *runningCtx) {
	if followRedirect {
		ctx.Client.Inner().CheckRedirect = nil
		return
	}
	ctx.Client.Inner().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

// isImproperStatusCode
func isImproperStatusCode(statusCode int) bool {
	invalidStatusCode := regexp.MustCompile("^[1-5][0-9][0-9]$")

	if invalidStatusCode.MatchString(strconv.FormatInt(int64(statusCode), 10)) {
		return true
	}
	return false
}

// is this function needed? the task can get the product handle and the product name from the monitor response
func retrieveCartInfo(ctx *runningCtx) error {
	resp, err := ctx.Client.Builder(fmt.Sprintf("%s/cart.js", ctx.StoreUrl.String())).
		Method(http.MethodGet).
		Header("content-type", "application/json").
		SendAndClose()

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("invalid status code")
	}

	var body CartApiResponse
	if err = resp.JsonBody(&body); err != nil {
		return errors.New("error retrieving cart information")
	}

	if len(body.Items) < 1 {
		return errors.New("cart empty")
	}

	return nil
}

func resetRunningCtx(ctx *runningCtx) {
	ctx.RetryCount = 0
	ctx.Restocking = false
	ctx.AuthToken = ""
	ctx.IsFree = false
	ctx.PaymentSessionId = ""
	ctx.Checked = false
	ctx.PollingUrl = false
	ctx.Fallback = false
}

// extractCheckoutHash Extracts the checkout hash found in the checkout url
func extractCheckoutHash(checkoutUrl string) (string, error) {
	checkoutHashRegexp := regexp.MustCompile(".*?checkouts\\/([0-9a-z]*)")

	if !checkoutHashRegexp.MatchString(checkoutUrl) {
		return "", checkoutHashNotFound
	}
	matches := checkoutHashRegexp.FindAllStringSubmatch(checkoutUrl, -1)

	if len(matches) < 1 {
		return "", checkoutHashNotFound
	}

	return matches[0][1], nil
}

func extractAuthToken(responseBody string) string {
	r := regexp.MustCompile(`<input\stype="hidden"\sname="authenticity_token"\svalue="(.+?)"`)
	authToken := r.FindAllStringSubmatch(responseBody, -1)[0][1]
	return authToken
}

func extractCaptchaData(responseBody string) string {
	r := regexp.MustCompile(`.*<noscript>.*<iframe\s.*src=.*\?k=(.*)"><\/iframe>`)
	sitekey := r.FindAllStringSubmatch(responseBody, -1)[0][1]
	return sitekey
}

func extractPaymentGateway(responseBody io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		return "", err
	}

	gateway, exists := doc.Find(`[data-gateway-name*="card"] input`).Attr("value")
	if !exists {
		return "", errors.New("could not find payment gateway")
	}

	return gateway, nil
}

func extractPrice(responseBody io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		return "", err
	}

	price, exists := doc.Find(`span.total-recap__final-price`).Attr("data-checkout-payment-due-target")
	if !exists {
		price, exists := doc.Find(`span.payment-due__price`).Attr("data-checkout-payment-due-target")
		if !exists {
			return "", errors.New("could not find price")
		}
		return price, nil
	}

	return price, nil
}

func extractProductFinalPrice(responseBody string) string {
	r := regexp.MustCompile(`<input\stype="hidden"\sname="checkout\[total_price]"\sid="checkout_total_price"\svalue="(.+?)"`)
	productPrice := r.FindAllStringSubmatch(responseBody, -1)[0][1]
	return productPrice
}

func extractProtection(responseBody io.Reader, query string, runningCtx *runningCtx) (*orderedmap.OrderedMap, error) {
	extractedForm := orderedmap.NewOrderedMap()
	var extractedBpHashes []string

	doc, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		return nil, err
	}

	doc.Find(query).Find("input").Each(func(i int, selection *goquery.Selection) {
		attrName, _ := selection.Attr("name")

		if (attrName == "field_start" ||
			attrName == "field_end") ||
			strings.Contains(attrName, "count") ||
			strings.Contains(attrName, "hosted_fields_redirect") ||
			strings.Contains(attrName, "checkout[billing_address]") {
			return
		}

		//log.Println(attrName)
		extractedForm.Set(attrName, "")
	})

	fsFormSelector := fmt.Sprintf(`#fs_%s`, runningCtx.CheckoutHash)

	doc.Find(fsFormSelector).Find("textarea").Each(func(i int, selection *goquery.Selection) {
		attrName, _ := selection.Attr("name")
		extractedBpHashes = append(extractedBpHashes, attrName)
	})
	// note: the fscount value counters always go after the bp hashes
	// if any hashes are found they should be added to the extracted form
	if len(extractedBpHashes) >= 1 {
		fsCountKey := fmt.Sprintf("%s-count", runningCtx.CheckoutHash)
		for _, h := range extractedBpHashes {
			extractedForm.Set(h, "")
		}
		// The same key needs to be added twice, one with the total amount of hashes and the other with "fs-count" as value
		extractedForm.Set(fsCountKey, fmt.Sprintf("%d", len(extractedBpHashes)))
		// todo add total length here and parse it in the loops
		extractedForm.Set(fsCountKey+"_REPLACE", "")
	}

	return extractedForm, nil
}

/* Safe and preload customer forms */
func getBrowserData(captchaToken string) *orderedmap.OrderedMap {
	form := orderedmap.NewOrderedMap()

	if captchaToken != "" {
		form.Set("g-recaptcha-response", captchaToken)
	}

	form.Set("checkout[client_details][browser_width]", "1920")
	form.Set("checkout[client_details][browser_height]", "1080")
	form.Set("checkout[client_details][javascript_enabled]", "1")
	form.Set("checkout[client_details][color_depth]", "24")
	form.Set("checkout[client_details][java_enabled]", "false")
	form.Set("checkout[client_details][browser_tz]", "240")

	return form
}

func getCustomerPayloadData(isEmpty bool, profile *profiles.ProfileData) *orderedmap.OrderedMap {
	form := orderedmap.NewOrderedMap()

	if isEmpty {
		form.Set("checkout[shipping_address][first_name]", "")
		form.Set("checkout[shipping_address][last_name]", "")
		form.Set("checkout[shipping_address][address1]", "")
		form.Set("checkout[shipping_address][address2]", "")
		form.Set("checkout[shipping_address][city]", "")
		form.Set("checkout[shipping_address][country]", "")
		form.Set("checkout[shipping_address][zip]", "")
		form.Set("checkout[shipping_address][phone]", "")
		return form
	}
	// TODO: Make a better way to get the name and last name
	nameParts := strings.Split(profile.Billing.Name, " ")

	form.Set("checkout[shipping_address][first_name]", nameParts[0])
	form.Set("checkout[shipping_address][last_name]", nameParts[1])
	form.Set("checkout[shipping_address][address1]", profile.Shipping.Line1)
	form.Set("checkout[shipping_address][address2]", profile.Shipping.Line2)
	form.Set("checkout[shipping_address][city]", profile.Shipping.City)
	form.Set("checkout[shipping_address][country]", profile.Shipping.Country.Name)
	form.Set("checkout[shipping_address][zip]", profile.Shipping.PostCode)
	form.Set("checkout[shipping_address][phone]", profile.Payment.Phone)

	return form
}

func generateSafeCustomerForm(runningCtx *runningCtx) string {
	// is there an extra payload for every site?
	emptyPayload := getCustomerPayloadData(false, runningCtx.ProfileData)
	profilePayload := getCustomerPayloadData(false, runningCtx.ProfileData)

	body := ""

	// Base data the payload will always have
	baseData := orderedmap.NewOrderedMap()
	baseData.Set("_method", "patch")
	baseData.Set("authenticity_token", runningCtx.AuthToken)
	baseData.Set("previous_step", "contact_information")
	baseData.Set("step", "shipping_method")
	baseData.Set("checkout[email]", runningCtx.ProfileData.Payment.Email)
	baseData.Set("checkout[buyer_accepts_marketing]", "0")

	browserData := getBrowserData(runningCtx.CaptchaToken)

	// Since shopify bot protection detects the order of the payload we have to add everything in order

	// first add the base data
	for _, b := range baseData.Keys() {
		value, exists := baseData.Get(b)
		if exists {
			body += url.QueryEscape(b.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	// add the first customer payload struct (do all sites want this payload empty?)
	for _, k := range emptyPayload.Keys() {
		value, exists := emptyPayload.Get(k)
		if exists {
			body += url.QueryEscape(k.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	// add the customer payload struct with all the profile data
	for _, k := range profilePayload.Keys() {
		value, exists := profilePayload.Get(k)
		if exists {
			body += url.QueryEscape(k.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	// add all the bot protection hashes, maybe we need to check for other keys that might not be protection hashes
	for _, k := range runningCtx.Form.Keys() {
		notMethod := k != "_method"
		notCaptcha := k != "g-recaptcha-response"
		notStep := k != "step"
		notPreviousStep := k != "previous_step"
		notAuthToken := k != "authenticity_token"

		if notMethod && notCaptcha && notStep && notPreviousStep && notAuthToken {
			if !strings.Contains(k.(string), "checkout[") {
				fsCountKey := fmt.Sprintf("%s-count", runningCtx.CheckoutHash)
				// if the current key is fs-count then add its value
				if k == fsCountKey {
					value, exists := runningCtx.Form.Get(k)
					if exists {
						body += url.QueryEscape(k.(string)) + "=" + value.(string) + "&"
						continue
					}
				}

				if k == fsCountKey+"_REPLACE" {
					body += url.QueryEscape(fsCountKey) + "=" + url.QueryEscape("fs_count") + "&"
					continue
				}
				body += url.QueryEscape(k.(string)) + "=" + "&"
			}
		}
	}

	// finally add the the browser data
	for _, browserDataKey := range browserData.Keys() {
		value, exists := browserData.Get(browserDataKey)
		if exists {
			if browserDataKey == "checkout[client_details][browser_tz]" {
				body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string))
				continue
			}
			body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}
	//log.Println(body)
	//for _, i2 := range strings.Split(body, "&") {
	//	log.Println(i2)
	//}
	return body
}

func generateSafeShippingForm(runningCtx *runningCtx) string {
	body := ""

	baseData := orderedmap.NewOrderedMap()
	baseData.Set("_method", "patch")
	baseData.Set("authenticity_token", runningCtx.AuthToken)
	baseData.Set("previous_step", "shipping_method")
	baseData.Set("step", "payment_method")
	baseData.Set("checkout[shipping_rate][id]", runningCtx.ShippingRate)

	// first add the base data
	for _, b := range baseData.Keys() {
		value, exists := baseData.Get(b)
		if exists {
			body += url.QueryEscape(b.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	// add all the bot protection hashes, maybe we need to check for other keys that might not be protection hashes
	for _, k := range runningCtx.Form.Keys() {
		notMethod := k != "_method"
		notShippingRate := k != "checkout[shipping_rate][id]"
		notStep := k != "step"
		notPreviousStep := k != "previous_step"
		notAuthToken := k != "authenticity_token"

		if notMethod && notShippingRate && notStep && notPreviousStep && notAuthToken {
			if !strings.Contains(k.(string), "checkout[") {
				fsCountKey := fmt.Sprintf("%s-count", runningCtx.CheckoutHash)
				// if the current key is fs-count then add its value
				if k == fsCountKey {
					value, exists := runningCtx.Form.Get(k)
					if exists {
						body += url.QueryEscape(k.(string)) + "=" + value.(string) + "&"
						continue
					}
				}

				if k == fsCountKey+"_REPLACE" {
					body += url.QueryEscape(fsCountKey) + "=" + url.QueryEscape("fs_count") + "&"
					continue
				}
				body += url.QueryEscape(k.(string)) + "=" + "&"
			}
		}
	}

	browserData := getBrowserData("")

	// finally add the the browser data
	for _, browserDataKey := range browserData.Keys() {
		value, exists := browserData.Get(browserDataKey)
		if exists {
			if browserDataKey == "checkout[client_details][browser_tz]" {
				body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string))
				continue
			}
			body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	return body
}

func generateSafePaymentForm(runningCtx *runningCtx) string {
	body := ""

	baseData := orderedmap.NewOrderedMap()
	baseData.Set("_method", "patch")
	baseData.Set("authenticity_token", runningCtx.AuthToken)
	baseData.Set("previous_step", "payment_method")
	baseData.Set("step", "")
	baseData.Set("s", runningCtx.PaymentSessionId)

	for _, b := range baseData.Keys() {
		value, exists := baseData.Get(b)
		if exists {
			body += url.QueryEscape(b.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	for _, k := range runningCtx.Form.Keys() {
		notMethod := k != "_method"
		notShippingRate := k != "s"
		notStep := k != "step"
		notPreviousStep := k != "previous_step"
		notAuthToken := k != "authenticity_token"

		if notMethod && notShippingRate && notStep && notPreviousStep && notAuthToken {
			if !strings.Contains(k.(string), "checkout[") && !strings.Contains(k.(string), "complete") {
				fsCountKey := fmt.Sprintf("%s-count", runningCtx.CheckoutHash)
				// if the current key is fs-count then add its value
				if k == fsCountKey {
					value, exists := runningCtx.Form.Get(k)
					if exists {
						body += url.QueryEscape(k.(string)) + "=" + value.(string) + "&"
						continue
					}
				}

				if k == fsCountKey+"_REPLACE" {
					body += url.QueryEscape(fsCountKey) + "=" + url.QueryEscape("fs_count") + "&"
					continue
				}
				body += url.QueryEscape(k.(string)) + "=" + "&"
			}
		}
	}

	paymentData := orderedmap.NewOrderedMap()
	isSameBillingAndShipping := false

	paymentData.Set("checkout[payment_gateway]", runningCtx.PaymentGateway)
	paymentData.Set("checkout[credit_card][vault]", "false")

	if strings.Contains(runningCtx.ProfileData.Shipping.Line1, runningCtx.ProfileData.Billing.Line1) {
		isSameBillingAndShipping = true
		paymentData.Set("checkout[different_billing_address]", "false")
	}

	if !isSameBillingAndShipping {
		paymentData.Set("checkout[different_billing_address]", "true")
		nameParts := strings.Split(runningCtx.ProfileData.Billing.Name, " ")

		paymentData.Set("checkout[billing_address][first_name]", nameParts[0])
		paymentData.Set("checkout[billing_address][last_name]", nameParts[1])
		paymentData.Set("checkout[billing_address][address1]", runningCtx.ProfileData.Billing.Line1)
		paymentData.Set("checkout[billing_address][address2]", runningCtx.ProfileData.Billing.Line2)
		paymentData.Set("checkout[billing_address][city]", runningCtx.ProfileData.Billing.City)
		paymentData.Set("checkout[billing_address][country]", runningCtx.ProfileData.Billing.Country.Code)
		paymentData.Set("checkout[billing_address][province]", runningCtx.ProfileData.Billing.State.Code)
		paymentData.Set("checkout[billing_address][zip]", runningCtx.ProfileData.Billing.PostCode)
		paymentData.Set("checkout[billing_address][phone]", runningCtx.ProfileData.Payment.Phone)
	}

	paymentData.Set("checkout[remember_me]", "false")
	paymentData.Set("checkout[remember_me]_REPLACE", "0")
	paymentData.Set("checkout[vault_phone]", runningCtx.ProfileData.Payment.Phone)
	paymentData.Set("checkout[total_price]", runningCtx.Product.Price)
	paymentData.Set("complete", "1")

	// add all the payment data
	for _, paymentDataKey := range paymentData.Keys() {
		value, exists := paymentData.Get(paymentDataKey)
		if exists {
			if paymentDataKey == "checkout[remember_me]_REPLACE" {
				body += url.QueryEscape("checkout[remember_me]") + "=" + url.QueryEscape(value.(string)) + "&"
				continue
			}
			body += url.QueryEscape(paymentDataKey.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	browserData := getBrowserData("")
	// finally add the the browser data
	for _, browserDataKey := range browserData.Keys() {
		value, exists := browserData.Get(browserDataKey)
		if exists {
			if browserDataKey == "checkout[client_details][browser_tz]" {
				body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string))
				continue
			}
			body += url.QueryEscape(browserDataKey.(string)) + "=" + url.QueryEscape(value.(string)) + "&"
		}
	}

	//fsCountKey := fmt.Sprintf("%s-count", runningCtx.CheckoutHash)
	////body = strings.Replace(body, fmt.Sprintf("&%s_REPLACE="), fmt.Sprintf("&%s=fs-count"), 1)
	//body = strings.Replace(body, fmt.Sprintf("&%s_REPLACE=", fsCountKey), fmt.Sprintf("&%s=fs-count", fsCountKey), 1)
	for _, i2 := range strings.Split(body, "&") {
		log.Println(i2)
	}
	return body
}

func generateCartForm(variant string) url.Values {
	// todo: handle extra properties in the cart (are those properties like the anti bot form hashes?)
	cartForm := make(url.Values)

	cartForm.Add("form_type", "product")
	cartForm.Add("utf8", "✓")
	cartForm.Add("id", variant)
	// todo: change quantity var for the task context data
	cartForm.Add("quantity", "1")
	return cartForm
}

func generateNextQueueForm(token string) url.Values {
	queueForm := make(url.Values)

	queueForm.Add("query", "'\n      {\n        poll(token: $token) {\n          token\n          pollAfter\n          queueEtaSeconds\n          productVariantAvailability {\n            id\n            available\n          }\n        }\n      }\n    ',")
	queueForm.Add("variables[token]", token)

	return queueForm
}

func generatePaymentForm(profile *profiles.ProfileData, shippingRate string, price string, shopId string, authToken string, paymentSession string) url.Values {
	form := make(url.Values)

	// this should always be in the form
	form.Add("_method", "patch")
	form.Add("authenticity_token", authToken)
	form.Add("previous_step", "payment_method")
	form.Add("step", "")
	form.Add("s", paymentSession)

	// todo: handle discounts
	gateway, gatewayExists := paymentGatewaysByShopId[shopId]

	if gatewayExists {
		form.Add("checkout[payment_gateway]", gateway)
	}

	isSameBillingAndShipping := false
	shippingLine1 := profile.Shipping.Line1
	billingLine1 := profile.Billing.Line1

	form.Add("checkout[credit_card][vault]", "false")

	// check if billing line 1 matches shipping line 1
	if strings.Contains(shippingLine1, billingLine1) {
		log.Println("same billing and shipping")
		isSameBillingAndShipping = true
		form.Add("checkout[different_billing_address]", "false")
	}

	if !isSameBillingAndShipping {
		form.Add("checkout[different_billing_address]", "true")

		nameParts := strings.Split(profile.Billing.Name, " ")
		form.Add("checkout[billing_address][first_name]", nameParts[0])
		form.Add("checkout[billing_address][last_name]", nameParts[1])
		form.Add("checkout[billing_address][address1]", profile.Billing.Line1)
		form.Add("checkout[billing_address][address2]", profile.Billing.Line2)
		form.Add("checkout[billing_address][city]", profile.Billing.City)
		form.Add("checkout[billing_address][country]", profile.Billing.Country.Code)
		form.Add("checkout[billing_address][province]", profile.Billing.State.Code)
		form.Add("checkout[billing_address][zip]", profile.Billing.PostCode)
		form.Add("checkout[billing_address][phone]", profile.Payment.Phone)
	}

	form.Add("checkout[total_price]", price)
	form.Add("complete", "1")
	form.Add("client_details[browser_width]", "1747")
	form.Add("client_details[browser_height]", "967")
	form.Add("client_details[javascript_enabled]", "1")
	form.Add("client_details[color_depth]", "30")
	form.Add("client_details[java_enabled]", "false")
	form.Add("client_details[browser_tz]", "300")

	form.Add("checkout[shipping_rate][id]", shippingRate)

	return form
}

func generateReviewForm(authToken string, price string, isFree bool, captchaToken string) url.Values {
	form := make(url.Values)

	form.Add("utf8", "✓")
	form.Add("_method", "patch")
	form.Add("authenticity_token", authToken)
	form.Add("checkout[total_price]", price)
	form.Add("complete", "1")
	form.Add("checkout[client_details][browser_width]", "927")
	form.Add("checkout[client_details][browser_height]", "967")
	form.Add("checkout[client_details][javascript_enabled]", "1")

	if isFree {
		form.Add("checkout[payment_gateway]", "free")
	}

	if captchaToken != "" {
		form.Add("g-recaptcha-response", captchaToken)
	}

	return form
}

// this function also sems like it's not needed since all it gets is the product image in multiple resolutions, the variant title, and the product name
func extractProductData(responseBody CustomerDataResponse) {
	//title, variantTitle, imageUrl := responseBody.Checkout.LineItems[0].Title, responseBody.Checkout.LineItems[0].VariantTitle, responseBody.Checkout.LineItems[0].ImageURL
	//
	//if
	//

}
