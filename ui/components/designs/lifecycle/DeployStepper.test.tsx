import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

const { Dummy, mockDispatch, mockOpenView } = vi.hoisted(() => ({
  Dummy: ({ children }: any) => <div>{children}</div>,
  mockDispatch: vi.fn(),
  mockOpenView: vi.fn(),
}));

vi.mock('@sistent/sistent', () => ({
  EventBus: class {
    on = vi.fn(() => ({ subscribe: vi.fn() }));
  },
  MESHERY_EXTENSION_EVENT: { DispatchToMesheryStore: 'DispatchToMesheryStore' },
  Box: Dummy,
  Typography: Dummy,
  ModalBody: Dummy,
  ModalFooter: Dummy,
  CustomizedStepper: Dummy,
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  ModalButtonPrimary: Dummy,
  ModalButtonSecondary: Dummy,
  styled: () => () => Dummy,
  useTheme: () => ({ palette: { background: { default: '#fff' } } }),
  useStepper: () => ({ activeStep: 0, steps: [] }),
  InfoIcon: () => <svg />,
  notificationColors: { info: {}, error: {}, warning: {}, success: {} },
  EnvironmentIcon: () => <svg />,
}));

vi.mock('./ValidateDesign', () => ({}));
vi.mock('./DryRun', () => ({}));
vi.mock('./SelectDeploymentTarget', () => ({ DeploymentTargetContext: React.createContext({}) }));
vi.mock('./finalizeDeployment', () => ({}));
vi.mock('machines/validator/designValidator', () => ({}));
vi.mock('../../environments/EnvironmentFormModal', () => ({}));
vi.mock('./common', () => ({
  Loading: ({ message }: any) => <div data-testid="loading-spinner">{message}</div>,
  DEPLOYMENT_TYPE: { DEPLOY: 'deploy', UNDEPLOY: 'undeploy' },
  CheckBoxField: () => <div />,
}));
vi.mock('./DeploymentSummary', () => ({
  DeploymentSummaryFormatter: ({ event }: any) => (
    <div data-testid="deployment-summary">{event.description}</div>
  ),
}));
vi.mock('@/utils/utils', () => ({
  openViewScopedToDesignInOperator: (...args: any[]) => mockOpenView(...args),
  useIsOperatorEnabled: () => false,
}));
vi.mock('next/router', () => ({ useRouter: () => ({ push: vi.fn() }) }));
vi.mock('react-redux', () => ({ useSelector: vi.fn(), useDispatch: () => mockDispatch }));
vi.mock('@/store/slices/events', () => ({
  default: (s = {}) => s,
  openNotificationCenter: (payload: any) => ({ type: 'events/openNotificationCenter', payload }),
}));
vi.mock('@/components/layout/NotificationCenter', () => ({
  NotificationCenterContext: React.createContext({ operationsCenterActorRef: null }),
}));

import { NotificationCenterContext } from '@/components/layout/NotificationCenter';
import { FinishDeploymentStep, DEPLOYMENT_EVENT_TIMEOUT_MS } from './DeployStepper';

