import React from 'react';
import { IconButton, CustomTooltip } from '@layer5/sistent';
import DeleteIcon from '@mui/icons-material/Delete';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { clearResultsSelection } from '../../lib/store';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const CustomToolbarSelect = ({
  selectedRows,
  patterns,
  showModal,
  deletePatterns,
  setSelectedRows,
}) => {
  const handleClickDelete = async () => {
    const toBeDeleted = selectedRows.data.map((idx) => ({
      id: patterns[idx.index]?.id,
      name: patterns[idx.index]?.name,
    }));
    let response = await showModal(
      toBeDeleted.length,
      toBeDeleted.map((p) => ' ' + p.name),
    );
    if (response.toLowerCase() === 'no') {
      return;
    }
    deletePatterns({ patterns: toBeDeleted }).then(() => {
      setSelectedRows([]);
    });
  };

  return (
    <div style={{ marginRight: '24px' }}>
      <CustomTooltip title={'Delete'}>
        <div>
          <IconButton
            onClick={handleClickDelete}
            disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
          >
            <DeleteIcon />
          </IconButton>
        </div>
      </CustomTooltip>
    </div>
  );
};

const mapDispatchToProps = (dispatch) => ({
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
});

const mapStateToProps = (state) => {
  const results_selection = state.get('results_selection').toObject();
  return { results_selection };
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomToolbarSelect);
