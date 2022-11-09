package harvesters

import (
	"errors"
	"fmt"
	"nebula/cmd/bot/windows"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nebulabots/go-astilectron"
	"github.com/asticode/go-astikit"
	"github.com/lithammer/shortuuid/v3"
)

const (
	Gmail_Login            = "GMAIL_LOGIN"
	HCaptcha               = "HCAPTCHA"
	ReCaptcha_V2           = "RECAPTCHA_V2"
	ReCaptcha_V2_Invisible = "RECAPTCHA_V2_INVISIBLE"
	ReCaptcha_V3           = "RECAPTCHA_V3"
	Shopify_Challenge      = "SHOPIFY_CHALLENGE"
	Shopify_Checkpoint     = "SHOPIFY_CHECKPOINT"
)

var (
	BrowserViews           = make(map[string]*astilectron.BrowserView)
	BrowserViewsMutex      = sync.RWMutex{}
	ErrBrowserViewNotFound = errors.New("browserview does not exist")
	ErrInvalidProxyFormat  = errors.New("invalid proxy format")
)

func init() {
	go func() {
		for {
			if windows.Windows.CollectiveWindow != nil {
				for _, h := range Harvesters {
					Create(h.ID)
				}
				return
			}
		}
	}()
}

// Create
func Create(HarvesterID string) (b *astilectron.BrowserView, err error) {
	BrowserViewsMutex.Lock()
	defer BrowserViewsMutex.Unlock()

	s := windows.App.NewSession()

	partition := fmt.Sprintf("persist:%s", shortuuid.New())

	s.FromPartition(partition)

	s.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36", "")

	b, err = windows.App.NewBrowserView("", &astilectron.WindowOptions{
		WebPreferences: &astilectron.WebPreferences{
			Session: &s.ID,
		},
	}, s)

	if err != nil {
		return
	}

	err = b.Create()

	if err != nil {
		return
	}

	err = windows.Windows.CollectiveWindow.AddBrowserView(b)

	if err != nil {
		return
	}

	BrowserViews[HarvesterID] = b

	return
}

// Remove
func Remove(HarvesterID string) {
	if BrowserViews[HarvesterID] == nil {
		return
	}

	BrowserViewsMutex.Lock()
	defer BrowserViewsMutex.Unlock()

	windows.Windows.CollectiveWindow.RemoveBrowserView(BrowserViews[HarvesterID])

	BrowserViews[HarvesterID] = nil
}

// Show
func Show(HarvesterID string) {
	if BrowserViews[HarvesterID] == nil {
		return
	}

	BrowserViews[HarvesterID].SetBounds(&astilectron.RectangleOptions{
		PositionOptions: astilectron.PositionOptions{
			X: astikit.IntPtr(0),
			Y: astikit.IntPtr(110),
		},
		SizeOptions: astilectron.SizeOptions{
			Width:  astikit.IntPtr(400),
			Height: astikit.IntPtr(400),
		},
	})
}

// Hide
func Hide(HarvesterID string) {
	if BrowserViews[HarvesterID] == nil {
		return
	}

	BrowserViews[HarvesterID].SetBounds(&astilectron.RectangleOptions{
		PositionOptions: astilectron.PositionOptions{
			X: astikit.IntPtr(0),
			Y: astikit.IntPtr(0),
		},
		SizeOptions: astilectron.SizeOptions{
			Width:  astikit.IntPtr(0),
			Height: astikit.IntPtr(0),
		},
	})
}

