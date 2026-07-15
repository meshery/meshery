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

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick }: any) => (
    <button onClick={onClick} data-testid={`btn-${String(children).trim()}`}>
      {children}
    </button>
  ),
  NoSsr: ({ children }: any) => <>{children}</>,
  AddCircleIcon: () => <svg data-testid="add-icon" />,
  ExternalLinkIcon: () => <svg data-testid="ext-icon" />,
  FileUploadIcon: () => <svg data-testid="upload-icon" />,
  useMediaQuery: () => false,
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: {
      common: { white: '#fff' },
      divider: 'rgba(0, 0, 0, 0.12)',
    },
    breakpoints: {
      down: () => '',
    },
  }),
}));

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
  useGetMeshModelsQuery: () => ({ data: { totalCount: 0 } }),
  useGetComponentsQuery: () => ({ data: { totalCount: 0 } }),
  useGetRelationshipsQuery: () => ({ data: { totalCount: 0 } }),
  useGetRegistrantsQuery: () => ({ data: { totalCount: 0 } }),
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

vi.mock('./CreateRelationshipModal', () => ({
  default: ({ isRelationshipModalOpen }: any) =>
    isRelationshipModalOpen ? <div data-testid="create-relationship-modal" /> : null,
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: {},
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
