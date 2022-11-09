import { useCallback, useEffect } from "react";

export function useKeyPress(
  targetKey: string,
  callback: (_?: any) => void,
  deps: any[] = []
) {
  // Add event listeners

  const handler = ({ key, code }: KeyboardEvent) => {
    if (key !== targetKey && code !== targetKey) {
      return;
    }

    callback();
  };

  useEffect(() => {
    window.addEventListener("keypress", handler);
    // Remove event listeners on cleanup
    return () => {
      window.removeEventListener("keypress", handler);
    };
  }, [...deps, callback]); // Empty array ensures that effect is only run on mount and unmount
}
