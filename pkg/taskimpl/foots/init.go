package foots

import "nebula/pkg/infra/tasks"

type State string

const (
	Setup             State = "setup"
	GetSession        State = "get-session"
	GetStock          State = "get-stock"
	AddToCart         State = "add-to-cart"
	SubmitInformation State = "submit-information"
	SubmitOrder       State = "submit-order"
)

func init() {
	handler, err := tasks.NewHandler(Setup, tasks.HandlerMap{
		Setup:             handleSetup,
		GetSession:        handleGetSession,
		GetStock:          handleGetStock,
		AddToCart:         handleAddToCart,
		SubmitInformation: handleSubmitInformation,
		SubmitOrder:       handleSubmitOrder,
	})

	if err != nil {
		panic(err)
	}

	h := tasks.Handlers{
		"normal": handler,
	}

	if err = tasks.RegisterType("foots", h); err != nil {
		panic(err)
	}
}
