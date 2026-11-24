import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { SistentThemeProvider } from '@sistent/sistent';

// Covers how the Filters container wires rows-per-page changes into the URL
// table state. Data hooks, children and the table option builder are stubbed;
// the builder is captured so the container's own setPageSize can be called.

// Hook results must be referentially stable across renders: the container
// keys effects on them.
const { buildTableOptions, updateTableState, TABLE_STATE, FILTERS_QUERY_RESULT, REDUX_STATE } =
  vi.hoisted(() => {
    const updateTableState = vi.fn();
    return {
      buildTableOptions: vi.fn((_options: any) => ({})),
      updateTableState,
      TABLE_STATE: {
        tableState: {
          page: 2,
          pageSize: 10,
          sortOrder: '',
          search: '',
          filters: { vis: '' },
        },
        updateTableState,
      },
      FILTERS_QUERY_RESULT: {
        data: { filters: [], total_count: 0 },
        isLoading: false,
        refetch: () => {},
      },
      REDUX_STATE: { ui: { user: {}, catalogVisibility: false } },
    };
  });

vi.mock('@/utils/hooks/useTableUrlState', () => ({ useTableUrlState: () => TABLE_STATE }));
vi.mock('@/utils/hooks/useColumnVisibilityPreference', () => {
  const setter = () => {};
  const result = {
    columnVisibility: {},
    setColumnVisibilityByUser: setter,
    setColumnVisibilityByResponsive: setter,
  };
  return { useColumnVisibilityPreference: () => result };
});
vi.mock('@/rtk-query/filter', () => {
  const mutation = [() => {}];
  return {
    useGetFiltersQuery: () => FILTERS_QUERY_RESULT,
    useCloneFilterMutation: () => mutation,
    usePublishFilterMutation: () => mutation,
    useUnpublishFilterMutation: () => mutation,
    useDeleteFilterMutation: () => mutation,
    useUpdateFilterFileMutation: () => mutation,
    useUploadFilterFileMutation: () => mutation,
  };
});
vi.mock('@/rtk-query/user', () => ({ useGetProviderCapabilitiesQuery: () => ({ data: {} }) }));
vi.mock('../../api/meshmodel', () => ({ getMeshModels: () => new Promise(() => {}) }));
vi.mock('@/graphql/queries/CatalogFilterQuery', () => ({
  default: () => ({ subscribe: () => ({ unsubscribe: () => {} }) }),
}));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) => selector(REDUX_STATE),
}));
vi.mock('./Filters.columns', () => ({ buildFiltersColumns: () => [] }));
vi.mock('./Filters.tableOptions', () => ({ buildFiltersTableOptions: buildTableOptions }));
vi.mock('./FiltersGrid', () => ({ default: () => null }));
vi.mock('./YAMLEditor', () => ({ default: () => null }));
vi.mock('./ImportModal', () => ({ default: () => null }));
vi.mock('./PublishModal', () => ({ default: () => null }));
vi.mock('../general/PromptComponent', () => ({ default: React.forwardRef(() => null) }));
vi.mock('../shared/LoadingState/LoadingComponent', () => ({ default: () => null }));
vi.mock('../shared/Modal/Information/InfoModal', () => ({ default: () => null }));
vi.mock('../general/error-404/index', () => ({ default: () => null }));
vi.mock('../../pages/_app', () => ({ mesheryExtensionRoute: '/' }));
vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    ResponsiveDataTable: () => null,
    useHasPermission: () => true,
  };
});

import MesheryFilters from './Filters';

describe('MesheryFilters rows per page', () => {
  it('resets to the first page in the same URL update as the new page size', () => {
    render(
      <SistentThemeProvider initialMode="light">
        <MesheryFilters />
      </SistentThemeProvider>,
    );
    // Ignore updates the page makes on mount.
    updateTableState.mockReset();
    const { setPageSize } = buildTableOptions.mock.calls.at(-1)[0];
    setPageSize(25);

    // One update carrying both keys: two separate updates would each be built
    // from the same URL snapshot, and the second would drop the first.
    expect(updateTableState).toHaveBeenCalledTimes(1);
    expect(updateTableState).toHaveBeenCalledWith({ pageSize: 25, page: 0 });
  });
});
