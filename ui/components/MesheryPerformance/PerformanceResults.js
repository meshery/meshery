//@ts-check
import React, { useEffect, useState, useRef } from "react";
import { NoSsr, TableCell, IconButton, Typography, Paper } from "@material-ui/core";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import { updateResultsSelection, clearResultsSelection, updateProgress } from "../../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../../lib/data-fetch";
import CustomToolbarSelect from "../CustomToolbarSelect";
import MesheryChart from "../MesheryChart";
import GrafanaCustomCharts from "../GrafanaCustomCharts";
import GenericModal from "../GenericModal";
import BarChartIcon from '@material-ui/icons/BarChart';

function generateResultsForDisplay(results) {
  if (Array.isArray(results)) {
    return results.map((record) => {
      const data = {
        name: record.name,
        mesh: record.mesh,
        test_start_time: record.runner_results.StartTime,
        qps: record.runner_results.ActualQPS.toFixed(1),
        duration: (record.runner_results.ActualDuration / 1000000000).toFixed(1),
        threads: record.runner_results.NumThreads,
      };

      if (record.runner_results?.DurationHistogram?.Percentiles) {
        record.runner_results.DurationHistogram.Percentiles.forEach(({ Percentile, Value }) => {
          data[`p${Percentile}`.replace(".", "_")] = Value.toFixed(3);
        });
      } else {
        data.p50 = 0;
        data.p75 = 0;
        data.p90 = 0;
        data.p99 = 0;
        data.p99_9 = 0;
      }

      return data;
    });
  }

  return [];
}

function generateColumnsForDisplay(sortOrder, setSelectedProfileIdx) {
  const columns = [
    {
      name: "name",
      label: "Name",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: "mesh",
      label: "Mesh",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: "test_start_time",
      label: "Start Time",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: "qps",
      label: "QPS",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name: "duration",
      label: "Duration",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },

    {
      name: "p50",
      label: "P50",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },

    {
      name: "p99_9",
      label: "P99.9",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name: "Details",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          return (
            <IconButton aria-label="more" color="inherit" onClick={() => setSelectedProfileIdx(tableMeta.rowIndex)}>
              <BarChartIcon />
            </IconButton>
          );
        },
      },
    },
  ];

  return columns.map((column) => {
    if (column.name === sortOrder.split(" ")[0]) {
      column.options.sortDirection = sortOrder.split(" ")[1];
    }

    return column;
  });
}

function generateSelectedRows(results_selection, page, pageSize) {
  const rowsSelected = [];
  Object.keys(results_selection).forEach((pg) => {
    if (parseInt(pg) !== page) {
      Object.keys(results_selection[parseInt(pg)]).forEach((ind) => {
        const val = (parseInt(pg) + 1) * pageSize + parseInt(ind) + 1;
        rowsSelected.push(val);
      });
    } else {
      Object.keys(results_selection[page]).forEach((ind) => {
        const val = parseInt(ind);
        rowsSelected.push(val);
      });
    }
  });

  return rowsSelected;
}

function ResultChart({ result }) {
  if (!result) return <div />;

  const row = result.runner_results;
  const boardConfig = result.server_board_config;
  const serverMetrics = result.server_metrics;
  const startTime = new Date(row.StartTime);
  const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
  return (
    <Paper
      style={{
        width: "100%",
        maxWidth: "90vw",
        padding: "0.5rem"
      }}
    >
      <div>
        <Typography variant="h6" gutterBottom align="center">Performance Graph</Typography>
        <MesheryChart data={[result && result.runner_results ? result.runner_results : {}]} />
      </div>
      {boardConfig && boardConfig !== null && Object.keys(boardConfig).length > 0 && (
        <div>
          <GrafanaCustomCharts
            boardPanelConfigs={[boardConfig]}
            // @ts-ignore
            boardPanelData={[serverMetrics]}
            startDate={startTime}
            from={startTime.getTime().toString()}
            endDate={endTime}
            to={endTime.getTime().toString()}
            liveTail={false}
          />
        </div>
      )}
    </Paper>
  );
}

