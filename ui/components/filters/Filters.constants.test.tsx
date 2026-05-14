import { describe, expect, it } from 'vitest';
import { ACTION_TYPES, COLUMN_VIEWS } from './Filters.constants';

describe('Filters.constants', () => {
  it('exposes ACTION_TYPES with name and error_msg for each entry', () => {
    expect(ACTION_TYPES.FETCH_FILTERS).toEqual({
      name: 'FETCH_FILTERS',
      error_msg: 'Failed to fetch filter',
    });
    expect(ACTION_TYPES.DELETE_FILTERS).toEqual({
      name: 'DELETE_FILTERS',
      error_msg: 'Failed to delete filter file',
    });
    expect(ACTION_TYPES.DEPLOY_FILTERS).toEqual({
      name: 'DEPLOY_FILTERS',
      error_msg: 'Failed to deploy filter file',
    });
    expect(ACTION_TYPES.UNDEPLOY_FILTERS).toEqual({
      name: 'UNDEPLOY_FILTERS',
      error_msg: 'Failed to undeploy filter file',
    });
    expect(ACTION_TYPES.UPLOAD_FILTERS).toEqual({
      name: 'UPLOAD_FILTERS',
      error_msg: 'Failed to upload filter file',
    });
    expect(ACTION_TYPES.CLONE_FILTERS).toEqual({
      name: 'CLONE_FILTER',
      error_msg: 'Failed to clone filter file',
    });
    expect(ACTION_TYPES.PUBLISH_CATALOG.error_msg).toBe('Failed to publish catalog');
    expect(ACTION_TYPES.UNPUBLISH_CATALOG.error_msg).toBe('Failed to publish catalog');
    expect(ACTION_TYPES.SCHEMA_FETCH).toEqual({
      name: 'SCHEMA_FETCH',
      error_msg: 'failed to fetch import schema',
    });
  });

  it('exposes COLUMN_VIEWS in [name, size] tuples', () => {
    expect(COLUMN_VIEWS).toEqual([
      ['name', 'xs'],
      ['created_at', 'm'],
      ['updated_at', 'l'],
      ['visibility', 's'],
      ['Actions', 'xs'],
    ]);
  });
});
