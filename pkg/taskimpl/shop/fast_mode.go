package shopify

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
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

	"github.com/google/uuid"
)

// handleSetup Sets all the global variables used throughout the module
func handleSetup(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	// These variables are supposed to come from the task ctx
	runningCtx.UsesPassword = false
	runningCtx.StorePassword = "supersecretpassword123"
	runningCtx.UsesAccount = false
	runningCtx.AccountCredentials = ""

	runningCtx.UserAgent = getUserAgent()

	runningCtx.Restocking = false
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

	return GetHomepage, nil
}

func handleEnterQueueFast(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

func handleGetQueueFast(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

func handleGetPastQueueFast(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

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
			runningCtx.Redirect = u.String()
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

func handleGetNextQueueFast(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

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

// handleWaitForProduct Starts a new monitor or retrieves a saved product from the monitor pool
func handleWaitForProduct(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

	return SubmitCustomer, nil
}

// handleGetHomepage Gets the site's homepage
func handleGetHomepage(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting homepage")

	runningCtx.StoreUrl, _ = url.Parse(staticCtx.Default.Store.Url)

	_, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()).
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

	return GetConfig, nil
}

// handleGetConfig Gets a site configuration which has the access token and the shop ID needed for checkout
func handleGetConfig(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting config")

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/payments/config").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "application/json").
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
		return SubmitPassword, nil
	}

	if runningCtx.UsesAccount {
		return GetAccount, nil
	}

	return CreateCheckout, nil
}

// handleSubmitPassword Submits the storefront password
func handleSubmitPassword(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting password")

	formValues := make(url.Values)
	formValues.Add("form_type", "storefront_password")
	formValues.Add("utf8", "✓")
	formValues.Add("password", runningCtx.StorePassword)

	// Stops the request from being redirected
	runningCtx.Client.Inner().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := runningCtx.Client.Builder(runningCtx.StoreUrl.String()+"/password").
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
		FormBody(formValues).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting password", util.ColorFailed)
		log.Println(err)
		return SubmitPassword, err
	}

	/** If the password handler was called from another handler
	then go to the next desired state instead of the default state to follow **/
	if runningCtx.NextState != "" {
		return runningCtx.NextState, nil
	}

	log.Println(resp.StatusCode, resp.Header)
	return CreateCheckout, nil
}

// handleGetAccount Gets a site's account page
func handleGetAccount(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

func handleGetCaptchaToken(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

func handleSubmitAccount(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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
			return CreateCheckout, nil
		}

	}
	taskCtx.SendStatusColored("Error submitting account", util.ColorFailed)
	time.Sleep(2 * time.Second)
	return SubmitAccount, nil
}

func handleGetChallenge(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

func handleSubmitChallenge(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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
			return CreateCheckout, nil
		}
	}
	taskCtx.SendStatusColored("Error submitting challenge", util.ColorFailed)
	return GetChallenge, nil
}

