// @ts-check
import { AppBar, Button, makeStyles, Tab, Tabs, Typography } from "@material-ui/core";
import React, { useEffect } from "react";
import { getMeshProperties } from "../../utils/nameMapper";
import PatternServiceFormCore from "./PatternServiceFormCore";
import SettingsIcon from '@material-ui/icons/Settings';

const useStyles = makeStyles(() => ({
  appBar : {
    // @ts-ignore
    boxShadow : `0px 2px 4px -1px "#677E88"`,
    // @ts-ignore
    background : "#677E88",
    position : "sticky",
    zIndex : 'auto',
  },
  tabPanel : {
    padding : "0px 2px"
  },
  formWrapper : {
    width : "100%"
  },
  settingsIcon : {
    color : "black"
  },
  appTabs : {
    width : 128,
    overflow : 'hidden',
    transition : 'width 0.5s',
    '&.Mui-disabled' : {
      width : 0,
    },
  },
  setIcon : {
    verticalAlign : 'middle',
    transform : "scale(0.8)"
  }
}));

function RJSFButton({ handler, text, ...restParams }) {
  return (
    <Button variant="contained" color="primary" style={{ margin : "0px 0.5rem 32px 0px" }} onClick={handler} {...restParams}>
      {text}
    </Button>
  );
}

function RJSFFormChildComponent({ onSubmit, onDelete }) {
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
 *  reference?: Record<any, any>;
 *  scroll?: Boolean; // If the window should be scrolled to zero after re-rendering
 * }} props
 * @returns
 */
function PatternServiceForm({ formData, schemaSet, onSubmit, onDelete, reference, namespace, onSettingsChange, onTraitsChange, scroll = false }) {
  const [tab, setTab] = React.useState(0);
  const classes = useStyles({ color : getMeshProperties(getMeshName(schemaSet))?.color });

  useEffect(() => {
    schemaSet.workload.properties.name = {
      description : "The Namespace for the service",
      default : "<Name of the Component>",
      type : "string"
    };
    schemaSet.workload.properties.namespace = {
      description : "The Name for the service",
      default : "default",
      type : "string",
    };
    schemaSet.workload.properties.label = {
      description : "The label for the service",
      additionalProperties : {
        "type" : "string"
      }
    };
    schemaSet.workload.properties.annotation = {
      description : "The annotation for the service",
      additionalProperties : {
        "type" : "string"
      }
    };
  }, [])

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
      scroll={scroll}
      tab={tab}
    >
      {(SettingsForm, TraitsForm) => {

        // For rendering addons without tabs
        if (schemaSet?.type === "addon") {
          return <SettingsForm />
        }

        // for rendering normal rjsf forms
        return (
          <div className={classes.formWrapper}>
            <AppBar className={classes.appBar}>
              <Tabs className={classes.appTabs} value={tab} onChange={handleTabChange} TabIndicatorProps={{
                style : {
                  display : "none",
                },
              }} aria-label="Pattern Service" >
                <Tab label={<div> <SettingsIcon className={classes.setIcon} />Settings</div>} {...a11yProps(0)} />
                {
                  renderTraits()
                    ? <Tab label="Traits" {...a11yProps(1)} />
                    : null
                }
              </Tabs>
            </AppBar>
            <TabPanel value={tab} index={0} className={classes.tabPanel}>
              <SettingsForm RJSFFormChildComponent={RJSFFormChildComponent} />
            </TabPanel>
            <TabPanel value={tab} index={1} className={classes.tabPanel}>
              <TraitsForm />
            </TabPanel>
          </div>
        )
      }}
    </PatternServiceFormCore>
  )
}

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

/**
 * @param {{ workload: { [x: string]: string; }; }} schema
 * @returns {String} name
 */
function getMeshName(schema) {
  return schema?.workload?.["service-mesh"]?.toLowerCase() || "core";
}

export default PatternServiceForm;
