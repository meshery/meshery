import { useSelector, useDispatch } from 'react-redux';
import { updateTableState, selectTableState } from '../../store/slices/tableState';
import { useRouter } from 'next/router';
import { useEffect, useRef } from 'react';

// Shared debounce timer to avoid race conditions when multiple states update simultaneously
let routerReplaceTimeout: ReturnType<typeof setTimeout>;

export function useTableState<T>(
  tableId: string,
  key: string,
  defaultValue: T,
): [T, (val: T | ((prev: T) => T)) => void] {
  const dispatch = useDispatch();
  const fullTableState = useSelector((state: any) => selectTableState(state, tableId));
  const router = useRouter();

  const queryKey = `${tableId}_${key}`;
  let value = fullTableState[key] !== undefined ? fullTableState[key] : defaultValue;
  const initialized = useRef(false);

  useEffect(() => {
    if (!router.isReady) return;

    const queryVal = router.query[queryKey];

    // 1. On first mount, if URL has data but Redux is empty, pull URL into Redux
    if (!initialized.current) {
      initialized.current = true;
      if (queryVal !== undefined && fullTableState[key] === undefined) {
        const stringVal = Array.isArray(queryVal) ? queryVal[0] : queryVal;
        let parsedVal: any = stringVal;
        if (typeof defaultValue === 'number') parsedVal = Number(stringVal);
        else if (typeof defaultValue === 'boolean') parsedVal = stringVal === 'true';

        dispatch(updateTableState({ tableId, update: { [key]: parsedVal } }));
        return; // State injected to Redux, skip URL update for this render
      }
    }

    // 2. Sync Redux state into the URL
    if (fullTableState[key] !== undefined) {
      const stringVal = String(fullTableState[key]);

      // If URL does not match Redux state, update the URL
      if (queryVal !== stringVal) {
        clearTimeout(routerReplaceTimeout);
        routerReplaceTimeout = setTimeout(() => {
          const urlParams = new URLSearchParams(window.location.search);

          Object.keys(fullTableState).forEach((k) => {
            if (fullTableState[k] !== undefined) {
              urlParams.set(`${tableId}_${k}`, String(fullTableState[k]));
            }
          });

          router.replace(
            {
              pathname: router.pathname,
              query: Object.fromEntries(urlParams.entries()),
            },
            undefined,
            { shallow: true },
          );
        }, 100);
      }
    }
  }, [router.isReady, fullTableState[key]]);

  const setValue = (newValue: T | ((prev: T) => T)) => {
    const resolvedValue = typeof newValue === 'function' ? (newValue as Function)(value) : newValue;
    dispatch(
      updateTableState({
        tableId,
        update: { [key]: resolvedValue },
      }),
    );
  };

  return [value, setValue];
}
