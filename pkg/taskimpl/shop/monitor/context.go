package monitor

import (
	"golang.org/x/net/proxy"
	"log"
	"math/rand"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/third_party/http"
	tls "nebula/third_party/utls"
	"time"
)

type runningCtx struct {
	Client *http.SimpleClient
	ProxyData   *proxies.ProxyData
	UserAgent string
}

var (
	Ctx = runningCtx{Client: nil}
	UserAgents = [...]string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	}
)
func getUserAgent() string {
	return UserAgents[rand.Intn(len(UserAgents))]
}

func (c *runningCtx) getProxyDialer(taskCtx *tasks.Context) (proxy.ContextDialer, error) {
	proxyData, err := proxies.GetProxyForTask(taskCtx.Identifier())
	if err == nil {
		Ctx.ProxyData = proxyData
		return http.NewConnectDialer(proxyData.Url)
	} else if err == proxies.NotLinkedErr {
		return proxy.Direct, nil
	} else {
		return nil, err
	}
}

func (c *runningCtx) Init(taskContext *tasks.Context) error {
	log.Println("[DEBUG] Initializing monitor context")

	Ctx.UserAgent = getUserAgent()
	clientConfig := http.DefaultClientConfig.Clone()

	proxyDialer, err := Ctx.getProxyDialer(taskContext)
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

	client := http.NewSimpleClient(taskContext.Context, clientConfig, http.DefaultFingerprint)
	client.Inner().Jar = newSingleJar()
	client.Inner().Timeout = 15 * time.Second

	Ctx.Client = client
	return nil
}