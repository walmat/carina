package ys

import (
	"nebula/pkg/infra/tasks"
)

type State string

const (
	GetHomepage         State = "get-homepage"
	GetPixel            State = "get-pixel"
	GetPayload          State = "get-payload"
	SubmitPixel         State = "submit-pixel"
	GetBloom            State = "get-bloom"
	GetProductPage      State = "get-product-page"
	GetWaitingRoom      State = "get-waiting-room"
	WaitInSplash        State = "wait-in-splash"
	CaptchaRequestState State = "captcha-request"
	CaptchaWaitState    State = "captcha-wait"
	GetProductDetails   State = "get-product-details"
	GetProductStock     State = "get-product-stock"
	GetAkamai           State = "get-akamai"
	GetSensor           State = "get-sensor"
	SubmitSensor        State = "submit-sensor"
	GetBasket           State = "get-basket"
	AddToCart           State = "add-to-cart"
	SubmitInformation   State = "submit-information"
	SubmitOrder         State = "submit-order"
	LaunchBrowser       State = "solve3ds-browser"
	WaitFor3DS          State = "wait-for-close"
	CompleteOrder       State = "complete-order"
)

func init() {
	handler, err := tasks.NewHandler(GetHomepage, tasks.HandlerMap{
		GetHomepage:         handleGetHomepage,
		GetPixel:            handleGetPixel,
		GetPayload:          handleGetPayload,
		GetAkamai:           handleGetAkamai,
		GetBloom:            handleGetBloom,
		SubmitPixel:         handleSubmitPixel,
		GetProductPage:      handleGetProductPage,
		GetWaitingRoom:      handleGetWaitingRoom,
		WaitInSplash:        handleWaitInSplash,
		CaptchaRequestState: handleCaptchaRequest,
		CaptchaWaitState:    handleCaptchaWait,
		GetProductDetails:   handleGetProductDetails,
		GetProductStock:     handleGetProductStock,
		GetSensor:           handleGetSensor,
		SubmitSensor:        handleSubmitSensor,
		GetBasket:           handleGetBasket,
		AddToCart:           handleAddToCart,
		SubmitInformation:   handleSubmitInformation,
		SubmitOrder:         handleSubmitOrder,
		LaunchBrowser:       handleLaunchBrowser,
		WaitFor3DS:          handleWaitFor3DS,
		CompleteOrder:       handleCompleteOrder,
	})

	if err != nil {
		panic(err)
	}

	h := tasks.Handlers{
		"normal": handler,
	}

	if err = tasks.RegisterType("ys", h); err != nil {
		panic(err)
	}
}
