import { graphql } from "react-relay";

export const ShipmentTrackerAdd = graphql`
  mutation ShipmentTrackerAddMutation(
    $input: ShipmentTrackerAddInput!
    $connections: [ID!]!
  ) {
    ShipmentTrackerAdd(input: $input) {
      error
      shipmentTrackerEdge @appendEdge(connections: $connections) {
        __typename
        cursor
        node {
          id
          name
          trackingID
        }
      }
    }
  }
`;
