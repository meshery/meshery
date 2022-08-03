import React, {useState} from 'react'
import {TableCell, IconButton, Paper,  TableSortLabel} from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import BarChartIcon from '@mui/icons-material/BarChart';
import InfoIcon from '@mui/icons-material/Info';


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
            <IconButton aria-label="more" color="inherit" >
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
            <IconButton aria-label="more" color="inherit" >
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

function PerformanceResults({CustomHeader, user, elevation = 4, results_selection,}) {

  const [results, setResults] = useState([]);  
  const columns = generateColumnsForDisplay();

  // sortOrder, (idx) => {
  //   setSelectedRowChart(results[idx])
  // }, (idx) => {    setSelectedRowNodeDetails(results[idx])
  // }

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
    rowsPerPageOptions : [10, 20, 25],
    fixedHeader : true,
    // rowsSelected : generateSelectedRows(results_selection),
    print : false,
    download : false,

  };

  return (
    <>
    <MUIDataTable 
    data={generateResultsForDisplay(results)}
    columns={columns}
    title={CustomHeader}
    options={options}
    />
    </>
  )
}

export default PerformanceResults