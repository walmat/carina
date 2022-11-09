package gql

import "github.com/graphql-go/graphql"

var StoreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Store",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"mode": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var StoreEdgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StoreEdge",
	Fields: graphql.Fields{
		"node": &graphql.Field{
			Type: StoreType,
		},
		"cursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var StoreConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StoreConnection",
	Fields: graphql.Fields{
		"pageInfo": &graphql.Field{
			Type: graphql.NewNonNull(PageInfoType),
		},
		"edges": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(StoreEdgeType))),
		},
	},
})
