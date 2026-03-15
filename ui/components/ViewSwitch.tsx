import React from 'react';
import { CustomTooltip, IconButton, useTheme } from '@sistent/sistent';
import GridOnIcon from '@mui/icons-material/GridOn';
import TableChartIcon from '@mui/icons-material/TableChart';

type ViewSwitchProps = {
  view: 'grid' | 'table';
  changeView: (_view: 'grid' | 'table') => void;
};

function ViewSwitch({ view, changeView }: ViewSwitchProps) {
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