describe('FinishDeploymentStep', () => {
  let eventCallback: ((event: any) => void) | null = null;
  const unsubscribeMock = vi.fn();
  const mockOperationsCenterActorRef = {
    on: vi.fn((_type: string, cb: any) => {
      eventCallback = cb;
      return { unsubscribe: unsubscribeMock };
    }),
  };

  const defaultProps = {
    design: { id: 'd-1', name: 'my-design' },
    perform_deployment: vi.fn().mockResolvedValue({}),
    deployment_type: 'deploy',
    autoOpenView: false,
  };

  const emit = (e: any) =>
    act(() => eventCallback?.({ data: { event: { action: 'deploy', ...e } } }));

  const renderWithContext = (props = {}, strict = false) => {
    const el = (p: any) => (
      <NotificationCenterContext.Provider
        value={{ operationsCenterActorRef: mockOperationsCenterActorRef } as any}
      >
        <FinishDeploymentStep {...defaultProps} {...p} />
      </NotificationCenterContext.Provider>
    );
    const content = el(props);
    const res = render(strict ? <React.StrictMode>{content}</React.StrictMode> : content);
    return { ...res, rerenderStep: (p = {}) => res.rerender(el(p)) };
  };

  beforeEach(() => {
    vi.clearAllMocks();
    eventCallback = null;
  });
  afterEach(() => vi.useRealTimers());

  it('renders summary when event arrives before or after HTTP response resolves, and ignores other designs', async () => {
    let resolveHttp: any;
    const { unmount } = renderWithContext({
      perform_deployment: vi.fn(() => new Promise((res) => (resolveHttp = res))),
    });
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    emit({ actedUpon: 'd-1', description: 'Early event' });
    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('Early event');
    await act(async () => resolveHttp({}));
    expect(screen.getByTestId('deployment-summary')).toBeInTheDocument();
    unmount();

    renderWithContext();
    await act(async () => {});
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    emit({ actedUpon: 'other', description: 'Other' });
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
    emit({ metadata: { designId: 'd-1' }, description: 'Late event' });
    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('Late event');
  });

  const errData = { error: 'Failed to deploy design', code: 'meshery-server-1234' };
  it.each([
    [401, 'Auth', 'Your session has expired. Please log in again.', 'Log In', true],
    [403, 'Forbidden', 'Your session has expired. Please log in again.', 'Log In', true],
    [500, errData, 'Error 500: Failed to deploy design', 'Retry', false],
    ['PARSER', errData, 'Error 502: Failed to deploy design', 'Retry', false],
    ['NO_FIELD', { code: '1' }, 'Error 500: An unexpected error occurred.', 'Retry', false],
    [0, 'Network fail', 'Network fail', 'Retry', false],
  ])('handles rejection: %s', async (code, data, exp, btn, isLogin) => {
    const err =
      code === 'PARSER'
        ? { status: 'PARSING_ERROR', originalStatus: 502, data }
        : code === 0
          ? new Error(data as string)
          : { status: typeof code === 'number' ? code : 500, data };

    renderWithContext({ perform_deployment: vi.fn().mockRejectedValue(err) });
    expect(await screen.findByText(exp)).toBeInTheDocument();
    expect(screen.getByText(btn)).toBeInTheDocument();

    if (isLogin) {
      const orig = window.location;
      delete (window as any).location;
      window.location = { ...orig, host: 'localhost:3000', href: 'http://localhost:3000/' } as any;
      fireEvent.click(screen.getByText('Log In'));
      expect(window.location.href).toBe('/user/login');
      window.location = orig;
    }
  });

  it('handles timeout: fires after rerender with new perform_deployment, opens notification center, and late event recovers', () => {
    vi.useFakeTimers();
    const [p1, p2] = [vi.fn(() => new Promise(() => {})), vi.fn()];
    const { rerenderStep } = renderWithContext({ perform_deployment: p1 });
    rerenderStep({ perform_deployment: p2 });

    act(() => vi.advanceTimersByTime(DEPLOYMENT_EVENT_TIMEOUT_MS));
    expect(screen.getByText(/Deployment status unknown/i)).toBeInTheDocument();
    expect(p1).toHaveBeenCalledTimes(1);
    expect(p2).not.toHaveBeenCalled();

    fireEvent.click(screen.getByText('Open Notification Center'));
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'events/openNotificationCenter',
      payload: {},
    });

    emit({ actedUpon: 'd-1', description: 'Arrived after timeout' });
    expect(screen.getByTestId('deployment-summary')).toHaveTextContent('Arrived after timeout');
  });

  it('clears timer and unsubscribes on unmount without state update', () => {
    vi.useFakeTimers();
    const { unmount } = renderWithContext({
      perform_deployment: vi.fn(() => new Promise(() => {})),
    });
    unmount();
    expect(unsubscribeMock).toHaveBeenCalled();
    act(() => vi.advanceTimersByTime(DEPLOYMENT_EVENT_TIMEOUT_MS));
    expect(screen.queryByText(/Deployment status unknown/i)).not.toBeInTheDocument();
  });

  it('resets state, restarts timer and sends second request on Retry', async () => {
    let calls = 0;
    const perform = vi.fn(() =>
      calls++ === 0 ? Promise.reject({ status: 500, data: 'Failed' }) : new Promise(() => {}),
    );
    renderWithContext({ perform_deployment: perform });
    expect(await screen.findByText('Error 500: Failed')).toBeInTheDocument();

    fireEvent.click(screen.getByText('Retry'));
    expect(perform).toHaveBeenCalledTimes(2);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('shows error and Retry when request rejects after timeout has fired', async () => {
    vi.useFakeTimers();
    let rejectHttp: any;
    renderWithContext({
      perform_deployment: vi.fn(() => new Promise((_, rej) => (rejectHttp = rej))),
    });

    act(() => vi.advanceTimersByTime(DEPLOYMENT_EVENT_TIMEOUT_MS));
    expect(screen.getByText(/Deployment status unknown/i)).toBeInTheDocument();

    await act(async () => {
      rejectHttp({ status: 504, data: { error: 'Gateway Timeout' } });
    });

    expect(screen.getByText('Error 504: Gateway Timeout')).toBeInTheDocument();
    expect(screen.getByText('Retry')).toBeInTheDocument();
    expect(screen.queryByText(/Deployment status unknown/i)).not.toBeInTheDocument();
  });

  it('calls new perform_deployment on Retry after rerender with error', async () => {
    const p1 = vi.fn().mockRejectedValue({ status: 500, data: { error: 'Initial fail' } });
    const p2 = vi.fn().mockReturnValue(new Promise(() => {}));
    const { rerenderStep } = renderWithContext({ perform_deployment: p1 });

    expect(await screen.findByText('Error 500: Initial fail')).toBeInTheDocument();
    rerenderStep({ perform_deployment: p2 });

    fireEvent.click(screen.getByText('Retry'));
    expect(p2).toHaveBeenCalledTimes(1);
    expect(p1).toHaveBeenCalledTimes(1);
  });

  it('sends exactly one request under React StrictMode and fires timeout after DEPLOYMENT_EVENT_TIMEOUT_MS', () => {
    vi.useFakeTimers();
    const perform = vi.fn(() => new Promise(() => {}));
    renderWithContext({ perform_deployment: perform }, true);
    expect(perform).toHaveBeenCalledTimes(1);

    act(() => vi.advanceTimersByTime(DEPLOYMENT_EVENT_TIMEOUT_MS));
    expect(screen.getByText(/Deployment status unknown/i)).toBeInTheDocument();
  });

  it('handles auto-open: fires for deploy+Informational, skips for Warning/Error/undeploy', () => {
    const send = (severity: string, action = 'deploy') =>
      emit({ action, actedUpon: 'd-1', severity, metadata: { designName: 'd', designId: 'd-1' } });
    renderWithContext({ autoOpenView: true, deployment_type: 'deploy' });
    send('warning');
    send('error');
    expect(mockOpenView).not.toHaveBeenCalled();

    send('informational');
    expect(mockOpenView).toHaveBeenCalledWith('d', 'd-1', expect.anything());

    mockOpenView.mockClear();
    renderWithContext({ autoOpenView: true, deployment_type: 'undeploy' });
    send('success', 'undeploy');
    expect(mockOpenView).not.toHaveBeenCalled();
  });
});
