import React from 'react';

type StateCallback<T> = (state: T) => void;
type SetStateAction<T> = T | ((prevState: T) => T);
type UseStateCBReturn<T> = [
  T,
  (nextState: SetStateAction<T>, callback?: StateCallback<T>) => void,
  () => T,
];

function useStateCB<T>(
  initState: T | (() => T),
  changeTrackCB?: StateCallback<NoInfer<T>>,
): UseStateCBReturn<T> {
  const [state, _setState] = React.useState<T>(initState);
  const stateRef = React.useRef<T>(state);

  const callbackRef = React.useRef<StateCallback<T> | undefined>(undefined);
  const changeTrackCBRef = React.useRef<StateCallback<T> | undefined>(changeTrackCB);
  const isFirstCBCall = React.useRef(true);

  React.useEffect(() => {
    changeTrackCBRef.current = changeTrackCB;
  }, [changeTrackCB]);

  React.useEffect(() => {
    if (isFirstCBCall.current) {
      isFirstCBCall.current = false;
      return;
    }

    const pendingCallback = callbackRef.current;
    callbackRef.current = undefined;
    if (typeof pendingCallback === 'function') {
      pendingCallback(state);
    }

    if (typeof changeTrackCBRef.current === 'function') {
      changeTrackCBRef.current(state);
    }
  }, [state]);

  const setState = React.useCallback(
    (nextState: SetStateAction<T>, callback?: StateCallback<T>) => {
      callbackRef.current = callback;

      const resolvedState =
        typeof nextState === 'function'
          ? (nextState as (prevState: T) => T)(stateRef.current)
          : nextState;

      stateRef.current = resolvedState;
      _setState(() => resolvedState);
    },
    [],
  );

  const getStateRefValue = React.useCallback(() => stateRef.current, []);

  return [state, setState, getStateRefValue];
}

export default useStateCB;
