import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import dataFetch, { DataFetchError, promisifiedDataFetch } from '../data-fetch';

describe('dataFetch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('calls successFn with parsed JSON on success', async () => {
    const mockData = { name: 'test' };
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: false,
      text: () => Promise.resolve(JSON.stringify(mockData)),
    });

    const successFn = vi.fn();
    const errorFn = vi.fn();

    dataFetch('/api/test', {}, successFn, errorFn);

    await vi.waitFor(() => {
      expect(successFn).toHaveBeenCalledWith(mockData);
    });
    expect(errorFn).not.toHaveBeenCalled();
  });

  it('does not call successFn or errorFn on 401', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 401,
      redirected: false,
    });

    const successFn = vi.fn();
    const errorFn = vi.fn();

    dataFetch('/api/test', {}, successFn, errorFn);

    // Give the promise chain time to settle
    await new Promise((r) => setTimeout(r, 50));

    expect(successFn).not.toHaveBeenCalled();
    expect(errorFn).not.toHaveBeenCalled();
  });

  it('calls errorFn on non-ok response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      redirected: false,
      text: () => Promise.resolve('Internal Server Error'),
    });

    const errorFn = vi.fn();

    dataFetch('/api/test', {}, vi.fn(), errorFn);

    await vi.waitFor(() => expect(errorFn).toHaveBeenCalled());
    const err = errorFn.mock.calls[0][0];
    expect(err).toBeInstanceOf(DataFetchError);
    expect(err).toMatchObject({
      message: 'Internal Server Error',
      status: 500,
      body: 'Internal Server Error',
      url: '/api/test',
    });
  });

  it('returns plain text when JSON parsing fails', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: false,
      text: () => Promise.resolve('not json'),
    });

    const successFn = vi.fn();
    dataFetch('/api/test', {}, successFn, vi.fn());

    await vi.waitFor(() => {
      expect(successFn).toHaveBeenCalledWith('not json');
    });
  });
});

describe('promisifiedDataFetch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('resolves with data on success', async () => {
    const mockData = { result: 'ok' };
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: false,
      text: () => Promise.resolve(JSON.stringify(mockData)),
    });

    const result = await promisifiedDataFetch('/api/test');
    expect(result).toEqual(mockData);
  });

  it('rejects on error response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      redirected: false,
      text: () => Promise.resolve('Server Error'),
    });

    await expect(promisifiedDataFetch('/api/test')).rejects.toMatchObject({
      name: 'DataFetchError',
      message: 'Server Error',
      status: 500,
    });
  });

  it('rejects with an Error so consumers can read error.message', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(new TypeError('Failed to fetch'));

    const err = await promisifiedDataFetch('/api/test').catch((e) => e);
    expect(err).toBeInstanceOf(Error);
    expect(err.message).toContain('Unable to reach the Meshery server');
  });
});

