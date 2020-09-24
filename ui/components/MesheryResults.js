import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, TableRow, TableCell, IconButton, Paper, Table, TableBody, TableHead, 
} from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MUIDataTable from 'mui-datatables';
import Moment from 'react-moment';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import {
  updateResultsSelection, clearResultsSelection, updateProgress,
} from '../lib/store';
import dataFetch from '../lib/data-fetch';
import CustomToolbarSelect from './CustomToolbarSelect';
import MesheryChart from './MesheryChart';
import GrafanaCustomCharts from './GrafanaCustomCharts';
import MesheryResultDialog from './MesheryResultDialog';


const styles = (theme) => ({
  grid: {
    padding: theme.spacing(2),
  },
  table: {
    margin: theme.spacing(10),
  },
  chartContent: {
    // minHeight: window.innerHeight * 0.7,
  },
});

class MesheryResults extends Component {
  constructor(props) {
    super(props);
    // const {results_selection} = props;
    this.state = {
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

      // results_selection,

      selectedRowData: null,

      smi_page: 0,
      smi_pageSize: 10,
      smi_search: '',
      smi_sortOrder: '',
      smi_results: [],
    };
  }

  // static getDerivedStateFromProps(props, state){
  //   const { results_selection } = props;
  //   return { results_selection };
  // }

    componentDidMount = () => {
      const {
        page, pageSize, search, sortOrder, smi_page, smi_pageSize, smi_search, smi_sortOrder,
      } = this.state;
      this.fetchResults(page, pageSize, search, sortOrder);
      this.fetchSMIResults(smi_page, smi_pageSize, smi_search, smi_sortOrder);
    }

