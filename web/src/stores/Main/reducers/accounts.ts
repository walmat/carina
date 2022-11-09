import {
  createSlice,
  createSelector,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { RootState } from ".";

import { Store } from "./stores";
import { notify } from "../../../utils";

export interface Account {
  id: string;
  store: Store;
  username: string;
  password: string;
}

export interface AccountGroup {
  id: string;
  name: string;
  byId: {
    [id: string]: Account;
  };
}

export interface Accounts {
  [id: string]: AccountGroup;
}

export const initialAccountsState: Accounts = {
  default: {
    id: "default",
    name: "Default",
    byId: {},
  },
};

export const createAccountsGroup = createAsyncThunk(
  "accounts/createGroup",
  async (groupName: string) => {
    if (!groupName.trim()) {
      return null;
    }

    const id = await window.RPCAction('accounts:addGroup', [groupName.trim()]);
    return { id, groupName: groupName.trim() };
  }
);

export const removeAccountsGroup = createAsyncThunk(
  "accounts/removeGroup",
  async (groups: string[]) => {
    if (!groups?.length) {
      return null;
    }

    groups.forEach((id) => {
      try {
        window.RPCAction('accounts:deleteGroup', [id]);
      } catch (e) {
        console.log(e);
        // noop...
      }
    });

    return groups;
  }
);

export const addAccounts = createAsyncThunk(
  "accounts/add",
  async ({ group, store, accounts }: any) => {
    const accs = [];
    for (const account of accounts) {
      const [username, password, ...rest] = account.split(":");
      if (!username || !password || (rest || []).length) {
        return;
      }

      accs.push({ store, username, password });
    }

    const ids = await window.RPCAction('accounts:add', [group.id, accs]);

    return { group, accounts: accs, ids };
  }
);

export const copyAccount = createAsyncThunk(
  "accounts/copy",
  async ({ ids, group }: any, { getState }) => {
    const { Accounts } = getState() as RootState;
    const state = Accounts[group];

    const accounts: any[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      accounts.push(original);
    }

    const newIds = await window.RPCAction('accounts:add', [group, accounts.map(a => ({ ...a, store: a.store.name }))])

    return { group, accounts, ids: newIds };
  }
);

export const removeAccount = createAsyncThunk(
  "accounts/remove",
  async ({ ids, group }: any) => {

    window.RPCAction('accounts:delete', [ids]);

    return { ids, group };
  }
);

export const editAccounts = createAsyncThunk(
  "accounts/edit",
  async ({ group, id, store, account }: any, { getState }: any) => {
    const { Accounts } = getState() as RootState;
    const state = Accounts[group.id];

    const [username, password, ...rest] = account.split(":");
    if (!username || !password || (rest || []).length) {
      return;
    }

    const old = state.byId[id];
    const base: any = { ...old };
    if (store) {
      base.store = store;
    }

    if (username) {
      base.username = username;
    }

    if (password) {
      base.password = password;
    }

    window.RPCAction('accounts:edit', [{ ...base }]);

    return { group, account: base };
  }
);

export const copyToGroup = createAsyncThunk(
  "accounts/copyToGroup",
  async ({ oldGroup, newGroup, ids }: any, { getState }) => {
    const { Accounts } = getState() as RootState;
    const state = Accounts[oldGroup];

    let accs: any[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      accs.push(original);
    }

    const newIds = await window.RPCAction('accounts:add', [newGroup, accs.map(a => ({ ...a, store: a.store.name }))])

    return { group: newGroup, accounts: accs, ids: newIds };
  }
);

export const moveToGroup = createAsyncThunk(
  "accounts/moveToGroup",
  async ({ oldGroup, newGroup, ids }: any) => {

    window.RPCAction('accounts:move', [ids, newGroup]);

    return { oldGroup, newGroup, ids };
  }
);

export const accounts = createSlice({
  name: "accounts",
  initialState: initialAccountsState,
  reducers: {},
  extraReducers: {
    // @ts-ignore
    [createAccountsGroup.fulfilled]: (state, action) => {
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
    [removeAccountsGroup.fulfilled]: (state, action) => {
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
    [addAccounts.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, accounts, ids } = payload;
      accounts.forEach((account: Account, i: number) => {
        state[group.id].byId[ids[i]] = {
          ...account,
          id: ids[i]
        };
      });

      notify(`Added ${ids.length} ${ids.length > 1 ? "accounts" : "account"}`, "success", { duration: 1500 });

    },
    // @ts-ignore
    [removeAccount.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { ids, group } = payload;

      ids.forEach((id: string) => {
        delete state[group].byId[id];
      });

      notify(`Removed ${ids.length} ${ids.length > 1 ? "accounts" : "account"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [editAccounts.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { group, account } = payload;

      state[group.id].byId[account.id] = account;

      notify(`Edited ${account.id} account`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [copyAccount.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, accounts, ids } = payload;
      accounts.forEach((account: Account, i: number) => {
        state[group].byId[ids[i]] = {
          ...account,
          id: ids[i],
        }
      });

      notify(`Copied ${ids.length} ${ids.length > 1 ? "accounts" : "account"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [copyToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, accounts, ids } = payload;
      accounts.forEach((account: Account, i: number) => {
        state[group].byId[ids[i]] = {
          ...account,
          id: ids[i]
        }
      });

      notify(`Copied ${accounts.length} ${accounts.length > 1 ? "accounts" : "account"}`, "success", { duration: 1500 });
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

      notify(`Moved ${ids.length} ${ids.length > 1 ? "accounts" : "account"}`, "success", { duration: 1500 });
    },
  },
});

export const selectAccounts = ({ Accounts }: RootState) => Accounts;

export const makeAccounts = createSelector(
  selectAccounts,
  (state: Accounts) => state
);

const { reducer } = accounts;

export default reducer;
