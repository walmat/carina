import { combineReducers } from '@reduxjs/toolkit';

import Theme from './theme';

export type RootState = {
  Theme: ReturnType<typeof Theme>;
};

const combinedReducer = combineReducers({
  Theme
});

const rootReducer = (state: any, action: any) => {
  return combinedReducer(state, action);
};

export const Reducers = {
  Theme
}

export default rootReducer;
