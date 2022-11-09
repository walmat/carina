package gql

import "github.com/graphql-go/graphql"

var shipment1 = shipmentTrackerEntry{
	ID:         "1",
	Name:       "Shipment1",
	TrackingID: "TrackingID1",
}

var shipment2 = shipmentTrackerEntry{
	ID:         "2",
	Name:       "Shipment2",
	TrackingID: "TrackingID2",
}

var ShipmentTrackerList = []shipmentTrackerEntry{
	shipment1, shipment2,
}

// array to mock shipment tracker
type shipmentTrackerEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TrackingID string `json:"trackingID"`
}

var ShipmentTrackerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShipmentTracker",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"trackingID": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var ShipmentTrackerEdgeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShipmentTrackerEdge",
	Fields: graphql.Fields{
		"node": &graphql.Field{
			Type: ShipmentTrackerType,
		},
		"cursor": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ShipmentTrackerConnectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ShipmentTrackerConnection",
	Fields: graphql.Fields{
		"pageInfo": &graphql.Field{
			Type: graphql.NewNonNull(PageInfoType),
		},
		"edges": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(ShipmentTrackerEdgeType))),
		},
	},
})
