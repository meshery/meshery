import { describe, expect, it } from 'vitest';
import {
  getDesignEventErrorSummary,
  getDesignEventIdentity,
  isMeshkitErrorCode,
  meshkitErrorDocLink,
} from './design_event_metadata';

describe('getDesignEventIdentity', () => {
  it('reads the camelCase identity keys', () => {
    expect(
      getDesignEventIdentity({
        designId: 'd1',
        designName: 'My Design',
        designVersion: '1.0.0',
        designLink: 'http://meshery/design/d1',
        viewLink: 'http://meshery/view/d1',
      }),
    ).toEqual({
      designId: 'd1',
      designName: 'My Design',
      designVersion: '1.0.0',
      designLink: 'http://meshery/design/d1',
      viewLink: 'http://meshery/view/d1',
    });
  });

  it('falls back to the legacy snake_case keys when the camelCase key is absent', () => {
    expect(
      getDesignEventIdentity({
        design_id: 'legacy-id',
        design_name: 'Legacy Design',
        view_link: 'http://meshery/view/legacy-id',
      }),
    ).toEqual({
      designId: 'legacy-id',
      designName: 'Legacy Design',
      designVersion: undefined,
      designLink: undefined,
      viewLink: 'http://meshery/view/legacy-id',
    });
  });

  it('prefers the camelCase key when both spellings are present', () => {
    const identity = getDesignEventIdentity({
      designId: 'new-id',
      design_id: 'old-id',
      designName: 'New Name',
      design_name: 'Old Name',
      viewLink: 'http://meshery/view/new-id',
      view_link: 'http://meshery/view/old-id',
    });
    expect(identity.designId).toBe('new-id');
    expect(identity.designName).toBe('New Name');
    expect(identity.viewLink).toBe('http://meshery/view/new-id');
  });

  it('treats an empty or null camelCase value as absent', () => {
    const identity = getDesignEventIdentity({
      designName: '',
      design_name: 'Legacy Design',
      designId: null,
      design_id: 'legacy-id',
    });
    expect(identity.designName).toBe('Legacy Design');
    expect(identity.designId).toBe('legacy-id');
  });

  it('tolerates missing metadata', () => {
    expect(getDesignEventIdentity(undefined)).toEqual({
      designId: undefined,
      designName: undefined,
      designVersion: undefined,
      designLink: undefined,
      viewLink: undefined,
    });
  });
});

describe('getDesignEventErrorSummary', () => {
  it('returns null for an event without status or errors', () => {
    expect(getDesignEventErrorSummary({ designId: 'd1' })).toBeNull();
    expect(getDesignEventErrorSummary(null)).toBeNull();
  });

  it('reads status, errorCount, errors and truncated', () => {
    expect(
      getDesignEventErrorSummary({
        status: 'partial',
        errorCount: 25,
        errors: [{ code: 'meshery-server-1480', componentId: 'c1' }],
        truncated: true,
      }),
    ).toEqual({
      status: 'partial',
      errorCount: 25,
      errors: [{ code: 'meshery-server-1480', componentId: 'c1' }],
      truncated: true,
    });
  });

  it('derives errorCount from errors when it is missing and drops malformed entries', () => {
    expect(
      getDesignEventErrorSummary({
        status: 'failed',
        errorCount: null,
        errors: [{ code: 'required' }, null, 'not-an-error'],
      }),
    ).toEqual({
      status: 'failed',
      errorCount: 1,
      errors: [{ code: 'required', componentId: undefined }],
      truncated: false,
    });
  });
});

describe('isMeshkitErrorCode', () => {
  it('accepts MeshKit codes only', () => {
    expect(isMeshkitErrorCode('meshery-server-1480')).toBe(true);
    expect(isMeshkitErrorCode('Meshery-Server-1480')).toBe(true);
    expect(isMeshkitErrorCode('meshkit-11001')).toBe(true);
    expect(isMeshkitErrorCode('meshkit-')).toBe(false);
    expect(isMeshkitErrorCode('meshkit-core-11001')).toBe(false);
    expect(isMeshkitErrorCode('required')).toBe(false);
    expect(isMeshkitErrorCode('type')).toBe(false);
    expect(isMeshkitErrorCode('FieldValueInvalid')).toBe(false);
    expect(isMeshkitErrorCode('meshery-server-')).toBe(false);
    expect(isMeshkitErrorCode(undefined)).toBe(false);
  });

  it('links a MeshKit code to its error code reference entry', () => {
    expect(meshkitErrorDocLink('meshery-server-1480')).toBe(
      'https://docs.meshery.io/reference/error-codes#meshery-server-1480',
    );
  });
});