func Load(HarvesterID string) (err error) {
	h := Harvesters[HarvesterID]
	b := BrowserViews[HarvesterID]

	if b == nil {
		return
	}

	if Harvesters[HarvesterID].Focused {
		Show(HarvesterID)
	}

	if h.Solver.Proxy != nil {
		SetProxy(HarvesterID, *h.Solver.Proxy)
	}

	if h.Account != nil && h.Account.Cookies != nil {

		var t []astilectron.SessionCookie
		for _, c := range *h.Account.Cookies {
			if c.Domain == ".youtube.com" {
				c.Url = "https://www.youtube.com/"
				t = append(t, c)
			}

			if c.Domain == ".google.com" {
				c.Url = "https://www.google.com/"
				t = append(t, c)
			}

			if c.Domain == "accounts.google.com" {
				c.Url = "https://accounts.google.com/"
				t = append(t, c)
			}
		}
		b.Session.SetCookies(t)
	}

	// b.OpenDevTools()

	switch h.Solver.Type {
	case Gmail_Login:
		var email string
		b.Session.OnBeforeRequest(astilectron.FilterOptions{Urls: []string{"*://*/*"}}, func(i astilectron.Event) (bool, string, bool) {
			if strings.Contains(i.Request.URL, "https://accounts.google.com/_/lookup/accountlookup") {
				r, _ := b.ExecuteJavaScript(`document.querySelector("[type=email").value;`)
				email = r.CodeResult
			}

			if i.Request.URL == "https://www.youtube.com/" {
				c, _ := b.Session.GetCookies()

				ReceiveSolve(SolverResponse{
					ID:        h.Solver.ID,
					Success:   true,
					CreatedAt: time.Now().Unix(),
					Data: ResponseData{
						Cookies: c.Cookies,
						Email:   email,
					},
				})
			}
			return false, "", false
		})

		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{UserAgent: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:70.0) Gecko/20100101 Firefox/70.0"})
		return
	case ReCaptcha_V2:
		template := fmt.Sprintf(`
		<html>
			<head>
				<script src="https://www.google.com/recaptcha/api.js?render=explicit&hl=en" async defer></script>
				<script>
				let onCaptchaSuccess = (token) => {
					let solveEvent = new CustomEvent("solve", {detail: token});
					window.dispatchEvent(solveEvent);
				};

				let onloadCallback = async () => {
				  if (
					typeof window.grecaptcha === 'undefined' ||
					typeof window.grecaptcha.render === 'undefined'
				  ) {
					await new Promise(resolve => setTimeout(resolve, 50));

					return onloadCallback();
				  }

				  const parameters = {
					sitekey: '%s',
					size: 'normal',
					callback: onCaptchaSuccess,
				  };

				  window.grecaptcha.render('g-recaptcha', parameters);
				};

				onloadCallback();
				</script>
			</head>
			<body>
				<div id="captcha">
					<div class="g-recaptcha" id="g-recaptcha"></div>
				</div>
			</body>
		</html>`, h.Solver.ReCaptcha.SiteKey)

		b.InterceptStringProtocol("https", func(i astilectron.Event) (mimeType string, data string, deleteListener bool) {
			if i.Request.URL == h.Solver.Store.Url {
				b.UninterceptProtocol("https")
				return "text/html", template, true
			}

			return
		})

		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{})

		go RetrieveToken(HarvesterID)

		return
	case ReCaptcha_V2_Invisible:
		template := fmt.Sprintf(`
		<html>
			<head>
				<script src="https://www.google.com/recaptcha/api.js" async defer></script>
				<script>
				let onCaptchaSuccess = (token) => {
					let solveEvent = new CustomEvent("solve", {detail: token});
					window.dispatchEvent(solveEvent);
				};

				let onloadCallback = async () => {
				  if (
					typeof window.grecaptcha === 'undefined' ||
					typeof window.grecaptcha.render === 'undefined'
				  ) {
					await new Promise(resolve => setTimeout(resolve, 50));

					return onloadCallback();
				  }
		
				  const parameters = {
					sitekey: %s,
					size: 'invisible',
					callback: onCaptchaSuccess,
				  };
		
				  window.grecaptcha.render('g-recaptcha', parameters);

				  window.grecaptcha.execute();
				};

				onloadCallback();
				</script>
			</head>
			<body>
				<div id="captcha">
					<div class="g-recaptcha" id="g-recaptcha"></div>
				</div>
			</body>
		</html>`, h.Solver.ReCaptcha.SiteKey)

		b.InterceptStringProtocol("https", func(i astilectron.Event) (mimeType string, data string, deleteListener bool) {
			if i.Request.URL == h.Solver.Store.Url {
				b.UninterceptProtocol("https")
				return "text/html", template, true
			}

			return
		})

		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{})

		go RetrieveToken(HarvesterID)
		return
	case ReCaptcha_V3:
		template := fmt.Sprintf(`
		<html>
			<head>
				<script src="https://www.google.com/recaptcha/api.js?render=%s"></script>
				<script>
					grecaptcha.ready(function() {
						grecaptcha.execute('%s', {action: '%s'}).then(function(token) {
							let str = "token " + token;
							// alert(str);
							let solveEvent = new CustomEvent("solve", {detail: token});
							window.dispatchEvent(solveEvent);
						});
					});
				</script>
			</head>
			<body>
				<div id="captcha">
					<div class="g-recaptcha" id="g-recaptcha"></div>
				</div>
			</body>
		</html>`, h.Solver.ReCaptcha.SiteKey, h.Solver.ReCaptcha.SiteKey, h.Solver.ReCaptcha.Action)

		b.InterceptStringProtocol("https", func(i astilectron.Event) (mimeType string, data string, deleteListener bool) {
			if i.Request.URL == h.Solver.Store.Url {
				b.UninterceptProtocol("https")
				return "text/html", template, true
			}

			return
		})
		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{})

		go RetrieveToken(HarvesterID)

		return
	case Shopify_Checkpoint:
		// todo: is the POST including the query or without?
		u, _ := url.Parse(h.Solver.Store.Url)
		p := fmt.Sprintf("*://%s/*", u.Host)

		urls := []string{p}
		b.Session.OnBeforeRequest(astilectron.FilterOptions{Urls: urls}, func(i astilectron.Event) (cancel bool, redirectUrl string, deleteListener bool) {
			if i.Request.Method == "POST" && strings.Contains(i.Request.URL, "checkpoint") {
				if data, ok := i.Request.UploadData[0]; ok {
					dataLen := len(data.Bytes)
					dataArr := make([]byte, dataLen)

					for i := 0; i < dataLen; i++ {
						dataArr[i] = data.Bytes[strconv.Itoa(i)]
					}

					formValues, err := url.ParseQuery(string(dataArr))

					if err != nil {
						return
					}

					ev, err := b.Session.GetCookies()

					if err != nil {
						return
					}

					deleteListener = true
					cancel = true

					go ReceiveSolve(SolverResponse{
						ID:        Harvesters[HarvesterID].Solver.ID,
						Success:   true,
						CreatedAt: time.Now().Unix(),
						Data: ResponseData{
							Form:    formValues,
							Cookies: ev.Cookies,
						},
					})
				}
			}
			return
		})
		// listen for the request to /checkpoint and intercept it.
		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{})

		// return the form to be submitted.
		return
	default:
		err = b.LoadURL(h.Solver.Store.Url, &astilectron.Load{})
		return
	}
}

