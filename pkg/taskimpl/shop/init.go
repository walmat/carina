package shopify

import (
	"nebula/pkg/infra/tasks"
	_ "nebula/pkg/taskimpl/shop/monitor"
)

type State string

const (
	Setup             State = "setup"
	GetHomepage       State = "get-homepage"
	GetConfig         State = "get-config"
	GetPreloadProduct State = "get-preload-product"
	SubmitCartSafe    State = "submit-preload-cart"
	ClearPreloadCart  State = "clear-preload-cart"
	GetAccount        State = "get-account"
	SubmitAccount     State = "submit-account"
	SubmitPassword    State = "submit-password"
	GetCheckpoint     State = "get-checkpoint"
	SubmitCheckpoint  State = "submit-checkpoint"
	GetCaptchaToken   State = "get-captcha-token"
	GetChallenge      State = "get-challenge"
	SubmitChallenge   State = "submit-challenge"
	CreateCheckout    State = "create-checkout"
	EnterQueue        State = "enter-queue"
	GetQueue          State = "get-queue"
	GetPastQueue      State = "get-past-queue"
	GetNextQueue      State = "get-next-queue"
	WaitForProduct    State = "wait-for-product"
	GetCart           State = "get-cart"
	SubmitCart        State = "submit-cart"
	GetCustomer       State = "get-customer"
	SubmitCustomer    State = "submit-customer"
	GetShipping       State = "get-shipping"
	SubmitShipping    State = "submit-shipping"
	GetPaymentSession State = "get-payment-session"
	GetPayment        State = "get-payment"
	SubmitPayment     State = "submit-payment"
	GetReview         State = "get-review"
	SubmitReview      State = "submit-review"
	GetOrder          State = "get-order"
)

func init() {
	fastModeHandler, err := tasks.NewHandler(Setup, tasks.HandlerMap{
		Setup:             handleSetup,
		WaitForProduct:    handleWaitForProduct,
		GetHomepage:       handleGetHomepage,
		GetConfig:         handleGetConfig,
		GetCaptchaToken:   handleGetCaptchaToken,
		SubmitPassword:    handleSubmitPassword,
		GetAccount:        handleGetAccount,
		SubmitAccount:     handleSubmitAccount,
		GetChallenge:      handleGetChallenge,
		SubmitChallenge:   handleSubmitChallenge,
		CreateCheckout:    handleCreateCheckout,
		EnterQueue:        handleEnterQueueFast,
		GetQueue:          handleGetQueueFast,
		GetPastQueue:      handleGetPastQueueFast,
		GetNextQueue:      handleGetNextQueueFast,
		SubmitCustomer:    handleSubmitCustomer,
		GetCustomer:       handleGetCustomer,
		GetCart:           handleGetCart,
		SubmitCart:        handleSubmitCart,
		GetShipping:       handleGetShipping,
		SubmitShipping:    handleSubmitShipping,
		GetPayment:        handleGetPayment,
		GetPaymentSession: handleGetPaymentSession,
		SubmitPayment:     handleSubmitPayment,
		GetReview:         handleGetReview,
		SubmitReview:      handleSubmitReview,
		GetOrder:          handleGetOrder,
	})

	if err != nil {
		panic(err)
	}

	safeModeHandler, err := tasks.NewHandler(Setup, tasks.HandlerMap{
		Setup:             handleSetupSafe,
		GetHomepage:       handleGetHomepageSafe,
		GetConfig:         handleGetConfigSafe,
		GetCaptchaToken:   handleGetCaptchaTokenSafe,
		GetAccount:        handleGetAccountSafe,
		SubmitAccount:     handleSubmitAccountSafe,
		GetChallenge:      handleGetChallengeSafe,
		SubmitChallenge:   handleSubmitChallengeSafe,
		GetCheckpoint:     handleGetCheckpoint,
		SubmitCheckpoint:  handleSubmitCheckpoint,
		SubmitCartSafe:    handleSubmitCartSafe,
		EnterQueue:        handleEnterQueueSafe,
		GetQueue:          handleGetQueue,
		GetPastQueue:      handleGetPastQueue,
		GetNextQueue:      handleGetNextQueue,
		ClearPreloadCart:  handleClearPreloadCart,
		CreateCheckout:    handleCreateCheckoutSafe,
		WaitForProduct:    handleWaitForProductSafe,
		GetCustomer:       handleGetCustomerSafe,
		SubmitCustomer:    handleSubmitCustomerSafe,
		GetShipping:       handleGetShippingSafe,
		SubmitShipping:    handleSubmitShippingSafe,
		GetPayment:        handleGetPaymentSafe,
		GetPaymentSession: handleGetPaymentSessionSafe,
		SubmitPayment:     handleSubmitPaymentSafe,
		GetOrder:          handleGetOrderSafe,
	})
	preloadModeHandler, err := tasks.NewHandler(Setup, tasks.HandlerMap{
		Setup:             handleSetupPreload,
		GetHomepage:       handleGetHomepagePreload,
		GetConfig:         handleGetConfigPreload,
		GetCaptchaToken:   handleGetCaptchaTokenPreload,
		GetAccount:        handleGetAccountPreload,
		SubmitAccount:     handleSubmitAccountPreload,
		GetChallenge:      handleGetChallengePreload,
		SubmitChallenge:   handleSubmitChallengePreload,
		GetPreloadProduct: handleGetPreloadProduct,
		GetCheckpoint:     handleGetCheckpointPreload,
		SubmitCheckpoint:  handleSubmitCheckpointPreload,
		SubmitCartSafe:    handleSubmitCartPreload,
		EnterQueue:        handleEnterQueuePreload,
		GetQueue:          handleGetQueuePreload,
		GetPastQueue:      handleGetPastQueuePreload,
		GetNextQueue:      handleGetNextQueuePreload,
		ClearPreloadCart:  handleClearPreloadCart,
		CreateCheckout:    handleCreateCheckoutPreload,
		WaitForProduct:    handleWaitForProductPreload,
		GetCustomer:       handleGetCustomerPreload,
		SubmitCustomer:    handleSubmitCustomerPreload,
		GetShipping:       handleGetShippingPreload,
		SubmitShipping:    handleSubmitShippingPreload,
		GetPayment:        handleGetPaymentPreload,
		GetPaymentSession: handleGetPaymentSessionPreload,
		SubmitPayment:     handleSubmitPaymentPreload,
		GetOrder:          handleGetOrderPreload,
	})
	if err != nil {
		panic(err)
	}

	h := tasks.Handlers{
		"fast":    fastModeHandler,
		"preload": preloadModeHandler,
		"safe":    safeModeHandler,
	}

	if err = tasks.RegisterType("shopify", h); err != nil {
		panic(err)
	}
}
