import { describe, expect, it } from 'vitest';
import {
  normalizeKubernetesContextsResponse,
  normalizePaginatedCollectionResponse,
  normalizeProviderCapabilities,
} from '../transforms';

describe('normalizeProviderCapabilities', () => {
  it('normalizes snake_case provider fields to camelCase', () => {
    expect(
      normalizeProviderCapabilities({
        provider_name: 'Meshery Cloud',
        provider_type: 'remote',
        provider_url: 'https://cloud.meshery.io',
        provider_description: ['Hosted provider'],
        capabilities: [],
      }),
    ).toEqual(
      expect.objectContaining({
        providerName: 'Meshery Cloud',
        providerType: 'remote',
        providerUrl: 'https://cloud.meshery.io',
        providerDescription: ['Hosted provider'],
      }),
    );
  });
});

describe('normalizePaginatedCollectionResponse', () => {
  it('normalizes total_count and page_size to camelCase', () => {
    expect(
      normalizePaginatedCollectionResponse(
        {
          page: 0,
          page_size: 10,
          total_count: 42,
          resources: [{ id: '1' }],
        },
        'resources',
      ),
    ).toEqual(
      expect.objectContaining({
        page: 0,
        pageSize: 10,
        totalCount: 42,
        resources: [{ id: '1' }],
      }),
    );
  });
});

describe('normalizeKubernetesContextsResponse', () => {
  it('normalizes context snake_case fields to camelCase', () => {
    expect(
      normalizeKubernetesContextsResponse({
        total_count: 1,
        contexts: [
          {
            id: 'ctx-1',
            connection_id: 'conn-1',
            is_current_context: true,
            created_by: 'meshery',
            deployment_type: 'in_cluster',
          },
        ],
      }),
    ).toEqual(
      expect.objectContaining({
        totalCount: 1,
        contexts: [
          expect.objectContaining({
            id: 'ctx-1',
            connectionId: 'conn-1',
            isCurrentContext: true,
            createdBy: 'meshery',
            deploymentType: 'in_cluster',
          }),
        ],
      }),
    );
  });
});
