//import useState from 'react';
import React, { useState, useEffect } from "react";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from "next/router";
import { withStyles } from "@material-ui/core/styles";
import { FormControl, FormLabel, FormGroup, FormControlLabel, Switch } from "@material-ui/core";
import NoSsr from "@material-ui/core/NoSsr";
import dataFetch from "../lib/data-fetch";
import { updateUser, updateProgress, toggleCatalogContent } from "../lib/store";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import { Paper, Tooltip } from "@material-ui/core";
import SettingsRemoteIcon from "@material-ui/icons/SettingsRemote";
import SettingsCellIcon from "@material-ui/icons/SettingsCell";
import ExtensionSandbox from "./ExtensionSandbox";
import RemoteComponent from "./RemoteComponent";
import ExtensionPointSchemaValidator from "../utils/ExtensionPointSchemaValidator";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTachometerAlt } from "@fortawesome/free-solid-svg-icons";
import MesherySettingsPerformanceComponent from "./MesherySettingsPerformanceComponent";
import { ctxUrl } from "../utils/multi-ctx";
import { iconMedium } from "../css/icons.styles";
import { getTheme, setTheme } from "../utils/theme";
import { isExtensionOpen } from "../pages/_app";
import { withNotify } from "../utils/hooks/useNotification";
import { EVENT_TYPES } from "../lib/event-types";

const styles = (theme) => ({
  statsWrapper : {
    // padding : theme.spacing(2),
    maxWidth : "100%",
    height : "auto",
    borderTopLeftRadius : 0,
    borderTopRightRadius : 0,
    borderBottomLeftRadius : 3,
    borderBottomRightRadius : 3,
  },
  paperRoot : {
    flexGrow : 1,
    maxWidth : "100%",
    marginLeft : 0,
    borderTopLeftRadius : 3,
    borderTopRightRadius : 3,
  },
  tabs : {
    marginLeft : 0,
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === "dark" ? "#00B39F" : theme.palette.primary,
    },
  },
  tab : {
    maxWidth : "min(33%, 200px)",
    minWidth : "50px",
    margin : 0,
    "&.Mui-selected" : {
      color : theme.palette.type === "dark" ? "#00B39F" : theme.palette.primary,
    },
  },
  icon : {
    display : "inline",
    verticalAlign : "text-top",
    width : theme.spacing(1.75),
    marginLeft : theme.spacing(0.5),
  },
  iconText : {
    display : "inline",
    verticalAlign : "middle",
  },
  backToPlay : { margin : theme.spacing(2) },
  link : { cursor : "pointer" },
  formContainer : {
    display : "flex",
    "flex-wrap" : "wrap",
    "justify-content" : "space-evenly",
    padding : 50,
  },
  formGrp : {
    padding : 20,
    border : "1.5px solid #969696",
    display : "flex",
    width : "70%",
  },
  formLegend : { fontSize : 20 },
  formLegendSmall : { fontSize : 16 },
  switchBase : {
    color : "#647881",
    "&$checked" : { color : "#00b39f" },
    "&$checked + $track" : { backgroundColor : "rgba(0,179,159,0.5)" },
  },
  track : { backgroundColor : "rgba(100,120,129,0.5)" },
  checked : {},
  tabLabel : {
    [theme.breakpoints.up("sm")] : {
      fontSize : "1em",
    },
    [theme.breakpoints.between("xs", "sm")] : {
      fontSize : "0.8em",
    },
  },
});

function ThemeToggler({ theme, themeSetter, notify, classes }) {
  const [themeToggle, setthemeToggle] = useState(false);
  const defaultTheme = "light";
  const handle = () => {
    if (isExtensionOpen()) {
      return;
    }

    theme === "dark" ? setthemeToggle(true) : setthemeToggle(false);
    setTheme(theme);
  };

  useEffect(() => {
    if (isExtensionOpen()) {
      if (getTheme() && getTheme() !== defaultTheme) {
        themeSetter(defaultTheme);
      }
      return;
    }

    themeSetter(getTheme() || defaultTheme);
  }, []);

  useEffect(handle, [theme]);

  const themeToggler = () => {
    if (isExtensionOpen()) {
      notify({ message : "Toggling between themes is not supported in MeshMap", event_type : EVENT_TYPES.INFO });
      return;
    }
    theme === "light" ? themeSetter("dark") : themeSetter("light");
  };

  return (
    <div onClick={themeToggler}>
      <Switch
        color="primary"
        classes={{
          switchBase : classes.switchBase,
          track : classes.track,
          checked : classes.checked,
          font : classes.checked,
        }}
        checked={themeToggle}
        onChange={themeToggler}
      />{" "}
      Dark Mode
    </div>
  );
}

