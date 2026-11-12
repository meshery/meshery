import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import { Keys } from '@meshery/schemas/permissions';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const {
  mockUseHasPermission,
  mockUseGetMeshModelsQuery,
  mockUseGetComponentsQuery,
  mockUseGetRelationshipsQuery,
  mockUseGetRegistrantsQuery,
  mockUseGetConnectionDefinitionsQuery,
  mockGetMeshModelsData,
  mockGetComponentsData,
  mockGetRelationshipsData,
  mockGetRegistrantsData,
  mockGetConnectionDefinitionsData,
} = vi.hoisted(() => ({
  mockUseHasPermission: vi.fn(),
  mockUseGetMeshModelsQuery: vi.fn(),
  mockUseGetComponentsQuery: vi.fn(),
  mockUseGetRelationshipsQuery: vi.fn(),
  mockUseGetRegistrantsQuery: vi.fn(),
  mockUseGetConnectionDefinitionsQuery: vi.fn(),
  mockGetMeshModelsData: vi.fn(),
  mockGetComponentsData: vi.fn(),
  mockGetRelationshipsData: vi.fn(),
  mockGetRegistrantsData: vi.fn(),
  mockGetConnectionDefinitionsData: vi.fn(),
}));

let modelsResState: any;
let componentsResState: any;
let relationshipsResState: any;
let registrantsResState: any;
let connectionsResState: any;

function createControllableQueryPromise<T>() {
  let resolve!: (value: T) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((res, rej) => {
    resolve = res;
    reject = rej;
  }) as Promise<T> & { abort: ReturnType<typeof vi.fn>; unsubscribe: ReturnType<typeof vi.fn> };
  promise.abort = vi.fn();
  promise.unsubscribe = vi.fn();
  return { promise, resolve, reject };
}

vi.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/settings', push: vi.fn(), isReady: true }),
}));

vi.mock('@sistent/sistent', () => ({
  Button: ({ startIcon: _startIcon, children, ...rest }: any) => (
    <button {...rest}>{children}</button>
  ),
  NoSsr: ({ children }: any) => <>{children}</>,
  AddCircleIcon: () => <svg />,
  ExternalLinkIcon: () => <svg />,
  FileUploadIcon: () => <svg />,
  useHasPermission: (permissionKey: any) => mockUseHasPermission(permissionKey),
  useMediaQuery: () => false,
}));

vi.mock('@/theme', () => ({
  styled:
    (Component: any) =>
    () =>
    ({ children, ...rest }: any) =>
      React.createElement(Component || 'div', rest, children),
  useTheme: () => ({
    breakpoints: { down: () => '@media (max-width:600px)' },
  }),
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: {},
}));

vi.mock('../general/error-404/index', () => ({
  default: ({ permissionKey }: any) => (
    <div data-testid="default-error" data-permission-key={JSON.stringify(permissionKey)} />
  ),
}));

vi.mock('@/assets/styles/general/tool.styles', () => {
  const Div = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
  return {
    MeshModelToolbar: Div,
    MainContainer: Div,
    TreeWrapper: Div,
    DetailsContainer: ({ isEmpty: _isEmpty, children, ...rest }: any) => (
      <div {...rest}>{children}</div>
    ),
    InnerContainer: Div,
    WorkloadsContainer: Div,
  };
});

vi.mock('./TabCard', () => ({
  default: ({ label, count }: any) => (
    <div data-testid={`tab-card-${label}`} data-count={count}>
      {label}: {count}
    </div>
  ),
}));

vi.mock('./MesheryTreeView', () => ({
  default: ({ data, view }: any) => (
    <div data-testid="meshery-tree-view" data-view={view} data-count={data?.length ?? 0} />
  ),
}));

vi.mock('./MeshModelDetails', () => ({
  default: () => <div data-testid="mesh-model-details" />,
}));

vi.mock('./MeshModelMobileDetails', () => ({
  default: () => <div data-testid="mesh-model-mobile-details" />,
}));

vi.mock('./ImportModelModal', () => ({
  default: () => null,
}));

vi.mock('./CreateModelModal', () => ({
  default: () => null,
}));

vi.mock('./CreateRelationshipModal', () => ({
  default: () => null,
}));

vi.mock('./hooks', () => ({
  useInfiniteScrollRef: () => vi.fn(),
  useMeshModelComponentRouter: () => ({
    searchQuery: null,
    selectedPageSize: 25,
    selectedTab: 'Models',
  }),
}));

