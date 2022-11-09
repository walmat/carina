import React from "react";
import {
  useFragment,
  graphql,
  useMutation,
  ConnectionHandler,
} from "react-relay";
import {
  Row,
  Col,
  ShipmentName,
  ShipmentStatus,
  TrackingNumber,
  RemoveIcon,
} from "./Shipments";
import { ShipmentRow_shipmentTracker$key } from "./__generated__/ShipmentRow_shipmentTracker.graphql";
import { ShipmentTrackerDelete } from "./ShipmentTrackerDeleteMutation";
import { ShipmentTrackerDeleteMutation } from "./__generated__/ShipmentTrackerDeleteMutation.graphql";
import { ROOT_ID } from "relay-runtime";

type Props = {
  shipmentTracker: ShipmentRow_shipmentTracker$key | null;
};
const ShipmentRow = (props: Props) => {
  const shipmentTracker = useFragment<ShipmentRow_shipmentTracker$key>(
    graphql`
      fragment ShipmentRow_shipmentTracker on ShipmentTracker {
        id
        name
        trackingID
      }
    `,
    props.shipmentTracker
  );

  const [shipmentTrackerDelete] = useMutation<ShipmentTrackerDeleteMutation>(
    ShipmentTrackerDelete
  );

  if (!shipmentTracker) {
    return null;
  }

  const delivered = /delivered/i.test(shipmentTracker.name || "");

  const onRemove = () => {
    const connectionID = ConnectionHandler.getConnectionID(
      ROOT_ID,
      "Shipments_shipmentTrackers"
    );

    const config = {
      variables: {
        input: {
          id: shipmentTracker.id,
        },
        connections: [connectionID],
      },
    };

    shipmentTrackerDelete(config);
  };

  return (
    <Row m="8px 0" useBorder={false}>
      <Col m="0 auto 0 0">
        <Row useBorder={false}>
          <ShipmentName>{shipmentTracker.name?.slice(0, 15)}</ShipmentName>
        </Row>
        <Row useBorder={false}>
          <TrackingNumber>
            {shipmentTracker.trackingID?.slice(0, 30)}
          </TrackingNumber>
        </Row>
      </Col>
      <Col m="0 8px">
        <ShipmentStatus delivered={delivered}>
          {shipmentTracker.name?.slice(0, 5)}
        </ShipmentStatus>
      </Col>
      <Col>
        <RemoveIcon onClick={onRemove} />
      </Col>
    </Row>
  );
};

export default ShipmentRow;
