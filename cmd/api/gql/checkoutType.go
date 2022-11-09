package gql

import "github.com/graphql-go/graphql"

var CheckoutList = []checkoutEntry{}

type Product struct {
	Image string `json:"image"`
}

// array to mock checkout
type checkoutEntry struct {
	ID      string  `json:"id"`
	Date    string  `json:"date"`
	Amount  float32 `json:"amount"`
	Product Product `json:"product"`
}

var ProductType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"image": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"url": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"size": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"price": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var CheckoutMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CheckoutMetadata",
	Fields: graphql.Fields{
		"orderNumber": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"profileName": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"cardType": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"email": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var CheckoutType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Checkout",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"status": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"date": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"amount": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"product": &graphql.Field{
			Type: graphql.NewNonNull(ProductType),
		},
		"proxy": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"stores": &graphql.Field{
			Type: graphql.NewNonNull(StoreType),
		},
		"metadata": &graphql.Field{
			Type: graphql.NewNonNull(CheckoutMetadataType),
		},
	},
})

var CheckoutEdgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CheckoutEdge",
	Fields: graphql.Fields{
		"node": &graphql.Field{
			Type: CheckoutType,
		},
		"cursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var CheckoutConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CheckoutConnection",
	Fields: graphql.Fields{
		"pageInfo": &graphql.Field{
			Type: graphql.NewNonNull(PageInfoType),
		},
		"edges": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(CheckoutEdgeType))),
		},
	},
})
