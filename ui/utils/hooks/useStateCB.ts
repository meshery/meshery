import React from 'react';

type StateUpdater<S> = (_state: S) => void;

function useStateCB<S>(
  initState: S,
  changeTrackCB?: StateUpdater<S>,
): [S, (_state: S, _callback?: StateUpdater<S>) => void, () => S] {
  const [state, _setState] = React.useState<S>(initState);
  const stateRef = React.useRef<S>(initState);

  const callbackRef = React.useRef<StateUpdater<S> | undefined>(undefined);
  const changeTrackCBRef = React.useRef<StateUpdater<S> | undefined>(changeTrackCB);
  const isFirstCBCall = React.useRef(true);

  React.useEffect(() => {
    if (isFirstCBCall.current) isFirstCBCall.current = false;
    else {
      callbackRef.current?.(state);
      changeTrackCBRef.current?.(state);
    }
  }, [state]);

  const setState = (state: S, callback?: StateUpdater<S>) => {
    callbackRef.current = callback;

    stateRef.current = state;
    _setState(state);
  };

  const getStateRefValue = () => stateRef.current;

  return [state, setState, getStateRefValue];
}

export default useStateCB;
