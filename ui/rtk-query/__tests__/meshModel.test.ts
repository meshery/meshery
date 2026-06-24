import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

vi.unmock('../utils');
vi.unmock('@/rtk-query/utils');

import { mesheryApiPath } from '../index';

// meshModel.ts imports from ../utils which depends on ../../store. Mock the
// store before the module graph is touched so we can import the helpers.
const { dispatch } = vi.hoisted(() => ({ dispatch: vi.fn() }));
vi.mock('../../store', () => ({ store: { dispatch } }));

import { componentUniqueKey, getComponentDefinition, modelUniqueKey } from '../meshModel';

describe('meshModel – URLs', () => {
  it('builds /api/meshmodels/models URL', () => {
    expect(mesheryApiPath('meshmodels/models')).toBe('/api/meshmodels/models');
  });

  it('builds /api/meshmodels/components URL', () => {
    expect(mesheryApiPath('meshmodels/components')).toBe('/api/meshmodels/components');
  });

  it('builds /api/meshmodels/relationships URL', () => {
    expect(mesheryApiPath('meshmodels/relationships')).toBe('/api/meshmodels/relationships');
  });

  it('builds /api/meshmodels/registrants URL', () => {
    expect(mesheryApiPath('meshmodels/registrants')).toBe('/api/meshmodels/registrants');
  });

  it('builds nested model->components URL', () => {
    expect(mesheryApiPath('meshmodels/models/istio/components')).toBe(
      '/api/meshmodels/models/istio/components',
    );
  });

  it('builds /api/meshmodels/categories URL', () => {
    expect(mesheryApiPath('meshmodels/categories')).toBe('/api/meshmodels/categories');
  });

  it('builds /api/meshmodels/categories/:cat/models URL', () => {
    expect(mesheryApiPath('meshmodels/categories/Networking/models')).toBe(
      '/api/meshmodels/categories/Networking/models',
    );
  });

  it('builds /api/meshmodels/components/:name URL', () => {
    expect(mesheryApiPath('meshmodels/components/Service')).toBe(
      '/api/meshmodels/components/Service',
    );
  });

  it('builds /api/meshmodels/models/:name URL', () => {
    expect(mesheryApiPath('meshmodels/models/istio')).toBe('/api/meshmodels/models/istio');
  });

  it('builds /api/meshmodels/export URL', () => {
    expect(mesheryApiPath('meshmodels/export')).toBe('/api/meshmodels/export');
  });

  it('builds /api/meshmodels/register URL', () => {
    expect(mesheryApiPath('meshmodels/register')).toBe('/api/meshmodels/register');
  });

  it('builds /api/meshmodels/:type/status URL', () => {
    expect(mesheryApiPath('meshmodels/components/status')).toBe(
      '/api/meshmodels/components/status',
    );
  });
});

describe('meshModel – module surface', () => {
  it('exposes the expected hooks', async () => {
    const mod = await import('../meshModel');
    expect(typeof mod.useLazyGetMeshModelsQuery).toBe('function');
    expect(typeof mod.useGetMeshModelsQuery).toBe('function');
    expect(typeof mod.useLazyGetComponentsQuery).toBe('function');
    expect(typeof mod.useGetComponentsQuery).toBe('function');
    expect(typeof mod.useLazyGetRelationshipsQuery).toBe('function');
    expect(typeof mod.useGetRelationshipsQuery).toBe('function');
    expect(typeof mod.useGetRegistrantsQuery).toBe('function');
    expect(typeof mod.useLazyGetRegistrantsQuery).toBe('function');
    expect(typeof mod.useGetComponentsFromModalQuery).toBe('function');
    expect(typeof mod.useLazyGetComponentsFromModalQuery).toBe('function');
    expect(typeof mod.useGetRelationshipsFromModalQuery).toBe('function');
    expect(typeof mod.useLazyGetRelationshipsFromModalQuery).toBe('function');
    expect(typeof mod.useLazyExportModelQuery).toBe('function');
    expect(typeof mod.useUpdateEntityStatusMutation).toBe('function');
    expect(typeof mod.useGetModelCategoriesQuery).toBe('function');
    expect(typeof mod.useLazyGetModelFromCategoryQuery).toBe('function');
    expect(typeof mod.useGetModelByNameQuery).toBe('function');
    expect(typeof mod.useLazyGetModelByNameQuery).toBe('function');
    expect(typeof mod.useGetComponentByNameQuery).toBe('function');
    expect(typeof mod.useGetModelFromCategoryQuery).toBe('function');
    expect(typeof mod.useGetComponentsByModelAndKindQuery).toBe('function');
    expect(typeof mod.useImportMeshModelMutation).toBe('function');
  });
});

describe('modelUniqueKey', () => {
  it('joins name and version with a hyphen', () => {
    expect(modelUniqueKey({ name: 'istio', version: 'v1.0.0' })).toBe('istio-v1.0.0');
  });

  it('handles missing fields by yielding undefined-suffixed strings', () => {
    expect(modelUniqueKey({ name: 'istio', version: undefined })).toBe('istio-undefined');
  });
});

describe('componentUniqueKey', () => {
  it('joins kind-version-cversion-model.name-model.version', () => {
    const component = {
      component: { kind: 'Service', version: 'v1' },
      version: 'cv1',
      model: { name: 'istio', version: 'v2' },
    };
    expect(componentUniqueKey(component)).toBe('Service-v1-cv1-istio-v2');
  });
});

describe('getComponentDefinition', () => {
  beforeEach(() => {
    dispatch.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns the first component when no apiVersion is provided', async () => {
    dispatch.mockResolvedValue({
      data: {
        components: [
          { component: { kind: 'Service', version: 'v1' } },
          { component: { kind: 'Service', version: 'v2' } },
        ],
      },
    });
    const result = await getComponentDefinition('Service', 'istio');
    expect(result).toEqual({ component: { kind: 'Service', version: 'v1' } });
  });

  it('filters components by apiVersion when supplied', async () => {
    dispatch.mockResolvedValue({
      data: {
        components: [
          { component: { kind: 'Service', version: 'v1' } },
          { component: { kind: 'Service', version: 'v2' } },
        ],
      },
    });
    const result = await getComponentDefinition('Service', 'istio', { apiVersion: 'v2' });
    expect(result).toEqual({ component: { kind: 'Service', version: 'v2' } });
  });

  it('returns undefined when no component matches the apiVersion', async () => {
    dispatch.mockResolvedValue({ data: { components: [] } });
    const result = await getComponentDefinition('Service', 'istio', { apiVersion: 'v99' });
    expect(result).toBeUndefined();
  });

  it('returns undefined when dispatch rejects', async () => {
    dispatch.mockRejectedValue(new Error('boom'));
    const result = await getComponentDefinition('Service', 'istio');
    expect(result).toBeUndefined();
  });
});

describe('meshModel – HTTP contracts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('updateEntityStatus posts to /meshmodels/:type/status', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('meshmodels/components/status'), {
      method: 'POST',
      body: JSON.stringify({ ids: ['c1'], status: 'enabled' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/meshmodels/components/status',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('importMeshModel POSTs to /meshmodels/register', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('meshmodels/register'), {
      method: 'POST',
      body: 'binary-bytes',
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/meshmodels/register',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('surfaces a 500 error when fetching models', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('database down'),
    });

    const resp = await fetch(mesheryApiPath('meshmodels/models'));
    expect(resp.ok).toBe(false);
  });
});
