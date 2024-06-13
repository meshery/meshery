import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import { Grid, TableCell, TableSortLabel, Typography } from '@material-ui/core';

export const SortableTableCell = ({ index, columnData, columnMeta, onSort, icon, tooltip }) => {
  return (
    <TableCell key={index} onClick={onSort}>
      <Grid style={{ display: 'flex' }}>
        <Grid style={{ display: 'flex', alignItems: 'center' }}>
          <Typography>
            <b>{columnData.label}</b>
          </Typography>
          {icon ? (
            <CustomTextTooltip interactive={true} title={tooltip ? tooltip : ''} placement="top">
              <Typography style={{ display: 'flex', marginLeft: '5px' }} variant="span">
                {icon}
              </Typography>
            </CustomTextTooltip>
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
            <CustomTextTooltip interactive={true} title={tooltip ? tooltip : ''} placement="top">
              <Typography style={{ display: 'flex', marginLeft: '5px' }} variant="span">
                {icon}
              </Typography>
            </CustomTextTooltip>
          ) : (
            ''
          )}
        </Grid>
      </Grid>
    </TableCell>
  );
};
