// @ts-check
import React, { useEffect } from "react";
import { Tab, Tabs, AppBar, Typography, Box, Card } from "@material-ui/core";
import PatternService from "./PatternService";
import useStateCB from "../../utils/hooks/useStateCB";
import PascalCaseToKebab from "../../utils/PascalCaseToKebab";

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box p={3}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

function getPatternAttributeName(jsonSchema) {
  return jsonSchema?._internal?.patternAttributeName || "NA";
}

function recursiveCleanObject(obj) {
  for (const k in obj) {
    if (!obj[k] || typeof obj[k] !== "object") continue;

    recursiveCleanObject(obj[k]);

    if (Object.keys(obj[k]).length === 0) delete obj[k];
  }
}

function recursiveCleanObjectExceptEmptyArray(obj) {

  for (const k in obj) {
    if (!obj[k] || typeof obj[k] !== "object" || Array.isArray(obj[k])) continue;

    recursiveCleanObjectExceptEmptyArray(obj[k]);

    if (Object.keys(obj[k]).length === 0) delete obj[k];
  }
}

/**
 * createPatternFromConfig will take in the form data
 * and will create a valid pattern from it
 *
 * It will/may also perform some sanitization on the
 * given inputs
 * @param {*} config
 */
function createPatternFromConfig(config, namespace, partialClean = false) {
  const pattern = {
    name: `pattern-${Math.random().toString(36).substr(2, 5)}`,
    services: {},
  };

  partialClean ? recursiveCleanObjectExceptEmptyArray(config) : recursiveCleanObject(config);

  Object.keys(config).forEach((key) => {
    // Add it only if the settings are non empty or "true"
    if (config[key].settings) {
      const name = PascalCaseToKebab(key);
      pattern.services[name] = config[key];

      pattern.services[name].type = key;
      pattern.services[name].namespace = namespace;
    }
  });

  Object.keys(pattern.services).forEach((key) => {
    // Delete the settings attribute/field if it is set to "true"
    if (pattern.services[key].settings === true) delete pattern.services[key].settings;
  });

  return pattern;
}

/**
 * PatternServiceForm renders a form from the workloads schema and
 * traits schema
 * @param {{
 *  schemaSet: { workload: any, traits: any[], type: string };
 *  onSubmit: Function;
 *  onDelete: Function;
 *  namespace: string;
 *  onChange?: Function
 * }} props
 * @returns
 */
function PatternServiceForm({ schemaSet, onChange, onSubmit, onDelete, namespace }) {
  const [tab, setTab] = React.useState(0);
  const [settings, setSettings, getSettingsRefValue] = useStateCB({});
  const [traits, setTraits, getTraitsRefValue] = useStateCB({});

  console.log(getTraitsRefValue());

  console.log({ settings, traits })

  const handleTabChange = (_, newValue) => {
    setTab(newValue);
  };

  const renderTraits = () => !!schemaSet.traits?.length;

  const submitHandler = (val) => {
    onSubmit?.(createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)]: val }, namespace))
  };

  const deleteHandler = (val) => {
    onDelete?.(createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)]: val }, namespace), true)
  };

  if (schemaSet.type === "addon") {
    return (
      <PatternService
        formData={settings}
        type="workload"
        jsonSchema={schemaSet.workload}
        onChange={setSettings}
        onSubmit={() => submitHandler({ settings: getSettingsRefValue() })}
        onDelete={() => deleteHandler({ settings: getSettingsRefValue() })}
      />
    );
  }

  return (
    <Card style={{ padding: "1rem" }}>
      <Typography variant="h6" gutterBottom>
        {schemaSet.workload.title}
      </Typography>
      <AppBar position="static">
        <Tabs value={tab} onChange={handleTabChange} aria-label="Pattern Service">
          <Tab label="Settings" {...a11yProps(0)} />
          {renderTraits() ? <Tab label="Traits" {...a11yProps(1)} /> : null}
        </Tabs>
      </AppBar>
      <TabPanel value={tab} index={0}>
        <PatternService
          type="workload"
          formData={settings}
          jsonSchema={schemaSet.workload}
          onChange={(val) => {
            console.log("rana", val)
            onChange(
              createPatternFromConfig(
                { [getPatternAttributeName(schemaSet.workload)]: { settings: val, traits } }, namespace, true), "");
            setSettings(val);
          }}
          onSubmit={() => submitHandler({ settings: getSettingsRefValue(), traits })}
          onDelete={() => deleteHandler({ settings: getSettingsRefValue(), traits })}
        />
      </TabPanel>
      {renderTraits() ? (
        <TabPanel value={tab} index={1}>
          {schemaSet.traits?.map((trait) => (
            <PatternService
              formData={traits[getPatternAttributeName(trait)]}
              type="trait"
              jsonSchema={trait}
              onChange={(val) => setTraits({ ...traits, [getPatternAttributeName(trait)]: val })}
            />
          ))}
        </TabPanel>
      ) : null}
    </Card>
  );
}

export default PatternServiceForm;
