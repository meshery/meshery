import React from 'react';
import { NoSsr, TableCell, } from '@material-ui/core';
import { MuiThemeProvider } from '@material-ui/core/styles';
import {  createTheme } from '@material-ui/core/styles';
import MUIDataTable from "mui-datatables";
import TableSortLabel from '@material-ui/core/TableSortLabel'

function NodeDetails(props) {

  const getMuiTheme = () => createTheme({
    shadows : ["none"],
    overrides : {
      MuiInput : {
        underline : {
          "&:hover:not(.Mui-disabled):before" : {
            borderBottom : "2px solid #222"
          },
          "&:after" : {
            borderBottom : "2px solid #222"
          }
        }
      },
      MUIDataTableSearch : {
        searchIcon : {
          color : "#607d8b" ,
          marginTop : "7px",
          marginRight : "8px",
        },
        clearIcon : {
          "&:hover" : {
            color : "#607d8b"
          }
        },
      },
      MUIDataTableSelectCell : {
        checkboxRoot : {
          '&$checked' : {
            color : '#607d8b',
          },
        },
      },
      MUIDataTableToolbar : {
        iconActive : {
          color : "#222"
        },
        icon : {
          "&:hover" : {
            color : "#607d8b"
          }
        },
      }
    }
  })
  const chartData = props.result;

  const columns = [
    { name : "hostname",
      label : "Hostname",
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
    { name : "cpu",
      label : "CPU",
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
    { name : "memory",
      label : "Memory",
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
    { name : "arch",
      label : "Arch",
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
    { name : "os",
      label : "OS",
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
    { name : "kubeletVersion",
      label : "Kubelet Version",
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
    { name : "containerRuntime",
      label : "Container Runtime",
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
      }, },]

  let data = []

  const options = {
    filter : false,
    selectableRows : false,
  };

  let server = chartData?.kubernetes?.server_version;

  chartData?.kubernetes?.nodes.map((node) => {
    let arr = []
    let m = node?.allocatable_memory
    const mem = (String(m).slice(0, String(m).length-2)*0.000001024).toPrecision(5);
    arr.push(node?.hostname);
    arr.push(node?.allocatable_cpu);
    arr.push(mem+"Gi");
    arr.push(node?.architecture);
    arr.push(node?.operating_system);
    arr.push(node?.kubelet_version);
    arr.push(node?.container_runtime_version);

    data.push(arr)
  })

  return (
    <NoSsr>
      <MuiThemeProvider theme={getMuiTheme()}>
        <MUIDataTable
          title={<div style={{ fontSize : 18 }}> Kubernetes Server Version: {server}</div>}
          data={data}
          options={options}
          columns={columns}
        />
      </MuiThemeProvider>
    </NoSsr>
  )
}

export default NodeDetails;
