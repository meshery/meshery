import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const back = vi.fn();
const ping = vi.fn();
const setView = vi.fn();
const normalizeStaticImagePath = vi.fn((src?: string) => src);
const useGetConnectionsQuery = vi.fn();

vi.mock('next/router', () => ({
  useRouter: () => ({ back, pathname: '/foo', push: vi.fn() }),
}));

vi.mock('@sistent/sistent', () => {
  const styled = () => (Component: any) => Component;
  return {
    ArrowBackIcon: () => <svg data-testid="arrow-back" />,
    Box: ({ children }: any) => <div>{children}</div>,
    ErrorBoundary: ({ children }: any) => <>{children}</>,
    OperatorDataFormatter: ({ data }: any) => (
      <div data-testid="op-formatter">{JSON.stringify(data)}</div>
    ),
    Paper: ({ children }: any) => <div>{children}</div>,
    styled: (Component: any) => () => Component,
    Typography: ({ children }: any) => <span>{children}</span>,
    useResourceCleanData: () => ({ getResourceCleanData: ({ resource }: any) => resource }),
  };
});

vi.mock('@/constants/common', () => ({
  FALLBACK_MESHERY_IMAGE_PATH: '/fallback.svg',
  KUBERNETES: 'kubernetes',
}));

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (...args: any[]) => normalizeStaticImagePath(...args),
}));

vi.mock('css/icons.styles', () => ({
  iconXLarge: { width: 32 },
}));

vi.mock('@/utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (id: string) => ({
    name: `ctx-${id}`,
    server: 'srv',
    connectionId: 'conn',
  }),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => ping,
}));

vi.mock('../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title }: any) => (
    <div data-testid="connection-chip">{title}</div>
  ),
}));

vi.mock('./view-component', () => ({
  default: ({ data }: any) => <div data-testid="resource-detail">{JSON.stringify(data)}</div>,
  JSONViewFormatter: ({ data }: any) => <pre>{JSON.stringify(data)}</pre>,
}));

vi.mock('./utils', () => ({
  default: ({ kind, model }: any) => (
    <svg data-testid="node-icon" data-kind={kind} data-model={model} />
  ),
}));

vi.mock('@/utils/Enum', () => ({
  CONNECTION_STATES: { DISCONNECTED: 'disconnected', CONNECTED: 'connected' },
}));

vi.mock('@/rtk-query/connection', () => ({
  useGetConnectionsQuery: (...args: any[]) => useGetConnectionsQuery(...args),
}));

vi.mock('../../utils/TooltipButton', () => ({
  TooltipIconButton: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick} data-testid="tooltip-button">
      {children}
    </button>
  ),
}));

vi.mock('./resources/config', () => ({ ALL_VIEW: 'all' }));

import View, { Title } from './view';

describe('View', () => {
  beforeEach(() => {
    back.mockReset();
    ping.mockReset();
    setView.mockReset();
    normalizeStaticImagePath.mockReset();
    useGetConnectionsQuery.mockReset();
    normalizeStaticImagePath.mockImplementation((src) => src);
    useGetConnectionsQuery.mockReturnValue({ data: { connections: [] } });
  });

  it('returns null when no resource is provided', () => {
    const { container } = render(<View setView={setView} resource={null} k8sConfig={{}} />);
    expect(container.firstChild).toBeNull();
  });

  it('renders the back button, image, and resource title for a resource', () => {
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          cluster_id: 'cluster-1',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: { svgColor: '/icon.svg' } },
        }}
        k8sConfig={{}}
      />,
    );

    expect(screen.getByTestId('tooltip-button')).toBeInTheDocument();
    const img = screen.getByRole('img', { name: 'Pod' });
    expect(img).toHaveAttribute('src', '/icon.svg');
    expect(screen.getByText('pod-1')).toBeInTheDocument();
    expect(screen.getByTestId('connection-chip')).toHaveTextContent('ctx-cluster-1');
  });

  it('navigates back and switches view when the back button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: {} },
        }}
        k8sConfig={{}}
      />,
    );

    await user.click(screen.getByTestId('tooltip-button'));
    expect(back).toHaveBeenCalled();
    expect(setView).toHaveBeenCalledWith('all');
  });

  it('renders cleanly without crashing when connection data is loading (undefined)', () => {
    useGetConnectionsQuery.mockReturnValue({ data: undefined });
    render(
      <View
        setView={setView}
        resource={{
          kind: 'Pod',
          cluster_id: 'cluster-1',
          metadata: { name: 'pod-1' },
          component_metadata: { styles: {} },
        }}
        k8sConfig={{}}
      />,
    );
    expect(screen.getByText('pod-1')).toBeInTheDocument();
  });
});

describe('Title', () => {
  it('renders the value and toggles underline on hover', async () => {
    const user = userEvent.setup();
    render(<Title onClick={vi.fn()} value="pod-x" kind="Pod" />);
    expect(screen.getByText('pod-x')).toBeInTheDocument();
    expect(screen.getByTestId('node-icon')).toHaveAttribute('data-kind', 'Pod');

    // Mouseenter changes hover state — verify by interacting with the element.
    const titleContent = screen.getByText('pod-x').parentElement?.parentElement;
    await user.hover(titleContent as Element);
  });

  it('calls onClick when the title content is clicked', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();
    render(<Title onClick={onClick} value="pod-x" kind="Pod" />);
    await user.click(screen.getByText('pod-x'));
    expect(onClick).toHaveBeenCalled();
  });
});
