import { createSlice, createSelector, createAsyncThunk } from '@reduxjs/toolkit';
import { RootState } from './index';

export const setTheme = createAsyncThunk(
    'theme/set',
    async ({ theme }: any) => {
      let newTheme = (theme + 1) % 2;

      await window.RPCAction('theme:set', [newTheme]);

      return { theme: newTheme }
    }
)

export const initialThemeState = 0;
export const theme = createSlice({
  name: 'theme',
  // @ts-ignore
  initialState: initialThemeState,
  reducers: {
    resetTheme: () => initialThemeState,
  },
  extraReducers: {
    // @ts-ignore
    [setTheme.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { theme } = payload;

      return theme;
    }
  }
});

export const selectTheme = ({ Theme }: RootState) => Theme;

export const makeTheme = createSelector(
  selectTheme,
  (state: RootState['Theme']) => state
);

const { actions, reducer } = theme;

export const { resetTheme } = actions;

export default reducer;
