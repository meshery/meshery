import { useSelector, useDispatch } from 'react-redux';
import { updateTableState, selectTableState } from '../../store/slices/tableState';
import { useRouter } from 'next/router';
import { useEffect, useRef } from 'react';

// Per-table debounce timers — keyed by tableId so updates on one table
// never cancel a pending URL write from another table.
const routerReplaceTimeouts = new Map<string, ReturnType<typeof setTimeout>>();

export function useTableState<T>(
  tableId: string,
  key: string,
  defaultValue: T,
): [T, (val: T | ((prev: T) => T)) => void] {
  const dispatch = useDispatch();
  const fullTableState = useSelector((state: any) => selectTableState(state, tableId));
  const router = useRouter();

  const queryKey = `${tableId}_${key}`;
  const value = fullTableState[key] !== undefined ? fullTableState[key] : defaultValue;
  
  const initialized = useRef(false);
  const lastSeenQueryVal = useRef<string | undefined>(undefined);

  useEffect(() => {
    if (!router.isReady) return;

    const queryVal = router.query[queryKey] as string | undefined;
    const stringQueryVal = Array.isArray(queryVal) ? queryVal[0] : queryVal;

    // 1. On first mount: if Redux is empty and the URL has a value, seed Redux from the URL.
    //    This handles hard refreshes — the URL is the persistence layer.
    if (!initialized.current) {
      initialized.current = true;
      lastSeenQueryVal.current = stringQueryVal;
      
      if (fullTableState[key] === undefined && stringQueryVal !== undefined) {
        let parsedVal: any = stringQueryVal;
        if (typeof defaultValue === 'number') parsedVal = Number(stringQueryVal);
        else if (typeof defaultValue === 'boolean') parsedVal = stringQueryVal === 'true';
        else if (stringQueryVal === 'null') parsedVal = null;

        dispatch(updateTableState({ tableId, update: { [key]: parsedVal } }));
        return; // State injected into Redux; skip URL update this render
      }
    }

    const stateVal = fullTableState[key];
    if (stateVal === undefined) return;

    const stateStringVal = stateVal === null ? '' : String(stateVal);
    const urlStringVal = stringQueryVal === undefined ? '' : String(stringQueryVal);

    // 2. External URL Change Detection (e.g. Back button or Sidebar click)
    // If the URL value changed from what we last saw, AND it doesn't match Redux,
    // the user used browser navigation. We must sync URL -> Redux.
    if (stringQueryVal !== lastSeenQueryVal.current && urlStringVal !== stateStringVal) {
      lastSeenQueryVal.current = stringQueryVal;
      let parsedVal: any = stringQueryVal;
      
      if (stringQueryVal === undefined) parsedVal = defaultValue;
      else if (typeof defaultValue === 'number') parsedVal = Number(stringQueryVal);
      else if (typeof defaultValue === 'boolean') parsedVal = stringQueryVal === 'true';
      else if (stringQueryVal === 'null') parsedVal = null;
      
      dispatch(updateTableState({ tableId, update: { [key]: parsedVal } }));
      return;
    }

    lastSeenQueryVal.current = stringQueryVal;

    // 3. Keep the URL in sync with Redux state so the URL always reflects the
    //    current table configuration and persists across hard refreshes.
    if (stateStringVal !== urlStringVal) {
      clearTimeout(routerReplaceTimeouts.get(tableId));
      routerReplaceTimeouts.set(
        tableId,
        setTimeout(() => {
          const urlParams = new URLSearchParams(window.location.search);

          Object.keys(fullTableState).forEach((k) => {
            const v = fullTableState[k];
            if (v === null || v === '' || v === undefined) {
              // Remove the param from URL when the value is cleared/reset to default
              urlParams.delete(`${tableId}_${k}`);
            } else {
              urlParams.set(`${tableId}_${k}`, String(v));
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
        }, 100),
      );
    }
  }, [router.isReady, router.query, fullTableState[key]]);

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

