//@ts-check
import { Grid } from "@material-ui/core";
import React from "react";
import PerformanceCard from "./PerformanceCard";

/**
 * PerformanceProfileGrid is the react component for rendering
 * grid
 * @param {{
 *  profiles: Array<{
 *    id: string,
 *    created_at: string,
 *    updated_at: string,
 *    endpoints: Array<string>,
 *    load_generators: Array<string>,
 *    name: string,
 *    user_id: string,
 *    duration: string,
 *  }>
 * }} props props
 */
function PerformanceProfileGrid({ profiles = [] }) {
  return (
    <Grid container spacing={2} style={{ padding: "1rem" }} >
      {profiles.map((profile) => (
        <Grid item md={6} xs={12}>
          <PerformanceCard
            name={profile.name}
            handleEdit={() => console.log("edit on", profile.id)}
            handleDelete={() => console.log("delete on", profile.id)}
            handleRunTest={() => console.log("run test on", profile.id)}
          />
        </Grid>
      ))}
    </Grid>
  );
}

export default PerformanceProfileGrid;
