import { ToggleButton, ToggleButtonGroup } from "@material-ui/lab";
import GridOnIcon from "@material-ui/icons/GridOn";
import TableChartIcon from "@material-ui/icons/TableChart";
import { Button } from "@material-ui/core";
import VisibilityIcon from '@material-ui/icons/Visibility';
import VisibilityOffIcon from '@material-ui/icons/VisibilityOff';

function ViewSwitch({ view, changeView, catalogVisibility, handleCatalogVisibility, hideCatalog }) {
  return (
    <ToggleButtonGroup
      size="small"
      value={view}
      exclusive
      onChange={(_, newView) => {
        if (newView !==null){
          changeView(newView)
        }
      }
      }
      aria-label="Switch View"
    >
      { !hideCatalog &&  // In application we  don't have catalog, hence this check
        <Button onClick = {handleCatalogVisibility}>
          {catalogVisibility ? <VisibilityIcon /> : <VisibilityOffIcon />}
        </Button>
      }
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
