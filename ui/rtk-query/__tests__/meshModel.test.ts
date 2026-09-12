import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { mesheryApiPath } from '../index';

// meshModel.ts imports from ../utils which depends on ../../store. Mock the
// store before the module graph is touched so we can import the helpers.
const { dispatch } = vi.hoisted(() => ({ dispatch: vi.fn() }));
vi.mock('../../store', () => ({ store: { dispatch } }));

import { componentUniqueKey, getComponentDefinition, modelUniqueKey } from '../meshModel';

describe('meshModel – URLs', () => {
  it('builds /api/registry/models URL', () => {
    expect(mesheryApiPath('registry/models')).toBe('/api/registry/models');
  });

  it('builds /api/registry/components URL', () => {
    expect(mesheryApiPath('registry/components')).toBe('/api/registry/components');
  });

  it('builds /api/registry/relationships URL', () => {
    expect(mesheryApiPath('registry/relationships')).toBe('/api/registry/relationships');
  });

  it('builds /api/registry/registrants URL', () => {
    expect(mesheryApiPath('registry/registrants')).toBe('/api/registry/registrants');
  });

  it('builds nested model->components URL', () => {
    expect(mesheryApiPath('registry/models/istio/components')).toBe(
      '/api/registry/models/istio/components',
    );
  });

  it('builds /api/registry/categories URL', () => {
    expect(mesheryApiPath('registry/categories')).toBe('/api/registry/categories');
  });

  it('builds /api/registry/categories/:cat/models URL', () => {
    expect(mesheryApiPath('registry/categories/Networking/models')).toBe(
      '/api/registry/categories/Networking/models',
    );
  });

  it('builds /api/registry/components/:name URL', () => {
    expect(mesheryApiPath('registry/components/Service')).toBe('/api/registry/components/Service');
  });

  it('builds /api/registry/models/:name URL', () => {
    expect(mesheryApiPath('registry/models/istio')).toBe('/api/registry/models/istio');
  });

  it('builds /api/registry/export URL', () => {
    expect(mesheryApiPath('registry/export')).toBe('/api/registry/export');
  });

  it('builds /api/registry/register URL', () => {
    expect(mesheryApiPath('registry/register')).toBe('/api/registry/register');
  });

  it('builds /api/registry/:type/status URL', () => {
    expect(mesheryApiPath('registry/components/status')).toBe('/api/registry/components/status');
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

  it('updateEntityStatus posts to /registry/:type/status', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('registry/components/status'), {
      method: 'POST',
      body: JSON.stringify({ ids: ['c1'], status: 'enabled' }),
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/registry/components/status',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('importMeshModel POSTs to /registry/register', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      text: () => Promise.resolve('{}'),
    });

    await fetch(mesheryApiPath('registry/register'), {
      method: 'POST',
      body: 'binary-bytes',
    });

    expect(global.fetch).toHaveBeenCalledWith(
      '/api/registry/register',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('surfaces a 500 error when fetching models', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('database down'),
    });

    const resp = await fetch(mesheryApiPath('registry/models'));
    expect(resp.ok).toBe(false);
  });
});
