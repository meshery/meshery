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
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
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
  MeshModelToolbar: ({ children }: any) => <div>{children}</div>,
  MainContainer: ({ children }: any) => <div>{children}</div>,
  TreeWrapper: ({ children }: any) => <div>{children}</div>,
  DetailsContainer: ({ children }: any) => <div>{children}</div>,
  InnerContainer: ({ children }: any) => <div>{children}</div>,
  CardStyle: ({ children, onClick }: any) => <div onClick={onClick}>{children}</div>,
}));

vi.mock('./MesheryTreeView', () => ({
  default: ({ view, lastItemRef }: any) => (
    <button onClick={() => lastItemRef[view]?.({})}>load next</button>
  ),
}));

vi.mock('./MeshModelDetails', () => ({
  default: () => null,
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
  groupRelationshipsByKind: (relationships: any[]) => relationships,
  removeDuplicateVersions: (models: any[]) => models,
}));

vi.mock('./hooks', () => ({
  useInfiniteScrollRef: (loadNextPage: any) => loadNextPage,
  useMeshModelComponentRouter: () => ({
    searchQuery: null,
    selectedPageSize: 25,
    selectedTab: queryMocks.selectedTab,
  }),
}));

vi.mock('./ImportModelModal', () => ({ default: () => null }));
vi.mock('./CreateModelModal', () => ({ default: () => null }));
vi.mock('./CreateRelationshipModal', () => ({ default: () => null }));
vi.mock('css/icons.styles', () => ({ iconSmall: {} }));

import MeshModelComponent from './MeshModelComponent';

describe('MeshModelComponent pagination', () => {
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

  it('loads the next components page when the components scroll trigger fires', async () => {
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
