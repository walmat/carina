import {
  createAsyncThunk,
  createSelector,
  createSlice,
} from "@reduxjs/toolkit";

import { RootState } from ".";
import { Harvester, Harvesters } from "../../../containers/Collective/types";

export const add = createAsyncThunk("harvesters/add", async () => {
  return window.RPCAction("collective:add");
});

export const edit = createAsyncThunk(
  "harvesters/edit",
  async ({
    harvester,
    fieldsToUpdate,
  }: {
    harvester: Harvester;
    fieldsToUpdate: Partial<Harvester>;
  }) => {
    return window.RPCAction("collective:edit", [
      { ...harvester, ...fieldsToUpdate },
    ]);
  }
);

export const focus = createAsyncThunk(
  "harvesters/focus",
  async (id: string) => {
    return window.RPCAction("collective:focus", [id]);
  }
);

export const remove = createAsyncThunk(
  "harvesters/delete",
  async (id: string) => {
    return window.RPCAction("collective:delete", [id]);
  }
);

export const initialHarvestersState: Harvesters = {
  byId: {
    default: {
      id: "default",
      index: 0,
      focused: true,
      name: "Default",
      proxy: "",
    },
  },
};

export const harvesters = createSlice({
  name: "harvesters",
  initialState: initialHarvestersState,
  reducers: {
    logout: () => {},
  },
  extraReducers: {
    // @ts-ignore
    [add.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      // const focused = Object.values(state.byId).find(({ focused }) => focused);
      // if (focused) {
      //   focused.focused = false;
      // }

      state.byId[payload.id] = {
        ...payload,
        // focused: true,
      };
    },
    // @ts-ignore
    [edit.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      // @ts-ignore
      state.byId[payload.id] = payload;
    },
    // @ts-ignore
    [focus.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      // @ts-ignore
      Object.values(state.byId).forEach((h) => (h.focused = false));

      // eslint-disable-next-line no-param-reassign
      state.byId[payload].focused = true;
    },
    // @ts-ignore
    [remove.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      // TODO: We'll have to worry about removing the requester from the harvester too

      // NOTE: Shift focus to the default solver
      state.byId["default"].focused = true;
      delete state.byId[payload];

      return state;
    },
  },
});

const selectHarvesters = ({ Harvesters }: RootState) => Harvesters;

export const makeHarvesters = createSelector(
  selectHarvesters,
  (state: RootState["Harvesters"]) => {
    return Object.values(state.byId);
  }
);

export const makeFocusedHarvester = createSelector(
  selectHarvesters,
  (state: RootState["Harvesters"]) =>
    Object.values(state.byId).find(({ focused }) => focused)
);

export const HARVESTER_FIELDS = {
  NAME: "name",
  STORE: "store",
  ACCOUNT: "account",
  PROXY: "proxy",
  PLATFORM: "platform",
};

const { actions, reducer } = harvesters;

export const { logout } = actions;

export default reducer;
