import React from 'react';
import PropTypes from 'prop-types';
import Modal from './Modal';

/**
 * MesheryChartDialog provides a modal wrapper for displaying charts.
 */
function MesheryChartDialog(props) {
  const { open, title, handleClose, content } = props;
  return (
    <Modal
      open={open}
      closeModal={handleClose}
      title={title && title.length ? title : 'Comparison'}
      maxWidth="md"
    >
      {content}
    </Modal>
  );
}

MesheryChartDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  content: PropTypes.node.isRequired,
  title: PropTypes.string,
};

export default MesheryChartDialog;
