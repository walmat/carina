import { configureStore } from "@reduxjs/toolkit";

import { createLogger } from "redux-logger";
import rootReducer from "../reducers";

const getMiddleware = (getDefaultMiddleware: any) => {
  // TODO: Disable logger when packaged

  const logger = createLogger({
    level: "info",
    collapsed: true,
  });

  return getDefaultMiddleware().concat(logger);
};

const config = (initialState = {}) => {
  const store = configureStore({
    preloadedState: initialState,
    reducer: rootReducer,
    middleware: (getDefaultMiddleware) => getMiddleware(getDefaultMiddleware),
    devTools: true, // TODO: Disable when packaged
  });

  return { store };
};

export { config };
