import { describe, expect, it } from 'vitest';
import { k8sContextMatchesConnectionId } from '../k8sContext';

describe('k8sContextMatchesConnectionId', () => {
  const connectionId = 'e03d62a8-a52d-7d62-8799-0f181a5d06dd';
  const contextId = 'e03d62a8a52d7d6287990f181a5d06dd';

  it('matches on connectionId', () => {
    expect(k8sContextMatchesConnectionId({ connectionId }, connectionId)).toBe(true);
  });

  it('matches on context id without dashes', () => {
    expect(k8sContextMatchesConnectionId({ id: contextId }, connectionId)).toBe(true);
  });

  it('does not match unrelated contexts', () => {
    expect(
      k8sContextMatchesConnectionId(
        { id: 'other', connectionId: '228b65ec-320e-3930-43c7-1a29fe9a2756' },
        connectionId,
      ),
    ).toBe(false);
  });

  it('returns false when ctx or connectionId is missing', () => {
    expect(k8sContextMatchesConnectionId(undefined as never, connectionId)).toBe(false);
    expect(k8sContextMatchesConnectionId({ id: contextId }, '')).toBe(false);
  });
});
