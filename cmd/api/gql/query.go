package gql

import "github.com/graphql-go/graphql"

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/api/graphql?query={lastTodo{id,text,done}}'
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"shipmentTrackers": &graphql.Field{
			Type: graphql.NewNonNull(ShipmentTrackerConnectionType),
			Args: ConnectionArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert args map[string]interface into ConnectionArguments
				args := NewConnectionArguments(p.Args)

				data := make([]interface{}, len(ShipmentTrackerList))
				for i, v := range ShipmentTrackerList {
					data[i] = v
				}
				// let relay library figure out the result, given
				// - the list of ships for this faction
				// - and the filter arguments (i.e. first, last, after, before)
				return ConnectionFromArray(data, args), nil
			},
		},
		"checkouts": &graphql.Field{
			Type: graphql.NewNonNull(CheckoutConnectionType),
			Args: ConnectionArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert args map[string]interface into ConnectionArguments
				args := NewConnectionArguments(p.Args)

				data := make([]interface{}, len(CheckoutList))
				for i, v := range CheckoutList {
					data[i] = v
				}
				// let relay library figure out the result, given
				// - the list of ships for this faction
				// - and the filter arguments (i.e. first, last, after, before)
				return ConnectionFromArray(data, args), nil
			},
		},
		"releases": &graphql.Field{
			Type: graphql.NewNonNull(ReleaseConnectionType),
			Args: graphql.FieldConfigArgument{
				"before": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"after": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"first": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"last": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"filters": &graphql.ArgumentConfig{
					Type: ReleaseFilterInput,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// convert args map[string]interface into ConnectionArguments
				args := NewConnectionArguments(p.Args)

				data := make([]interface{}, len(ReleaseList))
				for i, v := range ReleaseList {
					data[i] = v
				}
				// let relay library figure out the result, given
				// - the list of ships for this faction
				// - and the filter arguments (i.e. first, last, after, before)
				return ConnectionFromArray(data, args), nil
			},
		},
	},
})
