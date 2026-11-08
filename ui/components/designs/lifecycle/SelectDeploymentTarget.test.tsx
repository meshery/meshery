import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getEnvironmentsQuery = vi.fn();

vi.mock('@/rtk-query/environments', () => ({
  useGetEnvironmentsQuery: (...args) => getEnvironmentsQuery(...args),
  useGetEnvironmentConnectionsQuery: vi.fn(() => ({ data: undefined, isLoading: true })),
}));

vi.mock('../../layout/Header/Header', () => ({
  K8sContextConnectionChip: () => null,
}));

// ./common imports utils/utils, which reaches pages/_app and the navigator's
// JSX-in-.js icon files that the test transform cannot parse.
vi.mock('./common', () => ({
  Loading: () => <div data-testid="loading" />,
  StepHeading: ({ children }) => <h2>{children}</h2>,
}));

import { DeploymentTargetContext, SelectTargetEnvironments } from './SelectDeploymentTarget';

// SelectDeploymentTarget consumes two RTK Query hooks
// (useGetEnvironmentConnectionsQuery, useGetEnvironmentsQuery), a
// React context (DeploymentTargetContext) with three required fields,
// useSelector / useDispatch from react-redux, and the K8sContextConnectionChip
// header component. Rendering it as a leaf would require mocking the entire
// environments slice, the meshsync controller state, and a connection-
// metadata map. End-to-end behaviour is covered by playwright; the
// callback/selector wiring it depends on is exercised through the redux
// store tests.
describe.skip('SelectDeploymentTarget', () => {
  it('skipped - depends on RTK Query, redux store, and multiple contexts', () => {});
});

// SelectTargetEnvironments only reads the organization from context before it
// renders its loading state, so its environments query can be checked without
// the redux store.
describe('SelectTargetEnvironments', () => {
  beforeEach(() => {
    getEnvironmentsQuery.mockReset();
    getEnvironmentsQuery.mockReturnValue({ data: undefined, isLoading: true, isError: false });
  });

  // The server returns only its default page of 10 environments when no page
  // size is sent, so any environment past the tenth could not be chosen as a
  // deployment target.
  it('requests every environment of the organization', () => {
    render(
      <DeploymentTargetContext.Provider
        value={{
          organization: { id: 'org-1' } as any,
          meshsyncControllerState: null,
          connectionMetadataState: null,
        }}
      >
        <SelectTargetEnvironments setIsEnvrionmentModalOpen={vi.fn()} />
      </DeploymentTargetContext.Provider>,
    );

    expect(screen.getByTestId('loading')).toBeTruthy();
    expect(getEnvironmentsQuery).toHaveBeenCalledWith({ orgId: 'org-1', pagesize: 'all' });
  });
});
