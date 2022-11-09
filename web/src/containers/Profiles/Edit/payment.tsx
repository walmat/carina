/* eslint-disable @typescript-eslint/no-unused-vars */
import React from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import { NumberFormatter, PhoneInput, IndicatorSeparator, Control } from '../../../elements';
import { makeProfiles } from "../../../stores/Main/reducers/profiles";
import InputFormik from "../../../elements/Input/InputFormik";
import SelectFormik from "../../../elements/Select/SelectFormik";


interface PaymentProps {
  values: any;
  errors: any;
  touched: any;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const Payment = ({
  values,
  errors,
  touched,
  handleChange,
  setFieldTouched,
  setFieldValue,
}: PaymentProps) => {
  const groups = useSelector(makeProfiles);

  return (
    <motion.div>
      {/* Row 1 - Email / Phone */}
      <Row m="0">
        <Col m="0 8px 0 0">
          <InputFormik
            autoFocus
            useLabel
            type="email"
            id="payment.email"
            name="payment.email"
            placeholder="Email Address"
            key="payment-email"
          />
        </Col>

        <Col m="0 0 0 8px">
          <PhoneInput
            id="payment.phone"
            name="payment.phone"
            useLabel
            country={values.shipping.country ? values.shipping.country.code : undefined}
            error={!!errors?.payment?.phone && !!touched?.payment?.phone}
            onBlur={() => setFieldTouched("payment.phone", true)}
            onChange={handleChange}
            value={values.payment.phone}
            placeholder="Phone Number"
            key="payment-phone"
          />
        </Col>
      </Row>

      {/* Row 2 - Holder / Number */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <InputFormik
            useLabel
            id="payment.name"
            name="payment.name"
            textTransform="capitalize"
            placeholder="Card Holder"
            key="payment-name"
          />
        </Col>

        <Col m="0 0 0 8px">
          <NumberFormatter
            id="payment.number"
            name="payment.number"
            useLabel
            error={!!errors?.payment?.number && !!touched?.payment?.number}
            onBlur={() => setFieldTouched("payment.number", true)}
            onChange={handleChange}
            format="#### #### #### #### ###"
            value={values.payment.number}
            placeholder="Card Number"
            key="payment-card-number"
          />
        </Col>
      </Row>

      {/* Row 1 - Country / Province */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <NumberFormatter
            id="payment.exp"
            name="payment.exp"
            useLabel
            error={!!errors?.payment?.exp && !!touched?.payment?.exp}
            onBlur={() => setFieldTouched("payment.exp", true)}
            onChange={handleChange}
            format="##/##"
            mask={["M", "M", "Y", "Y"]}
            value={values.payment.exp}
            placeholder="Card Expiration"
            key="payment-exp"
          />
        </Col>

        <Col m="0 0 0 8px">
          <InputFormik
            useLabel
            id="payment.cvv"
            name="payment.cvv"
            restriction="numerical"
            maxLength={4}
            placeholder="Security Code"
            key="payment-cvv"
          />
        </Col>
      </Row>

      {/* Row 1 - City / Zip */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <InputFormik
            useLabel
            id="name"
            name="name"
            textTransform="capitalize"
            maxLength={25}
            placeholder="Profile Name"
            key="profile-name"
          />
        </Col>

        <Col m="0 0 0 8px">
          <SelectFormik
            key="group"
            name="group"
            required
            isMulti={false}
            isClearable={false}
            label="Profile Group"
            placeholder="Profile Group"
            components={{ IndicatorSeparator, Control }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.id}
            options={Object.values(groups)}
            onChange={(event: any) => {
              if (!event) {
                return setFieldValue("group", null);
              }

              return setFieldValue("group", event);
            }}
          />
        </Col>
      </Row>
    </motion.div>
  );
};

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
`;

const Col = styled.div<{ m: string; f?: number }>`
  display: flex;
  flex: 1;
  flex-direction: column;
  margin: ${({ m }) => m};
  ${({ f }) => (typeof f !== undefined ? `flex: ${f};` : "")}
`;

export default Payment;
