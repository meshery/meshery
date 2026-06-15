import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const queryMocks = vi.hoisted(() => ({
  getMeshModelsData: vi.fn(),
  getComponentsData: vi.fn(),
  getRelationshipsData: vi.fn(),
  getRegistrantsData: vi.fn(),
  selectedTab: 'Models',
}));

const lazyMock = (queryFn: any, data: any = undefined) => [
  queryFn,
  { data, isLoading: false, isFetching: false },
];

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
  default: ({ view, data, lastItemRef }: any) => (
    <div data-testid="meshery-tree-view" data-view={view} data-len={data.length}>
      <button onClick={() => lastItemRef[view]?.({})}>load next</button>
    </div>
  ),
}));

vi.mock('./MeshModelDetails', () => ({
  default: ({ view }: any) => <div data-testid="mesh-model-details" data-view={view} />,
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useLazyGetMeshModelsQuery: () =>
    lazyMock(queryMocks.getMeshModelsData, { models: [], totalCount: 0, pageSize: 25, page: 0 }),
  useLazyGetComponentsQuery: () =>
    lazyMock(queryMocks.getComponentsData, {
      components: [{ id: 'component-1' }],
      totalCount: 50,
      pageSize: 25,
      page: 1,
    }),
  useLazyGetRelationshipsQuery: () =>
    lazyMock(queryMocks.getRelationshipsData, {
      relationships: [{ id: 'relationship-1', kind: 'edge' }],
      totalCount: 50,
      pageSize: 25,
      page: 1,
    }),
  useLazyGetRegistrantsQuery: () =>
    lazyMock(queryMocks.getRegistrantsData, {
      registrants: [],
      totalCount: 0,
      pageSize: 25,
      page: 0,
    }),
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
  useInfiniteScrollRef: (loadNextPage: any) => loadNextPage,
  useMeshModelComponentRouter: () => ({
    searchQuery: null,
    selectedPageSize: 25,
    selectedTab: queryMocks.selectedTab,
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
    queryMocks.selectedTab = 'Models';
    queryMocks.getMeshModelsData.mockResolvedValue({ data: { models: [] } });
    queryMocks.getComponentsData.mockResolvedValue({
      data: { components: [{ id: 'component-1' }], totalCount: 50, pageSize: 25, page: 1 },
    });
    queryMocks.getRelationshipsData.mockResolvedValue({
      data: {
        relationships: [{ id: 'relationship-1', kind: 'edge' }],
        totalCount: 50,
        pageSize: 25,
        page: 1,
      },
    });
    queryMocks.getRegistrantsData.mockResolvedValue({ data: { registrants: [] } });
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

  it('loads the next components page when the components tab scroll trigger fires', async () => {
    queryMocks.selectedTab = 'Components';

    render(<MeshModelComponent />);

    await waitFor(() => {
      expect(queryMocks.getComponentsData).toHaveBeenCalledWith(
        expect.objectContaining({
          params: expect.objectContaining({ page: 0 }),
        }),
        true,
      );
    });

    fireEvent.click(screen.getByText('load next'));

    await waitFor(() => {
      expect(queryMocks.getComponentsData).toHaveBeenCalledWith(
        expect.objectContaining({
          params: expect.objectContaining({ page: 1 }),
        }),
        true,
      );
    });
  });

  it('loads the next relationships page from relationship pagination metadata', async () => {
    queryMocks.selectedTab = 'Relationships';

    render(<MeshModelComponent />);

    await waitFor(() => {
      expect(queryMocks.getRelationshipsData).toHaveBeenCalledWith(
        expect.objectContaining({
          params: expect.objectContaining({ page: 0 }),
        }),
        true,
      );
    });

    fireEvent.click(screen.getByText('load next'));

    await waitFor(() => {
      expect(queryMocks.getRelationshipsData).toHaveBeenCalledWith(
        expect.objectContaining({
          params: expect.objectContaining({ page: 1 }),
        }),
        true,
      );
    });
  });
});
