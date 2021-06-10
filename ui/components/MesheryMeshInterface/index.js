// @ts-check
import { Grid, Paper, Typography } from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { promisifiedDataFetch } from "../../lib/data-fetch";
import PatternServiceForm from "./PatternServiceForm";

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
    const res = await promisifiedDataFetch("/api/experimental/oam/workload");

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
    const res = await promisifiedDataFetch("/api/experimental/oam/trait");

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
 * }>>}
 */
async function createWorkloadTraitSets(adapter) {
  const workloads = await getWorkloadDefinitionsForAdapter(adapter);
  const traits = await getTraitDefinitionsForAdapter(adapter);

  const sets = [];
  workloads?.forEach((w) => {
    const item = { workload: w, traits: [] };

    item.traits = traits?.filter((t) => {
      if (Array.isArray(t?.oam_definition?.spec?.appliesToWorkloads))
        return t?.oam_definition?.spec?.appliesToWorkloads?.includes(w?.oam_definition?.metadata?.name);

      return false;
    });

    sets.push(item);
  });

  return sets;
}

/**
 * getJSONSchemaSets
 * @param {string} adapter
 * @returns
 */
async function getJSONSchemaSets(adapter) {
  const wtSets = await createWorkloadTraitSets(adapter);

  return wtSets?.map((s) => {
    const item = {
      workload: JSON.parse(s.workload?.oam_ref_schema),
      traits: s.traits?.map((t) => {
        const trait = JSON.parse(t?.oam_ref_schema);

        // Attaching internal metadata to the json schema
        trait._internal = {
          patternAttributeName: t?.oam_definition.metadata.name,
        };

        return trait;
      }),
      type: s.workload?.metadata?.["ui.meshery.io/category"],
    };

    // Attaching internal metadata to the json schema
    item.workload._internal = {
      patternAttributeName: s.workload?.oam_definition.metadata.name,
    };

    return item;
  });
}

function isObject(obj) {
  return obj === Object(obj) && Object.prototype.toString.call(obj) !== "[object Array]";
}

function recursiveCleanObject(obj) {
  if (isObject(obj)) {
    const objMap = Object.keys(obj);

    if (objMap.length === 0) return true;

    objMap.forEach((k) => {
      if (recursiveCleanObject(obj[k])) delete obj[k];
    });
  }

  return false;
}

/**
 * createPatternFromConfig will take in the form data
 * and will create a valid pattern from it
 *
 * It will/may also perform some sanitization on the
 * given inputs
 * @param {*} config
 */
function createPatternFromConfig(config) {
  const pattern = {};

  recursiveCleanObject(config);

  Object.keys(config).forEach((key) => {
    // Add it only if the settings are non empty or "true"
    if (config[key].settings) pattern[key] = config[key];
  });

  return pattern;
}

function MesheryMeshInterface({ adapter }) {
  const [schemeSets, setSchemaSets] = useState([]);
  const [config, setConfig] = useState({});

  useEffect(() => {
    console.log(config);
  }, [config]);

  const updateConfig = (val) => {
    setConfig((cfg) => {
      return createPatternFromConfig({ ...cfg, ...val });
    });
  };

  useEffect(() => {
    getJSONSchemaSets(adapter).then((res) => setSchemaSets(res));
  }, []);

  return (
    <Grid container spacing={1}>
      <Grid item md={6} xs={12}>
        <div>
          <Grid container spacing={1}>
            {schemeSets
              .filter((s) => s.type !== "addon")
              .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
              .map((s) => (
                <Grid item xs={12}>
                  <PatternServiceForm schemaSet={s} onChange={updateConfig} />
                </Grid>
              ))}
          </Grid>
        </div>
      </Grid>
      <Grid item md={6} xs={12}>
        <Paper style={{ padding: "1rem" }}>
          <Typography variant="h6" gutterBottom>
            Addons
          </Typography>
          <Grid container spacing={1}>
            {schemeSets
              .filter((s) => s.type === "addon")
              .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
              .map((s) => (
                <Grid item>
                  <PatternServiceForm schemaSet={s} onChange={updateConfig} />
                </Grid>
              ))}
          </Grid>
        </Paper>
      </Grid>
    </Grid>
  );
}

export default MesheryMeshInterface;
