import { beforeEach, describe, expect, it, vi } from 'vitest';

const commitMutationMock = vi.fn();

vi.mock('react-relay', () => ({
  graphql: () => ({}),
  commitMutation: (...args: any[]) => commitMutationMock(...args),
}));

vi.mock('lib/relayEnvironment', () => ({
  createRelayEnvironment: () => ({}),
}));

import changeAdapterState from './AdapterStatusMutation';

const variables = { status: 'ENABLED', targetPort: 'localhost:10000', adapter: 'meshery-istio' };

describe('changeAdapterState', () => {
  beforeEach(() => {
    commitMutationMock.mockReset();
  });

  it('passes the response and errors from onCompleted to the callback', () => {
    const onComplete = vi.fn();
    changeAdapterState(onComplete, variables);

    const config = commitMutationMock.mock.calls[0][1];
    config.onCompleted({ adapterStatus: 'ENABLED' }, null);

    expect(onComplete).toHaveBeenCalledWith({ adapterStatus: 'ENABLED' }, null);
  });

  it('reports a failed request to the callback instead of swallowing it', () => {
    const onComplete = vi.fn();
    changeAdapterState(onComplete, variables);

    const error = new Error('network down');
    commitMutationMock.mock.calls[0][1].onError(error);

    expect(onComplete).toHaveBeenCalledWith(null, [error]);
  });
});
