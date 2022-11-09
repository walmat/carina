import {
  createAsyncThunk,
  createSelector,
  createSlice,
} from "@reduxjs/toolkit";
import valid from "card-validator";

import {capitalize, notify, stripSpaces} from '../../../utils';

import { RootState } from '.';

export interface NameCodePair {
  name: string;
  code: string;
}

export interface Address {
  name: string;
  line1: string;
  line2: string;
  line3: string;
  postCode: string;
  city: string;
  country: NameCodePair;
  state: NameCodePair;
}

export interface Payment {
  email: string;
  phone: string;
  name: string;
  number: string;
  type: string;
  expMonth: string;
  expYear: string;
  cvv: string;
}

export interface Profile {
  id: string | null;
  name: string;
  shipping: Address;
  billing: Address;
  payment: Payment;
  maxCheckouts: number;
}

export interface ProfileGroup {
  id: string;
  name: string;
  byId: {
    [id: string]: Profile;
  };
}

export interface Profiles {
  [id: string]: ProfileGroup;
}

export const initialProfilesState: Profiles = {
  default: {
    id: "default",
    name: "Default",
    byId: {},
  },
};

const transformProfile = (values: Profile, matches: boolean) => {
  const profile: Profile = { ...values };

  profile.name = capitalize(profile.name);

  profile.shipping.name = capitalize(profile.shipping.name);
  profile.shipping.line1 = capitalize(profile.shipping.line1);
  profile.shipping.line2 = capitalize(profile.shipping.line2);
  profile.shipping.line3 = capitalize(profile.shipping.line3);
  profile.shipping.city = capitalize(profile.shipping.city);

  if (matches) {
    profile.billing = { ...profile.shipping };
  }

  profile.payment.name = capitalize(profile.payment.name);
  profile.payment.number = stripSpaces(profile.payment.number);
  profile.payment.email = stripSpaces(profile.payment.email);
  profile.payment.cvv = stripSpaces(profile.payment.cvv);

  const validator: any = valid.number(profile.payment.number);
  profile.payment.type = validator?.card?.type;

  // @ts-ignore
  const [expMonth, expYear] = profile.payment.exp.split("/");
  // @ts-ignore
  delete profile.payment.exp;

  profile.payment.expMonth = expMonth;
  profile.payment.expYear = `20${expYear}`;

  return profile;
};

const digestChanges = (values: Profile, matches: boolean) => {
  let changes: any = {};

  const profile: Profile = { ...values };

  if (profile.name.trim()) {
    changes.name = capitalize(profile.name.trim());
  }

  changes.shipping = {};
  if (profile.shipping.name.trim()) {
    changes.shipping.name = capitalize(profile.shipping.name.trim());
  }

  if (profile.shipping.line1.trim()) {
    changes.shipping.line1 = capitalize(profile.shipping.line1.trim());
  }

  if (profile.shipping.line2.trim()) {
    changes.shipping.line2 = capitalize(profile.shipping.line2.trim());
  }

  if (profile.shipping.line3.trim()) {
    changes.shipping.line3 = capitalize(profile.shipping.line3.trim());
  }

  if (profile.shipping.postCode.trim()) {
    changes.shipping.postCode = capitalize(profile.shipping.postCode.trim());
  }

  if (profile.shipping.city.trim()) {
    changes.shipping.city = capitalize(profile.shipping.city.trim());
  }

  if (profile.shipping.country) {
    changes.shipping.country = profile.shipping.country;
  }

  if (profile.shipping.state) {
    changes.shipping.state = profile.shipping.state;
  }

  if (matches) {
    changes.billing = { ...changes.shipping };
  } else {
    changes.billing = {};
    if (profile.billing.name.trim()) {
      changes.billing.name = capitalize(profile.billing.name.trim());
    }

    if (profile.billing.line1.trim()) {
      changes.billing.line1 = capitalize(profile.billing.line1.trim());
    }

    if (profile.billing.line2.trim()) {
      changes.billing.line2 = capitalize(profile.billing.line2.trim());
    }

    if (profile.billing.line3.trim()) {
      changes.billing.line3 = capitalize(profile.billing.line3.trim());
    }

    if (profile.billing.postCode.trim()) {
      changes.billing.postCode = capitalize(profile.billing.postCode.trim());
    }

    if (profile.billing.city.trim()) {
      changes.billing.city = capitalize(profile.billing.city.trim());
    }

    if (profile.billing.country) {
      changes.billing.country = profile.billing.country;
    }

    if (profile.billing.state) {
      changes.billing.state = profile.billing.state;
    }
  }

  changes.payment = {};
  if (profile.payment.name.trim()) {
    changes.payment.name = capitalize(profile.payment.name.trim());
  }

  if (profile.payment.number.trim()) {
    changes.payment.number = stripSpaces(profile.payment.number.trim());

    const validator: any = valid.number(changes.payment.number);
    changes.payment.type = validator?.card?.type;
  }

  if (profile.payment.email.trim()) {
    changes.payment.email = stripSpaces(profile.payment.email.trim());
  }

  if (profile.payment.cvv.trim()) {
    changes.payment.cvv = stripSpaces(profile.payment.cvv.trim());
  }

  // @ts-ignore
  if (profile.payment.exp) {
    // @ts-ignore
    const [expMonth, expYear] = profile.payment.exp.split("/");
    // @ts-ignore
    delete profile.payment.exp;

    changes.payment.expMonth = expMonth;
    changes.payment.expYear = `20${expYear}`;
  }

  return changes;
};

