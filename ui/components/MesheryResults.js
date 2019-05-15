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
        search: '',
        sortOrder: '',
        // pageMap: {
        //   0: '',
        // },
        count: 0,
        pageSize: 10,
        results: [],
        // startKey: '',
    }

    componentDidMount = () => {
      const {page, pageSize, search, sortOrder} = this.state;
      this.fetchResults(page, pageSize, search, sortOrder);
    }

    fetchResults = (page, pageSize, search, sortOrder) => {
          let self = this;
          let query = '';
          if (typeof search === 'undefined' || search === null) {
            search = '';
          }
          if (typeof sortOrder === 'undefined' || sortOrder === null) {
            sortOrder = '';
          }
          query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(sortOrder)}`;
          dataFetch(`/api/results${query}`, { 
            credentials: 'same-origin',
            method: 'GET',
            credentials: 'include',
          }, result => {
            // console.log(`received results: ${JSON.stringify(result)}`);
            if (typeof result !== 'undefined'){
              this.setState({
                results: result.results,
                search: search,
                sortOrder: sortOrder,
                page: result.page,
                pageSize: result.page_size,
                count: result.total_count,
              });
            }
          }, self.handleError);
    }

    handleError = error => {
        console.log(`error fetching results: ${error}`);
      }

    render() {
        const { classes, results_selection } = this.props;
        const { results, page, count, pageSize, search, sortOrder } = this.state;
        const self = this;
        const resultsForDisplay = [];
        results.forEach((record) => {
          const row = {
            name: record.name,
            mesh: record.mesh,
            test_start_time: record.runner_results.StartTime,
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
            name: "name",
            label: "Name",
            options: {
             filter: false,
             sort: true,
             searchable: true,
            }
           },
          {
           name: "mesh",
           label: "Mesh",
           options: {
            filter: false,
            sort: true,
            searchable: true,
           }
          },
          {
           name: "test_start_time",
           label: "StartTime",
           options: {
              filter: false,
              sort: true,
              searchable: false,
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
            options: {
              filter: false,
              sort: false,
              searchable: false,
             }
           },
           {
            name: "duration",
            label: "Duration",
            options: {
              filter: false,
              sort: false,
              searchable: false,
             }
           },
           {
            name: "threads",
            label: "Threads",
            options: {
              filter: false,
              sort: false,
              searchable: false,
             }
           },
           {
            name: "p50",
            label: "P50",
            options: {
              filter: false,
              sort: false,
              searchable: false,
             }
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
            options: {
              filter: false,
              sort: false,
              searchable: false,
             }
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
          sort: true,
          search: true,
          filterType: 'textField',
          responsive: 'stacked',
          // resizableColumns: true,
          // selectableRows: true,
          serverSide: true,
          count: count,
          rowsPerPage: pageSize,
          rowsPerPageOptions: [10, 20, 25],
          fixedHeader: true,
          page: page,
          rowsSelected,
          print: false,
          download: false,
          onRowsSelect: (currentRowsSelected, allRowsSelected) => {
            // const rs = self.props.results_selection;
            
            let res = {};
            allRowsSelected.forEach(({dataIndex}) => {
              if (dataIndex < self.state.pageSize) {
                if (typeof res[dataIndex] !== 'undefined'){
                    delete res[dataIndex];
                  } else {
                    res[dataIndex] = self.state.results[dataIndex];
                  }  
              }});

              // let rsk = 0;
              // Object.keys(rs).forEach((k1) => {
              //     const pg = parseInt(k1);
              //     if(pg !== page){ // skipping count for this page
              //       Object.keys(rs[k1]).forEach((k2) => {
              //         rsk++;
              //       });
              //     }
              // })

              // if (rsk + Object.keys(res).length > 4) { // count from other pages + this page
              //   return null;
              // }

              this.props.updateResultsSelection({page, results: res});
          },
          onTableChange: (action, tableState) => {

            // console.log(action, tableState);
          
            switch (action) {
              case 'changePage':
                this.fetchResults(tableState.page, self.state.pageSize, self.state.search, self.state.sortOrder);
                break;
              case 'changeRowsPerPage':
                this.fetchResults(self.state.page, tableState.rowsPerPage, self.state.search, self.state.sortOrder);
                break;
              case 'search':
                if (self.searchTimeout) {
                  clearTimeout(self.searchTimeout);
                }
                self.searchTimeout = setTimeout(function(){
                  if (self.state.search !== tableState.searchText){
                    self.fetchResults(self.state.page, self.state.pageSize, tableState.searchText !== null?tableState.searchText:'', self.state.sortOrder);
                  }
                }, 500);
                break;
              case 'sort':
                const sortInfo = tableState.announceText.split(' : ');
                let order = 'asc';
                if (sortInfo.length == 2) {
                  if (sortInfo[1] === 'descending'){
                    order = 'desc';
                  }
                }
                if (order !== sortOrder){
                  this.fetchResults(self.state.page, self.state.pageSize, self.state.search, columns[tableState.activeColumn].name+ ' ' + order);
                }
                break;
            }
          }, 
          customToolbarSelect: (selectedRows, displayData, setSelectedRows) => {
            return (
              <CustomToolbarSelect selectedRows={selectedRows} displayData={displayData} setSelectedRows={setSelectedRows} results={self.state.results} />
            );
          },
          expandableRows: true,
          renderExpandableRow: (rowData, rowMeta) => {
            const row = self.state.results[rowMeta.dataIndex].runner_results;
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
        // updateMeshResults: bindActionCreators(updateMeshResults, dispatch),
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
  