import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import adapterReducer, { updateAdaptersInfo, setAdapter } from '../adapter';

describe('adapter slice', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2024-06-01T00:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns the initial state for an unknown action', () => {
    const state = adapterReducer(undefined, { type: 'unknown/action' } as any);
    expect(state.meshAdapters).toEqual([]);
    expect(state.selectedAdapter).toBe('');
    expect(typeof state.meshAdaptersts).toBe('number');
  });

  it('updateAdaptersInfo replaces meshAdapters and updates the timestamp', () => {
    vi.setSystemTime(new Date('2024-06-01T00:00:00Z'));
    const initial = adapterReducer(undefined, { type: 'init' } as any);
    // Seed the state through the reducer to capture a fake-timer timestamp.
    const seeded = adapterReducer(initial, updateAdaptersInfo({ meshAdapters: [] }));
    const seededTs = seeded.meshAdaptersts;

    vi.setSystemTime(new Date('2024-06-01T00:01:00Z'));

    const adapters = [{ name: 'istio' }, { name: 'linkerd' }];
    const state = adapterReducer(seeded, updateAdaptersInfo({ meshAdapters: adapters }));
    expect(state.meshAdapters).toEqual(adapters);
    expect(state.meshAdaptersts).toBeGreaterThan(seededTs);
  });

  it('updateAdaptersInfo replaces empty arrays correctly', () => {
    const seeded = adapterReducer(
      undefined,
      updateAdaptersInfo({ meshAdapters: [{ name: 'istio' }] }),
    );
    expect(seeded.meshAdapters).toHaveLength(1);

    const cleared = adapterReducer(seeded, updateAdaptersInfo({ meshAdapters: [] }));
    expect(cleared.meshAdapters).toEqual([]);
  });

  it('updateAdaptersInfo is immutable', () => {
    const initial = adapterReducer(undefined, { type: 'init' } as any);
    const next = adapterReducer(initial, updateAdaptersInfo({ meshAdapters: [{ name: 'x' }] }));
    expect(next).not.toBe(initial);
    expect(initial.meshAdapters).toEqual([]);
  });

  it('setAdapter sets the selected adapter', () => {
    const initial = adapterReducer(undefined, { type: 'init' } as any);
    const state = adapterReducer(initial, setAdapter({ selectedAdapter: 'istio' }));
    expect(state.selectedAdapter).toBe('istio');
  });

  it('setAdapter preserves meshAdapters and meshAdaptersts', () => {
    let state = adapterReducer(
      undefined,
      updateAdaptersInfo({ meshAdapters: [{ name: 'istio' }] }),
    );
    const adaptersBefore = state.meshAdapters;
    const tsBefore = state.meshAdaptersts;
    state = adapterReducer(state, setAdapter({ selectedAdapter: 'linkerd' }));
    expect(state.meshAdapters).toEqual(adaptersBefore);
    expect(state.meshAdaptersts).toBe(tsBefore);
    expect(state.selectedAdapter).toBe('linkerd');
  });

  it('setAdapter is immutable', () => {
    const initial = adapterReducer(undefined, { type: 'init' } as any);
    const next = adapterReducer(initial, setAdapter({ selectedAdapter: 'foo' }));
    expect(next).not.toBe(initial);
    expect(initial.selectedAdapter).toBe('');
  });

  it('setAdapter accepts empty string to clear selection', () => {
    let state = adapterReducer(undefined, setAdapter({ selectedAdapter: 'foo' }));
    state = adapterReducer(state, setAdapter({ selectedAdapter: '' }));
    expect(state.selectedAdapter).toBe('');
  });

  it('actions have the expected type prefixes', () => {
    expect(updateAdaptersInfo({ meshAdapters: [] }).type).toBe('adapter/updateAdaptersInfo');
    expect(setAdapter({ selectedAdapter: '' }).type).toBe('adapter/setAdapter');
  });
});
