/* eslint-disable @typescript-eslint/no-unused-vars */
import React from "react";
import { motion } from "framer-motion";
import styled from "styled-components";

import { Input, Select, IndicatorSeparator, Control } from "../../../elements";
import { ProfileUtils } from "../../../utils";

import Address from "./address";
import InputFormik from "../../../elements/Input/InputFormik";
import SelectFormik from "../../../elements/Select/SelectFormik";

interface BillingProps {
  values: any;
  errors: any;
  touched: any;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const { buildCountryOptions, buildProvinceOptions } = ProfileUtils;

const Billing = ({
  values,
  errors,
  touched,
  setFieldTouched,
  setFieldValue,
}: BillingProps) => {
  return (
    <motion.div>
      {/* Row 1 - Full Name */}
      <Row m="0">
        <InputFormik
          id="billing.name"
          name="billing.name"
          autoFocus
          useLabel
          textTransform="capitalize"
          maxLength={45}
          placeholder="Full Name"
          key="billing-name"
        />
      </Row>

      {/* Row 2 - Address / Apt */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0" f={3}>
          <Address
            values={values}
            errors={errors}
            touched={touched}
            type="billing"
            value={values.billing.line1}
            setFieldValue={setFieldValue}
            setFieldTouched={setFieldTouched}
          />
        </Col>

        <Col m="0 0 0 8px" f={1}>
          <InputFormik
            id="billing.line2"
            name="billing.line2"
            useLabel
            textTransform="capitalize"
            placeholder="Apt / Suite"
            key="billing-line2"
          />
        </Col>
      </Row>

      {/* Row 1 - Country / Province */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <SelectFormik
            key="billing-country"
            name="billing.country"
            required
            isClearable={false}
            label="Country / Region"
            placeholder="Country / Region"
            components={{
              IndicatorSeparator,
              Control,
            }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.code}
            options={buildCountryOptions()}
            onChange={(event: any) => {
              if (event !== values.billing.country) {
                if (values.payment.phone) {
                  setFieldValue("payment.phone", "");
                }
                setFieldValue("billing.state", null);
              }

              return setFieldValue("billing.country", event);
            }}
          />
        </Col>

        <Col m="0 0 0 8px">
          <SelectFormik
            key="billing-state"
            name="billing.state"
            label="State / Province"
            placeholder="State / Province"
            components={{
              IndicatorSeparator,
              Control,
            }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.code}
            options={buildProvinceOptions(values.billing.country)}
          />
        </Col>
      </Row>

      {/* Row 1 - City / Zip */}
      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <InputFormik
            id="billing.city"
            name="billing.city"
            useLabel
            textTransform="capitalize"
            placeholder="City"
            key="billing-city"
          />
        </Col>

        <Col m="0 0 0 8px">
          <InputFormik
            id="billing.postCode"
            name="billing.postCode"
            maxLength={10}
            useLabel
            placeholder="Zip Code"
            key="billing-zip"
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

export default Billing;
