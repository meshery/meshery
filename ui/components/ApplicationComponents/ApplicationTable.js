import React from 'react'
import { IconButton, TableCell, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import UndeployIcon from "../../public/static/img/UndeployIcon";

function ApplicationTable({applications= [], setSelectedRowData, user}) {

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
        // responsive : "scrollFullHeight",
        resizableColumns : true,
        serverSide : true,
        rowsPerPageOptions : [10, 20, 25],
        fixedHeader : true,
        print : false,
        download : false,
        textLabels : {
          selectedRows : {
            text : "application(s) selected"
          }
        },
        onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(applications[meta.rowIndex]),
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