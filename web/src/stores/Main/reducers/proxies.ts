import {
  createAsyncThunk,
  createSelector,
  createSlice,
} from "@reduxjs/toolkit";

import { RootState } from ".";
import { notify } from "../../../utils";

export type Proxy = {
  id: string;
  active: boolean;
  testing: boolean;
  speed: null | number;
  host: string;
  port: string;
  raw: string;
  username: null | string;
  password: null | string;
};

export type ProxyGroup = {
  id: string;
  name: string;
  byId: {
    [id: string]: Proxy;
  };
};

export type Proxies = {
  [id: string]: ProxyGroup;
};

export const initialProxiesState: Proxies = {
  default: {
    id: "default",
    name: "Default",
    byId: {},
  },
};

export const testProxies = createAsyncThunk(
  "proxies/test",
  async (payload: any) => {
    const { ids, platforms, groupId } = payload;

    window.RPCAction('proxies:test', [ids, groupId, platforms]);

    return { ids, group: groupId };
  }
);

export const addProxiesGroup = createAsyncThunk(
  "proxies/addGroup",
  async (groupName: string) => {
    if (!groupName.trim()) {
      return null;
    }

    const id = await window.RPCAction('proxies:addGroup', [groupName.trim()]);
    return { id, groupName: groupName.trim() };
  }
);

export const removeProxiesGroup = createAsyncThunk(
  "proxies/removeGroup",
  async (groups: string[]) => {
    if (!groups?.length) {
      return null;
    }

    groups.forEach((id) => {
      try {
        window.RPCAction("proxies:deleteGroup", [id]);
      } catch (e) {
        console.log(e);
        // noop...
      }
    });

    return groups;
  }
);

export const addProxies = createAsyncThunk(
  "proxies/add",
  async (payload: any) => {
    const { group, proxies } = payload;
    if (!group || !proxies?.length) {
      return;
    }

    const proxiesArr: Proxy[] = [];
    for (const proxy of proxies) {
      let p: Proxy;
      if (proxy.startsWith("http")) {
        const [host, port, username = null, password = null] = proxy.split('http://')[1].split(':')
        p = {
          id: '',
          active: false,
          testing: false,
          speed: null,
          raw: proxy,
          host,
          port,
          username,
          password
        }
      } else {
        const [host, port, username = null, password = null] = proxy.split(":");

        if (!username || !password) {
          p = {
            id: '',
            active: false,
            testing: false,
            speed: null,
            raw: `http://${host}:${port}`,
            host,
            port,
            username,
            password
          }
        } else {
          p = {
            id: '',
            active: false,
            testing: false,
            speed: null,
            raw: `http://${username}:${password}@${host}:${port}`,
            host,
            port,
            username,
            password
          }
        }
      }
      proxiesArr.push(p);
    }

    const ids = await window.RPCAction('proxies:add', [group.id, proxiesArr.map(p => p.raw)]);

    return {
      group,
      proxies: proxiesArr,
      ids
    };
  }
);

export const editProxies = createAsyncThunk(
  "proxies/edit",
  async (payload: any) => {
    const { group, id, proxy } = payload;

    const [host, port, username = null, password = null] = proxy.split(":");

    let p: any;
    if (!username || !password) {
      p = {
        id,
        active: false,
        testing: false,
        speed: null,
        raw: `http://${host}:${port}`,
        host,
        port,
        username,
        password
      }
    } else {
      p = {
        id,
        active: false,
        testing: false,
        speed: null,
        raw: `http://${username}:${password}@${host}:${port}`,
        host,
        port,
        username,
        password
      }
    }

    window.RPCAction('proxies:edit', [[{ id, url: p.raw }]]);

    return { group, proxy: p };
  }
);

export const removeProxies = createAsyncThunk(
  "proxies/remove",
  async ({ ids, group }: any) => {

    window.RPCAction('proxies:delete', [ids]);

    return { ids, group };
  }
);

export const copyToGroup = createAsyncThunk(
  "proxies/copyToGroup",
  async ({ oldGroup, newGroup, ids }: any, { getState }) => {
    const { Proxies } = getState() as RootState;
    const state = Proxies[oldGroup];

    let proxiesArr: Proxy[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      proxiesArr.push(original);
    }

    const newIds = await window.RPCAction('proxies:add', [newGroup, proxiesArr.map(p => p.raw)]);

    return { group: newGroup, proxies: proxiesArr, ids: newIds };
  }
);

export const moveToGroup = createAsyncThunk(
  "proxies/moveToGroup",
  async ({ oldGroup, newGroup, ids }: any) => {

    window.RPCAction('proxies:move', [ids, newGroup]);

    return { oldGroup, newGroup, ids };
  }
);

export const proxies = createSlice({
  name: "proxies",
  initialState: initialProxiesState,
  reducers: {},
  extraReducers: {
    // @ts-ignore
    [addProxiesGroup.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { id, groupName } = payload;

      notify(`Added ${groupName} group`, "success", { duration: 1500 });
      // eslint-disable-next-line no-param-reassign
      state[id] = {
        id: id,
        name: groupName,
        byId: {},
      };
    },
    // @ts-ignore
    [removeProxiesGroup.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload.length) {
        return;
      }

      payload.forEach((id: string) => {
        notify(`Removed ${state[id].name} group`, "success", { duration: 1500 });
        delete state[id];
      });
    },
    // @ts-ignore
    [addProxies.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { group, proxies, ids } = payload;
      proxies.forEach((proxy: Proxy, i: number) => {
        state[group.id].byId[ids[i]] = {
          ...proxy,
          id: ids[i]
        };
      });

      notify(`Added ${ids.length} ${ids.length > 1 ? "proxies" : "proxy"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [editProxies.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { group, proxy } = payload;
      state[group.id].byId[proxy.id] = proxy;

      notify(`Edited ${proxy.id} proxy`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [removeProxies.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { ids, group } = payload;
      ids.forEach((id: string) => {
        delete state[group].byId[id];
      });

      notify(`Removed ${ids.length} ${ids.length > 1 ? "proxies" : "proxy"}`, "success", { duration: 1500 });
    },
    //@ts-ignore
    [testProxies.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { ids, group } = payload;

      ids.forEach((id: string) => {
        state[group].byId[id].speed = -1;
        state[group].byId[id].testing = true;
      });
    },
    // @ts-ignore
    [copyToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, proxies, ids } = payload;

      // @ts-ignore
      Object.values(proxies).forEach((proxy: Proxy, i: number) => {
        state[group].byId[ids[i]] = {
          ...proxy,
          id: ids[i]
        };
      });

      notify(`Copied ${proxies.length} ${proxies.length > 1 ? "proxies" : "proxy"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [moveToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { oldGroup, ids, newGroup } = payload;
      ids.forEach((id: string) => {
        const old = state[oldGroup].byId[id];
        delete state[oldGroup].byId[id];
        state[newGroup].byId[id] = old;
      });

      notify(`Moved ${ids.length} ${ids.length > 1 ? "proxies" : "proxy"}`, "success", { duration: 1500 });
    },
  },
});

export const selectProxies = ({ Proxies }: RootState) => Proxies;

export const makeProxies = createSelector(
  selectProxies,
  (state: Proxies) => state
);

const { reducer } = proxies;

export default reducer;
