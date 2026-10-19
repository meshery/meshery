import { useCallback } from 'react';
import { useDispatch } from 'react-redux';
import { useNotification } from './useNotification';
import { getErrorMessage } from '../../components/connections/ConnectionTable.constants';
import { EVENT_TYPES } from '../../lib/event-types';
import { updateProgressAction } from '@/store/slices/mesheryUi';

/** Per-attempt budget for the reachability probe (both attempts share it). */
const PING_TIMEOUT_MS = 8000;

/** Unauthenticated endpoint every Meshery Server exposes (see ServerVersionHandler). */
export const MESHERY_VERSION_PATH = '/api/system/version';

export type MesheryServerProbe =
  | {
      reachable: true;
      /** Build reported by the live server (e.g. `stable-v0.8.76`), when readable. */
      build?: string;
      commitsha?: string;
      /**
       * True when the body could not be read (cross-origin target without CORS
       * headers) — reachability still holds, the version just isn't visible.
       */
      bodyUnavailable: boolean;
    }
  | { reachable: false; reason: string };

/**
 * Resolve `<serverLocation>/api/system/version`, or null when the recorded
 * value isn't an http(s) URL.
 */
export const mesheryVersionURL = (serverLocation: string): string | null => {
  try {
    const url = new URL(serverLocation.trim());
    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return null;
    }
    url.pathname = `${url.pathname.replace(/\/+$/, '')}${MESHERY_VERSION_PATH}`;
    url.search = '';
    url.hash = '';
    return url.toString();
  } catch {
    return null;
  }
};

// `Parameters<typeof fetch>[1]` rather than the `RequestInit` identifier:
// eslint's no-undef flags the DOM lib type name in this config (the one prior
// usage in the repo sits in a test file, where no-undef is off).
const fetchWithTimeout = (url: string, init: Parameters<typeof fetch>[1]): Promise<Response> => {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), PING_TIMEOUT_MS);
  return fetch(url, { ...init, cache: 'no-store', signal: controller.signal }).finally(() =>
    clearTimeout(timer),
  );
};

/**
 * probeMesheryServer checks whether the Meshery Server recorded on a
 * `meshery`-kind connection is reachable.
 *
 * Why a browser-side probe rather than a server endpoint: no schemas-defined
 * API pings a connection by id (`pingKubernetes` is the only ping operation in
 * @meshery/schemas), and per the contributing guide a new endpoint has to be
 * defined in meshery/schemas first — so this stays a pure UI-side check
 * against the connection's own recorded `serverLocation`.
 *
 * Two attempts, mirroring what a browser can learn about an arbitrary origin:
 *
 *  1. A normal GET of `/api/system/version`. When the body is readable (same
 *     origin, or a CORS-enabled target) the live build/commit come back too.
 *  2. On a network/CORS failure, retry with `mode: 'no-cors'`. The promise
 *     resolving proves *something* answered at that address — the response is
 *     opaque (status 0) but that is exactly the signal a reachability ping
 *     needs; a rejection means nothing answered.
 *
 * Any HTTP status (401, 404, …) still counts as reachable: a server is
 * listening there. Only an invalid URL, no answer within PING_TIMEOUT_MS, or a
 * failed second attempt reports unreachable.
 */
export const probeMesheryServer = async (serverLocation: string): Promise<MesheryServerProbe> => {
  const location = serverLocation?.trim();
  if (!location) {
    return { reachable: false, reason: 'No server location recorded for this connection' };
  }

  const url = mesheryVersionURL(location);
  if (!url) {
    return {
      reachable: false,
      reason: `"${location}" is not a valid http(s) server address`,
    };
  }

  // Attempt 1: readable request — gives us the live version when permitted.
  try {
    const res = await fetchWithTimeout(url, { method: 'GET' });
    let body: { build?: string; commitsha?: string } | null = null;
    if (res.ok) {
      body = await res.json().catch(() => null);
    }
    return {
      reachable: true,
      build: typeof body?.build === 'string' ? body.build : undefined,
      commitsha: typeof body?.commitsha === 'string' ? body.commitsha : undefined,
      bodyUnavailable: false,
    };
  } catch (err) {
    // A timeout means nothing answered in budget — don't burn a second timeout.
    if (err instanceof DOMException && err.name === 'AbortError') {
      return {
        reachable: false,
        reason: `No response from ${location} within ${PING_TIMEOUT_MS}ms`,
      };
    }
    // Fall through: most likely a CORS/network rejection on the readable form.
  }

  // Attempt 2: opaque probe — resolves for any HTTP answer regardless of CORS.
  try {
    await fetchWithTimeout(url, { method: 'GET', mode: 'no-cors' });
    return { reachable: true, bodyUnavailable: true };
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      return {
        reachable: false,
        reason: `No response from ${location} within ${PING_TIMEOUT_MS}ms`,
      };
    }
    return {
      reachable: false,
      reason: getErrorMessage(err, `No response from ${location}`),
    };
  }
};

/**
 * useMesheryPingHook returns an imperative `ping(serverLocation, name)` that
 * probes a `meshery`-kind connection's recorded Meshery Server and surfaces
 * the outcome as a notification, mirroring useGrafanaPingHook /
 * usePrometheusPingHook. The global progress loader is shown while the probe
 * runs, so a click on the connection chip gives immediate feedback.
 */
export default function useMesheryPingHook() {
  const { notify } = useNotification();
  const dispatch = useDispatch();

  // Memoized so consumers can list `ping` in hook dep arrays without
  // invalidating their memos every render.
  const ping = useCallback(
    async (serverLocation?: string, name?: string) => {
      const label = name?.trim() || 'Meshery Server';
      const location = serverLocation?.trim() || 'unknown server';

      dispatch(updateProgressAction({ showProgress: true }));
      try {
        const probe = await probeMesheryServer(serverLocation ?? '');

        if (!probe.reachable) {
          notify({
            message: `Connection failed for ${label} (${location}) - unable to reach Meshery Server`,
            details: probe.reason,
            event_type: EVENT_TYPES.ERROR,
          });
          return;
        }

        let message = `Connected successfully to ${label} (${location}) - Meshery Server is reachable`;
        if (probe.build) {
          message = `Connected successfully to ${label} (${location}) - Meshery ${probe.build}`;
        } else if (probe.bodyUnavailable) {
          message = `Connected successfully to ${label} (${location}) - Meshery Server is reachable (build unknown)`;
        }
        notify({ message, event_type: EVENT_TYPES.SUCCESS });
      } catch (err) {
        notify({
          message: `Connection failed for ${label} (${location}) - unable to reach Meshery Server`,
          details: getErrorMessage(err, 'Unable to reach Meshery Server'),
          event_type: EVENT_TYPES.ERROR,
        });
      } finally {
        dispatch(updateProgressAction({ showProgress: false }));
      }
    },
    [dispatch, notify],
  );

  return ping;
}
