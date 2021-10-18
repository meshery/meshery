import React, { useState, useEffect } from "react";
import { makeStyles } from "@material-ui/core/styles";
import {
  IconButton,
  Tooltip,
  Grid,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  CardContent,
  Card,
  CardActions,
  AppBar,
  Toolbar,
} from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";
import DeleteIcon from "@material-ui/icons/Delete";
import SaveIcon from "@material-ui/icons/Save";
import { promisifiedDataFetch } from "../../lib/data-fetch";
import { CircularProgress } from "@material-ui/core";
import PatternServiceForm from "../MesheryMeshInterface/PatternServiceForm";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import { Button } from "@material-ui/core";
import jsYaml from "js-yaml";
import ArrowBackIcon from "@material-ui/icons/ArrowBack";
import PascalCaseToKebab from "../../utils/PascalCaseToKebab";

const useStyles = makeStyles((theme) => ({
  codeMirror : {
    "& .CodeMirror" : {
      minHeight : "300px",
      height : "60vh",
    },
  },
  backButton : {
    marginRight : theme.spacing(2),
  },
  appBar : {
    marginBottom : "16px",
  },
  yamlDialogTitle : {
    display : "flex",
    alignItems : "center",
  },
  yamlDialogTitleText : {
    flexGrow : 1,
  },
  fullScreenCodeMirror : {
    height : "100%",
    "& .CodeMirror" : {
      minHeight : "300px",
      height : "100%",
    },
  },
}));



