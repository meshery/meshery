import React from 'react';
import { CustomTooltip, DeleteIcon, IconButton } from '@sistent/sistent';

import { Keys } from '@meshery/schemas/permissions';

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
            permissionKey={Keys.CatalogManagementDeleteADesign}
          >
            <DeleteIcon />
          </IconButton>
        </div>
      </CustomTooltip>
    </div>
  );
};

export default CustomToolbarSelect;
