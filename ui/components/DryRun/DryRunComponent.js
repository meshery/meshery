import React from "react";
import {
  ListItem,
  ListItemText,
  ListItemIcon,
  ListSubheader,
  List,
  Typography,
  withStyles,
  Collapse,
  alpha,
  CircularProgress
} from "@material-ui/core";
import { useEffect } from "react";
import { v4 } from "uuid";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import { useState } from "react";
import ErrorIcon from "../../assets/icons/ErrorIcon";
import { dryRunPattern } from "../../api/patterns"
import { NOTIFICATIONCOLORS } from "../../themes"
import { useNotification } from "../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../lib/event-types";

const styles = theme => {
  const error_color = NOTIFICATIONCOLORS.ERROR_DARK;
  return {
    singleErrorRoot : {
      backgroundColor : theme.palette.secondary.mainBackground2,
      cursor : "pointer",
      "&:hover" : {
        backgroundColor : alpha(error_color, 0.25)
      }
    },
    singleError : {
      paddingInline : theme.spacing(1),
      paddingBlock : theme.spacing(1),
      marginInline : theme.spacing(0.5)
    },

    componentLabel : {
      backgroundColor : error_color,
      color : "white",
      "&:hover" : {
        backgroundColor : error_color
      }
    },
    component : {
      backgroundColor : theme.palette.secondary.mainBackground2,
      color : theme.palette.secondary.text3,
      fontFamily : "Qanelas Soft, sans-serif",
      marginBlock : "0.5rem"
    },
    errorList : {
      border : `solid 2px ${error_color}`
    },

    root : {
      width : "100%",
      // maxHeight: "18rem",
      position : "relative",
      marginBottom : "0.5rem"
    },
    subHeader : {
      marginTop : "1rem",
      display : "flex",
      padding : 0,
      justifyContent : "space-between",
      width : "100%"
    }
  };
};

/**
 *
 * @param {Array} errors
 * @returns
 */
function getTotalCountOfDeploymentErrors(errors) {
  if (!errors) {
    return 0;
  }

  return errors.reduce((preCount, currEle) => preCount + currEle.errors.length, 0);
}

function getFieldPathString(fieldPath) {
  if (!fieldPath) {
    return "";
  }

  return fieldPath.split(".").splice(2).join(" > ");
}

// errors - [{type, fieldPath, message}]
// 'component' refers to MeshModel Component
// 'componentName' is assumed to be unique
const ExpandableComponentErrors = withStyles(styles)(({
  errors,
  componentName,
  classes,
}) => {

  const [isComponentAccordionOpen, setIsComponentAccordionOpen] = useState(false);

  if (!errors.length) return null;

  return (
    <div
      aria-labelledby="nested-list-subheader"
      className={classes.component}
    >
      <ListItem
        button
        className={classes.componentLabel}
        onClick={() => setIsComponentAccordionOpen(p => !p)}
      >
        <ListItemText primary={componentName} />
        {isComponentAccordionOpen ? <ExpandLess /> : <ExpandMore />}
      </ListItem>

      <Collapse
        className={classes.errorList}
        in={isComponentAccordionOpen}
        timeout="auto"
        unmountOnExit
        style={{ height : "100%" }}
      >
        <List style={{ height : "100%" }}>
          {errors?.length > 0 &&
            errors?.map((err, index) => (
              <ListItem
                disablePadding
                key={index}
                className={classes.singleErrorRoot}
                style={{ cursor : "pointer" }}
              >
                <ListItemIcon>
                  {" "}
                  <ErrorIcon
                    height="24px"
                    width="24px"
                    // bangFill="#fff"
                    fill={NOTIFICATIONCOLORS.ERROR_DARK}
                  />{" "}
                </ListItemIcon>

                {typeof err === "string" ? (
                  <ListItemText
                    disableTypography
                    primary={
                      <Typography variant="subtitle2" disablePadding>
                        {err}
                      </Typography>
                    }
                    disablePadding
                    className={classes.singleError}
                  ></ListItemText>
                ) : (
                  <>
                    <Typography
                      variant="subtitle2"
                      disablePadding
                      className={classes.nested}
                    >
                      {`Error Type: ${err.Type}`}
                    </Typography>
                    <Typography
                      variant="subtitle2"
                      disablePadding
                      className={classes.nested}
                    >
                      {`Location: ${getFieldPathString(err.FieldPath)}`}
                    </Typography>
                    <Typography
                      variant="subtitle2"
                      disablePadding
                      className={classes.nested}
                    >
                      {`Error Message: ${err.Message}`}
                    </Typography>
                  </>
                )}
              </ListItem>
            ))}
        </List>
      </Collapse>
    </div>
  );
});