    fetchSMIResults= ( page, pageSize, search, sortOrder) => {
      const self = this;
      let query = '';
      if (typeof search === 'undefined' || search === null) {
        search = '';
      }
      if (typeof sortOrder === 'undefined' || sortOrder === null) {
        sortOrder = '';
      }
      query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(sortOrder)}`;
      // console.log(`/api/smi/results${query}`)
      dataFetch(`/api/smi/results${query}`, {
        credentials: 'same-origin',
        method: 'GET',
        credentials: 'include',
      }, (result) => {
        console.log(result)
        if (typeof result !== 'undefined' && result.results) {
          self.setState({smi_results: result});
        }
      }, self.handleError('Could not fetch SMI results.'));
    }

    fetchResults = (page, pageSize, search, sortOrder) => {
      const self = this;
      let query = '';
      if (typeof search === 'undefined' || search === null) {
        search = '';
      }
      if (typeof sortOrder === 'undefined' || sortOrder === null) {
        sortOrder = '';
      }
      query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(sortOrder)}`;
      self.props.updateProgress({ showProgress: true });
      dataFetch(`/api/perf/results${query}`, {
        credentials: 'same-origin',
        method: 'GET',
        credentials: 'include',
      }, (result) => {
        console.log("Results API",`/api/perf/results${query}`)
        self.props.updateProgress({ showProgress: false });
        // console.log(`received results: ${JSON.stringify(result)}`);
        if (typeof result !== 'undefined') {
          this.setState({
            results: result.results,
            search,
            sortOrder,
            page: result.page,
            pageSize: result.page_size,
            count: result.total_count,
          });
        }
      }, self.handleError);
    }

    handleError = (error) => {
      this.props.updateProgress({ showProgress: false });
      // console.log(`error fetching results: ${error}`);
      const self = this;
      this.props.enqueueSnackbar(`There was an error fetching results: ${error}`, {
        variant: 'error',
        action: (key) => (
          <IconButton
            key="close"
            aria-label="Close"
            color="inherit"
            onClick={() => self.props.closeSnackbar(key)}
          >
            <CloseIcon />
          </IconButton>
        ),
        autoHideDuration: 8000,
      });
    }

    resetSelectedRowData() {
      const self = this;
      return () => {
        self.setState({ selectedRowData: null });
      };
    }

    render() {
      const { classes, results_selection, user } = this.props;
      const {
        results, page, count, pageSize, selectedRowData, smi_results,
      } = this.state;
      const self = this;
      const resultsForDisplay = [];
      const smi_resultsForDisplay = [];
      results.forEach((record) => {
        const row = {
          name: record.name,
          mesh: record.mesh,
          test_start_time: record.runner_results.StartTime,
          qps: record.runner_results.ActualQPS.toFixed(1),
          duration: (record.runner_results.ActualDuration / 1000000000).toFixed(1),
          threads: record.runner_results.NumThreads,
        };
        if (record.runner_results.DurationHistogram && record.runner_results.DurationHistogram.Percentiles) {
          record.runner_results.DurationHistogram.Percentiles.forEach(({ Percentile, Value }) => {
            row[(`p${Percentile}`).replace('.', '_')] = Value.toFixed(3);
          });
        } else {
          row.p50 = 0;
          row.p75 = 0;
          row.p90 = 0;
          row.p99 = 0;
          row.p99_9 = 0;
        }
        resultsForDisplay.push(row);
        // console.log(`adding custom row: ${JSON.stringify(row)}`);
      });

      const smi_columns = ["ID","Date", "Service Mesh","Service Mesh Version", "% Passed","Status"];
      if(smi_results&&smi_results.results) {
        smi_results.results.map((val) => {
          smi_resultsForDisplay.push([val.id,val.date,val.mesh_name,val.mesh_version,val.passing_percentage,val.status])
        }) 
      }
      const smi_options = {
        filterType: 'textField',
        expandableRows: true,
        selectableRows: false,
        rowsPerPage: pageSize,
        rowsPerPageOptions: [10, 20, 25],
        fixedHeader: true,
        renderExpandableRow: (rowData, rowMeta) => {
          console.log("Rox Data",rowData,rowMeta)
          const column = ["SMI Specification","Assertions", "Time","SMI Version", "Capability", "Result", "Reason"]
          const data = smi_results.results[rowMeta.dataIndex].more_details.map((val) => {
            return [val.smi_specification,val.assertions,val.time,"alpha1/v1",val.capability,val.status,val.reason] 
          })
          const colSpan = rowData.length + 1
          return (
            <TableRow>
              <TableCell colSpan={colSpan}>
                <Paper elevation={4} >
                  <Table aria-label="a dense table">
                    <TableHead>
                      <TableRow>
                        {column.map((val) => (<TableCell colSpan={colSpan}>{val}</TableCell>))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {data.map((row) => (
                        <TableRow >
                          {row.map(val => (<TableCell colSpan={colSpan}>{val}</TableCell>))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </Paper>
              </TableCell>
            </TableRow>
          );
        },
        onTableChange: (action, tableState) => {
          const sortInfo = tableState.announceText? tableState.announceText.split(' : '):[];
          let order='';
          if(tableState.activeColumn){
            order = `${columns[tableState.activeColumn].name} desc`;
          }

          switch (action) {
            case 'changePage':
              self.fetchSMIResults(tableState.page, self.state.smi_pageSize, self.state.smi_search, self.state.smi_sortOrder);
              break;
            case 'changeRowsPerPage':
              self.fetchSMIResults(self.state.smi_page, tableState.rowsPerPage, self.state.smi_search, self.state.smi_sortOrder);
              break;
            case 'search':
              if (self.searchTimeout) {
                clearTimeout(self.searchTimeout);
              }
              self.searchTimeout = setTimeout(() => {
                if (self.state.search !== tableState.searchText) {
                  self.fetchSMIResults(self.state.smi_page, self.state.smi_pageSize, tableState.searchText !== null ? tableState.searchText : '', self.state.smi_sortOrder);
                }
              }, 500);
              break;
            case 'sort':
              if (sortInfo.length == 2) {
                if (sortInfo[1] === 'ascending') {
                  order = `${columns[tableState.activeColumn].name} asc`;
                } else {
                  order = `${columns[tableState.activeColumn].name} desc`;
                }
              }
              if (order !== this.state.sortOrder) {
                self.fetchSMIResults(self.state.smi_page, self.state.smi_pageSize, self.state.smi_search, order);
              }
              break;
          }
        },
      }

      const columns = [
        {
          name: 'name',
          label: 'Name',
          options: {
            filter: false,
            sort: true,
            searchable: true,
          },
        },
        {
          name: 'mesh',
          label: 'Mesh',
          options: {
            filter: false,
            sort: true,
            searchable: true,
          },
        },
        {
          name: 'test_start_time',
          label: 'StartTime',
          options: {
            filter: false,
            sort: true,
            searchable: false,
            customBodyRender: (value) => (
              <Moment format="LLLL">{value}</Moment>
            ),
          },
        },
        {
          name: 'qps',
          label: 'QPS',
          options: {
            filter: false,
            sort: false,
            searchable: false,
          },
        },
        {
          name: 'duration',
          label: 'Duration',
          options: {
            filter: false,
            sort: false,
            searchable: false,
          },
        },

        {
          name: 'p50',
          label: 'P50',
          options: {
            filter: false,
            sort: false,
            searchable: false,
          },
        },

        {
          name: 'p99_9',
          label: 'P99.9',
          options: {
            filter: false,
            sort: false,
            searchable: false,
          },
        },
        {
          name: 'Details',
          options: {
            filter: false,
            sort: false,
            searchable: false,
            customBodyRender: (value, tableMeta) => (
              <IconButton
                aria-label="more"
                color="inherit"
                onClick={() => self.setState({ selectedRowData: self.state.results[tableMeta.rowIndex] })}
              >
                <MoreHorizIcon />
              </IconButton>
            ),
          },
        },
      ];

      columns.forEach((column, idx) => {
        if(column.name === this.state.sortOrder.split(' ')[0]) {
          columns[idx].options.sortDirection = this.state.sortOrder.split(' ')[1];
        }
      })

      const rowsSelected = [];
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
        sort: !(user && user.user_id === 'meshery'),
        search: !(user && user.user_id === 'meshery'),
        filterType: 'textField',
        responsive: 'scrollFullHeight',
        // resizableColumns: true,
        // selectableRows: true,
        serverSide: true,
        count,
        rowsPerPage: pageSize,
        rowsPerPageOptions: [10, 20, 25],
        fixedHeader: true,
        page,
        rowsSelected,
        print: false,
        download: false,
        onRowsSelect: (currentRowsSelected, allRowsSelected) => {
          // const rs = self.props.results_selection;
          const res = {};
          allRowsSelected.forEach(({ dataIndex }) => {
            if (dataIndex < self.state.pageSize) {
              if (typeof res[dataIndex] !== 'undefined') {
                delete res[dataIndex];
              } else {
                res[dataIndex] = self.state.results[dataIndex];
              }
            }
          });

          self.props.updateResultsSelection({ page, results: res });
        },

        onTableChange: (action, tableState) => {
          const sortInfo = tableState.announceText? tableState.announceText.split(' : '):[];
          let order='';
          if(tableState.activeColumn){
            order = `${columns[tableState.activeColumn].name} desc`;
          }

          switch (action) {
            case 'changePage':
              self.fetchResults(tableState.page, self.state.pageSize, self.state.search, self.state.sortOrder);
              break;
            case 'changeRowsPerPage':
              self.fetchResults(self.state.page, tableState.rowsPerPage, self.state.search, self.state.sortOrder);
              break;
            case 'search':
              if (self.searchTimeout) {
                clearTimeout(self.searchTimeout);
              }
              self.searchTimeout = setTimeout(() => {
                if (self.state.search !== tableState.searchText) {
                  self.fetchResults(self.state.page, self.state.pageSize, tableState.searchText !== null ? tableState.searchText : '', self.state.sortOrder);
                }
              }, 500);
              break;
            case 'sort':
              if (sortInfo.length == 2) {
                if (sortInfo[1] === 'ascending') {
                  order = `${columns[tableState.activeColumn].name} asc`;
                } else {
                  order = `${columns[tableState.activeColumn].name} desc`;
                }
              }
              if (order !== this.state.sortOrder) {
                self.fetchResults(self.state.page, self.state.pageSize, self.state.search, order);
              }
              break;
          }
        },
        customToolbarSelect: (selectedRows, displayData, setSelectedRows) => (
          <CustomToolbarSelect selectedRows={selectedRows} displayData={displayData} setSelectedRows={setSelectedRows} results={self.state.results} />
        ),
        expandableRows: true,
        renderExpandableRow: (rowData, rowMeta) => {
          const row = self.state.results[rowMeta.dataIndex].runner_results;
          const boardConfig = self.state.results[rowMeta.dataIndex].server_board_config;
          const serverMetrics = self.state.results[rowMeta.dataIndex].server_metrics;
          const startTime = new Date(row.StartTime);
          const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
          const colSpan = rowData.length + 1;
          return (
            <TableRow>
              <TableCell colSpan={colSpan}>
                <div className={classes.chartContent}>
                  <MesheryChart data={[row]} hideTitle />
                </div>
                {boardConfig && boardConfig !== null && Object.keys(boardConfig).length > 0 && (
                  <div>
                    <GrafanaCustomCharts
                      boardPanelConfigs={[boardConfig]}
                      boardPanelData={[serverMetrics]}
                      startDate={startTime}
                      from={startTime.getTime().toString()}
                      endDate={endTime}
                      to={endTime.getTime().toString()}
                      liveTail={false}
                    />
                  </div>
                )}
              </TableCell>
            </TableRow>
          );
        },
      };

      return (
        <NoSsr>
          {selectedRowData && selectedRowData !== null && Object.keys(selectedRowData).length > 0
              && (
                <MesheryResultDialog
                  rowData={selectedRowData}
                  close={self.resetSelectedRowData()}
                />
              )}
          <MUIDataTable className={classes.table} title="Performance Test Results" data={resultsForDisplay} columns={columns} options={options} />
          <MUIDataTable className={classes.table} title="SMI Test Results" data={smi_resultsForDisplay} columns={smi_columns} options={smi_options} />
        </NoSsr>
      );
    }
}
MesheryResults.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  // updateMeshResults: bindActionCreators(updateMeshResults, dispatch),
  updateResultsSelection: bindActionCreators(updateResultsSelection, dispatch),
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (state) => {
  const startKey = state.get('results').get('startKey');
  const results = state.get('results').get('results').toArray();
  const results_selection = state.get('results_selection').toObject();
  const user = state.get('user').toObject();
  if (typeof results !== 'undefined') {
    return {
      startKey, results, results_selection, user,
    };
  }
  return { results_selection, user };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withSnackbar(MesheryResults)));