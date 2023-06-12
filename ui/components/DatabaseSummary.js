import React, { useEffect, useState } from 'react';
import { updateProgress } from '../lib/store';
import { Button, IconButton, Paper, Typography, withStyles } from '@material-ui/core';
import CloseIcon from "@material-ui/icons/Close";
import dataFetch from '../lib/data-fetch';
import DataTable from "mui-datatables";
import { connect } from 'react-redux';
import { withSnackbar } from 'notistack';
import { bindActionCreators } from 'redux';
import PropTypes from "prop-types";
import resetDatabase from './graphql/queries/ResetDatabaseQuery';

const styles = (theme) => ({
  textCenter : {
    textAlign : "center"
  },
  textEnd : {
    textAlign : "end"
  },
  gapBottom : {
    paddingBottom : "0.5rem"
  },
  DBBtn : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5,
    backgroundColor : "#8F1F00",
    "&:hover" : {
      backgroundColor : "#B32700",
    },
  },
  container : {
    display : "flex",
    justifyContent : "center",
    marginTop : theme.spacing(2),
  },
})

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

  const getDatabaseSummary = () => {
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
  }

  useEffect(() => {
    getDatabaseSummary()
  }, [])

  const handleResetDatabase = () => {
    return async () => {
      let responseOfResetDatabase = await props.promptRef.current.show({
        title : "Reset Meshery Database?",
        subtitle : "Are you sure that you want to purge all data?",
        options : ["RESET", "CANCEL"]
      });
      if (responseOfResetDatabase === "RESET") {
        props.updateProgress({ showProgress : true });
        const self = this;
        resetDatabase({
          selector : {
            clearDB : "true",
            ReSync : "true",
            hardReset : "true",
          },
          k8scontextID : ""
        }).subscribe({
          next : (res) => {
            props.updateProgress({ showProgress : false });
            if (res.resetStatus === "PROCESSING") {
              props.enqueueSnackbar(`Database reset successful.`, {
                variant : "success",
                action : (key) => (
                  <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                    <CloseIcon />
                  </IconButton>
                ),
                autohideduration : 3000,
              })

              getDatabaseSummary()
            }
          },
          error : handleError("Database is not reachable, try restarting server.")
        });
      }
    }
  }

  return (<>
      <DataTable
        title={<>
          <Typography>Database Overview</Typography>
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
          customToolbar: () => (<Button
            type="submit"
            variant="contained"
            color="primary"
            size="medium"
            onClick={handleResetDatabase()}
            className={classes.DBBtn}
            data-cy="btnResetDatabase"
          >
            <Typography align="center" variant="subtitle2"> RESET DATABASE </Typography>
          </Button>)
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
  </>)
};


const mapStateToProps = () => ({});

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

DatabaseSummary.propTypes = {
  promptRef : PropTypes.object.isRequired
}

export default withStyles(styles, { withTheme : true })(
  connect(mapStateToProps, mapDispatchToProps)(withSnackbar(DatabaseSummary))
);