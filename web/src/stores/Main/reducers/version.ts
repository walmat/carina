import {
	createSlice,
	createSelector,
} from "@reduxjs/toolkit";

import { RootState } from ".";

export const initialVersionState = 'unknown';
export const theme = createSlice({
	name: "version",
	// @ts-ignore
	initialState: initialVersionState,
	reducers: {}
});

export const selectVersion = ({ Version }: RootState) => Version;

export const makeVersion = createSelector(
	selectVersion,
	(state: RootState["Version"]) => state
);

const { reducer } = theme;

export default reducer;
