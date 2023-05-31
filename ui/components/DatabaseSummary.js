import React, { useEffect, useState } from 'react';
import { updateProgress } from '../lib/store';
import { IconButton, Paper, Typography, withStyles } from '@material-ui/core';
import CloseIcon from "@material-ui/icons/Close";
import dataFetch from '../lib/data-fetch';
import DataTable from "mui-datatables";
import { connect } from 'react-redux';
import { withSnackbar } from 'notistack';
import { bindActionCreators } from 'redux';

const styles = {
  textCenter : {
    textAlign : "center"
  },
  textEnd : {
    textAlign : "end"
  },
  gapBottom : {
    paddingBottom : "0.5rem"
  }
}

const DatabaseSummary = (props) => {
  const { classes } = props
  const [databaseSummary, setDatabaseSummary] = useState({ tables : [], totalRecords : 0, totalSize : 0, totalTables : 0 })

  const handleError = (msg) => (error) => {
    props.updateProgress({ showProgress : false });
    const self = this;
    props.enqueueSnackbar(`${msg}: ${error}`, {
      variant : "error",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000,
    });
  };

  useEffect(() => {
    dataFetch(
      "/api/system/database",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (typeof result !== "undefined") {
          setDatabaseSummary({
            tables : result?.tables,
            totalRecords : result?.totalRecords,
            totalSize : result?.totalSize,
            totalTables : result?.totalTables
          })
        }
      },
      handleError("Unable to fetch database summary.")
    );
  }, [])

  return (<>
    <Paper elevation={1} style={{ padding : "2rem" }}>
      <Typography variant="h6" className={classes.textCenter}>Database Overview</Typography>
      <Typography className={`${classes.textEnd} ${classes.gapBottom}`}>Total Records : {databaseSummary?.totalRecords}</Typography>
      <Typography className={`${classes.textEnd} ${classes.gapBottom}`}>Total Size : {databaseSummary?.totalSize}</Typography>
      <DataTable
        title={<>
          <Typography>Tables</Typography>
        </>
        }
        data={databaseSummary?.tables}
        options={{
          filter : false,
          selectableRows : "none",
          responsive : "scrollMaxHeight",
          print : false,
          download : false,
          viewColumns : false,
          pagination : false,
          fixedHeader : true,
        }}
        columns={[
          {
            name : "name",
            label : "Name"
          },
          {
            name : "count",
            label : "Count"
          }
        ]}
      />
    </Paper>
  </>)
};


const mapStateToProps = () => ({});

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default withStyles(styles, { withTheme : true })(
  connect(mapStateToProps, mapDispatchToProps)(withSnackbar(DatabaseSummary))
);