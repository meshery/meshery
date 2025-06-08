import { Grid2, Tooltip, Typography, TableCell, TableSortLabel } from '@sistent/sistent';

export const SortableTableCell = ({ index, columnData, columnMeta, onSort, icon, tooltip }) => {
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

export const DefaultTableCell = ({ columnData, icon, tooltip }) => {
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
