import { ToggleButton, ToggleButtonGroup, } from "@mui/material";
import GridOnIcon from '@mui/icons-material/GridOn';
import TableChartIcon from "@mui/icons-material/TableChart";

function ViewSwitch({ view, changeView }) {

  return (
    <ToggleButtonGroup
      size="small"
      value={view}
      exclusive
      onChange={(_, newView) => changeView(newView)}
      aria-label="Switch View"
    >
      <ToggleButton value="grid" data-cy="grid-view">
        <GridOnIcon />
      </ToggleButton>
      <ToggleButton value="table" data-cy="table-view">
        <TableChartIcon />
      </ToggleButton>
    </ToggleButtonGroup>    
  )
}

export default ViewSwitch