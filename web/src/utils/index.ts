import { ProfileUtils } from "./profiles";
import { notify } from './notify';

export const minimize = () => {
  astilectron.sendMessage({ type: 'minimize' })
};

export const hide = () => {
  astilectron.sendMessage({ type: 'hide' })
}

export const close = () => {
  astilectron.sendMessage({ type: 'close' })
};
export const relaunch = () => {};

export const noop = () => {};

export const capitalize = (input: string) =>
  input
    .replace(
      /\w\S*/g,
      (txt) => txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase()
    )
    .trim();

export const stripSpaces = (input: string) => input.replace(/\s+/g, "");

export { ProfileUtils, notify };
