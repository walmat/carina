package ys

import (
	"context"
	"fmt"
	"math"
	"nebula/cmd/bot/app"
	"nebula/cmd/bot/inject"
	"nebula/third_party/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/Nebulabots/go-astilectron"
	"github.com/asticode/go-astikit"
	"golang.org/x/sync/semaphore"
)

var (
	browserSem = semaphore.NewWeighted(int64(math.Max(float64(runtime.NumCPU()/2), 5)))
)

const ThreeDSFormFmt = `
		<html>
			<body>
				<form method="%s" action="%s" id="Cardinal-CCA-Form">
					<input type="hidden" name="PaReq" value="%s" />
					<input type="hidden" name="MD" value="%s" />
					<input type="hidden" name="TermUrl" value="%s" />
				</form>
				<script>
					document.getElementById("Cardinal-CCA-Form").submit();
				</script>
			</body>
		</html>
	`

func solve3ds(ctx context.Context, cookies []*http.Cookie, storeUrl, proxy, userAgent string, form ThreeDSForm, termUrl string) (chan ThreeDSData, error) {
	if err := browserSem.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	dataUrl := fmt.Sprintf(ThreeDSFormFmt,
		form.FormMethod,
		form.FormAction,
		form.FormFields.PaReq,
		form.FormFields.MD,
		termUrl,
	)

	fmt.Println(dataUrl)

	var windowProxyOptions *astilectron.WindowProxyOptions = nil
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		windowProxyOptions = &astilectron.WindowProxyOptions{
			Rules: fmt.Sprintf("http=%s;https=%s", proxyUrl.Host, proxyUrl.Host),
		}
	}

	w, err := app.App.NewWindow(storeUrl, &astilectron.WindowOptions{
		Center:           astikit.BoolPtr(true),
		Title:            astikit.StrPtr("Yeezy Supply - 3DSecure"),
		Height:           astikit.IntPtr(600),
		Width:            astikit.IntPtr(450),
		Resizable:        astikit.BoolPtr(false),
		Frame:            astikit.BoolPtr(true),
		Show:             astikit.BoolPtr(false),
		Transparent:      astikit.BoolPtr(false),
		AcceptFirstMouse: astikit.BoolPtr(true),
		Proxy:            windowProxyOptions,
		WebPreferences: &astilectron.WebPreferences{
			WebviewTag:                  astikit.BoolPtr(true),
			AllowRunningInsecureContent: astikit.BoolPtr(true),
			BackgroundThrottling:        astikit.BoolPtr(true),
			ContextIsolation:            astikit.BoolPtr(false),
			WebSecurity:                 astikit.BoolPtr(false),
		},
		Load: &astilectron.WindowLoadOptions{
			UserAgent: userAgent,
		},
	})
	if err != nil {
		return nil, err
	}

	w.OnLogin(func(i astilectron.Event) (username string, password string, err error) {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return "", "", err
		}

		if proxyUrl.User == nil {
			return "", "", nil
		}

		proxyPass, _ := proxyUrl.User.Password()
		return proxyUrl.User.Username(), proxyPass, nil
	})

	w.On(astilectron.EventNameWindowEventReadyToShow, func(e astilectron.Event) (deleteListener bool) {
		_ = w.Show()

		if inject.IsDev == "yes" {
			_ = w.OpenDevTools()
		}

		u := fmt.Sprintf("data:text/html,%s", strings.ReplaceAll(url.QueryEscape(dataUrl), "+", "%20"))

		w.ExecuteJavaScript(fmt.Sprintf(`document.write("<html><iframe src='%s'></iframe></html>");`, u))

		return true
	})

	w.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
		if inject.IsDev == "yes" {
			_ = w.CloseDevTools()
		}

		browserSem.Release(1)
		return true
	})

	if err = w.Create(); err != nil {
		return nil, err
	}

	if len(cookies) > 0 {
		var cooks []astilectron.SessionCookie
		for _, c := range cookies {
			cooks = append(cooks, astilectron.SessionCookie{
				Url:            storeUrl,
				Name:           c.Name,
				Value:          c.Value,
				Domain:         c.Domain,
				Path:           c.Path,
				Secure:         &c.Secure,
				HttpOnly:       &c.HttpOnly,
				Session:        nil,
				ExpirationDate: nil,
				SameSite:       "",
			})
		}

		if err = w.Session.SetCookies(cooks); err != nil {
			return nil, err
		}
	}

	callbackChan := make(chan ThreeDSData)
	_ = w.Session.OnBeforeRequest(astilectron.FilterOptions{}, func(event astilectron.Event) (bool, string, bool) {
		if event.Request.URL != termUrl {
			return false, "", false
		}

		if data, ok := event.Request.UploadData[0]; ok {
			dataLen := len(data.Bytes)
			dataArr := make([]byte, dataLen)

			for i := 0; i < dataLen; i++ {
				dataArr[i] = data.Bytes[strconv.Itoa(i)]
			}

			threeData := make(map[string]interface{})
			d, _ := url.ParseQuery(string(dataArr))
			for k, v := range d {
				threeData[k] = strings.Join(v, " ")
			}

			callbackChan <- ThreeDSData{
				Data:       threeData,
				OrderID:    form.OrderId,
				PaymentUrl: fmt.Sprintf("%sapi/checkout/payment-verification/%s", storeUrl, form.FormFields.EncodedData),
				TermUrl:    termUrl,
			}

			_ = w.Destroy()
		}

		// TODO: figure out why window still sends onBeforeRequest (for possibly other windows) but
		// still has this target id and requires a response or app locks up?????????????????????
		return true, "", false
	})

	return callbackChan, nil
}
