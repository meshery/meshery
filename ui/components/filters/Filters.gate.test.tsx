import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Filters is a large container; its split-out pieces are tested in their own
// files. This file covers only the page access gate: what a viewer without
// the view permission sees, and that no filter list is requested for them.
const hoisted = vi.hoisted(() => ({
  canView: true,
  filtersQuery: vi.fn(),
}));

vi.mock('@sistent/sistent', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@sistent/sistent')>()),
  useHasPermission: () => hoisted.canView,
}));
vi.mock('@/rtk-query/filter', () => {
  const mutation = () => [vi.fn()];
  return {
    useGetFiltersQuery: (...args: unknown[]) => hoisted.filtersQuery(...args),
    useCloneFilterMutation: mutation,
    usePublishFilterMutation: mutation,
    useUnpublishFilterMutation: mutation,
    useDeleteFilterMutation: mutation,
    useUpdateFilterFileMutation: mutation,
    useUploadFilterFileMutation: mutation,
  };
});
vi.mock('@/rtk-query/user', () => ({ useGetProviderCapabilitiesQuery: () => ({}) }));
vi.mock('@/utils/hooks/useTableUrlState', () => ({
  useTableUrlState: () => ({
    tableState: { page: 0, pageSize: 10, sortOrder: '', search: '', filters: { vis: '' } },
    updateTableState: vi.fn(),
  }),
}));
vi.mock('@/utils/hooks/useColumnVisibilityPreference', () => ({
  useColumnVisibilityPreference: () => ({
    columnVisibility: {},
    setColumnVisibilityByUser: vi.fn(),
    setColumnVisibilityByResponsive: vi.fn(),
  }),
}));
vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('../../utils/dimension', () => ({ useWindowDimensions: () => ({ width: 1280 }) }));
vi.mock('../../api/meshmodel', () => ({ getMeshModels: () => Promise.resolve({ models: [] }) }));
vi.mock('@/graphql/queries/CatalogFilterQuery', () => ({
  default: () => ({ subscribe: () => ({ unsubscribe: vi.fn() }) }),
}));
vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { user: {}, catalogVisibility: false } }),
}));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('@/assets/icons', () => ({ Publish: () => null }));
vi.mock('../../utils/utils', () => ({ modifyRJSFSchema: (schema: unknown) => schema }));
vi.mock('../general/error-404/index', () => ({
  default: () => <div data-testid="default-error" />,
}));
// Children only rendered on the allowed path; their own tests cover them.
vi.mock('./FiltersGrid', () => ({ default: () => null }));
vi.mock('./Filters.columns', () => ({ buildFiltersColumns: () => [] }));
vi.mock('./Filters.tableOptions', () => ({ buildFiltersTableOptions: () => ({}) }));
vi.mock('./YAMLEditor', () => ({ default: () => null }));
vi.mock('./ImportModal', () => ({ default: () => null }));
vi.mock('./PublishModal', () => ({ default: () => null }));
vi.mock('../general/ViewSwitch', () => ({ default: () => null }));
vi.mock('../general/PromptComponent', () => ({ default: () => null }));
vi.mock('../shared/Modal/Information/InfoModal', () => ({ default: () => null }));

import { SistentThemeProvider } from '@sistent/sistent';
import Filters from './Filters';

const renderFilters = () =>
  render(
    <SistentThemeProvider>
      <Filters />
    </SistentThemeProvider>,
  );

describe('Filters access gate', () => {
  beforeEach(() => {
    hoisted.filtersQuery.mockReset().mockReturnValue({
      data: undefined,
      isLoading: false,
      refetch: vi.fn(),
    });
  });

  it('shows DefaultError and does not request filters without the view permission', () => {
    hoisted.canView = false;

    renderFilters();

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(hoisted.filtersQuery).toHaveBeenCalledWith(expect.anything(), { skip: true });
    expect(hoisted.filtersQuery).not.toHaveBeenCalledWith(expect.anything(), { skip: false });
  });

  it('requests filters with the view permission', () => {
    hoisted.canView = true;

    renderFilters();

    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
    expect(hoisted.filtersQuery).toHaveBeenCalledWith(expect.anything(), { skip: false });
  });
});
