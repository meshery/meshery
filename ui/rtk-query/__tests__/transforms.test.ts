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

  it('can normalize collection items alongside pagination metadata', () => {
    expect(
      normalizePaginatedCollectionResponse(
        {
          total_count: 1,
          designs: [
            {
              id: 'design-1',
              user_id: 'user-1',
              catalog_data: { type: 'pattern' },
              pattern_file: 'apiVersion: v1',
              created_at: '2026-05-08T00:00:00Z',
              updated_at: '2026-05-08T01:00:00Z',
            },
          ],
        },
        'designs',
        (design) => ({
          ...design,
          userId: design.userId ?? design.user_id,
          catalogData: design.catalogData ?? design.catalog_data,
          patternFile: design.patternFile ?? design.pattern_file,
          createdAt: design.createdAt ?? design.created_at,
          updatedAt: design.updatedAt ?? design.updated_at,
        }),
      ),
    ).toEqual(
      expect.objectContaining({
        totalCount: 1,
        designs: [
          expect.objectContaining({
            id: 'design-1',
            userId: 'user-1',
            catalogData: { type: 'pattern' },
            patternFile: 'apiVersion: v1',
            createdAt: '2026-05-08T00:00:00Z',
            updatedAt: '2026-05-08T01:00:00Z',
          }),
        ],
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
