import { ErrorBoundary as ReactErrorBoundary } from "react-error-boundary";

function Fallback({ error }) {
  // Call resetErrorBoundary() to reset the error boundary and retry the render.

  return (
    <div role="alert">
      <p>Something went wrong:</p>
      <pre style={{ color : "red" }}>{error.message}</pre>
    </div>
  )
}

export const ErrorBoundary = ({ children ,...props }) => {
  return (
    <ReactErrorBoundary FallbackComponent={Fallback} {...props}>
      {children}
    </ReactErrorBoundary>
  );
}

export const withErrorBoundary = (Component) => {
  const WrappedWithErrorBoundary = (props) => (
    <ErrorBoundary FallbackComponent>
      <Component {...props} />
    </ErrorBoundary>
  );

  return WrappedWithErrorBoundary;
}
