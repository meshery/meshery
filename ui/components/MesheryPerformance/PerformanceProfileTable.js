import React, {useState, useEffect } from 'react'
import { Avatar, Box, Button, Divider, Dialog, DialogActions, DialogContent, DialogTitle, IconButton,TableRow, TableCell, Typography, Tooltip, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import EditIcon from '@mui/icons-material/Edit';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PerformanceResults from "./PerformanceResults";

function PerformanceProfileTable({
    user,
    page = 0,
    count = 0,
    pageSize = 10,
    testProfiles = [],
    setProfileForModal,
    showModal,}) {

      const [selectedProfile, setSelectedProfile] = useState();

      useEffect(() => {
        setProfileForModal(selectedProfile);
      }, [selectedProfile]);

    const columns = [
        { name : "name",
          label : "Profile",
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

        {
            name : "endpoints",
            label : "Endpoints",
            options: {
                filter: false,
                sort: true,
                customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
                return(
                    <TableCell key={index} onClick={() => sortColumn(index)}>
                        <TableSortLabel  active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                            <b>{column.label}</b>
                        </TableSortLabel>
                    </TableCell>
                )
                }
            }
        },
        { name : "last_run",
          label : "Last Run",
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
        { name : "updated_on",
          label : "Updated On",
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
        { name : "Actions",
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
                return (
                  <>
                    <IconButton
                      aria-label="edit"
                      onClick={(ev) => {
                        ev.stopPropagation();
                        setSelectedProfile(testProfiles[tableMeta.rowIndex]);
                      }}
                    >
                      <EditIcon  />
                    </IconButton>
                    <IconButton
                      aria-label="run"
                    >
                      <PlayArrowIcon  />
                    </IconButton>
                  </>
                );
              },
          }, },
      ];
      
      
      const options = {
        filter : false,
        sort : !(user && user.user_id === "meshery"),
        search : !(user && user.user_id === "meshery"),
        filterType : "textField",
        responsive : "scrollFullHeight",
        resizableColumns : true,
        serverSide : true,
        selectableRows : true,
        count,
        rowsPerPage : pageSize,
        rowsPerPageOptions : [10, 20, 25],
        fixedHeader : true,
        page,
        print : false,
        download : false,
        textLabels : {
          selectedRows : {
            text : "profile(s) selected"
          }
        },
    
    
        expandableRows : true,
        renderExpandableRow : function ExpandableRow(rowData, rowMeta) {
          const colSpan = rowData.length;
          return (
            <TableRow>
              <TableCell />
              <TableCell colSpan={colSpan}>
                <PerformanceResults
                  // @ts-ignore
                  CustomHeader={<Typography variant="h6">Test Results</Typography>}
                  // @ts-ignore
                  // endpoint={`/api/user/performance/profiles/${testProfiles[rowMeta.rowIndex].id}/results`}
                  // @ts-ignore
                  elevation={0}
                />
              </TableCell>
            </TableRow>
          );
        },
      };
  return (
    <div style={{marginTop : "1rem", marginLeft: "1rem"}}>
         <MUIDataTable
    title={<div>Profiles</div>}
    data={testProfiles}
    columns={columns}
    options={options}   

  />
    </div>
  )
}

export default PerformanceProfileTable