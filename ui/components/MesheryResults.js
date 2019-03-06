import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, TableRow, TableCell } from '@material-ui/core';
import MesheryResult from './MesheryResult';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { updateMeshResults } from '../lib/store';
import dataFetch from '../lib/data-fetch';
import MUIDataTable from "mui-datatables";
import CustomToolbarSelect from './CustomToolbarSelect';
import CustomTableFooter from './CustomTableFooter';


const styles = theme => ({
  grid: {
    padding: theme.spacing(2),
  },
});

class MesheryResults extends Component {
    state = {
        page: 0,
        pageMap: {
          0: '',
        },
        count: 10,
        results: [],
        // startKey: '',
    }

    componentDidMount = () => {
      this.fetchResults(0);
    }

    fetchResults = (page) => {
          let self = this;
          let query = '';
          const startKey = (typeof this.state.pageMap[page] !== 'undefined')?this.state.pageMap[page]:'';
          if (startKey !== ''){
            query = `?startKey=${encodeURIComponent(startKey)}`;
          }
          dataFetch(`/api/results${query}`, { 
            credentials: 'same-origin',
            method: 'GET',
            credentials: 'include',
          }, result => {
            // console.log(`received results: ${JSON.stringify(result)}`);
            if (typeof result !== 'undefined'){
              let res = [];
              if (typeof result.results !== 'undefined'){
                res = result.results.map((r) => {
                  return r.runner_results;
                });
              }
              const {pageMap} = this.state;
              let {count} = this.state;
              if (typeof result.last_key !== 'undefined'){
                pageMap[page+1] = result.last_key;
                count += 10;
              }
              this.setState({
                results: res,
                page,
                pageMap,
                count,
              });
              // this.props.updateMeshResults({startKey: result.last_key, results: result.results});
            }
          }, self.handleError);
    }

    handleError = error => {
        // this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Load test did not run successfully with msg: ${error}`});
        console.log(`error fetching results: ${error}`);
      }

    render() {
        const { classes } = this.props; // data here maps to the MesheryResult model

        // const columns = ["RunType", "StartTime", "RequestedQPS", "RequestedDuration"];
        const columns = [
          {
           name: "RunType",
           label: "RunType",
           options: {
            filter: false,
            sort: false,
           }
          },
          {
           name: "StartTime",
           label: "StartTime",
           options: {
            filter: false,
            sort: false,
           }
          },
          {
            name: "RequestedQPS",
            label: "RequestedQPS",
            options: {
              filter: false,
              sort: false,
             }
           },
           {
            name: "RequestedDuration",
            label: "RequestedDuration",
            options: {
              filter: false,
              sort: false,
             }
           },
        ];
        const { results, page, count } = this.state;

        const options = {
          filter: true,
          filterType: 'textField',
          responsive: 'stacked',
          // pagination: true, // default
          // resizableColumns: true,
          // selectableRows: true,
          serverSide: true,
          count: count,
          // rowsPerPage: count,
          rowsPerPageOptions: [10],
          fixedHeader: true,
          page: page,
          print: false,
          download: false,
          onTableChange: (action, tableState) => {

            console.log(action, tableState);
            // a developer could react to change on an action basis or
            // examine the state as a whole and do whatever they want

            switch (action) {
              case 'changePage':
                this.fetchResults(tableState.page);
                break;
            }
          }, 
          customFooter: (count, page, rowsPerPage, changeRowsPerPage, changePage) => {
            return (  
              <CustomTableFooter changePage={changePage} rowsPerPage={rowsPerPage} page={page} count={count} />
            );
          },
          customToolbarSelect: (selectedRows, displayData, setSelectedRows) => (
            <CustomToolbarSelect selectedRows={selectedRows} displayData={displayData} setSelectedRows={setSelectedRows} />
          ),
          expandableRows: true,
          renderExpandableRow: (rowData, rowMeta) => {
            const colSpan = rowData.length + 1;
            return (
              <TableRow>
                <TableCell colSpan={colSpan}>
                  Custom expandable row option. Data: {JSON.stringify(rowData)}
                </TableCell>
              </TableRow>
            );
          },
        };

        return (
            <NoSsr>
            {/* <Grid container spacing={5} className={classes.grid}>
                {results.map((result) => (
                <Grid item xs={12} sm={6}>
                    <MesheryResult key={result.meshery_id} data={result} />
                </Grid>
                ))}
            </Grid> */}
              <MUIDataTable title={"Meshery Results"} data={results} columns={columns} options={options} />
            </NoSsr>
        );
    }
}
MesheryResults.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateMeshResults: bindActionCreators(updateMeshResults, dispatch)
    }
  }
  const mapStateToProps = state => {
    const startKey = state.get("results").get('startKey');
    const results =  state.get("results").get('results').toArray();
    if (typeof results !== 'undefined'){
        return {startKey: startKey, results: results};
    }
    return {};
  }
  
export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
)(MesheryResults));
  