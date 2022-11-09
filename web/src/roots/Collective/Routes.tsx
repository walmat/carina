import React, { useMemo } from "react";
import Loadable from 'react-loadable';
import { Switch, Route } from "react-router-dom";
import { AnimatePresence } from "framer-motion";

export const routes = [
	{
		path: "/",
		name: "Collective",
		component: Loadable({
			loader: () => import('../../containers/Collective'),
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
