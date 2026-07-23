import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);
const useGetFiltersQuery = vi.fn();
const useGetProviderCapabilitiesQuery = vi.fn();
const getMeshModels = vi.fn();
const routerReplace = vi.fn();

const router = {
  query: {} as Record<string, unknown>,
  pathname: '/configuration/filters',
  push: vi.fn(),
  replace: (...args: unknown[]) => {
    const [{ query }] = args as [{ query: Record<string, string> }];
    if (query) {
      const params = new URLSearchParams();
      for (const [key, value] of Object.entries(query)) {
        if (value !== undefined) params.set(key, value);
      }
      const qs = params.toString();
      window.history.replaceState({}, '', qs ? `${router.pathname}?${qs}` : router.pathname);
    }
    routerReplace(...args);
  },
};

vi.mock('next/router', () => ({
  useRouter: () => router,
}));

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('@/rtk-query/filter', () => ({
  useGetFiltersQuery: (...args: unknown[]) => useGetFiltersQuery(...args),
  useCloneFilterMutation: () => [vi.fn()],
  usePublishFilterMutation: () => [vi.fn()],
  useUnpublishFilterMutation: () => [vi.fn()],
  useDeleteFilterMutation: () => [vi.fn()],
  useUpdateFilterFileMutation: () => [vi.fn()],
  useUploadFilterFileMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/user', () => ({
  useGetProviderCapabilitiesQuery: (...args: unknown[]) => useGetProviderCapabilitiesQuery(...args),
}));

vi.mock('../../api/meshmodel', () => ({
  getMeshModels: (...args: unknown[]) => getMeshModels(...args),
}));

