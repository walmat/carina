import { graphql } from "react-relay";

export const ShipmentTrackerDelete = graphql`
  mutation ShipmentTrackerDeleteMutation(
    $input: ShipmentTrackerDeleteInput!
    $connections: [ID!]!
  ) {
    ShipmentTrackerDelete(input: $input) {
      error
      deletedShipmentTrackerId @deleteEdge(connections: $connections)
    }
  }
`;
