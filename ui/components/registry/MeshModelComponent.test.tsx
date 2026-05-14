import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const lazyMock = () => [vi.fn(), { data: undefined, isLoading: false, isFetching: false }];

vi.mock('next/router', () => ({
  useRouter: () => ({
    query: {},
    pathname: '/settings',
    push: vi.fn(),
    route: '/settings',
  }),
}));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<any>();
  return {
    ...actual,
    createTheme: (theme: any = {}) => ({
      ...theme,
      breakpoints: theme.breakpoints ?? {
        values: { xs: 0, sm: 600, md: 900, lg: 1200, xl: 1536 },
        up: () => '',
        down: () => '',
        between: () => '',
      },
    }),
    EventBus: class {
      publish() {}
      on() {
        return { subscribe: () => ({ unsubscribe() {} }) };
      }
      onAny() {
        return { subscribe: () => ({ unsubscribe() {} }) };
      }
    },
    Button: ({ children, onClick }: any) => (
      <button onClick={onClick} data-testid={`btn-${String(children).trim()}`}>
        {children}
      </button>
    ),
    NoSsr: ({ children }: any) => <>{children}</>,
    AddCircleIcon: () => <svg data-testid="add-icon" />,
    ExternalLinkIcon: () => <svg data-testid="ext-icon" />,
    FileUploadIcon: () => <svg data-testid="upload-icon" />,
    InfoIcon: () => <svg data-testid="info-icon" />,
  };
});

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
  COMPONENTS: 'Components',
  RELATIONSHIPS: 'Relationships',
  REGISTRANTS: 'Registrants',
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  MeshModelToolbar: ({ children }: any) => <div data-testid="toolbar">{children}</div>,
  MainContainer: ({ children }: any) => <div>{children}</div>,
  TreeWrapper: ({ children }: any) => <div data-testid="tree-wrapper">{children}</div>,
  DetailsContainer: ({ children }: any) => <div>{children}</div>,
  InnerContainer: ({ children }: any) => <div>{children}</div>,
  CardStyle: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('./MesheryTreeView', () => ({
  default: ({ view, data }: any) => (
    <div data-testid="meshery-tree-view" data-view={view} data-len={data.length} />
  ),
}));

vi.mock('./MeshModelDetails', () => ({
  default: ({ view }: any) => <div data-testid="mesh-model-details" data-view={view} />,
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: () => lazyMock(),
  useLazyGetComponentsQuery: () => lazyMock(),
  useLazyGetRelationshipsQuery: () => lazyMock(),
  useLazyGetRegistrantsQuery: () => lazyMock(),
}));

vi.mock('./helper', () => ({
  groupRelationshipsByKind: (r: any[]) => r,
  removeDuplicateVersions: (m: any[]) => m,
}));

vi.mock('./hooks', () => ({
  useInfiniteScrollRef: () => ({ current: null }),
  useMeshModelComponentRouter: () => ({
    searchQuery: null,
    selectedPageSize: 25,
    selectedTab: 'Models',
  }),
}));

vi.mock('./ImportModelModal', () => ({
  default: ({ isImportModalOpen }: any) =>
    isImportModalOpen ? <div data-testid="import-modal" /> : null,
}));

vi.mock('./CreateModelModal', () => ({
  default: ({ isCreateModalOpen }: any) =>
    isCreateModalOpen ? <div data-testid="create-modal" /> : null,
}));

vi.mock('@/components/relationship-builder/CreateRelationshipModal', () => ({
  default: ({ open }: any) => (open ? <div data-testid="create-relationship-modal" /> : null),
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: {},
}));

vi.mock('../layout/Navigator/navigatorComponents', () => ({
  drawerIconsStyle: {},
  getNavigatorComponents: () => [],
}));

vi.mock('../AppComponents', () => ({
  Footer: () => null,
  KubernetesSubscription: () => null,
  NavigationBar: () => null,
}));

vi.mock('../layout/Header/Header', () => ({
  default: () => null,
  K8sContextConnectionChip: () => null,
}));

vi.mock('../subscription/helpers', () => ({
  MESHERY_CONTROLLER_SUBSCRIPTION: 'MESHERY_CONTROLLER_SUBSCRIPTION',
  fnMapping: {},
  isControllerObjectEqual: () => true,
}));

vi.mock('../general/error-404', () => ({
  default: () => null,
}));

vi.mock('../general/error-404/index', () => ({
  default: () => null,
}));

vi.mock('../../ui.config', () => ({
  default: {
    loadingComponent: null,
    loadingComponentDark: null,
  },
}));

vi.mock('machines/operationsCenter', () => ({
  OPERATION_CENTER_EVENTS: { EVENT_RECEIVED_FROM_SERVER: 'event' },
}));

import MeshModelComponent from './MeshModelComponent';

describe('MeshModelComponent', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the toolbar with create / import buttons', () => {
    render(<MeshModelComponent />);
    expect(screen.getByTestId('toolbar')).toBeInTheDocument();
  });

  it('renders the MesheryTreeView and MeshModelDetails subcomponents', () => {
    render(<MeshModelComponent />);
    expect(screen.getByTestId('meshery-tree-view')).toHaveAttribute('data-view', 'Models');
    expect(screen.getByTestId('mesh-model-details')).toHaveAttribute('data-view', 'Models');
  });

  it('uses externalView when provided', () => {
    render(<MeshModelComponent externalView="Components" />);
    expect(screen.getByTestId('meshery-tree-view')).toHaveAttribute('data-view', 'Components');
  });
});
