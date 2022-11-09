import { StrictMode } from "react";
import { render } from "react-dom";

import { Main } from "../../stores";
import Roots from "../../roots";
import { attachRpc } from "../../utils/rpc";
import { GenericErrorBoundary } from "../../error/ErrorFallback";

const MOUNT_POINT = document.getElementById("root");

export const mount = async () => {
  const { Harvesters, ...initialState } = await window.RPCAction("state:get", [
    "main",
  ]);
  const { store } = Main.config(initialState);

  return render(
    <StrictMode>
      <GenericErrorBoundary>
        <Roots.Main store={store} />
      </GenericErrorBoundary>
    </StrictMode>,
    MOUNT_POINT
  );
};

attachRpc().then(() => {
  mount()
    .then(() => {})
    .catch(() => {});
});