// handleCreateCheckout
// Creates a checkout url by sending a post request with an empty cart payload to a site's cart endpoint
func handleCreateCheckout(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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
		Header("referer", runningCtx.StoreUrl.String()+"/").
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
	log.Println(location)
	if strings.Contains(location, "checkpoint") {
		taskCtx.SendStatusColored("Checkpoint up, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return CreateCheckout, nil
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

	toggleRedirect(true, runningCtx)

	return WaitForProduct, nil
}

// handleSubmitCart
// Adds a product to cart via the frontend cart endpoint
func handleSubmitCart(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting cart")

	form := generateCartForm(runningCtx.Product.Variant.ID)

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/cart/add.js", runningCtx.StoreUrl.String())).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("content-type", "application/x-www-form-urlencoded; charset=UTF-8").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "application/json, text/javascript, */*; q=0.01").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-user", "?1").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting cart", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCart, err
	}

	// if status == 404 then the variant is not live
	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Variant not live", util.ColorWarning)
		time.Sleep(3 * time.Second)
		return SubmitCart, nil
	}

	var body SubmitCartResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing cart, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCart, nil
	}

	runningCtx.Product.Url = runningCtx.StoreUrl.String() + body.URL
	runningCtx.Product.Size = body.VariantTitle
	runningCtx.Product.Image = body.Image

	// if palace skateboards then get cart to extract the data needed to checkout?
	if strings.Contains(staticCtx.Default.Store.Name, "palaceskateboards") {
		return GetCart, nil
	}

	return GetCustomer, nil
}

// handleSubmitCustomer
// Sends all the profile data to the backend api endpoint, this payload can also include the product's variant as a line item
func handleSubmitCustomer(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Preloading checkout")

	nameParts := strings.Split(runningCtx.ProfileData.Billing.Name, " ")
	if len(nameParts) <= 1 {
		taskCtx.SendStatusColored("Invalid profile data", util.ColorFailed)
		return "", nil
	}
	checkoutProperties := map[string]string{}

	if strings.Contains(staticCtx.Default.Store.Name, "kawsone") {
		checkoutProperties["_hkey"] = "789gfd78934hfk74jml849320"
	}

	customerData := CustomerDataApiForm{Checkout: CustomerDataApiCheckoutForm{
		Email:  runningCtx.ProfileData.Payment.Email,
		Secret: true,
		ShippingAddress: CustomerDataApiAddress{
			FirstName:    nameParts[0],
			LastName:     nameParts[1],
			Address1:     runningCtx.ProfileData.Shipping.Line1,
			Address2:     runningCtx.ProfileData.Shipping.Line2,
			CountryCode:  runningCtx.ProfileData.Shipping.Country.Code,
			ProvinceCode: runningCtx.ProfileData.Shipping.State.Code,
			City:         runningCtx.ProfileData.Shipping.City,
			Zip:          runningCtx.ProfileData.Shipping.PostCode,
			Phone:        runningCtx.ProfileData.Payment.Phone,
		},
		BillingAddress: CustomerDataApiAddress{
			FirstName:    nameParts[0],
			LastName:     nameParts[1],
			Address1:     runningCtx.ProfileData.Shipping.Line1,
			Address2:     runningCtx.ProfileData.Shipping.Line2,
			CountryCode:  runningCtx.ProfileData.Shipping.Country.Code,
			ProvinceCode: runningCtx.ProfileData.Shipping.State.Code,
			City:         runningCtx.ProfileData.Shipping.City,
			Zip:          runningCtx.ProfileData.Shipping.PostCode,
			Phone:        runningCtx.ProfileData.Payment.Phone,
		},
		LineItems: []LineItemsApiForm{},
	}}
	// if the task is not restocking and not falling back (which should be the first run) add the variant to the checkout
	if !runningCtx.Restocking && !runningCtx.Fallback {
		v, _ := strconv.ParseInt(runningCtx.Product.Variant.ID, 10, 64)
		customerData.Checkout.LineItems = append(customerData.Checkout.LineItems, LineItemsApiForm{
			VariantID:  v,
			Quantity:   1,
			Properties: checkoutProperties,
		})
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(runningCtx.AccessToken + "::"))

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/api/checkouts/%s.json", runningCtx.StoreUrl.String(), runningCtx.CheckoutHash)).
		Method(http.MethodPatch).
		Header("authorization", fmt.Sprintf("Basic %s", authToken)).
		Header("content-type", "application/json").
		Header("X-Shopify-Checkout-Version", "2021-01-04").
		Header("X-Shopify-UniqueToken", uuid.New().String()).
		Header("X-Shopify-VisitToken", uuid.New().String()).
		JsonBody(customerData).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error initializing checkout, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCustomer, nil
	}

	if !isImproperStatusCode(resp.StatusCode) {
		taskCtx.SendStatusColored("Error preloading checkout, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCustomer, nil
	}

	var body interface{}
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error preloading checkout, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, nil
	}
	// todo: handle invalid zip code in profile
	//data, _ := body.(CustomerDataResponse)

	if resp.StatusCode == 422 {
		stringBody := resp.StringBody()
		if strings.Contains(stringBody, "zip") && strings.Contains(stringBody, "invalid") {
			taskCtx.SendStatusColored("Invalid zip code", util.ColorFailed)
			return "", nil
		}
		if strings.Contains(stringBody, "line_items") && strings.Contains(stringBody, "invalid") {
			log.Println("falling back")
			runningCtx.Fallback = true
			return SubmitCustomer, nil
		}

		if strings.Contains(stringBody, "domain") {
			taskCtx.SendStatusColored("Invalid email domain", util.ColorFailed)
			return SubmitCustomer, nil
		}

		if strings.Contains(stringBody, "not_enough_in_stock") {
			log.Println("product oos while submitting customer data")
			runningCtx.Restocking = true
			if strings.Contains(staticCtx.Default.Store.Name, "eflash-us") {
				return SubmitCart, nil
			}
			log.Println("retrying with oos product...")
			return SubmitCustomer, nil
		}
		time.Sleep(3 * time.Second)
		return SubmitCustomer, err
	}

	//log.Println(len(data.Checkout.LineItems))
	//if len(data.Checkout.LineItems) == 0 {
	//	log.Println("ell tah")
	//	return SubmitCart, nil
	//}

	if strings.Contains(staticCtx.Default.Store.Name, "eflash-us") || strings.Contains(staticCtx.Default.Store.Name, "palaceskateboards") {
		if runningCtx.ShippingRate != "" {
			return SubmitShipping, nil
		}
		return GetShipping, nil
	}

	if runningCtx.Restocking || runningCtx.Fallback {
		log.Println("activating restock mode...")
		return SubmitCart, nil
	}

	return GetCustomer, nil
}

// handleGetCustomer
// Gets the customer information frontend page, this is used to extract the auth token needed for checkout
// and also checks what state to go to depending on the task's configuration
func handleGetCustomer(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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
		log.Println("being redirected", resp.Header.Get("location"))
		location := resp.Header.Get("location")

		if strings.Contains(location, "checkpoint") {
			taskCtx.SendStatusColored("Checkpoint up, retrying...", util.ColorWarning)
			time.Sleep(3 * time.Second)
			return GetCustomer, nil
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
			runningCtx.Restocking = true
			return GetShipping, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}
	}

	authToken := extractAuthToken(resp.StringBody())

	runningCtx.AuthToken = authToken

	if runningCtx.ShippingRate != "" {
		if runningCtx.Restocking {
			return SubmitShipping, nil
		}
		// todo: handle paypal mode > get payment page
		return GetPaymentSession, nil
	}
	return GetShipping, nil
}

// TODO
func handleGetCart(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting cart")
	return "", nil
}

// handleGetShipping
// Gets the task's shipping rates
func handleGetShipping(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting shipping")
	authToken := base64.StdEncoding.EncodeToString([]byte(runningCtx.AccessToken))

	resp, err := runningCtx.Client.Builder(fmt.Sprintf("%s/api/checkouts/%s/shipping_rates.json", runningCtx.StoreUrl.String(), runningCtx.CheckoutHash)).
		Header("authorization", fmt.Sprintf("Basic %s", authToken)).
		Header("content-type", "application/json").
		Header("X-Shopify-Checkout-Version", "2021-01-04").
		Header("X-Shopify-UniqueToken", uuid.New().String()).
		Header("X-Shopify-VisitToken", uuid.New().String()).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting shipping", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetShipping, err
	}

	var body ShippingApiResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing config, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetConfig, nil
	}

	// NOTE: Possibly a free order, let's proceed to payment session
	if resp.StatusCode == 412 {
		runningCtx.IsFree = true
		return GetPaymentSession, nil
	}

	if len(body.ShippingRates) == 0 {
		taskCtx.SendStatus("Polling rates...")
		time.Sleep(500 * time.Millisecond)
		return GetShipping, nil
	}

	runningCtx.ShippingRate = body.ShippingRates[0].ID
	runningCtx.Product.Price = body.ShippingRates[0].Price
	//if resp.StatusCode == 200 && len(body.ShippingRates) == 0 {
	//
	//}

	//todo: parse price from checkout and check the task's max price

	if runningCtx.Restocking {
		return SubmitShipping, nil
	}

	// todo: handle paypal mode > get payment

	return GetPaymentSession, nil
}

// handleSubmitShipping
// Submits the task's shipping rate to the backend api endpoint
func handleSubmitShipping(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting shipping")

	payload := map[string]map[string]map[string]string{
		"checkout": {
			"shipping_line": {
				"handle": runningCtx.ShippingRate,
			},
		},
	}
	authToken := base64.StdEncoding.EncodeToString([]byte(runningCtx.AccessToken + "::"))

	_, err := runningCtx.Client.Builder(fmt.Sprintf("%s/api/checkouts/%s.json", runningCtx.StoreUrl.String(), runningCtx.CheckoutHash)).
		Method(http.MethodPatch).
		Header("authorization", fmt.Sprintf("Basic %s", authToken)).
		Header("content-type", "application/json").
		Header("X-Shopify-Checkout-Version", "2021-01-04").
		Header("X-Shopify-UniqueToken", uuid.New().String()).
		Header("X-Shopify-VisitToken", uuid.New().String()).
		JsonBody(payload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error initializing checkout, retrying...", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return SubmitCustomer, nil
	}

	//log.Println(resp.StatusCode)
	// todo: handle paypal mode > get payment
	return GetPaymentSession, nil
}

// handleGetPayment
// Gets the payment frontend page, this is used to extract the auth token needed for checkout
func handleGetPayment(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting payment page")
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(true, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint+"?previous_step=shipping_method&step=payment_method").
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
		Header("referer", endpoint+"?previous_step=shipping_method&step=payment_method").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error getting payment session", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	// does the task really need to handle being redirected here?

	//if resp.StatusCode == 302 {
	//	log.Println("being redirected after payment")
	//}

	runningCtx.AuthToken = extractAuthToken(resp.StringBody())

	price, err := extractPrice(strings.NewReader(resp.StringBody()))
	if err != nil {
		taskCtx.SendStatusColored("Error parsing price", util.ColorFailed)
		return GetPayment, err
	}

	runningCtx.Product.Price = price

	return SubmitPayment, nil
}

// handleGetPaymentSession
// Submits the task's payment data and sets a new payment session token
func handleGetPaymentSession(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Getting payment session")
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

	if runningCtx.Checked && runningCtx.Restocking {
		return GetPayment, nil
	}

	return SubmitPayment, nil
}

// handleSubmitPayment
// Submits the payment form to the frontend payment endpoint
func handleSubmitPayment(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting order")

	form := generatePaymentForm(runningCtx.ProfileData, runningCtx.ShippingRate, runningCtx.Product.Price, strconv.Itoa(runningCtx.ShopID), runningCtx.AuthToken, runningCtx.PaymentSessionId)
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("content-type", "application/x-www-form-urlencoded").
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
		Header("referer", fmt.Sprintf("%s?previous_step=shipping_method&step=payment_method", endpoint)).
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting order", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	if (resp.StatusCode == 429 || resp.StatusCode == 430) && strings.Contains(resp.StringBody(), "too many requests") {
		runningCtx.RetryCount += 1

		if runningCtx.RetryCount >= 5 {
			taskCtx.SendStatusColored("Checkout expired", util.ColorWarning)
			resetRunningCtx(runningCtx)
			return Setup, nil
		}
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "checkpoint") {
			taskCtx.SendStatusColored("Checkpoint up, retrying...", util.ColorFailed)
			time.Sleep(3 * time.Second)
			return CreateCheckout, nil
		}

		if strings.Contains(location, "stock_problems") {
			if !runningCtx.Checked {
				runningCtx.Checked = true
				return SubmitPayment, nil
			}
			taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
			time.Sleep(3000 * time.Millisecond)
			return GetPaymentSession, nil
		}
		// handle account login

		if strings.Contains(location, "previous_step=payment_method") {
			if !runningCtx.Checked {
				runningCtx.Checked = true
				return SubmitPayment, nil
			}

			taskCtx.SendStatusColored("Out of stock", util.ColorFailed)
			time.Sleep(3000 * time.Millisecond)
			return GetPaymentSession, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "processing") {

			return GetOrder, nil
		}
	}

	if runningCtx.Checked {
		return GetPaymentSession, nil
	}

	body := resp.StringBody()

	if regexp.MustCompile("Calculating taxes").MatchString(body) {
		taskCtx.SendStatus("Calculating taxes")
		return GetReview, nil
	}

	// free order
	if regexp.MustCompile("No payment is required").MatchString(body) {
		runningCtx.IsFree = true
		return SubmitReview, nil
	}
	/*
		todo: if body includes captcha class > extract recaptcha data > update auth token >
	*/

	// todo: extract the current page the task is in (customer, review, shipping, etc...) and proceed to next state accordingly

	return GetPaymentSession, nil
}

// handleGetReview
// Gets the order review page, this is used to extract a new auth token for the checkout, calculate taxes
// and also handle multiple redirect possibilities
func handleGetReview(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint+"?step=review").
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
		taskCtx.SendStatusColored("Error getting review", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			if runningCtx.Restocking {
				return GetPaymentSession, nil
			}
			time.Sleep(3 * time.Second)
			return GetReview, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}
	}

	body := resp.StringBody()

	if regexp.MustCompile("Calculating taxes").MatchString(body) {
		taskCtx.SendStatus("Calculating taxes")
		return GetReview, nil
	}

	// get product price here?
	runningCtx.AuthToken = extractAuthToken(body)
	runningCtx.Product.Price = extractProductFinalPrice(body)

	return SubmitReview, nil
}

