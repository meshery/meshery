import React, {useState, useEffect, useRef } from 'react'
import { IconButton,TableRow, TableCell, Typography, TableSortLabel } from "@mui/material";
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
    setPage,
    setPageSize,
    setSortOrder,
    sortOrder = "",
    setSearch,
    search = "",
    testProfiles = [],
    setProfileForModal,
    showModal,
    fetchTestProfiles,
    handleDelete,
  }) {

      const [selectedProfile, setSelectedProfile] = useState();

      useEffect(() => {
        setProfileForModal(selectedProfile);
      }, [selectedProfile]);
     
      const searchTimeout = useRef(null);

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
      
      columns.forEach((column, idx) => {
        if (column.name === sortOrder.split(" ")[0]) {
          columns[idx].options.sortDirection = sortOrder.split(" ")[1];
        }
      });
      
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
    
        onRowsDelete : async function handleDeleteRow(row) {
          let response = await showModal(Object.keys(row.lookup).length)
          console.log(response)
          if (response === "Yes") {
            const pids = Object.keys(row.lookup).map(idx => testProfiles[idx]?.id)
            pids.forEach(pid => handleDelete(pid))
          }
          if (response === "No") {
            fetchTestProfiles(page, pageSize, search, sortOrder);
          }
        },
        onTableChange : (action, tableState) => {
          const sortInfo = tableState.announceText
            ? tableState.announceText.split(" : ")
            : [];
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
              if (searchTimeout.current) clearTimeout(searchTimeout.current);
              searchTimeout.current = setTimeout(() => {
                if (search !== tableState.searchText) setSearch(tableState.searchText);
              }, 500);
              break;
            case "sort":
              if (sortInfo.length == 2) {
                if (sortInfo[1] === "ascending") {
                  order = `${columns[tableState.activeColumn].name} asc`;
                } else {
                  order = `${columns[tableState.activeColumn].name} desc`;
                }
              }
              if (order !== sortOrder) setSortOrder(order);
              break;
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
                  endpoint={`/api/user/performance/profiles/${testProfiles[rowMeta.rowIndex].id}/results`}
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