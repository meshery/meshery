import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  TableRow, TableCell, IconButton, Table, TableBody, TableHead, Tooltip
} from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import MUIDataTable from 'mui-datatables';
import Moment from 'react-moment';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';
import { updateProgress, } from '../lib/store';
import dataFetch from '../lib/data-fetch';


const styles = (theme) => ({
  grid : { padding : theme.spacing(2), },
  secondaryTable : { borderRadius : 10, backgroundColor : theme.palette.secondary.elevatedComponents2, },
  tableHeader : {
    fontWeight : 'bolder',
    fontSize : 18,
  },
});

class MesherySMIResults extends Component {
  constructor(props) {
    super(props);
    this.state = {

      smi_page : 0,
      smi_pageSize : 10,
      smi_search : '',
      smi_sortOrder : '',
      smi_results : [],
      count : 0,
    };
  }

  componentDidMount = () => {
    const {
      smi_page, smi_pageSize, smi_search, smi_sortOrder,
    } = this.state;
    this.fetchSMIResults(smi_page, smi_pageSize, smi_search, smi_sortOrder);
  }

  fetchSMIResults = (page, pageSize, search, sortOrder) => {
    const self = this;
    let query = '';
    if (typeof search === 'undefined' || search === null) {
      search = '';
    }
    if (typeof sortOrder === 'undefined' || sortOrder === null) {
      sortOrder = '';
    }

    query = `?page=${page}&pageSize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(sortOrder)}`;
    dataFetch(`/api/smi/results${query}`, {
      method : 'GET',
      credentials : 'include',
    }, (result) => {
      if (typeof result !== 'undefined' && result.results) {
        self.setState({ smi_results : result });
        self.setState({ count : result?.total_count })
        self.setState({ smi_pageSize : result?.page_size })
      }
    }, () => self.handleError('Could not fetch SMI results.'));
  }

  handleError = (error) => {
    this.props.updateProgress({ showProgress : false });
    // console.log(`error fetching results: ${error}`);
    const self = this;
    this.props.enqueueSnackbar(`There was an error fetching results: ${error}`, {
      variant : 'error',
      action : (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => self.props.closeSnackbar(key)}
        >
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 8000,
    });
  }

  resetSelectedRowData() {
    const self = this;
    return () => {
      self.setState({ selectedRowData : null });
    };
  }

  render() {
    const self = this;
    const { classes, user } = this.props;
    const { smi_pageSize, smi_results, count } = this.state;

    const smi_resultsForDisplay = [];
    if (smi_results && smi_results?.results) {
      smi_results?.results?.map((val) => {
        smi_resultsForDisplay.push([val.id, val.date, val.mesh_name, val.mesh_version, val.passing_percentage, val.status])
      })
    }
    const smi_columns = [
      {
        name : 'ID',
        label : 'ID',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
          customBodyRender : (value) => (
            <Tooltip title={value} placement="top">
              <div>{value.slice(0, 5) + "..."}</div>
            </Tooltip>
          )
        },
      },
      {
        name : 'Date',
        label : 'Date',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
          customBodyRender : (value) => (
            <Moment format="LLLL">{value}</Moment>
          ),
        },
      },
      {
        name : 'Service Mesh',
        label : 'Service Mesh',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
        },
      },
      {
        name : 'Service Mesh Version',
        label : 'Service Mesh Version',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
        },
      },
      {
        name : '% Passed',
        label : '% Passed',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
        },
      },
      {
        name : 'status',
        label : 'Status',
        options : {
          sort : true,
          searchable : true,
          customHeadRender : ({ index, ...column }) => {
            return (
              <TableCell key={index}>
                <b>{column.label}</b>
              </TableCell>

            )
          },
        },
      }

    ]

    const smi_options = {
      sort : !(user && user.user_id === 'meshery'),
      search : !(user && user.user_id === 'meshery'),
      filter : false,
      expandableRows : true,
      selectableRows : false,
      serverSide : true,
      count,
      rowsPerPage : smi_pageSize,
      rowsPerPageOptions : [10, 20, 25],
      fixedHeader : true,
      print : false,
      download : false,
      renderExpandableRow : (rowData, rowMeta) => {
        const column = ["Specification", "Assertions", "Time", "Version", "Capability", "Result", "Reason"]
        const data = smi_results?.results[rowMeta.dataIndex]?.more_details?.map((val) => {
          return [val.smi_specification, val.assertions, val.time, val.smi_version, val.capability, val.status, val.reason]
        })
        const colSpan = rowData.length + 1
        return (
          <TableRow>
            <TableCell colSpan={colSpan}>
              <div className={classes.secondaryTable}>
                <Table aria-label="a dense table">
                  <TableHead>
                    <TableRow>
                      {column.map((val) => (<TableCell colSpan={colSpan} key={val} >{val}</TableCell>))}
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {data?.map((row) => (
                      <TableRow key={row?.uniqueID}>
                        {row?.map(val => {
                          if (val && val.match(/[0-9]+m[0-9]+.+[0-9]+s/i) != null) {
                            const time = val.split(/m|s/)
                            return <TableCell colSpan={colSpan} key={val}>{time[0] + "m " + parseFloat(time[1]).toFixed(1) + "s"}</TableCell>
                          } else {
                            return <TableCell colSpan={colSpan} key={val}>{val}</TableCell>
                          }
                        }
                        )}
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </TableCell>
          </TableRow>
        );
      },
      onTableChange : (action, tableState) => {
        const sortInfo = tableState.announceText
          ? tableState.announceText.split(' : ')
          : [];
        let order = '';
        if (tableState.activeColumn) {
          order = `${smi_columns[tableState.activeColumn].name} desc`;
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
                self.fetchSMIResults(self.state.smi_page, self.state.smi_pageSize, tableState.searchText !== null
                  ? tableState.searchText
                  : '', self.state.smi_sortOrder);
              }
            }, 500);
            break;
          case 'sort':
            if (sortInfo.length === 2) {
              if (sortInfo[1] === 'ascending') {
                order = `${smi_columns[tableState.activeColumn].name} asc`;
              } else {
                order = `${smi_columns[tableState.activeColumn].name} desc`;
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
      <>
        <MUIDataTable
          title={<div className={classes.tableHeader}>Service Mesh Interface Conformance Results</div>}
          data={smi_resultsForDisplay}
          columns={smi_columns}
          options={smi_options}
        />
      </>
    );
  }
}
MesherySMIResults.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });
const mapStateToProps = (state) => {
  const user = state.get('user')?.toObject();
  return { user };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(withSnackbar(MesherySMIResults)));