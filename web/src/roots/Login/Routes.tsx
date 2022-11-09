import React, { useMemo } from "react";
import Loadable from "react-loadable";
import { Switch, Route } from "react-router-dom";
import { AnimatePresence } from "framer-motion";

export const routes = [
  {
    path: "/",
    name: "Login",
    component: Loadable({
      loader: () => import("../../containers/Login"),
      loading: () => null,
    }),
  },
  {
    path: "/register",
    name: "Register",
    component: Loadable({
      loader: () => import("../../containers/Register"),
      loading: () => null,
    }),
  },
  {
    path: "/forgot",
    name: "Forgot Password",
    component: Loadable({
      loader: () => import("../../containers/ForgotPassword"),
      loading: () => null,
    }),
  },
  {
    path: "/change",
    name: "Change Password",
    component: Loadable({
      loader: () => import("../../containers/ChangePassword"),
      loading: () => null,
    }),
  },
  {
    path: "/2fa",
    name: "2-Factor Authentication",
    component: Loadable({
      loader: () => import("../../containers/2FA"),
      loading: () => null,
    }),
  },
];

export const isActive = (route: string, pathname: string) => pathname === route;

const Routes = () => {
  return useMemo(
    () => (
      <Route
        render={({ location }) => (
          // @ts-ignore
          <AnimatePresence key={location} exitBeforeEnter initial={false}>
            <Switch location={location} key={location.pathname}>
              {routes.map(({ path, component }) => (
                <Route key={path} exact path={path} component={component} />
              ))}
            </Switch>
          </AnimatePresence>
        )}
      />
    ),
    []
  );
};

export default Routes;
