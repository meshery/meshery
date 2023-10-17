import React from 'react';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import { withStyles } from '@material-ui/core/styles';
import DeleteIcon from '@material-ui/icons/Delete';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { clearResultsSelection } from '../../lib/store';
const defaultToolbarSelectStyles = {
  iconButton: {},
  iconContainer: {
    marginRight: '24px',
  },
};

class CustomToolbarSelect extends React.Component {
  handleClickDelete = async () => {
    const toBeDeleted = this.props.selectedRows.data.map((idx) => ({
      id: this.props.patterns[idx.index]?.id,
      name: this.props.patterns[idx.index]?.name,
    }));
    let response = await this.props.showModal(
      toBeDeleted.length,
      toBeDeleted.map((p) => ' ' + p.name),
    );
    if (response.toLowerCase() == 'no') {
      return;
    }
    this.props.deletePatterns({ patterns: toBeDeleted }).then(() => {
      this.props.setSelectedRows([]);
    });
  };

  render() {
    const { classes } = this.props;
    return (
      <div className={classes.iconContainer}>
        <Tooltip title={'Delete'}>
          <IconButton className={classes.iconButton} onClick={this.handleClickDelete}>
            <DeleteIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      </div>
    );
  }
}
const mapDispatchToProps = (dispatch) => ({
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
});

const mapStateToProps = (state) => {
  const results_selection = state.get('results_selection').toObject();
  return { results_selection };
};
export default withStyles(defaultToolbarSelectStyles, { name: 'CustomToolbarSelect' })(
  connect(mapStateToProps, mapDispatchToProps)(CustomToolbarSelect),
);
