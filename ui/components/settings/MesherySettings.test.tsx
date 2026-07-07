import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const routerState = {
  query: {} as Record<string, any>,
  pathname: '/settings',
  route: '/settings',
  push: vi.fn(),
};

vi.mock('next/router', () => ({
  useRouter: () => routerState,
}));

vi.mock('next/link', () => ({
  default: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

vi.mock('@sistent/sistent', () => ({
  NoSsr: ({ children }: any) => <>{children}</>,
  CustomTooltip: ({ children, title, value }: any) => (
    <span data-tooltip-title={String(title)} data-tooltip-value={String(value)}>
      {children}
    </span>
  ),
  AppBar: ({ children }: any) => <div>{children}</div>,
  Typography: ({ children, sx, color, variant, ...rest }: any) => (
    <span data-variant={variant} {...rest}>
      {children}
    </span>
  ),
  Tabs: ({ children, value, onChange }: any) => (
    <div data-testid="tabs" data-value={value} onClick={(e) => onChange?.(e, 'mock-tab')}>
      {children}
    </div>
  ),
  Tab: ({ label, value, disabled }: any) => (
    <button data-testid={`tab-${value}`} disabled={disabled}>
      {label}
    </button>
  ),
  Paper: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  LeftArrowIcon: () => <svg />,
  PollIcon: () => <svg />,
  DatabaseIcon: () => <svg />,
  MendeleyIcon: () => <svg />,
  FileIcon: () => <svg />,
  SettingsIcon: () => <svg />,
  useTheme: () => ({
    palette: {
      icon: { default: 'icon' },
      primary: { main: 'primary' },
    },
    spacing: (n: number) => n,
  }),
  styled: (Component: any) => () => {
    const StyledComponent = ({ children, ...rest }: any) =>
      typeof Component === 'string'
        ? React.createElement(Component, rest, children)
        : React.createElement(Component, rest, children);
    return StyledComponent;
  },
}));

vi.mock('../dashboard/charts/DashboardMeshModelGraph', () => ({
  default: () => <div data-testid="dashboard-graph" />,
}));

vi.mock('../MeshAdapterConfigComponent', () => ({
  default: () => <div data-testid="adapter-config" />,
}));

vi.mock('../PromptComponent', () => ({
  default: () => <div data-testid="prompt-component" />,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('../DatabaseSummary', () => ({
  default: () => <div data-testid="database-summary" />,
}));

vi.mock('../../api/meshmodel', () => ({
  getComponentsDetail: vi.fn().mockResolvedValue({ totalCount: 0 }),
  getMeshModels: vi.fn().mockResolvedValue({ models: [] }),
  getRelationshipsDetail: vi.fn().mockResolvedValue({ totalCount: 0 }),
  getMeshModelRegistrants: vi.fn().mockResolvedValue({ totalCount: 0 }),
}));

let canMockReturn = true;
vi.mock('@/utils/can', () => ({
  default: () => canMockReturn,
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: {
    VIEW_SETTINGS: { action: 'a', subject: 's' },
    VIEW_CLOUD_NATIVE_INFRASTRUCTURE: { action: 'a', subject: 's' },
    VIEW_METRICS: { action: 'a', subject: 's' },
    VIEW_REGISTRY: { action: 'a', subject: 's' },
    VIEW_OVERVIEW: { action: 'a', subject: 's' },
  },
}));

vi.mock('@/constants/navigator', () => ({
  METRICS: 'Metrics',
  ADAPTERS: 'Adapters',
  RESET: 'Reset',
  GRAFANA: 'Grafana',
  PROMETHEUS: 'Prometheus',
  OVERVIEW: 'Overview',
  REGISTRY: 'Registry',
  CONTROLLERS: 'Controllers',
}));

vi.mock('../registry/helper', () => ({
  removeDuplicateVersions: (m: any[]) => m,
}));

vi.mock('../registry/MeshModelComponent', () => ({
  default: () => <div data-testid="mesh-model-component" />,
}));

vi.mock('./MesheryControllersConfig', () => ({
  default: () => <div data-testid="controllers-config" />,
}));

vi.mock('../general/error-404', () => ({
  default: () => <div data-testid="error-404" />,
}));

vi.mock('../dashboard/charts/MesheryConfigurationCharts', () => ({
  default: () => <div data-testid="config-charts" />,
}));

vi.mock('../dashboard/charts/ConnectionCharts', () => ({
  default: () => <div data-testid="connection-charts" />,
}));

vi.mock('../dashboard/style', () => ({
  SecondaryTab: ({ label }: any) => <button>{label}</button>,
  SecondaryTabs: ({ children }: any) => <div>{children}</div>,
}));

let selectorReturn: any = {
  ui: { selectedK8sContexts: [], k8sConfig: { clusterConfigured: false } },
  telemetry: { prometheus: {}, grafana: {} },
  adapter: { meshAdapters: [] },
};

vi.mock('react-redux', () => ({
  useSelector: (fn: any) => fn(selectorReturn),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetProviderCapabilitiesQuery: () => ({
    data: { providerName: 'TestProv', providerType: 'cloud' },
  }),
}));

import MesherySettings from './MesherySettings';

describe('MesherySettings', () => {
  beforeEach(() => {
    canMockReturn = true;
    routerState.query = {};
    routerState.push.mockReset();
    selectorReturn = {
      ui: { selectedK8sContexts: [], k8sConfig: { clusterConfigured: false } },
      telemetry: { prometheus: {}, grafana: {} },
      adapter: { meshAdapters: [] },
    };
  });

  it('renders the main settings tab bar with the 5 settings categories', () => {
    render(<MesherySettings />);

    expect(screen.getByTestId('tab-Overview')).toBeInTheDocument();
    expect(screen.getByTestId('tab-Adapters')).toBeInTheDocument();
    expect(screen.getByTestId('tab-Registry')).toBeInTheDocument();
    expect(screen.getByTestId('tab-Controllers')).toBeInTheDocument();
    expect(screen.getByTestId('tab-Reset')).toBeInTheDocument();
    expect(screen.queryByTestId('tab-Metrics')).not.toBeInTheDocument();
  });

  it('renders the controllers configuration tab when selected', () => {
    routerState.query = { settingsCategory: 'Controllers' };
    render(<MesherySettings />);
    expect(screen.getByTestId('controllers-config')).toBeInTheDocument();
  });

  it('renders the Overview tab content (dashboard graph) by default', () => {
    render(<MesherySettings />);
    expect(screen.getByTestId('dashboard-graph')).toBeInTheDocument();
    expect(screen.getByTestId('config-charts')).toBeInTheDocument();
    expect(screen.getByTestId('connection-charts')).toBeInTheDocument();
  });

  it('returns null when VIEW_SETTINGS permission is denied', () => {
    canMockReturn = false;
    const { container } = render(<MesherySettings />);
    expect(container.textContent).toBe('');
  });
});
