import React from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import {
  IndicatorSeparator,
  Control,
} from "../../../../elements";
import { Defaults, makeDefaults } from "../../../../stores/Main/reducers/defaults";
import { StoreObject, makeStores } from "../../../../stores/Main/reducers/stores";

import { stores as defaultStores } from "../../../../constants/stores";
import { modesForStoreUrl } from "../../../../utils/tasks";
import SelectFormik from '../../../../elements/Select/SelectFormik';
import InputFormik from '../../../../elements/Input/InputFormik';

interface GenericProps {
  values: any;
  errors: any;
  touched: any;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const compareArrays = (a1: any, a2: any) =>
  (a1 = new Set(a1)) &&
  (a2 = new Set(a2)) &&
  a1.size === a2.size &&
  [...a1].every((v) => a2.has(v));

const extractStores = (defaults: Defaults, stores: StoreObject[]) => {
  const existingUrls = Object.keys(defaults);

  return stores
    .reduce((acc: any[], curr: any) => [...acc, ...curr.options], [])
    .filter(({ url }) => !existingUrls.some((u) => u === url));
};

const Generic = ({
  values,
  setFieldValue,
}: GenericProps) => {
  const defaults = useSelector(makeDefaults);
  const stores = useSelector(makeStores);

  const storeOptions = extractStores(defaults, stores || defaultStores);

  return (
    <motion.div>
      <Row m="16px 0">
        <SelectFormik
          key="store"
          name="store"
          required
          label="Store / Platform"
          isClearable={false}
          placeholder="Choose Store"
          components={{ IndicatorSeparator, Control }}
          getOptionLabel={(option: any) => option.name}
          getOptionValue={(option: any) => option.url}
          onChange={(event: any) => {
            // NOTE: If we're changing the store and already have a mode/rate selected, we need to undo those
            if (values.store) {
              const oldModes = modesForStoreUrl(values.store.url);
              const newModes = modesForStoreUrl(event.url);

              // if they're the same, we don't have to do anything extra
              if (compareArrays(oldModes, newModes)) {
                return setFieldValue("store", event);
              }

              // if they're different, we need to unset the rate/mode
              setFieldValue("mode", "");
              setFieldValue("rate", null);
            }

            return setFieldValue("store", event);
          }}
          options={storeOptions}
        />
      </Row>

      <Row>
        <InputFormik
          key="name"
          id="name"
          name="name"
          textTransform="capitalize"
          useLabel
          placeholder="Default Name"
        />
      </Row>
    </motion.div>
  );
};

const Row = styled.div<{ m?: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => (m ? m : "0")};
`;

export default Generic;
