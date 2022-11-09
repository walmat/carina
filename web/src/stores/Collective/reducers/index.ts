import { combineReducers } from '@reduxjs/toolkit';

import Theme from './theme';
import Harvesters from './harvester';
import Accounts from './accounts';

export type RootState = {
  Theme: ReturnType<typeof Theme>;
  Accounts: ReturnType<typeof Accounts>;
  Harvesters: ReturnType<typeof Harvesters>;
};

const combinedReducer = combineReducers({
  Theme,
  Accounts,
  Harvesters
});

const rootReducer = (state: any, action: any) => {
  return combinedReducer(state, action);
};

export const Reducers = {
  Theme,
  Accounts,
  Harvesters
}

export default rootReducer;
