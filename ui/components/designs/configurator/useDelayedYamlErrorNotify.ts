import { useCallback, useEffect, useRef } from 'react';
import * as jsYaml from 'js-yaml';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';

/** Idle window before the design editor shows an invalid-YAML popup. */
export const DESIGN_YAML_ERROR_IDLE_MS = 400;

export default function useDelayedYamlErrorNotify() {
  const { notify } = useNotification();
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingYamlRef = useRef<string | null>(null);

  const cancelPendingYamlError = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    pendingYamlRef.current = null;
  }, []);

  const scheduleInvalidYamlCheck = useCallback(
    (yamlData: string) => {
      pendingYamlRef.current = yamlData;
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      timeoutRef.current = setTimeout(() => {
        timeoutRef.current = null;
        const text = pendingYamlRef.current;
        pendingYamlRef.current = null;
        if (text === null) {
          return;
        }
        try {
          jsYaml.load(text);
        } catch (err) {
          const details = err instanceof Error ? err.toString() : String(err);
          notify({
            message: `Invalid Yaml Data`,
            event_type: EVENT_TYPES.ERROR,
            details,
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
