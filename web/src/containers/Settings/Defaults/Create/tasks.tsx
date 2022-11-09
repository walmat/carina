import React from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import {
  Input,
  Select,
  IndicatorSeparator,
  Control,
} from "../../../../elements";
import { sizes } from "../../../../constants";

import { modesForStoreUrl} from '../../../../utils/tasks';
import { makeProxies } from "../../../../stores/Main/reducers/proxies";
import {makeAccounts} from "../../../../stores/Main/reducers/accounts";

interface TasksProps {
  values: any;
  errors: any;
  touched: any;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const Tasks = ({
  values,
  errors,
  touched,
  handleChange,
  setFieldTouched,
  setFieldValue,
}: TasksProps) => {
  const proxies = useSelector(makeProxies);
  const accounts = useSelector(makeAccounts);

  return (
    <motion.div>
      <Row m="0">
        <Col m="0 8px 0 0">
          <Select
            key="tasks-mode"
            name="mode"
            isMulti={false}
            isClearable={false}
            label="Task Mode"
            placeholder="Task Mode"
            error={!!errors?.mode}
            touched={!!touched?.mode}
            components={{ IndicatorSeparator, Control }}
            value={
              values.mode ? { label: values.mode, value: values.mode } : null
            }
            onBlur={() => setFieldTouched("mode", true)}
            onChange={(event: any) => {
              if (!event) {
                return setFieldValue("mode", null);
              }

              return setFieldValue("mode", event.value);
            }}
            options={
              values.store
                ? modesForStoreUrl(values.store.url).map((mode: string) => ({
                    label: mode,
                    value: mode,
                  }))
                : []
            }
          />
        </Col>

        <Col m="0 0 0 8px">
          <Select
            key="proxies"
            name="proxies"
            isMulti={false}
            isClearable
            label="Proxy Group"
            placeholder="Proxy Group"
            error={!!errors?.proxies}
            touched={!!touched?.proxies}
            components={{ IndicatorSeparator, Control }}
            value={values.proxies}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.id}
            onBlur={() => setFieldTouched("proxies", true)}
            onChange={(event: any) => {
              if (!event) {
                return setFieldValue("proxies", null);
              }

              return setFieldValue("proxies", event);
            }}
            options={Object.values(proxies)}
          />
        </Col>
      </Row>

      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <Select
            required
            isMulti
            isClearable
            isCreatable
            name="sizes"
            placeholder="Sizes"
            label="Desired Sizes"
            error={!!errors?.sizes}
            touched={!!touched?.sizes}
            components={{ Control }}
            value={values.sizes.map((size: string) => ({
              label: size,
              value: size,
            }))}
            onFocus={() => setFieldTouched("sizes", true)}
            onBlur={() => setFieldTouched("sizes", true)}
            onChange={(event: any) => {
              if (!event) {
                return setFieldValue("sizes", []);
              }

              return setFieldValue(
                "sizes",
                event.map(({ value }: any) => value)
              );
            }}
            options={sizes}
          />
        </Col>

        <Col m="0 0 0 8px">
          <Input
            autoFocus
            useLabel
            key="tasks-amount"
            id="amount"
            name="amount"
            error={!!errors?.amount}
            touched={!!touched?.amount}
            onBlur={() => setFieldTouched("amount", true)}
            onChange={handleChange}
            value={values.amount}
            placeholder="Task Amount"
          />
        </Col>
      </Row>

      <Row m="16px 0 0 0">
        <Col m="0 8px 0 0">
          <Select
            required
            isMulti
            isClearable
            closeMenuOnSelect={false}
            name="accounts"
            placeholder="Accounts"
            label="Accounts"
            error={!!errors?.accounts}
            touched={!!touched?.accounts}
            components={{ Control }}
            value={values.accounts}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.id}
            onFocus={() => setFieldTouched("accounts", true)}
            onBlur={() => setFieldTouched("accounts", true)}
            onChange={(event: any) => {
              if (!event) {
                return setFieldValue("accounts", []);
              }

              return setFieldValue("accounts", event);
            }}
            options={Object.values(accounts)}
          />
        </Col>

        <Col m="0 0 0 8px">
          <Input
            autoFocus
            useLabel
            key="accountLimit"
            id="accountLimit"
            name="accountLimit"
            error={!!errors?.accountLimit}
            touched={!!touched?.accountLimit}
            onBlur={() => setFieldTouched("accountLimit", true)}
            onChange={handleChange}
            value={values.accountLimit}
            placeholder="Checkout Limit (Account)"
          />
        </Col>
      </Row>

      <Row m="16px 0 0 0">
        <Select
          required
          isClearable
          closeMenuOnSelect={false}
          name="rate"
          placeholder="Shipping Rate"
          label="Shipping Rate"
          error={!!errors?.rate}
          touched={!!touched?.rate}
          components={{ Control }}
          value={values.rate}
          onFocus={() => setFieldTouched("rate", true)}
          onBlur={() => setFieldTouched("rate", true)}
          onChange={(event: any) => {
            if (!event) {
              return setFieldValue("rate", []);
            }

            return setFieldValue("rate", event);
          }}
          options={[]}
        />
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

export default Tasks;
