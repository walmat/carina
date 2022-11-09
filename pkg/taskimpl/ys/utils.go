package ys

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/avct/uasurfer"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/third_party/http"
	"regexp"
	"strconv"
	"strings"
)

/*
	Builds the `sec-ch-ua` header based on the used User-Agent

	should return => `sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`
*/
func buildSecurityHeader(userAgent string) string {
	ua := uasurfer.Parse(userAgent)

	return fmt.Sprintf("\"Chromium\";v=\"%d\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"%d\"", ua.Browser.Version.Major, ua.Browser.Version.Major)
}

func buildCorrelationHeaders() (map[string]string, error) {
	headers := make(map[string]string)

	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		return headers, err
	}
	id := hex.EncodeToString(buf)

	headers["x-instana-l"] = fmt.Sprintf("1,correlationType=web;correlationId=%s", id)
	headers["x-instana-s"] = id
	headers["x-instana-t"] = id

	return headers, nil
}

/*
	Extracts the pixel script URL value from the homepage body

	e.g. - given `<img src="https://www.yeezysupply.com/akam/11/pixel_7e63a165?a=dD03OWNhZGM4NzZkMDM2MzgxMThiN2I2OWQ4NWE2YjJlMTRkMTVmMGQ3JmpzPW9mZg==" style="visibility: hidden; position: absolute; left: -999px; top: -999px;"/>`
	should return => `7e63a165`
*/
func extractPixelScript(body string) (string, error) {
	regex := regexp.MustCompile(`akam/[\d]+/pixel_([^?|&"].*)\?`)

	match := regex.FindStringSubmatch(body)

	if len(match) == 0 {
		return "", errors.New("unable to extract pixel script")
	}

	return match[1], nil
}

/*
	Extracts the pixel script hash value from the pixel page body

	NOTE: the `_` array is hex encoded, so we have to decode that as well

	e.g. - given `var _ = [...]; {...} g = _[67]`
	should return => `_[67]` => `35b9f0cd2c485530bbb00b7a57e343f0`
*/
func extractPixelValue(body string) (string, error) {
	// TODO: need to find a better way to extract this data

	idx0 := strings.LastIndex(body, "_=")
	if idx0 == -1 {
		return "", errors.New("unable to extract idx0")
	}

	idx1 := strings.Index(body, "\"];")
	if idx1 == -1 {
		return "", errors.New("unable to extract idx1")
	}

	idx1 += 3

	array := body[idx0:idx1]

	idx0 = strings.LastIndex(body, "g=_[")
	if idx0 == -1 {
		return "", errors.New("unable to extract g value")
	}

	idx0 += 4

	idx1 = strings.LastIndex(body, "],m=_[")
	if idx1 == -1 {
		return "", errors.New("unable to extract g value")
	}

	index, err := strconv.Atoi(body[idx0:idx1])
	if err != nil {
		return "", errors.New("unable to extract index")
	}

	parts := strings.Split(array, ",")

	unicode := strings.ReplaceAll(parts[index], "\"", "")
	unicode = strings.ReplaceAll(unicode, "\\x", "")
	unicode = strings.ReplaceAll(unicode, "\"", "")

	bytes, err := hex.DecodeString(unicode)
	if err != nil {
		return "", errors.New("unable to extract bytes")
	}

	return string(bytes), nil
}

/*
	Extracts the (currently parsed as) `bazadebezolkohpepadr` value from the homepage body

	e.g. - given `<script>bazadebezolkohpepadr="2120458319"</script>`
	should return => `2120458319`
*/
func extractPixelHash(body string) (string, error) {
	Parts := strings.Split(body, "bazadebezolkohpepadr=\"")
	Final := strings.Split(Parts[1], "\"")

	return Final[0], nil
}

