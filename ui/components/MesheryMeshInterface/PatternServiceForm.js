// @ts-check
import React from "react";
import { Tab, Tabs, AppBar, Typography, Button } from "@material-ui/core";
import PatternServiceFormCore from "./PatternServiceFormCore";

function TabPanel(props) {
  const {
    children, value, index, ...other
  } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Typography>{children}</Typography>
      )}
    </div>
  );
}

function a11yProps(index) {
  return {
    id : `simple-tab-${index}`,
    "aria-controls" : `simple-tabpanel-${index}`,
  };
}

function RJSFButton({ handler, text, ...restParams }) {
  return (
    <Button variant="contained" color="primary" style={{ marginRight : "0.5rem" }} onClick={handler} {...restParams}>
      {text}
    </Button>
  );
}

function RJSFFormChildComponent({ onSubmit, onDelete }){
  return (
    <>
      <RJSFButton handler={onSubmit} text="Submit" />
      <RJSFButton handler={onDelete} text="Delete" />
    </>
  )
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
 *  onSettingsChange?: Function;
 *  onTraitsChange?: Function;
 *  formData?: Record<String, unknown>
 *  renderAsTooltip: boolean;
 *  reference?: Record<any, any>;
 *  appBarColor?: any;
 * }} props
 * @returns
 */
function PatternServiceForm({ formData, schemaSet, onSubmit, onDelete, reference, namespace, onSettingsChange, onTraitsChange }) {
  const [tab, setTab] = React.useState(0);

  const handleTabChange = (_, newValue) => {
    setTab(newValue);
  };
  const renderTraits = () => !!schemaSet.traits?.length;

  return (
    <PatternServiceFormCore
      formData={formData}
      schemaSet={schemaSet}
      onSubmit={onSubmit}
      onDelete={onDelete}
      reference={reference}
      namespace={namespace}
      onSettingsChange={onSettingsChange}
      onTraitsChange={onTraitsChange}
    >
      {(SettingsForm, TraitsForm) => {
        return (
          <div>
            <Typography variant="h6" gutterBottom>{schemaSet.workload?.title}</Typography>
            <div style={{ height : '57vh', overflowY : "scroll", overflowX : "hidden", scrollbarWidth : 'thin' }}>
              <AppBar position="static" >
                <Tabs value={tab} onChange={handleTabChange} aria-label="Pattern Service" >
                  <Tab label="Settings" {...a11yProps(0)} />
                  {
                    renderTraits()
                      ? <Tab label="Traits" {...a11yProps(1)} />
                      : null
                  }
                </Tabs>
              </AppBar>
              <TabPanel value={tab} index={0} style={{ marginTop : "1.1rem" }}>
                <SettingsForm RJSFFormChildComponent={RJSFFormChildComponent} />
              </TabPanel>
              <TabPanel value={tab} index={0} style={{ marginTop : "1.1rem" }}>
                <TraitsForm />
              </TabPanel>
            </div>
          </div>
        )
      }}
    </PatternServiceFormCore>
  )
}

export default PatternServiceForm;
