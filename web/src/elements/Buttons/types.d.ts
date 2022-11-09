import { ReactElement } from "react";

interface ButtonProps {
  variant?: "Icon" | "Button" | "IconButton";
  command?: string;
  text?: string | ReactElement;
  loading?: boolean;
  disabled?: boolean;
  width: number | string;
  height: number | string;
  onClick: any;
  children?: any;
}
