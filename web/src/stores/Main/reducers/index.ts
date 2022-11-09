import { combineReducers } from "@reduxjs/toolkit";

import User from "./user";
import Stores from "./stores";
import Theme from "./theme";
import Version from "./version";

// main reducers
import Tasks from "./tasks";
import Workflows from "./workflows";
import Proxies from "./proxies";
import Profiles from "./profiles";

// settings reducers
import Settings from "./settings";
import Defaults from "./defaults";
import Accounts from "./accounts";
import Rates from "./rates";
import Webhooks from "./webhooks";
import Integrations from "./integrations";

// delays reducers
import Delays from "./delays";

import Update from "./update";

export type RootState = {
  Theme: ReturnType<typeof Theme>;
  Version: ReturnType<typeof Version>;
  User: ReturnType<typeof User>;
  Stores: ReturnType<typeof Stores>;
  Tasks: ReturnType<typeof Tasks>;
  Workflows: ReturnType<typeof Workflows>;
  Proxies: ReturnType<typeof Proxies>;
  Profiles: ReturnType<typeof Profiles>;
  Settings: ReturnType<typeof Settings>;
  Defaults: ReturnType<typeof Defaults>;
  Accounts: ReturnType<typeof Accounts>;
  Rates: ReturnType<typeof Rates>;
  Webhooks: ReturnType<typeof Webhooks>;
  Integrations: ReturnType<typeof Integrations>;
  Delays: ReturnType<typeof Delays>;
  Update: boolean;
};

const combinedReducer = combineReducers({
  Theme,
  Version,
  User,
  Stores,
  Tasks,
  Workflows,
  Proxies,
  Profiles,
  Settings,
  Defaults,
  Accounts,
  Rates,
  Webhooks,
  Integrations,
  Delays,
  Update,
});

const rootReducer = (state: any, action: any) => {
  if (action.type === "user/logout") {
    // eslint-disable-next-line no-param-reassign
    state = undefined;
  }

  if (action.type === "global/import") {
    const { payload } = action;

    Object.entries(payload).forEach(([key, values]) => {
      state[key] = values;
    });
  }

  return combinedReducer(state, action);
};

export const Reducers = {
  Theme,
  Version,
  User,
  Stores,
  Tasks,
  Workflows,
  Proxies,
  Profiles,
  Settings,
  Defaults,
  Accounts,
  Rates,
  Webhooks,
  Integrations,
  Delays,
  Update,
};

export default rootReducer;
