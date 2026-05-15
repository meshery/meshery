import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const fetchSystemSyncMock = vi.fn();
const unwrapMock = vi.fn();
const notifyMock = vi.fn();
const manageAdapterMock = vi.fn();
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

vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: notifyMock }),
}));

vi.mock('../../../lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

vi.mock('../../../rtk-query/system', () => ({
  useLazyGetSystemSyncQuery: () => [fetchSystemSyncMock],
  useManageAdapterMutation: () => [manageAdapterMock],
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
    manageAdapterMock.mockReset();
    updateProgressMock.mockReset();

    fetchSystemSyncMock.mockReturnValue({
      unwrap: () =>
        Promise.resolve({
          meshAdapters: [],
        }),
    });
    // Default: manageAdapter resolves so success path runs without throwing.
    manageAdapterMock.mockReturnValue({
      unwrap: () => Promise.resolve({}),
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
    await act(async () => {
      render(<Adapters />);
    });

    const istioSwitch = screen.getAllByTestId('switch-OperatorSwitch')[0];

    await act(async () => {
      fireEvent.click(istioSwitch);
    });

    expect(manageAdapterMock).toHaveBeenCalledTimes(1);
    // ENABLED path passes { meshLocationURL } (DELETE path passes a method).
    const arg = manageAdapterMock.mock.calls[0][0];
    expect(arg.meshLocationURL).toBeDefined();
    expect(arg.method).toBeUndefined();
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

  it('passes the deploy payload through to manageAdapter on toggle', async () => {
    // Reject so the catch path runs and the switch state is rolled back.
    manageAdapterMock.mockReturnValue({
      unwrap: () => Promise.reject(new Error('fail')),
    });

    await act(async () => {
      render(<Adapters />);
    });

    const istioSwitch = screen.getAllByTestId('switch-OperatorSwitch')[0];

    await act(async () => {
      fireEvent.click(istioSwitch);
    });

    expect(manageAdapterMock).toHaveBeenCalled();
    // ENABLED toggle: argument is { meshLocationURL: '<adapter url>' }
    const arg = manageAdapterMock.mock.calls[0][0];
    expect(arg.meshLocationURL).toBeDefined();
    expect(typeof arg.meshLocationURL).toBe('string');
  });
});
