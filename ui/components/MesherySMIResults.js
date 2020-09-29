import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, TableRow, TableCell, IconButton, Paper, Table, TableBody, TableHead, 
} from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MUIDataTable from 'mui-datatables';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';
import {
  updateProgress,
} from '../lib/store';
import dataFetch from '../lib/data-fetch';


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

class MesherySMIResults extends Component {
  constructor(props) {
    super(props);
    this.state = {
    
      smi_page: 0,
      smi_pageSize: 10,
      smi_search: '',
      smi_sortOrder: '',
      smi_results: [],
    };
  }

    componentDidMount = () => {
      const {
        smi_page, smi_pageSize, smi_search, smi_sortOrder,
      } = this.state;
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
      }, console.log('Could not fetch SMI results.'));
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
      const self = this;
      const { user } = this.props;
      const { smi_pageSize, smi_results } = this.state;

      const smi_resultsForDisplay = [];
      if(smi_results&&smi_results.results) {
        smi_results.results.map((val) => {
          smi_resultsForDisplay.push([val.id,val.date,val.mesh_name,val.mesh_version,val.passing_percentage,val.status])
        }) 
      }

      const smi_columns = ["ID","Date", "Service Mesh","Service Mesh Version", "% Passed","Status"];

      const smi_options = {
        sort: !(user && user.user_id === 'meshery'),
        search: !(user && user.user_id === 'meshery'),
        filterType: 'textField',
        expandableRows: true,
        selectableRows: false,
        rowsPerPage: smi_pageSize,
        rowsPerPageOptions: [10, 20, 25],
        fixedHeader: true,
        print: false,
        download: false,
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

      return (
        <NoSsr>
          <MUIDataTable title="SMI Test Results" data={smi_resultsForDisplay} columns={smi_columns} options={smi_options} />
        </NoSsr>
      );
    }
}
MesherySMIResults.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = (state) => {
  const user = state.get('user').toObject();
  return { user };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withSnackbar(MesherySMIResults)));