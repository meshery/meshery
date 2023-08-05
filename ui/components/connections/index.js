import {
  NoSsr,
  TableCell,
  Button,
  Tooltip,
} from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import EditIcon from "@material-ui/icons/Edit";
import CloseIcon from "@material-ui/icons/Close";
import YoutubeSearchedForIcon from '@mui/icons-material/YoutubeSearchedFor';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import MUIDataTable from "mui-datatables";
import { withSnackbar } from "notistack";
import React, { useEffect, useState } from "react";
import Moment from "react-moment";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateProgress } from "../../lib/store";
import { iconMedium } from "../../css/icons.styles";
import { Avatar, Chip, FormControl, IconButton } from "@mui/material";
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import ExploreIcon from '@mui/icons-material/Explore';
import RemoveCircleIcon from '@mui/icons-material/RemoveCircle';
import classNames from "classnames";
import ReactSelectWrapper from "../ReactSelectWrapper";
import dataFetch from "../../lib/data-fetch";

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
  },
  createButton : {
    display : "flex",
    justifyContent : "flex-start",
    alignItems : "center",
    whiteSpace : "nowrap",
  },
  topToolbar : {
    margin : "2rem auto",
    display : "flex",
    justifyContent : "space-between",
    paddingLeft : "1rem",
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : "1rem",
  },
  statusCip : {
    "& .MuiChip-label" : {
      paddingTop : "3px",
      fontWeight : "400",
    },
    borderRadius : "3px !important",
    display : "flex",
    width : "117px",
    padding : "6px 8px",
    alignItems : "center",
    gap : "5px",
  },
  ignored : {
    "& .MuiChip-label" : {
      color : `${theme.palette.secondary.default}`,
    },
    background : `${theme.palette.secondary.default}15 !important`,
  },
  connected : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.success,
    },
    background : `${theme.palette.secondary.success}15 !important`,
  },
  registered : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.primary,
    },
    background : `${theme.palette.secondary.primary}15 !important`,
  },
  discovered : {
    "& .MuiChip-label" : {
      color : theme.palette.secondary.warning,
    },
    background : `${theme.palette.secondary.warning}15 !important`,
  },
});

const ACTION_TYPES = {
  FETCH_CONNECTIONS : {
    name : "FETCH_CONNECTIONS",
    error_msg : "Failed to fetch connections"
  },
};

function Connections({ classes, updateProgress, closeSnackbar, enqueueSnackbar }) {
  const [page] = useState(0);
  const [count] = useState(0);
  const [pageSize] = useState(10);
  const [connections, setConnections] = useState([]);


  const status = (value) => {
    switch (value) {
      case 'Ignored':
        return <Chip className={classNames(classes.statusCip, classes.ignored)} avatar={<RemoveCircleIcon style={{ color : "#51636B" }} />} label={value} />
      case 'Connected':
        return <Chip className={classNames(classes.statusCip, classes.connected)} avatar={<CheckCircleIcon style={{ color : "#00B39F" }}/>} label={value} />
      case 'Registered':
        return <Chip className={classNames(classes.statusCip, classes.registered)} avatar={<AssignmentTurnedInIcon style={{ color : "#477E96" }} />} label={value} />
      case 'Discovered':
        return <Chip className={classNames(classes.statusCip, classes.discovered)} avatar={<ExploreIcon style={{ color : "#EBC017" }} />} label={value} />
      default:
        return "-"
    }
  }

  const columns = [
    {
      name : "element",
      label : "Element",
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
      name : "cluster",
      label : "Cluster",
      options : {
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(/* value */) {
          return <Chip avatar={<Avatar>M</Avatar>} label={'value'} />;
        },
      },
    },
    {
      name : "environment",
      label : "Environment",
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
            <FormControl sx={{ m : 1, minWidth : 120 }} size="small">
              <ReactSelectWrapper
                onChange={handleChange}
                options={[{ value : "environment 1", label : "environment 1" }, { value : "environment 2", label : "environment 2" }]}
                value={{ value : value, label : value }}
              />
            </FormControl>
          );
        },
      },
    },
    {
      name : "updated_at",
      label : "Update At",
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
      label : "Discoverd At",
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
      name : "asdf",
      label : "asdf",
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
    },
  ];

  const handleChange = () => {
    // Select change
  }

  const options = {
    filter : false,
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
    }
  };

  useEffect(() => {
    getConnections()
  },[])

  const getConnections = () => {
    dataFetch(
      `/api/integrations/connections`,
      {
        credentials : "include",
        method : "GET",
      },
      (res) => {
        console.log("🚀 ~ file: index.js:306 ~ getConnections ~ res:", res)
        setConnections(res?.connection)
      },
      handleError(ACTION_TYPES.FETCH_CONNECTIONS)
    );
  }

  const handleError = (action) => (error) => {
    updateProgress({ showProgress : false });

    enqueueSnackbar(`${action.error_msg}: ${error}`, {
      variant : "error",
      action : function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon style={iconMedium} />
          </IconButton>
        );
      },
      autoHideDuration : 8000,
    });
  };

  return (
    <>
      <NoSsr>
        <div className={classes.topToolbar}>
          <div className={classes.createButton}>
            <div>
              <Button
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
              </Button>
            </div>
          </div>
          <div
            className={classes.searchAndView}
            style={{
              display : "flex",
              alignItems : "center",
              justifyContent : "flex-end",
              height : "5ch",
            }}
          >
            <Button
              aria-label="Edit"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              onClick={() => {}}
              style={{ marginRight : "0.5rem" }}
            >
              <EditIcon style={iconMedium} />
            </Button>
            <Button
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
            </Button>
          </div>
        </div>
        <MUIDataTable
          data={connections}
          columns={columns}
          // @ts-ignore
          options={options}
          className={classes.muiRow}
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
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(Connections)));

