// @ts-check
import {
  Grid, Paper, TextField, Typography
} from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { promisifiedDataFetch } from "../../lib/data-fetch";
import { getPatternServiceName, getPatternServiceType } from "./helpers";
import PatternServiceForm from "./LazyPatternServiceForm";

// ********************************** TYPE DEFINTIONS **********************************

/**
 * @typedef {Object} OAMDefinition
 * @property {string} kind
 * @property {string} apiVersion
 * @property {Record<string, any>} metadata
 * @property {Record<string, any>} spec
 */

/**
 * @typedef {string} OAMRefSchema
 */

/**
 * @typedef {Object} OAMGenericResponse
 * @property {OAMDefinition} oam_definition
 * @property {OAMRefSchema} oam_ref_schema
 * @property {string} host
 * @property {Record<string, any>} metadata
 */

// ******************************************************************************************

/**
 * getWorkloadDefinitionsForAdapter will fetch workloads for the given
 * adapter from meshery server
 * @param {string} adapter
 * @returns {Promise<Array<OAMGenericResponse>>}
 */
async function getWorkloadDefinitionsForAdapter(adapter) {
  try {
    const res = await promisifiedDataFetch("/api/oam/workload?trim=true");

    return res?.filter((el) => el?.metadata?.["adapter.meshery.io/name"] === adapter);
  } catch (error) {
    console.error(error);
    return [];
  }
}

/**
 * getTraitDefinitionsForAdapter will fetch tratis for the given
 * adapter from meshery server
 * @param {string} adapter
 * @returns {Promise<Array<OAMGenericResponse>>}
 */
async function getTraitDefinitionsForAdapter(adapter) {
  try {
    const res = await promisifiedDataFetch("/api/oam/trait?trim=true");

    return res?.filter((el) => el?.metadata?.["adapter.meshery.io/name"] === adapter);
  } catch (error) {
    console.error(error);
    return [];
  }
}

/**
 * createWorkloadTraitSets returns an array of workloads and traits object
 * which are interrelated
 * @param {string} adapter
 * @returns {Promise<Array<{
 *  workload: OAMGenericResponse;
 *  traits: Array<OAMGenericResponse>;
 *  type?: string;
 * }>>}
 */
async function createWorkloadTraitSets(adapter) {
  const workloads = await getWorkloadDefinitionsForAdapter(adapter);
  const traits = await getTraitDefinitionsForAdapter(adapter);

  const sets = [];
  workloads?.forEach((w) => {
    const item = { workload : w, traits : [], type : getPatternServiceType(w?.metadata) };

    item.traits = traits?.filter((t) => {
      if (Array.isArray(t?.oam_definition?.spec?.appliesToWorkloads))
        return t?.oam_definition?.spec?.appliesToWorkloads?.includes(w?.oam_definition?.metadata?.name);

      return false;
    });

    sets.push(item);
  });

  return sets;
}

async function submitPattern(pattern, del = false) {
  const res = await fetch(
    "/api/pattern/deploy", {
      headers : { "Content-Type" : "application/json", },
      method : del ? "DELETE" : "POST",
      body : JSON.stringify(pattern),
    });

  return res.text();
}


function MesheryMeshInterface({ adapter }) {
  const [workloadTraitsSet, setWorkloadTraitsSet] = useState([]);
  const [ns, setNS] = useState("default");

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
    createWorkloadTraitSets(adapter).then(res => setWorkloadTraitsSet(res))
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
                <Grid item xs={12}>
                  <PatternServiceForm schemaSet={s} onSubmit={handleSubmit} onDelete={handleDelete} namespace={ns} renderAsTooltip={false} />
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
                <Grid item>
                  <PatternServiceForm schemaSet={s} onSubmit={handleSubmit} onDelete={handleDelete} namespace={ns} renderAsTooltip={false} />
                </Grid>
              ))}
          </Grid>
        </Paper>
      </Grid>
    </Grid>
  );
}

export default MesheryMeshInterface;
