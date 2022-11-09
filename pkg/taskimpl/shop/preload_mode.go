package shopify

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"nebula/pkg/infra/accounts"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/webhooks"
	"nebula/pkg/taskimpl/shop/monitor"
	"nebula/pkg/taskutil"
	"nebula/pkg/util"
	"nebula/third_party/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Preload mode
// Get home -> Get config -> Get preload product -> submit preload cart -> create checkout -> clear cart -> wait for product -> submit cart -> get customer

// Safe mode
// Get home -> Get config -> wait for product -> submit cart -> create checkout -> get customer
func handleSetupPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	runningCtx.StoreUrl, _ = url.Parse(staticCtx.Default.Store.Url)

	// These variables are supposed to come from the task ctx
	runningCtx.UsesPassword = false
	runningCtx.StorePassword = "supersecretpassword123"
	runningCtx.UsesAccount = false
	runningCtx.AccountCredentials = ""

	runningCtx.UserAgent = getUserAgent()

	runningCtx.Restocking = false
	runningCtx.Preloading = true
	runningCtx.Fallback = false
	runningCtx.PollingUrl = false
	runningCtx.RetryCount = 0
	runningCtx.ShippingRate = ""

	err := monitor.Ctx.Init(taskCtx)
	if err != nil {
		taskCtx.SendStatus("Error initializing task")
		time.Sleep(1 * time.Second)
		return Setup, err
	}
	go func() {
		monitor.Pool.GetProduct(taskCtx, staticCtx.Default)
	}()
	return GetHomepage, nil
}

func handleGetHomepagePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting homepage")

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error visiting homepage", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetHomepage, err
	}

	if resp.StatusCode == 401 {
		taskCtx.SendStatus("Password page")
		time.Sleep(3 * time.Second)
		return GetHomepage, nil
	}

	if location, err := resp.Location(); location != nil && err != nil {
		runningCtx.Redirect = location.String()
	}

	return GetConfig, nil
}

func handleGetConfigPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting config")

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/payments/config").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting config", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, err
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Invalid shopify store", util.ColorFailed)
		return "", errors.New("invalid shopify store")
	}

	var body GetConfigResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing config, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, nil
	}

	if body.PaymentInstruments.AccessToken == "" {
		taskCtx.SendStatusColored("Invalid shopify store", util.ColorFailed)
		return "", errors.New("invalid shopify store")
	}

	runningCtx.AccessToken = body.PaymentInstruments.AccessToken
	runningCtx.ShopID = body.PaymentInstruments.CheckoutConfig.ShopID

	if runningCtx.UsesPassword {
		log.Println("submit password")
		//return SubmitPassword, nil
	}

	if runningCtx.UsesAccount {
		return GetAccount, nil
	}

	return GetPreloadProduct, nil
}

func handleGetAccountPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting account")

	_, err := accounts.GetAccountForTask(taskCtx.Identifier())
	if err != nil {
		if err == accounts.NotLinkedErr {
			taskCtx.SendStatusColored("Account needed", util.ColorFailed)
			return "", nil
		}
		return "", nil
	}


	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/account/login", runningCtx.StoreUrl.String())).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting account", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetAccount, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")
		if strings.Contains(location, "password") {
			taskCtx.SendStatus("Password page")
			if runningCtx.UsesPassword {
				runningCtx.NextState = GetAccount
				return SubmitPassword, nil
			}
			time.Sleep(3 * time.Second)
			return GetAccount, nil
		}

		if strings.Contains(location, "challenge") {
			return GetChallenge, nil
		}
	}

	if strings.Contains(resp.StringBody(), "window.Shopify.recaptchaV3") {
		runningCtx.Sitekey = "6LcCR2cUAAAAANS1Gpq_mDIJ2pQuJphsSQaUEuc9"
		runningCtx.NextState = SubmitAccount
		return GetCaptchaToken, nil
	}

	return SubmitAccount, nil
}

func handleGetCaptchaTokenPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	var p string
	var captchaUrl string
	var captchaType string

	if runningCtx.ProxyData != nil {
		p = runningCtx.ProxyData.Url.String()
	}

	if runningCtx.NextState == SubmitAccount {
		taskCtx.SendStatusColored("Login captcha", util.ColorWarning)
		captchaUrl = runningCtx.StoreUrl.String() + "/account/login"
		captchaType = harvesters.ReCaptcha_V3
	}

	if runningCtx.NextState == SubmitChallenge {
		taskCtx.SendStatusColored("Challenge captcha", util.ColorWarning)
		captchaUrl = runningCtx.StoreUrl.String() + "/challenge"
		captchaType = harvesters.ReCaptcha_V2
	}

	r, err := harvesters.RequestSolve(harvesters.Solver{
		ID:   taskCtx.Identifier(),
		Type: captchaType,
		Store: &harvesters.Store{
			Url: captchaUrl,
		},
		Cookies: runningCtx.Client.Inner().Jar.Cookies(runningCtx.StoreUrl),
		Proxy:   &p,
		ReCaptcha: &harvesters.ReCaptcha{
			SiteKey: runningCtx.Sitekey,
			Action:  "",
		},
		RequestedAt: time.Now().Unix(),
	})

	if err != nil {
		taskCtx.SendStatusColored("Harvester error, retrying...", util.ColorFailed)
		return GetHomepage, nil
	}

	s := <-r.Collective

	runningCtx.CaptchaToken = s.Data.Token

	if runningCtx.NextState == SubmitAccount {
		return SubmitAccount, nil
	}

	if runningCtx.NextState == SubmitChallenge {
		return SubmitChallenge, nil
	}

	// return to the next state
	// final state = create checkout
	return "", nil
}

func handleSubmitAccountPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	account, err := accounts.GetAccountForTask(taskCtx.Identifier())
	if err != nil {
		if err == accounts.NotLinkedErr {
			taskCtx.SendStatusColored("Account needed", util.ColorFailed)
			return "", nil
		}
		return "", nil
	}

	form := url.Values{}
	form.Add("form_type", "customer_login")
	form.Add("utf8", "✓")
	form.Add("customer[email]", account.Username)
	form.Add("customer[password]", account.Password)
	form.Add("recaptcha-v3-token", runningCtx.CaptchaToken)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/account/login").
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", runningCtx.StoreUrl.String()+"/account/login?return_url=%2Faccount").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting account", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetChallenge, err
	}

	runningCtx.CaptchaToken = ""

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "challenge") {
			return GetChallenge, nil
		}

		if strings.Contains(location, "login") {
			taskCtx.SendStatusColored("Invalid credential", util.ColorFailed)
			time.Sleep(2 * time.Second)
			return SubmitAccount, nil
		}

		if strings.Contains(location, "account") {
			// if the submit account handler was called from another handler other than the config handler
			if runningCtx.NextState != SubmitAccount && runningCtx.NextState != SubmitChallenge {
				return runningCtx.NextState, nil
			}
			return GetPreloadProduct, nil
		}

	}
	taskCtx.SendStatusColored("Error submitting account", util.ColorFailed)
	time.Sleep(2 * time.Second)
	return SubmitAccount, nil
}

func handleGetChallengePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting challenge")

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/challenge", runningCtx.StoreUrl.String())).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", fmt.Sprintf("%s/account/login", runningCtx.StoreUrl.String())).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting challenge", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetChallenge, err
	}

	if resp.StatusCode != 200 {
		taskCtx.SendStatusColored(fmt.Sprintf("Error getting challenge %d", resp.StatusCode), util.ColorFailed)
		time.Sleep(3 * time.Minute)
		return GetChallenge, nil
	}

	runningCtx.NextState = SubmitChallenge
	runningCtx.Sitekey = extractCaptchaData(resp.StringBody())
	runningCtx.AuthToken = extractAuthToken(resp.StringBody())

	return GetCaptchaToken, nil
}

func handleSubmitChallengePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting challenge")
	form := url.Values{}
	form.Add("authenticity_token", runningCtx.AuthToken)
	form.Add("g-recaptcha-response", runningCtx.CaptchaToken)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/account/login").
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", runningCtx.StoreUrl.String()+"/challenge").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting challenge", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetChallenge, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "challenge") {
			return GetChallenge, nil
		}

		if strings.Contains(location, "password") {
			if runningCtx.UsesPassword {
				runningCtx.NextState = SubmitChallenge
				return SubmitPassword, nil
			}
			taskCtx.SendStatusColored("Password page, retrying...", util.ColorWarning)
			time.Sleep(2 * time.Second)
			return SubmitChallenge, nil
		}

		if strings.Contains(location, "login") {
			taskCtx.SendStatusColored("Invalid credentials", util.ColorFailed)
			time.Sleep(2 * time.Second)
			return GetAccount, nil
		}

		if strings.Contains(location, "account") {
			runningCtx.CaptchaToken = ""
			return GetPreloadProduct, nil
		}
	}
	taskCtx.SendStatusColored("Error submitting challenge", util.ColorFailed)
	return GetChallenge, nil
}

func handleWaitForProductPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	result := <-monitor.Pool.GetProduct(taskCtx, staticCtx.Default)
	if result.Product == nil {
		return WaitForProduct, nil
	}
	// if the variant is a string then it means it came from the task static ctx
	taskCtxVariant, isTaskCtxVariant := result.Product.Variants.(string)

	if isTaskCtxVariant {
		runningCtx.Product.Variant = taskutil.Variant{
			ID:      taskCtxVariant,
			Size:    "",
			InStock: false,
		}
	} else {
		runningCtx.Product.Variant = result.Product.Variants.(taskutil.Variant)
	}

	runningCtx.Product.Name = result.Product.ProductName

	return SubmitCartSafe, nil
}

func handleGetPreloadProduct(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Preloading product")

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/products.json?limit=250&order=%d", runningCtx.StoreUrl.String(), rand.Intn(7))).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error preloading product", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, err
	}

	if resp.StatusCode == 401 {
		// todo: change mode to safe mode
		taskCtx.SendStatus("Password page")
		time.Sleep(3 * time.Second)
		return GetPreloadProduct, nil
	}

	if resp.StatusCode != 200 {
		taskCtx.SendStatusColored(fmt.Sprintf("Error preloading product %d, retrying...", resp.StatusCode), util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPreloadProduct, nil
	}

	var body GetProductResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error preloading product, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPreloadProduct, nil
	}

	if len(body.Products) == 0 {
		// todo: change mode to safe mode
		taskCtx.SendStatusColored("No products loaded", util.ColorWarning)
		// find checkout product
		return "", nil
	}

	// Select the first product and the first variant
	runningCtx.PreloadProduct.Variant = body.Products[0].Variants[0]

	runningCtx.PreloadProduct.ProductUrl = fmt.Sprintf("%s/products/%s", runningCtx.StoreUrl.String(), body.Products[0].Handle)

	return SubmitCartSafe, nil
}

func handleEnterQueuePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	cookies := runningCtx.Client.Inner().Jar.Cookies(runningCtx.StoreUrl)

	var isNewQueue *http.Cookie
	for _, c := range cookies {
		if c.Name == "_checkout_queue_token" {
			isNewQueue = c
			break
		}
	}

	if isNewQueue != nil {
		q, err := url.QueryUnescape(isNewQueue.Value)

		if err != nil {
			taskCtx.SendStatusColored("Error parsing queue", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return runningCtx.PrevState, nil
		}

		runningCtx.NextQueueToken = q
		runningCtx.UseNewQueue = true
	}

	var ctd *http.Cookie
	for _, c := range cookies {
		if strings.Contains(c.Name, "_ctd") {
			ctd = c
			break
		}
	}

	if ctd != nil {
		q, err := url.QueryUnescape(ctd.Value)

		if err != nil {
			taskCtx.SendStatusColored("Error parsing CTD", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return runningCtx.PrevState, nil
		}

		runningCtx.Ctd = q
	}

	toggleRedirect(false, runningCtx)

	var endpoint string
	if runningCtx.Ctd != "" {
		endpoint = fmt.Sprintf("/throttle/queue?_ctd=%s&_ctd_update=", runningCtx.Ctd)
	} else {
		endpoint = "/throttle/queue"
	}

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+endpoint).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", runningCtx.StoreUrl.String()+"/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error entering queue", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return EnterQueue, err
	}

	if resp.StatusCode == 302 {
		location, err := resp.Location()

		if err != nil {
			taskCtx.SendStatusColored("Error getting redirect", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return runningCtx.PrevState, err
		}

		runningCtx.CheckoutHash, err = extractCheckoutHash(location.String())
		if err != nil {
			taskCtx.SendStatusColored("Error getting checkout hash", util.ColorWarning)
			time.Sleep(3 * time.Second)
			return CreateCheckout, err
		}

		if strings.Contains(location.String(), "processing") {
			return GetOrder, nil
		}

		if runningCtx.CheckoutHash != "" {
			return CreateCheckout, nil
		}

		if runningCtx.NextState != "" {
			nextState := runningCtx.NextState
			runningCtx.NextState = ""
			return nextState, nil
		}

		return GetCustomer, nil
	}

	return GetQueue, nil
}

func handleGetQueuePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	if runningCtx.UseNewQueue {
		return GetNextQueue, nil
	}

	referrer := fmt.Sprintf("%s/throttle/queue?_ctd=%s&_ctd_update", runningCtx.StoreUrl.String(), url.QueryEscape(runningCtx.Ctd))

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/checkout/poll?js_poll=1").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", referrer).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting queue", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetQueue, err
	}

	if resp.StatusCode != 200 {
		taskCtx.SendStatusColored("Error getting queue", util.ColorFailed)
		var retryAfter int
		h := resp.Header.Get("retry-after")
		if h != "" {
			i, err := strconv.Atoi(h)
			if err != nil || i == 0 {
				retryAfter = 5
			} else {
				retryAfter = i
			}
		}

		time.Sleep(time.Duration(retryAfter) * time.Second)
		return GetQueue, nil
	}

	return GetPastQueue, nil
}

func handleGetPastQueuePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	toggleRedirect(false, runningCtx)

	ctdUrl := fmt.Sprintf("%s/throttle/queue?_ctd=%s&_ctd_update", runningCtx.StoreUrl.String(), url.QueryEscape(runningCtx.Ctd))

	resp, err := runningCtx.Client.Builder(ctdUrl).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error passing queue", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetQueue, err
	}

	if resp.StatusCode == 302 {
		u, err := resp.Location()

		if err != nil {
			taskCtx.SendStatusColored("Error passing queue", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetPastQueue, err
		}

		runningCtx.CheckoutHash, err = extractCheckoutHash(u.String())

		if err != nil {
			taskCtx.SendStatusColored("Error passing queue", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetPastQueue, err
		}

		if strings.Contains(u.String(), "checkpoint") {
			return GetCheckpoint, nil
		}

		if strings.Contains(u.String(), "processing") {
			return GetOrder, nil

		}
	}

	if runningCtx.NextState != "" {
		nextState := runningCtx.NextState
		runningCtx.NextState = ""
		return nextState, nil
	}

	return GetCustomer, nil
}

func handleGetNextQueuePreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	form := generateNextQueueForm(runningCtx.NextQueueToken)

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/queue/poll").
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", runningCtx.StoreUrl.String()+"/throttle/queue").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting next queue", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetCheckpoint, err
	}

	var body GetNextQueueResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing next queue, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, nil
	}

	if strings.Contains(body.Data.Poll.__Typename, "continue") {
		runningCtx.NextQueueToken = body.Data.Poll.Token

		c := []*http.Cookie{
			&http.Cookie{
				Name:  "_checkout_queue_token",
				Value: url.QueryEscape(runningCtx.NextQueueToken),
			},
		}

		runningCtx.Client.SetCookies(runningCtx.StoreUrl, c)

		if body.Data.Poll.QueueEtaSeconds != 0 {
			runningCtx.QueueETA = body.Data.Poll.QueueEtaSeconds
		}

		if len(body.Data.Poll.ProductVariantAvailability) > 0 {
			prod := body.Data.Poll.ProductVariantAvailability[0]

			if !prod.Available {
				runningCtx.QueueAvailability = "OOS"
			}
		}

		message := "Waiting in queue"
		if runningCtx.QueueETA != 0 {
			message = `Waiting in queue [${this.queueEta}s]`
		}

		if runningCtx.QueueAvailability != "" {
			message = `Waiting in queue [${this.queueAvailability}]`
		}

		taskCtx.SendStatus(message)
		time.Sleep(10 * time.Second)

		return GetNextQueue, nil
	}

	if strings.Contains(body.Data.Poll.__Typename, "complete") {
		t, err := url.QueryUnescape(runningCtx.NextQueueToken)

		if err != nil {
			taskCtx.SendStatusColored("Error parsing queue token, retrying...", util.ColorFailed)
			time.Sleep(5 * time.Second)

			return GetNextQueue, err
		}

		c := []*http.Cookie{
			&http.Cookie{
				Name:  "_checkout_queue_token",
				Value: t,
			},
		}

		runningCtx.Client.SetCookies(runningCtx.StoreUrl, c)

		if runningCtx.CheckoutHash != "" {
			if runningCtx.NextState != "" {
				nextState := runningCtx.NextState
				runningCtx.NextState = ""
				return nextState, nil
			}

			return GetCustomer, nil
		}

		return CreateCheckout, nil
	}

	time.Sleep(5 * time.Second)

	return GetQueue, nil
}

func handleGetCheckpointPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	if runningCtx.Redirect == "" {
		taskCtx.SendStatusColored("Error getting checkpoint", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return runningCtx.PrevState, errors.New("no redirect url specified")
	}

	var p string
	if runningCtx.ProxyData != nil {
		p = runningCtx.ProxyData.Url.String()
	}
	r, err := harvesters.RequestSolve(harvesters.Solver{
		ID:   taskCtx.Identifier(),
		Type: harvesters.Shopify_Checkpoint,
		Store: &harvesters.Store{
			Url: runningCtx.Redirect,
		},
		Cookies:     runningCtx.Client.Inner().Jar.Cookies(runningCtx.StoreUrl),
		Proxy:       &p,
		RequestedAt: time.Now().Unix(),
	})

	if err != nil {
		taskCtx.SendStatusColored("Error getting checkpoint data", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetCheckpoint, err
	}

	s := <-r.Collective

	runningCtx.CheckpointForm = s.Data.Form

	cookies := util.ConvertAstilectronCookies(s.Data.Cookies)

	runningCtx.Client.Inner().Jar.SetCookies(runningCtx.StoreUrl, cookies)

	return SubmitCheckpoint, nil
}

func handleSubmitCheckpointPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(runningCtx.Redirect).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("content-type", "application/x-www-form-urlencoded; charset=UTF-8").
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", runningCtx.StoreUrl.String()+"/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(runningCtx.CheckpointForm).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting checkpoint", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCheckpoint, err
	}

	if resp.StatusCode == 403 {
		taskCtx.SendStatusColored("Invalid checkpoint data", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetCheckpoint, nil
	}

	if resp.StatusCode == 404 {
		if runningCtx.CheckoutHash != "" {
			return GetCustomer, nil
		}

		// todo: is this correct?
		return CreateCheckout, nil
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")
		if strings.Contains(location, "checkpoint") {
			taskCtx.SendStatusColored("Invalid checkpoint data", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetCheckpoint, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "checkouts") {
			runningCtx.Form = nil

			if runningCtx.NextState != "" && !runningCtx.Preloading {
				next := runningCtx.NextState
				runningCtx.NextState = ""

				return next, nil
			}

			if runningCtx.Preloading {
				return ClearPreloadCart, nil
			}

			return GetCustomer, nil
		}

		if strings.Contains(location, "cart") || strings.Contains(location, "checkout") {
			runningCtx.Form = nil

			if !runningCtx.Preloading && runningCtx.CheckoutHash != "" {
				return ClearPreloadCart, nil
			}

			if runningCtx.CheckoutHash != "" {
				return GetCustomer, nil
			}

			return CreateCheckout, nil
		}
	}

	taskCtx.SendStatusColored("Error submitting checkpoint", util.ColorFailed)
	time.Sleep(3 * time.Second)

	return SubmitCheckpoint, nil
}

func handleSubmitCartPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Adding to cart")

	form := generateCartForm(runningCtx.Product.Variant.ID)
	referer := runningCtx.Product.Url

	if runningCtx.Preloading {
		referer = runningCtx.PreloadProduct.ProductUrl
		form = generateCartForm(strconv.FormatInt(runningCtx.PreloadProduct.Variant.ID, 10))
	}

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/cart/add.js", runningCtx.StoreUrl.String())).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("content-type", "application/x-www-form-urlencoded; charset=UTF-8").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "application/json, text/javascript, */*; q=0.01").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", referer).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		log.Println(err)
		taskCtx.SendStatusColored("Error adding to cart", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCart, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")
		if strings.Contains(location, "checkpoint") {
			runningCtx.Redirect = location
			runningCtx.PrevState = SubmitCart
			return GetCheckpoint, nil
		}
	}

	if resp.StatusCode != 200 {
		time.Sleep(3 * time.Second)
		return SubmitCart, nil
	}

	if strings.Contains(runningCtx.StoreUrl.String(), "palaceskate") {
		return GetCart, nil
	}

	// if the task added to cart the wanted product and the checkout url has been found then proceed to the customer page
	if !runningCtx.Preloading && runningCtx.CheckoutHash != "" {
		var body SubmitCartResponse
		if err = resp.JsonBody(&body); err != nil {
			taskCtx.SendStatusColored("Error parsing cart, retrying...", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetConfig, nil
		}

		runningCtx.Product.Url = runningCtx.StoreUrl.String() + body.URL
		runningCtx.Product.Size = body.VariantTitle
		runningCtx.Product.Image = body.Image
		return GetCustomer, nil
	}

	// reset preloading for the real atc
	runningCtx.Preloading = false
	return CreateCheckout, nil
}

func handleClearPreloadCart(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Clearing cart")
	form := url.Values{}

	toggleRedirect(false, runningCtx)

	_, err := runningCtx.Client.Builder(fmt.Sprintf("%s/cart/clear?line=1&quantity=0", runningCtx.StoreUrl.String())).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("content-type", "application/x-www-form-urlencoded; charset=UTF-8").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "application/json, text/javascript, */*; q=0.01").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", runningCtx.PreloadProduct.ProductUrl).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error clearing cart", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return ClearPreloadCart, err
	}

	return WaitForProduct, nil
}

func handleCreateCheckoutPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Creating checkout")

	checkoutForm := make(url.Values)
	checkoutForm.Add("updates[]", "1")
	checkoutForm.Add("attributes[checkout_clicked]", "true")
	checkoutForm.Add("checkout", "")

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/cart").
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", runningCtx.PreloadProduct.ProductUrl).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(checkoutForm).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error initializing checkout", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return CreateCheckout, err
	}

	if resp.StatusCode != 302 {
		taskCtx.SendStatusColored(fmt.Sprintf("Error getting checkout url [%d]", resp.StatusCode), util.ColorFailed)
		time.Sleep(3 * time.Second)
		return CreateCheckout, nil
	}

	location := resp.Header.Get("location")

	if strings.Contains(location, "checkpoint") {
		runningCtx.Redirect = location
		return GetCheckpoint, nil
	}

	if strings.Contains(location, "throttle") {
		return EnterQueue, nil
	}

	if strings.Contains(location, "password") {
		if runningCtx.UsesPassword {
			runningCtx.NextState = SubmitCustomer
			return SubmitPassword, nil
		}
		// If the task doesn't have the site password then keep trying to make a checkout
		return CreateCheckout, nil
	}

	runningCtx.CheckoutHash, err = extractCheckoutHash(location)
	if err != nil {
		taskCtx.SendStatusColored("Error getting checkout hash", util.ColorWarning)
		time.Sleep(3 * time.Second)
		return CreateCheckout, nil
	}

	return ClearPreloadCart, nil
}

func handleGetCustomerPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting customer")

	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", runningCtx.StoreUrl.String()+"/cart").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting customer", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetCustomer, err
	}

	// Handle redirect possibilities
	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "checkpoint") {
			runningCtx.Redirect = location
			return GetCheckpoint, nil
		}

		if strings.Contains(location, "password") {
			if runningCtx.UsesPassword {
				runningCtx.NextState = GetCustomer
				return SubmitPassword, nil
			}
			time.Sleep(3 * time.Second)
			return GetCustomer, nil
		}

		if strings.Contains(location, "account") {
			return GetAccount, nil
		}

		if strings.Contains(location, "stock_problems") {
			taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetCustomer, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}
	}

	runningCtx.AuthToken = extractAuthToken(resp.StringBody())

	// customer page: form:has(input[id=checkout_shipping_address_first_name]
	form, err := extractProtection(strings.NewReader(resp.StringBody()), "form:has(input[id=checkout_shipping_address_first_name])", runningCtx)
	if err != nil {
		taskCtx.SendStatusColored("Error parsing form, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetCustomer, err
	}
	runningCtx.Form = form

	return SubmitCustomer, nil
}

func handleSubmitCustomerPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting customer")

	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	nameParts := strings.Split(runningCtx.ProfileData.Billing.Name, " ")
	if len(nameParts) <= 1 {
		taskCtx.SendStatusColored("Invalid profile data", util.ColorFailed)
		return "", nil
	}
	// only pass in the running ctx pointer
	payload := generateSafeCustomerForm(runningCtx)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", endpoint).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		StringBody(payload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting customer", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCustomer, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetCustomer, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "previous_step=contact_information") {
			return GetShipping, nil
		}
	}
	// todo extract next step

	taskCtx.SendStatusColored("Error submitting customer", util.ColorFailed)
	return GetCustomer, nil
}

func handleGetShippingPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	baseUrl := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)
	endpoint := ""

	if runningCtx.PollingUrl {
		taskCtx.SendStatus("Polling rates")
		endpoint = baseUrl + "/shipping_rates?step=shipping_method"
	} else {
		taskCtx.SendStatus("Visiting rates")
		endpoint = baseUrl + "?previous_step=contact_information&step=shipping_method"
	}

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", baseUrl+"?step=contact_information").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error visiting rates", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetShipping, err
	}

	runningCtx.PollingUrl = false

	if resp.StatusCode == 202 {
		runningCtx.PollingUrl = true
		time.Sleep(1 * time.Second)
		return GetShipping, nil
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return GetShipping, nil
		}

	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.StringBody()))
	if err != nil {
		taskCtx.SendStatusColored("Error parsing rates", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetShipping, err
	}

	shippingRate, rateExists := doc.Find("div.content-box__row .radio-wrapper").Attr("data-shipping-method")

	if !rateExists {
		taskCtx.SendStatusColored("No rates available", util.ColorWarning)
		time.Sleep(1 * time.Second)
		return GetShipping, nil
	}

	if strings.Contains(resp.StringBody(), "order can’t be shipped to your location") {
		taskCtx.SendStatusColored("Unsupported country", util.ColorWarning)
		return "", nil
	}

	if strings.Contains(resp.StringBody(), "no shipping methods") {
		taskCtx.SendStatusColored("No rates available", util.ColorWarning)
		return "", nil
	}

	// Getting available shipping
	if strings.Contains(resp.StringBody(), "Getting available shipping") {
		runningCtx.PollingUrl = true
		taskCtx.SendStatusColored("No rates available", util.ColorWarning)
		time.Sleep(1 * time.Second)
		return GetShipping, nil
	}

	runningCtx.AuthToken = extractAuthToken(resp.StringBody())
	runningCtx.ShippingRate = shippingRate

	form, err := extractProtection(strings.NewReader(resp.StringBody()), "form:has(input[id^=checkout_shipping_rate_id])", runningCtx)
	if err != nil {
		taskCtx.SendStatusColored("Error parsing rates", util.ColorWarning)
		return GetShipping, nil
	}

	runningCtx.Form = form

	return SubmitShipping, nil
}

func handleSubmitShippingPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting rates")

	payload := generateSafeShippingForm(runningCtx)
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("referer", endpoint+"?previous_step=contact_information&step=shipping_method").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		StringBody(payload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting rates", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetShipping, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			// should the task be changed to fast mode?
			return GetPayment, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "previous_step=shipping_method") {
			return GetPayment, nil
		}
	}
	// todo extract next stepBody())

	return GetShipping, nil
}

func handleGetPaymentPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	baseUrl := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)
	endpoint := ""
	referer := ""

	if runningCtx.PollingUrl {
		taskCtx.SendStatus("Calculating taxes")
		endpoint = baseUrl + "?step=payment_method"
		referer = baseUrl + "?previous_step=shipping_method&step=payment_method"
	} else {
		taskCtx.SendStatus("Visiting payment")
		endpoint = baseUrl + "?previous_step=shipping_method&step=payment_method"
		referer = baseUrl + "?previous_step=contact_information&step=shipping_method"
	}

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", referer).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error visiting payment", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPayment, err
	}

	runningCtx.PollingUrl = false

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			// should the task be changed to fast mode?
			return GetPayment, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "checkpoint") {
			runningCtx.NextState = GetPayment
			return GetCheckpoint, nil
		}
	}

	if regexp.MustCompile("Calculating taxes").MatchString(resp.StringBody()) {
		runningCtx.PollingUrl = true
		return GetPayment, nil
	}

	form, err := extractProtection(strings.NewReader(resp.StringBody()), "form:has(input[id=s])", runningCtx)
	if err != nil {
		taskCtx.SendStatusColored("Error parsing payment form", util.ColorWarning)
		return GetShipping, nil
	}

	runningCtx.AuthToken = extractAuthToken(resp.StringBody())

	gateway, err := extractPaymentGateway(strings.NewReader(resp.StringBody()))
	if err != nil {
		taskCtx.SendStatusColored("Error parsing gateway", util.ColorWarning)
		return GetShipping, nil
	}

	price, err := extractPrice(strings.NewReader(resp.StringBody()))
	if err != nil {
		taskCtx.SendStatusColored("Error parsing gateway", util.ColorWarning)
		return GetShipping, nil
	}

	runningCtx.Product.Price = price
	runningCtx.PaymentGateway = gateway
	runningCtx.Form = form
	// if paypal get paypal, get session slug default
	return GetPaymentSession, nil
}

func handleGetPaymentSessionPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting session")
	payload := PaymentSessionPayload{
		CreditCard: PaymentSessionCardData{
			Number:            runningCtx.ProfileData.Payment.Number,
			Name:              runningCtx.ProfileData.Payment.Name,
			Month:             runningCtx.ProfileData.Payment.ExpMonth,
			Year:              runningCtx.ProfileData.Payment.ExpYear,
			VerificationValue: runningCtx.ProfileData.Payment.Cvv,
		},
		PaymentSessionScope: runningCtx.StoreUrl.Host,
	}

	resp, err := runningCtx.Client.Builder("https://deposit.us.shopifycs.com/sessions").
		Method(http.MethodPost).
		Header("accept", "application/json").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("content-type", "application/json").
		Header("dnt", "1").
		Header("origin", "https://checkout.shopifycs.com").
		Header("pragma", "no-cache").
		Header("referer", "https://checkout.shopifycs.com/").
		Header("user-agent", runningCtx.UserAgent).
		JsonBody(payload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting payment session", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	if resp.StatusCode != 200 {
		taskCtx.SendStatusColored("Error getting payment session", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, nil
	}

	var body PaymentSessionResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing payment session, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, nil
	}

	if body.Id == "" {
		taskCtx.SendStatusColored("Invalid payment session", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, nil
	}
	runningCtx.PaymentSessionId = body.Id

	//if runningCtx.Checked && runningCtx.Restocking {
	//	return GetPayment, nil
	//}

	return SubmitPayment, nil
}

func handleSubmitPaymentPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting payment")

	payload := generateSafePaymentForm(runningCtx)
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("content-type", "application/x-www-form-urlencoded").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("referer", endpoint+"?previous_step=payment_method&step=payment_method").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		StringBody(payload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting rates", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetShipping, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "processing") {
			return GetOrder, nil
		}
	}

	return GetPayment, nil
}

func handleGetOrderPreload(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Processing...")
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	if runningCtx.PollingUrl {
		endpoint = endpoint + "/processing?from_processing_page=1"
	} else {
		endpoint = endpoint + "/processing"
	}

	toggleRedirect(false, runningCtx)
	resp, err := runningCtx.Client.Builder(endpoint).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("content-type", "application/x-www-form-urlencoded").
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("content-type", "application/x-www-form-urlencoded").
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("referer", fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)).
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting order", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	runningCtx.PollingUrl = false

	if resp.StatusCode != 200 {
		if resp.StatusCode == 302 {
			location := resp.Header.Get("location")
			log.Println(location)

			if strings.Contains(location, "validate=true") || strings.Contains(location, "checkoutnow") {
				taskCtx.SendStatusColored("Declined", util.ColorFailed)
				// send declined webhook
				addWebhook(taskCtx, runningCtx, staticCtx, webhooks.OrderStatusFailed, "Preload")
				resetRunningCtx(runningCtx)

				return "", nil
			}

			if strings.Contains(location, "thank_you") {
				// send success webhook
				taskCtx.SendStatusColored("Success", util.ColorSuccess)
				addWebhook(taskCtx, runningCtx, staticCtx, webhooks.OrderStatusPlaced, "Preload")
				return "", nil
			}
		}
	}

	if strings.Contains(resp.StringBody(), "out of stock") {
		if runningCtx.Restocking {
			return GetConfig, nil
		}
		return GetHomepage, nil
	}

	return GetOrder, nil
}