vi.mock('./helper', () => ({
  groupRelationshipsByKind: (items: any[]) => items,
  removeDuplicateVersions: (items: any[]) => items,
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useGetMeshModelsQuery: (...args: any[]) => mockUseGetMeshModelsQuery(...args),
  useGetComponentsQuery: (...args: any[]) => mockUseGetComponentsQuery(...args),
  useGetRelationshipsQuery: (...args: any[]) => mockUseGetRelationshipsQuery(...args),
  useGetRegistrantsQuery: (...args: any[]) => mockUseGetRegistrantsQuery(...args),
  useGetConnectionDefinitionsQuery: (...args: any[]) =>
    mockUseGetConnectionDefinitionsQuery(...args),
  useLazyGetMeshModelsQuery: () => [mockGetMeshModelsData, modelsResState],
  useLazyGetComponentsQuery: () => [mockGetComponentsData, componentsResState],
  useLazyGetRelationshipsQuery: () => [mockGetRelationshipsData, relationshipsResState],
  useLazyGetRegistrantsQuery: () => [mockGetRegistrantsData, registrantsResState],
  useLazyGetConnectionDefinitionsQuery: () => [
    mockGetConnectionDefinitionsData,
    connectionsResState,
  ],
}));

import MeshModelComponent from './MeshModelComponent';

