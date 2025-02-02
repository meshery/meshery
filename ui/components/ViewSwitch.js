import { CustomTooltip, ToggleButton, useTheme } from '@layer5/sistent';
import GridOnIcon from '@mui/icons-material/GridOn';
import TableChartIcon from '@mui/icons-material/TableChart';

function ViewSwitch({ view, changeView }) {
  const theme = useTheme();
  return (
    <ToggleButton
      style={{ border: 'none' }}
      size="small"
      value={view}
      onChange={() => {
        changeView(view === 'grid' ? 'table' : 'grid');
      }}
      aria-label="Switch View"
      sx={{
        border: 'none',
      }}
    >
      {view === 'grid' ? (
        <CustomTooltip title="Grid View">
          <div>
            <TableChartIcon style={{ color: theme.palette.icon.default }} />
          </div>
        </CustomTooltip>
      ) : (
        <CustomTooltip title="Table View">
          <div>
            <GridOnIcon style={{ color: theme.palette.icon.default }} />
          </div>
        </CustomTooltip>
      )}
    </ToggleButton>
  );
}

export default ViewSwitch;