export const dryRunAndFormatErrors = (design,selectedContexts) => {
  function getErrors(error) {
    if (error?.Causes && error?.Causes.length > 0) {
      return error.Causes;
    }
    // if causes aren't present use the status
    // status are strings, and causes are objects
    // so they're handled seprately in UI
    if (error?.Status) {
      return [error.Status];
    }
    return [];
  }

  return new Promise((resolve, rej) => {
    dryRunPattern(design, selectedContexts)
      .then(res => {
        const dryRunResponse = res.data.dryRunResponse;
        let errorList = [];
        if (dryRunResponse) {
          Object.keys(dryRunResponse).forEach(compName => {
            const contextsErrors = dryRunResponse?.[compName];

            if (!contextsErrors) {
              return;
            }

            Object.keys(contextsErrors).forEach(contextKey => {
              const errorAndMeta = contextsErrors[contextKey];

              if (!errorAndMeta.success) {
                errorList.push({
                  compName,
                  contextId : contextKey,
                  errors : getErrors(errorAndMeta.error)
                });
              }
            });
          });
          resolve(errorList);
          return;
        }
        resolve([]);
      })
      .catch(rej);
  });
};


const DryRunComponent = props => {
  const { design, classes, handleClose, numberOfElements ,selectedContexts } = props;
  const [deploymentErrors, setDeploymentErrors] = useState();
  const [isLoading, setIsLoading] = useState(true);
  const { notify } = useNotification()

  useEffect(() => {
    dryRunAndFormatErrors(design,selectedContexts)
      .then(errs => {
        setDeploymentErrors(errs);
      })
      .catch(e => {
        notify({
          message : "error while doing a dry run",
          event_type : EVENT_TYPES.ERROR ,
          details : e.toString()
        })
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  if (isLoading) {
    return (
      <div>
        <CircularProgress />
        <Typography variant="caption" style={{ display : "block", marginBottom : 8 }}>
          Checking for Deployment Errors
        </Typography>
      </div>
    );
  }


  const errorCount = deploymentErrors?.length;
  if (!errorCount) return null;
  return (
    <List
      aria-labelledby="nested-list-subheader"
      subheader={
        <ListSubheader
          disableSticky="true"
          component="div"
          id="nested-list-subheader"
          className={classes.subHeader}
        >
          <Typography varaint="h6" disablePadding>
            {numberOfElements} component{numberOfElements > 1 && "s"}
          </Typography>
          <Typography
            varaint="h6"
            disablePadding
            className={classes.error}
            style={{
              color : `${
                errorCount > 0
                  ? NOTIFICATIONCOLORS.ERROR_DARK
                  : NOTIFICATIONCOLORS.SUCCESS_V2
              }`
            }}
          >
            {getTotalCountOfDeploymentErrors(deploymentErrors)} error
            {getTotalCountOfDeploymentErrors(deploymentErrors) > 0 && "s"}
          </Typography>
        </ListSubheader>
      }
      className={classes.root}
    >
      {errorCount > 0 ? (
        deploymentErrors?.map(err => (
          <ExpandableComponentErrors
            key={v4()}
            componentName={err.compName}
            errors={err.errors}
            handleModalClose={handleClose}
          />
        ))
      ) : (
        <Typography varaint="h6" align="center" disablePadding>
          No deployment errors.
        </Typography>
      )}
    </List>
  );
};

export default withStyles(styles)(DryRunComponent);