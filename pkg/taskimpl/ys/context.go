package ys

import (
	"context"
	"math/rand"
	"nebula/cmd/bot/inject"
	"nebula/pkg/api/model"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/third_party/http"
	"net/url"
	"time"

	"github.com/cloudflare/ahocorasick"

	"golang.org/x/net/proxy"
)

var (
	_ = tasks.StaticContext((*staticCtx)(nil))

	ysUrl, _ = url.Parse("https://www.yeezysupply.com/")
)

var UserAgents = [...]string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
}

type staticCtx struct {
	Default tasks.DefaultData `json:"default"`
}

func (ctx *staticCtx) Metadata() (meta tasks.Metadata) {
	meta = tasks.Metadata{
		DisplayName: "Yeezy Supply",
		Stores: []tasks.Store{
			{
				Name: "Yeezy Supply",
				Url:  "https://www.yeezysupply.com",
			},
		},
		UsesAccounts: false,
	}

	for i := range meta.Stores {
		meta.Stores[i].Modes = []tasks.Mode{
			{
				Label: "Normal",
				Name:  "normal",
			},
		}
	}

	return
}

func (ctx *staticCtx) Validate() []error {
	return nil
}

func (ctx *staticCtx) GetProductInfo() *model.CheckoutProduct {
	return &model.CheckoutProduct{Name: &ctx.Default.Sku}
}

type runningCtx struct {
	Client           *http.SimpleClient
	StoreUrl         *url.URL
	CaptchaShareChan chan harvesters.SolverResponse
	CaptchaToken     *harvesters.SolverResponse

	Authorization string
	BasketId      string
	ProductName   string
	ProductImage  string
	ProductPrice  int64

	Size string
	Sku  string

	UserAgent       string
	SecUAHeader     string
	ProductId       string
	PixelScript     string
	PixelValue      string
	PixelHash       string
	PixelPayload    string
	SensorPayload   string
	SensorUserAgent string

	FirewallRetries int
	AbckRetries     int
	MaxRetries      int

	RecaptchaCookie string
	CaptchaExpires  time.Time

	NumberPosts       int
	ReturnState       State
	NextState         State
	InfoSubmitted     bool
	IsFirewallBlocked bool
	IsFirstBasket     bool
	NeedsAuth         bool
	AkamaiUrl         string
	WaitingRoom       string
	Referer           string
	Abck              string
	ScriptName        string
	TermUrl           string
	Form              ThreeDSForm
	ThreeDsData       *ThreeDSData
	SplashFlag        string
	SplashUrl         string

	ThreeDsCallback chan ThreeDSData

	BannedMatcher *ahocorasick.Matcher

	ConfigCancelFunc context.CancelFunc

	ProfileData *profiles.ProfileData
	ProxyData   *proxies.ProxyData
}

func getUserAgent() string {
	return UserAgents[rand.Intn(len(UserAgents))]
}

// TODO: possibly move this to a `taskutil` package?
func (ctx *runningCtx) getProxyDialer(taskCtx *tasks.Context) (proxy.ContextDialer, error) {
	proxyData, err := proxies.GetProxyForTask(taskCtx.Identifier())
	if err == nil {
		ctx.ProxyData = proxyData
		return http.NewConnectDialer(proxyData.Url)
	} else if err == proxies.NotLinkedErr {
		return proxy.Direct, nil
	} else {
		return nil, err
	}
}

func (ctx *runningCtx) Init(taskCtx *tasks.Context) error {
	ctx.StoreUrl = ysUrl
	ctx.UserAgent = getUserAgent()
	ctx.Authorization = "null"
	ctx.RecaptchaCookie = "xhwUqgFqfW88H50"
	ctx.AkamaiUrl = "https://www.yeezysupply.com/c_8G/W4wQ/O8J/QpX/oUXQ/EDYtrm2r9X/aRc-Ag/TSM/GQQFqc1Q"
	ctx.WaitingRoom = "https://www.yeezysupply.com/wrgen_orig_assets/ddad474a32813e165e22.js"

	ctx.SplashFlag = "wrgen_orig_assets"
	ctx.FirewallRetries = 0
	ctx.AbckRetries = 0
	ctx.MaxRetries = 15
	ctx.IsFirstBasket = true

	ctx.BannedMatcher = ahocorasick.NewStringMatcher([]string{"UNFORTUNATELY WE ARE UNABLE TO GIVE YOU ACCESS TO OUR SITE AT THIS TIME", "If you're on a personal connection you can run"})

	clientConfig := http.DefaultClientConfig.Clone()
	proxyDialer, err := ctx.getProxyDialer(taskCtx)
	if err != nil {
		return err
	}

	if inject.IsDev == "yes" {
		if proxyDialer == proxy.Direct {
			proxyDialer = taskutil.NewCharlesDialer()
		}
	}

	clientConfig.ProxyDialer = proxyDialer

	profileData, err := profiles.GetProfileForTask(taskCtx.Identifier())
	if err != nil {
		return err
	}
	ctx.ProfileData = profileData

	client := http.NewSimpleClient(taskCtx.Context, clientConfig, http.FingerprintChrome)
	client.Inner().Jar = newSingleJar()
	client.Inner().Timeout = 15 * time.Second

	ctx.Client = client
	return nil
}
