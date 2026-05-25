import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const fetchSystemSyncMock = vi.fn();
const unwrapMock = vi.fn();
const notifyMock = vi.fn();
const changeAdapterStateMock = vi.fn();
const updateProgressMock = vi.fn();

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children }: any) => <div data-testid="grid">{children}</div>,
  Switch: ({ checked, onChange, name }: any) => (
    <input data-testid={`switch-${name}`} type="checkbox" checked={!!checked} onChange={onChange} />
  ),
  Typography: ({ children }: any) => <div>{children}</div>,
  useTheme: () => ({
    palette: { text: { brand: '#123456' } },
  }),
}));

vi.mock('../../../css/icons.styles', () => ({
  CardContainer: ({ children }: any) => <div data-testid="card-container">{children}</div>,
  FrontSideDescription: ({ children }: any) => <div data-testid="description">{children}</div>,
  ImageWrapper: ({ src }: any) => <img data-testid="image" alt="adapter" src={src} />,
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../../css/grid.style', () => ({
  LARGE_6_MED_12_GRID_STYLE: { xs: 12 },
}));

vi.mock('@/graphql/mutations/AdapterStatusMutation', () => ({
  default: (...args: any[]) => changeAdapterStateMock(...args),
}));

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: notifyMock }),
}));

vi.mock('../../../lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

vi.mock('../../../rtk-query/system', () => ({
  useLazyGetSystemSyncQuery: () => [fetchSystemSyncMock],
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: (...args: any[]) => updateProgressMock(...args),
}));

import Adapters from './adapters';

describe('Adapters component', () => {
  beforeEach(() => {
    fetchSystemSyncMock.mockReset();
    unwrapMock.mockReset();
    notifyMock.mockReset();
    changeAdapterStateMock.mockReset();
    updateProgressMock.mockReset();

    fetchSystemSyncMock.mockReturnValue({
      unwrap: () =>
        Promise.resolve({
          meshAdapters: [],
        }),
    });
  });

  it('renders adapter cards for each adapter in the list', async () => {
    await act(async () => {
      render(<Adapters />);
    });

    // title and description both reference adapter names; use getAllByText
    expect(screen.getAllByText(/Meshery Adapter for Istio/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/Meshery Adapter for Linkerd/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/Meshery Adapter for Consul/i).length).toBeGreaterThan(0);
  });

  it('renders docs links for each adapter with a slugified test id', async () => {
    await act(async () => {
      render(<Adapters />);
    });

    expect(screen.getByTestId('adapter-docs-istio')).toBeInTheDocument();
    expect(screen.getByTestId('adapter-docs-istio')).toHaveAttribute(
      'href',
      'https://docs.meshery.io/concepts/architecture/adapters',
    );
  });

  it('enables an adapter when its switch is toggled on', async () => {
    changeAdapterStateMock.mockImplementation((cb: any) => cb({}, null));

    await act(async () => {
      render(<Adapters />);
    });

    const istioSwitch = screen.getAllByTestId('switch-OperatorSwitch')[0];

    await act(async () => {
      fireEvent.click(istioSwitch);
    });

    expect(changeAdapterStateMock).toHaveBeenCalledTimes(1);
    const payload = changeAdapterStateMock.mock.calls[0][1];
    expect(payload.status).toBe('ENABLED');
    expect(payload.adapter).toBe('meshery-istio');
    expect(notifyMock).toHaveBeenCalledWith(expect.objectContaining({ event_type: 'success' }));
  });

  it('handles a rejected sync request without throwing', async () => {
    fetchSystemSyncMock.mockReturnValue({
      unwrap: () => Promise.reject(new Error('boom')),
    });

    await act(async () => {
      render(<Adapters />);
    });

    await act(async () => {
      await Promise.resolve();
    });

    // The handler doesn't notify directly (handleError is curried); just verify the
    // initial render still produced cards and didn't crash.
    expect(screen.getAllByText(/Meshery Adapter for Istio/i).length).toBeGreaterThan(0);
  });

  it('marks adapters returned by system sync as enabled', async () => {
    fetchSystemSyncMock.mockReturnValue({
      unwrap: () =>
        Promise.resolve({
          meshAdapters: [{ name: 'ISTIO', adapter_location: 'localhost:10000' }],
        }),
    });

    await act(async () => {
      render(<Adapters />);
    });

    // Allow microtasks to flush
    await act(async () => {
      await Promise.resolve();
    });

    const switches = screen.getAllByTestId('switch-OperatorSwitch');
    // First adapter is ISTIO which should now be enabled
    expect(switches[0]).toBeChecked();
  });

  it('passes the deploy payload through to changeAdapterState on toggle', async () => {
    changeAdapterStateMock.mockImplementation((cb: any) => cb(null, new Error('fail')));

    await act(async () => {
      render(<Adapters />);
    });

    const istioSwitch = screen.getAllByTestId('switch-OperatorSwitch')[0];

    await act(async () => {
      fireEvent.click(istioSwitch);
    });

    // mutation called with a payload describing the deploy action
    expect(changeAdapterStateMock).toHaveBeenCalled();
    const payload = changeAdapterStateMock.mock.calls[0][1];
    expect(payload.adapter).toBe('meshery-istio');
    expect(payload.status).toBe('ENABLED');
  });
});
