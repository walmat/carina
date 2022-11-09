package foots

import (
	"nebula/pkg/api/model"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/third_party/http"
	"net/url"

	"golang.org/x/net/proxy"
)

var (
	_ = tasks.StaticContext((*staticCtx)(nil))
)

type staticCtx struct {
	Default tasks.DefaultData `json:"default"`
}

func (ctx *staticCtx) Metadata() (meta tasks.Metadata) {
	meta = tasks.Metadata{
		DisplayName: "Footsites",
		Stores: []tasks.Store{
			{
				Name: "Footlocker US",
				Url:  "https://www.footlocker.com",
			},
			{
				Name: "Footaction",
				Url:  "https://www.footaction.com",
			},
			{
				Name: "Eastbay",
				Url:  "https://www.eastbay.com",
			},
			{
				Name: "Champssports",
				Url:  "https://www.champssports.com",
			},
			{
				Name: "Kids Footlocker",
				Url:  "https://www.kidsfootlocker.com",
			},
		},
		UsesAccounts: true,
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
	return nil
}

type runningCtx struct {
	Client   *http.SimpleClient
	StoreUrl *url.URL

	ProductId string
	SessionId string
	CartId    string
	CsrfToken string
	OrderId   string
	OrderUrl  string

	Cookies map[string]string

	ProductName  string
	ProductPrice int64
	ProductUrl   string
	ProductImage string
	Sku          string
	Size         string

	FootsiteType

	ProfileData *profiles.ProfileData
	ProxyData   *proxies.ProxyData
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
	profile, err := profiles.GetProfileForTask(taskCtx.Identifier())
	if err != nil {
		return err
	}
	ctx.ProfileData = profile

	return nil
}
