import { renderHook, act } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
const dispatch = vi.fn();
const fetchMock = vi.fn();

vi.mock('../useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('react-redux', () => ({
  useDispatch: () => dispatch,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgressAction: (payload: unknown) => ({ type: 'ui/progress', payload }),
}));

import useMesheryPingHook, { mesheryVersionURL, MESHERY_VERSION_PATH } from '../useMesheryPingHook';
import { EVENT_TYPES } from '../../../lib/event-types';

const PROBE_URL = `https://meshery.example.com${MESHERY_VERSION_PATH}`;

describe('mesheryVersionURL', () => {
  it('appends the version endpoint and strips trailing slashes, query and hash', () => {
    expect(mesheryVersionURL('https://meshery.example.com/')).toBe(PROBE_URL);
    expect(mesheryVersionURL('https://meshery.example.com')).toBe(PROBE_URL);
    expect(mesheryVersionURL('  https://meshery.example.com  ')).toBe(PROBE_URL);
    expect(mesheryVersionURL('https://meshery.example.com/base/?a=1#x')).toBe(
      `https://meshery.example.com/base${MESHERY_VERSION_PATH}`,
    );
  });

  it('rejects values that are not http(s) addresses', () => {
    expect(mesheryVersionURL('')).toBeNull();
    expect(mesheryVersionURL('not a url')).toBeNull();
    expect(mesheryVersionURL('ftp://meshery.example.com')).toBeNull();
    expect(mesheryVersionURL('javascript:alert(1)')).toBeNull();
  });
});

describe('useMesheryPingHook', () => {
  beforeEach(() => {
    notify.mockReset();
    dispatch.mockReset();
    fetchMock.mockReset();
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('reports the live build when the version endpoint is readable', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ build: 'stable-v1.2.3' }),
    });
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('https://meshery.example.com', 'meshery-local');
    });

    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(fetchMock).toHaveBeenCalledWith(
      PROBE_URL,
      expect.objectContaining({ method: 'GET', cache: 'no-store' }),
    );
    expect(notify).toHaveBeenCalledWith({
      message:
        'Connected successfully to meshery-local (https://meshery.example.com) - Meshery stable-v1.2.3',
      event_type: EVENT_TYPES.SUCCESS,
    });
    // Progress loader is shown for the probe and dismissed afterwards.
    expect(dispatch).toHaveBeenNthCalledWith(1, {
      type: 'ui/progress',
      payload: { showProgress: true },
    });
    expect(dispatch).toHaveBeenLastCalledWith({
      type: 'ui/progress',
      payload: { showProgress: false },
    });
  });

  it('falls back to an opaque no-cors probe when the readable attempt is blocked', async () => {
    fetchMock
      .mockRejectedValueOnce(new TypeError('Failed to fetch'))
      .mockResolvedValueOnce({ type: 'opaque' });
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('https://other.example.com', 'remote-meshery');
    });

    expect(fetchMock).toHaveBeenCalledTimes(2);
    expect(fetchMock).toHaveBeenLastCalledWith(
      'https://other.example.com/api/system/version',
      expect.objectContaining({ method: 'GET', mode: 'no-cors' }),
    );
    expect(notify).toHaveBeenCalledWith({
      message:
        'Connected successfully to remote-meshery (https://other.example.com) - Meshery Server is reachable (build unknown)',
      event_type: EVENT_TYPES.SUCCESS,
    });
  });

  it('treats a non-2xx answer as reachable — something is listening there', async () => {
    fetchMock.mockResolvedValueOnce({ ok: false, status: 404 });
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('https://meshery.example.com');
    });

    expect(notify).toHaveBeenCalledWith({
      message:
        'Connected successfully to Meshery Server (https://meshery.example.com) - Meshery Server is reachable',
      event_type: EVENT_TYPES.SUCCESS,
    });
  });

  it('reports failure when both attempts are refused', async () => {
    fetchMock
      .mockRejectedValueOnce(new TypeError('Failed to fetch'))
      .mockRejectedValueOnce(new TypeError('Failed to fetch'));
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('https://down.example.com', 'down-meshery');
    });

    expect(notify).toHaveBeenCalledWith({
      message:
        'Connection failed for down-meshery (https://down.example.com) - unable to reach Meshery Server',
      details: 'Failed to fetch',
      event_type: EVENT_TYPES.ERROR,
    });
    // The loader still comes down on the failure path.
    expect(dispatch).toHaveBeenLastCalledWith({
      type: 'ui/progress',
      payload: { showProgress: false },
    });
  });

  it('reports failure when the probe times out', async () => {
    fetchMock.mockRejectedValueOnce(new DOMException('Aborted', 'AbortError'));
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('https://slow.example.com');
    });

    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ event_type: EVENT_TYPES.ERROR }));
  });

  it('fails without fetching when no server location is recorded', async () => {
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current(undefined, 'meshery-local');
    });

    expect(fetchMock).not.toHaveBeenCalled();
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        event_type: EVENT_TYPES.ERROR,
        details: 'No server location recorded for this connection',
      }),
    );
  });

  it('fails without fetching for a non-http(s) server location', async () => {
    const { result } = renderHook(() => useMesheryPingHook());

    await act(async () => {
      await result.current('ftp://meshery.example.com', 'meshery-local');
    });

    expect(fetchMock).not.toHaveBeenCalled();
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        event_type: EVENT_TYPES.ERROR,
        details: '"ftp://meshery.example.com" is not a valid http(s) server address',
      }),
    );
  });
});