describe('MeshModelComponent RBAC guard', () => {
  beforeEach(() => {
    modelsResState = { data: undefined, isLoading: false, isFetching: false };
    componentsResState = { data: undefined, isLoading: false, isFetching: false };
    relationshipsResState = { data: undefined, isLoading: false, isFetching: false };
    registrantsResState = { data: undefined, isLoading: false, isFetching: false };
    connectionsResState = { data: undefined, isLoading: false, isFetching: false };

    mockUseHasPermission.mockReset();
    mockUseHasPermission.mockReturnValue(true);

    mockUseGetMeshModelsQuery.mockReset();
    mockUseGetMeshModelsQuery.mockReturnValue({ data: { totalCount: 12 } });
    mockUseGetComponentsQuery.mockReset();
    mockUseGetComponentsQuery.mockReturnValue({ data: { totalCount: 34 } });
    mockUseGetRelationshipsQuery.mockReset();
    mockUseGetRelationshipsQuery.mockReturnValue({ data: { totalCount: 5 } });
    mockUseGetRegistrantsQuery.mockReset();
    mockUseGetRegistrantsQuery.mockReturnValue({ data: { totalCount: 2 } });
    mockUseGetConnectionDefinitionsQuery.mockReset();
    mockUseGetConnectionDefinitionsQuery.mockReturnValue({ data: { totalCount: 4 } });

    mockGetMeshModelsData.mockReset();
    mockGetMeshModelsData.mockResolvedValue({ data: { models: [{ id: 'model-1' }] } });
    mockGetComponentsData.mockReset();
    mockGetComponentsData.mockResolvedValue({ data: { components: [] } });
    mockGetRelationshipsData.mockReset();
    mockGetRelationshipsData.mockResolvedValue({ data: { relationships: [] } });
    mockGetRegistrantsData.mockReset();
    mockGetRegistrantsData.mockResolvedValue({ data: { registrants: [] } });
    mockGetConnectionDefinitionsData.mockReset();
    mockGetConnectionDefinitionsData.mockResolvedValue({ data: { connectionDefinitions: [] } });
  });

  it('skips all count queries, fires no lazy triggers, ignores cached data, and renders DefaultError when denied', () => {
    mockUseHasPermission.mockReturnValue(false);
    modelsResState = {
      data: { totalCount: 99, models: [{ id: 'cached-model' }] },
      isLoading: false,
      isFetching: false,
    };

    render(<MeshModelComponent />);

    expect(mockUseHasPermission).toHaveBeenCalledWith(Keys.MesherySystemViewRegistry);
    expect(mockUseGetMeshModelsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: true });
    expect(mockUseGetComponentsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: true });
    expect(mockUseGetRelationshipsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: true });
    expect(mockUseGetRegistrantsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: true });
    expect(mockUseGetConnectionDefinitionsQuery).toHaveBeenCalledWith(expect.any(Object), {
      skip: true,
    });

    expect(mockGetMeshModelsData).not.toHaveBeenCalled();
    expect(mockGetComponentsData).not.toHaveBeenCalled();
    expect(mockGetRelationshipsData).not.toHaveBeenCalled();
    expect(mockGetRegistrantsData).not.toHaveBeenCalled();
    expect(mockGetConnectionDefinitionsData).not.toHaveBeenCalled();

    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.MesherySystemViewRegistry),
    );
    expect(screen.queryByTestId('meshery-tree-view')).not.toBeInTheDocument();
    expect(screen.queryByTestId('mesh-model-details')).not.toBeInTheDocument();
    expect(screen.queryByTestId('tab-card-Models')).not.toBeInTheDocument();
  });

  it('enables count queries and fetches active view data when permitted', async () => {
    render(<MeshModelComponent />);

    expect(mockUseGetMeshModelsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: false });
    expect(mockUseGetComponentsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: false });
    expect(mockUseGetRelationshipsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: false });
    expect(mockUseGetRegistrantsQuery).toHaveBeenCalledWith(expect.any(Object), { skip: false });
    expect(mockUseGetConnectionDefinitionsQuery).toHaveBeenCalledWith(expect.any(Object), {
      skip: false,
    });

    await waitFor(() => expect(mockGetMeshModelsData).toHaveBeenCalledTimes(1));
    await waitFor(() =>
      expect(screen.getByTestId('meshery-tree-view')).toHaveAttribute('data-count', '1'),
    );
  });

  it('aborts in-flight fetch, renders DefaultError, and drops stale completions on mid-session revocation', async () => {
    const inflight = createControllableQueryPromise<any>();
    mockGetMeshModelsData.mockReturnValueOnce(inflight.promise);

    const { rerender } = render(<MeshModelComponent />);
    await waitFor(() => expect(mockGetMeshModelsData).toHaveBeenCalledTimes(1));

    mockUseHasPermission.mockReturnValue(false);
    rerender(<MeshModelComponent />);

    expect(inflight.promise.unsubscribe).toHaveBeenCalledTimes(1);
    expect(inflight.promise.abort).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(screen.queryByTestId('meshery-tree-view')).not.toBeInTheDocument();
    expect(mockUseGetMeshModelsQuery).toHaveBeenLastCalledWith(expect.any(Object), { skip: true });

    await act(async () => {
      inflight.resolve({ data: { models: [{ id: 'stale-model' }] } });
      await Promise.resolve();
    });

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(screen.queryByTestId('meshery-tree-view')).not.toBeInTheDocument();
  });

  it('halts subsequent registrant model queries if permission is revoked mid-loop', async () => {
    const firstRegistrantModels = createControllableQueryPromise<any>();
    mockGetRegistrantsData.mockResolvedValueOnce({
      data: {
        registrants: [{ hostname: 'ArtifactHub' }, { hostname: 'Kubernetes' }],
      },
    });
    mockGetMeshModelsData.mockReturnValueOnce(firstRegistrantModels.promise);

    const { rerender } = render(<MeshModelComponent externalView="Registrants" />);
    await waitFor(() => expect(mockGetRegistrantsData).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(mockGetMeshModelsData).toHaveBeenCalledTimes(1));

    mockUseHasPermission.mockReturnValue(false);
    rerender(<MeshModelComponent externalView="Registrants" />);

    expect(firstRegistrantModels.promise.abort).toHaveBeenCalledTimes(1);

    await act(async () => {
      firstRegistrantModels.resolve({ data: { models: [{ id: 'r1-m1' }] } });
      await Promise.resolve();
    });

    expect(mockGetMeshModelsData).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
  });

  it('fetches registry data once permission transitions from denied to granted', async () => {
    mockUseHasPermission.mockReturnValue(false);
    const { rerender } = render(<MeshModelComponent />);

    expect(mockGetMeshModelsData).not.toHaveBeenCalled();
    expect(screen.getByTestId('default-error')).toBeInTheDocument();

    mockUseHasPermission.mockReturnValue(true);
    rerender(<MeshModelComponent />);

    await waitFor(() => expect(mockGetMeshModelsData).toHaveBeenCalledTimes(1));
    await waitFor(() =>
      expect(screen.getByTestId('meshery-tree-view')).toHaveAttribute('data-count', '1'),
    );
  });

  it('unsubscribes settled single and registrant handles when permission is revoked', async () => {
    const settledRegistrants = createControllableQueryPromise<any>();
    const settledRegistrantModels = createControllableQueryPromise<any>();
    settledRegistrants.resolve({ data: { registrants: [{ hostname: 'ArtifactHub' }] } });
    settledRegistrantModels.resolve({ data: { models: [{ id: 'r1-m1' }] } });

    mockGetRegistrantsData.mockReturnValueOnce(settledRegistrants.promise);
    mockGetMeshModelsData.mockReturnValueOnce(settledRegistrantModels.promise);

    const { rerender } = render(<MeshModelComponent externalView="Registrants" />);

    await waitFor(() =>
      expect(screen.getByTestId('meshery-tree-view')).toHaveAttribute('data-count', '1'),
    );
    expect(settledRegistrants.promise.unsubscribe).not.toHaveBeenCalled();
    expect(settledRegistrantModels.promise.unsubscribe).not.toHaveBeenCalled();

    mockUseHasPermission.mockReturnValue(false);
    rerender(<MeshModelComponent externalView="Registrants" />);

    expect(settledRegistrants.promise.unsubscribe).toHaveBeenCalledTimes(1);
    expect(settledRegistrantModels.promise.unsubscribe).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
  });
});
