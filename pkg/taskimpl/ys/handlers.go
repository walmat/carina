package ys

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/integrations/autosolve"
	"nebula/pkg/taskutil"
	"nebula/pkg/util"
	"nebula/third_party/adyen2"
	"nebula/third_party/http"
	"net/url"
	"strings"
	"time"
)

func resetContext(staticCtx *staticCtx, runningCtx *runningCtx) {
	runningCtx.UserAgent = getUserAgent()
	runningCtx.Authorization = "null"
	runningCtx.BasketId = ""
	runningCtx.Size = ""
	runningCtx.Sku = ""
	runningCtx.PixelScript = ""
	runningCtx.PixelValue = ""
	runningCtx.PixelHash = ""
	runningCtx.PixelPayload = ""

	runningCtx.ProductId = staticCtx.Default.Sku
	runningCtx.Referer = fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)
	runningCtx.IsFirstBasket = true
	runningCtx.FirewallRetries = 0
	runningCtx.AbckRetries = 0
}

func resetRetries(runningCtx *runningCtx) {
	runningCtx.FirewallRetries = 0
	runningCtx.AbckRetries = 0
}

func handleGetHomepage(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting homepage")

	runningCtx.ProductId = staticCtx.Default.Sku
	runningCtx.Referer = fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)
	runningCtx.IsFirstBasket = true

	resp, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/").
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
		taskCtx.SendStatusColored("Error visiting homepage, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetHomepage, nil
	}

	banned := isBanned(runningCtx, resp.StringBody(), 300000)
	if banned == true {
		taskCtx.SendStatusColored("Proxy banned, resetting...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		_ = rotateProxy(taskCtx, runningCtx)
		resetContext(staticCtx, runningCtx)

		return GetHomepage, nil
	}

	extractAkamaiScript(runningCtx, resp.StringBody())

	hasPixel := checkPixel(resp, runningCtx)
	if hasPixel {
		runningCtx.NextState = GetAkamai
		runningCtx.ReturnState = GetHomepage
		return GetPixel, nil
	}

	return GetAkamai, nil
}

func handleGetPixel(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting pixel")

	resp, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/akam/11/%s", runningCtx.PixelScript)).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "*/*").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "document").
		Header("referer", "https://www.yeezysupply.com/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting pixel, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetPixel, nil
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Pixel not found, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		if runningCtx.ReturnState != "" {
			return runningCtx.ReturnState, nil
		}

		resetContext(staticCtx, runningCtx)

		return GetHomepage, nil
	}

	PixelValue, PixelValueError := extractPixelValue(resp.StringBody())
	if PixelValueError != nil {
		taskCtx.SendStatusColored("Nil pixel value, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetPixel, nil
	}

	runningCtx.PixelValue = PixelValue

	if runningCtx.NextState != "" {
		nextState := runningCtx.NextState
		runningCtx.NextState = ""
		return nextState, nil
	}

	return GetPayload, nil
}

func handleGetAkamai(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting script")

	_, err := runningCtx.Client.
		Builder(runningCtx.AkamaiUrl).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "*/*").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "no-cors").
		Header("sec-fetch-dest", "script").
		Header("referer", "https://www.yeezysupply.com/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting script, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetAkamai, nil
	}

	return GetBloom, nil
}

