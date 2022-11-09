import { createSelector, createSlice } from "@reduxjs/toolkit";
import { startCase } from "lodash";
import { RootState } from ".";

export type Platform = "Shopify" | "Supreme" | string;

export type Mode = {
  name: string;
  label: string;
};

export type Store = {
  id: string;
  platform: string;
  name: string;
  url: string;
  modes: Mode[];
};

export type StoreObject = {
  options: Store[];
  type: string;
  index: number;
  label: string;
  usesRates: boolean;
  usesAccounts: boolean;
};

export const initialStoresState: StoreObject[] = [];
export const stores = createSlice({
  name: "stores",
  initialState: initialStoresState,
  reducers: {
    resetStores: () => initialStoresState,
    addStore: (state, action) => {
      if (!action) {
        return;
      }

      // @ts-ignore
      const { id, platform, name, url, modes } = action;
      if (!name || !url) {
        return;
      }

      const index = state.findIndex(({ label }) => label === "Shopify");
      if (index === -1) {
        return;
      }

      if (state[index].options.some(({ url: _url }) => _url === url)) {
        return;
      }

      state[index].options.push({
        id,
        platform,
        name: startCase(name),
        url,
        modes,
      });

      state[index].options.sort((a, b) => {
        const nameA = a.name.toLowerCase();
        const nameB = b.name.toLowerCase();
        if (nameA < nameB) return -1;
        if (nameA > nameB) return 1;
        return 0;
      });
    },
    setStores: (_, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      // eslint-disable-next-line no-param-reassign
      return payload;
    },
  },
});

export const selectStores = (state: RootState) => state.Stores;

export const makeStores = createSelector(
  selectStores,
  (state: RootState["Stores"]) => state
);

const { actions, reducer } = stores;

export const { resetStores, addStore, setStores } = actions;

export default reducer;
