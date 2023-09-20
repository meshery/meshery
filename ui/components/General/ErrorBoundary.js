import { ErrorBoundary as ReactErrorBoundary } from "react-error-boundary";

function Fallback({ error }) {
  // Call resetErrorBoundary() to reset the error boundary and retry the render.

  return (

    <div role="alert" >
      <h2>Uh-oh!😔 Please pardon our mesh.</h2>
      <div
        style={{
          backgroundColor : "#1E2117",
          color : "#FFFFFF",
          padding : ".85rem",
          borderRadius : ".2rem"
        }}
      >
        <code>{error.message}</code>
      </div>
    </div>

  )
}

const reportError = (error,info) => {
  // This is where you'd send the error to Sentry,etc
  console.log("Error Caught Inside Boundary --reportError", error,"Info",info);
}

export const ErrorBoundary = ({ children, ...props  }) => {
  return (
    <ReactErrorBoundary FallbackComponent={Fallback} onError={reportError}  {...props}>
      {children}
    </ReactErrorBoundary>
  );
}

export const withErrorBoundary = (Component,errorHandlingProps={}) => {
  const WrappedWithErrorBoundary = (props) => (
    <ErrorBoundary  {...errorHandlingProps}>
      <Component {...props} />
    </ErrorBoundary>
  );

  return WrappedWithErrorBoundary;
}

export const withSuppressedErrorBoundary = (Component) => {
  const WrappedWithErrorBoundary = (props) => (
    <ErrorBoundary  FallbackComponent={() => null }>
      <Component {...props} />
    </ErrorBoundary>
  );

  return WrappedWithErrorBoundary;
}