/*
	Extracts the `Authorization` header needed to checkout

	e.g. => `Bearer eyJfdiI6IjEiLCJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdiI6IjEiLCJleHAiOjE2MjMwNzUxMDcsImlhdCI6MTYyMzA3MzMwNywic3ViIjoie1wiX3ZcIjpcIjFcIixcImN1c3RvbWVyX2luZm9cIjp7XCJjdXN0b21lcl9pZFwiOlwiYmNVeEphZ0NmZ2hRaGpHWHFEZUpvNjVseW5cIixcImd1ZXN0XCI6dHJ1ZSxcInZpc2l0X2lkXCI6XCJiMDcxNjk2NjNhZmVmMThkNTA0MjYxYWJlYVwifX0ifQ.dl8XgldAVo8SGRDrrSAdnbD_tnRnfYwIrohjhsPW78JgTif2kukQqnB74RgKHRx6U5CTBee8ktTVwqnmtguRT66QGPPilmzCGCjqq8D3oLf7HLPBHmPTpb9zHVcDNUNECoq3jySr4R2W6sOrKOzP-5v8g31qffCvWfVKO8KDzdYGNtxg8wIAaAAmbDAmeW-DgdK81asHkE5wpJX1A4Q-Utp8OC-lNBI4dYDyqck5vjt6Y2NpG_mOX_gOmYm_m7hJzKuY4zU90SGc-GYkbBKfRPK3GthTr0LNXVsknydirpsZDI1hlBjrCxNz689-oguljGCDXIqn1-zI5hvS9DPvaB_ACr6YXESufUdYUwMUpbHPXM0pF-GYWPh8nq7Hw5Xi6h0oBMECPKh_kp90uxy_dk_2yEjnvWV1qcPTol5xK-8DJ1Q8aYih7b4guPWnpmQxMcX8CmrcSq01hjQ4ksoVGSGm-ApAAMiCcsO7wk1l3yKANw4gnAeGJoSDpI2-_ScOZIpkKh5ElnWJFFXBshXFx03JIItOn7veZqnHVKEraPGdt4TeL2E96AFZ4XSZNiJr`
*/
func extractAuthorizationHeader(resp *http.Response) string {
	for key, value := range resp.Header {
		if strings.Contains(strings.ToLower(key), "authorization") {
			return value[0]
		}
	}

	return ""
}

func extractAbckCookie(runningCtx *runningCtx) *http.Cookie {
	cookies := runningCtx.Client.Cookies(runningCtx.StoreUrl)

	for _, cookie := range cookies {
		if strings.Contains(strings.ToLower(cookie.Name), "_abck") {
			return cookie
		}
	}

	return nil
}

func extractAkamaiScript(runningCtx *runningCtx, body string) {
	regex := regexp.MustCompile(`<script\s+type="text/javascript"\s+src="/(.*)">`)

	match := regex.FindStringSubmatch(body)

	if len(match) > 1 {
		url := fmt.Sprintf("%s%s", runningCtx.StoreUrl, match[1])
		runningCtx.AkamaiUrl = url
	}
}

func extractProductDetails(runningCtx *runningCtx, body ProductDetailsResponse) {
	runningCtx.ProductName = body.Name
	runningCtx.ProductName += fmt.Sprintf("– %s", body.AttributeList.SearchColor)
	runningCtx.ProductImage = body.ViewList[0].ImageURL
}

func isFakeSplash(runningCtx *runningCtx) bool {
	cookies := runningCtx.Client.Cookies(runningCtx.StoreUrl)
	regex := regexp.MustCompile(fmt.Sprintf("%s_u", runningCtx.RecaptchaCookie))

	for _, cookie := range cookies {
		found := regex.Find([]byte(cookie.Name))

		if strings.Contains(strings.ToLower(cookie.Value), "hmac") && found != nil {
			return strings.Contains(strings.ToLower(cookie.Value), "data=1")
		}
	}

	// assume if we didn't find an hmac, it was fake
	return true
}

func getHmacCookie(runningCtx *runningCtx) *http.Cookie {
	cookies := runningCtx.Client.Cookies(runningCtx.StoreUrl)

	regex := regexp.MustCompile(fmt.Sprintf("%s_u", runningCtx.RecaptchaCookie))

	for _, cookie := range cookies {

		found := regex.Find([]byte(cookie.Name))

		if strings.Contains(strings.ToLower(cookie.Value), "hmac") && !strings.Contains(strings.ToLower(cookie.Value), "data=1") && found != nil {
			return cookie
		}
	}

	return nil
}

func removePassedCookies(runningCtx *runningCtx) {
	cookies := runningCtx.Client.Cookies(runningCtx.StoreUrl)

	for _, cookie := range cookies {
		if strings.Contains(strings.ToLower(cookie.Value), "remove_me") {
			// runningCtx.Client.DeleteCookie(runningCtx.StoreUrl, cookie.Name)
			runningCtx.Client.Inner().Jar.(*SingleJar).DeleteCookie(cookie.Name)
		}
	}
}

func setSplashCookie(runningCtx *runningCtx) {
	cookie := http.Cookie{
		Name:  "PH0ENIX",
		Value: "false",
	}

	runningCtx.Client.AddCookie(runningCtx.StoreUrl, &cookie)
}

