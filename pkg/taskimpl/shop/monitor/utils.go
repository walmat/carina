package monitor

import (
	"errors"
	"fmt"
	"github.com/avct/uasurfer"
	"nebula/pkg/taskutil"
	"regexp"
	"strings"
)

// getMonitorType Checks if the task product input is an url, variant, or keywords
func getMonitorType(sku string) string {
	variant, url := regexp.MustCompile("^\\d+$"), regexp.MustCompile("https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")

	if variant.MatchString(sku) {
		return "variant"
	} else if url.MatchString(sku) {
		return "url"
	}
	return "keywords"
}

func matchKeyword(productTile string, keywords string) bool  {
	_keywords := strings.Split(keywords, ",")

	positiveKeywords := 0
	keywordsFound := 0
	result := false

	for _, k := range _keywords {
		if strings.Index(k, "+") == 0 {
			cleanKeyword := strings.Replace(k, "+", "", 1)
			positiveKeywords += 1
			if strings.Contains(productTile, cleanKeyword) {
				keywordsFound += 1
			}
		} else if strings.Index(k, "-") == 0 {
			cleanKeyword := strings.Replace(k, "-", "", 1)
			if strings.Contains(productTile, cleanKeyword) {
 				return false
			}
		}
	}

	if keywordsFound == positiveKeywords {
		 result = true
	}

	return result
}

func getSize(variant interface{}) string {
	productUrlResponseVariant, isProductUrlResponseVariant := variant.(ProductResponseVariant)
	variantData, isVariantData := variant.(GetJsonProductsVariant)

	var size string

	if isVariantData {
		possibleVariantSizes := map[string]string{
			"title": variantData.Title,
		}

		if variantData.Option1 != "" {
			possibleVariantSizes["option1"] = variantData.Option1
		}

		if variantData.Option2 != nil && variantData.Option2.(string) != "" {
			possibleVariantSizes["option2"] = variantData.Option2.(string)
		}

		if variantData.Option3 != nil && variantData.Option3.(string) != ""{
			possibleVariantSizes["option3"] = variantData.Option3.(string)
		}


		for _, v := range possibleVariantSizes {
			if v == "" {
				continue
			}
			if regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(v){
				size = v
			}
		}

		return size
	}

	if isProductUrlResponseVariant {
		// first check the options array which seems to have all the available variant options
		if len(productUrlResponseVariant.Options) >= 1 {
			for _, o := range productUrlResponseVariant.Options {
				// since some of the variants unwanted characters like "-, _, \" we only want an option that does not contain those, dots are allowed
				if regexp.MustCompile(`^[a-zA-Z0-9.]+$`).MatchString(o){
					size = o
				}
			}
		}
		possibleVariantSizes := map[string]string{
			"title": productUrlResponseVariant.Title,
		}

		if productUrlResponseVariant.Option1 != "" {
			possibleVariantSizes["option1"] = productUrlResponseVariant.Option1
		}

		if productUrlResponseVariant.Option2 != nil && productUrlResponseVariant.Option2.(string) != "" {
			possibleVariantSizes["option2"] = productUrlResponseVariant.Option2.(string)
		}

		if productUrlResponseVariant.Option3 != nil && productUrlResponseVariant.Option3.(string) != ""{
			possibleVariantSizes["option3"] = productUrlResponseVariant.Option3.(string)
		}


		for _, v := range possibleVariantSizes {
			if v == "" {
				continue
			}
			if regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(v){
				size = v
			}
		}
	}

	return size
}

func createResult(product *StoreProduct, err error) Result {
	result := Result{
		Product: product,
		Error:   nil,
	}

	if err != nil {
		result.Error = err
		result.Product = nil
		return result
	}

	return result
}

func buildSecurityHeader(userAgent string) string {
	ua := uasurfer.Parse(userAgent)
	return fmt.Sprintf("\"Chromium\";v=\"%d\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"%d\"", ua.Browser.Version.Major, ua.Browser.Version.Major)
}

func predicate(size string, variants []taskutil.Variant) (taskutil.Variant, error) {

	trimmed := taskutil.TrimLeadingZeroes(size)
	numerical, _ := regexp.MatchString("[0-9]+", trimmed)

	for _, variant := range variants {
		if numerical {
			regex := regexp.MustCompile(fmt.Sprintf("^%s([^.]*)", trimmed))
			replacer := regexp.MustCompile(`/^[^0-9]+/`)

			//fmt.Printf("Matching size: %s, Desired size: %s\n", replacer.ReplaceAllString(variant.Size, ""), trimmed)

			if regex.MatchString(taskutil.TrimLeadingZeroes(replacer.ReplaceAllString(variant.Size, ""))) {
				fmt.Printf("Matched variant size: %s\n", variant.Size)
				return variant, nil
			}
		} else {
			regex := regexp.MustCompile(fmt.Sprintf("^%s", trimmed))
			if regex.MatchString(strings.TrimSpace(variant.Size)) {
				fmt.Printf("Matched variant size: %s\n", variant.Size)
				return variant, nil
			}
		}
	}

	return taskutil.Variant{}, errors.New("no match")
}