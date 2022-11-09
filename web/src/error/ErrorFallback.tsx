import { ErrorBoundary, FallbackProps } from 'react-error-boundary';
import { ReactNode } from 'react';

export const ErrorFallback = ({error, resetErrorBoundary}: FallbackProps) => {
  return (
    <div role="alert">
      <p>Something went wrong:</p>
      <pre>{error.message}</pre>
      <button onClick={resetErrorBoundary}>Try again</button>
    </div>
  )
}

type GenericErrorBoundaryProps = {
  children: ReactNode,
}
export const GenericErrorBoundary = ({ children }: GenericErrorBoundaryProps) => {
  return (
    <ErrorBoundary FallbackComponent={ErrorFallback}>
      {children}
    </ErrorBoundary>
  );
};
