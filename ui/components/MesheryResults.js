import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, TableRow, TableCell } from '@material-ui/core';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { updateMeshResults, updateResultsSelection, clearResultsSelection } from '../lib/store';
import dataFetch from '../lib/data-fetch';
import MUIDataTable from "mui-datatables";
import CustomToolbarSelect from './CustomToolbarSelect';
import CustomTableFooter from './CustomTableFooter';
import Moment from 'react-moment';
import MesheryChart from './MesheryChart';


const styles = theme => ({
  grid: {
    padding: theme.spacing(2),
  },
  chartContent: {
    minHeight: window.innerHeight * 0.7,
  },
});

class MesheryResults extends Component {
    state = {
        page: 0,
        pageMap: {
          0: '',
        },
        count: 10,
        pageSize: 10,
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
            //   let res = [];
            //   if (typeof result.results !== 'undefined'){
            //     res = result.results.map((r) => {
            //       return r.runner_results;
            //     });
            //   }
              const {pageMap} = this.state;
              let {count} = this.state;
              if (typeof result.last_key !== 'undefined'){
                pageMap[page+1] = result.last_key;
                count = (page + 1) * self.state.pageSize + self.state.pageSize;
              }
              this.setState({
                // results: res,
                results: result.results,
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
        const { classes, results_selection } = this.props; // data here maps to the MesheryResult model
        const { results, page, count, pageSize } = this.state;

        const resultsForDisplay = [];
        results.forEach((record) => {
          const row = {
            mesh: record.mesh,
            start_time: record.runner_results.StartTime,
            qps: record.runner_results.ActualQPS.toFixed(1),
            duration: (record.runner_results.ActualDuration / 1000000000).toFixed(1),
            threads: record.runner_results.NumThreads,
          }
          if (record.runner_results.DurationHistogram && record.runner_results.DurationHistogram.Percentiles) {
            record.runner_results.DurationHistogram.Percentiles.forEach(({Percentile, Value}) => {
              row['p'+Percentile] = Value.toFixed(3);
            });
          } else {
            row['p50'] = 0;
            row['p75'] = 0;
            row['p90'] = 0;
            row['p99'] = 0;
            row['p99.9'] = 0;
          }
          resultsForDisplay.push(row);
          // console.log(`adding custom row: ${JSON.stringify(row)}`);
        });
        
        const columns = [
          {
           name: "mesh",
           label: "Mesh",
          //  options: {
          //   filter: false,
          //   sort: false,
          //   searchable: false,
          //  }
          },
          {
           name: "start_time",
           label: "StartTime",
           options: {
          //   filter: false,
          //   sort: false,
          //   searchable: false,
              customBodyRender: (value, tableMeta, updateValue) => {
                return (
                  <Moment format="LLLL">{value}</Moment>
                );
              }
            }
          },
          {
            name: "qps",
            label: "QPS",
            // options: {
            //   filter: false,
            //   sort: false,
            //   searchable: false,
            //  }
           },
           {
            name: "duration",
            label: "Duration",
            // options: {
            //   filter: false,
            //   sort: false,
            //   searchable: false,
            //  }
           },
           {
            name: "threads",
            label: "Threads",
            // options: {
            //   filter: false,
            //   sort: false,
            //   searchable: false,
            //  }
           },
           {
            name: "p50",
            label: "P50",
            // options: {
            //   filter: false,
            //   sort: false,
            //   searchable: false,
            //  }
           },
          //  {
          //   name: "p75",
          //   label: "P75",
          //   // options: {
          //   //   filter: false,
          //   //   sort: false,
          //   //   searchable: false,
          //   //  }
          //  },
          //  {
          //   name: "p90",
          //   label: "P90",
          //   // options: {
          //   //   filter: false,
          //   //   sort: false,
          //   //   searchable: false,
          //   //  }
          //  },
          //  {
          //   name: "p99",
          //   label: "P99",
          //   // options: {
          //   //   filter: false,
          //   //   sort: false,
          //   //   searchable: false,
          //   //  }
          //  },
           {
            name: "p99.9",
            label: "P99.9",
            // options: {
            //   filter: false,
            //   sort: false,
            //   searchable: false,
            //  }
           }
        ];

        let rowsSelected = [];
        Object.keys(results_selection).forEach((pg) => {
          if (parseInt(pg) !== page) {
            Object.keys(results_selection[parseInt(pg)]).forEach((ind) => {
              const val = ((parseInt(pg) + 1) * pageSize + parseInt(ind) + 1);
              rowsSelected.push(val);
            });
          } else {
            Object.keys(results_selection[page]).forEach((ind) => {
              rowsSelected.push(ind);
            });
          }
        });
        // console.log(`selected rows after adjustments: ${JSON.stringify(rowsSelected)}`);

        const options = {
          filter: false,
          sort: false,
          search: false,
          filterType: 'textField',
          responsive: 'stacked',
          // pagination: true, // default
          // resizableColumns: true,
          // selectableRows: true,
          serverSide: true,
          count: count,
          // rowsPerPage: count,
          rowsPerPageOptions: [pageSize],
          fixedHeader: true,
          page: page,
          rowsSelected,
          print: false,
          download: false,
          onRowsSelect: (currentRowsSelected, allRowsSelected) => {
            // console.log(`currentRowsSelected: ${JSON.stringify(currentRowsSelected)}`);
            // console.log(`allRowsSelected: ${JSON.stringify(allRowsSelected)}`);
            let res = {};
            allRowsSelected.forEach(({dataIndex}) => {
              if (dataIndex < pageSize) {
                if (typeof res[dataIndex] !== 'undefined'){
                    delete res[dataIndex];
                  } else {
                    res[dataIndex] = results[dataIndex];
                  }  
              }});
              this.props.updateResultsSelection({page, results: res});
          },
          // onRowsDelete: (rowsDeleted) => {
          //   console.log(`delete rows: ${JSON.stringify(rowsDeleted)}`);
          // },
          onTableChange: (action, tableState) => {

            // console.log(action, tableState);
            // a developer could react to change on an action basis or
            // examine the state as a whole and do whatever they want

            switch (action) {
              case 'changePage':
                this.fetchResults(tableState.page);
                break;
            //   case 'rowsSelect':
            //     // TODO: get the updated list of items from state
            //     console.log(`current table state: ${JSON.stringify(tableState.selectedRows.data)}`);
            //     tableState.selectedRows.data = tableState.selectedRows.data.filter(({dataIndex}) => {
            //       return (typeof dataIndex !== 'undefined' && parseInt(dataIndex) < pageSize);
            //     })
            //     console.log(`current table state after cleanup: ${JSON.stringify(tableState.selectedRows.data)}`);
            //     console.log(`selected rows: ${JSON.stringify(tableState.selectedRows.data)}`);
            //     Object.keys(results_selection).forEach((pg) => {
            //       if (parseInt(pg) !== page) {
            //         Object.keys(results_selection[parseInt(pg)]).forEach((ind) => {
            //           const val = ((parseInt(pg) + 1) * pageSize + parseInt(ind) + 1);
            //           tableState.selectedRows.data.push({index: val+'', dataIndex: val+''});
            //           tableState.selectedRows.lookup[val+''] = true;
            //         });
            //       }
            //     });
            //     console.log(`selected rows after adjustments: ${JSON.stringify(tableState.selectedRows.data)}`);
            //     break;
            }
          }, 
          customFooter: (count, page, rowsPerPage, changeRowsPerPage, changePage) => {
            return (  
              <CustomTableFooter changePage={changePage} rowsPerPage={rowsPerPage} page={page} count={count} />
            );
          },
          customToolbarSelect: (selectedRows, displayData, setSelectedRows) => {
            return (
              <CustomToolbarSelect selectedRows={selectedRows} displayData={displayData} setSelectedRows={setSelectedRows} results={results} />
            );
          },
          expandableRows: true,
          renderExpandableRow: (rowData, rowMeta) => {
            const row = results[rowMeta.dataIndex].runner_results;
            const colSpan = rowData.length + 1;
            return (
              <TableRow>
                <TableCell colSpan={colSpan}>
                  <div className={classes.chartContent}>
                    <MesheryChart data={[row]} />
                  </div>
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
              <MUIDataTable title={"Meshery Results"} data={resultsForDisplay} columns={columns} options={options} />
            </NoSsr>
        );
    }
}
MesheryResults.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateMeshResults: bindActionCreators(updateMeshResults, dispatch),
        updateResultsSelection: bindActionCreators(updateResultsSelection, dispatch),
        clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
    }
  }
  const mapStateToProps = state => {
    const startKey = state.get("results").get('startKey');
    const results =  state.get("results").get('results').toArray();
    const results_selection = state.get("results_selection").toObject();
    if (typeof results !== 'undefined'){
        return {startKey: startKey, results: results, results_selection};
    }
    return {results_selection};
  }
  
export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
)(MesheryResults));
  