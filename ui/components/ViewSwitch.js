import { ToggleButton } from "@material-ui/lab";
import GridOnIcon from "@material-ui/icons/GridOn";
import TableChartIcon from "@material-ui/icons/TableChart";

function ViewSwitch({ view, changeView }) {
  return (
    <ToggleButton
      size="small"
      value={view}
      onChange={() => {
        changeView(view === "grid" ? "table" : "grid");
      }}
      aria-label="Switch View"
    >
      {view === "grid" ? <GridOnIcon /> : <TableChartIcon />}
    </ToggleButton>
  );
}

export default ViewSwitch;
