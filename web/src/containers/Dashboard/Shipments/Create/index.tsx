import React, { Fragment, Dispatch, SetStateAction, useCallback } from "react";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";

import { Buttons } from "../../../../elements";
import { Modal } from "../../../../components";

import { Shipment } from "../../../../forms";
import { ConnectionHandler, useMutation } from "react-relay";
import { ROOT_ID } from "relay-runtime";
import { ShipmentTrackerAdd } from "../ShipmentTrackerAddMutation";
import { ShipmentTrackerAddMutation } from "../__generated__/ShipmentTrackerAddMutation.graphql";
import InputFormik from "../../../../elements/Input/InputFormik";

interface AddShipmentProps {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

interface AddShipmentForm {
  name: string;
  tracking: string;
}

const AddShipment = ({ open, setOpen }: AddShipmentProps) => {
  const [shipmentTrackerAdd] =
    useMutation<ShipmentTrackerAddMutation>(ShipmentTrackerAdd);

  const onSubmit = async (values: AddShipmentForm) => {
    const { name, tracking } = values;

    const connectionID = ConnectionHandler.getConnectionID(
      ROOT_ID,
      "Shipments_shipmentTrackers"
    );

    const config = {
      variables: {
        input: {
          name,
          trackingID: tracking,
        },
        connections: [connectionID],
      },
      onCompleted: () => {
        handleClear();
      },
    };

    shipmentTrackerAdd(config);
  };

  const formikbag = useFormik<AddShipmentForm>({
    isInitialValid: false,
    validateOnMount: true,
    ...Shipment,
    onSubmit,
  });

  const { resetForm, handleSubmit } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
  }, []);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        show={open}
        setShow={setOpen}
        title="Add shipment"
        width={325}
        height="auto"
      >
        <Fragment>
          <TextContainer useMargin={false}>
            <InputFormik
              autoFocus
              useLabel
              id="name"
              name="name"
              textTransform="capitalize"
              placeholder="Shipment Name"
            />
          </TextContainer>

          <TextContainer useMargin>
            <InputFormik
              useLabel
              id="tracking"
              name="tracking"
              placeholder="Tracking #"
              onSubmit={handleSubmit}
            />
          </TextContainer>

          <Row m="16px 0 0 0">
            <Buttons.Tertiary
              variant="IconButton"
              command="⌘ C"
              text="Clear"
              width={84}
              height={39}
              onClick={handleClear}
            />
            <Buttons.Primary
              variant="IconButton"
              command="↩︎"
              text="Save"
              width={84}
              height={39}
              onClick={handleSubmit}
            />
          </Row>
        </Fragment>
      </Modal>
    </FormikProvider>
  );
};

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
`;

const TextContainer = styled.div<{ useMargin: boolean }>`
  margin-top: ${({ useMargin }) => (useMargin ? 16 : 0)}px;
`;

export default AddShipment;