func RetrieveToken(harvesterID string) (err error) {
	h := Harvesters[harvesterID]
	b := BrowserViews[harvesterID]

	if b == nil {
		return ErrBrowserViewNotFound
	}

	r, err := b.ExecuteJavaScript(`
	new Promise((resolve) => {
		window.addEventListener("solve", (e) => {
			resolve(e.detail);
		});
	});
	`)

	if err != nil {
		return
	}

	ReceiveSolve(SolverResponse{
		ID:        h.Solver.ID,
		Success:   true,
		CreatedAt: time.Now().Unix(),
		Data: ResponseData{
			Token: r.CodeResult,
		},
	})

	return
}

func SetProxy(harvesterID string, proxy string) (err error) {
	b := BrowserViews[harvesterID]

	if b == nil {
		return ErrBrowserViewNotFound
	}

	if proxy == "" {
		b.Session.CloseAllConnections()

		return b.Session.SetProxy(astilectron.WindowProxyOptions{
			Rules: "",
		})
	}

	r := regexp.MustCompile("http://|https://")
	s := r.ReplaceAllString(proxy, "")
	sp := strings.Split(s, ":")

	formatted := ""
	if len(sp) == 2 {
		formatted = fmt.Sprintf(`//%s:%s`, sp[0], sp[1])
	} else if len(sp) == 4 {
		formatted = fmt.Sprintf(`//%s:%s@%s:%s`, sp[2], sp[3], sp[0], sp[1])
	} else {
		return ErrInvalidProxyFormat
	}

	parsed, err := url.Parse(formatted)

	if err != nil {
		return ErrInvalidProxyFormat
	}

	if parsed.User != nil {
		b.OnLogin(func(i astilectron.Event) (username string, password string, err error) {
			proxyPass, _ := parsed.User.Password()
			return parsed.User.Username(), proxyPass, nil
		})
	}

	b.Session.CloseAllConnections()

	err = b.Session.SetProxy(astilectron.WindowProxyOptions{
		Rules: fmt.Sprintf("http=%s;https=%s", parsed.Host, parsed.Host),
	})

	return
}
