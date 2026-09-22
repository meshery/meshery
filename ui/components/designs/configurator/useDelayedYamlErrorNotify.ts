import { useCallback, useEffect, useRef } from 'react';
import * as jsYaml from 'js-yaml';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';

/** Idle window before the design editor shows an invalid-YAML popup. */
export const DESIGN_YAML_ERROR_IDLE_MS = 400;

export default function useDelayedYamlErrorNotify() {
  const { notify } = useNotification();
  const timeoutRef = useRef(null);
  const pendingYamlRef = useRef(null);

  const cancelPendingYamlError = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    pendingYamlRef.current = null;
  }, []);

  const scheduleInvalidYamlCheck = useCallback(
    (yamlData) => {
      pendingYamlRef.current = yamlData;
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        timeoutRef.current = null;
        const text = pendingYamlRef.current;
        pendingYamlRef.current = null;
        try {
          jsYaml.load(text);
        } catch (err) {
          notify({
            message: `Invalid Yaml Data`,
            event_type: EVENT_TYPES.ERROR,
            details: err.toString(),
          });
        }
      }, DESIGN_YAML_ERROR_IDLE_MS);
    },
    [notify],
  );

  useEffect(() => {
    return () => cancelPendingYamlError();
  }, [cancelPendingYamlError]);

  return { scheduleInvalidYamlCheck, cancelPendingYamlError };
}
