import { Suspense, ComponentType, ReactNode } from "react";
import { ErrorViewBoundary } from "./ErrorBoundaryWithRetry";

type Config = {
  loadingView?: NonNullable<ReactNode> | null;
};
export const withRelayBoundary = (
  Component: ComponentType<any>,
  config: Config = {}
): ComponentType<any> => {
  const { loadingView = <span>Loading</span> } = config;

  const Wrapper = (props: any) => {
    return (
      <ErrorViewBoundary>
        {({ fetchKey }) => (
          <Suspense fallback={loadingView}>
            <Component {...props} fetchKey={fetchKey} />
          </Suspense>
        )}
      </ErrorViewBoundary>
    );
  };
  Wrapper.displayName = `withRelayBoundary(${
    Component.displayName || Component.name
  })`;

  return Wrapper;
};