func handleGetBloom(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting bloom")

	_, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/api/yeezysupply/products/bloom").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "*/*").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "document").
		Header("referer", "https://www.yeezysupply.com/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting bloom, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetBloom, nil
	}

	runningCtx.NextState = GetPayload
	runningCtx.ReturnState = ""
	return GetBasket, nil
}

func handleGetBasket(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting basket")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return GetBasket, nil
	}

	resp, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/api/checkout/customer/baskets").
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("checkout-authorization", runningCtx.Authorization).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting basket, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetBasket, nil
	}

	if runningCtx.IsFirewallBlocked {
		return GetProductStock, nil
	}

	if runningCtx.IsFirstBasket {
		runningCtx.IsFirstBasket = false

		if runningCtx.NextState != "" {
			nextState := runningCtx.NextState
			runningCtx.NextState = ""
			return nextState, nil
		}

		return GetPayload, nil
	}

	if resp.StatusCode == 400 && strings.Contains(strings.ToLower(resp.StringBody()), "invalid url") {
		if runningCtx.AbckRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		taskCtx.SendStatusColored("Akamai block, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.AbckRetries++
		runningCtx.NeedsAuth = true
		return GetSensor, nil
	}

	resetRetries(runningCtx)

	Authorization := extractAuthorizationHeader(resp)
	if Authorization != "" {
		runningCtx.Authorization = Authorization
	}

	if runningCtx.NextState != "" {
		nextState := runningCtx.NextState
		runningCtx.NextState = ""
		return nextState, nil
	}

	runningCtx.NextState = AddToCart
	return GetSensor, nil
}

func handleGetPayload(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Generating payload")

	values := make(url.Values)
	values.Set("scriptVal", runningCtx.PixelValue)
	values.Set("ua", runningCtx.UserAgent)
	values.Set("pixelID", runningCtx.PixelHash)
	values.Set("mode", "PIXEL")
	values.Set("key", "nebula-23rwertydrser%E$TWRAW-esfgdhrtsetrhrthdgsefGEHRJYTGD")

	resp, err := http.NewSimpleClient(taskCtx.Context, http.DefaultClientConfig, http.DefaultFingerprint).
		Builder("https://akam-b429.ganeshbot.cloud/Akamai").
		Method(http.MethodPost).
		Header("accept", "*/*").
		Header("content-type", "application/x-www-form-urlencoded").
		FormBody(values).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error generating payload, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetPayload, nil
	}

	if resp.StringBody() == "" {
		taskCtx.SendStatusColored("Nil payload, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetPayload, nil
	}

	Parts := strings.Split(resp.StringBody(), "*")
	if len(Parts) == 0 {
		taskCtx.SendStatusColored("Nil payload, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetPayload, nil
	}

	runningCtx.PixelPayload = Parts[0]

	return SubmitPixel, nil
}

func handleSubmitPixel(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Submitting pixel")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return SubmitPixel, nil
	}

	setExtraCookies(runningCtx)

	_, err = runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/akam/11/pixel_%s", runningCtx.PixelScript)).
		Method(http.MethodPost).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/x-www-form-urlencoded").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("origin", "https://www.yeezysupply.com").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", "https://www.yeezysupply.com/").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		StringBody(runningCtx.PixelPayload).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting pixel, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		return SubmitPixel, nil
	}

	if runningCtx.ReturnState != "" {
		returnState := runningCtx.ReturnState
		runningCtx.ReturnState = ""
		return returnState, nil
	}

	return GetProductPage, nil
}

func handleGetProductPage(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	if runningCtx.ConfigCancelFunc != nil {
		runningCtx.ConfigCancelFunc()
		runningCtx.ConfigCancelFunc = nil
	}

	taskCtx.SendStatus("Visiting product")

	resp, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		Header("dnt", "1").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("sec-fetch-site", "empty").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "none").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("referer", "https://www.yeezysupply.com/").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting product, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetProductPage, nil
	}

	if runningCtx.IsFirewallBlocked {
		return GetProductDetails, nil
	}

	banned := isBanned(runningCtx, resp.StringBody(), -1)
	if banned {
		taskCtx.SendStatusColored("Akamai block, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 5)

		return GetProductPage, nil
	}

	hasPixel := checkPixel(resp, runningCtx)
	if hasPixel {
		runningCtx.NextState = ""
		runningCtx.ReturnState = GetProductPage
		return GetPixel, nil
	}

	hmac := getHmacCookie(runningCtx)
	if hmac != nil {
		taskCtx.SendStatusColored("Passed splash", util.ColorSuccess)
		time.Sleep(time.Second * 5)

		return GetProductDetails, nil
	}

	flashSale := isFlashSale(resp.StringBody(), runningCtx)
	if flashSale == true {
		return GetProductDetails, nil
	}

	return GetWaitingRoom, nil
}

