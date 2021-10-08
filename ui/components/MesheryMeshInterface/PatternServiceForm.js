// @ts-check
import React from "react";
import { Tab, Tabs, AppBar, Typography, IconButton, Toolbar } from "@material-ui/core";
import { Delete } from "@material-ui/icons";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import Tooltip from '@material-ui/core/Tooltip';
import PatternService from "./PatternService";
import useStateCB from "../../utils/hooks/useStateCB";
import { pSBCr } from "../../utils/lightenOrDarkenColor"
import { CamelCaseToSentanceCase } from "../../utils/camelCaseToSentanceCase.js";
import { getPatternAttributeName, createPatternFromConfig } from "./helpers";

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
function PatternServiceForm({ formData, schemaSet, onSubmit, onDelete, reference, namespace, renderAsTooltip, appBarColor, onSettingsChange, onTraitsChange }) {
  const [tab, setTab] = React.useState(0);
  const [settings, setSettings, getSettingsRefValue] = useStateCB(formData?.settings ? formData.settings : {}, onSettingsChange);
  const [traits, setTraits, getTraitsRefValue] = useStateCB(formData?.traits ? formData.traits : {}, onTraitsChange);
  const handleTabChange = (_, newValue) => setTab(newValue);

  const renderTraits = () => !!schemaSet.traits?.length;

  const submitHandler = (val) => {
    onSubmit?.(createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)] : val }, namespace))
  };

  const deleteHandler = (val) => {
    onDelete?.(createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)] : val }, namespace), true)
  };

  if (reference){
    if (reference.current == null) reference.current = {}

    reference.current.submit = (cb) => {
      submitHandler(cb(getSettingsRefValue(), getTraitsRefValue()))
    }
    reference.current.getSettings = () => getSettingsRefValue()
    reference.current.getTraits = () => getTraitsRefValue()
  }

  if (schemaSet.type === "addon") {
    return (
      <PatternService
        formData={settings}
        type="workload"
        jsonSchema={schemaSet.workload}
        onChange={setSettings}
        onSubmit={() => submitHandler({ settings : getSettingsRefValue() })}
        onDelete={() => deleteHandler({ settings : getSettingsRefValue() })}
        renderAsTooltip={renderAsTooltip}
      />
    );
  }

  return (
    <div>
      {!renderAsTooltip ? (<Typography variant="h6" gutterBottom>
        {schemaSet.workload?.title}
      </Typography>) : (
        <AppBar style={{ boxShadow : `0px 2px 4px -1px ${pSBCr(appBarColor, -30)}`, position: "sticky" }}>
          <Toolbar variant="dense" style={{ padding : "0 0px", background : `linear-gradient(115deg, ${pSBCr( appBarColor, -20)} 0%, ${appBarColor} 100%)`, height : "0.7rem !important" }}>
            <p style={{ margin : "auto auto auto 10px", fontSize : "16px" }}>{schemaSet.workload.title || CamelCaseToSentanceCase(schemaSet.workload["object-type"])}</p>
            {schemaSet?.workload?.description && (
              <label htmlFor="help-button" >
                <Tooltip title={schemaSet?.workload?.description} >
                  <IconButton component="span" style={{ paddingRight : "0.1rem" }} >
                    <HelpOutlineIcon style={{ color : '#fff' }} fontSize="small" />
                  </IconButton>
                </Tooltip>
              </label>
            )}
            <IconButton style={{ paddingLeft : '0.1rem' }} onClick={() => deleteHandler({ settings : getSettingsRefValue(), traits : getTraitsRefValue() })}>
              <Delete style={{ color : "#ffffff" }} fontSize="small"/>
            </IconButton>
          </Toolbar>
        </AppBar>
      )}
      <div style={{ maxHeight : '300px', scrollbarWidth : 'thin' }}>
        {!renderAsTooltip && (<AppBar position="static" >
          <Tabs value={tab} onChange={handleTabChange} aria-label="Pattern Service" >
            <Tab label="Settings" {...a11yProps(0)} />
            {renderTraits()
              ? <Tab label="Traits" {...a11yProps(1)} />
              : null}
          </Tabs>
        </AppBar>)}
        {renderAsTooltip && renderTraits() && (
          <AppBar style={{ background : 'inherit', boxShadow : 'none', color : 'black' }} position="static">
            <Tabs TabIndicatorProps={{ style : { background : '#00b39f' } }} style={{ margin : 0 }} value={tab} onChange={handleTabChange} aria-label="Pattern Service">
              <Tab label="Settings" style={{ minWidth : "50%", margin : 0 }} {...a11yProps(0)} />
              <Tab label="Traits" style={{ minWidth : "50%", margin : 0 }} {...a11yProps(1)} />
            </Tabs>
          </AppBar>)}
        <TabPanel value={tab} index={0} style={{ marginTop : "1.1rem" }}>
          <PatternService
            type="workload"
            formData={settings}
            jsonSchema={schemaSet.workload}
            onChange={setSettings}
            onSubmit={() => submitHandler({ settings : getSettingsRefValue(), traits })}
            onDelete={() => deleteHandler({ settings : getSettingsRefValue(), traits })}
            renderAsTooltip={renderAsTooltip}
          />
        </TabPanel>
        {renderTraits() ? (
          <TabPanel value={tab} index={1} style={{ marginTop : "1.1rem" }}>
            {schemaSet.traits?.map((trait) => (
              <PatternService
                formData={traits[getPatternAttributeName(trait)]}
                type="trait"
                jsonSchema={trait}
                onChange={(val) => setTraits({ ...traits, [getPatternAttributeName(trait)] : val })}
                renderAsTooltip={renderAsTooltip}
              />
            ))}
          </TabPanel>
        ) : null}
      </div>
    </div>
  );
}

export default PatternServiceForm;
