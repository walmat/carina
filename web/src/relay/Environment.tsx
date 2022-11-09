import { Environment, Network, RecordSource, Store } from "relay-runtime";
import fetchGraphQL from "./fetchGraphQL";
import { LogFunction } from "relay-runtime/lib/store/RelayStoreTypes";
import { FetchFunction } from "relay-runtime/lib/network/RelayNetworkTypes";

const isDev = process.env.NODE_ENV === "development";

// Relay passes a "params" object with the query name and text. So we define a helper function
// to call our fetchGraphQL utility with params.text.
const fetchRelay: FetchFunction = async (params, variables) => {
  console.log(
    `fetching query ${params.name} with ${JSON.stringify(variables)}`
  );
  return fetchGraphQL(params.text, variables);
};

export const relayTransactionLogger: LogFunction = (event) => {
  // eslint-disable-next-line
  console.log("RELAY: ", event);
  return;
};

// Export a singleton instance of Relay Environment configured with our network function:
const env = new Environment({
  network: Network.create(fetchRelay),
  store: new Store(new RecordSource(), {
    gcReleaseBufferSize: 10,
  }),
  log: isDev ? relayTransactionLogger : null,
});

// TODO - enable this
// if (isDev) {
//   if (typeof window !== 'undefined') {
//     window.relayEnvironment = env;
//     window.debugRelayStore = () => env.getStore().getSource().toJSON();
//   }
// }

export default env;