func handleGetConfig(_ *tasks.Context, _ *staticCtx, runningCtx *runningCtx) {
	_, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/hpl/content/yeezy-supply/config/US/waitingRoomConfig.json").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("accept", "application/json, text/plain, */*").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		return
	}
}

func handleGetRelease(_ *tasks.Context, _ *staticCtx, runningCtx *runningCtx) {
	_, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/hpl/content/yeezy-supply/releases/%s/en_US.json", runningCtx.ProductId)).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "document").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		return
	}
}

func handleGetShared(_ *tasks.Context, _ *staticCtx, runningCtx *runningCtx) {
	_, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/hpl/content/yeezy-supply/releases/%s/shared.json", runningCtx.ProductId)).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "document").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		return
	}
}

func handleGetWaitingRoom(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Visiting asset")

	_, err := runningCtx.Client.
		Builder(runningCtx.WaitingRoom).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("sec-ch-ua-mobile", "?0").
		Header("dnt", "1").
		Header("upgrade-insecure-requests", "1").
		Header("user-agent", runningCtx.UserAgent).
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Header("sec-fetch-site", "none").
		Header("sec-fetch-mode", "navigate").
		Header("sec-fetch-dest", "document").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored(err.Error(), util.ColorFailed)
		time.Sleep(1 * time.Second)
		return GetWaitingRoom, nil
	}

	handleGetConfig(taskCtx, staticCtx, runningCtx)
	handleGetRelease(taskCtx, staticCtx, runningCtx)
	handleGetShared(taskCtx, staticCtx, runningCtx)

	if runningCtx.ConfigCancelFunc == nil {
		confCtx, cancelFunc := context.WithCancel(taskCtx.Context)
		runningCtx.ConfigCancelFunc = cancelFunc

		go func() {
			for {
				select {
				case <-confCtx.Done():
					return
				default:
					<-time.After(5 * time.Second)
					handleGetConfig(taskCtx, staticCtx, runningCtx)
				}
			}
		}()
	}

	return WaitInSplash, nil
}

func handleCaptchaRequest(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatusColored("Solving captcha", util.ColorWarning)

	if cachedCap := getCachedCaptchaToken(); cachedCap != nil {
		runningCtx.CaptchaToken = cachedCap

		cookie := http.Cookie{
			Name:  runningCtx.RecaptchaCookie,
			Value: runningCtx.CaptchaToken.Data.Token,
		}

		runningCtx.CaptchaExpires = time.Now().Add(time.Second*110 - (time.Second * time.Duration(time.Now().Unix()-runningCtx.CaptchaToken.CreatedAt)))

		runningCtx.Client.AddCookie(runningCtx.StoreUrl, &cookie)

		return WaitInSplash, nil
	}

	runningCtx.CaptchaShareChan = newCaptchaShareChan()

	go fetchNewToken(taskCtx)

	return CaptchaWaitState, nil
}

func handleCaptchaWait(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	// why? why would you do this?
	if runningCtx.CaptchaShareChan == nil {
		return CaptchaRequestState, nil
	}

	waitCtx, cancelFunc := context.WithTimeout(taskCtx.Context, time.Second*15)
	defer func() {
		cancelFunc()
		runningCtx.CaptchaShareChan = nil
	}()

	select {
	case captchaToken := <-runningCtx.CaptchaShareChan:
		runningCtx.CaptchaToken = &captchaToken
	case <-waitCtx.Done():
		if waitCtx.Err() == context.DeadlineExceeded {
			taskCtx.SendStatusColored("AutoSolve Timeout, Retrying...", util.ColorFailed)
			_ = autosolve.CancelCaptchaRequest(taskCtx.Identifier())
			return CaptchaRequestState, nil
		} else {
			break
		}
	}

	if runningCtx.CaptchaToken != nil {
		cookie := http.Cookie{
			Name:  runningCtx.RecaptchaCookie,
			Value: runningCtx.CaptchaToken.Data.Token,
		}

		runningCtx.CaptchaExpires = time.Now().Add(time.Second*120 - (time.Second * time.Duration(time.Now().Unix()-runningCtx.CaptchaToken.CreatedAt)))

		runningCtx.Client.AddCookie(runningCtx.StoreUrl, &cookie)
	}

	return WaitInSplash, nil
}

func handleWaitInSplash(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	setSplashCookie(runningCtx)
	removeInSplashCookie(runningCtx)

	if time.Now().After(runningCtx.CaptchaExpires) {
		return CaptchaRequestState, nil
	}

	taskCtx.SendStatusColored("In Splash", util.ColorWarning)

	resp, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/__queue/yzysply").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "application/json, text/plain, */*").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("dnt", "1").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()

	if err != nil {
		time.Sleep(time.Second * 3)
		return WaitInSplash, nil
	}

	if resp.StatusCode == 200 {
		removePassedCookies(runningCtx)

		if isFakeSplash(runningCtx) {
			taskCtx.SendStatusColored("Invalid session, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 5)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		return GetProductPage, nil
	}

	time.Sleep(time.Second * 3)

	return WaitInSplash, nil
}

func handleGetSensor(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Generating sensor")

	abck := extractAbckCookie(runningCtx)

	values := make(url.Values)
	values.Set("abck", abck.Value)
	values.Set("ua", runningCtx.UserAgent)
	values.Set("site", "www.yeezysupply.com")
	values.Set("mode", "API")
	values.Set("key", "nebula-23rwertydrser%E$TWRAW-esfgdhrtsetrhrthdgsefGEHRJYTGD")

	resp, err := http.NewSimpleClient(taskCtx.Context, http.DefaultClientConfig, http.DefaultFingerprint).
		Builder("https://akam-b429.ganeshbot.cloud/Akamai").
		Method(http.MethodPost).
		Header("accept", "*/*").
		Header("content-type", "application/x-www-form-urlencoded").
		FormBody(values).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error generating sensor, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetSensor, nil
	}

	if resp.StringBody() == "" {
		taskCtx.SendStatusColored("Error generating sensor, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		return GetSensor, nil
	}

	SensorPayload := taskutil.DecryptAkamaiSensor(resp.StringBody())

	if len(SensorPayload) < 2 {
		taskCtx.SendStatusColored("Error generating sensor, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		return GetSensor, nil
	}

	runningCtx.SensorPayload = SensorPayload[0]
	runningCtx.SensorUserAgent = SensorPayload[1]

	return SubmitSensor, nil
}

func handleSubmitSensor(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Submitting sensor")

	body := fmt.Sprintf(`{"sensor_data":"%s"}`, runningCtx.SensorPayload)

	_, err := runningCtx.Client.
		Builder(runningCtx.AkamaiUrl).
		Method(http.MethodPost).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("content-type", "text/plain;charset=UTF-8").
		Header("accept", "*/*").
		Header("origin", "https://www.yeezysupply.com").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", runningCtx.Referer).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		StringBody(body).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error submitting sensor, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitSensor, nil
	}

	if runningCtx.NumberPosts >= 5 {
		runningCtx.NumberPosts = 0

		if runningCtx.NeedsAuth {
			runningCtx.NeedsAuth = false
			return GetBasket, nil
		}

		if runningCtx.NextState != "" {
			Next := runningCtx.NextState
			runningCtx.NextState = ""

			return Next, nil
		}

		if runningCtx.Authorization == "" {
			return GetBasket, nil
		}

		if runningCtx.BasketId == "" {
			return AddToCart, nil
		}

		if runningCtx.InfoSubmitted == false {
			return SubmitInformation, nil
		}

		return SubmitOrder, nil
	}

	runningCtx.NumberPosts += 1
	return GetSensor, nil
}

func handleGetProductDetails(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Visiting details")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return GetProductDetails, nil
	}

	resp, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/api/products/%s", runningCtx.ProductId)).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("content-type", "application/json").
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("accept", "*/*").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting details, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetProductDetails, nil
	}

	if runningCtx.IsFirewallBlocked {
		return GetBasket, nil
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Product not loaded, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 5)

		return GetProductDetails, nil
	}

	fmt.Println(len(resp.StringBody()))

	if isBanned(runningCtx, resp.StringBody(), -1) {
		taskCtx.SendStatusColored("Akamai block, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 5)

		runningCtx.NextState = GetProductDetails
		return GetSensor, nil
	}

	var body ProductDetailsResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing details, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetProductDetails, nil
	}

	extractProductDetails(runningCtx, body)

	runningCtx.NextState = GetProductStock
	return GetBasket, nil
}