describe('dataFetch error reporting', () => {
  const errorResponse = (status, statusText, text) => ({
    ok: false,
    status,
    statusText,
    redirected: false,
    text,
  });

  const fetchError = async (options = {}) => {
    const errorFn = vi.fn();
    dataFetch('/api/test', options, vi.fn(), errorFn);
    await vi.waitFor(() => expect(errorFn).toHaveBeenCalled());
    return errorFn.mock.calls[0][0];
  };

  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('describes the status when the error body is empty', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
      errorResponse(500, 'Internal Server Error', () => Promise.resolve('')),
    );

    const err = await fetchError();
    expect(err.message).toBe('Request to /api/test failed with status 500 Internal Server Error');
    expect(err).toMatchObject({ status: 500, body: '' });
  });

  it('omits statusText when the response has none', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
      errorResponse(503, '', () => Promise.resolve('  ')),
    );

    expect((await fetchError()).message).toBe('Request to /api/test failed with status 503');
  });

  it('keeps an HTML error page out of the message but on body', async () => {
    const html = '<!DOCTYPE html><html><body><h1>504 Gateway Time-out</h1></body></html>';
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
      errorResponse(504, 'Gateway Timeout', () => Promise.resolve(html)),
    );

    const err = await fetchError();
    expect(err.message).toBe('Request to /api/test failed with status 504 Gateway Timeout');
    expect(err.body).toBe(html);
  });

  it('describes the status when the error body cannot be read', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(
      errorResponse(502, 'Bad Gateway', () => Promise.reject(new TypeError('body stream'))),
    );

    const err = await fetchError();
    expect(err).toBeInstanceOf(DataFetchError);
    expect(err.message).toBe('Request to /api/test failed with status 502 Bad Gateway');
  });

  it('names the URL when the server cannot be reached', async () => {
    const cause = new TypeError('Failed to fetch');
    (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(cause);

    const err = await fetchError();
    expect(err).toBeInstanceOf(DataFetchError);
    expect(err.message).toBe('Unable to reach the Meshery server (/api/test): Failed to fetch');
    expect(err.status).toBeUndefined();
    expect(err.cause).toBe(cause);
  });

  it('does not reclassify an error thrown by successFn', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: false,
      text: () => Promise.resolve('{}'),
    });
    const thrown = new TypeError('consumer bug');
    const errorFn = vi.fn();

    dataFetch(
      '/api/test',
      {},
      () => {
        throw thrown;
      },
      errorFn,
    );

    await vi.waitFor(() => expect(errorFn).toHaveBeenCalledWith(thrown));
  });

  it('does not add a signal or timer when no timeout is given', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: false,
      text: () => Promise.resolve('{}'),
    });
    const successFn = vi.fn();

    dataFetch('/api/test', { method: 'POST' }, successFn, vi.fn());

    await vi.waitFor(() => expect(successFn).toHaveBeenCalled());
    expect(global.fetch).toHaveBeenCalledWith('/api/test', { method: 'POST' });
  });

  describe('opt-in timeout', () => {
    // A fetch that never settles on its own and rejects like the browser when aborted.
    const hangingFetch = () =>
      (global.fetch as ReturnType<typeof vi.fn>).mockImplementation(
        (_url, init) =>
          new Promise((_resolve, reject) => {
            init.signal.addEventListener('abort', () =>
              reject(new DOMException('The operation was aborted.', 'AbortError')),
            );
          }),
      );

    it('aborts and reports a timeout after the given milliseconds', async () => {
      vi.useFakeTimers();
      hangingFetch();
      const errorFn = vi.fn();

      dataFetch('/api/test', { timeout: 1000 }, vi.fn(), errorFn);
      await vi.advanceTimersByTimeAsync(999);
      expect(errorFn).not.toHaveBeenCalled();
      await vi.advanceTimersByTimeAsync(1);

      expect(errorFn).toHaveBeenCalledTimes(1);
      const err = errorFn.mock.calls[0][0];
      expect(err).toBeInstanceOf(DataFetchError);
      expect(err.message).toBe('Request to /api/test timed out after 1000ms');
      expect(err.cause.name).toBe('AbortError');
    });

    it('does not pass the timeout option through to fetch', async () => {
      hangingFetch();

      dataFetch('/api/test', { timeout: 1000, method: 'GET' }, vi.fn(), vi.fn());

      const init = (global.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
      expect(init).not.toHaveProperty('timeout');
      expect(init.method).toBe('GET');
      expect(init.signal).toBeInstanceOf(AbortSignal);
    });

    it('clears the timer once the request completes', async () => {
      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        status: 200,
        redirected: false,
        text: () => Promise.resolve('{}'),
      });
      const setTimeoutSpy = vi.spyOn(globalThis, 'setTimeout');
      const clearTimeoutSpy = vi.spyOn(globalThis, 'clearTimeout');
      const successFn = vi.fn();

      dataFetch('/api/test', { timeout: 1000 }, successFn, vi.fn());
      const requestTimer = setTimeoutSpy.mock.results[0].value;

      await vi.waitFor(() => expect(clearTimeoutSpy).toHaveBeenCalledWith(requestTimer));
      expect(successFn).toHaveBeenCalled();
    });

    it('clears the timer when the session has expired (401)', async () => {
      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: false,
        status: 401,
        redirected: false,
      });
      const setTimeoutSpy = vi.spyOn(globalThis, 'setTimeout');
      const clearTimeoutSpy = vi.spyOn(globalThis, 'clearTimeout');

      dataFetch('/api/test', { timeout: 1000 }, vi.fn(), vi.fn());
      const requestTimer = setTimeoutSpy.mock.results[0].value;

      // The 401 promise never settles, so `finally` never runs; the timer must
      // be cleared on the 401 path itself.
      await vi.waitFor(() => expect(clearTimeoutSpy).toHaveBeenCalledWith(requestTimer));
    });

    it('detaches its listener from a reused caller signal once the request settles', async () => {
      (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
        ok: true,
        status: 200,
        redirected: false,
        text: () => Promise.resolve('{}'),
      });
      const caller = new AbortController();
      const addSpy = vi.spyOn(caller.signal, 'addEventListener');
      const removeSpy = vi.spyOn(caller.signal, 'removeEventListener');
      const successFn = vi.fn();

      dataFetch('/api/test', { timeout: 1000, signal: caller.signal }, successFn, vi.fn());

      await vi.waitFor(() => expect(successFn).toHaveBeenCalled());
      await vi.waitFor(() => expect(removeSpy).toHaveBeenCalled());
      expect(removeSpy.mock.calls[0][1]).toBe(addSpy.mock.calls[0][1]);
    });

    it('passes a caller abort through as-is instead of reporting a timeout', async () => {
      hangingFetch();
      const caller = new AbortController();
      const errorFn = vi.fn();

      dataFetch('/api/test', { timeout: 1000, signal: caller.signal }, vi.fn(), errorFn);
      caller.abort();

      await vi.waitFor(() => expect(errorFn).toHaveBeenCalled());
      const err = errorFn.mock.calls[0][0];
      expect(err).not.toBeInstanceOf(DataFetchError);
      expect(err.name).toBe('AbortError');
    });

    it('honours a caller signal that is already aborted', async () => {
      hangingFetch();
      const caller = new AbortController();
      caller.abort();

      dataFetch('/api/test', { timeout: 1000, signal: caller.signal }, vi.fn(), vi.fn());

      const init = (global.fetch as ReturnType<typeof vi.fn>).mock.calls[0][1];
      expect(init.signal.aborted).toBe(true);
    });
  });
});
