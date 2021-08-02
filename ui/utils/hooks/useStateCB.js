import React from "react";

function useStateCB(initState) {
  const [state, _setState] = React.useState(initState);
  const stateRef = React.useRef(initState);

  const callbackRef = React.useRef();
  const isFirstCBCall = React.useRef(true);

  React.useEffect(() => {
    if (isFirstCBCall.current) isFirstCBCall.current = false;
    else callbackRef.current?.(state);
  }, [state]);

  const setState = (state, callback) => {
    callbackRef.current = callback;

    stateRef.current = state;
    _setState(state);
  };

  const getStateRefValue = () => stateRef.current

  return [state, setState, getStateRefValue];
}

export default useStateCB;
