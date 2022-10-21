// @ts-check
import {
  Grid, Paper, TextField, Typography
} from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { connect } from "react-redux";
import { ctxUrl } from "../../utils/multi-ctx";
import { getPatternServiceName, createWorkloadTraitSets } from "./helpers";
import PatternServiceForm from "./LazyPatternServiceForm";


function MesheryMeshInterface({ adapter, selectedK8sContexts }) {
  const [workloadTraitsSet, setWorkloadTraitsSet] = useState([]);
  const [ns, setNS] = useState("default");


  async function submitPattern(pattern, del = false) {
    const res = await fetch(
      ctxUrl("/api/pattern/deploy", selectedK8sContexts), {
        headers : { "Content-Type" : "application/json", },
        method : del ? "DELETE" : "POST",
        body : JSON.stringify(pattern),
      });

    return res.text();
  }

  const handleSubmit = (cfg) => {
    submitPattern(cfg)
      .then((res) => console.log(res))
      .catch((err) => console.error(err));
  };

  const handleDelete = (cfg) => {
    submitPattern(cfg, true)
      .then((res) => console.log(res))
      .catch((err) => console.error(err));
  };

  useEffect(() => {
    createWorkloadTraitSets(adapter).then(res => setWorkloadTraitsSet(res));
  }, []);

  return (
    <Grid container spacing={1}>
      <Grid item xs={12}>
        <TextField
          label="Namespace"
          variant="filled"
          value={ns}
          onChange={(e) => setNS(e.target.value)}
          fullWidth
          required
        />
      </Grid>
      <Grid item md={8} xs={12}>
        <div>
          <Grid container spacing={1}>
            {workloadTraitsSet
              .filter((s) => s.type !== "addon")
              .sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
              .map((s) => (
                <Grid item xs={12} key={s}>
                  <PatternServiceForm schemaSet={s} onSubmit={handleSubmit} onDelete={handleDelete} namespace={ns} />
                </Grid>
              ))}
          </Grid>
        </div>
      </Grid>
      <Grid item md={4} xs={12}>
        <Paper style={{ padding : "1rem" }}>
          <Typography variant="h6" gutterBottom>
            Addons
          </Typography>
          <Grid container spacing={1}>
            {workloadTraitsSet
              .filter((s) => s.type === "addon")
              .sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
              .map((s) => (
                <Grid item key={s}>
                  <PatternServiceForm schemaSet={s} onSubmit={handleSubmit} onDelete={handleDelete} namespace={ns} />
                </Grid>
              ))}
          </Grid>
        </Paper>
      </Grid>
    </Grid>
  );
}


const mapStateToProps = (state) => {
  return { selectedK8sContexts : state.get("selectedK8sContexts") };
};

export default connect(mapStateToProps)(MesheryMeshInterface);
