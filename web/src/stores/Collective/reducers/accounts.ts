import {
  createSelector,
  createSlice,
  createAsyncThunk,
} from "@reduxjs/toolkit";

import { RootState } from ".";

export interface Account {
  id: string;
  email: string;
  cookies: string[];
}

export interface Accounts {
  byId: {
    [id: string]: Account;
  };
}

export const add = createAsyncThunk("harvesterAccount/add", async () => {
  const harvester = await window.RPCAction("gmail:add");

  return { harvester };
});

export const select = createAsyncThunk(
  "harvesterAccount/select",
  async (id: string) => {
    const harvester = await window.RPCAction("gmail:select", [id]);

    return { harvester };
  }
);

export const remove = createAsyncThunk(
  "harvesterAccount/delete",
  async (id: string) => {
    const newHarvester = await window.RPCAction("gmail:delete", [id]);

    return { newHarvester };
  }
);

export const initialHarvesterAccountsState: Accounts = {
  byId: {},
};

export const harvesters = createSlice({
  name: "harvesterAccounts",
  initialState: initialHarvesterAccountsState,
  reducers: {},
  extraReducers: {
    // @ts-ignore
    [add.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { harvester } = payload;

      state.byId[harvester.id] = harvester;
    },
    // @ts-ignore
    [select.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { harvester } = payload;

      state.byId[harvester.id] = harvester;
    },
    // @ts-ignore
    [remove.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      delete state.byId[payload.newHarvester];

      return state;
    },
  },
});

const selectHarvesterAccounts = ({ Accounts }: RootState) => Accounts;

export const makeAccounts = createSelector(
  selectHarvesterAccounts,
  (state: RootState["Accounts"]) => {
    return state;
  }
);

export const makeAccountsOptions = createSelector(
  selectHarvesterAccounts,
  (state: RootState["Accounts"]) => Object.values(state.byId)
);

const { reducer } = harvesters;

export default reducer;
