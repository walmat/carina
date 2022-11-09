import { StrictMode } from "react";
import { render } from "react-dom";

import { Collective } from "../../stores";
import Roots from "../../roots";
import { attachRpc } from "../../utils/rpc";
import { GenericErrorBoundary } from "../../error/ErrorFallback";
import { notify } from "../../utils";

const MOUNT_POINT = document.getElementById("root");

export const mount = async () => {
  const { Theme, Harvesters, Gmails } = await window.RPCAction("state:get", [
    "collective",
  ]);
  const { store } = Collective.config({ Theme, Harvesters, Accounts: Gmails });

  console.log("Store", store.getState());

  return render(
    <StrictMode>
      <GenericErrorBoundary>
        <Roots.Harvester store={store} />
      </GenericErrorBoundary>
    </StrictMode>,
    MOUNT_POINT
  );
};

attachRpc().then(() => {
  mount()
    .then(() => { })
    .catch(() => { });
});