func removeInSplashCookie(runningCtx *runningCtx) {
	// runningCtx.Client.DeleteCookie(runningCtx.StoreUrl, "akavpwr_ys_us")
	runningCtx.Client.Inner().Jar.(*SingleJar).DeleteCookie("akavpwr_ys_us")
}

func setExtraCookies(runningCtx *runningCtx) {
	cookies := []*http.Cookie{
		{
			Name:  "UserSignUpAndSave",
			Value: "1",
		},
		{
			Name:  "UserSignUpAndSaveOverlay",
			Value: "0",
		},
		{
			Name:  "default_searchTerms_CustomizeSearch",
			Value: "%5B%5D",
		},
		{
			Name:  "geoRedirectionAlreadySuggested",
			Value: "false",
		},
		{
			Name:  "wishlist",
			Value: "%5B%5D",
		},
		{
			Name:  "persistentBasketCount",
			Value: "0",
		},
		{
			Name:  "userBasketCount",
			Value: "0",
		},
		{
			Name:  "s_cc",
			Value: "true",
		},
		{
			Name:  "utag_main",
			Value: "v_id:0179da7d7d510072a1dfc37fade003072002a06a00bd0$_sn:1$_se:3$_ss:0$_st:1622870133816$ses_id:1622868327761%3Bexp-session$_pn:1%3Bexp-session$_prevpage:PRODUCT%7CYEEZY%20BOOST%20380%20ADULTS%20(GW0304)%3Bexp-1622871928166",
		},
		{
			Name:  "_ga",
			Value: "GA1.2.2013592599.1622868328",
		},
		{
			Name:  "_gid",
			Value: "GA1.2.1602442703.1622868328",
		},
		{
			Name:  "_gat_tealium_0",
			Value: "1",
		},
		{
			Name:  "_gcl_au",
			Value: "1.1.1027174275.1622868328",
		},
		{
			Name:  "_fbp",
			Value: "fb.1.1622868328402.623832784",
		},
		{
			Name:  "AMCVS_7ADA401053CCF9130A490D4C%40AdobeOrg",
			Value: "1",
		},
		{
			Name:  "AMCV_7ADA401053CCF9130A490D4C%40AdobeOrg",
			Value: "-227196251%7CMCIDTS%7C18784%7CMCMID%7C53841739473916663753859934600661933645%7CMCAAMLH-1623473127%7C7%7CMCAAMB-1623473127%7C6G1ynYcLPuiQxYZrsz_pkqfLG9yMXBpb2zX5dvJdYQJzPXImdj0y%7CMCOPTOUT-1622875528s%7CNONE%7CMCAID%7CNONE",
		},
	}

	runningCtx.Client.SetCookies(runningCtx.StoreUrl, cookies)
}

/*
	Weak check to see if the response from Yeezy Supply gives us the banned response

	should return => true || false
*/
func isBanned(runningCtx *runningCtx, body string, length int) bool {
	if length != -1 {
		return runningCtx.BannedMatcher.Contains([]byte(body)) && len(body) <= length
	}

	return runningCtx.BannedMatcher.Contains([]byte(body))
}

/*
	Check to see if the product is a flash sale and not waiting room

	return => true || false
*/
func isFlashSale(body string, runningCtx *runningCtx) bool {
	regex := regexp.MustCompile(runningCtx.SplashFlag)

	match := regex.Find([]byte(body))

	return match == nil
}

func rotateProxy(taskCtx *tasks.Context, runningCtx *runningCtx) error {
	newDialer, err := runningCtx.getProxyDialer(taskCtx)
	if err != nil {
		return err
	}
	clientConfig := http.DefaultClientConfig.Clone()
	clientConfig.ProxyDialer = newDialer
	newTransport := http.NewTlsTransport(http.DefaultClientConfig, http.FingerprintChrome)
	runningCtx.Client.Inner().Transport = newTransport
	runningCtx.Client.Inner().Jar = newSingleJar()
	return nil
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

func checkPixel(resp *http.Response, runningCtx *runningCtx) bool {
	PixelScript, PixelErr := extractPixelScript(resp.StringBody())
	if PixelErr != nil {
		return false
	}

	runningCtx.PixelScript = PixelScript

	PixelHash, PixelHashError := extractPixelHash(resp.StringBody())
	if PixelHashError != nil {
		return false
	}

	runningCtx.PixelHash = PixelHash

	return true
}
