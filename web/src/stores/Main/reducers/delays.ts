import { createSlice } from "@reduxjs/toolkit";
import { RootState } from ".";

type DelayObject = {
  monitor: number;
  retry: number;
};

export interface Delays {
  [id: string]: DelayObject;
}

export const initialDelaysState: Delays = {};

export const delays = createSlice({
  name: "delays",
  initialState: initialDelaysState,
  reducers: {}, // TODO;
});

export const selectDelays = ({ Delays }: RootState) => Delays;

const { reducer } = delays;

// export const {  } = actions;

export default reducer;
