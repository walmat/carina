import Account from "./accounts";
import { Create, Remove } from "./Group";
import Tasks, { EditTaskForm, CreateTaskForm } from "./tasks";
import Profiles from "./profiles";
import Proxies from "./proxies";
import { Rates } from "./rates";
import { Webhooks } from "./webhooks";
import { Shipment } from "./shipment";
import { Defaults } from "./defaults";
import { Login } from "./login";
import { Register } from "./register";
import { ForgotPassword } from "./forgot";
import { TwoFactor } from "./2fa";
import { ChangePassword } from "./change";

export {
  Account,
  Create,
  Remove,
  Tasks,
  Profiles,
  Proxies,
  Rates,
  Webhooks,
  Shipment,
  Defaults,
  Login,
  Register,
  ForgotPassword,
  ChangePassword,
  TwoFactor,
};

export type { EditTaskForm, CreateTaskForm };
