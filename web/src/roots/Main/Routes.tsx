import React, { useMemo } from "react";
import Loadable from 'react-loadable';
import { Switch, Route } from "react-router-dom";
import { AnimatePresence } from "framer-motion";

export const top = [
  {
    path: "/",
    name: "Dashboard",
    component: Loadable({
      loader: () => import('../../containers/Dashboard'),
      loading: () => null,
    }),
  },
];

export const middle = [
  // {
  //   path: "/calendar",
  //   name: "Calendar",
  //   component: Loadable({
  //     loader: () => import('../../containers/Calendar'),
  //     loading: () => null,
  //   }),
  // },
  {
    path: "/tasks",
    name: "Tasks",
    component: Loadable({
      loader: () => import('../../containers/Tasks'),
      loading: () => null,
    }),
  },
  // {
  //   path: "/workflows",
  //   name: "Workflows",
  //   component: Loadable({
  //     loader: () => import('../../containers/Workflows'),
  //     loading: () => null,
  //   }),
  // },
];

export const bottom = [
  {
    path: "/profiles",
    name: "Profiles",
    component: Loadable({
      loader: () => import('../../containers/Profiles'),
      loading: () => null,
    }),
  },
  {
    path: "/proxies",
    name: "Proxies",
    component: Loadable({
      loader: () => import('../../containers/Proxies'),
      loading: () => null,
    }),
  },
];

export const settings = [
  {
    path: "/general",
    name: "General",
    component: Loadable({
      loader: () => import('../../containers/Settings/General'),
      loading: () => null,
    }),
  },
  {
    path: "/defaults",
    name: "Defaults",
    component: Loadable({
      loader: () => import('../../containers/Settings/Defaults'),
      loading: () => null,
    }),
  },
  {
    path: "/accounts",
    name: "Accounts",
    component: Loadable({
      loader: () => import('../../containers/Settings/Accounts'),
      loading: () => null,
    }),
  },
  {
    path: "/shipping",
    name: "Shipping Rates",
    component: Loadable({
      loader: () => import('../../containers/Settings/Rates'),
      loading: () => null,
    }),
  },
  {
    path: "/webhooks",
    name: "Webhooks",
    component: Loadable({
      loader: () => import('../../containers/Settings/Webhooks'),
      loading: () => null,
    }),
  },
  {
    path: "/integrations",
    name: "Integrations",
    component: Loadable({
      loader: () => import('../../containers/Settings/Integrations'),
      loading: () => null,
    }),
  },
];

export const isActive = (route: string, pathname: string) => pathname === route;

const sections = [...top, ...middle, ...bottom, ...settings];

const Routes = () => {
  return useMemo(
    () => (
      <Route
        render={({ location }) => (
          // @ts-ignore
          <AnimatePresence key={location} exitBeforeEnter initial={false}>
            <Switch location={location} key={location.pathname}>
              {sections.map(({ path, component }) => (
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
