import { StrictMode } from "react";
import { render } from "react-dom";

import { Login } from "../../stores";
import Roots from "../../roots";

import { attachRpc } from "../../utils/rpc";
import { GenericErrorBoundary } from '../../error/ErrorFallback';

const MOUNT_POINT = document.getElementById("root");

export const mount = async () => {
  const { Theme } = await window.RPCAction('state:get', ['auth']);
  const { store } = Login.config({ Theme });

  return render(
    <StrictMode>
      <GenericErrorBoundary>
        <Roots.Login store={store} />
      </GenericErrorBoundary>
    </StrictMode>,
    MOUNT_POINT
  );
}

attachRpc().then(() => {
  mount().then(() => {}).catch(() => {})
});
