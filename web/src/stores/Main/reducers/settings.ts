import { createSelector, createSlice } from "@reduxjs/toolkit";
import { RootState } from ".";

interface Behaviors {
  notifications: boolean;
  retryCheckouts: boolean;
  monitorPooling: boolean;
  autoClick: boolean;
  autoLaunch: boolean;
}

export interface Settings {
  language: {
    label: string;
    value: string;
  };
  currency: {
    label: string;
    value: string;
  };
  sounds: {
    checkout: {
      name: string;
      volume: number;
      default: boolean;
    };
    harvester: {
      name: string;
      volume: number;
      default: boolean;
    };
  };
  behaviors: Behaviors;
  isOpen: boolean;
}

export const initialSettingsState: Settings = {
  language: {
    label: "English (US)",
    value: "en-US",
  },
  currency: {
    label: "USD ($)",
    value: "usd",
  },
  sounds: {
    checkout: {
      name: "checkout.mp3",
      volume: 100,
      default: true,
    },
    harvester: {
      name: "harvester.mp3",
      volume: 100,
      default: true,
    },
  },
  behaviors: {
    notifications: true,
    retryCheckouts: false,
    monitorPooling: false,
    autoClick: true,
    autoLaunch: true,
  },
  isOpen: false
};

export const settings = createSlice({
  name: "settings",
  initialState: initialSettingsState,
  reducers: {
    editSound: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { type, name, ...rest } = payload;

      // @ts-ignore
      state.sounds[type] = { name, default: false, ...rest };
    },
    revertSound: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { type } = payload;

      // @ts-ignore
      state.sounds[type] = initialSettingsState.sounds[type];
    },
    changeVolume: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { type, volume } = payload;

      // @ts-ignore
      state.sounds[type].volume = volume;
    },
    editLanguage: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      state.language = payload;
    },
    editCurrency: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      state.currency = payload;
    },
    toggle: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { base, field } = payload;

      if (base) {
        // @ts-ignore
        state[base][field] = !state[base][field];
      } else {
        // @ts-ignore
        state[field] = !state[field];
      }
    },
    stash: (state) => {
      (async () => {
        // @ts-ignore
        window.RPCAction('preferences:set', [state]);
      })();
    },
  },
});

export const BEHAVIORS = {
  NOTIFICATIONS: "notifications",
  RETRY_CHECKOUTS: "retryCheckouts",
  MONITOR_POOLS: "monitorPooling",
  AUTO_CLICK: "autoClick",
  AUTO_LAUNCH: "autoLaunch",
};

export const SETTINGS = {
  IS_OPEN: 'isOpen'
};

const selectSettings = ({ Settings }: RootState) => Settings;

export const makeSettings = createSelector(
  selectSettings,
  (state: RootState["Settings"]) => state
);

export const makeSounds = createSelector(
  selectSettings,
  (state: RootState["Settings"]) => state.sounds
);
const { actions, reducer } = settings;

export const {
  toggle,
  stash,
  editSound,
  revertSound,
  changeVolume,
  editLanguage,
  editCurrency,
} = actions;

export default reducer;
