import * as CSS from "csstype";

export interface Theme {
  background: CSS.Property.Color;
  backdrop: CSS.Property.Color;
  linear: CSS.Property.Color;
  fade: CSS.Property.Color;
  tooltip: CSS.Property.Color;
  lightRow: CSS.Property.Color;
  viewing: CSS.Property.Color;
  arrow: CSS.Property.Color;
  sidebar: CSS.Property.Color;
  primary: CSS.Property.Color;
  secondary: CSS.Property.Color;
  tertiary: CSS.Property.Color;
  inactive: CSS.Property.Color;
  h1: CSS.Property.Color;
  h2: CSS.Property.Color;
  h3: CSS.Property.Color;
  paragraph: CSS.Property.Color;
  placeholder: CSS.Property.Color;
  link: CSS.Property.Color;
  table: CSS.Property.Color;
  disabled: CSS.Property.Color;
  border: CSS.Property.Color;
  edit: CSS.Property.Color;
  separator: CSS.Property.Color;
  altLogin: CSS.Property.Color;
  accent: CSS.Property.Color;
  view: CSS.Property.Color;
  muted: CSS.Property.Color;
  expand: CSS.Property.Color;
  subHeading: CSS.Property.Color;
  description: CSS.Property.Color;
  tertiaryInner: CSS.Property.Color;
  tertiaryHue: CSS.Property.Color;
  lightHue: CSS.Property.Color;
  scrollbar: CSS.Property.Color;
  unscrollbar: CSS.Property.Color;
  day: CSS.Property.Color;
  dayText: CSS.Property.Color;
  dayActive: CSS.Property.Color;
  success: CSS.Property.Color;
  warning: CSS.Property.Color;
  failed: CSS.Property.Color;
  stop: CSS.Property.Color;
}

export const basePalette = {
  success: "#87B38D",
  warning: "#FFC857",
  failed: "#F26E86",
  stop: "#FFB15E",
};

export const lightPalette: Theme = {
  background: "#F4F4F4",
  backdrop: "rgba(244, 244, 244, 0.6)",
  linear: "linear-gradient(180deg, rgba(255, 255, 255, 0) 0%, #FFFFFF 100%)",
  fade: "rgba(255, 255, 255",
  tooltip: "rgba(0, 0, 0, 0.4)",
  lightRow: "#EFEDFE",
  viewing: "#786EF2",
  arrow: "#fff",
  sidebar: "#FFFFFF",
  primary: "#786EF2",
  secondary: "#EFEDFE",
  tertiary: "#616161",
  inactive: "#786EF2",
  h1: "#000000",
  h2: "#202126",
  h3: "#979797",
  paragraph: "#616161",
  placeholder: "hsl(0,0%,70%)",
  link: "#786EF2",
  table: "#202126",
  disabled: "rgba(97, 97, 97, 0.4)",
  border: "#979797",
  edit: "#d8d8d8",
  separator: "#F0F0F0",
  altLogin: "#C4C4C4",
  accent: "#F26E6E",
  view: "#d8d8d8",
  muted: "#616161",
  expand: "#4c4c4c",
  subHeading: "#202126",
  description: "#979797",
  tertiaryHue: "rgba(0, 0, 0, 0.35)",
  tertiaryInner: "rgba(0, 0, 0, 0.15)",
  lightHue: "rgba(0, 0, 0",
  scrollbar: "#786EF2",
  unscrollbar: "#d8d8d8",
  day: "rgba(240, 240, 240, 1)",
  dayText: "#786EF2",
  dayActive: "#F2F1FE",
  ...basePalette,
};

export const darkPalette: Theme = {
  background: "#202126",
  backdrop: "rgba(32, 33, 38, 0.6)",
  linear: "linear-gradient(180deg, rgba(0, 0, 0, 0) 0%, #2E2F34 100%)",
  fade: "rgba(46, 47, 52",
  tooltip: "rgba(255, 255, 255, 0.4)",
  lightRow: "rgba(120, 110, 242, 0.2)",
  viewing: "#fff",
  arrow: "#fff",
  sidebar: "#2E2F34",
  primary: "#786EF2",
  secondary: "rgba(120, 110, 242, 0.2)",
  tertiary: "#616161",
  inactive: "#FFFFFF",
  h1: "#ffffff",
  h2: "#f4f4f4",
  h3: "#979797",
  paragraph: "#979797",
  placeholder: "hsl(0,0%,40%)",
  link: "#786EF2",
  table: "#d8d8d8",
  disabled: "#d8d8d8",
  border: "#616161",
  edit: "#616161",
  separator: "#484848",
  altLogin: "#414247",
  accent: "#F26E6E",
  view: "#202126",
  muted: "#FFFFFF",
  expand: "#d8d8d8",
  subHeading: "#e6e6e6",
  description: "#f0f0f0",
  scrollbar: "#786EF2",
  tertiaryHue: "rgba(255, 255, 255, 0.15)",
  tertiaryInner: "rgba(255, 255, 255, 0.25)",
  lightHue: "rgba(255, 255, 255",
  unscrollbar: "rgba(120, 110, 242, 0.2)",
  day: "rgba(240, 240, 240, 0.2)",
  dayText: "#fff",
  dayActive: "#464184",
  ...basePalette,
};
