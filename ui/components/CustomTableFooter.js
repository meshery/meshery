//NOTE: Component not used anywhere
import React from 'react';
import PropTypes from 'prop-types';
import { withStyles, useTheme } from '@material-ui/core/styles';
import TableFooter from '@material-ui/core/TableFooter';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import IconButton from '@material-ui/core/IconButton';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import { iconSmall } from '../css/icons.styles';

const actionsStyles = (theme) => ({
  wrapper: {
    flexShrink: 0,
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(2.5),
  },
});

const TablePaginationActions = ({ classes, onChangePage, count, page, rowsPerPage }) => {
  const theme = useTheme();

  const handleBackButtonClick = () => {
    onChangePage(page - 1);
  };

  const handleNextButtonClick = () => {
    onChangePage(page + 1);
  };

  return (
    <div className={classes.wrapper}>
      <IconButton onClick={handleBackButtonClick} disabled={page === 0} aria-label="Previous Page">
        {theme.direction === 'rtl' ? (
          <KeyboardArrowRight style={iconSmall} />
        ) : (
          <KeyboardArrowLeft style={iconSmall} />
        )}
      </IconButton>
      <IconButton
        onClick={handleNextButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage) - 1}
        aria-label="Next Page"
      >
        {theme.direction === 'rltr' ? (
          <KeyboardArrowLeft style={iconSmall} />
        ) : (
          <KeyboardArrowRight style={iconSmall} />
        )}
      </IconButton>
    </div>
  );
};

TablePaginationActions.propTypes = {
  classes: PropTypes.object.isRequired,
  onChangePage: PropTypes.func.isRequired,
  page: PropTypes.number.isRequired,
  count: PropTypes.number.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
};

const TablePaginationActionsWrapper = withStyles(actionsStyles, {
  withTheme: true,
})(TablePaginationActions);

const defaultFooterStyles = {};

const CustomTableFooter = ({ changePage, rowsPerPage, page, count }) => {
  const customLabelDisplayedRows = () => `Page ${page + 1}`;

  return (
    <TableFooter>
      <TableRow>
        <TablePagination
          labelRowsPerPage=""
          labelDisplayedRows={customLabelDisplayedRows}
          rowsPerPageOptions={[10]}
          colSpan={3}
          count={count}
          rowsPerPage={rowsPerPage}
          page={page}
          SelectProps={{ native: true }}
          onChangePage={changePage}
          ActionsComponent={TablePaginationActionsWrapper}
        />
      </TableRow>
    </TableFooter>
  );
};

CustomTableFooter.propTypes = {
  changePage: PropTypes.func.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
  page: PropTypes.number.isRequired,
  count: PropTypes.number.isRequired,
};

export default withStyles(defaultFooterStyles, { name: 'CustomFooter' })(CustomTableFooter);
