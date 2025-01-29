import { ToggleButton, useTheme } from '@layer5/sistent';
import GridOnIcon from '@mui/icons-material/GridOn';
import TableChartIcon from '@mui/icons-material/TableChart';

function ViewSwitch({ view, changeView }) {
  const theme = useTheme();
  return (
    <ToggleButton
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
        <TableChartIcon style={{ fill: theme.palette.icon.default }} />
      ) : (
        <GridOnIcon style={{ fill: theme.palette.icon.default }} />
      )}
    </ToggleButton>
  );
}

export default ViewSwitch;