export const addProfileGroup = createAsyncThunk(
  "profiles/addGroup",
  async (groupName: string) => {
    if (!groupName.trim()) {
      return null;
    }

    const id = await window.RPCAction('profiles:addGroup', [groupName.trim()]);
    return { id, groupName: groupName.trim() };
  }
);

export const removeProfileGroup = createAsyncThunk(
  "profile/removeGroup",
  async (groups: string[]) => {
    if (!groups?.length) {
      return null;
    }

    groups.forEach((id) => {
      try {
        window.RPCAction('profiles:deleteGroup', [id]);
      } catch (e) {
        // noop...
      }
    });

    return groups;
  }
);

export const addProfiles = createAsyncThunk(
  "profiles/add",
  async (values: any) => {
    const { group, matches, ...rest } = values;

    const profile = transformProfile(rest, matches);
    const ids = await window.RPCAction('profiles:add', [group.id, [profile]]);

    return { ids, profile, group };
  }
);

export const copyProfiles = createAsyncThunk(
  "profiles/copy",
  async ({ ids, group }: any, { getState }) => {
    const { Profiles } = getState() as RootState;
    const state = Profiles[group];

    const profiles: any[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      profiles.push(original);
    }

    const newIds = await window.RPCAction('profiles:add', [group, profiles]);

    return { group, profiles, ids: newIds };
  }
);

export const editProfiles = createAsyncThunk(
  "profiles/edit",
  async (values: any, { getState }) => {
    const { Profiles } = getState() as RootState;

    const { ids, group, matches, ...rest } = values;
    const profiles = Profiles[group.id];

    const changes = digestChanges(rest, matches);
    let profs: any[] = [];
    for (const id of ids) {
      const p = profiles.byId[id];

      profs.push({
        ...p,
        ...changes,
        shipping: {
          ...p.shipping,
          ...changes.shipping,
        },
        billing: {
          ...p.billing,
          ...changes.billing,
        },
        payment: {
          ...p.payment,
          ...changes.payment,
        },
        id
      });
    }

    window.RPCAction('profiles:edit', [profs]);

    return { profiles: profs, group };
  }
);

export const removeProfiles = createAsyncThunk(
  "profiles/remove",
  async ({ ids, group }: ProfileProps) => {
    window.RPCAction('profiles:delete', [ids]);

    return { ids, group };
  }
);

export const copyToGroup = createAsyncThunk(
  "profiles/copyToGroup",
  async ({ oldGroup, newGroup, ids }: any, { getState }) => {
    const { Profiles } = getState() as RootState;
    const state = Profiles[oldGroup];

    const profs: any[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      profs.push(original);
    }

    const newIds = await window.RPCAction('profiles:add', [newGroup, profs]);

    return { group: newGroup, profiles: profs, ids: newIds };
  }
);

export const moveToGroup = createAsyncThunk(
  "profiles/moveToGroup",
  async ({ oldGroup, newGroup, ids }: any) => {
    window.RPCAction('profiles:move', [ids, newGroup]);

    return { oldGroup, newGroup, ids };
  }
);

interface ProfileProps {
  ids: string[];
  group: string;
}

export const profiles = createSlice({
  name: "profiles",
  initialState: initialProfilesState,
  reducers: {
    importProfiles: (state, action) => {},
    loadProfile: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }
      Object.keys(payload).forEach((profileGroup) => {
        state[profileGroup] = payload[profileGroup];
      });
    },
  },
  extraReducers: {
    // @ts-ignore
    [addProfileGroup.fulfilled]: (state, action) => {
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
    [removeProfileGroup.fulfilled]: (state, action) => {
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
    [addProfiles.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { profile, group, ids } = payload;
      ids.forEach((id: string) => {
        state[group.id].byId[id] = {
          ...profile,
          id
        }
      });

      notify(`Added ${ids.length} ${ids.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [removeProfiles.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { ids, group } = payload;

      ids.forEach((id: string) => {
        delete state[group].byId[id];
      });

      notify(`Removed ${ids.length} ${ids.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [editProfiles.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { profiles, group } = payload;
      profiles.forEach((profile: Profile) => {
        // @ts-ignore
        state[group.id].byId[profile.id] = profile;
      });

      notify(`Added ${profiles.length} ${profiles.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [copyProfiles.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, profiles, ids } = payload;
      profiles.forEach((profile: Profile, i: number) => {
        state[group].byId[ids[i]] = {
          ...profile,
          id: ids[i],
          name: `${profile.name} Copy`
        }
      });

      notify(`Copied ${ids.length} ${ids.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
    },
    // @ts-ignore
    [copyToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, profiles, ids } = payload;
      profiles.forEach((profile: Profile, i: number) => {
        state[group].byId[ids[i]] = {
          ...profile,
          id: ids[i]
        }
      });

      notify(`Copied ${profiles.length} ${profiles.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
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

      notify(`Moved ${ids.length} ${ids.length > 1 ? "profiles" : "profile"}`, "success", { duration: 1500 });
    },
  },
});

export const selectProfiles = ({ Profiles }: RootState) => Profiles;

export const makeProfiles = createSelector(
  selectProfiles,
  (state: RootState["Profiles"]) => state
);

export const makeProfileGroupOptions = createSelector(
  selectProfiles,
  (state: RootState["Profiles"]) => Object.values(state).map(({ byId, ...rest }) => ({ ...rest }))
);

const { actions, reducer } = profiles;

export const { importProfiles, loadProfile } = actions;

export default reducer;
