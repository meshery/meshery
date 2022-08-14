import React, {useState, useEffect} from 'react'
import {Button, TableCell, IconButton, Paper,Tabs, Tab,  TableSortLabel} from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import BarChartIcon from '@mui/icons-material/BarChart';
import InfoIcon from '@mui/icons-material/Info';
import CustomModal from "@/components/Modal";
import fetchPerformanceResults from "@/features/performance/graphql/queries/PerformanceResultQuery";
import MesheryChart from "@/components/MesheryChart";

function generateResultsForDisplay(results) {
  if (Array.isArray(results)) {
    return results.map((record) => {
      const data = {
        name : record.name,
        mesh : record.mesh,
        test_start_time : record.runner_results.StartTime,
        qps : record.runner_results.ActualQPS.toFixed(1),
        duration : (record.runner_results.ActualDuration / 1000000000).toFixed(1),
        threads : record.runner_results.NumThreads,
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

function generateColumnsForDisplay(sortOrder, setSelectedProfileIdxForChart, setSelectedProfileIdxForNodeDetails) {
  const columns = [
    { name : "name",
      label : "Name",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      }, },
    { name : "mesh",
      label : "Mesh",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      }, },
    { name : "test_start_time",
      label : "Start Time",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      }, },
    { name : "qps",
      label : "QPS",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      }, },
    { name : "duration",
      label : "Duration",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      }, },

    { name : "p50",
      label : "P50",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      }, },

    { name : "p99_9",
      label : "P99.9",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      }, },
    { name : "Chart",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value, tableMeta) {
          return (
            <IconButton aria-label="more" color="inherit" onClick={() => {
              setSelectedProfileIdxForChart(tableMeta.rowIndex)}} >
              <BarChartIcon />
            </IconButton>
          );
        },
      }, },
    { name : "Node Details",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value, tableMeta) {
          return (
            <IconButton aria-label="more" color="inherit" onClick={() => setSelectedProfileIdxForNodeDetails(tableMeta.rowIndex)} >
              <InfoIcon />
            </IconButton>
          );
        },
      }, },
  ];

  return columns.map((column) => {

    return column;
  });
}

function ResultChart({ result, handleTabChange, tabValue }) {
  // if (!result) return <div />;

  // const row = result.runner_results;
  // const boardConfig = result.server_board_config;
  // const serverMetrics = result.server_metrics;
  // const startTime = new Date(row.StartTime);
  // const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
  

  return (
    <Paper
      sx={{ width : "100%",
        maxWidth : "90vw",
        padding : "0.5rem" }}
    >
      <Tabs
        value={tabValue}
        onChange={handleTabChange}
        TabIndicatorProps={{
          sx : {
            backgroundColor : "#00B39F"
          }
        }}
      >
        <Tab label="Performance Chart" />
        <Tab label="Node Details" />
      </Tabs>

      {
        (tabValue == 0) ?
          <div>
        <MesheryChart
                rawdata={[result && result.runner_results ? result : {}]}
                data={[result && result.runner_results ? result.runner_results : {}]}
              />
            {/* {boardConfig && boardConfig !== null && Object.keys(boardConfig).length > 0 && ( */}
                <h1>GrafanaCustomCharts</h1>
                {/* <GrafanaCustomCharts
                  boardPanelConfigs={[boardConfig]}
                  // @ts-ignore
                  boardPanelData={[serverMetrics]}
                  startDate={startTime}
                  from={startTime.getTime().toString()}
                  endDate={endTime}
                  to={endTime.getTime().toString()}
                  liveTail={false}
                /> */}
            {/* )} */}
          </div>
          : (tabValue == 1) ?
            <div>
              <h1>NodeDetails</h1>
              {/* <NodeDetails result={row}/> */}
            </div>
            : <div />
      }
    </Paper>
  );
}

function PerformanceResults({CustomHeader, endpoint, user, elevation = 4, results_selection, updateResultsSelection,}) {

  const [results, setResults] = useState([]);  
  const [sortOrder, setSortOrder] = useState("");
  const [selectedRowChart, setSelectedRowChart] = useState();
  const [selectedRowNodeDetails, setSelectedRowNodeDetails] = useState();
  const [tabValue, setTabValue] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");
  

  

  const columns = generateColumnsForDisplay(sortOrder, (idx) => {
    setSelectedRowChart(results[idx])
    setTabValue(0)
  }, (idx) => {
    setSelectedRowChart(results[idx])
    setTabValue(1)
  } );

  const options = {
    elevation : elevation,
    filter : false,
    sort : !(user?.user_id === "meshery"),
    search : !(user?.user_id === "meshery"),
    filterType : "textField",
    responsive : "scrollFullHeight",
    resizableColumns : true,
    selectableRows : true,
    serverSide : true,
    count,
    page,
    rowsPerPage : pageSize,
    rowsPerPageOptions : [10, 20, 25],
    fixedHeader : true,
    print : false,
    download : false,
     
  };
  
  function handleTabChange(event, newValue) {
    setTabValue(newValue);
  }
 
  useEffect(() => {
    fetchResults(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  function fetchResults(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    fetchPerformanceResults({ selector : {
      pageSize : `${pageSize}`,
      page : `${page}`,
      search : `${encodeURIComponent(search)}`,
      order : `${encodeURIComponent(sortOrder)}`
    },
    profileID : endpoint.split("/")[endpoint.split("/").length - 2] }).subscribe({ next : (res) => {
      // @ts-ignore
      let result = res?.fetchResults
      if (typeof result !== "undefined") {

        if (result) {
          setCount(result.total_count);
          setPageSize(result.page_size);
          setSortOrder(sortOrder);
          setSearch(search);
          setResults(result.results);
          setPageSize(result.page_size);
        }
      }
    },  
 });
  }

  return (
    <>
    <MUIDataTable 
    data={generateResultsForDisplay(results)}
    columns={columns}
    title={CustomHeader}
    options={options}
    />

    <CustomModal
    open={!!selectedRowChart}
    Content={
      <div>
      <ResultChart
        result={selectedRowChart}
        handleTabChange={handleTabChange}
        tabValue={tabValue}
      />
      </div>}
   handleClose={() => setSelectedRowChart(undefined)}
     />
    </>
  )
}

export default PerformanceResults