package ys

import (
	"context"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/util"
	"sync"
	"sync/atomic"
	"time"
)

var (
	fetchingCaptcha    int32
	captchaPublishChan = make(chan harvesters.SolverResponse)

	captchaCallbacks      []chan harvesters.SolverResponse
	captchaCallbacksMutex sync.Mutex

	cachedCaptchaResponse     *harvesters.SolverResponse
	cachedCaptchaResponseTime time.Time
)

func init() {
	go captchaShareHandler()
}

func fetchNewToken(taskCtx *tasks.Context) {
	if atomic.CompareAndSwapInt32(&fetchingCaptcha, 0, 1) {
		defer atomic.SwapInt32(&fetchingCaptcha, 0)
		for {
			solveChans, err := harvesters.RequestSolve(harvesters.Solver{
				ID: taskCtx.Identifier(),
				ReCaptcha: &harvesters.ReCaptcha{
					SiteKey: "6Lf34M8ZAAAAANgE72rhfideXH21Lab333mdd2d-",
					Action:  "yzysply_wr_pageview",
				},
				Store: &harvesters.Store{
					Id:       "yeezysupply",
					Platform: "YEEZYSUPPLY",
					Name:     "Yeezy Supply",
					Url:      "https://www.yeezysupply.com/",
				},
				Type:        harvesters.ReCaptcha_V3,
				RequestedAt: time.Now().Unix(),
			})

			if err != nil {
				taskCtx.SendStatusColored("Harvester Error, Retrying...", util.ColorFailed)
				time.Sleep(time.Second * 1)
				continue
			}

			waitCtx, cancelFunc := context.WithTimeout(taskCtx.Context, time.Second*30)

			tokenChan := make(chan harvesters.SolverResponse)

			go func() {
				token, err := harvesters.GetSolve(solveChans)
				if err != nil {
					taskCtx.SendStatusColored("Harvester Error", util.ColorFailed)
				}

				tokenChan <- token
			}()

			select {
			case c := <-tokenChan:
				captchaPublishChan <- c
				cancelFunc()
				return
			case <-waitCtx.Done():
				cancelFunc()
				return
			}
		}
	}
}

func getCachedCaptchaToken() *harvesters.SolverResponse {
	if cachedCaptchaResponse != nil && time.Since(cachedCaptchaResponseTime).Seconds() > 100 {
		cachedCaptchaResponse = nil
	}
	return cachedCaptchaResponse
}

func newCaptchaShareChan() chan harvesters.SolverResponse {
	newChan := make(chan harvesters.SolverResponse)
	captchaCallbacksMutex.Lock()
	captchaCallbacks = append(captchaCallbacks, newChan)
	captchaCallbacksMutex.Unlock()
	return newChan
}

func captchaShareHandler() {
	for {
		c := <-captchaPublishChan

		cachedCaptchaResponse = &c
		// TODO: base this off of captcha token time
		cachedCaptchaResponseTime = time.Now()

		captchaCallbacksMutex.Lock()
		for _, out := range captchaCallbacks {
			out <- c
		}
		captchaCallbacks = nil
		captchaCallbacksMutex.Unlock()
	}
}
