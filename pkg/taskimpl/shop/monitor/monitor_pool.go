package monitor

// TODO: Move to general utils package: buildSecurityHeader, getUserAgent and ua variable, single jar

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/pkg/util"
	"strconv"
	"strings"
	"time"
)

var (
	// note: change types string to type store name and product name. map[StoreName]map[ProductName]*StoreProduct{}
	Pool = StoresPool{ProductsByStore: map[string]map[string]*StoreProduct{}}
	invalidMonitorType = errors.New("invalid monitor type")
)

func (m *StoresPool) GetProduct(taskCtx *tasks.Context, taskData tasks.DefaultData) <- chan Result {
	taskCtx.SendStatus("Monitoring")

	p := make(chan Result)

	monitorType := getMonitorType(taskData.Sku)

	_, productExists := m.ProductsByStore[taskData.Store.Name][taskData.Sku]

	if productExists && monitorType != "variant" {
		variantData, err := m.getVariantData(m.ProductsByStore[taskData.Store.Name][taskData.Sku].Variants, taskData)

		if err != nil {
			taskCtx.SendStatusColored("No variations matched, retrying...", util.ColorWarning)
			time.Sleep(2 * time.Second)
			go func() {
				p <- createResult(nil, nil)
			}()
			return p
		}

		go func() {
			p <- createResult(variantData, nil)
		}()
		return p
	}

	switch monitorType {
	case "variant":
		go func() {
			p <- createResult(&StoreProduct{
				ProductName: taskData.Sku,
				Sku:         taskData.Sku,
				Variants:    taskData.Sku,
			}, nil)
		}()

		return p
	case "url":
		go m.findProductByUrl(p, taskData, taskCtx)
		return p
	case "keywords":
		log.Println("keyword sku, product needs to be found")
		go m.findProductByKeywords(p, taskData, taskCtx)
		return p
	default:
		go func() {
			p <- createResult(nil, invalidMonitorType)
		}()
		return p
	}
}

