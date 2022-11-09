package shopify

import (
	"math/rand"
	"nebula/pkg/api/model"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/taskutil"
	"nebula/third_party/http"
	tls "nebula/third_party/utls"
	"net/url"
	"time"

	"github.com/elliotchance/orderedmap"

	"golang.org/x/net/proxy"
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

type runningCtx struct {
	NextState State
	PrevState State

	Client    *http.SimpleClient
	StoreUrl  *url.URL
	UserAgent string

	ProxyData   *proxies.ProxyData
	ProfileData *profiles.ProfileData
	Mode        string

	// AccessToken Token used to make requests to the shopify backend api
	AccessToken string
	ShopID      int
	// CheckoutHash An unique checkout session
	CheckoutHash string
	ShippingRate string
	// AuthToken Token used to submit forms in the shopify checkout stages
	AuthToken        string
	PaymentSessionId string
	PaymentGateway   string
	Sitekey          string
	Product          struct {
		Variant taskutil.Variant
		Size string
		Name    string
		Url     string
		Image   string
		Price   string
	}
	PreloadProduct struct {
		Variant    VariantOptionResponse
		ProductUrl string
	}
	Form *orderedmap.OrderedMap

	RetryCount   int64
	Fallback     bool
	IsFree       bool
	CaptchaToken string
	Restocking   bool
	Checked      bool
	// PollingUrl Used whenever an URL is being requested multiple times
	PollingUrl bool

	// note: remove when these values appear on the task's default data
	UsesPassword       bool
	StorePassword      string
	UsesAccount        bool
	AccountCredentials string

	Redirect       string
	CheckpointForm url.Values

	Preloading bool

	NextQueueToken string
	UseNewQueue    bool
	Ctd            string

	QueueETA          int
	QueueAvailability string
}

func getUserAgent() string {
	return UserAgents[rand.Intn(len(UserAgents))]
}

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
	ctx.UserAgent = getUserAgent()

	clientConfig := http.DefaultClientConfig.Clone()

	proxyDialer, err := ctx.getProxyDialer(taskCtx)
	if err != nil {
		return err
	}

	//if inject.IsDev == "yes" {
	//	if proxyDialer == proxy.Direct {
	//		proxyDialer = taskutil.NewCharlesDialer()
	//	}
	//}

	clientConfig.ClientHello = tls.HelloGolang
	clientConfig.ProxyDialer = proxyDialer

	profileData, err := profiles.GetProfileForTask(taskCtx.Identifier())
	if err != nil {
		return err
	}
	ctx.ProfileData = profileData

	client := http.NewSimpleClient(taskCtx.Context, clientConfig, http.DefaultFingerprint)

	client.Inner().Jar = newSingleJar()
	client.Inner().Timeout = 15 * time.Second

	ctx.Client = client
	return nil
}

