// @ts-check
import React from "react";
import { Tab, Tabs, AppBar, Typography, Box, Card } from "@material-ui/core";
import PatternService from "./PatternService";

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

/**
 * PatternServiceForm renders a form from the workloads schema and
 * traits schema
 * @param {{
 *  schemaSet: { workload: any, traits: any[], type: string };
 *  onChange: Function;
 *  onSubmit: Function;
 * }} props
 * @returns
 */
function PatternServiceForm({ schemaSet, onChange, onSubmit }) {
  const [tab, setTab] = React.useState(0);
  const [settings, setSettings] = React.useState({});
  const [traits, setTraits] = React.useState({});

  React.useEffect(() => {
    onChange?.({ [getPatternAttributeName(schemaSet.workload)]: { settings, traits } });
  }, [settings, traits]);

  const handleTabChange = (_, newValue) => {
    setTab(newValue);
  };

  const renderTraits = () => !!schemaSet.traits?.length;

  if (schemaSet.type === "addon") {
    return (
      <PatternService type="workload" jsonSchema={schemaSet.workload} onChange={setSettings} onSubmit={onSubmit} />
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
        <PatternService type="workload" jsonSchema={schemaSet.workload} onChange={setSettings} onSubmit={onSubmit} />
      </TabPanel>
      {renderTraits() ? (
        <TabPanel value={tab} index={1}>
          {schemaSet.traits?.map((trait) => (
            <PatternService
              type="trait"
              jsonSchema={trait}
              onSubmit={onSubmit}
              onChange={(val) => setTraits((t) => ({ ...t, [getPatternAttributeName(trait)]: val }))}
            />
          ))}
        </TabPanel>
      ) : null}
    </Card>
  );
}

export default PatternServiceForm;
