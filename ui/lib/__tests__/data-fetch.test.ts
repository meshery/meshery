import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock the store before importing data-fetch
vi.mock('../../store', () => ({
  store: {
    dispatch: vi.fn(),
  },
}));

import dataFetch, { promisifiedDataFetch } from '../data-fetch';
import { store } from '../../store';

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

    // Wait for async chain
    await vi.waitFor(() => {
      expect(successFn).toHaveBeenCalledWith(mockData);
    });
    expect(errorFn).not.toHaveBeenCalled();
  });

  it('dispatches SESSION_EXPIRED on 401 instead of reloading', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 401,
      redirected: false,
    });

    const successFn = vi.fn();
    const errorFn = vi.fn();

    dataFetch('/api/test', {}, successFn, errorFn);

    await vi.waitFor(() => {
      expect(store.dispatch).toHaveBeenCalledWith({ type: 'SESSION_EXPIRED' });
    });
    expect(successFn).not.toHaveBeenCalled();
  });

  it('dispatches SESSION_EXPIRED on redirected response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      redirected: true,
    });

    dataFetch('/api/test', {}, vi.fn(), vi.fn());

    await vi.waitFor(() => {
      expect(store.dispatch).toHaveBeenCalledWith({ type: 'SESSION_EXPIRED' });
    });
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
