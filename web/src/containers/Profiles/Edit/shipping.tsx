/* eslint-disable @typescript-eslint/no-unused-vars */
import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { Select, IndicatorSeparator, Control } from "../../../elements";

import Address from "./address";

import { ProfileUtils } from "../../../utils";
import InputFormik from "../../../elements/Input/InputFormik";
import SelectFormik from "../../../elements/Select/SelectFormik";

interface ShippingProps {
  values: any;
  errors: any;
  touched: any;
  matches: boolean;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const { buildCountryOptions, buildProvinceOptions } = ProfileUtils;

const Shipping = ({
  values,
  errors,
  touched,
  matches,
  setFieldTouched,
  setFieldValue,
}: ShippingProps) => {
  return (
    <motion.div>
      {/* Row 1 - Full Name */}
      <Row m="0">
        <InputFormik
          id="shipping.name"
          name="shipping.name"
          autoFocus
          useLabel
          textTransform="capitalize"
          maxLength={45}
          placeholder="Full Name"
          key="shipping-name"
        />
      </Row>

      {/* Row 2 - Address / Apt */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0" f={3}>
          <Address
            values={values}
            errors={errors}
            touched={touched}
            type="shipping"
            value={values.shipping.line1}
            setFieldValue={setFieldValue}
            setFieldTouched={setFieldTouched}
          />
        </Col>

        <Col m="0 0 0 8px" f={1}>
          <InputFormik
            id="shipping.line2"
            name="shipping.line2"
            textTransform="capitalize"
            maxLength={20}
            useLabel
            placeholder="Apt / Suite"
            key="shipping-apt"
          />
        </Col>
      </Row>

      {/* Row 1 - Country / Province */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <SelectFormik
            key="shipping-country"
            name="shipping.country"
            required
            label="Country / Region"
            isClearable={false}
            placeholder="Country / Region"
            components={{ IndicatorSeparator, Control }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.code}
            options={buildCountryOptions()}
            onChange={(event: any) => {
              if (event !== values.shipping.country) {
                if (values.payment.phone && matches) {
                  setFieldValue("payment.phone", "");
                }

                setFieldValue("shipping.state", null);
              }

              return setFieldValue("shipping.country", event);
            }}
          />
        </Col>

        <Col m="0 0 0 8px">
          <SelectFormik
            key="shipping-state"
            name="shipping.state"
            label="State / Province"
            placeholder="State / Province"
            components={{ IndicatorSeparator, Control }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.code}
            options={buildProvinceOptions(values.shipping.country)}
          />
        </Col>
      </Row>

      {/* Row 1 - City / Zip */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <InputFormik
            id="shipping.city"
            name="shipping.city"
            textTransform="capitalize"
            restriction="alpha"
            maxLength={30}
            useLabel
            placeholder="City"
            key="shipping-city"
          />
        </Col>

        <Col m="0 0 0 8px">
          <InputFormik
            id="shipping.postCode"
            name="shipping.postCode"
            maxLength={10}
            useLabel
            placeholder="Postal Code"
            key="shipping-zip"
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

export default Shipping;
