import React, { useState, useEffect, useCallback, useMemo } from 'react';
import PropTypes from 'prop-types';
import { NoSsr, TableRow, TableCell, TableSortLabel } from '@mui/material';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MUIDataTable from 'mui-datatables';
import Moment from 'react-moment';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import { updateResultsSelection, clearResultsSelection, updateProgress } from '../lib/store';
import CustomToolbarSelect from './CustomToolbarSelect';
import MesheryChart from './MesheryChart';
import GrafanaCustomCharts from './telemetry/grafana/GrafanaCustomCharts';
import MesheryResultDialog from './MesheryResultDialog';
import { withNotify } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { Box, IconButton } from '@layer5/sistent';
import { useLazyGetResultsQuery } from '@/rtk-query/meshResult';

const DEFAULT_PAGE_SIZE = 10;
const ROWS_PER_PAGE_OPTIONS = [10, 20, 25];

const MesheryResults = ({
  classes,
  results_selection,
  user,
  endpoint,
  updateProgress,
  updateResultsSelection,
  notify,
  customHeader,
}) => {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [results, setResults] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [searchTimeout, setSearchTimeout] = useState(null);

  // RTK Query hook for fetching results
  const [trigger, { isFetching, error }] = useLazyGetResultsQuery();

  // Fetch results using RTK Query
  const fetchResults = useCallback(
    async (page, pageSize, search, sortOrder) => {
      updateProgress({ showProgress: true });

      try {
        const result = await trigger({
          endpoint,
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
    return results.map((record) => {
      const row = {
        name: record.name,
        mesh: record.mesh,
        test_start_time: record.runner_results.StartTime,
        qps: record.runner_results.ActualQPS.toFixed(1),
        duration: (record.runner_results.ActualDuration / 1000000000).toFixed(1),
        threads: record.runner_results.NumThreads,
      };

      const percentiles = record.runner_results.DurationHistogram?.Percentiles || [];
      percentiles.forEach(({ Percentile, Value }) => {
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
              onClick={() => setSelectedRowData(results[tableMeta.rowIndex])}
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
    (currentRowsSelected, allRowsSelected) => {
      const selectedResults = {};
      allRowsSelected.forEach(({ dataIndex }) => {
        if (dataIndex < pageSize) {
          selectedResults[dataIndex] = results[dataIndex];
        }
      });
      updateResultsSelection({ page, results: selectedResults });
    },
    [page, pageSize, results, updateResultsSelection],
  );

  // Handle table state changes
  const handleTableChange = useCallback(
    (action, tableState) => {
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
            }, 500),
          );
          break;
        case 'sort': {
          const sortInfo = tableState.announceText?.split(' : ') || [];
          const order = sortInfo[1] === 'ascending' ? 'asc' : 'desc';
          setSortOrder(`${columns[tableState.activeColumn].name} ${order}`);
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
        <CustomToolbarSelect
          selectedRows={selectedRows}
          displayData={displayData}
          setSelectedRows={setSelectedRows}
          results={results}
        />
      ),
      expandableRows: true,
      renderExpandableRow: (rowData, rowMeta) => {
        const row = results[rowMeta.dataIndex].runner_results;
        const boardConfig = results[rowMeta.dataIndex].server_board_config;
        const serverMetrics = results[rowMeta.dataIndex].server_metrics;
        const startTime = new Date(row.StartTime);
        const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
        const colSpan = rowData.length + 1;

        return (
          <TableRow>
            <TableCell colSpan={colSpan}>
              <div className={classes.chartContent}>
                <MesheryChart rawdata={[results[rowMeta.dataIndex]]} data={[row]} hideTitle />
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
      classes.chartContent,
    ],
  );

  return (
    <NoSsr>
      {selectedRowData && (
        <MesheryResultDialog rowData={selectedRowData} close={resetSelectedRowData} />
      )}
      <MUIDataTable
        title={
          customHeader || (
            <Box fontWeight="bolder" fontSize="18">
              Performance Test Results
            </Box>
          )
        }
        data={resultsForDisplay}
        columns={columns}
        options={options}
      />
    </NoSsr>
  );
};

MesheryResults.propTypes = {
  classes: PropTypes.object.isRequired,
  results_selection: PropTypes.object.isRequired,
  user: PropTypes.object,
  endpoint: PropTypes.string,
  updateProgress: PropTypes.func.isRequired,
  updateResultsSelection: PropTypes.func.isRequired,
  notify: PropTypes.func.isRequired,
  customHeader: PropTypes.node,
};

const mapDispatchToProps = (dispatch) => ({
  updateResultsSelection: bindActionCreators(updateResultsSelection, dispatch),
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => ({
  startKey: state.get('results').get('startKey'),
  results: state.get('results').get('results').toArray(),
  results_selection: state.get('results_selection').toObject(),
  user: state.get('user')?.toObject(),
});

export default connect(mapStateToProps, mapDispatchToProps)(withNotify(MesheryResults));
