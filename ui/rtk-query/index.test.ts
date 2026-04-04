import { describe, expect, it } from 'vitest';
import { normalizeMesheryRequest } from './index';

describe('normalizeMesheryRequest', () => {
  it('prefixes relative Meshery paths with /api', () => {
    expect(normalizeMesheryRequest('system/sync')).toBe('/api/system/sync');
    expect(normalizeMesheryRequest({ url: 'workspaces?page=0', method: 'GET' })).toEqual({
      url: '/api/workspaces?page=0',
      method: 'GET',
    });
  });

  it('keeps already-qualified and extension paths unchanged', () => {
    expect(normalizeMesheryRequest('/api/system/sync')).toBe('/api/system/sync');
    expect(normalizeMesheryRequest('/evaluate')).toBe('/evaluate');
    expect(normalizeMesheryRequest('extensions/api/workspaces/1/views')).toBe(
      'extensions/api/workspaces/1/views',
    );
    expect(normalizeMesheryRequest('https://example.com/system/sync')).toBe(
      'https://example.com/system/sync',
    );
  });
});
