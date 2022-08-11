import React from 'react'
import { IconButton, TableCell, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import UndeployIcon from "../../public/static/img/UndeployIcon";

function ApplicationTable({applications= [], count, pageSize, page, setSelectedRowData, handleAppDownload, user}) {

    const columns = [
        {
          name : "name",
          label : "Application Name",
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
          },
        },
        {
          name : "created_at",
          label : "Upload Timestamp",
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
          },
        },
        {
          name : "updated_at",
          label : "Update Timestamp",
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
          },
        },
        {
          name : "source_type",
          label : "Source Type",
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
            customBodyRender : function CustomBody(_, tableMeta) {
              const rowData = applications[tableMeta.rowIndex];
              console.log(rowData);
              return (
                <>
                  <IconButton
                    title="click to download"
                    onClick={() => handleAppDownload(rowData.id ,rowData.type.String, rowData.name)}
                  >
                    {/* <img src={`/static/img/${(rowData.type.String).replaceAll(" ", "_").toLowerCase()}.svg`} width="45px" height="45px" /> */}
                  </IconButton>
                </>
              );
            },
          },
        },
        {
          name : "Actions",
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
            customBodyRender : function CustomBody(_, tableMeta) {
              const rowData = applications[tableMeta.rowIndex];
              return (
                <>
                  <IconButton
                    title="Deploy"
                    onClick={() => handleModalOpen(rowData.application_file, true)}
                  >
                    <DoneAllIcon data-cy="deploy-button" />
                  </IconButton>
                  <IconButton
                    title="Undeploy"
                    onClick={() => handleModalOpen(rowData.application_file, false)}
                  >
                    <UndeployIcon fill="rgba(0, 0, 0, 0.54)" data-cy="undeploy-button" />
                  </IconButton>
                </>
              );
            },
          },
        },
      ];
    
      const options = {
        filter : false,
        sort : !(user && user.user_id === "meshery"),
        search : !(user && user.user_id === "meshery"),
        filterType : "textField",
        responsive : "scrollFullHeight",
        resizableColumns : true,
        serverSide : true,
        count,
        rowsPerPage : pageSize,
        rowsPerPageOptions : [10, 20, 25],
        fixedHeader : true,
        page,
        print : false,
        download : false,
        textLabels : {
          selectedRows : {
            text : "application(s) selected"
          }
        },
        onCellClick : (_, meta) => meta.colIndex !== 3 &&  meta.colIndex !== 4 && setSelectedRowData(applications[meta.rowIndex]),
        setRowProps : (row, dataIndex, rowIndex) => {
          return {
            "data-cy" : `config-row-${rowIndex}`
          }
        },
        setTableProps : () => {
          return {
            "data-cy" : "applications-grid"
          }
        }
    
      }

  return (
    <MUIDataTable
    title={<div>Applications</div>}
    data={applications}
    columns={columns}
    options={options}    
  />
  )
}

export default ApplicationTable