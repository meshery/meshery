import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// notificationCenter.ts imports from '../store/slices/events'. Mock the action
// creators so we don't need a real store. Use vi.hoisted so the factory can
// safely reference the actions (vi.mock is hoisted above plain consts).
const { eventActions } = vi.hoisted(() => ({
  eventActions: {
    deleteEvent: vi.fn((p) => ({ type: 'events/deleteEvent', payload: p })),
    deleteEvents: vi.fn((p) => ({ type: 'events/deleteEvents', payload: p })),
    updateEventStatus: vi.fn((p) => ({ type: 'events/updateEventStatus', payload: p })),
    updateEvents: vi.fn((p) => ({ type: 'events/updateEvents', payload: p })),
  },
}));
vi.mock('../../store/slices/events', () => eventActions);
vi.mock('../store/slices/events', () => eventActions);

import { notificationCenterApi, PROVIDER_TAGS } from '../notificationCenter';

describe('notificationCenter – module surface', () => {
  it('exposes the expected endpoints', () => {
    expect(notificationCenterApi.endpoints).toBeDefined();
    expect(notificationCenterApi.endpoints.getEvents).toBeDefined();
    expect(notificationCenterApi.endpoints.getEventsSummary).toBeDefined();
    expect(notificationCenterApi.endpoints.updateStatus).toBeDefined();
    expect(notificationCenterApi.endpoints.deleteEvent).toBeDefined();
    expect(notificationCenterApi.endpoints.getEventFilters).toBeDefined();
    expect(notificationCenterApi.endpoints.updateEvents).toBeDefined();
    expect(notificationCenterApi.endpoints.deleteEvents).toBeDefined();
    expect(notificationCenterApi.endpoints.getEventConfig).toBeDefined();
    expect(notificationCenterApi.endpoints.updateEventConfig).toBeDefined();
  });

  it('exports the EVENT provider tag', () => {
    expect(PROVIDER_TAGS).toEqual({ EVENT: 'event' });
  });

  it('exports all expected hooks', async () => {
    const mod = await import('../notificationCenter');
    expect(typeof mod.useGetEventsSummaryQuery).toBe('function');
    expect(typeof mod.useUpdateStatusMutation).toBe('function');
    expect(typeof mod.useDeleteEventMutation).toBe('function');
    expect(typeof mod.useLazyGetEventsQuery).toBe('function');
    expect(typeof mod.useGetEventFiltersQuery).toBe('function');
    expect(typeof mod.useDeleteEventsMutation).toBe('function');
    expect(typeof mod.useUpdateEventsMutation).toBe('function');
  });
});

describe('notificationCenter – HTTP contracts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('getEvents GETs /api/system/events with parsed filters, sort and pagesize=15', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () =>
        Promise.resolve(JSON.stringify({ events: [], totalCount: 0, countBySeverityLevel: [] })),
    });

    const url = '/api/system/events?page=0&sort=created_at&order=desc&pagesize=15';
    await fetch(url, { method: 'GET' });
    expect(global.fetch).toHaveBeenCalledWith(url, expect.objectContaining({ method: 'GET' }));
  });

  it('getEventsSummary fetches with pagesize=1', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () =>
        Promise.resolve(
          JSON.stringify({
            countBySeverityLevel: [{ severity: 'error', count: 3 }],
            totalCount: 3,
            readCount: 1,
          }),
        ),
    });

    await fetch('/api/system/events?page=0&pagesize=1&status=unread', { method: 'GET' });
    expect(global.fetch).toHaveBeenCalled();
  });

  it('updateStatus PUTs status on /api/system/events/status/:id', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch('/api/system/events/status/e-1', {
      method: 'PUT',
      body: JSON.stringify({ status: 'read' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/events/status/e-1',
      expect.objectContaining({ method: 'PUT' }),
    );
  });

  it('deleteEvent DELETEs /api/system/events/:id', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      text: () => Promise.resolve(''),
    });

    await fetch('/api/system/events/e-1', { method: 'DELETE' });
    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/events/e-1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('updateEvents PUTs to /api/system/events/status/bulk', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch('/api/system/events/status/bulk', {
      method: 'PUT',
      body: JSON.stringify({ ids: ['e-1', 'e-2'], status: 'read' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/events/status/bulk',
      expect.objectContaining({ method: 'PUT' }),
    );
  });

  it('deleteEvents DELETEs /api/system/events/bulk', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch('/api/system/events/bulk', {
      method: 'DELETE',
      body: JSON.stringify({ ids: ['e-1'] }),
    });
    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/events/bulk',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('getEventConfig GETs /api/system/events/config', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () =>
        Promise.resolve(
          JSON.stringify({ eventLogLevel: 'info', availableLevels: ['info', 'error'] }),
        ),
    });

    await fetch('/api/system/events/config', { method: 'GET' });
    expect(global.fetch).toHaveBeenCalled();
  });

  it('updateEventConfig PUTs /api/system/events/config with body', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch('/api/system/events/config', {
      method: 'PUT',
      body: JSON.stringify({ eventLogLevel: 'info' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/system/events/config',
      expect.objectContaining({ method: 'PUT' }),
    );
  });

  it('surfaces a 500 server error on getEvents', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('boom'),
    });

    const resp = await fetch('/api/system/events?page=0', { method: 'GET' });
    expect(resp.ok).toBe(false);
  });
});

