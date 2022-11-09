import { createAsyncThunk, createSelector, createSlice } from "@reduxjs/toolkit";

import { RootState } from ".";

export type Credentials = {
  label: string;
  value: string;
};

export type Integration = {
  id: string;
  name: string;
  active: boolean;
  credentials?: Credentials[];
};

export type Integrations = {
  [id: string]: Integration;
};

export const initialIntegrationsState: Integrations = {
  aycd: {
    id: "aycd",
    active: false,
    credentials: [
      {
        label: "Access Token",
        value: "",
      },
      {
        label: "Api Key",
        value: "",
      },
    ],
    name: "AYCD AutoSolve",
  },
  "2captcha": {
    id: "2captcha",
    active: false,
    credentials: [
      {
        label: "Api Key",
        value: "",
      },
    ],
    name: "2Captcha",
  },
  capmonster: {
    id: "capmonster",
    active: false,
    credentials: [
      {
        label: "Api Key",
        value: "",
      },
    ],
    name: "CapMonster",
  },
  scout: {
    id: "scout",
    active: false,
    credentials: [],
    name: "Scout Analytics",
  },
};

export const integrations = createSlice({
  name: "integrations",
  initialState: initialIntegrationsState,
  reducers: {
    importIntegrations: (state, action) => {},
    saveIntegration: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { id, credentials } = payload;
      if (!id || !state[id]) {
        return;
      }

      switch (payload) {
        default:
        case "aycd":
          window.RPCAction('autosolve:connect', [...credentials.map((credentials: Credentials) => credentials.value)])
          break;
        case "2captcha":
          break;
        case "capmonster":
          break;
        case "scout":
          break;
      }

      state[id].active = true;
      state[id].credentials = credentials;
    },
    disableIntegration: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      switch (payload) {
        default:
        case "aycd":
          window.RPCAction('autosolve:close')
          break;
        case "2captcha":
          break;
        case "capmonster":
          break;
        case "scout":
          break;
      }

      state[payload].active = false;
    },
    stash: (state) => {
      (async () => {
        // @ts-ignore
        window.RPCAction('integrations:set', [state]);
      })();
    },
  },
});

export const selectIntegrations = ({ Integrations }: RootState) => Integrations;

export const makeIntegrations = createSelector(
  selectIntegrations,
  (state: Integrations) => state
);

const { actions, reducer } = integrations;

export const { importIntegrations, saveIntegration, disableIntegration, stash } =
  actions;

export default reducer;
