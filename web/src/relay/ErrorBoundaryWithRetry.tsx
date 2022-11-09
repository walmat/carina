import { Component, ReactNode } from "react";

type RetryFn = () => void;

type FetchKeyProp = {
  fetchKey: number;
};

type ChildrenFunction = (fetchKeyProps: FetchKeyProp) => ReactNode;

type Props = {
  children: ReactNode | ChildrenFunction;
  fallback?: (error: string | null, retry: RetryFn) => ReactNode | ReactNode;
};
type State = {
  error: string | null;
  fetchKey: number;
};
class ErrorBoundaryWithRetry extends Component<Props, State> {
  state = { error: null, fetchKey: 0 };

  static getDerivedStateFromError(error: string | null) {
    return { error: error };
  }

  _retry = () => {
    this.setState((prev) => ({ error: null, fetchKey: prev.fetchKey + 1 }));
  };

  render() {
    const { children, fallback } = this.props;
    const { error, fetchKey } = this.state;

    if (error) {
      if (typeof fallback === "function") {
        return fallback(error, this._retry);
      }
      return fallback;
    }
    if (typeof children === "function") {
      // @ts-ignore
      return children({ fetchKey });
    }
    return children;
  }
}

export default ErrorBoundaryWithRetry;

export const ErrorViewBoundary = ({
  children,
}: {
  children: ChildrenFunction;
}) => {
  return (
    <ErrorBoundaryWithRetry
      fallback={(error) => <span>{JSON.stringify(error)}</span>}
    >
      {({ fetchKey }) => children({ fetchKey })}
    </ErrorBoundaryWithRetry>
  );
};
