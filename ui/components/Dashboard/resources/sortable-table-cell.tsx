import React from 'react';
import { Grid2, Tooltip, Typography, TableCell, TableSortLabel } from '@sistent/sistent';

type ColumnData = { name: string; label: React.ReactNode };
type ColumnMeta = { name: string; direction?: 'asc' | 'desc' };

type SortableTableCellProps = {
  index: number;
  columnData: ColumnData;
  columnMeta: ColumnMeta;
  onSort?: () => void;
  icon?: React.ReactNode;
  tooltip?: React.ReactNode;
};

export const SortableTableCell = ({
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
        <Grid2 style={{ display: 'flex', alignItems: 'center' }}>
          <Typography>
            <b>{columnData.label}</b>
          </Typography>
          {icon ? (
            <Tooltip title={tooltip} placement="top">
              <Typography style={{ display: 'flex', marginLeft: '2px' }} variant="span">
                {icon}
              </Typography>
            </Tooltip>
          ) : (
            ''
          )}
        </Grid2>
        <TableSortLabel
          active={columnMeta.name === columnData.name}
          direction={columnMeta.direction || 'asc'}
        ></TableSortLabel>
      </Grid2>
    </TableCell>
  );
};

type DefaultTableCellProps = {
  columnData: ColumnData;
  icon?: React.ReactNode;
  tooltip?: React.ReactNode;
};

export const DefaultTableCell = ({ columnData, icon, tooltip }: DefaultTableCellProps) => {
  return (
    <TableCell>
      <Grid2 style={{ display: 'flex' }}>
        <Grid2 style={{ display: 'flex', alignItems: 'center' }}>
          <Typography>
            <b>{columnData.label}</b>
          </Typography>
          {icon ? (
            <Tooltip title={tooltip ? tooltip : ''} placement="top">
              <Typography style={{ display: 'flex', marginLeft: '5px' }} variant="span">
                {icon}
              </Typography>
            </Tooltip>
          ) : (
            ''
          )}
        </Grid2>
      </Grid2>
    </TableCell>
  );
};
