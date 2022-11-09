import { createSlice, createSelector } from "@reduxjs/toolkit";

import { RootState } from ".";

export const initialUpdateState = false;
export const theme = createSlice({
  name: "update",
  // @ts-ignore
  initialState: initialUpdateState,
  reducers: {
    wsUpdate: (state, action) => {
      state = true;
      return state;
    },
  },
});

export const selectUpdate = ({ Update }: RootState) => Update;

export const makeUpdate = createSelector(
  selectUpdate,
  (state: RootState["Update"]) => state
);

const { actions, reducer } = theme;

export const { wsUpdate } = actions;

export default reducer;
