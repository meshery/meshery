//@ts-check
import React from "react";
import Card from "@material-ui/core/Card";
import { Button, Grid, IconButton, Typography } from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import EditIcon from "@material-ui/icons/Edit";

function PerformanceCard({ name, handleDelete, handleEdit, handleRunTest }) {
  return (
    <Card style={{ padding: "1rem" }} raised>
      <Grid style={{ marginBottom: "3rem" }} container spacing={1} alignContent="space-between" alignItems="center">
        <Grid item xs={8}>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
        </Grid>
        <Grid item xs={4}>
          <div style={{ width: "fit-content", margin: "0 0 0 auto" }}>
            <IconButton onClick={() => handleDelete()}>
              <DeleteIcon />
            </IconButton>
            <IconButton onClick={() => handleEdit()}>
              <EditIcon />
            </IconButton>
          </div>
        </Grid>
      </Grid>
      <div>
        <Button
          color="primary"
          variant="contained"
          onClick={() => handleRunTest()}
          style={{ display: "block", marginLeft: "auto" }}
        >
          Run Test
        </Button>
      </div>
    </Card>
  );
}

export default PerformanceCard;
