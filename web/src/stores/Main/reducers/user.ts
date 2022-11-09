import { createSelector, createSlice } from "@reduxjs/toolkit";

import { RootState } from ".";

export type User = {
  id: string;
  hash: string;
  avatar: string;
  email: string;
  type: "F&F" | "Staff" | "Lifetime" | "Renewal" | string;
  createdAt: string;
  instances: number;
  maxInstances: number;
};

export const initialUserState: User = {
  id: '',
  hash: '',
  avatar: '',
  email: '',
  type: '',
  createdAt: '',
  instances: 0,
  maxInstances: 0
};

export const user = createSlice({
  name: "user",
  initialState: initialUserState,
  reducers: {
    logout: () => {},
    setUser: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      return { ...state, ...payload };
    },
  },
});

export const selectUser = ({ User }: RootState) => User;

export const makeUser = createSelector(
  selectUser,
  (state: RootState["User"]) => state
);

const { actions, reducer } = user;

export const { logout, setUser } = actions;

export default reducer;
