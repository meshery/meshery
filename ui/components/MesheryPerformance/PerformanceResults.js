import React, {useState} from 'react'
import {TableCell, IconButton, Paper,Tabs, Tab,  TableSortLabel} from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import BarChartIcon from '@mui/icons-material/BarChart';
import InfoIcon from '@mui/icons-material/Info';
import CustomModal from "@/components/Modal";
import MesheryChart from "@/components/MesheryChart";

const resultstest = [
  {
  meshery_id: "e94247f1-5b8e-4735-a4f4-0631d4f618d1",
  name: "app mesh_1647676936903",
  test_start_time: "2022-03-19T13:32:20.529138Z",
  mesh: "app mesh",
  user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
  runner_results: {
  AbortOn: 0,
  ActualDuration: 30762529800,
  ActualQPS: 1.0402265420966776,
  DurationHistogram: {
  Avg: 0.96132905625,
  Count: 32,
  Data: [
  {
  Count: 7,
  End: 0.9,
  Percent: 21.875,
  Start: 0.8091348
  },
  {
  Count: 15,
  End: 1,
  Percent: 68.75,
  Start: 0.9
  },
  {
  Count: 10,
  End: 1.2735744,
  Percent: 100,
  Start: 1
  }
  ],
  Max: 1.2735744,
  Min: 0.8091348,
  Percentiles: [
  {
  Percentile: 50,
  Value: 0.96
  },
  {
  Percentile: 75,
  Value: 1.05471488
  },
  {
  Percentile: 90,
  Value: 1.186030592
  },
  {
  Percentile: 99,
  Value: 1.2648200192
  },
  {
  Percentile: 99.9,
  Value: 1.27269896192
  }
  ],
  StdDev: 0.08486755232774887,
  Sum: 30.7625298
  },
  Exactly: 0,
  HeaderSizes: {
  Avg: 0,
  Count: 32,
  Data: [
  {
  Count: 32,
  End: 0,
  Percent: 100,
  Start: 0
  }
  ],
  Max: 0,
  Min: 0,
  Percentiles: null,
  StdDev: 0,
  Sum: 0
  },
  Jitter: false,
  Labels: "app mesh_1647676936903 -_- https://youtu.be/JNoL5CLrY68",
  NumThreads: 1,
  RequestedDuration: "30s",
  RequestedQPS: "max",
  RetCodes: {
  200: 32
  },
  RunID: 0,
  RunType: "HTTP",
  Sizes: {
  Avg: 753774.28125,
  Count: 32,
  Data: [
  {
  Count: 27,
  End: 750000,
  Percent: 84.375,
  Start: 675404
  },
  {
  Count: 5,
  End: 1136885,
  Percent: 100,
  Start: 1000000
  }
  ],
  Max: 1136885,
  Min: 675404,
  Percentiles: null,
  StdDev: 155880.00727828167,
  Sum: 24120777
  },
  SocketCount: 0,
  StartTime: "2022-03-19T13:32:20.5291384+05:30",
  URL: "https://youtu.be/JNoL5CLrY68",
  Version: "dev",
  kubernetes: {
  nodes: null,
  server_version: ""
  },
  },
  performance_profile: "acbff055-721d-4e57-a363-cdc51f9d43cc",
  created_at: "2022-03-19T08:02:50.903797Z",
  updated_at: "2022-03-19T08:02:50.903807Z"
  }
  ]

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
              console.log("Chart clicked", tableMeta.rowIndex);
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
             <MesheryChart  />
              {/* <MesheryChart
                rawdata={[result && result.runner_results ? result : {}]}
                data={[result && result.runner_results ? result.runner_results : {}]}
              /> */}
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


function PerformanceResults({CustomHeader, user, elevation = 4, results_selection,}) {

  const [results, setResults] = useState([]);  
  const [sortOrder, setSortOrder] = useState("");
  const [selectedRowChart, setSelectedRowChart] = useState();
  const [selectedRowNodeDetails, setSelectedRowNodeDetails] = useState();
  const [tabValue, setTabValue] = useState(0);
  const [open, setOpen] = useState(false);
   

  const handleOpenModal = () => setOpen(true);
  const handleCloseModal = () => setOpen(false);

  const columns = generateColumnsForDisplay(sortOrder, (idx) => {
    setSelectedRowChart(results[idx])
    handleOpenModal();
    console.log("Selected row chart", results[idx]);
    setTabValue(0)
  }, (idx) => {
    setSelectedRowNodeDetails(results[idx])
    handleOpenModal();
    console.log("Selected row chart", results[idx]);
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
    rowsPerPageOptions : [10, 20, 25],
    fixedHeader : true,
    // rowsSelected : generateSelectedRows(results_selection),
    print : false,
    download : false,

  };
  
  function handleTabChange(event, newValue) {
    setTabValue(newValue);
  }

  return (
    <>
    <MUIDataTable 
    data={resultstest}
    columns={columns}
    title={CustomHeader}
    options={options}
    />

    <CustomModal
    open={open}
    Content={
      <div>
      <ResultChart
        // result={selectedRowChart}
        handleTabChange={handleTabChange}
        tabValue={tabValue}
      />
      </div>}
    handleClose={handleCloseModal}
     />
    </>
  )
}

export default PerformanceResults