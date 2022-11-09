package gql

import (
	"github.com/graphql-go/graphql"
)

// root mutation
var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"ShipmentTrackerDelete": ShipmentTrackerDeleteMutation,
		"ShipmentTrackerAdd":    ShipmentTrackerAddMutation,
	},
})
