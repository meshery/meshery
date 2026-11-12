import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import dataFetch, { promisifiedDataFetch } from '../data-fetch';

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

    await vi.waitFor(() => {
      expect(errorFn).toHaveBeenCalledWith('Internal Server Error');
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

    await expect(promisifiedDataFetch('/api/test')).rejects.toBe('Server Error');
  });
});