// handleSubmitReview
// Submits the order review, this can be considered as the final step in the checkout, also handles multiple redirect possibilities
func handleSubmitReview(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Completing order")
	form := generateReviewForm(runningCtx.AuthToken, runningCtx.Product.Price, runningCtx.IsFree, runningCtx.CaptchaToken)

	endpoint := fmt.Sprintf("%s/%d/checkouts/%s", runningCtx.StoreUrl.String(), runningCtx.ShopID, runningCtx.CheckoutHash)

	toggleRedirect(false, runningCtx)

	resp, err := runningCtx.Client.Builder(endpoint).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("content-type", "application/x-www-form-urlencoded").
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
		Header("referer", endpoint).
		Header("accept-language", "en-US,en;q=0.9").
		FormBody(form).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting review", util.ColorFailed)
		time.Sleep(3 * time.Second)
		return GetPaymentSession, err
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")

		if strings.Contains(location, "stock_problems") {
			time.Sleep(2 * time.Second)
			return SubmitReview, nil
		}

		if strings.Contains(location, "throttle") {
			return EnterQueue, nil
		}

		if strings.Contains(location, "processing") {
			return GetOrder, nil
		}
	}

	return "", nil
}

func handleGetOrder(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
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

	if resp.StatusCode != 200 {
		if resp.StatusCode == 302 {
			location := resp.Header.Get("location")
			log.Println(location)

			if strings.Contains(location, "validate=true") || strings.Contains(location, "checkoutnow") {
				taskCtx.SendStatusColored("Declined", util.ColorFailed)
				// send declined webhook
				addWebhook(taskCtx, runningCtx, staticCtx, webhooks.OrderStatusFailed, "Fast")
				resetRunningCtx(runningCtx)

				return GetConfig, nil
			}

			if strings.Contains(location, "thank_you") {
				// send success webhook
				taskCtx.SendStatusColored("Success", util.ColorSuccess)
				addWebhook(taskCtx, runningCtx, staticCtx, webhooks.OrderStatusPlaced, "Fast")

				return "", nil
			}
		}
	}

	runningCtx.PollingUrl = false

	if strings.Contains(resp.StringBody(), "out of stock") {
		if runningCtx.Restocking {
			return GetConfig, nil
		}
		return GetHomepage, nil
	}

	return GetOrder, nil
}
