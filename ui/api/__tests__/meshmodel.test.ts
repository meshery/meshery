import { describe, expect, it, vi, beforeEach } from 'vitest';

const promisifiedDataFetch = vi.fn();

vi.mock('../../lib/data-fetch', () => ({
  promisifiedDataFetch: (...args: unknown[]) => promisifiedDataFetch(...args),
}));

vi.mock('../../constants/endpoints', () => ({
  REGISTRY_MODELS_ENDPOINT: '/api/registry/models',
  REGISTRY_ENDPOINT: '/api/registry',
}));

import {
  getMeshModels,
  getComponentFromModelApi,
  getDuplicateModels,
  getDuplicateComponents,
  getMeshModelRegistrants,
  getVersionedComponentFromModel,
  getComponentsDetail,
  getRelationshipsDetail,
  getMeshModelComponent,
  getMeshModelComponentByName,
  fetchCategories,
  getModelFromCategoryApi,
  getModelByName,
} from '../meshmodel';

describe('api/meshmodel', () => {
  beforeEach(() => {
    promisifiedDataFetch.mockReset();
    promisifiedDataFetch.mockResolvedValue({ ok: true });
  });

  describe('getMeshModels', () => {
    it('builds URL with default page, pageSize, and trim option', async () => {
      await getMeshModels();
      expect(promisifiedDataFetch).toHaveBeenCalledTimes(1);
      const url = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(url).toContain('/api/registry/models');
      expect(url).toContain('page=1');
      expect(url).toContain('pagesize=all');
      expect(url).toContain('trim=true');
    });

    it('supports custom page and pageSize', async () => {
      await getMeshModels(3, 20);
      const url = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(url).toContain('page=3');
      expect(url).toContain('pagesize=20');
    });

    it('forwards components and relationships flags through options', async () => {
      await getMeshModels(1, 'all', {
        pageSize: 'all',
        page: 2,
        trim: true,
        components: true,
        relationships: true,
      } as never);
      const url = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(url).toContain('components=true');
      expect(url).toContain('relationships=true');
      // override pages via options
      expect(url).toContain('page=2');
    });
  });

  describe('getComponentFromModelApi', () => {
    it('builds component-by-model URL with defaults', async () => {
      await getComponentFromModelApi('istio');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components?pagesize=all&trim=true',
      );
    });

    it('accepts custom pageSize and trim values', async () => {
      await getComponentFromModelApi('linkerd', 50, false);
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/linkerd/components?pagesize=50&trim=false',
      );
    });
  });

  describe('getDuplicateModels', () => {
    it('queries the model endpoint with a version query string and no trailing whitespace', async () => {
      await getDuplicateModels('istio', '1.0.0');
      const arg = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(arg).toBe('/api/registry/models/istio?version=1.0.0');
    });
  });

  describe('getDuplicateComponents', () => {
    it('encodes componentKind, apiVersion and modelName', async () => {
      await getDuplicateComponents('Pod', 'v1', 'k8s');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/components/Pod?apiVersion=v1&model=k8s',
      );
    });

    it('URL-encodes special characters in apiVersion and modelName', async () => {
      await getDuplicateComponents('Pod', 'v1/beta', 'k8s core');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/components/Pod?apiVersion=v1%2Fbeta&model=k8s+core',
      );
    });
  });

  describe('getMeshModelRegistrants', () => {
    it('uses default pagination values', async () => {
      await getMeshModelRegistrants();
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/registrants?page=1&pageSize=all',
      );
    });

    it('forwards custom pagination values', async () => {
      await getMeshModelRegistrants(7, 10);
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/registrants?page=7&pageSize=10',
      );
    });
  });

  describe('getVersionedComponentFromModel', () => {
    it('builds versioned URL with all parameters', async () => {
      await getVersionedComponentFromModel('istio', '1.18', 25, false);
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components?version=1.18&pagesize=25&trim=false',
      );
    });

    it('uses defaults for pageSize and trim', async () => {
      await getVersionedComponentFromModel('istio', '1.18');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components?version=1.18&pagesize=all&trim=true',
      );
    });
  });

  describe('paginated detail endpoints', () => {
    it('getComponentsDetail builds /api/registry/components URL', async () => {
      await getComponentsDetail(3);
      expect(promisifiedDataFetch).toHaveBeenCalledWith('/api/registry/components?page=3');
    });

    it('getRelationshipsDetail builds /api/registry/relationships URL', async () => {
      await getRelationshipsDetail(2);
      expect(promisifiedDataFetch).toHaveBeenCalledWith('/api/registry/relationships?page=2');
    });
  });

  describe('getMeshModelComponent', () => {
    it('emits no version/apiVersion query when both are absent', async () => {
      await getMeshModelComponent('istio', 'Service', undefined, undefined);
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components/Service',
      );
    });

    it('emits a leading question mark when only version is provided', async () => {
      await getMeshModelComponent('istio', 'Service', '1.0', undefined);
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components/Service?version=1.0',
      );
    });

    it('emits a leading question mark when only apiVersion is provided', async () => {
      await getMeshModelComponent('istio', 'Service', undefined, 'v1');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components/Service?apiVersion=v1',
      );
    });

    it('uses ampersand to join apiVersion when both version and apiVersion present', async () => {
      await getMeshModelComponent('istio', 'Service', '1.0', 'v1');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/models/istio/components/Service?version=1.0&apiVersion=v1',
      );
    });
  });

  describe('component-by-name and category endpoints', () => {
    it('getMeshModelComponentByName builds URL', async () => {
      await getMeshModelComponentByName('Deployment');
      expect(promisifiedDataFetch).toHaveBeenCalledWith('/api/registry/components/Deployment');
    });

    it('fetchCategories hits the categories endpoint', async () => {
      await fetchCategories();
      expect(promisifiedDataFetch).toHaveBeenCalledWith('/api/registry/categories');
    });

    it('getModelFromCategoryApi builds URL with all pagesize', async () => {
      await getModelFromCategoryApi('Orchestration');
      expect(promisifiedDataFetch).toHaveBeenCalledWith(
        '/api/registry/categories/Orchestration/models?pagesize=all',
      );
    });
  });

  describe('getModelByName', () => {
    it('uses default options when none provided', async () => {
      await getModelByName('istio');
      const url = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(url).toContain('/api/registry/models/istio');
      expect(url).toContain('trim=true');
      expect(url).toContain('pagesize=all');
    });

    it('forwards custom options', async () => {
      await getModelByName('istio', {
        pageSize: '5',
        page: 2,
        trim: false,
        components: true,
        relationships: false,
      } as never);
      const url = promisifiedDataFetch.mock.calls[0][0] as string;
      expect(url).toContain('pagesize=5');
      expect(url).toContain('page=2');
      // trim is false -> not appended
      expect(url).not.toContain('trim=');
      expect(url).toContain('components=true');
      // relationships is false -> not appended
      expect(url).not.toContain('relationships=');
    });

    it('returns whatever promisifiedDataFetch resolves with', async () => {
      promisifiedDataFetch.mockResolvedValueOnce({ name: 'istio' });
      await expect(getModelByName('istio')).resolves.toEqual({ name: 'istio' });
    });
  });
});
