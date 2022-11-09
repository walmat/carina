package gql

import (
	"github.com/graphql-go/graphql"
)

var ReleaseList = []releaseEntry{}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

// array to mock release
type releaseEntry struct {
	ID     string  `json:"id"`
	Date   string  `json:"date"`
	Amount float32 `json:"amount"`
}

var ReleaseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Release",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"image": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"date": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"stores": &graphql.Field{
			Type: graphql.NewNonNull(StoreConnectionType),
			Args: ConnectionArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert args map[string]interface into ConnectionArguments
				args := NewConnectionArguments(p.Args)

				stores := []interface{}{}
				// if release, ok := p.Source.(*releaseEntry); ok {
				// for _, stores := range release. {
				// 	stores = append(stores, stores)
				// }
				// }

				// let relay library figure out the result, given
				// - the list of ships for this faction
				// - and the filter arguments (i.e. first, last, after, before)
				return ConnectionFromArray(stores, args), nil
			},
		},
	},
})

var ReleaseEdgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ReleaseEdge",
	Fields: graphql.Fields{
		"node": &graphql.Field{
			Type: ReleaseType,
		},
		"cursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ReleaseConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ReleaseConnection",
	Fields: graphql.Fields{
		"pageInfo": &graphql.Field{
			Type: graphql.NewNonNull(PageInfoType),
		},
		"edges": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(ReleaseEdgeType))),
		},
	},
})

var DateFilterInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "DateFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"begin": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"end": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var ReleaseFilterInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ReleaseFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"date": &graphql.InputObjectFieldConfig{
			Type: DateFilterInput,
		},
	},
})
