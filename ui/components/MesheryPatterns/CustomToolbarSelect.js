import React from 'react';
import { IconButton, CustomTooltip } from '@layer5/sistent';
import DeleteIcon from '@mui/icons-material/Delete';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { clearResultsSelection } from '../../lib/store';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

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
    return (
      <div style={{ marginRight: '24px' }}>
        <CustomTooltip title={'Delete'}>
          <div>
            <IconButton
              onClick={this.handleClickDelete}
              disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
            >
              <DeleteIcon />
            </IconButton>
          </div>
        </CustomTooltip>
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
export default connect(mapStateToProps, mapDispatchToProps)(CustomToolbarSelect);