export function MesheryConfigurationForm({ application, pattern, onSubmit, show }) {
  let formType ;
  let formHeading ;
  if (application){
    formType = application;
    formHeading = "Application"
  } else {
    formType = pattern;
    formHeading = "Pattern"
  }
  const [schemaSet, setSchemaSet] = useState();
  const [deployServiceConfig, setDeployServiceConfig] = useState(getPatternJson() || {});
  const [yaml, setYaml] = useState("");
  const [expanded, setExpanded] = useState([]);
  // const [changedYaml, setChangedYaml] = useState("");
  const classes = useStyles();



  function getPatternJson() {
    let patternString;
    if (application){
      patternString = application.application_file;
    } else {
      patternString = pattern.pattern_file
    }
    return jsYaml.load(patternString).services;
  }

  async function fetchWorkloadAndTraitsSchema() {
    try {
      const workloads = await promisifiedDataFetch("/api/oam/workload");
      const traits = await promisifiedDataFetch("/api/oam/trait");

      console.log({ workloads, traits });

      const workloadTraitSets = createWorkloadTraitSets(workloads, traits);

      return workloadTraitSets;
    } catch (e) {
      console.log("Error in Fetching Workload or traits", e);
      return {};
    }
  }

  function createWorkloadTraitSets(workloads, traits) {
    const sets = [];
    workloads?.forEach((w) => {
      const item = { workload : w, traits : [] };

      item.traits = traits?.filter((t) => {
        if (Array.isArray(t?.oam_definition?.spec?.appliesToWorkloads))
          return t?.oam_definition?.spec?.appliesToWorkloads?.includes(w?.oam_definition?.metadata?.name);

        return false;
      });

      sets.push(item);
    });

    return sets;
  }

  async function getJSONSchemaSets() {
    const wtSets = await fetchWorkloadAndTraitsSchema();

    return wtSets?.map((s) => {
      const item = {
        workload : JSON.parse(s.workload?.oam_ref_schema),
        traits : s.traits?.map((t) => {
          const trait = JSON.parse(t?.oam_ref_schema);

          // Attaching internal metadata to the json schema
          trait._internal = {
            patternAttributeName : t?.oam_definition.metadata.name,
          };

          return trait;
        }),
        type : s.workload?.metadata?.["ui.meshery.io/category"],
      };

      // Attaching internal metadata to the json schema
      item.workload._internal = {
        patternAttributeName : s.workload?.oam_definition.metadata.name,
      };

      return item;
    });
  }

  function getPatternAttributeName(jsonSchema) {
    if (application){
      return PascalCaseToKebab(jsonSchema?._internal?.patternAttributeName || "NA");
    }
    return jsonSchema?._internal?.patternAttributeName || "NA";
  }

  function getPatternKey(cfg) {
    return Object.keys(cfg?.services)?.[0] || undefined;
  }

  const handleSubmit = (cfg, patternName) => {
    console.log("submitted", { cfg, patternName });
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    if (key) setDeployServiceConfig({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    handleExpansion(patternName);
  };

  const handleChangeData = (cfg, patternName) => {
    console.log("Ran Changed", { cfg, patternName });
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    if (key) setDeployServiceConfig({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
  };

  const handleDelete = (cfg, patternName) => {
    console.log("deleted", cfg);
    const newCfg = schemaSet.filter((schema) => schema.workload.title !== patternName);
    setSchemaSet(newCfg);
  };

  const handleDeploy = (cfg) => {
    const deployConfig = {};
    deployConfig.name = formType.name;
    deployConfig.services = cfg;
    const deployConfigYaml = jsYaml.dump(deployConfig);
    setYaml(deployConfigYaml);
  };

  const handleExpansion = (item) => {
    let expandedItems = [...expanded];
    if (expandedItems.includes(item)) {
      expandedItems = expandedItems.filter((el) => el !== item);
    } else {
      expandedItems.push(item);
    }
    setExpanded(expandedItems);
  };

  function handleSubmitFinalPattern(yaml, id, name, action) {
    onSubmit(yaml, id, name, action);
    show(false);
  }

  const ns = "default";

  function saveCodeEditorChanges(data) {
    setYaml(data.valueOf().getValue());
  }

  function insertPattern(workload) {
    const attrName = getPatternAttributeName(workload);
    var returnValue = {};
    Object.keys(deployServiceConfig).find((key) => {
      if (deployServiceConfig[key]["type"] === attrName) {
        returnValue = deployServiceConfig[key];
        return true;
      }
    });

    return returnValue;
  }

  useEffect(() => {
    getJSONSchemaSets().then((res) => setSchemaSet(res));
  }, []);

  if (!schemaSet) {
    return <CircularProgress />;
  }

  return (
    <>
      <AppBar position="static" className={classes.appBar} elevation={0}>
        <Toolbar>
          <IconButton edge="start" className={classes.backButton} color="inherit" onClick={() => show(false)}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            Edit {`${formHeading}`} Configuration of <i>{`${formType.name}`}</i>
          </Typography>
        </Toolbar>
      </AppBar>
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          {schemaSet
            .filter((s) => s.type !== "addon")
            .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
            .map((s) => accordion(s))}
          <Accordion
            expanded={expanded.includes("addon")}
            onChange={() => handleExpansion("addon")}
            style={{ width : "100%" }}
          >
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              <Typography variant="h6">Configure Addons</Typography>
            </AccordionSummary>
            <AccordionDetails>
              {schemaSet
                .filter((s) => s.type === "addon")
                .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
                .map((s) => (
                  <Grid item>
                    <PatternServiceForm
                      formData={deployServiceConfig[s.workload?.title]}
                      onChange={handleChangeData}
                      schemaSet={s}
                      onSubmit={handleSubmit}
                      onDelete={handleDelete}
                      namespace={ns}
                    />
                  </Grid>
                ))}
            </AccordionDetails>
          </Accordion>
        </Grid>
        <Grid item xs={12} md={6}>
          <CodeEditor />
        </Grid>
      </Grid>
    </>
  );

  function CustomButton({ title, onClick }) {
    return (
      <Button
        fullWidth
        color="primary"
        variant="contained"
        onClick={onClick}
        style={{
          marginTop : "16px",
          padding : "10px",
        }}
      >
        {title}
      </Button>
    );
  }

  function accordion(schema) {
    const patternName = schema?.workload?.title;

    return (
      <Accordion
        expanded={expanded.includes(patternName)}
        onChange={() => handleExpansion(patternName)}
        style={{ width : "100%" }}
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant="h6">{patternName || "Expand More"}</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <PatternServiceForm
            formData={ application?deployServiceConfig[getPatternAttributeName(schema.workload)]:insertPattern(schema.workload)}
            onChange={handleChangeData}
            schemaSet={schema}
            onSubmit={(val) => handleSubmit(val, patternName)}
            onDelete={(val) => handleDelete(val, patternName)}
            namespace={ns}
          />
        </AccordionDetails>
      </Accordion>
    );
  }

  function CodeEditor() {
    const cardStyle = { marginBottom : "16px", position : "sticky", float : "right", minWidth : "100%" };
    const cardcontentStyle = { margin : "16px" };

    const classes = useStyles();

    return (
      <div>
        <Card style={cardStyle}>
          <CardContent style={cardcontentStyle}>
            <CodeMirror
              value={yaml}
              className={classes.codeMirror}
              options={{
                theme : "material",
                lineNumbers : true,
                lineWrapping : true,
                gutters : ["CodeMirror-lint-markers"],
                lint : true,
                mode : "text/x-yaml",
              }}
              onBlur={(a) => saveCodeEditorChanges(a)}
            />
            <CustomButton
              title= {`Save ${formHeading}`}
              onClick={() => handleSubmitFinalPattern(yaml, "", `meshery_${Math.floor(Math.random() * 100)}`, "upload")}
            />
            <CardActions style={{ justifyContent : "flex-end" }}>
              <Tooltip title= {`Update ${formHeading}`}>
                <IconButton
                  aria-label="Update"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, formType.id, formType.name, "update")}
                >
                  <SaveIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title= {`Delete ${formHeading}`}>
                <IconButton
                  aria-label="Delete"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, formType.id, formType.name, "delete")}
                >
                  <DeleteIcon />
                </IconButton>
              </Tooltip>
            </CardActions>
          </CardContent>
        </Card>
      </div>
    );
  }
}