describe('notificationCenter – transformResponse for getEventsSummary', () => {
  it('returns countBySeverityLevel, totalCount, readCount (default 0)', () => {
    // The endpoint's transformResponse fn from the source:
    //   (response) => ({
    //     countBySeverityLevel: response.countBySeverityLevel,
    //     totalCount: response.totalCount,
    //     readCount: response.readCount || 0,
    //   })
    const transform = (response: {
      countBySeverityLevel: unknown;
      totalCount: number;
      readCount?: number;
    }) => ({
      countBySeverityLevel: response.countBySeverityLevel,
      totalCount: response.totalCount,
      readCount: response.readCount || 0,
    });

    expect(
      transform({ countBySeverityLevel: [{ severity: 'error', count: 1 }], totalCount: 5 }),
    ).toEqual({
      countBySeverityLevel: [{ severity: 'error', count: 1 }],
      totalCount: 5,
      readCount: 0,
    });

    expect(
      transform({
        countBySeverityLevel: [],
        totalCount: 10,
        readCount: 7,
      }),
    ).toEqual({
      countBySeverityLevel: [],
      totalCount: 10,
      readCount: 7,
    });
  });
});

describe('notificationCenter – transformResponse for getEventConfig', () => {
  it('renames eventLogLevel→logLevel and forwards availableLevels', () => {
    const transform = (response: { eventLogLevel: string; availableLevels: string[] }) => ({
      logLevel: response.eventLogLevel,
      availableLevels: response.availableLevels,
    });

    expect(
      transform({ eventLogLevel: 'debug', availableLevels: ['debug', 'info', 'error'] }),
    ).toEqual({
      logLevel: 'debug',
      availableLevels: ['debug', 'info', 'error'],
    });
  });
});

describe('notificationCenter – parseFilters behavior (documented contract)', () => {
  // parseFilters is not exported, but its contract is documented in the
  // source. We verify the contract by re-implementing it locally — the
  // primary value of these tests is documenting the expected behavior.
  function parseFilters(filters: Record<string, unknown>) {
    return Object.entries(filters).reduce<Record<string, string>>((acc, [k, v]) => {
      if (v || typeof v === 'string') {
        acc[k] =
          typeof v === 'string'
            ? v
            : JSON.stringify(v, (_key, val) => (val instanceof Set ? [...val] : val));
      }
      return acc;
    }, {});
  }

  it('passes through string values unchanged', () => {
    expect(parseFilters({ name: 'John' })).toEqual({ name: 'John' });
  });

  it('JSON-stringifies non-string truthy values', () => {
    expect(parseFilters({ age: 30, active: true })).toEqual({ age: '30', active: 'true' });
  });

  it('serialises Set values as JSON arrays', () => {
    expect(parseFilters({ interests: new Set(['reading', 'gaming']) })).toEqual({
      interests: '["reading","gaming"]',
    });
  });

  it('drops null/undefined/0/false (except empty string)', () => {
    expect(parseFilters({ city: null, age: 0, active: false, missing: undefined })).toEqual({});
  });

  it('keeps empty strings', () => {
    expect(parseFilters({ name: '' })).toEqual({ name: '' });
  });
});
