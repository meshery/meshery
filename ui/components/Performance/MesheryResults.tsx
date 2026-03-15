import React, { useState, useEffect, useCallback, useMemo } from 'react';
import MUIDataTable from '@sistent/mui-datatables';
import Moment from 'react-moment';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import CustomToolbarSelect from '../CustomToolbarSelect';
import MesheryChart from '../MesheryChart';
import GrafanaCustomCharts from '../telemetry/grafana/GrafanaCustomCharts';
import MesheryResultDialog from '../MesheryResultDialog';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { Box, IconButton, NoSsr, TableRow, TableCell, TableSortLabel } from '@sistent/sistent';
import { useLazyGetResultsQuery } from '@/rtk-query/meshResult';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import { updateResultsSelection } from '@/store/slices/prefTest';

const DEFAULT_PAGE_SIZE = 10;
const ROWS_PER_PAGE_OPTIONS = [10, 20, 25];

const MesheryResults = () => {
  const { notify } = useNotification();
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [results, setResults] = useState<any[]>([]);
  const [selectedRowData, setSelectedRowData] = useState<any | null>(null);
  const [searchTimeout, setSearchTimeout] = useState<ReturnType<typeof setTimeout> | null>(null);
  const { user } = useSelector((state: any) => state.ui);
  const dispatch = useDispatch();
  const { results_selection } = useSelector((state: any) => state.prefTest);
  // RTK Query hook for fetching results
  const [trigger, { isFetching, error }] = useLazyGetResultsQuery();

  // Fetch results using RTK Query
  const fetchResults = useCallback(
    async (page, pageSize, search, sortOrder) => {
      updateProgress({ showProgress: true });

      try {
        const result = await trigger({
          page,
          pagesize: pageSize,
          search,
          sortOrder,
        }).unwrap(); // Unwrap the result to handle errors properly

        if (result) {
          setResults(result.results || []);
          setSearch(search);
          setSortOrder(sortOrder);
          setPage(result.page || 0);
          setPageSize(result.page_size || DEFAULT_PAGE_SIZE);
          setCount(result.total_count || 0);
        }
      } catch (error) {
        handleError(error);
      } finally {
        updateProgress({ showProgress: false });
      }
    },
    [trigger, updateProgress],
  );

  // Handle API errors
  const handleError = useCallback(
    (error) => {
      notify({
        message: `There was an error fetching results: ${error.message || 'Unknown error'}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    },
    [notify],
  );

  // Fetch results on component mount or when dependencies change
  useEffect(() => {
    fetchResults(page, pageSize, search, sortOrder);
  }, [fetchResults, page, pageSize, search, sortOrder]);

  // Handle RTK Query errors
  useEffect(() => {
    updateProgress({ showProgress: isFetching });
    if (error) {
      handleError(error);
    }
  }, [error, handleError, isFetching]);

  // Reset selected row data
  const resetSelectedRowData = () => setSelectedRowData(null);

  // Format results for display in the table
  const resultsForDisplay = useMemo(() => {
    return results.map((record: any) => {
      const row: any = {
        name: record.name,
        mesh: record.mesh,
        test_start_time: record.runner_results.StartTime,
        qps: record.runner_results.ActualQPS.toFixed(1),
        duration: (record.runner_results.ActualDuration / 1000000000).toFixed(1),
        threads: record.runner_results.NumThreads,
      };

      const percentiles = record.runner_results.DurationHistogram?.Percentiles || [];
      percentiles.forEach(({ Percentile, Value }: any) => {
        row[`p${Percentile}`.replace('.', '_')] = Value.toFixed(3);
      });

      // Default values for percentiles if not available
      if (percentiles.length === 0) {
        row.p50 = 0;
        row.p75 = 0;
        row.p90 = 0;
        row.p99 = 0;
        row.p99_9 = 0;
      }

      return row;
    });
  }, [results]);

  // Columns configuration for the table
  const columns = useMemo(
    () => [
      {
        name: 'name',
        label: 'Name',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }, sortColumn) => (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          ),
        },
      },
      {
        name: 'mesh',
        label: 'Mesh',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }, sortColumn) => (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          ),
        },
      },
      {
        name: 'test_start_time',
        label: 'Start Time',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          customHeadRender: ({ index, ...column }, sortColumn) => (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          ),
          customBodyRender: (value) => <Moment format="LLLL">{value}</Moment>,
        },
      },
      {
        name: 'qps',
        label: 'QPS',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: ({ index, ...column }) => (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          ),
        },
      },
      {
        name: 'duration',
        label: 'Duration',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: ({ index, ...column }) => (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          ),
        },
      },
      {
        name: 'p50',
        label: 'P50',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: ({ index, ...column }) => (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          ),
        },
      },
      {
        name: 'p99_9',
        label: 'P99.9',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: ({ index, ...column }) => (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          ),
        },
      },
      {
        name: 'Details',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: ({ index, ...column }) => (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          ),
          customBodyRender: (value, tableMeta) => (
            <IconButton
              aria-label="more"
              color="inherit"
              onClick={() => {
                const rowData = results[tableMeta.rowIndex];
                if (rowData) {
                  setSelectedRowData(rowData);
                }
              }}
            >
              <MoreHorizIcon />
            </IconButton>
          ),
        },
      },
    ],
    [results],
  );

  // Handle row selection
  const handleRowSelection = useCallback(
    (currentRowsSelected: any, allRowsSelected: any) => {
      const selectedResults: any = {};
      allRowsSelected.forEach(({ dataIndex }: any) => {
        if (dataIndex < pageSize && results[dataIndex]) {
          selectedResults[dataIndex] = results[dataIndex];
        }
      });
      dispatch(updateResultsSelection({ page, results: selectedResults }));
    },
    [page, pageSize, results, dispatch],
  );

  // Handle table state changes
  const handleTableChange = useCallback(
    (action: string, tableState: any) => {
      switch (action) {
        case 'changePage':
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout) clearTimeout(searchTimeout);
          setSearchTimeout(
            setTimeout(() => {
              if (search !== tableState.searchText) {
                setSearch(tableState.searchText || '');
              }
            }, 500) as ReturnType<typeof setTimeout>,
          );
          break;
        case 'sort': {
          const sortInfo = tableState.announceText?.split(' : ') || [];
          const order = sortInfo[1] === 'ascending' ? 'asc' : 'desc';
          const activeColumn = columns[tableState.activeColumn];
          if (activeColumn) {
            setSortOrder(`${activeColumn.name} ${order}`);
          }
          break;
        }
        default:
          break;
      }
    },
    [columns, search, searchTimeout],
  );

  // Table options
  const options = useMemo(
    () => ({
      filter: false,
      sort: !(user && user.user_id === 'meshery'),
      search: !(user && user.user_id === 'meshery'),
      filterType: 'textField',
      responsive: 'standard',
      resizableColumns: true,
      selectableRows: true,
      serverSide: true,
      count,
      rowsPerPage: pageSize,
      rowsPerPageOptions: ROWS_PER_PAGE_OPTIONS,
      fixedHeader: true,
      page,
      rowsSelected: Object.keys(results_selection).flatMap((pg) =>
        parseInt(pg) === page
          ? Object.keys(results_selection[pg]).map((ind) => parseInt(ind))
          : Object.keys(results_selection[pg]).map(
              (ind) => (parseInt(pg) + 1) * pageSize + parseInt(ind) + 1,
            ),
      ),
      print: false,
      download: false,
      onRowsSelect: handleRowSelection,
      onTableChange: handleTableChange,
      customToolbarSelect: (selectedRows, displayData, setSelectedRows) => (
        <CustomToolbarSelect setSelectedRows={setSelectedRows} />
      ),
      expandableRows: true,
      renderExpandableRow: (rowData: any, rowMeta: any) => {
        const result = results[rowMeta.dataIndex];
        if (!result || !result.runner_results) {
          return null;
        }
        const row = result.runner_results;
        const boardConfig = result.server_board_config;
        const serverMetrics = result.server_metrics;
        const startTime = new Date(row.StartTime);
        const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
        const colSpan = rowData.length + 1;

        return (
          <TableRow>
            <TableCell colSpan={colSpan}>
              <div>
                <MesheryChart rawdata={[result]} data={[row]} hideTitle />
              </div>
              {boardConfig && Object.keys(boardConfig).length > 0 && (
                <div>
                  <GrafanaCustomCharts
                    boardPanelConfigs={[boardConfig]}
                    boardPanelData={[serverMetrics]}
                    startDate={startTime}
                    from={startTime.getTime().toString()}
                    endDate={endTime}
                    to={endTime.getTime().toString()}
                    liveTail={false}
                  />
                </div>
              )}
            </TableCell>
          </TableRow>
        );
      },
    }),
    [
      user,
      count,
      pageSize,
      page,
      results_selection,
      handleRowSelection,
      handleTableChange,
      results,
    ],
  );

  return (
    <NoSsr>
      {selectedRowData && (
        <MesheryResultDialog rowData={selectedRowData} close={resetSelectedRowData} />
      )}
      <MUIDataTable
        title={
          <Box fontWeight="bolder" fontSize="18">
            Performance Test Results
          </Box>
        }
        data={resultsForDisplay}
        columns={columns}
        options={options}
      />
    </NoSsr>
  );
};

export default MesheryResults;
