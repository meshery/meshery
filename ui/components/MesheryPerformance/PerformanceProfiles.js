//@ts-check
import { ToggleButton, ToggleButtonGroup } from "@material-ui/lab";
import React, { useState } from "react";
import GridOnIcon from "@material-ui/icons/GridOn";
import TableChartIcon from "@material-ui/icons/TableChart";
import PerformanceProfileTable from "./PerformanceProfileTable";

/**
 * Type Definition for View Type
 * @typedef {"grid" | "table"} TypeView
 */

/**
 * ViewSwitch component renders a switch for toggling between
 * grid and table views
 * @param {{ view: TypeView, changeView: (view: TypeView) => void }} props
 */
function ViewSwitch({ view, changeView }) {
  return (
    <ToggleButtonGroup value={view} exclusive onChange={(_, newView) => changeView(newView)} aria-label="Switch View">
      <ToggleButton value="grid">
        <GridOnIcon />
      </ToggleButton>
      <ToggleButton value="table">
        <TableChartIcon />
      </ToggleButton>
    </ToggleButtonGroup>
  );
}

function PerformanceProfile() {
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );

  return (
    <div>
      <div style={{ margin: "0 0 0 auto" }}><ViewSwitch view={viewType} changeView={setViewType}/></div>
      {viewType === "grid" ? "Hello" : <PerformanceProfileTable />}
    </div>
  );
}

export default PerformanceProfile;
