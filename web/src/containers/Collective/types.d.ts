import { Store } from "../../stores/Main/reducers/stores";

export interface Account {
  id: string;
  name: string;
}

export type Harvester = {
  id: string;
  index: number; // used for sorting the solvers
  focused: boolean; // is it the focused tab
  name: string; // name of solver
  store?: Store;
  account?: Account;
  proxy: string;
};

export interface Harvesters {
  byId: {
    [id: string]: Harvester;
  };
}
