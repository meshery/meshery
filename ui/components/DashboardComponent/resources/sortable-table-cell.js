import { UsesSistent } from '@/components/SistentWrapper';
import { Grid, Tooltip, Typography } from '@layer5/sistent';
import { TableCell, TableSortLabel } from '@mui/material';

export const SortableTableCell = ({ index, columnData, columnMeta, onSort, icon, tooltip }) => {
  return (
    <UsesSistent>
      <TableCell key={index} onClick={onSort}>
        <Grid style={{ display: 'flex' }}>
          <Grid style={{ display: 'flex', alignItems: 'center' }}>
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
          </Grid>
          <TableSortLabel
            active={columnMeta.name === columnData.name}
            direction={columnMeta.direction || 'asc'}
          ></TableSortLabel>
        </Grid>
      </TableCell>
    </UsesSistent>
  );
};

export const DefaultTableCell = ({ columnData, icon, tooltip }) => {
  return (
    <TableCell>
      <Grid style={{ display: 'flex' }}>
        <Grid style={{ display: 'flex', alignItems: 'center' }}>
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
        </Grid>
      </Grid>
    </TableCell>
  );
};
