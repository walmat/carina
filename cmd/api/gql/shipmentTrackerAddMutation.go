package gql

import (
	"github.com/graphql-go/graphql"
	"math/rand"
)

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var shipmentTrackerAddPayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShipmentTrackerAddPayload",
	Fields: graphql.Fields{
		"shipmentTrackerEdge": &graphql.Field{
			Type: ShipmentTrackerEdgeType,
		},
		"error": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var shipmentTrackerAddInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ShipmentTrackerAddInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"trackingID": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var ShipmentTrackerAddMutation = &graphql.Field{
	Type:        shipmentTrackerAddPayload,
	Description: "Add a Shipment Tracker",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(shipmentTrackerAddInput),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := map[string]interface{}{}
		if inputVal, ok := p.Args["input"]; ok {
			if inputVal, ok := inputVal.(map[string]interface{}); ok {
				input = inputVal
			}
		}

		name := input["name"].(string)
		trackingID := input["trackingID"].(string)

		newID := randStringRunes(8)

		newShipmentTracker := shipmentTrackerEntry{
			ID:         newID,
			Name:       name,
			TrackingID: trackingID,
		}

		ShipmentTrackerList = append(ShipmentTrackerList, newShipmentTracker)

		return map[string]interface{}{
			"shipmentTrackerEdge": map[string]interface{}{
				"node":   newShipmentTracker,
				"cursor": "cursor",
			},
			"error": nil,
		}, nil
	},
}
