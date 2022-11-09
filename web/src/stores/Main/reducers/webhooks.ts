import {
  createSelector,
  createSlice,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import moment from 'moment';

import { RootState } from ".";

export interface Field {
  name: string;
  enabled: boolean;
}

interface ProfileStripped {
  id: string;
  name: string;
}

export type Webhook = {
  id: string;
  name: string;
  active: boolean;
  type: string;
  url: string;
  profiles: ProfileStripped[];
  fields: string[];
  declines: boolean;
  sensitivity: boolean;
};

export type Webhooks = {
  [id: string]: Webhook;
};

export const initialWebhooksState: Webhooks = {};

export const addWebhook = createAsyncThunk(
  "webhooks/add",
  async (values: any) => {
    const enabledFields = values.fields.filter((f: Field) => f.enabled).map((f: Field) => f.name);
    const wh: Webhook = { ...values, fields: enabledFields };

    const id = await window.RPCAction('webhooks:add', [wh]);
    return { ...wh, id };
  }
);

export const testWebhook = createAsyncThunk(
  "webhooks/test",
  async (id: string) => {
    window.RPCAction('webhooks:sendToWebhook', [id, {
      date: moment().format('MM/DD/YYYY'),
      mode: "Normal",
      success: true,
      profile: "Test Profile",
      product: {
        name: "Test Product",
        price: "$100",
        image: "https://images-ext-2.discordapp.net/external/ouj9wWdCNDBTjuD_DHaeU9iBU7sdkmhRSw2XdfwCA-4/https/nebulabots.s3.amazonaws.com/nebula-logo.png",
        quantity: "1",
        size: "Random",
        url: "https://example.com",
      },
      store: {
        name: "Test Site",
        url: "https://example.com",
      },
    }]);

    return { id };
  }
);

export const enableWebhook = createAsyncThunk(
  "webhooks/enable",
  async (values: any) => {
    window.RPCAction('webhooks:edit', [values.id, values]);

    return { ...values };
  }
);

export const editWebhook = createAsyncThunk(
  "webhooks/edit",
  async (values: any) => {
    const enabledFields = values.fields.filter((f: Field) => f.enabled).map((f: Field) => f.name);
    const wh: Webhook = { ...values, fields: enabledFields };

    window.RPCAction('webhooks:edit', [wh.id, wh]);

    return { ...wh };
  }
);

export const removeWebhook = createAsyncThunk(
  "webhooks/remove",
  async (id: any) => {

    window.RPCAction('webhooks:delete', [id]);

    return { id };
  }
);

export const webhooks = createSlice({
  name: "webhooks",
  initialState: initialWebhooksState,
  reducers: {
    importWebhooks: (state, action) => {},
  },
  extraReducers: {
    //@ts-ignore
    [addWebhook.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      state[payload.id] = {
        ...payload
      };
    },
    //@ts-ignore
    [testWebhook.fulfilled]: (state, action) => {},
    //@ts-ignore
    [editWebhook.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      state[payload.id] = payload;
    },
    //@ts-ignore
    [enableWebhook.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      state[payload.id] = payload;
    },
    //@ts-ignore
    [removeWebhook.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { id } = payload;

      delete state[id];
    }
  },
});

export const selectWebhooks = ({ Webhooks }: RootState) => Webhooks;

export const makeWebhooks = createSelector(
  selectWebhooks,
  (state: Webhooks) => state
);

const { actions, reducer } = webhooks;

export const { importWebhooks } = actions;

export default reducer;
