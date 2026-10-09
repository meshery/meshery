import { describe, expect, it, vi, beforeEach } from 'vitest';
import { buildFiltersTableOptions } from './Filters.tableOptions';

const filters = [
  { id: 'f-1', name: 'filter1' },
  { id: 'f-2', name: 'filter2' },
];

const columns = [
  { name: 'name' },
  { name: 'created_at' },
  { name: 'updated_at' },
  { name: 'visibility' },
  { name: 'Actions' },
];

describe('buildFiltersTableOptions', () => {
  let setPage: ReturnType<typeof vi.fn>;
  let setPageSize: ReturnType<typeof vi.fn>;
  let setSearch: ReturnType<typeof vi.fn>;
  let setSortOrder: ReturnType<typeof vi.fn>;
  let setSelectedRowData: ReturnType<typeof vi.fn>;
  let showmodal: ReturnType<typeof vi.fn>;
  let deleteFilter: ReturnType<typeof vi.fn>;
  let searchTimeout: { current: ReturnType<typeof setTimeout> | null };

  beforeEach(() => {
    setPage = vi.fn();
    setPageSize = vi.fn();
    setSearch = vi.fn();
    setSortOrder = vi.fn();
    setSelectedRowData = vi.fn();
    showmodal = vi.fn();
    deleteFilter = vi.fn();
    searchTimeout = { current: null };
  });

  const build = (overrides: any = {}) =>
    buildFiltersTableOptions({
      isLocalProvider: false,
      count: 10,
      page: 0,
      pageSize: 25,
      search: '',
      sortOrder: 'updated_at desc',
      filters,
      columns,
      searchTimeout: searchTimeout as any,
      setPage,
      setPageSize,
      setSearch,
      setSortOrder,
      setSelectedRowData,
      showmodal,
      deleteFilter,
      ...overrides,
    });

  it('disables sort on the local provider and enables it on a remote provider', () => {
    expect(build({ isLocalProvider: true }).sort).toBe(false);
    expect(build().sort).toBe(true);
  });

  it('exposes serverSide, count, rowsPerPage and page from arguments', () => {
    const options = build({ count: 99, page: 4, pageSize: 50 });
    expect(options.serverSide).toBe(true);
    expect(options.count).toBe(99);
    expect(options.page).toBe(4);
    expect(options.rowsPerPage).toBe(50);
  });

  it('selects row data when cell click is NOT on visibility / actions columns', () => {
    const options = build();
    options.onCellClick(null, { colIndex: 0, rowIndex: 1 });
    expect(setSelectedRowData).toHaveBeenCalledWith(filters[1]);
  });

  it('does NOT select row data for the visibility column (colIndex 3)', () => {
    const options = build();
    options.onCellClick(null, { colIndex: 3, rowIndex: 0 });
    expect(setSelectedRowData).not.toHaveBeenCalled();
  });

  it('does NOT select row data for the actions column (colIndex 4)', () => {
    const options = build();
    options.onCellClick(null, { colIndex: 4, rowIndex: 0 });
    expect(setSelectedRowData).not.toHaveBeenCalled();
  });

  it('onRowsDelete deletes only when modal returns "Delete"', async () => {
    showmodal.mockResolvedValueOnce('Delete');
    const options = build();
    await options.onRowsDelete({ lookup: { 0: true, 1: true } });

    expect(showmodal).toHaveBeenCalledWith(2);
    expect(deleteFilter).toHaveBeenCalledTimes(2);
    expect(deleteFilter).toHaveBeenNthCalledWith(1, 'f-1');
    expect(deleteFilter).toHaveBeenNthCalledWith(2, 'f-2');
  });

  it('onRowsDelete is a no-op when modal returns something else', async () => {
    showmodal.mockResolvedValueOnce('Cancel');
    const options = build();
    await options.onRowsDelete({ lookup: { 0: true } });

    expect(deleteFilter).not.toHaveBeenCalled();
  });

  it('onTableChange "changePage" updates page', () => {
    const options = build();
    options.onTableChange('changePage', { page: 3 });
    expect(setPage).toHaveBeenCalledWith(3);
  });

  it('onTableChange "changeRowsPerPage" updates page size', () => {
    const options = build();
    options.onTableChange('changeRowsPerPage', { rowsPerPage: 50 });
    expect(setPageSize).toHaveBeenCalledWith(50);
  });

  it('onTableChange "search" debounces and eventually calls setSearch when text differs', async () => {
    const options = build();
    options.onTableChange('search', { searchText: 'foo', announceText: null });
    expect(searchTimeout.current).toBeTruthy();

    // Wait for the debounce timer
    await new Promise((resolve) => setTimeout(resolve, 600));
    expect(setSearch).toHaveBeenCalledWith('foo');
  });

  it('onTableChange "sort" updates sortOrder using announceText direction', () => {
    const options = build({ sortOrder: 'updated_at desc' });
    options.onTableChange('sort', { activeColumn: 0, announceText: 'col : ascending' });
    expect(setSortOrder).toHaveBeenCalledWith('name asc');
  });

  it('onTableChange "sort" is a no-op when the computed order matches current sortOrder', () => {
    const options = build({ sortOrder: 'updated_at desc' });
    options.onTableChange('sort', { activeColumn: 2, announceText: 'col : descending' });
    // updated_at desc was already the sortOrder, so the setter is not called.
    expect(setSortOrder).not.toHaveBeenCalled();
  });

  it('setRowProps and setTableProps expose data-cy attributes', () => {
    const options = build();
    expect(options.setRowProps({}, 0, 2)).toEqual({ 'data-cy': 'config-row-2' });
    expect(options.setTableProps()).toEqual({ 'data-cy': 'filters-grid' });
  });
});