function UserPreference(props) {
  const [anonymousStats, setAnonymousStats] = useState(props.anonymousStats);
  const [perfResultStats, setPerfResultStats] = useState(props.perfResultStats);
  const [tabVal, setTabVal] = useState(0);
  const [userPrefs, setUserPrefs] = useState(ExtensionPointSchemaValidator("user_prefs")());
  const [providerType, setProviderType] = useState("");
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [capabilitiesLoaded, setCapabilitiesLoaded] = useState(false);

  const handleCatalogContentToggle = () => {
    props.toggleCatalogContent({ catalogVisibility : !catalogContent });

    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  };

  const handleCatalogPreference = (catalogContent) => {
    let body = Object.assign({}, extensionPreferences);
    body["catalogContent"] = catalogContent;
    dataFetch(
      "/api/user/prefs",
      { credentials : "include", method : "POST", body : JSON.stringify({ usersExtensionPreferences : body }) },
      () => {
        const notify = props.notify;
        notify({
          message : `Catalog Content was ${catalogContent ? "enab" : "disab"}led`,
          event_type : EVENT_TYPES.SUCCESS,
        });
      },
      handleError("There was an error sending your preference")
    );
  };

  const handleToggle = (name) => () => {
    if (name === "anonymousUsageStats") {
      setAnonymousStats(!anonymousStats);
      handleChange(name, !anonymousStats);
    } else {
      setPerfResultStats(!perfResultStats);
      handleChange(name, !perfResultStats);
    }
  };

  const handleError = (name) => () => {
    props.updateProgress({ showProgress : false });
    const notify = props.notify;
    notify({ message : name, event_type : EVENT_TYPES.ERROR });
  };

  const handleChange = (name, resultState) => {
    let val = resultState,
        msg;
    if (name === "anonymousUsageStats") {
      msg = val ? "Sending anonymous usage statistics was enabled" : "Sending anonymous usage statistics was disabled";
    } else {
      msg = val
        ? "Sending anonymous performance results was enabled"
        : "Sending anonymous performance results was disabled";
    }

    const requestBody = JSON.stringify({
      anonymousUsageStats : name === "anonymousUsageStats" ? val : anonymousStats,
      anonymousPerfResults : name === "anonymousPerfResults" ? val : perfResultStats,
    });

    // console.log(requestBody,anonymousStats,perfResultStats);

    props.updateProgress({ showProgress : true });
    dataFetch(
      ctxUrl("/api/user/prefs", props.selectedK8sContexts),
      {
        credentials : "include",
        method : "POST",
        headers : { "Content-Type" : "application/json;charset=UTF-8" },
        body : requestBody,
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          console.log(result);
          const notify = props.notify;
          notify({ message : msg, event_type : val ? EVENT_TYPES.SUCCESS : EVENT_TYPES.INFO });
        }
      },
      handleError("There was an error sending your preference")
    );
  };

  const handleTabValChange = (event, newVal) => {
    setTabVal(newVal);
  };

  useEffect(() => {
    if (props.capabilitiesRegistry && !capabilitiesLoaded) {
      setCapabilitiesLoaded(true); // to prevent re-compute
      setUserPrefs(ExtensionPointSchemaValidator("user_prefs")(props.capabilitiesRegistry?.extensions?.user_prefs));
      setProviderType(props.capabilitiesRegistry?.provider_type);
    }
  }, [props.capabilitiesRegistry]);

  useEffect(() => {
    dataFetch(
      "/api/user/prefs",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (result) {
          console.log(result);
          setExtensionPreferences(result?.usersExtensionPreferences);
          setCatalogContent(result?.usersExtensionPreferences?.catalogContent);
        }
      },
      (err) => console.error(err)
    );
  }, []);

  return (
    <NoSsr>
      <Paper square className={props.classes.paperRoot}>
        <Tabs
          value={tabVal}
          onChange={handleTabValChange}
          variant="fullWidth"
          indicatorColor="primary"
          textColor="primary"
          className={props.classes.tabs}
        >
          <Tooltip title="General preferences" placement="top">
            <Tab
              className={props.classes.tab}
              icon={<SettingsCellIcon style={iconMedium} />}
              label={<span className={props.classes.tabLabel}>General</span>}
            />
          </Tooltip>
          <Tooltip title="Choose Performance Test Defaults" placement="top">
            <Tab
              className={props.classes.tab}
              icon={<FontAwesomeIcon icon={faTachometerAlt} style={iconMedium} />}
              label={<span className={props.classes.tabLabel}>Performance</span>}
            />
          </Tooltip>
          {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
          {userPrefs && providerType != "local" && (
            <Tooltip title="Remote Provider preferences" placement="top">
              <Tab
                className={props.classes.tab}
                icon={<SettingsRemoteIcon style={iconMedium} />}
                label={<span className={props.classes.tabLabel}>Remote Provider</span>}
              />
            </Tooltip>
          )}
        </Tabs>
      </Paper>
      <Paper className={props.classes.statsWrapper}>
        {tabVal === 0 && (
          <>
            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Extensions
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="CatalogContentPreference"
                    control={
                      <Switch
                        checked={catalogContent}
                        onChange={handleCatalogContentToggle}
                        color="primary"
                        classes={{
                          switchBase : props.classes.switchBase,
                          track : props.classes.track,
                          checked : props.classes.checked,
                        }}
                        data-cy="CatalogContentPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Meshery Catalog Content"
                  />
                </FormGroup>
              </FormControl>
            </div>
            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Analytics and Improvement Program
                </FormLabel>
                <FormGroup>
                  <FormControlLabel
                    key="UsageStatsPreference"
                    control={
                      <Switch
                        checked={anonymousStats}
                        onChange={handleToggle("anonymousUsageStats")}
                        color="primary"
                        classes={{
                          switchBase : props.classes.switchBase,
                          track : props.classes.track,
                          checked : props.classes.checked,
                        }}
                        data-cy="UsageStatsPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Send Anonymous Usage Statistics"
                  />
                  <FormControlLabel
                    key="PerfResultPreference"
                    control={
                      <Switch
                        checked={perfResultStats}
                        onChange={handleToggle("anonymousPerfResults")}
                        color="primary"
                        classes={{
                          switchBase : props.classes.switchBase,
                          track : props.classes.track,
                          checked : props.classes.checked,
                        }}
                        data-cy="PerfResultPreference"
                      />
                    }
                    labelPlacement="end"
                    label="Send Anonymous Performance Results"
                  />
                </FormGroup>
              </FormControl>
            </div>

            <div className={props.classes.formContainer}>
              <FormControl component="fieldset" className={props.classes.formGrp}>
                <FormLabel component="legend" className={props.classes.formLegend}>
                  Theme
                </FormLabel>

                <FormGroup>
                  <FormControlLabel
                    key="ThemePreference"
                    control={
                      <ThemeToggler
                        classes={props.classes}
                        theme={props.theme}
                        themeSetter={props.themeSetter}
                        notify={props.notify}
                      />
                    }
                    labelPlacement="end"
                    // label="Theme"
                  />
                </FormGroup>
              </FormControl>
            </div>
          </>
        )}
        {tabVal === 1 && <MesherySettingsPerformanceComponent />}
        {tabVal === 2 && userPrefs && providerType !== "local" && (
          <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
        )}
      </Paper>
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateUser : bindActionCreators(updateUser, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch),
  toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch),
});

const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get("selectedK8sContexts");
  const catalogVisibility = state.get("catalogVisibility");
  const capabilitiesRegistry = state.get("capabilitiesRegistry");
  return {
    selectedK8sContexts,
    catalogVisibility,
    capabilitiesRegistry,
  };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(UserPreference))));
