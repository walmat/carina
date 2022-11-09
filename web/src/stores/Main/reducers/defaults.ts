import {
  createAsyncThunk,
  createSlice,
  createSelector,
  nanoid,
} from "@reduxjs/toolkit";
import { RootState } from ".";
import { ProfileGroup } from "./profiles";
import { Proxies } from "./proxies";
import { Accounts } from "./accounts";
import { Rates } from "./rates";
import { Store } from "./stores";

export interface DefaultEntry {
  id: string;
  name: string;
  store: Store;
  sizes: string[];
  profiles: ProfileGroup;
  proxies: Proxies;
  accounts: Accounts;
  rates: Rates;
}

export interface Default {
  name: string;
  id: string;
  byId: {
    [id: string]: DefaultEntry;
  };
}

export interface Defaults {
  [url: string]: Default;
}

export const initialDefaultsState: Defaults = {};

export const addDefault = createAsyncThunk(
  "default/add",
  async (values: any) => {
    // TODO: Replace by backend id
    let id = nanoid();
    try {
      id = await window.RPCAction('defaults:add', [values]);
    } catch (e) {
      // noop...
    }

    return { id, values: { ...values, sizes: values.sizes.join(", ") } };
  }
);

export const removeDefault = createAsyncThunk(
  "default/remove",
  async ({ ids }: any, store) => {
    const { Defaults }: any = store.getState();

    const groupMap: any = {};

    for (const id of ids) {
      const group: any = Object.values(Defaults).find(
        (group: any) => group.byId[id]
      );
      if (group) {
        try {
          await window.RPCAction('defaults:delete', [id, group.id]);
        } catch (e) {
          // noop...
        }

        if (!groupMap[group.id]) {
          groupMap[group.id] = [];
        }

        groupMap[group.id].push(id);
      }
    }

    return groupMap;
  }
);

export const copyDefault = createAsyncThunk(
  'default/copy',
  async ({ ids }: any, store) => {
    const { Defaults }: any = store.getState();

    const groupMap: any = {};

    for (const id of ids) {
      const group: any = Object.values(Defaults).find(
        (group: any) => group.byId[id]
      );

      if (group) {
        // TODO: Replace by backend id
        let _id = nanoid();
        try {
          console.log(group.byId[id])

          // _id = await defaultsHelper.addDefault(group.b);
        } catch (e) {
          // noop...
        }

        if (!groupMap[group.id]) {
          groupMap[group.id] = [];
        }

        groupMap[group.id].push({
          ...group.byId[id],
          id: _id
        });
      }
    }

    return { groupMap };
  },
)

export const defaults = createSlice({
  name: "defaults",
  initialState: initialDefaultsState,
  reducers: {},
  extraReducers: {
    // @ts-ignore
    [addDefault.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { id, values } = payload;

      // if store exists already
      if (state[values.store.url]) {
        state[values.store.url].byId[id] = {
          ...values,
          id,
        };
      } else {
        // create the base
        state[values.store.url] = {
          name: values.store.name,
          id: values.store.url,
          byId: {
            [id]: {
              ...values,
              id,
            },
          },
        };
      }
    },
    // @ts-ignore
    [removeDefault.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      Object.entries(payload).map(([url, ids]) => {
        // @ts-ignore
        ids.forEach((id: string) => {
          delete state[url].byId[id];

          if (!Object.values(state[url].byId).length) {
            delete state[url];
          }
        });
      });
    },
    // @ts-ignore
    [copyDefault.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { groupMap } = payload;

      console.log(groupMap);

      Object.keys(groupMap).forEach(group => {
        groupMap[group].forEach((d: Default) => {
          console.log(state[group])
          // @ts-ignore
          state[group].byId[d.id] = d
        });
      });
    },
  },
});

export const selectDefaults = ({ Defaults }: RootState) => Defaults;

export const makeDefaults = createSelector(
  selectDefaults,
  (state: Defaults) => state
);

const { reducer } = defaults;

export default reducer;
