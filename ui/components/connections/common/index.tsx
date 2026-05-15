import React, { memo } from 'react';
import { CustomTextTooltip } from '@/components/meshery-mesh-interface/PatternService/CustomTextTooltip';
import { Grid2, Typography, TableCell, TableSortLabel } from '@sistent/sistent';

const HeaderLabel = ({ label, icon, tooltip }) => (
  <Grid2 style={{ display: 'flex', alignItems: 'center' }}>
    <Typography>
      <b>{label}</b>
    </Typography>
    {icon && (
      <CustomTextTooltip interactive={true} title={tooltip || ''} placement="top">
        <Typography style={{ display: 'flex', marginLeft: '5px' }} variant="span">
          {icon}
        </Typography>
      </CustomTextTooltip>
    )}
  </Grid2>
);

const SortableTableCell_ = ({ index, columnData, columnMeta, onSort, icon, tooltip }) => {
  return (
    <TableCell key={index} onClick={onSort}>
      <Grid2 style={{ display: 'flex' }}>
        <HeaderLabel label={columnData.label} icon={icon} tooltip={tooltip} />
        <TableSortLabel
          active={columnMeta.name === columnData.name}
          direction={columnMeta.direction || 'asc'}
        />
      </Grid2>
    </TableCell>
  );
};

export const SortableTableCell = memo(SortableTableCell_);

const DefaultTableCell_ = ({ columnData, icon, tooltip }) => {
  return (
    <TableCell>
      <Grid2 style={{ display: 'flex' }}>
        <HeaderLabel label={columnData.label} icon={icon} tooltip={tooltip} />
      </Grid2>
    </TableCell>
  );
};

export const DefaultTableCell = memo(DefaultTableCell_);
