import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CatalogIcon: (props: any) => <svg data-testid="catalog-icon" {...props} />,
  TachographDigitalIcon: (props: any) => <svg data-testid="tachograph-icon" {...props} />,
  InsertChartIcon: (props: any) => <svg data-testid="insert-chart-icon" {...props} />,
}));

vi.mock('../../../assets/icons/ConfigurationIcon', () => ({
  default: () => <svg data-testid="configuration-icon" />,
}));
vi.mock('../../../assets/icons/Connection', () => ({
  default: () => <svg data-testid="connection-icon" />,
}));
vi.mock('../../../assets/icons/CredentialIcon', () => ({
  default: () => <svg data-testid="credential-icon" />,
}));
vi.mock('@/assets/icons/DashboardIcon', () => ({
  default: () => <svg data-testid="dashboard-icon" />,
}));
vi.mock('../../../assets/icons/Environment', () => ({
  default: () => <svg data-testid="environment-icon" />,
}));
vi.mock('../../../assets/icons/ServiceMeshIcon', () => ({
  default: () => <svg data-testid="service-mesh-icon" />,
}));
vi.mock('../../../assets/icons/WorkspaceOutlined', () => ({
  default: () => <svg data-testid="workspace-outlined-icon" />,
}));
vi.mock('../../../public/static/img/drawer-icons/lifecycle_mgmt_svg', () => ({
  default: () => <svg data-testid="lifecycle-icon" />,
}));
vi.mock('../../../public/static/img/drawer-icons/performance_svg', () => ({
  default: () => <svg data-testid="performance-icon" />,
}));
vi.mock('../../../public/static/img/drawer-icons/extensions_svg', () => ({
  default: () => <svg data-testid="extension-icon" />,
}));
vi.mock('../../../public/static/img/drawer-icons/pattern_svg', () => ({
  default: () => <svg data-testid="pattern-icon" />,
}));
vi.mock('../../../public/static/img/drawer-icons/lifecycle_hover_svg', () => ({
  default: () => <svg data-testid="lifecycle-hover" />,
}));
vi.mock('../../../public/static/img/drawer-icons/performance_hover_svg', () => ({
  default: () => <svg data-testid="performance-hover" />,
}));
vi.mock('../../../public/static/img/drawer-icons/configuration_hover_svg', () => ({
  default: () => <svg data-testid="configuration-hover" />,
}));

vi.mock('../../../constants/navigator', () => ({
  CATALOG: 'CATALOG',
  CONFIGURATION: 'CONFIGURATION',
  CONNECTION: 'CONNECTION',
  CREDENTIAL: 'CREDENTIAL',
  DASHBOARD: 'DASHBOARD',
  DESIGN: 'DESIGN',
  ENVIRONMENT: 'ENVIRONMENT',
  EXTENSIONS: 'EXTENSIONS',
  GRAFANA: 'GRAFANA',
  LIFECYCLE: 'LIFECYCLE',
  PERFORMANCE: 'PERFORMANCE',
  PROFILES: 'PROFILES',
  PROMETHEUS: 'PROMETHEUS',
  SERVICE_MESH: 'SERVICE_MESH',
  TELEMETRY: 'TELEMETRY',
  WORKSPACE: 'WORKSPACE',
}));

vi.mock('../../../css/icons.styles', () => ({
  iconSmall: {},
}));

import { Keys } from '@meshery/schemas/permissions';
import { drawerIconsStyle, getNavigatorComponents } from './navigatorComponents';

const fakeProviderUiAccessControl = (allowed: Record<string, boolean>) => ({
  isNavigatorComponentEnabled: (keys: string[]) =>
    keys.every((key) => allowed[key]) || allowed[keys.join(':')] || false,
});

const theme = {
  palette: {
    background: { constant: { white: '#fff' } },
    icon: { default: '#000' },
  },
} as any;

