import type React from 'react';

type BuildFiltersTableOptionsArgs = {
  isLocalProvider: boolean;
  count: number;
  page: number;
  pageSize: number;
  search: string;
  sortOrder: string;
  filters: any[];
  columns: any[];
  searchTimeout: React.MutableRefObject<ReturnType<typeof setTimeout> | null>;
  setPage: (_page: number) => void;
  setPageSize: (_size: number) => void;
  setSearch: (_search: string) => void;
  setSortOrder: (_sortOrder: string) => void;
  setSelectedRowData: (_row: any) => void;
  initFiltersSubscription: (
    _pageNo?: string,
    _pagesize?: string,
    _searchText?: string,
    _order?: string,
  ) => void;
  showmodal: (_count: number) => Promise<string | undefined>;
  deleteFilter: (_id: string) => void;
};

export function buildFiltersTableOptions({
  isLocalProvider,
  count,
  page,
  pageSize,
  search,
  sortOrder,
  filters,
  columns,
  searchTimeout,
  setPage,
  setPageSize,
  setSearch,
  setSortOrder,
  setSelectedRowData,
  initFiltersSubscription,
  showmodal,
  deleteFilter,
}: BuildFiltersTableOptionsArgs) {
  return {
    filter: false,
    viewColumns: false,
    sort: !isLocalProvider,
    search: false,
    filterType: 'textField',
    responsive: 'standard',
    resizableColumns: true,
    serverSide: true,
    count,
    rowsPerPage: pageSize,
    fixedHeader: true,
    page,
    print: false,
    download: false,
    sortOrder: {
      name: 'updated_at',
      direction: 'desc',
    },
    textLabels: {
      selectedRows: {
        text: 'filter(s) selected',
      },
    },

    onCellClick: (_, meta) =>
      meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(filters[meta.rowIndex]),

    onRowsDelete: async function handleDelete(row) {
      let response = await showmodal(Object.keys(row.lookup).length);

      if (response === 'Delete') {
        const fid = Object.keys(row.lookup).map((idx) => filters[idx]?.id);
        fid.forEach((fid) => deleteFilter(fid));
      }
      // if (response === "No")
      // fetchFilters(page, pageSize, search, sortOrder);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          initFiltersSubscription(
            tableState.page.toString(),
            pageSize.toString(),
            search,
            sortOrder,
          );
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          initFiltersSubscription(
            page.toString(),
            tableState.rowsPerPage.toString(),
            search,
            sortOrder,
          );
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              setSearch(tableState.searchText);
            }
          }, 500);
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            initFiltersSubscription(page.toString(), pageSize.toString(), search, order);
            setSortOrder(order);
          }
          break;
      }
    },
    setRowProps: (row, dataIndex, rowIndex) => {
      return {
        'data-cy': `config-row-${rowIndex}`,
      };
    },
    setTableProps: () => {
      return {
        'data-cy': 'filters-grid',
      };
    },
  };
}
