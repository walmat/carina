export interface Action {
  type: "Primary" | "First" | "Second" | "Action" | "Button";
  popover?: boolean;
  title: string;
  Icon: any;
  onClick?: any;
}
