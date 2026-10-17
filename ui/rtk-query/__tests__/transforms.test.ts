import { describe, expect, it } from 'vitest';
import {
  normalizeKubernetesContextsResponse,
  normalizeLoggedInUser,
  normalizePaginatedCollectionResponse,
  normalizeProviderCapabilities,
} from '../transforms';

describe('normalizeLoggedInUser', () => {
  it('backfills userId from id when the v1beta3 response omits userId', () => {
    expect(normalizeLoggedInUser({ id: 'uuid-1', email: 'owner@meshery.io' })).toEqual({
      id: 'uuid-1',
      email: 'owner@meshery.io',
      userId: 'uuid-1',
    });
  });

  it('preserves an existing userId (does not overwrite it with id)', () => {
    expect(
      normalizeLoggedInUser({ id: 'uuid-1', userId: 'legacy-user', status: 'active' }),
    ).toEqual({
      id: 'uuid-1',
      userId: 'legacy-user',
      status: 'active',
    });
  });

  it('returns undefined for null / non-object responses', () => {
    expect(normalizeLoggedInUser(undefined)).toBeUndefined();
    expect(normalizeLoggedInUser(null as unknown as undefined)).toBeUndefined();
    expect(normalizeLoggedInUser('not-an-object' as unknown as undefined)).toBeUndefined();
  });
});

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

  it('normalizes non-array capabilities to an empty array', () => {
    expect(
      normalizeProviderCapabilities({
        provider_name: 'Meshery Cloud',
        capabilities: { feature: 'persist-meshery-patterns' },
      }),
    ).toEqual(
      expect.objectContaining({
        providerName: 'Meshery Cloud',
        capabilities: [],
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
