package foots

import (
	"errors"
	"fmt"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/pkg/util"
	"nebula/third_party/http"
	tls "nebula/third_party/utls"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func getSkuCode(ctx *staticCtx, stockResp StockResponse) string {
	for _, variantAttr := range stockResp.VariantAttributes {
		if variantAttr.Sku == ctx.Default.Sku {
			return variantAttr.Code
		}
	}

	return ""
}

func extractStyleCodes(style string, sellableUnits []SellableUnit) []SellableUnit {
	var variants []SellableUnit
	for _, sellableUnit := range sellableUnits {
		styleId := ""
		for _, attr := range sellableUnit.Attributes {
			if strings.Contains(attr.Type, "style") {
				styleId = attr.ID
				break
			}
		}

		if style == styleId {
			variants = append(variants, sellableUnit)
		}
	}

	return variants
}

func predicate(size string, variants []taskutil.Variant) (taskutil.Variant, error) {

	trimmed := taskutil.TrimLeadingZeroes(size)
	numerical, _ := regexp.MatchString("[0-9]+", trimmed)

	for _, variant := range variants {
		if numerical {
			regex := regexp.MustCompile(fmt.Sprintf("^%s([^.]*)", trimmed))
			replacer := regexp.MustCompile(`/^[^0-9]+/`)

			fmt.Printf("Matching size: %s, Desired size: %s\n", replacer.ReplaceAllString(variant.Size, ""), trimmed)

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

func getFootsiteType(runningCtx *runningCtx) {
	switch runningCtx.StoreUrl.Hostname() {
	case "www.footlocker.com":
		runningCtx.FootsiteType = FootsiteFTL
	case "www.footaction.com":
		runningCtx.FootsiteType = FootsiteFootaction
	case "www.eastbay.com":
		runningCtx.FootsiteType = FootsiteEastBay
	case "www.champssports.com":
		runningCtx.FootsiteType = FootsiteChamps
	case "www.kidsfootlocker.com":
		runningCtx.FootsiteType = FootsiteKFTL
	default:
		runningCtx.FootsiteType = FootsiteFTL
	}
}

func setStore(footsiteType FootsiteType, runningCtx *runningCtx) {
	switch footsiteType {
	case FootsiteFTL:
		runningCtx.StoreUrl, _ = url.Parse("https://www.footlocker.com")
	case FootsiteFootaction:
		runningCtx.StoreUrl, _ = url.Parse("https://www.footaction.com")
	case FootsiteEastBay:
		runningCtx.StoreUrl, _ = url.Parse("https://www.eastbay.com")
	case FootsiteChamps:
		runningCtx.StoreUrl, _ = url.Parse("https://www.champssports.com")
	case FootsiteKFTL:
		runningCtx.StoreUrl, _ = url.Parse("https://www.kidsfootlocker.com")
	default:
		runningCtx.StoreUrl, _ = url.Parse("https://www.footlocker.com")
	}
}

func configureClient(taskCtx *tasks.Context, runningCtx *runningCtx) {
	clientConfig := http.InsecureClientConfig.Clone()

	clientConfig.ClientHello = tls.HelloGolang
	clientConfig.TlsConfig.ServerName = runningCtx.StoreUrl.Hostname()

	runningCtx.Cookies = make(map[string]string)

	proxyDialer, err := runningCtx.getProxyDialer(taskCtx)
	if err == nil {
		clientConfig.ProxyDialer = proxyDialer
	}

	client := http.NewSimpleClient(taskCtx.Context, clientConfig, http.FingerprintChrome)
	runningCtx.Client = client
}

func extractCookies(headers http.Header, runningCtx *runningCtx) {
	for header := range headers {
		if strings.Contains(strings.ToLower(header), "set-cookie") {
			for _, cookie := range headers[header] {
				data := strings.Split(cookie, ";")[0]
				parts := strings.Split(data, "=")

				if len(parts) > 0 {
					name := parts[0]
					value := parts[1]

					if name != "" {
						if name == "JSESSIONID" {
							runningCtx.SessionId = value
						}

						runningCtx.Cookies[name] = value
					}
				}
			}
		}
	}
}

func makeCookies(runningCtx *runningCtx) string {
	cookieHeader := ""
	for key, value := range runningCtx.Cookies {
		cookieHeader += fmt.Sprintf("%s=%s; ", key, value)
	}

	return cookieHeader
}

func handleQueue(taskCtx *tasks.Context, state State) State {
	taskCtx.SendStatusColored("In queue", util.ColorWarning)
	time.Sleep(5 * time.Second)
	return state
}