/**
 *
 * @param {{
 *  updateProgress?: any,
 *  enqueueSnackbar?: any,
 *  closeSnackbar?: any,
 *  results_selection?: any,
 *  user?: any
 *  updateResultsSelection?: any,
 *  classes?: any
 *  endpoint: string,
 *  CustomHeader?: JSX.Element
 *  elevation?: Number
 * }} props
 */
function MesheryResults({
  updateProgress,
  enqueueSnackbar,
  closeSnackbar,
  endpoint,
  updateResultsSelection,
  results_selection,
  user,
  CustomHeader = <div />,
  elevation = 4,
}) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  const [sortOrder, setSortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [results, setResults] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState();

  const searchTimeout = useRef();

  useEffect(() => {
    fetchResults(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  function fetchResults(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;
    updateProgress({ showProgress: true });

    dataFetch(
      `${endpoint}${query}`,
      {
        method: "GET",
        credentials: "include",
      },
      (result) => {
        console.log("Results API", `${endpoint}${query}`);

        updateProgress({ showProgress: false });

        if (result) {
          setCount(result.total_count);
          setPageSize(result.page_size);
          setSortOrder(sortOrder);
          setSearch(search);
          setResults(result.results);
          setPageSize(result.page_size);
        }
      },
      handleError
    );
  }

  function handleError(error) {
    updateProgress({ showProgress: false });

    enqueueSnackbar(`There was an error fetching results: ${error}`, {
      variant: "error",
      action: function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration: 8000,
    });
  }

  const columns = generateColumnsForDisplay(sortOrder, (idx) => setSelectedRowData(results[idx]));

  const options = {
    elevation: elevation,
    filter: false,
    sort: !(user?.user_id === "meshery"),
    search: !(user?.user_id === "meshery"),
    filterType: "textField",
    responsive: "scrollFullHeight",
    resizableColumns: true,
    selectableRows: true,
    serverSide: true,
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    rowsSelected: generateSelectedRows(results_selection, page, pageSize),
    print: false,
    download: false,
    onRowsSelect: (_, allRowsSelected) => {
      // const rs = self.props.results_selection;
      const res = {};
      allRowsSelected.forEach(({ dataIndex }) => {
        if (dataIndex < pageSize) {
          if (res[dataIndex]) delete res[dataIndex];
          else res[dataIndex] = results[dataIndex];
        }
      });

      updateResultsSelection({ page, results: res });
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          setPage(tableState.page);
          break;
        case "changeRowsPerPage":
          setPageSize(tableState.rowsPerPage);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          // @ts-ignore
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              setSearch(tableState.searchText || "");
            }
          }, 500);
          break;
        case "sort":
          if (sortInfo.length == 2) {
            if (sortInfo[1] === "ascending") order = `${columns[tableState.activeColumn].name} asc`;
            else order = `${columns[tableState.activeColumn].name} desc`;
          }

          if (order !== sortOrder) setSortOrder(order);
          break;
      }
    },
    customToolbarSelect: function CustomToolbarSelectComponent(selectedRows, displayData, setSelectedRows) {
      return (
        <CustomToolbarSelect
          selectedRows={selectedRows}
          displayData={displayData}
          setSelectedRows={setSelectedRows}
          results={results}
        />
      );
    },
  };

  return (
    <NoSsr>
      <MUIDataTable
        title={CustomHeader}
        data={generateResultsForDisplay(results)}
        columns={columns}
        // @ts-ignore
        options={options}
      />

      <GenericModal
        open={!!selectedRowData}
        // @ts-ignore
        Content={<ResultChart result={selectedRowData} />}
        handleClose={() => setSelectedRowData(undefined)}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateResultsSelection: bindActionCreators(updateResultsSelection, dispatch),
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const startKey = state.get("results").get("startKey");
  const results = state.get("results").get("results").toArray();
  const results_selection = state.get("results_selection").toObject();
  const user = state.get("user").toObject();
  if (typeof results !== "undefined") {
    return {
      startKey,
      results,
      results_selection,
      user,
    };
  }
  return { results_selection, user };
};

// @ts-ignore
export default connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryResults));