func handleGetProductStock(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Visiting stock")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return GetProductStock, nil
	}

	resp, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/api/products/%s/availability", runningCtx.ProductId)).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("origin", "https://www.yeezysupply.com").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error visiting stock, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return GetProductStock, nil
	}

	if runningCtx.IsFirewallBlocked {
		runningCtx.IsFirewallBlocked = false
		return GetSensor, nil
	}

	fmt.Println(len(resp.StringBody()))

	if isBanned(runningCtx, resp.StringBody(), -1) {
		taskCtx.SendStatusColored("Akamai block, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 5)

		runningCtx.NextState = GetProductStock
		return GetSensor, nil
	}

	var body ProductStockResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing stock, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 1)
		return GetProductStock, nil
	}

	if body.VariationList == nil {
		taskCtx.SendStatusColored("Stock not loaded, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 10)

		return GetProductStock, nil
	}

	var vars []taskutil.Variant
	for _, variant := range body.VariationList {
		trimmed := taskutil.Variant{
			ID:      variant.Sku,
			Size:    variant.Size,
			InStock: variant.Availability > 0,
		}

		vars = append(vars, trimmed)
	}

	matches, err := taskutil.MatchVariants(staticCtx.Default.Sizes, vars, predicate)
	if err != nil {
		taskCtx.SendStatusColored("No variations matched, retrying...", util.ColorFailed)
		time.Sleep(time.Millisecond * 5000)
		return GetProductStock, nil
	}

	chosenVariant := matches[rand.Intn(len(matches))]
	runningCtx.Size = chosenVariant.Size
	runningCtx.Sku = chosenVariant.ID

	runningCtx.NextState = AddToCart
	return GetSensor, nil
}

