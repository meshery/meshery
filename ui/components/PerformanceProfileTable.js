// @ts-check
import React, { useState, useEffect, useRef } from "react";
import { withStyles, createMuiTheme, MuiThemeProvider } from "@material-ui/core/styles";
import { NoSsr, TableCell, IconButton, Slide, Paper } from "@material-ui/core";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import MoreHorizIcon from "@material-ui/icons/MoreHoriz";
import ArrowBackIcon from "@material-ui/icons/ArrowBack";
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";
import MesheryResults from "./MesheryResults";
import GenericModal from "./GenericModal";
import MesheryPerformanceComponent from "./MesheryPerformanceComponent";

const styles = (theme) => ({
  grid: {
    padding: theme.spacing(2),
  },
  tableHeader: {
    fontWeight: "bolder",
    fontSize: 18,
  },
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
});

function MesheryTestProfiles({ updateProgress, enqueueSnackbar, closeSnackbar, user, classes }) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [testProfiles, setTestProfiles] = useState([]);
  const [isProfileTableActive, setIsProfileTableActive] = useState(true);
  const [selectedProfile, setSelectedProfile] = useState();
  const [profileForModal, setProfileForModal] = useState();

  function getMuiDatatableTheme() {
    return createMuiTheme({
      overrides: {
        MuiTableRow: {
          root: {
            cursor: "pointer",
          },
        },
      },
    });
  }

  const searchTimeout = useRef(null);

  /**
   * fetch performance profiles when the page loads
   */
  useEffect(() => {
    fetchTestProfiles(page, pageSize, search, sortOrder);
  }, []);

  /**
   * fetchTestProfiles constructs the queries based on the parameters given
   * and fetches the performance profiles
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchTestProfiles(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    const val = "f";
    if (val) {
      updateProgress({ showProgress: true });

      dataFetch(
        `/api/user/performance/profiles${query}`,
        {
          credentials: "include",
        },
        (result) => {
          updateProgress({ showProgress: false });
          if (result) {
            setTestProfiles(result.profiles || []);
            setPage(result.page || 0);
            setPageSize(result.page_size || 0);
            setCount(result.total_count || 0);
          }
        },
        handleError
      );
    }
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

  const columns = [
    {
      name: "name",
      label: "Profile",
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
      name: "endpoints",
      label: "Endpoints",
      options: {
        filter: false,
        sort: true,
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
      name: "last_run",
      label: "Last Run",
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
      name: "Actions",
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
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <IconButton
              onClick={(ev) => {
                ev.stopPropagation();
                setProfileForModal(testProfiles[tableMeta.rowIndex]);
                console.log(testProfiles[tableMeta.rowIndex]);
              }}
              aria-label="more"
              color="inherit"
            >
              <MoreHorizIcon />
            </IconButton>
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
    filter: false,
    sort: !(user && user.user_id === "meshery"),
    search: !(user && user.user_id === "meshery"),
    filterType: "textField",
    responsive: "scrollFullHeight",
    resizableColumns: true,
    serverSide: true,
    selectableRows: "none",
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    onRowClick: (_, rowMetadata) => {
      setIsProfileTableActive(false);
      setSelectedProfile(testProfiles[rowMetadata.rowIndex]);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          fetchTestProfiles(tableState.page, pageSize, search, sortOrder);
          break;
        case "changeRowsPerPage":
          fetchTestProfiles(page, tableState.rowsPerPage, search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchTestProfiles(page, pageSize, tableState.searchText !== null ? tableState.searchText : "", sortOrder);
            }
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
          if (order !== sortOrder) {
            fetchTestProfiles(page, pageSize, search, order);
          }
          break;
      }
    },
  };

  return (
    <NoSsr>
      {isProfileTableActive ? (
        <MuiThemeProvider theme={getMuiDatatableTheme}>
          <Slide direction="left" in={isProfileTableActive} mountOnEnter unmountOnExit timeout={200}>
            <MUIDataTable
              title={<div className={classes.tableHeader}>Profiles</div>}
              data={testProfiles}
              columns={columns}
              // @ts-ignore
              options={options}
            />
          </Slide>
        </MuiThemeProvider>
      ) : null}
      {!isProfileTableActive ? (
        <Slide direction="right" in={!isProfileTableActive} mountOnEnter unmountOnExit timeout={200}>
          <div>
            <MesheryResults
              endpoint={
                // @ts-ignore
                `/api/user/performance/profiles/${selectedProfile?.id}/results`
              }
              customHeader={
                <div
                  style={{
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    width: "fit-content",
                  }}
                >
                  <IconButton aria-label="Back" color="inherit" onClick={() => setIsProfileTableActive(true)}>
                    <ArrowBackIcon />
                  </IconButton>
                  <div className={classes.tableHeader}>
                    Performance Test Results For Profile -{" "}
                    {
                      // @ts-ignore
                      selectedProfile?.name
                    }
                  </div>
                </div>
              }
            />
          </div>
        </Slide>
      ) : null}

      <GenericModal
        open={!!profileForModal}
        Content={
          <Paper className={classes.paper}>
            <MesheryPerformanceComponent
              // @ts-ignore
              loadAsPerformanceProfile
              // @ts-ignore
              performanceProfileID={profileForModal?.id}
              // @ts-ignore
              profileName={profileForModal?.name}
              // @ts-ignore
              meshName={profileForModal?.service_mesh}
              // @ts-ignore
              url={profileForModal?.endpoints?.[0]}
              // @ts-ignore
              qps={profileForModal?.qps}
              // @ts-ignore
              loadGenerator={profileForModal?.load_generators?.[0]}
              // @ts-ignore
              t={profileForModal?.duration}
              // @ts-ignore
              c={profileForModal?.concurrent_request}
              // @ts-ignore
              reqBody={profileForModal?.request_body}
              // @ts-ignore
              headers={profileForModal?.request_headers}
              // @ts-ignore
              cookies={profileForModal?.request_cookies}
              // @ts-ignore
              contentType={profileForModal?.content_type}
            />
          </Paper>
        }
        handleClose={() => {
          fetchTestProfiles(page, pageSize, search, sortOrder);
          setProfileForModal(undefined);
        }}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get("user").toObject(),
  };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryTestProfiles)));
