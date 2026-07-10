import { describe, expect, it } from 'vitest';
import { describeSessionError } from '../errors';

describe('describeSessionError', () => {
  it('explains a 404 as a stale MeshSync view, not a generic failure', () => {
    // This is the case a user actually hits: Kanvas draws a pod from MeshSync's
    // cache, the cluster deleted it, and the capabilities call 404s.
    const described = describeSessionError(
      {
        status: 404,
        meshkit: { message: 'Session target not found', code: 'meshery-server-1445' },
      },
      'meshery-nats-0',
    );

    expect(described.message).toContain('meshery-nats-0');
    expect(described.message).toContain('no longer exists');
    expect(described.detail).toContain('MeshSync');
    expect(described.code).toBe('meshery-server-1445');
  });

  it('surfaces the server’s own description and remediation', () => {
    const described = describeSessionError({
      status: 400,
      meshkit: {
        message: 'Session kind is not supported for this target',
        code: 'meshery-server-1444',
        suggestedRemediation: ['Request a session kind the target supports.'],
        probableCause: ['The pod is not running.'],
      },
    });

    expect(described.message).toBe('Session kind is not supported for this target');
    expect(described.detail).toBe('Request a session kind the target supports.');
    expect(described.code).toBe('meshery-server-1444');
  });

  it('falls back to the probable cause when no remediation is offered', () => {
    const described = describeSessionError({
      status: 500,
      meshkit: {
        message: 'Failed to open logs session',
        probableCause: ['RBAC forbids pods/log.'],
      },
    });

    expect(described.detail).toBe('RBAC forbids pods/log.');
  });

  it('names a transport failure as such', () => {
    expect(describeSessionError({ status: 'FETCH_ERROR' }).message).toContain(
      'Could not reach the Meshery server',
    );
  });

  it('degrades to a generic message when the server said nothing useful', () => {
    expect(describeSessionError(undefined).message).toBe(
      'Could not determine what this resource supports.',
    );
    expect(describeSessionError({ status: 500 }).message).toBe(
      'Could not determine what this resource supports.',
    );
  });
});