func handleAddToCart(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Adding to cart")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return AddToCart, nil
	}

	resp, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/api/checkout/baskets/-/items").
		Method(http.MethodPost).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("checkout-authorization", runningCtx.Authorization).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("origin", "https://www.yeezysupply.com").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId)).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		JsonBody([]AddToCartForm{
			{
				ProductId1:          runningCtx.ProductId,
				ProductVariationSku: runningCtx.Sku,
				ProductId2:          runningCtx.Sku,
				Quantity:            1,
				Size:                runningCtx.Size,
				DisplaySize:         runningCtx.Size,
			},
		}).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error adding to cart, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return AddToCart, nil
	}

	if resp.StatusCode == 400 {
		if runningCtx.AbckRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		taskCtx.SendStatusColored("Error adding to cart [400]", util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.AbckRetries++
		runningCtx.IsFirewallBlocked = false
		runningCtx.NextState = AddToCart

		return GetSensor, nil
	}

	if resp.StatusCode == 403 {
		if runningCtx.FirewallRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		taskCtx.SendStatusColored("Error adding to cart [403]", util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.FirewallRetries++
		runningCtx.IsFirewallBlocked = true
		runningCtx.NextState = AddToCart

		return GetProductPage, nil
	}

	resetRetries(runningCtx)

	var body AddToCartResponse
	if err = resp.JsonBody(&body); err != nil {
		taskCtx.SendStatusColored("Error parsing cart, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 1)
		return AddToCart, nil
	}

	Authorization := extractAuthorizationHeader(resp)
	if Authorization != "" {
		runningCtx.Authorization = Authorization
	}

	runningCtx.BasketId = body.BasketID

	if body.Pricing.BaseTotal > 0 {
		runningCtx.ProductPrice = body.Pricing.BaseTotal
	}

	runningCtx.Referer = "https://www.yeezysupply.com/delivery"
	runningCtx.NextState = SubmitInformation

	return GetSensor, nil
}

func handleSubmitInformation(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Submitting information")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return SubmitInformation, nil
	}

	billingNameParts := strings.Split(runningCtx.ProfileData.Billing.Name, " ")
	shippingNameParts := strings.Split(runningCtx.ProfileData.Shipping.Name, " ")

	var shippingState string
	if runningCtx.ProfileData.Shipping.State != nil {
		shippingState = runningCtx.ProfileData.Shipping.State.Code
	}

	var billingState string
	if runningCtx.ProfileData.Billing.State != nil {
		billingState = runningCtx.ProfileData.Billing.State.Code
	}

	resp, err := runningCtx.Client.
		Builder(fmt.Sprintf("https://www.yeezysupply.com/api/checkout/baskets/%s", runningCtx.BasketId)).
		Method(http.MethodPatch).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("checkout-authorization", runningCtx.Authorization).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", "https://www.yeezysupply.com/delivery").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		JsonBody(map[string]interface{}{
			"customer": map[string]interface{}{
				"email":             runningCtx.ProfileData.Payment.Email,
				"receiveSmsUpdates": false,
			},
			"shippingAddress": map[string]interface{}{
				"country":     "US",
				"firstName":   shippingNameParts[0],
				"lastName":    shippingNameParts[1],
				"address1":    runningCtx.ProfileData.Shipping.Line1,
				"address2":    runningCtx.ProfileData.Shipping.Line2,
				"city":        runningCtx.ProfileData.Shipping.City,
				"stateCode":   shippingState,
				"zipcode":     runningCtx.ProfileData.Shipping.PostCode,
				"phoneNumber": runningCtx.ProfileData.Payment.Phone,
			},
			"billingAddress": map[string]interface{}{
				"country":     "US",
				"firstName":   billingNameParts[0],
				"lastName":    billingNameParts[1],
				"address1":    runningCtx.ProfileData.Billing.Line1,
				"address2":    runningCtx.ProfileData.Billing.Line2,
				"city":        runningCtx.ProfileData.Billing.City,
				"stateCode":   billingState,
				"zipcode":     runningCtx.ProfileData.Billing.PostCode,
				"phoneNumber": runningCtx.ProfileData.Payment.Phone,
			},
			"methodList": []MethodList{
				{
					ID:               "2ndDay-1",
					ShipmentId:       "me",
					CollectionPeriod: "",
					DeliveryPeriod:   "",
				},
			},
			"newsletterSubscription": true,
		}).
		SendAndClose()

	if err != nil {
		taskCtx.SendStatusColored("Error submitting information, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitInformation, nil
	}

	if resp.StatusCode == 400 {
		if runningCtx.AbckRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		status := fmt.Sprintf("Error submitting information [%d]", resp.StatusCode)
		taskCtx.SendStatusColored(status, util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.AbckRetries++
		runningCtx.IsFirewallBlocked = true
		runningCtx.NextState = SubmitInformation

		return GetProductPage, nil
	}

	if resp.StatusCode == 403 {
		if runningCtx.FirewallRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		status := fmt.Sprintf("Error submitting information [%d]", resp.StatusCode)
		taskCtx.SendStatusColored(status, util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.FirewallRetries++
		runningCtx.IsFirewallBlocked = true
		runningCtx.NextState = SubmitInformation

		return GetProductPage, nil
	}

	resetRetries(runningCtx)
	runningCtx.InfoSubmitted = true
	runningCtx.Referer = "https://www.yeezysupply.com/payment"
	runningCtx.NextState = SubmitOrder

	return GetSensor, nil
}

func handleSubmitOrder(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {

	taskCtx.SendStatus("Submitting order")

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return SubmitOrder, nil
	}

	AdyenKey := "C4F415A1A41A283417FAB7EF8580E077284BCC2B06F8A6C1785E31F5ABFD38A3E80760E0CA6437A8DC95BA4720A83203B99175889FA06FC6BABD4BF10EEEF0D73EF86DD336EBE68642AC15913B2FC24337BDEF52D2F5350224BD59F97C1B944BD03F0C3B4CA2E093A18507C349D68BE8BA54B458DB63D01377048F3E53C757F82B163A99A6A89AD0B969C0F745BB82DA7108B1D6FD74303711065B61009BC8011C27D1D1B5B9FC5378368F24DE03B582FE3490604F5803E805AEEA8B9EF86C54F27D9BD3FC4138B9DC30AF43A58CFF7C6ECEF68029C234BBC0816193DF9BD708D10AAFF6B10E38F0721CF422867C8CC5C554A357A8F51BA18153FB8A83CCBED1"
	yen := adyen2.NewAdYen(AdyenKey)
	res, err := yen.Encrypt(adyen2.CardForm{
		Number:          runningCtx.ProfileData.Payment.Number,
		Holder:          runningCtx.ProfileData.Payment.Name,
		Cvv:             runningCtx.ProfileData.Payment.Cvv,
		ExpiryMonth:     runningCtx.ProfileData.Payment.ExpMonth,
		ExpiryYear:      runningCtx.ProfileData.Payment.ExpYear,
		PaymentMethodId: "CREDIT_CARD",
	})

	if err != nil {
		taskCtx.SendStatusColored("Error encrypting card", util.ColorFailed)
		time.Sleep(time.Second * 2)

		return SubmitOrder, nil
	}

	resp, err := runningCtx.Client.
		Builder("https://www.yeezysupply.com/api/checkout/orders").
		Method(http.MethodPost).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("checkout-authorization", runningCtx.Authorization).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", "https://www.yeezysupply.com/delivery").
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		JsonBody(map[string]interface{}{
			"basketId":            runningCtx.BasketId,
			"encryptedInstrument": res,
			"paymentInstrument": map[string]interface{}{
				"holder":          runningCtx.ProfileData.Payment.Name,
				"expirationMonth": runningCtx.ProfileData.Payment.ExpMonth,
				"expirationYear":  runningCtx.ProfileData.Payment.ExpYear,
				"lastFour":        runningCtx.ProfileData.Payment.Number[len(runningCtx.ProfileData.Payment.Number)-4:],
				"paymentMethodId": "CREDIT_CARD",
				"cardType":        strings.ToUpper(runningCtx.ProfileData.Payment.Type),
			},
			"fingerprint": "ryEGX8eZpJ0030000000000000bsx09CX6tD0089146776cVB94iKzBGOhnFsup2S25S16Goh5Mk0045zgp4q8JSa00000qZkTE00000q6IQbnyNfplgQ6OzVQDG:40",
		}).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error submitting order, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return SubmitOrder, nil
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Session expired", util.ColorFailed)
		time.Sleep(time.Second * 2)

		_ = rotateProxy(taskCtx, runningCtx)
		resetContext(staticCtx, runningCtx)

		return GetHomepage, nil
	}

	if resp.StatusCode == 400 && strings.Contains(strings.ToLower(resp.StringBody()), "invalid url") {
		if runningCtx.AbckRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		taskCtx.SendStatusColored("Akamai block, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.AbckRetries++
		runningCtx.IsFirewallBlocked = false
		runningCtx.NextState = SubmitOrder

		return GetProductPage, nil
	}

	if resp.StatusCode == 400 {
		body := ThreeDSResponse{}
		if err = resp.JsonBody(&body); err != nil {
			addWebhook(taskCtx, runningCtx, "failed")

			taskCtx.SendStatusColored("Order failed", util.ColorFailed)
			return "", tasks.CompletedErr
		}

		runningCtx.TermUrl = fmt.Sprintf("https://www.yeezysupply.com/payment/callback/CREDIT_CARD/%s/adyen?orderId=%s&encodedData=%s&result=AUTHORISED", runningCtx.BasketId, body.OrderID, body.PaRedirectForm.FormFields.EncodedData)
		runningCtx.Form = ThreeDSForm{
			FormMethod: body.PaRedirectForm.FormMethod,
			FormAction: body.PaRedirectForm.FormAction,
			FormFields: body.PaRedirectForm.FormFields,
			OrderId:    body.OrderID,
		}

		return LaunchBrowser, nil
	}

	if resp.StatusCode == 403 {
		if runningCtx.FirewallRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		status := fmt.Sprintf("Error submitting order [%d]", resp.StatusCode)
		taskCtx.SendStatusColored(status, util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.FirewallRetries++
		runningCtx.IsFirewallBlocked = true
		runningCtx.NextState = SubmitOrder

		return GetProductPage, nil
	}

	addWebhook(taskCtx, runningCtx, "placed")

	taskCtx.SendStatusColored("Order placed", util.ColorSuccess)
	return "", tasks.CompletedErr
}

func handleLaunchBrowser(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatusColored("Solving 3DSecure", util.ColorWarning)

	if runningCtx.ThreeDsCallback == nil {

		proxyUrl := ""
		if runningCtx.ProxyData != nil {
			proxyUrl = runningCtx.ProxyData.Url.String()
		}

		var err error
		runningCtx.ThreeDsCallback, err = solve3ds(
			taskCtx.Context,
			runningCtx.Client.Cookies(runningCtx.StoreUrl),
			runningCtx.StoreUrl.String(),
			proxyUrl,
			runningCtx.UserAgent,
			runningCtx.Form,
			runningCtx.TermUrl,
		)
		if err != nil {
			return "", err
		}
	}

	return WaitFor3DS, nil
}

func handleWaitFor3DS(taskCtx *tasks.Context, _ *staticCtx, runningCtx *runningCtx) (State, error) {
	// should never happen
	if runningCtx.ThreeDsCallback == nil {
		return "", errors.New("3ds callback nil")
	}

	select {
	case <-taskCtx.Context.Done():
		return "", nil
	case threeDsData := <-runningCtx.ThreeDsCallback:
		runningCtx.ThreeDsCallback = nil
		runningCtx.ThreeDsData = &threeDsData
	}

	return CompleteOrder, nil
}

func handleCompleteOrder(taskCtx *tasks.Context, staticCtx *staticCtx, runningCtx *runningCtx) (State, error) {
	taskCtx.SendStatus("Completing order")

	var reqBody = make(map[string]interface{})
	reqBody["orderId"] = runningCtx.ThreeDsData.OrderID
	for k, v := range runningCtx.ThreeDsData.Data {
		reqBody[k] = v
	}

	correlation, err := buildCorrelationHeaders()
	if err != nil {
		return CompleteOrder, nil
	}

	resp, err := runningCtx.Client.
		Builder(runningCtx.ThreeDsData.PaymentUrl).
		Method(http.MethodPost).
		Header("x-instana-t", correlation["x-instana-t"]).
		Header("dnt", "1").
		Header("sec-ch-ua-mobile", "?0").
		Header("user-agent", runningCtx.UserAgent).
		Header("x-instana-l", correlation["x-instana-l"]).
		Header("x-instana-s", correlation["x-instana-s"]).
		Header("content-type", "application/json").
		Header("checkout-authorization", runningCtx.Authorization).
		Header("sec-ch-ua", buildSecurityHeader(runningCtx.UserAgent)).
		Header("accept", "*/*").
		Header("origin", runningCtx.StoreUrl.String()).
		Header("sec-fetch-site", "same-origin").
		Header("sec-fetch-mode", "cors").
		Header("sec-fetch-dest", "empty").
		Header("referer", runningCtx.ThreeDsData.TermUrl).
		Header("accept-encoding", "gzip, deflate, br").
		Header("accept-language", "en-US,en;q=0.9").
		JsonBody(reqBody).
		SendAndClose()
	if err != nil {
		taskCtx.SendStatusColored("Error completing order, retrying...", util.ColorFailed)
		time.Sleep(time.Second * 2)
		return CompleteOrder, nil
	}

	if resp.StatusCode == 404 {
		taskCtx.SendStatusColored("Session expired, resetting...", util.ColorFailed)
		time.Sleep(time.Second * 2)

		_ = rotateProxy(taskCtx, runningCtx)
		resetContext(staticCtx, runningCtx)

		return GetHomepage, nil
	}

	if resp.StatusCode == 400 {
		addWebhook(taskCtx, runningCtx, "failed")

		taskCtx.SendStatusColored("Order failed", util.ColorFailed)
		return "", tasks.CompletedErr
	}

	if resp.StatusCode == 403 {
		if runningCtx.FirewallRetries > runningCtx.MaxRetries {
			taskCtx.SendStatusColored("Max attempts, resetting...", util.ColorFailed)
			time.Sleep(time.Second * 2)

			_ = rotateProxy(taskCtx, runningCtx)
			resetContext(staticCtx, runningCtx)

			return GetHomepage, nil
		}

		taskCtx.SendStatusColored("Error completing order [403]", util.ColorFailed)
		time.Sleep(time.Second * 2)

		runningCtx.FirewallRetries++
		runningCtx.IsFirewallBlocked = true
		runningCtx.NextState = CompleteOrder

		return GetProductPage, nil
	}

	addWebhook(taskCtx, runningCtx, "placed")

	taskCtx.SendStatusColored("Order placed", util.ColorSuccess)
	return "", tasks.CompletedErr
}
