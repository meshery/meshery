import { describe, expect, it } from 'vitest';
import { describeConsoleError } from '../errors';

describe('describeConsoleError', () => {
  it('explains a 404 as a stale MeshSync view, not a generic failure', () => {
    // This is the case a user actually hits: Kanvas draws a pod from MeshSync's
    // cache, the cluster deleted it, and the capabilities call 404s.
    const described = describeConsoleError(
      {
        status: 404,
        meshkit: { message: 'Console target not found', code: 'meshery-server-1445' },
      },
      'meshery-nats-0',
    );

    expect(described.message).toContain('meshery-nats-0');
    expect(described.message).toContain('no longer exists');
    expect(described.detail).toContain('MeshSync');
    expect(described.code).toBe('meshery-server-1445');
  });

  it('surfaces the server’s own description and remediation', () => {
    const described = describeConsoleError({
      status: 400,
      meshkit: {
        message: 'Console kind is not supported for this target',
        code: 'meshery-server-1444',
        suggestedRemediation: ['Request a console kind the target supports.'],
        probableCause: ['The pod is not running.'],
      },
    });

    expect(described.message).toBe('Console kind is not supported for this target');
    expect(described.detail).toBe('Request a console kind the target supports.');
    expect(described.code).toBe('meshery-server-1444');
  });

  it('falls back to the probable cause when no remediation is offered', () => {
    const described = describeConsoleError({
      status: 500,
      meshkit: {
        message: 'Failed to open logs console',
        probableCause: ['RBAC forbids pods/log.'],
      },
    });

    expect(described.detail).toBe('RBAC forbids pods/log.');
  });

  it('names a transport failure as such', () => {
    expect(describeConsoleError({ status: 'FETCH_ERROR' }).message).toContain(
      'Could not reach the Meshery server',
    );
  });

  it('degrades to a generic message when the server said nothing useful', () => {
    expect(describeConsoleError(undefined).message).toBe(
      'Could not determine what this resource supports.',
    );
    expect(describeConsoleError({ status: 500 }).message).toBe(
      'Could not determine what this resource supports.',
    );
  });
});