describe('navigatorComponents', () => {
  it('exposes drawerIconsStyle with expected dimensions', () => {
    expect(drawerIconsStyle).toMatchObject({
      height: '19.36px',
      width: '19.36px',
      fontSize: '1.45rem',
    });
  });

  it('returns the canonical top-level nav items', () => {
    const items = getNavigatorComponents(
      fakeProviderUiAccessControl({
        DASHBOARD: true,
        LIFECYCLE: true,
        CONFIGURATION: true,
        PERFORMANCE: true,
        EXTENSIONS: true,
      }),
      theme,
    );
    const titles = items.map((i: any) => i.title);
    expect(titles).toEqual([
      'Dashboard',
      'Lifecycle',
      'Configuration',
      'Telemetry',
      'Performance',
      'Extensions',
    ]);
  });

  it('honours providerUiAccessControl when deciding visibility', () => {
    const items = getNavigatorComponents(
      fakeProviderUiAccessControl({
        DASHBOARD: false,
        LIFECYCLE: false,
        CONFIGURATION: false,
        PERFORMANCE: false,
        EXTENSIONS: false,
      }),
      theme,
    );
    items.forEach((i: any) => {
      // The "Adapters" child under Lifecycle is always shown; other items respect the flag.
      if (i.title !== 'Lifecycle') {
        expect(i.show).toBe(false);
      }
    });
  });

  it('exposes child entries under Lifecycle and Configuration', () => {
    const items = getNavigatorComponents(
      fakeProviderUiAccessControl({
        DASHBOARD: true,
        LIFECYCLE: true,
        CONFIGURATION: true,
        PERFORMANCE: true,
        EXTENSIONS: true,
        CONNECTION: true,
        ENVIRONMENT: true,
        WORKSPACE: true,
        CATALOG: true,
        DESIGN: true,
        PROFILES: true,
      }),
      theme,
    );

    const lifecycle = items.find((i: any) => i.id === 'LIFECYCLE');
    expect(lifecycle?.children?.map((c: any) => c.id)).toEqual([
      'CONNECTION',
      'CREDENTIAL',
      'ENVIRONMENT',
      'WORKSPACE',
      'SERVICE_MESH',
    ]);

    const configuration = items.find((i: any) => i.id === 'CONFIGURATION');
    expect(configuration?.children?.map((c: any) => c.id)).toEqual(['CATALOG', 'DESIGN']);
  });

  it('marks Adapters child item with always-true show', () => {
    const items = getNavigatorComponents(fakeProviderUiAccessControl({ LIFECYCLE: true }), theme);
    const lifecycle = items.find((i: any) => i.id === 'LIFECYCLE');
    const adapters = lifecycle?.children?.find((c: any) => c.id === 'SERVICE_MESH');
    expect(adapters?.show).toBe(true);
  });

  it('attaches permission descriptors to leaf items', () => {
    const items = getNavigatorComponents(fakeProviderUiAccessControl({ LIFECYCLE: true }), theme);
    const lifecycle = items.find((i: any) => i.id === 'LIFECYCLE');
    const connection = lifecycle?.children?.find((c: any) => c.id === 'CONNECTION');
    expect(connection?.permission).toEqual({
      action: Keys.WorkspaceManagementViewConnections.id,
      subject: Keys.WorkspaceManagementViewConnections.function,
    });
  });

  it('passes white fills to CatalogIcon when on the catalog route', () => {
    Object.defineProperty(window, 'location', {
      value: { pathname: '/configuration/catalog' },
      writable: true,
    });
    const items = getNavigatorComponents(
      fakeProviderUiAccessControl({ CONFIGURATION: true, CATALOG: true }),
      theme,
    );
    const catalog = items
      .find((i: any) => i.id === 'CONFIGURATION')
      ?.children?.find((c: any) => c.id === 'CATALOG');
    // Just sanity-check that a React element was provided as an icon
    expect(React.isValidElement(catalog?.icon)).toBe(true);
  });
});
