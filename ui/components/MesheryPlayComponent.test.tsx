import React from 'react';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import MesheryPlayComponent from './MesheryPlayComponent';

const mockDispatch = vi.fn();
let mockRouter: any = {
  push: vi.fn(),
  events: {
    on: vi.fn(),
    off: vi.fn(),
  },
  query: {},
};

let mockAdapters = [
  { name: 'Istio', adapterLocation: 'meshery-istio:10000' },
  { name: 'Linkerd', adapterLocation: 'meshery-linkerd:10001' },
];

let mockHasPermission = true;

vi.mock('next/router', () => ({
  useRouter: () => mockRouter,
}));

vi.mock('react-redux', () => ({
  useDispatch: () => mockDispatch,
  useSelector: (selector: any) =>
    selector({
      adapter: { meshAdapters: mockAdapters },
    }),
}));

vi.mock('@sistent/sistent', async () => {
  const actual: any = await vi.importActual('@sistent/sistent');
  return {
    ...actual,
    useHasPermission: () => mockHasPermission,
    NoSsr: ({ children }: any) => <>{children}</>,
  };
});

vi.mock('./MesheryAdapterPlayComponent', () => ({
  default: ({ adapter }: any) => (
    <div data-testid="meshery-adapter-play">{adapter?.name || 'no-adapter'}</div>
  ),
}));

vi.mock('./general/error-404/index', () => ({
  default: () => <div data-testid="default-error">Permission Denied</div>,
}));

describe('MesheryPlayComponent', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockHasPermission = true;
    mockAdapters = [
      { name: 'Istio', adapterLocation: 'meshery-istio:10000' },
      { name: 'Linkerd', adapterLocation: 'meshery-linkerd:10001' },
    ];
    mockRouter = {
      push: vi.fn(),
      events: {
        on: vi.fn(),
        off: vi.fn(),
      },
      query: {},
    };
  });

  it('renders permission error when user does not have permission', () => {
    mockHasPermission = false;
    render(<MesheryPlayComponent />);
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
  });

  it('renders "Adapter Unavailable" state when no adapters are configured', () => {
    mockAdapters = [];
    render(<MesheryPlayComponent />);
    expect(screen.getByText('Adapter Unavailable')).toBeInTheDocument();
    expect(screen.getByText('Configure Settings')).toBeInTheDocument();

    fireEvent.click(screen.getByText('Configure Settings'));
    expect(mockRouter.push).toHaveBeenCalledWith('/settings?settingsCategory=Adapters');
  });

  it('renders adapter dropdown and adapter play component with default adapter', () => {
    render(<MesheryPlayComponent />);
    expect(screen.getByLabelText('Select Meshery Adapter')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Istio');
  });

  it('dispatches setAdapter and updates rendered UI when selecting a different adapter', () => {
    const { container } = render(<MesheryPlayComponent />);

    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Istio');

    const input = container.querySelector('input[name="adapter_name"]') as HTMLInputElement;
    expect(input).toBeInTheDocument();
    fireEvent.change(input, { target: { value: 'meshery-linkerd:10001' } });

    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'adapter/setAdapter',
      payload: { selectedAdapter: 'Linkerd' },
    });
    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Linkerd');
  });

  it('renders individual adapter when adapter query param is set', () => {
    mockRouter.query = { adapter: 'meshery-linkerd:10001' };
    render(<MesheryPlayComponent />);

    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Linkerd');
    expect(screen.queryByLabelText('Select Meshery Adapter')).not.toBeInTheDocument();
  });

  it('falls back to default adapter dropdown view when adapter query param is unknown', () => {
    mockRouter.query = { adapter: 'unknown-adapter:9999' };
    render(<MesheryPlayComponent />);

    expect(screen.getByLabelText('Select Meshery Adapter')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Istio');
  });

  it('resets to default adapter when routeChangeComplete introduces an unknown adapter query after selecting another adapter', () => {
    const { container } = render(<MesheryPlayComponent />);

    // Select Linkerd
    const input = container.querySelector('input[name="adapter_name"]') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'meshery-linkerd:10001' } });
    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Linkerd');

    // Simulate route change with unknown adapter query
    mockRouter.query = { adapter: 'unknown-adapter:9999' };
    const routeChangeHandler = mockRouter.events.on.mock.calls.find(
      (call: any) => call[0] === 'routeChangeComplete',
    )?.[1];
    expect(routeChangeHandler).toBeDefined();

    act(() => {
      routeChangeHandler();
    });

    expect(screen.getByLabelText('Select Meshery Adapter')).toBeInTheDocument();
    expect(screen.getByTestId('meshery-adapter-play')).toHaveTextContent('Istio');
  });
});
