package gql

import "github.com/graphql-go/graphql"

// ShipmentTrackerDeleteMutation
var shipmentTrackerDeletePayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShipmentTrackerDeletePayload",
	Fields: graphql.Fields{
		"deletedShipmentTrackerId": &graphql.Field{
			Type: graphql.ID,
		},
		"error": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var shipmentTrackerDeleteInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ShipmentTrackerDeleteInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
})

var ShipmentTrackerDeleteMutation = &graphql.Field{
	Type:        shipmentTrackerDeletePayload,
	Description: "Delete a Shipment Tracker",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(shipmentTrackerDeleteInput),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		input := map[string]interface{}{}
		if inputVal, ok := p.Args["input"]; ok {
			if inputVal, ok := inputVal.(map[string]interface{}); ok {
				input = inputVal
			}
		}

		id := input["id"].(string)

		newShipmentTrackerList := []shipmentTrackerEntry{}

		hasDeleted := false

		for i := range ShipmentTrackerList {
			if ShipmentTrackerList[i].ID != id {
				newShipmentTrackerList = append(newShipmentTrackerList, ShipmentTrackerList[i])
			} else {
				hasDeleted = true
			}
		}

		if !hasDeleted {
			return map[string]interface{}{
				"deletedShipmentTrackerId": nil,
				"error":                    "id not found",
			}, nil
		}

		ShipmentTrackerList = newShipmentTrackerList

		return map[string]interface{}{
			"deletedShipmentTrackerId": id,
			"error":                    nil,
		}, nil
	},
}
