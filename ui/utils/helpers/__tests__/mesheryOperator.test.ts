import { describe, expect, it } from 'vitest';

import { getOperatorStatusFromQueryResult, isMesheryOperatorConnected } from '../mesheryOperator';

describe('isMesheryOperatorConnected', () => {
  it('returns the operatorInstalled flag', () => {
    expect(isMesheryOperatorConnected({ operatorInstalled: true })).toBe(true);
    expect(isMesheryOperatorConnected({ operatorInstalled: false })).toBe(false);
    expect(isMesheryOperatorConnected({ operatorInstalled: undefined })).toBe(undefined);
  });
});

describe('getOperatorStatusFromQueryResult', () => {
  it('returns disconnected information when the operator entry carries an error', () => {
    const result = getOperatorStatusFromQueryResult({ operator: { error: 'boom' } });
    expect(result[0]).toBe(false);
    expect(result[1].operatorInstalled).toBe(false);
    expect(result[1].NATSInstalled).toBe(false);
    expect(result[1].meshSyncInstalled).toBe(false);
    expect(result[1].operatorVersion).toBe('N/A');
  });

  it('returns disconnected when status is not ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: { status: 'DISABLED', controllers: [] },
    });
    expect(result[0]).toBe(false);
    expect(result[1].operatorInstalled).toBe(false);
  });

  it('marks operator installed and aggregates broker/meshsync state when ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: {
        status: 'ENABLED',
        version: 'v0.7.0',
        controllers: [
          { name: 'broker', status: 'ENABLED', version: '2.10.0' },
          { name: 'meshsync', status: 'ENABLED', version: '0.7.5' },
        ],
      },
    });
    expect(result[0]).toBe(true);
    const info = result[1];
    expect(info.operatorInstalled).toBe(true);
    expect(info.operatorVersion).toBe('v0.7.0');
    expect(info.NATSInstalled).toBe(true);
    expect(info.NATSVersion).toBe('2.10.0');
    expect(info.meshSyncInstalled).toBe(true);
    expect(info.meshSyncVersion).toBe('0.7.5');
  });

  it('marks individual controllers as not installed when their status is not ENABLED', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: {
        status: 'ENABLED',
        version: 'v0.7.0',
        controllers: [
          { name: 'broker', status: 'DISABLED', version: '2.10.0' },
          { name: 'meshsync', status: 'DISABLED', version: '0.7.5' },
        ],
      },
    });
    expect(result[0]).toBe(true);
    expect(result[1].operatorInstalled).toBe(true);
    expect(result[1].NATSInstalled).toBe(false);
    expect(result[1].NATSVersion).toBe('N/A');
    expect(result[1].meshSyncInstalled).toBe(false);
    expect(result[1].meshSyncVersion).toBe('N/A');
  });

  it('handles ENABLED status with no controllers array', () => {
    const result = getOperatorStatusFromQueryResult({
      operator: { status: 'ENABLED', version: 'v0.7.0' },
    });
    expect(result[0]).toBe(true);
    expect(result[1].operatorInstalled).toBe(true);
    expect(result[1].operatorVersion).toBe('v0.7.0');
  });
});