vi.mock('@/graphql/queries/CatalogFilterQuery', () => ({
  default: () => ({
    subscribe: ({ next }: any) => {
      next({ catalogFilters: [] });
      return { unsubscribe: vi.fn() };
    },
  }),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: any) =>
    selector({
      ui: {
        user: { id: 'user-1' },
        catalogVisibility: false,
      },
    }),
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('../../utils/dimension', () => ({
  useWindowDimensions: () => ({ width: 1200 }),
}));

vi.mock('@/utils/hooks/useColumnVisibilityPreference', () => ({
  useColumnVisibilityPreference: () => ({
    columnVisibility: {},
    setColumnVisibilityByUser: vi.fn(),
    setColumnVisibilityByResponsive: vi.fn(),
  }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('../../utils/utils', () => ({
  modifyRJSFSchema: (schema: any) => schema,
}));

vi.mock('../../utils/Enum', () => ({
  MesheryFiltersCatalog: 'meshery-filters-catalog',
  VISIBILITY: { PUBLISHED: 'published', PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('../../lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'ERROR', SUCCESS: 'SUCCESS' },
}));

vi.mock('../general/ViewSwitch', () => ({
  default: () => <div data-testid="view-switch" />,
}));

vi.mock('./FiltersGrid', () => ({
  default: () => <div data-testid="filters-grid" />,
}));

vi.mock('./YAMLEditor', () => ({
  default: () => null,
}));

vi.mock('./ImportModal', () => ({
  default: () => null,
}));

vi.mock('./PublishModal', () => ({
  default: () => null,
}));

vi.mock('../shared/Modal/Information/InfoModal', () => ({
  default: () => null,
}));

vi.mock('../general/PromptComponent', () => ({
  default: React.forwardRef(() => null),
}));

vi.mock('../shared/LoadingState/LoadingComponent', () => ({
  default: () => <div data-testid="loading" />,
}));

vi.mock('../general/error-404/index', () => ({
  default: () => <div data-testid="error-404" />,
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('./Filters.styled', () => ({
  CreateButton: ({ children }: any) => <div>{children}</div>,
  ViewSwitchButton: ({ children }: any) => <div>{children}</div>,
  BtnText: ({ children }: any) => <span>{children}</span>,
}));

vi.mock('@/assets/icons', () => ({
  Publish: () => <svg data-testid="publish-icon" />,
}));

vi.mock('./Filters.columns', () => ({
  buildFiltersColumns: () => [{ name: 'name' }, { name: 'Actions' }],
}));

vi.mock('./Filters.tableOptions', () => ({
  buildFiltersTableOptions: () => ({
    serverSide: true,
    count: 0,
    page: 0,
    rowsPerPage: 10,
  }),
}));

vi.mock('./Filters.fileActions', () => ({
  createDeleteFilter: () => vi.fn(),
  createHandleClone: () => vi.fn(),
  createHandleDownload: () => vi.fn(),
  createHandleImportFilter: () => vi.fn(),
  createHandlePublish: () => vi.fn(),
  createHandleSubmit: () => vi.fn(),
  createHandleUnpublishModal: () => vi.fn(),
  createUploadHandler: () => vi.fn(),
}));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    NoSsr: ({ children }: any) => <>{children}</>,
    CustomColumnVisibilityControl: () => <div data-testid="column-visibility" />,
    ResponsiveDataTable: () => <div data-testid="filters-table" />,
    SearchBar: () => <div data-testid="search-bar" />,
    UniversalFilter: ({ handleApplyFilter, setSelectedFilters }: any) => (
      <div data-testid="universal-filter">
        <button
          type="button"
          data-testid="apply-public"
          onClick={() => {
            const next = { visibility: 'public' };
            setSelectedFilters(next);
            handleApplyFilter(next);
          }}
        >
          Apply Public
        </button>
        <button
          type="button"
          data-testid="apply-private"
          onClick={() => {
            const next = { visibility: 'private' };
            setSelectedFilters(next);
            handleApplyFilter(next);
          }}
        >
          Apply Private
        </button>
        <button
          type="button"
          data-testid="apply-all"
          onClick={() => {
            const next = { visibility: 'All' };
            setSelectedFilters(next);
            handleApplyFilter(next);
          }}
        >
          Apply All
        </button>
      </div>
    ),
    Button: ({ children, onClick }: any) => (
      <button type="button" onClick={onClick}>
        {children}
      </button>
    ),
    publishCatalogItemSchema: {},
    publishCatalogItemUiSchema: {},
    PROMPT_VARIANTS: { DANGER: 'danger' },
    styled,
  };
});

import MesheryFilters from './Filters';

const filtersData = [
  { id: '1', name: 'public-filter', visibility: 'public', userId: 'user-1' },
  { id: '2', name: 'private-filter', visibility: 'private', userId: 'user-1' },
];

describe('Filters visibility filter apply', () => {
  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
    router.query = {};
    routerReplace.mockReset();
    window.history.replaceState({}, '', '/configuration/filters');
    useGetFiltersQuery.mockReset();
    useGetProviderCapabilitiesQuery.mockReset();
    getMeshModels.mockReset();

    useGetFiltersQuery.mockReturnValue({
      data: { filters: filtersData, totalCount: filtersData.length },
      isLoading: false,
      refetch: vi.fn(),
    });
    useGetProviderCapabilitiesQuery.mockReturnValue({
      data: { capabilities: [], providerType: 'remote' },
    });
    getMeshModels.mockResolvedValue({ models: [] });
  });

  const renderComponent = () => render(<MesheryFilters />);

  it('writes fil_vis to the URL on the first Apply click', async () => {
    const user = userEvent.setup();
    renderComponent();

    await user.click(screen.getByTestId('apply-public'));

    await waitFor(() => {
      expect(routerReplace).toHaveBeenCalledWith(
        expect.objectContaining({
          pathname: '/configuration/filters',
          query: expect.objectContaining({ fil_vis: 'public' }),
        }),
        undefined,
        { shallow: true },
      );
    });
  });

  it('uses the current Apply payload instead of the previous filter value', async () => {
    const user = userEvent.setup();
    renderComponent();

    await user.click(screen.getByTestId('apply-public'));
    router.query = { fil_vis: 'public' };

    await user.click(screen.getByTestId('apply-private'));

    await waitFor(() => {
      expect(routerReplace).toHaveBeenLastCalledWith(
        expect.objectContaining({
          query: expect.objectContaining({ fil_vis: 'private' }),
        }),
        undefined,
        { shallow: true },
      );
    });
  });

  it('clears fil_vis from the URL when visibility is reset to All', async () => {
    const user = userEvent.setup();
    router.query = { fil_vis: 'public' };
    window.history.replaceState({}, '', '/configuration/filters?fil_vis=public');
    renderComponent();

    await user.click(screen.getByTestId('apply-all'));

    await waitFor(() => {
      expect(routerReplace).toHaveBeenLastCalledWith(
        expect.objectContaining({
          query: expect.not.objectContaining({ fil_vis: expect.anything() }),
        }),
        undefined,
        { shallow: true },
      );
    });
  });
});
