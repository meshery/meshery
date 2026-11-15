import React from 'react';
import { render, screen } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getEnvironmentsQuery = vi.fn();
const getEnvironmentConnectionsQuery = vi.fn();

vi.mock('@/rtk-query/environments', () => ({
  useGetEnvironmentsQuery: (...args) => getEnvironmentsQuery(...args),
  useGetEnvironmentConnectionsQuery: (...args) => getEnvironmentConnectionsQuery(...args),
}));

vi.mock('../../layout/Header/Header', () => ({
  K8sContextConnectionChip: ({ ctx }) => <div data-testid={`k8s-${ctx.connectionId}`} />,
}));

// ./common imports utils/utils, which reaches pages/_app and the navigator's
// JSX-in-.js icon files that the test transform cannot parse.
vi.mock('./common', () => ({
  Loading: () => <div data-testid="loading" />,
  StepHeading: ({ children }) => <h2>{children}</h2>,
}));

// Each EnvironmentCard reads its selection state from the store.
vi.mock('react-redux', () => ({
  useSelector: (selector) => selector({}),
  useDispatch: () => vi.fn(),
}));

vi.mock('@/store/slices/globalEnvironmentContext', () => ({
  selectIsEnvSelected: () => false,
  selectIsConnectionSelected: () => false,
  selectSelectedK8sConnections: () => [],
  toggleEnvSelection: vi.fn(),
  toggleConnection: vi.fn(),
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
const renderSelectTargetEnvironments = () =>
  render(
    <SistentThemeProvider initialMode="light">
      <DeploymentTargetContext.Provider
        value={{
          organization: { id: 'org-1' } as any,
          meshsyncControllerState: null,
          connectionMetadataState: null,
        }}
      >
        <SelectTargetEnvironments setIsEnvrionmentModalOpen={vi.fn()} />
      </DeploymentTargetContext.Provider>
    </SistentThemeProvider>,
  );

describe('SelectTargetEnvironments', () => {
  beforeEach(() => {
    getEnvironmentsQuery.mockReset();
    getEnvironmentsQuery.mockReturnValue({ data: undefined, isLoading: true, isError: false });
    getEnvironmentConnectionsQuery.mockReset();
    getEnvironmentConnectionsQuery.mockReturnValue({ data: undefined, isLoading: true });
  });

  // The server returns only its default page of 10 environments when no page
  // size is sent, so any environment past the tenth could not be chosen as a
  // deployment target.
  it('requests every environment of the organization', () => {
    renderSelectTargetEnvironments();

    expect(screen.getByTestId('loading')).toBeTruthy();
    expect(getEnvironmentsQuery).toHaveBeenCalledWith({ orgId: 'org-1', pagesize: 'all' });
  });

  it('offers every environment the query returns as a deployment target, past the tenth', () => {
    const names = Array.from({ length: 11 }, (_, i) => `staging-${String(i + 1).padStart(2, '0')}`);
    getEnvironmentsQuery.mockReturnValue({
      data: { environments: names.map((name) => ({ id: `id-${name}`, name })) },
      isLoading: false,
      isError: false,
    });

    renderSelectTargetEnvironments();

    for (const name of names) {
      expect(screen.getByTestId(`env-${name}`)).toBeTruthy();
    }
  });

  // Each environment card lists that environment's Kubernetes connections. With
  // no page size the server returns only the 10 most recently updated
  // connections of any kind, and the card keeps the Kubernetes ones from that
  // page, so older clusters could not be chosen as deployment targets.
  it("requests every connection of each environment, not only the server's first page", () => {
    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [{ id: 'env-1', name: 'prod' }] },
      isLoading: false,
      isError: false,
    });

    renderSelectTargetEnvironments();

    expect(getEnvironmentConnectionsQuery).toHaveBeenCalledWith({
      environmentId: 'env-1',
      pagesize: 'all',
    });
  });

  it('offers every Kubernetes connection of an environment, past the tenth connection', () => {
    getEnvironmentsQuery.mockReturnValue({
      data: { environments: [{ id: 'env-1', name: 'prod' }] },
      isLoading: false,
      isError: false,
    });
    const clusters = Array.from({ length: 11 }, (_, i) => ({
      id: `k8s-${i + 1}`,
      kind: 'kubernetes',
      metadata: {},
    }));
    getEnvironmentConnectionsQuery.mockReturnValue({
      data: { connections: [...clusters, { id: 'prom-1', kind: 'prometheus', metadata: {} }] },
      isLoading: false,
    });

    renderSelectTargetEnvironments();

    for (const { id } of clusters) {
      expect(screen.getByTestId(`k8s-${id}`)).toBeTruthy();
    }
    expect(screen.queryByTestId('k8s-prom-1')).toBeNull();
    expect(screen.getByText('(0/11)')).toBeTruthy();
  });
});
