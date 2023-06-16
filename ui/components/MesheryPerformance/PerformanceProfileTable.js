// @ts-check
import React, { useState, useRef, useEffect } from "react";
import { withStyles } from "@material-ui/core/styles";
import {
  NoSsr, TableCell, IconButton, TableRow, Typography
} from "@material-ui/core";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import { updateProgress } from "../../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import PerformanceResults from "./PerformanceResults";
import EditIcon from '@material-ui/icons/Edit';
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import { iconMedium } from "../../css/icons.styles";


const styles = (theme) => ({
  grid : { padding : theme.spacing(2), },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  paper : {
    maxWidth : "90%",
    margin : "auto",
    overflow : "hidden",
  },
});

/**
 *
 * @param {*} props
 */
function MesheryTestProfiles({
  user,
  classes,
  page = 0,
  setPage,
  search = "",
  setSearch,
  sortOrder = "",
  setSortOrder,
  count = 0,
  pageSize = 10,
  setPageSize,
  testProfiles = [],
  setProfileForModal,
  handleDelete,
  showModal,
  fetchTestProfiles
}) {
  const [selectedProfile, setSelectedProfile] = useState();

  useEffect(() => {
    setProfileForModal(selectedProfile);
  }, [selectedProfile]);

  const searchTimeout = useRef(null);



  const columns = [
    {
      name : "name",
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
      },
    },
    {
      name : "endpoints",
      label : "Endpoints",
      options : {
        filter : false,
        sort : true,
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
      name : "last_run",
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
      },
    },
    // {
    //   name: "next_run",
    //   label: "Next Run",
    //   options: {
    //     filter: false,
    //     sort: true,
    //     searchable: true,
    //     customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
    //       return (
    //         <TableCell key={index} onClick={() => sortColumn(index)}>
    //           <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
    //             <b>{column.label}</b>
    //           </TableSortLabel>
    //         </TableCell>
    //       );
    //     },
    //     customBodyRender: function CustomBody(value) {
    //       return <Moment format="LLLL">{value}</Moment>;
    //     },
    //   },
    // },
    {
      name : "updated_at",
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
          return (
            <div>
              <IconButton
                style={iconMedium}
                onClick={(ev) => {
                  ev.stopPropagation();
                  setSelectedProfile(testProfiles[tableMeta.rowIndex]);
                }}
                aria-label="edit"
                // @ts-ignore
                color="rgba(0, 0, 0, 0.54)"
              >
                <EditIcon style={iconMedium}/>
              </IconButton>
              <IconButton
                style={iconMedium}
                onClick={(ev) => {
                  ev.stopPropagation();
                  setSelectedProfile({ ...testProfiles[tableMeta.rowIndex], runTest : true });
                }}
                aria-label="run"
                // @ts-ignore
                color="rgba(0, 0, 0, 0.54)"
              >
                <PlayArrowIcon style={iconMedium} />
              </IconButton>
            </div>
          );
        },
      },
    },
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
    responsive : "standard",
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
          if (sortInfo.length === 2) {
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
    <NoSsr>

      <MUIDataTable
        title={<div className={classes.tableHeader}>Profiles</div>}
        data={testProfiles}
        columns={columns}
        // @ts-ignore
        options={options}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryTestProfiles)));
