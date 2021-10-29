
import {
  Accordion, AccordionDetails, AccordionSummary, AppBar, ButtonGroup, CircularProgress, Divider, FormControl, Grid, IconButton, makeStyles, MenuItem, Paper, Select, TextField, Toolbar, Tooltip, Typography
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import FileCopyIcon from '@material-ui/icons/FileCopy';
import ListAltIcon from '@material-ui/icons/ListAlt';
import SaveIcon from '@material-ui/icons/Save';
import { Autocomplete } from '@material-ui/lab';
import jsYaml from "js-yaml";
import React, { useContext, useEffect, useRef, useState } from "react";
import { trueRandom } from "../../lib/trueRandom";
import { SchemaContext } from "../../utils/context/schemaSet";
import { nameToIcon, getMeshProperties } from "../../utils/nameMapper";
import { groupWorkloadByVersion } from "../../utils/workloadFilter";
import { createPatternFromConfig, getPatternServiceName } from "../MesheryMeshInterface/helpers";
import LazyPatternServiceForm, { getWorkloadTraitAndType } from "../MesheryMeshInterface/LazyPatternServiceForm";
import PatternServiceForm from "../MesheryMeshInterface/PatternServiceForm";
import CodeEditor from "./CodeEditor";

const useStyles = makeStyles((theme) => ({
  backButton: {
    marginRight: theme.spacing(2),
  },
  appBar: {
    marginBottom: "16px",
    backgroundColor: "#fff",
    borderRadius: "8px"
  },
  yamlDialogTitle: {
    display: "flex",
    alignItems: "center"
  },
  yamlDialogTitleText: {
    flexGrow: 1
  },
  fullScreenCodeMirror: {
    height: '100%',
    '& .CodeMirror': {
      minHeight: "300px",
      height: '100%',
    }
  },
  formCtrl: {
    width: "90px",
    minWidth: "90px",
    maxWidth: "90px",
    marginRight: 8,
  },
  autoComplete: {
    width: "120px",
    minWidth: "120px",
    maxWidth: 150,
    marginRight: "auto"
  },
  btngroup: {
    marginLeft: "auto",
    overflowX: "auto",
    overflowY: "hidden"
  },
  paper: {
    backgroundColor: "#fcfcfc",
    padding: 12,
    height: "100%",
  },
}))


function PatternForm({ pattern, onSubmit, show: setSelectedPattern }) {
  const { workloadTraitSet, meshWorkloads } = useContext(SchemaContext);
  const [workloadTraitsSet, setWorkloadTraitsSet] = useState(workloadTraitSet);
  const [deployServiceConfig, setDeployServiceConfig] = useState(getPatternJson() || {});
  const [yaml, setYaml] = useState(pattern.pattern_file);
  const [selectedMeshType, setSelectedMeshType] = useState("core")
  const [selectedVersionMesh, setSelectedVersionMesh] = useState()
  const [selectedVersion, setSelectedVersion] = useState("")
  const [activeForm, setActiveForm] = useState()
  const classes = useStyles();
  const reference = useRef({});

  useEffect(() => {
    if (workloadTraitSet != workloadTraitsSet) {
      setWorkloadTraitsSet(workloadTraitSet)
    }
  }, [workloadTraitSet]);

  useEffect(() => {
    // core is not versioned
    if (selectedMeshType == "core") {
      console.log({ meshWorkloads })
    } else {
      const meshVersionsWithDetails = groupWlByVersion()
      setSelectedVersionMesh(meshVersionsWithDetails)
    }
  }, [selectedMeshType])

  useEffect(() => {
    if (selectedVersionMesh) {
      setSelectedVersion(Object.keys(selectedVersionMesh).sort().reverse()[0])
    }
  }, [selectedVersionMesh])


  function groupWlByVersion() {
    const mfw = meshWorkloads[selectedMeshType];
    return mfw ? groupWorkloadByVersion(mfw) : {};
  }


  function getPatternJson() {
    const patternString = pattern.pattern_file;
    // @ts-ignore
    return jsYaml.load(patternString).services;
  }

  function getPatternKey(cfg) {
    return Object.keys(cfg?.services)?.[0] || undefined;
  }

  const handleSubmit = (cfg, patternName) => {
    console.log("submitted", { cfg, patternName })
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [key]: cfg?.services?.[key] });
    if (key) setDeployServiceConfig({ ...deployServiceConfig, [key]: cfg?.services?.[key] });
  }

  const handleSettingsChange = (schemaSet) => () => {
    const config = createPatternFromConfig({
      [getPatternServiceName(schemaSet)]: {
        // @ts-ignore
        settings: reference.current?.getSettings(),
        // @ts-ignore
        traits: reference.current?.getTraits()
      }
    }, "default", true);

    handleChangeData(config, "");
  }

  const handleChangeData = (cfg, patternName) => {
    console.log("Ran Changed", { cfg, patternName })
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [getPatternKey(cfg)]: cfg?.services?.[key] });
    if (key)
      setDeployServiceConfig({ ...deployServiceConfig, [getPatternKey(cfg)]: cfg?.services?.[key] });
  }

  const handleDelete = (cfg, patternName) => {
    console.log("deleted", cfg);
    const newCfg = workloadTraitsSet?.filter(schema => schema.workload.title !== patternName)
    setWorkloadTraitsSet(newCfg);
  }

  const handleDeploy = (cfg) => {
    const deployConfig = {};
    deployConfig.name = pattern.name;
    deployConfig.services = cfg;
    const deployConfigYaml = jsYaml.dump(deployConfig);
    setYaml(deployConfigYaml);
  }

  function handleSubmitFinalPattern(yaml, id, name, action) {
    onSubmit(yaml, id, name, action);
    setSelectedPattern(resetSelectedPattern()); // Remove selected pattern
  }

  const ns = "default";

  function saveCodeEditorChanges(data) {
    setYaml(data.valueOf().getValue())
  }

  function insertPattern(workload) {
    const attrName = getPatternServiceName(workload);
    var returnValue = {}
    Object.keys(deployServiceConfig).find(key => {
      if (deployServiceConfig[key]['type'] === attrName) {
        returnValue = deployServiceConfig[key]
        return true
      }
    })

    return returnValue;
  }

  function getMeshOptions() {
    return meshWorkloads ? Object.keys(meshWorkloads) : []
  }

  function handleMeshSelection(event) {
    setSelectedMeshType(event.target.value);
  }

  function handleVersionChange(_, value) {
    setSelectedVersion(value)
  }

  async function getPatternProps(schema) {
    const refinedSchema = await getWorkloadTraitAndType(schema)
    setActiveForm(refinedSchema)
  }

  console.log({ selectedVersionMesh })
  console.log({ selectedMeshType })
  console.log({ selectedVersionMesh })

  if (!workloadTraitsSet) return <CircularProgress />

  return (
    <>
      <AppBar position="static" className={classes.appBar} elevation={0}>
        <Toolbar className={classes.toolbar}>
          <FormControl className={classes.formCtrl}>
            <Select
              labelId="service-mesh-selector"
              id="service-mesh-selector"
              value={selectedMeshType}
              onChange={handleMeshSelection}
              disableUnderline
            >
              {getMeshOptions().map(item => {
                const details = getMeshProperties(item)
                return (<MenuItem value={details.name}>
                  <li>
                    <img src={details.img} height="32px" />
                  </li>
                </MenuItem>)
              })}
            </Select>
          </FormControl>
          {
            selectedVersion &&
            <Autocomplete
              options={Object.keys(selectedVersionMesh).sort().reverse()}
              renderInput={(params) => <TextField {...params} variant="outlined" label="Version" />}
              value={selectedVersion}
              onChange={handleVersionChange}
              className={classes.autoComplete}
              disableClearable
            />
          }
          <ButtonGroup
            disableFocusRipple
            disableElevation
            className={classes.btngroup}
          >
            {selectedVersionMesh && selectedVersionMesh?.[selectedVersion]
              ?.sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
              .map((s) => {
                const name = s?.workload?.oam_definition?.spec?.metadata?.k8sKind
                return nameToIcon(name, () => getPatternProps(s))
              })
            }
            <Divider
              orientation="vertical"
            />
          </ButtonGroup>
          <Tooltip title="Save Pattern as New File">
            <IconButton
              aria-label="Save"
              color="primary"
              onClick={() => handleSubmitFinalPattern(yaml, "", `meshery_${Math.floor(trueRandom() * 100)}`, "upload")}
            >
              <FileCopyIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Update Pattern">
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "update")}
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Pattern">
            <IconButton
              aria-label="Delete"
              color="secondary"
              onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "delete")}
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="List View">
            <IconButton color="primary" onClick={() => setActiveForm(null)}>
              <ListAltIcon />
            </IconButton>
          </Tooltip>
          {/* <Typography variant="h6">
            Edit Pattern Configuration of <i>{`${pattern.name}`}</i>
          </Typography> */}
        </Toolbar>
      </AppBar>
      <Grid container spacing={3}>
        {
          activeForm
            ? (
              <Grid item xs={12} md={6}>
                <Paper className={classes.paper} >
                  <PatternServiceForm
                    schemaSet={activeForm}
                    jsonSchema={activeForm.workload}
                    formData={insertPattern(activeForm.workload)}
                    onSettingsChange={handleSettingsChange(activeForm.workload)}
                    onSubmit={(val) => handleSubmit(val, pattern.name)}
                    onDelete={(val) => handleDelete(val, pattern.name)}
                    namespace={ns}
                    reference={reference}
                  />
                </Paper>
              </Grid>
            ) : (
              <Grid item xs={12} md={6}>
                {
                  selectedMeshType === "core"
                  && meshWorkloads[selectedMeshType]
                    ?.filter((s) => s.type !== "addon")
                    .sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
                    .map((s, i) => (
                      <div style={{ marginBottom: "0.5rem" }} key={`svc-form-${i}`} >
                        <LazyPatternServiceForm
                          schemaSet={s}
                          formData={insertPattern(s.workload)}
                          onSettingsChange={handleSettingsChange(s.workload)}
                          onSubmit={(val) => handleSubmit(val, pattern.name)}
                          onDelete={(val) => handleDelete(val, pattern.name)}
                          namespace={ns}
                          reference={reference}
                          renderAsTooltip
                          appBarColor={getMeshProperties("core").color}
                        />
                      </div>))
                }
                {selectedVersionMesh && selectedVersionMesh?.[selectedVersion]
                  ?.filter((s) => s.type !== "addon")
                  .sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
                  .map((s, i) => (
                    <div style={{ marginBottom: "0.5rem" }} key={`svc-form-${i}`} >
                      <LazyPatternServiceForm
                        schemaSet={s}
                        formData={insertPattern(s.workload)}
                        onSettingsChange={handleSettingsChange(s.workload)}
                        onSubmit={(val) => handleSubmit(val, pattern.name)}
                        onDelete={(val) => handleDelete(val, pattern.name)}
                        namespace={ns}
                        reference={reference}
                        renderAsTooltip
                        appBarColor={getMeshProperties(selectedMeshType).color}
                      />
                    </div>))}
                <Accordion elevation={0} style={{ width: '100%' }}>
                  <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                    <Typography variant="h6">
                      Configure Addons
                    </Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    {selectedVersionMesh && selectedVersionMesh?.[selectedVersion]
                      ?.filter((s) => s.type === "addon")
                      .sort((a, b) => (getPatternServiceName(a.workload) < getPatternServiceName(b.workload) ? -1 : 1))
                      .map((s, i) => (
                        <Grid item key={`svc-form-addons-${i}`}>
                          <LazyPatternServiceForm
                            formData={deployServiceConfig[s.workload?.title]}
                            onSettingsChange={handleSettingsChange(s.workload)}
                            schemaSet={s}
                            onSubmit={handleSubmit}
                            onDelete={handleDelete}
                            namespace={ns}
                            reference={reference}
                          />
                        </Grid>
                      ))}
                  </AccordionDetails>
                </Accordion>
              </Grid>)
        }
        <Grid item xs={12} md={6} >
          {/* <CodeEditor yaml={yaml} pattern={pattern} handleSubmitFinalPattern={handleSubmitFinalPattern} saveCodeEditorChanges={saveCodeEditorChanges} /> */}
        </Grid>
      </Grid>
    </>
  );
}

export default PatternForm