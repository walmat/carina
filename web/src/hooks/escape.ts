const { useEffect } = require("react");

const ESCAPE_KEY = 27;

type Callback = (_: any) => void;
type Dependencies = any[];

export const useEscape = (
  callback: Callback,
  dependencies: Dependencies = []
) => {
  useEffect(() => {
    if (!window || !window.document || !callback) {
      return;
    }

    if (!Array.isArray(dependencies)) {
      console.warn("Dependencies must be an array!");
    }

    const onKeyPress = (event: any) =>
      event.keyCode === ESCAPE_KEY && callback(event);
    window.document.addEventListener("keydown", onKeyPress);
    return () => {
      window.document.removeEventListener("keydown", onKeyPress);
    };
  }, [callback, ...dependencies]);
};
