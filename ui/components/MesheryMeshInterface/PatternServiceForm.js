// @ts-check
import { AppBar, IconButton, makeStyles, Toolbar, Tooltip } from "@material-ui/core";
import { Delete, HelpOutline } from "@material-ui/icons";
import SettingsIcon from '@material-ui/icons/Settings';
import React, { useEffect } from "react";
import { iconSmall } from "../../css/icons.styles";
import { pSBCr } from "../../utils/lightenOrDarkenColor";
import PatternServiceFormCore from "./PatternServiceFormCore";
import { useTheme } from "@material-ui/core/styles";

const useStyles = makeStyles(() => ({

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

/**
 * PatternServiceForm renders a form from the workloads schema
 * @param {{
 *  schemaSet: { workload: any, type: string };
 *  onSubmit: Function;
 *  onDelete: Function;
 *  namespace: string;
 *  onChange?: Function
 *  onSettingsChange?: Function;
 *  formData?: Record<String, unknown>
 *  reference?: Record<any, any>;
 *  scroll?: Boolean; // If the window should be scrolled to zero after re-rendering
 *  color ?: string;
 * }} props
 * @returns
 */
function PatternServiceForm({ formData, schemaSet, onSubmit, onDelete, reference, namespace, onSettingsChange, scroll = false, color }) {
  console.log({ schemaSet })
  const classes = useStyles();
  const theme = useTheme();

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
    schemaSet.workload.properties.labels = {
      description : "The label for the service",
      additionalProperties : {
        "type" : "string"
      },
      type : "object"
    };
    schemaSet.workload.properties.annotations = {
      description : "The annotation for the service",
      additionalProperties : {
        "type" : "string"
      },
      "type" : "object"
    };
  }, [])

  return (
    <PatternServiceFormCore
      formData={formData}
      schemaSet={schemaSet}
      onSubmit={onSubmit}
      onDelete={onDelete}
      reference={reference}
      namespace={namespace}
      onSettingsChange={onSettingsChange}
      scroll={scroll}
    >
      {(SettingsForm) => {
        return (
          <div className={classes.formWrapper}>
            <AppBar style={{
              boxShadow : `0px 2px 4px -1px "#677E88"`,
              background : `${theme.palette.type === "dark" ? "#202020" : "#647881"}`,
              position : "sticky",
              zIndex : 'auto',
            }}>
              <Toolbar
                variant="dense"
                style={{
                  padding : "0 5px",
                  paddingLeft : 16,
                  background : `linear-gradient(115deg, ${pSBCr(
                    color,
                    -20
                  )} 0%, ${color} 100%)`,
                  height : "0.7rem !important",
                }}
              >
                <SettingsIcon style={iconSmall} />
                <p
                  style={{
                    margin : "auto auto auto 10px",
                    fontSize : "16px",
                    display : "flex",
                    alignItems : "center"
                  }}
                >
                  Settings
                </p>
                {schemaSet?.workload?.description && (
                  <label htmlFor="help-button">
                    <Tooltip title={schemaSet?.workload?.description} interactive>
                      <IconButton component="span">
                        <HelpOutline width="22px" style={{ color : "#fff" }} height="22px" />
                      </IconButton>
                    </Tooltip>
                  </label>
                )}
                <IconButton
                  component="span"
                  onClick={() =>
                    // @ts-ignore
                    reference.current.delete(settings => ({
                      settings
                    }))
                  }
                >
                  <Delete width="22px" height="22px" style={{ color : "#FFF" }} />
                </IconButton>
              </Toolbar>
            </AppBar>
            <SettingsForm />
          </div>
        )
      }}
    </PatternServiceFormCore>
  )
}

export default PatternServiceForm;