// Metadata Contains all the shopify stores data
func (ctx *staticCtx) Metadata() (meta tasks.Metadata) {
	meta = tasks.Metadata{
		DisplayName: "Shopify",
		Stores: []tasks.Store{
			{
				Name: "A Ma Maniere",
				Url:  "https://www.a-ma-maniere.com",
			},
			{
				Name: "Addict Miami",
				Url:  "https://www.addictmiami.com",
			},
			{
				Name: "Atmos USA",
				Url:  "https://www.atmosusa.com",
			},
			{
				Name: "Aime Leon Dore",
				Url:  "https://www.aimeleondore.com",
			},
			{
				Name: "Alife NY",
				Url:  "https://alifenewyork.com",
			},
			{
				Name: "APB Store",
				Url:  "https://www.apbstore.com",
			},
			{
				Name: "Bape US",
				Url:  "https://us.bape.com",
			},
			{
				Name: "Bape JP",
				Url:  "https://bape.com",
			},
			{
				Name: "BBBranded",
				Url:  "https://www.bbbranded.com",
			},
			{
				Name: "BBC Ice Cream",
				Url:  "https://www.bbcicecream.com",
			},
			{
				Name: "BlendsUS",
				Url:  "https://www.blendsus.com",
			},
			{
				Name: "Bodega",
				Url:  "https://bdgastore.com",
			},
			{
				Name: "Cactus Plant Flea Market",
				Url:  "https://cactusplantfleamarket.com",
			},
			{
				Name: "Chicago City Sports",
				Url:  "https://chicagocitysports.com",
			},
			{
				Name: "Concepts",
				Url:  "https://cncpts.com",
			},
			{
				Name: "Creme321",
				Url:  "https://creme321.com",
			},
			{
				Name: "Deadstock CA",
				Url:  "https://www.deadstock.ca",
			},
			{
				Name: "Drew House",
				Url:  "https://thehouseofdrew.com",
			},
			{
				Name: "DSMJP E-Flash",
				Url:  "https://eflash-jp.doverstreetmarket.com",
			},
			{
				Name: "DSMJP E-Shop",
				Url:  "https://shop-jp.doverstreetmarket.com",
			},
			{
				Name: "DSML E-Flash",
				Url:  "https://eflash.doverstreetmarket.com",
			},
			{
				Name: "DSML E-Shop",
				Url:  "https://shop.doverstreetmarket.com",
			},
			{
				Name: "DSMNY E-Flash",
				Url:  "https://eflash-us.doverstreetmarket.com",
			},
			{
				Name: "DSMNY E-Shop",
				Url:  "https://shop-us.doverstreetmarket.com",
			},
			{
				Name: "DSMSG E-Flash",
				Url:  "https://eflash-sg.doverstreetmarket.com",
			},
			{
				Name: "DSMSG E-Shop",
				Url:  "https://shop-sg.doverstreetmarket.com",
			},
			{
				Name: "DTLR",
				Url:  "https://www.dtlr.com",
			},
			{
				Name: "Eric Emanuel",
				Url:  "https://www.ericemanuel.com",
			},
			{
				Name: "Exclusivity",
				Url:  "https://shop.exclucitylife.com",
			},
			{
				Name: "Exclusive Fitted",
				Url:  "https://exclusivefitted.com",
			},
			{
				Name: "Fear of God",
				Url:  "https://fearofgod.com",
			},
			{
				Name: "Feature",
				Url:  "https://feature.com",
			},
			{
				Name: "Final Mouse",
				Url:  "https://finalmouse.com",
			},
			{
				Name: "Funko Shop",
				Url:  "https://www.funko.com",
			},
			{
				Name: "Hanon Shop",
				Url:  "https://www.hanon-shop.com",
			},
			{
				Name: "Hat Club",
				Url:  "https://www.hatclub.com",
			},
			{
				Name: "Haven Shop",
				Url:  "https://havenshop.com",
			},
			{
				Name: "Human Made JP",
				Url:  "https://humanmade.jp",
			},
			{
				Name: "Jimmy Jazz",
				Url:  "https://www.jimmyjazz.com",
			},
			{
				Name: "JJJound",
				Url:  "https://www.jjjjound.com",
			},
			{
				Name: "Juice Store",
				Url:  "https://juicestore.com",
			},
			{
				Name: "Kaws",
				Url:  "https://kawsone.com",
			},
			{
				Name: "Kith",
				Url:  "https://kith.com",
			},
			{
				Name: "Lust Mexico",
				Url:  "https://www.lustmexico.com",
			},
			{
				Name: "MCT Tokyo",
				Url:  "https://mct.tokyo",
			},
			{
				Name: "NRML CA",
				Url:  "https://nrml.ca",
			},
			{
				Name: "Oneness Boutique",
				Url:  "https://www.onenessboutique.com",
			},
			{
				Name: "Packer Shoes",
				Url:  "https://packershoes.com",
			},
			{
				Name: "Palace US",
				Url:  "https://shop-usa.palaceskateboards.com",
			},
			{
				Name: "RSVP Gallery",
				Url:  "https://rsvpgallery.com",
			},
			{
				Name: "Saint Alfred",
				Url:  "https://www.saintalfred.com",
			},
			{
				Name: "Shoe Palace",
				Url:  "https://www.shoepalace.com",
			},
			{
				Name: "Shop Nice Kicks",
				Url:  "https://shopnicekicks.com",
			},
			{
				Name: "Size? Canada",
				Url:  "https://size.ca",
			},
			{
				Name: "Sneaker Politics",
				Url:  "https://sneakerpolitics.com",
			},
			{
				Name: "Social Status",
				Url:  "https://www.socialstatuspgh.com",
			},
			{
				Name: "Solefly",
				Url:  "https://www.solefly.com",
			},
			{
				Name: "Sports World 165",
				Url:  "https://sportsworld165.com",
			},
			{
				Name: "Stussy",
				Url:  "https://www.stussy.com",
			},
			{
				Name: "Stussy UK",
				Url:  "https://www.stussy.co.uk",
			},
			{
				Name: "Telfar",
				Url:  "https://shop.telfar.net",
			},
			{
				Name: "The Darkside Initiative",
				Url:  "https://www.thedarksideinitiative.com",
			},
			{
				Name: "The Premier Store",
				Url:  "https://thepremierstore.com",
			},
			{
				Name: "Trophy Room",
				Url:  "https://www.trophyroomstore.com",
			},
			{
				Name: "Undefeated",
				Url:  "https://undefeated.com",
			},
			{
				Name: "Union LA",
				Url:  "https://store.unionlosangeles.com",
			},
			{
				Name: "YCMC",
				Url:  "https://www.ycmc.com",
			},
			{
				Name: "Xhibition",
				Url:  "https://www.xhibition.co",
			},
		},
		UsesAccounts: true,
	}

	for i := range meta.Stores {
		meta.Stores[i].Modes = []tasks.Mode{
			{
				Label: "Fast",
				Name:  "fast",
			},
			{
				Label: "Preload",
				Name:  "preload",
			},
			{
				Label: "Safe",
				Name:  "safe",
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
