import {
  NoSsr,
  TableCell,
  // Button,
  Tooltip,
  Link,
  // FormControl
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
// import EditIcon from "@material-ui/icons/Edit";
// import YoutubeSearchedForIcon from '@mui/icons-material/YoutubeSearchedFor';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import React, { useEffect, useRef, useState } from "react";
import Moment from "react-moment";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateProgress } from "../../lib/store";
import { /* Avatar, */ Chip, /* FormControl, */ } from "@mui/material";
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import classNames from "classnames";
// import ReactSelectWrapper from "../ReactSelectWrapper";
import dataFetch from "../../lib/data-fetch";
import LaunchIcon from '@mui/icons-material/Launch';
import TableRow from '@mui/material/TableRow';
import { useNotification } from "../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../lib/event-types";
import CustomColumnVisibilityControl from "../../utils/custom-column";
import SearchBar from "../../utils/custom-search";
import ResponsiveDataTable from "../../utils/data-table";
import useStyles from "../../assets/styles/general/tool.styles";
// import MySelectComponent from "./environment";

const styles = (theme) => ({
  grid : { padding : theme.spacing(2) },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  muiRow : {
    "& .MuiTableRow-root" : {
      cursor : "pointer",
    },
    "& .MuiTableCell-root" : {
      textTransform : "capitalize",
    },
  },
  createButton : {
    display : "flex",
    justifyContent : "flex-start",
    alignItems : "center",
    whiteSpace : "nowrap",
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : "1rem",
  },
  statusCip : {
    minWidth : "120px !important",
    maxWidth : "max-content !important",
    display : "flex !important",
    justifyContent : "flex-start !important",
    textTransform : "capitalize",
    borderRadius : "3px !important",
    padding : "6px 8px",
    "& .MuiChip-label" : {
      paddingTop : "3px",
      fontWeight : "400",
    },
    "&:hover" : {
      boxShadow : "0px 1px 2px 0px rgba(0, 0, 0, 0.25)"
    }
  },
  capitalize : {
    textTransform : "capitalize",
  },
  ignored : {
    "& .MuiChip-label" : {
      color : `${theme.palette.secondary.default}`,
    },
    background : `${theme.palette.secondary.default}30 !important`,
    "& .MuiSvgIcon-root" : {
      color : `${theme.palette.secondary.default} !important`,
    }
  },
  connected : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.success,
    },
    background : `${theme.palette.secondary.success}30 !important`,
    "& .MuiSvgIcon-root" : {
      color : `${theme.palette.secondary.success} !important`,
    }
  },
  registered : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.primary,
    },
    background : `${theme.palette.secondary.primary}30 !important`,
    "& .MuiSvgIcon-root" : {
      color : `${theme.palette.secondary.primary} !important`,
    }
  },
  discovered : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.warning,
    },
    background : `${theme.palette.secondary.warning}30 !important`,
    "& .MuiSvgIcon-root" : {
      color : `${theme.palette.secondary.warning} !important`,
    }
  },
  deleted : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.error,
    },
    background : `${theme.palette.secondary.lightError}30 !important`,
    "& .MuiSvgIcon-root" : {
      color : `${theme.palette.secondary.error} !important`,
    }
  },
  expandedRows : {
    background : `${theme.palette.secondary.default}10`
  }
});

const ACTION_TYPES = {
  FETCH_CONNECTIONS : {
    name : "FETCH_CONNECTIONS",
    error_msg : "Failed to fetch connections"
  },
};