func (m *StoresPool) getVariantData(product interface{}, taskData tasks.DefaultData) (*StoreProduct, error) {
	variantsData, isVariantsData := product.([]GetJsonProductsVariant)
	productUrlResponse, isProductUrlResponse := product.(ProductResponseJs)
	var vars []taskutil.Variant

	if isVariantsData {
		for _, variant := range variantsData {
			trimmed := taskutil.Variant{
				ID:      strconv.FormatInt(variant.ID, 10),
				Size:    getSize(variant),
				InStock: variant.Available,
			}
			vars = append(vars, trimmed)
		}

		// Find any matches for sizes
		matches, err := taskutil.MatchVariants(taskData.Sizes, vars, predicate)
		if err != nil {
			time.Sleep(2 * time.Second)
			return nil, err
		}

		rand.Seed(time.Now().UnixNano())
		chosenVariant := matches[rand.Intn(len(matches))]

		// Save all the product's variants for other tasks to find
		m.setProduct(productUrlResponse.Title, taskData, variantsData)

		// Since it's the first time this product is being found then the task has to return the matched size
		return &StoreProduct{
			ProductName: productUrlResponse.Title,
			Sku:         taskData.Sku,
			Variants:    chosenVariant,
		}, nil
	}

	// If it is the first time the task is getting a product's variants we set them all in our struct
	if isProductUrlResponse {
		// Format all the variant into our own variant type
		for _, variant := range productUrlResponse.Variants {
			trimmed := taskutil.Variant{
				ID:      strconv.FormatInt(variant.ID, 10),
				Size:    getSize(variant),
				InStock: variant.Available,
			}
			vars = append(vars, trimmed)
		}
		// Find any matches for sizes
		matches, err := taskutil.MatchVariants(taskData.Sizes, vars, predicate)

		if err != nil {
			time.Sleep(2 * time.Second)
			return nil, err
		}

		rand.Seed(time.Now().UnixNano())
		chosenVariant := matches[rand.Intn(len(matches))]

		// Save all the product's variants for other tasks to find
		m.setProduct(productUrlResponse.Title, taskData, productUrlResponse.Variants)

		// Since it's the first time this product is being found then the task has to return the matched size
		return &StoreProduct{
			ProductName: productUrlResponse.Title,
			Sku:         taskData.Sku,
			Variants:    chosenVariant,
		}, nil
	}

	// If the task gets to this stage it means the product has already been found
	variantArray := m.ProductsByStore[taskData.Store.Name][taskData.Sku].Variants.([]ProductResponseVariant)
	for _, variant := range variantArray {
		trimmed := taskutil.Variant{
			ID:      strconv.FormatInt(variant.ID, 10),
			// note: do option1 and option2 matter?
			Size:    getSize(variant),
			InStock: variant.Available,
		}
		vars = append(vars, trimmed)
	}

	// Find any matches for sizes
	matches, err := taskutil.MatchVariants(taskData.Sizes, vars, predicate)
	if err != nil {
		time.Sleep(2 * time.Second)
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	chosenVariant := matches[rand.Intn(len(matches))]

	return &StoreProduct{
		ProductName: m.ProductsByStore[taskData.Store.Name][taskData.Sku].ProductName,
		Sku:        taskData.Sku,
		Variants:    chosenVariant,
	}, nil
}

func (m *StoresPool) setProduct(productName string, taskData tasks.DefaultData, productData interface{}) {
	if Pool.ProductsByStore[taskData.Store.Name][taskData.Sku] == nil {
		Pool.ProductsByStore[taskData.Store.Name] = map[string]*StoreProduct{}
	}

	Pool.ProductsByStore[taskData.Store.Name][taskData.Sku] = &StoreProduct{
		ProductName: productName,
		Sku: taskData.Sku,
		Variants: productData,
	}
}

func (m *StoresPool) findProductByUrl(resultChan chan Result, taskData tasks.DefaultData, taskCtx *tasks.Context) {
	taskCtx.SendStatus("Getting product data")
	// note: the .json product data endpoint (product url + ".json") doesnt have the product availability key
	endpoint := fmt.Sprintf("%s.js?order=-%d", taskData.Sku, rand.Intn(99999 - 10000) + 10000)

	resp, err := Ctx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(Ctx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", Ctx.UserAgent).
		Header("accept", "application/json").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, err)
		return
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Product not live, retrying...", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	if resp.StatusCode == 401 {
		taskCtx.SendStatusColored("Password page", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	// Switch proxy?
	if resp.StatusCode == 429 || resp.StatusCode == 430 || resp.StatusCode == 403 {
		taskCtx.SendStatusColored("Banned, retrying...", util.ColorFailed)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	var body ProductResponseJs
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing product, retrying...", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}
	variantData, err := m.getVariantData(body, taskData)
	if err != nil {
		taskCtx.SendStatusColored("No variations matched, retrying...", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
	}

	// Return only the variant that was wanted
	resultChan <- createResult(variantData, nil)
}

func (m *StoresPool) findProductByKeywords(resultChan chan Result, taskData tasks.DefaultData, taskCtx *tasks.Context) {
	endpoint := fmt.Sprintf("%s/products.json?order=-%d", taskData.Store.Url, rand.Intn(99999 - 10000) + 10000)

	resp, err := Ctx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(Ctx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", Ctx.UserAgent).
		Header("accept", "application/json").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, err)
		return
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Products not found, retrying...", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	if resp.StatusCode == 401 {
		taskCtx.SendStatusColored("Password page", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	// Switch proxy?
	if resp.StatusCode == 429 || resp.StatusCode == 430 || resp.StatusCode == 403 {
		taskCtx.SendStatusColored("Banned, retrying...", util.ColorFailed)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	var body GetJsonProductsResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing product, retrying...", util.ColorWarning)
		time.Sleep(2 * time.Second)
		resultChan <- createResult(nil, nil)
		return
	}

	var product []GetJsonProductsVariant
	for _, p := range body.Products {
		if matchKeyword(strings.ToLower(p.Title), taskData.Sku) {
			product = p.Variants
		}
	}

	if len(product) >= 1 {
		variantData, err := m.getVariantData(product, taskData)
		if err != nil {
			taskCtx.SendStatusColored("No variations matched, retrying...", util.ColorWarning)
			time.Sleep(2 * time.Second)
			resultChan <- createResult(nil, nil)
		}

		resultChan <- createResult(variantData, nil)
	}

	taskCtx.SendStatus("Monitoring")
	time.Sleep(2 * time.Second)
	resultChan <- createResult(nil, nil)
}