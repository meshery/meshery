import React from 'react';
import { act, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const listeners: Array<(event: any) => void> = [];

vi.mock('@sistent/sistent', () => ({
  ModalBody: ({ children }: any) => <div>{children}</div>,
  ModalFooter: ({ children }: any) => <div>{children}</div>,
  useStepper: () => ({}),
  CustomizedStepper: ({ children }: any) => <div>{children}</div>,
  ModalButtonPrimary: ({ children }: any) => <button>{children}</button>,
  ModalButtonSecondary: ({ children }: any) => <button>{children}</button>,
  Box: ({ children }: any) => <div>{children}</div>,
  Typography: ({ children }: any) => <span data-testid="deploy-error">{children}</span>,
  styled: (Component: any) => () => Component,
  useTheme: () => ({ palette: {} }),
  EnvironmentIcon: () => null,
}));

vi.mock('./ValidateDesign', () => ({ ValidateDesign: () => null }));
vi.mock('./DryRun', () => ({
  DryRunDesign: () => null,
  getTotalCountOfDeploymentErrors: () => 0,
}));
vi.mock('./common', () => ({
  CheckBoxField: () => null,
  DEPLOYMENT_TYPE: { DEPLOY: 'deploy', UNDEPLOY: 'undeploy' },
  Loading: ({ message }: any) => <div data-testid="loading">{message}</div>,
}));
vi.mock('@/assets/icons/DryRunIcon', () => ({ default: () => null }));
vi.mock('@/assets/icons/DeploymentSelectorIcon', () => ({ DeploymentSelectorIcon: () => null }));
vi.mock('@/assets/icons/CheckIcon', () => ({ default: () => null }));
vi.mock('@/assets/icons/FinishFlagIcon', () => ({ default: () => null }));
vi.mock('./SelectDeploymentTarget', () => ({
  DeploymentTargetContext: (require('react') as typeof import('react')).createContext({}),
  SelectTargetEnvironments: () => null,
}));
vi.mock('./finalizeDeployment', () => ({ FinalizeDeployment: () => null }));
vi.mock('@/store/slices/globalEnvironmentContext', () => ({
  selectAllSelectedK8sConnections: () => [],
}));
vi.mock('machines/validator/designValidator', () => ({
  useDryRunValidationResults: () => ({}),
  useIsValidatingDryRun: () => false,
}));
vi.mock('react-redux', () => ({ useSelector: () => ({}) }));
vi.mock('../../layout/NotificationCenter', () => {
  const ctx = (require('react') as typeof import('react')).createContext({
    operationsCenterActorRef: {
      on: (_type: any, listener: any) => {
        listeners.push(listener);
        return { unsubscribe: () => {} };
      },
    },
  });
  return { NotificationCenterContext: ctx };
});
vi.mock('machines/operationsCenter', () => ({
  OPERATION_CENTER_EVENTS: { EVENT_RECEIVED_FROM_SERVER: 'event-received-from-server' },
}));
vi.mock('./DeploymentSummary', () => ({
  DeploymentSummaryFormatter: ({ event }: any) => (
    <div data-testid="deployment-summary">{event.description}</div>
  ),
}));
vi.mock('../../layout/NotificationCenter/constants', () => ({
  SEVERITY: { SUCCESS: 'success' },
}));
vi.mock('../../environments/EnvironmentFormModal', () => ({ default: () => null }));
vi.mock('@/utils/utils', () => ({ openViewScopedToDesignInOperator: vi.fn() }));
vi.mock('next/router', () => ({ useRouter: () => ({}) }));
vi.mock('@/store/ProviderStoreWrapper', () => ({
  default: ({ children }: any) => <>{children}</>,
}));

import { DEPLOYMENT_RESULT_TIMEOUT_MS, FinishDeploymentStep } from './DeployStepper';

const design = { id: 'design-1', name: 'my design' };

const receive = (event: any) =>
  act(() => {
    listeners.forEach((listener) => listener({ data: { event } }));
  });

const renderStep = (performDeployment: () => any) =>
  render(
    <FinishDeploymentStep
      perform_deployment={performDeployment}
      deployment_type="deploy"
      autoOpenView={false}
      design={design}
    />,
  );

// DeployStepper builds a sistent useStepper-driven multi-step modal that
// wires Redux selectors (selectAllSelectedK8sConnections), xstate actor
// hooks (useDryRunValidationResults / useIsValidatingDryRun), several
// child containers (SelectDeploymentTarget, DryRunDesign,
// FinalizeDeployment, ValidateDesign), the NotificationCenter context,
// and a router. End-to-end deploy/undeploy behavior is covered by
// playwright. The DEPLOYMENT_TYPE constants and individual sub-steps
// (DryRun summary, ValidateDesign, common helpers) are unit-tested
// separately.
describe.skip('DeployStepper / UpdateDeploymentStepper', () => {
  it('skipped - multi-step stepper bound to xstate + redux + router', () => {});
});

describe('FinishDeploymentStep', () => {
  beforeEach(() => {
    listeners.length = 0;
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('performs the deployment once and waits for its result', async () => {
    const performDeployment = vi.fn().mockResolvedValue({ data: {} });

    renderStep(performDeployment);
    await act(async () => {});

    expect(performDeployment).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('loading')).toHaveTextContent('Deploying design');
  });

  it('shows the summary of the event raised for this design', async () => {
    renderStep(vi.fn().mockResolvedValue({ data: {} }));
    await act(async () => {});

    receive({ action: 'deploy', actedUpon: design.id, description: 'Deployment completed' });

    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('Deployment completed');
  });

  it('ignores events raised for another design or another action', async () => {
    renderStep(vi.fn().mockResolvedValue({ data: {} }));
    await act(async () => {});

    receive({ action: 'deploy', actedUpon: 'another-design', description: 'other design' });
    receive({ action: 'undeploy', actedUpon: design.id, description: 'other action' });

    expect(screen.queryByTestId('deployment-summary')).not.toBeInTheDocument();
    expect(screen.getByTestId('loading')).toBeInTheDocument();
  });

  it('surfaces a request that failed before the server raised an event', async () => {
    renderStep(vi.fn().mockResolvedValue({ error: { status: 401, data: 'session expired' } }));
    await act(async () => {});

    expect(screen.queryByTestId('loading')).not.toBeInTheDocument();
    expect(screen.getByTestId('deploy-error')).toHaveTextContent('session expired');
  });

  it('surfaces a deployment that could not be started', async () => {
    renderStep(
      vi.fn(() => {
        throw new Error('no connection selected');
      }),
    );
    await act(async () => {});

    expect(screen.getByTestId('deploy-error')).toHaveTextContent('no connection selected');
  });

  it('prefers the event over a failed request, since it carries the summary', async () => {
    renderStep(vi.fn().mockResolvedValue({ error: { status: 500 } }));
    await act(async () => {});

    receive({ action: 'deploy', actedUpon: design.id, description: 'Failed to deploy design' });

    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('Failed to deploy design');
  });

  it('stops waiting when no result arrives in time', async () => {
    vi.useFakeTimers();
    renderStep(() => new Promise(() => {}));

    await act(async () => {
      vi.advanceTimersByTime(DEPLOYMENT_RESULT_TIMEOUT_MS);
    });

    expect(screen.queryByTestId('loading')).not.toBeInTheDocument();
    expect(screen.getByTestId('deploy-error')).toHaveTextContent('No deploy result was received');
  });

  it('does not time out once the result has arrived', async () => {
    vi.useFakeTimers();
    renderStep(() => new Promise(() => {}));

    receive({ action: 'deploy', actedUpon: design.id, description: 'Deployment completed' });
    await act(async () => {
      vi.advanceTimersByTime(DEPLOYMENT_RESULT_TIMEOUT_MS);
    });

    expect(screen.getByTestId('deployment-summary')).toBeInTheDocument();
  });
});