function Connections({ classes, updateProgress }) {
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [connections, setConnections] = useState([]);
  const [search,setSearch] = useState("");
  const { notify } = useNotification()
  const StyleClass = useStyles();

  const searchTimeout = useRef(null);

  const status = (value) => {
    switch (value) {
      case 'ignored':
        return <Chip className={classNames(classes.statusCip, classes.ignored)} avatar={<RemoveCircleIcon />} label={value} />
      case 'connected':
        return <Chip className={classNames(classes.statusCip, classes.connected)} avatar={<CheckCircleIcon />} label={value} />
      case 'REGISTERED':
        return <Chip className={classNames(classes.statusCip, classes.registered)} avatar={<AssignmentTurnedInIcon />} label={value.toLowerCase()} />
      case 'discovered':
        return <Chip className={classNames(classes.statusCip, classes.discovered)} avatar={<ExploreIcon />} label={value} />
      case 'deleted':
        return <Chip className={classNames(classes.statusCip, classes.deleted)} avatar={<DeleteForeverIcon />} label={value} />
      default:
        return "-"
    }
  }

  const columns = [
    {
      name : "name",
      label : "Element",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : (value, tableMeta) => {
          return (
            <Tooltip title={connections ? connections[tableMeta.rowIndex]?.metadata?.server_location : ''} placement="top" >
              <Link href={connections ? connections[tableMeta.rowIndex]?.metadata?.server_location : ''} target="_blank">
                {value}
                <sup>
                  <LaunchIcon sx={{ fontSize : "12px" }} />
                </sup>
              </Link>
            </Tooltip>
          );
        }
      }
    },
    {
      name : "type",
      label : "Type",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        // customBodyRender : function CustomBody(value) {
        //   return (
        //     <FormControl sx={{ m : 1, minWidth : 150, maxWidth : 150 }} size="small">
        //       <MySelectComponent />
        //     </FormControl>
        //   );
        // },
      },
    },
    {
      name : "sub_type",
      label : "Sub Type",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        // customBodyRender : function CustomBody(value) {
        //   return <Chip avatar={<Avatar>M</Avatar>} label={value} />;
        // },
      },
    },
    {
      name : "kind",
      label : "Kind",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name : "updated_at",
      label : "Updated At",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return (
            <Tooltip title={<Moment startOf="day" format="LLL">{value}</Moment>} placement="top" arrow interactive >
              <Moment format="LL">{value}</Moment>
            </Tooltip>
          );
        },
      },
    },
    {
      name : "discoverd_at",
      label : "Discovered At",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return (
            <Tooltip title={<Moment startOf="day" format="LLL">{value}</Moment>} placement="top" arrow interactive >
              <Moment format="LL">{value}</Moment>
            </Tooltip>
          );
        },
      },
    },
    {
      name : "status",
      label : "Status",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return (
            status(value)
          );
        },
      },
    }
  ];

  // const handleChange = () => {
  //   // Select change
  // }

  const options = {
    filter : false,
    viewColumns : false,
    search : false,
    responsive : "standard",
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
        text : "connection(s) selected",
      },
    },
    enableNestedDataAccess : '.',
    onSearchClose : () => {
      setSearch("")
    },
    onTableChange : (action, tableState) => {
      switch (action) {
        case "changePage":
          getConnections(tableState.page.toString(), pageSize.toString());
          break;
        case "changeRowsPerPage":
          getConnections(page.toString(), tableState.rowsPerPage.toString());
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              getConnections(page, pageSize, tableState.searchText !== null ? tableState.searchText : "");
              setSearch(tableState.searchText);
            }
          }, 500);
          break;
      }
    },
    expandableRows : true,
    expandableRowsHeader : false,
    expandableRowsOnClick : true,
    isRowExpandable : () => {
      return true;
    },
    // rowsExpanded : [0, 1],
    renderExpandableRow : (_, tableMeta) => {
      return (
        <TableRow>
          <TableCell>
          </TableCell>
          <TableCell colSpan={2}>
            <b>Server Build SHA:</b> {connections ? connections[tableMeta.rowIndex]?.metadata?.server_build_sha : '-'}
          </TableCell>
          <TableCell colSpan={2}>
            <b>Server Version:</b> {connections ? connections[tableMeta.rowIndex]?.metadata?.server_version : '-'}
          </TableCell>
          <TableCell colSpan={2}>
          </TableCell>
        </TableRow>
      );
    },
  };


  /**
   * fetch connections when the page loads
   */
  useEffect(() => {
    getConnections(page, pageSize,search)
  }, [page, pageSize, search]);

  const getConnections = (page, pageSize,search) => {
    if (!search) search = "";
    dataFetch(
      `/api/integrations/connections?page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(search)}`,
      {
        credentials : "include",
        method : "GET",
      },
      (res) => {
        setConnections(res?.connections)
        setPage(res?.page || 0)
        setCount(res?.total_count || 0)
        setPageSize(res?.page_size || 0)
      },
      handleError(ACTION_TYPES.FETCH_CONNECTIONS)
    );
  }

  const handleError = (action) => (error) => {
    updateProgress({ showProgress : false });
    notify({ message : `${action.error_msg}: ${error}`, event_type : EVENT_TYPES.ERROR, details : error.toString() })
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach(col => {
      initialVisibility[col.name] = col.options?.display !== false;
    });
    return initialVisibility;
  });

  return (
    <>
      <NoSsr>
        <div className={StyleClass.toolWrapper} >
          <div className={classes.createButton}>
            <div>
              {/* <Button
                aria-label="Rediscover"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={() => {}}
                style={{ marginRight : "2rem" }}
              >
                <YoutubeSearchedForIcon style={iconMedium} />
                Rediscover
              </Button> */}
            </div>
          </div>
          <div
            className={classes.searchAndView}
            style={{
              display : "flex",
            }}
          >
            {/* <Button
              aria-label="Edit"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              onClick={() => {}}
              style={{ marginRight : "0.5rem" }}
            >
              <EditIcon style={iconMedium} />
            </Button> */}

            {/* <Button
              aria-label="Delete"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              onClick={() => {}}
              style={{ background : "#8F1F00" }}
            >
              <DeleteForeverIcon style={iconMedium} />
              Delete
            </Button> */}

            <SearchBar
              onSearch={(value) => {
                setSearch(value);
                getConnections(page, pageSize, value);
              }}
              placeholder="Search"
            />

            <CustomColumnVisibilityControl
              columns={columns}
              customToolsProps={{ columnVisibility, setColumnVisibility }}
            />

          </div>
        </div>
        <ResponsiveDataTable
          data={connections}
          columns={columns}
          // @ts-ignore
          options={options}
          className={classes.muiRow}
          tableCols={tableCols}
          updateCols={updateCols}
          columnVisibility={columnVisibility}
        />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch) });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts") };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(Connections));
