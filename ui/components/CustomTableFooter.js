import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableFooter from '@material-ui/core/TableFooter';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import FirstPageIcon from '@material-ui/icons/FirstPage';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import LastPageIcon from '@material-ui/icons/LastPage';

const actionsStyles = theme => ({
  root: {
    flexShrink: 0,
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(2.5),
  },
});

class TablePaginationActions extends React.Component {
  handleFirstPageButtonClick = event => {
    this.props.onChangePage(0);
  };

  handleBackButtonClick = event => {
    this.props.onChangePage(this.props.page - 1);
  };

  handleNextButtonClick = event => {
    this.props.onChangePage(this.props.page + 1);
  };

  handleLastPageButtonClick = event => {
    this.props.onChangePage(
      // event,
      Math.max(0, Math.ceil(this.props.count / this.props.rowsPerPage) - 1),
    );
  };

  render() {
    const { classes, count, page, rowsPerPage, theme } = this.props;

    return (
      <div className={classes.root}>
        {/* <IconButton
          onClick={this.handleFirstPageButtonClick}
          disabled={page === 0}
          aria-label="First Page"
        >
          {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
        </IconButton> */}
        <IconButton
          onClick={this.handleBackButtonClick}
          disabled={page === 0}
          aria-label="Previous Page"
        >
          {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
        </IconButton>
        <IconButton
          onClick={this.handleNextButtonClick}
          disabled={page >= Math.ceil(count / rowsPerPage) - 1}
          aria-label="Next Page"
        >
          {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
        </IconButton>
        {/* <IconButton
          onClick={this.handleLastPageButtonClick}
          disabled={page >= Math.ceil(count / rowsPerPage) - 1}
          aria-label="Last Page"
        >
          {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
        </IconButton> */}
      </div>
    );
  }
}

TablePaginationActions.propTypes = {
  classes: PropTypes.object.isRequired,
  //   count: PropTypes.number.isRequired,
  onChangePage: PropTypes.func.isRequired,
  page: PropTypes.number.isRequired,
//   rowsPerPage: PropTypes.number.isRequired,
//   theme: PropTypes.object.isRequired,
};

const TablePaginationActionsWrapper = withStyles(actionsStyles, { withTheme: true })(
  TablePaginationActions,
);

const defaultFooterStyles = {
};

class CustomTableFooter extends Component {

  // state = {
  //     page: 0,
  //     rowsPerPage: 10,
  // }

  // handleChangePage = (event, page) => {
  //     this.setState({ page });
  // };

    // handleChangeRowsPerPage = event => {
    //     this.setState({ page: 0, rowsPerPage: event.target.value });
    // };
    customLabelDisplayedRows = ({ from, to, count }) => {
      return `Page ${this.props.page + 1}`;
    }

    render() {
      // const { rowsPerPage, page } = this.state;
      return (
        <TableFooter>
          <TableRow>
            <TablePagination 
              labelRowsPerPage={''}
              labelDisplayedRows={this.customLabelDisplayedRows}
              rowsPerPageOptions={[10]}
              colSpan={3}
              count={this.props.count}
              rowsPerPage={this.props.rowsPerPage}
              page={this.props.page}
              SelectProps={{
                native: true,
              }}
              onChangePage={this.props.changePage}
              // onChangeRowsPerPage={this.handleChangeRowsPerPage}
              ActionsComponent={TablePaginationActionsWrapper}
            />
          </TableRow>
        </TableFooter>
      );
    }
}

CustomTableFooter.propTypes = {
  // classes: PropTypes.object.isRequired,
  changePage: PropTypes.func.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
  page: PropTypes.number.isRequired,
};

export default withStyles(defaultFooterStyles, { name: "CustomFooter" })(CustomTableFooter);