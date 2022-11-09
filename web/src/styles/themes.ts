import { Theme } from "styled-system";
import { darkPalette, lightPalette } from "./colors";
import { space } from "./spacing";
import { breakpoints } from "./breakpoints";

export const light: Theme = {
  space: {
    ...space,
  },
  breakpoints,
  colors: {
    ...lightPalette,
  },
};

export const dark: Theme = {
  space: {
    ...space,
  },
  breakpoints,
  colors: {
    ...darkPalette,
  },
};

export const themes = [light, dark];
