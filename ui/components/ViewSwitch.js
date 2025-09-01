import { CustomTooltip, IconButton, useTheme } from '@sistent/sistent';
import GridOnIcon from '@mui/icons-material/GridOn';
import TableChartIcon from '@mui/icons-material/TableChart';

function ViewSwitch({ view, changeView }) {
  const theme = useTheme();
  return (
    <IconButton
      size="small"
      value={view}
      onClick={() => {
        changeView(view === 'grid' ? 'table' : 'grid');
      }}
      aria-label="Switch View"
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
    </IconButton>
  );
}

export default ViewSwitch;
