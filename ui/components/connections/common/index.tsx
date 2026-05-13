import React, { memo } from 'react';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { Grid2, Typography, TableCell, TableSortLabel } from '@sistent/sistent';

type ColumnData = {
  label: React.ReactNode;
  name?: string;
};

type ColumnMeta = {
  name?: string;
  direction?: 'asc' | 'desc';
};

type HeaderLabelProps = {
  label: React.ReactNode;
  icon?: React.ReactNode;
  tooltip?: React.ReactNode;
};

type SortableTableCellProps = {
  index: number;
  columnData: ColumnData;
  columnMeta: ColumnMeta;
  onSort: () => void;
  icon?: React.ReactNode;
  tooltip?: React.ReactNode;
};

type DefaultTableCellProps = {
  columnData: ColumnData;
  icon?: React.ReactNode;
  tooltip?: React.ReactNode;
};

const HeaderLabel = ({ label, icon, tooltip }: HeaderLabelProps) => (
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

const SortableTableCell_ = ({
  index,
  columnData,
  columnMeta,
  onSort,
  icon,
  tooltip,
}: SortableTableCellProps) => {
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

const DefaultTableCell_ = ({ columnData, icon, tooltip }: DefaultTableCellProps) => {
  return (
    <TableCell>
      <Grid2 style={{ display: 'flex' }}>
        <HeaderLabel label={columnData.label} icon={icon} tooltip={tooltip} />
      </Grid2>
    </TableCell>
  );
};

export const DefaultTableCell = memo(DefaultTableCell_);
