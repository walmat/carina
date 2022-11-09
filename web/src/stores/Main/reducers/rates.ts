import {
  createSlice,
  createSelector,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { RootState } from ".";

import { Store } from "./stores";

export type Rate = {
  id: string;
  store: Store;
  carrier: "UPS" | "USPS" | "FedEx" | "DHL" | string;
  location: string;
  price: number;
  rate: string;
};

export type Rates = {
  [id: string]: Rate;
};

export const initialShippingRatesState: Rates = {};

export const addRates = createAsyncThunk(
  "rates/add",
  async ({ store, profile, rates }: any) => {
    let ratesMap: any = {};
    for (const rate of rates) {
      const { id } = await window.RPCAction('rates:add', [{
        store: store.name,
        profile,
        rate,
      }]);

      ratesMap[id] = {
        id,
        store,
        profile,
        rate,
      };
    }

    return { ratesMap };
  }
);

export const removeRates = createAsyncThunk(
  "rates/remove",
  async ({ ids, group }: any) => {
    ids.forEach((id: string) => {
      try {
        window.RPCAction('rates:delete', [id]);
      } catch (e) {
        console.log(e);
        // noop...
      }
    });

    return { ids, group };
  }
);

const fetchRates = createAsyncThunk("rates/fetch", async ({}) => {});

export const rates = createSlice({
  name: "rates",
  initialState: initialShippingRatesState,
  reducers: {
    importRates: (state, action) => {},
    editRates: (state, action) => {},
  },
  extraReducers: {
    // @ts-ignore
    [addRates.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { ratesMap } = payload;

      Object.values(ratesMap).forEach(({ id, ...rest }: any) => {
        // @ts-ignore
        state[id] = {
          ...rest,
          id,
        };
      });
    },
    // @ts-ignore
    [removeRates.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { ids } = payload;

      ids.forEach((id: string) => {
        delete state[id];
      });
    },
  },
});

export const selectRates = ({ Rates }: RootState) => Rates;

export const makeRates = createSelector(selectRates, (state: Rates) => state);

const { reducer, actions } = rates;

export const { importRates, editRates } = actions;

export default reducer;
