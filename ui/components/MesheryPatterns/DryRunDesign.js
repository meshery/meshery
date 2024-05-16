import React from 'react';
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
  CircularProgress,
} from '@material-ui/core';
import { useEffect } from 'react';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import { useState } from 'react';
import ErrorIcon from '../../assets/icons/ErrorIcon';
import { NOTIFICATIONCOLORS } from '../../themes';
import { getComponentsinFile } from '@/utils/utils';
import { useDeployPatternMutation, useUndeployPatternMutation } from '@/rtk-query/design';

function breakCapitalizedWords(input) {
  // Use regular expression to split capitalized words
  // into separate words with space in between
  return input.replace(/([a-z])([A-Z])/g, '$1 $2');
}
const styles = (theme) => {
  const error_color = NOTIFICATIONCOLORS.ERROR_DARK;
  return {
    singleErrorRoot: {
      backgroundColor: theme.palette.secondary.mainBackground2,
      cursor: 'pointer',
      '&:hover': {
        backgroundColor: alpha(error_color, 0.25),
      },
    },
    singleError: {
      paddingInline: theme.spacing(1),
      paddingBlock: theme.spacing(1),
      marginInline: theme.spacing(0.5),
    },

    componentLabel: {
      backgroundColor: error_color,
      color: 'white',
      '&:hover': {
        backgroundColor: error_color,
      },
    },
    component: {
      backgroundColor: theme.palette.secondary.mainBackground2,
      color: theme.palette.secondary.text3,
      fontFamily: 'Qanelas Soft, sans-serif',
      marginBlock: '0.5rem',
    },
    errorList: {
      border: `solid 2px ${error_color}`,
    },

    root: {
      width: '100%',
      // maxHeight: "18rem",
      position: 'relative',
      marginBottom: '0.5rem',
    },
    subHeader: {
      marginTop: '1rem',
      display: 'flex',
      padding: 0,
      justifyContent: 'space-between',
      width: '100%',
    },
  };
};

/**
 *
 * @param {Array} errors
 * @returns {Number} Total count of errors in the deployment
 */
function getTotalCountOfDeploymentErrors(errors) {
  if (!errors) {
    return 0;
  }

  return errors.reduce((preCount, currEle) => preCount + currEle.errors.length, 0);
}

/**
 * returns the field path string from the field path by removing the first two elements
 * @param {String} fieldPath
 * @returns {String} fieldPathString
 */
function getFieldPathString(fieldPath) {
  if (!fieldPath) {
    return '';
  }

  return fieldPath.split('.').splice(2).join('.');
}

// errors - [{type, fieldPath, message}]
// 'component' refers to Model Component
// 'componentName' is assumed to be unique
const ExpandableComponentErrors = withStyles(styles)(({ errors, componentName, classes }) => {
  const [isComponentAccordionOpen, setIsComponentAccordionOpen] = useState(false);

  if (!errors.length) return null;

  return (
    <div aria-labelledby="nested-list-subheader" className={classes.component}>
      <ListItem
        button
        className={classes.componentLabel}
        onClick={() => setIsComponentAccordionOpen((p) => !p)}
      >
        <ListItemText primary={componentName} />
        {isComponentAccordionOpen ? <ExpandLess /> : <ExpandMore />}
      </ListItem>

      <Collapse
        className={classes.errorList}
        in={isComponentAccordionOpen}
        timeout="auto"
        unmountOnExit
        style={{ height: '100%' }}
      >
        <List style={{ height: '100%' }}>
          {errors?.length > 0 &&
            errors?.map((err) => (
              <ListItem
                disablePadding
                key={JSON.stringify(err)}
                className={classes.singleErrorRoot}
                style={{ cursor: 'pointer' }}
              >
                <ListItemIcon>
                  {' '}
                  <ErrorIcon
                    height="24px"
                    width="24px"
                    // bangFill="#fff"
                    fill={NOTIFICATIONCOLORS.ERROR_DARK}
                  />{' '}
                </ListItemIcon>

                {typeof err === 'string' ? (
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
                  <div
                    style={{
                      display: 'flex',
                      padding: '0.5rem',
                      flexDirection: 'column',
                    }}
                  >
                    <Typography
                      variant="subtitle2"
                      disablePadding
                      style={{ fontWeight: 'bold', marginBottom: '0.25rem' }}
                      className={classes.nested}
                    >
                      {breakCapitalizedWords(err.Type || '')}: {getFieldPathString(err.FieldPath)}
                    </Typography>

                    <Typography variant="subtitle2" disablePadding className={classes.nested}>
                      {err.Message}
                    </Typography>
                  </div>
                )}
              </ListItem>
            ))}
        </List>
      </Collapse>
    </div>
  );
});

export const formatDryRunResponse = (dryRunResponse) => {
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

  let errorList = [];
  if (dryRunResponse) {
    Object.keys(dryRunResponse).forEach((compName) => {
      const contextsErrors = dryRunResponse?.[compName];

      if (!contextsErrors) {
        return;
      }

      Object.keys(contextsErrors).forEach((contextKey) => {
        const errorAndMeta = contextsErrors[contextKey];

        if (!errorAndMeta.success) {
          errorList.push({
            compName,
            contextId: contextKey,
            errors: getErrors(errorAndMeta.error),
          });
        }
      });
    });
  }
  return errorList;
};

export const FormatDryRunResponse = withStyles(styles)(({
  dryRunResponse,
  numberOfComponentsInDesign,
  onErrorClick,
  classes,
}) => {
  const deploymentErrors = formatDryRunResponse(dryRunResponse);

  const totalErrors = getTotalCountOfDeploymentErrors(deploymentErrors);

  const errorCount = deploymentErrors?.length || 0;

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
            {numberOfComponentsInDesign !== null
              ? `${numberOfComponentsInDesign} component(s)`
              : ''}
          </Typography>
          <Typography
            varaint="h6"
            disablePadding
            className={classes.error}
            style={{
              color: `${
                errorCount > 0 ? NOTIFICATIONCOLORS.ERROR_DARK : NOTIFICATIONCOLORS.SUCCESS_V2
              }`,
              fontWeight: '600',
            }}
          >
            {totalErrors} error
            {totalErrors > 0 && 's'}
          </Typography>
        </ListSubheader>
      }
      className={classes.root}
    >
      {errorCount > 0 ? (
        deploymentErrors?.map((err) => (
          <ExpandableComponentErrors
            key={err.compName}
            componentName={err.compName}
            errors={err.errors}
            handleModalClose={onErrorClick}
          />
        ))
      ) : (
        <Typography varaint="h6" align="center" disablePadding>
          No dry run errors.
        </Typography>
      )}
    </List>
  );
});

const DryRunComponent = (props) => {
  const { pattern_file, pattern_id, handleErrors, handleClose, selectedContexts, dryRunType } =
    props;
  const [isLoading, setIsLoading] = useState(false);
  const numberOfElements = getComponentsinFile(pattern_file);
  const [dryRunResponse, setDryRunResponse] = useState(null);

  const useDryRunMutation =
    dryRunType == 'deploy' ? useDeployPatternMutation : useUndeployPatternMutation;
  const [dryRunMutation, { error: failedToPerformDryRun }] = useDryRunMutation();

  useEffect(async () => {
    setIsLoading(true);
    try {
      const dryRunResults = await dryRunMutation({
        pattern_file,
        pattern_id,
        selectedK8sContexts: selectedContexts,
        dryRun: true,
        verify: false,
      });
      setDryRunResponse(dryRunResults.data?.dryRunResponse);
      console.log('dry run response', dryRunResults.data?.dryRunResponse);
      const errors = formatDryRunResponse(dryRunResults.data?.dryRunResponse);
      handleErrors?.(errors);
    } finally {
      setIsLoading(false);
    }
  }, []);

  if (isLoading) {
    return (
      <div>
        <CircularProgress />
        <Typography variant="caption" style={{ display: 'block', marginBottom: 8 }}>
          Performing dry run for {dryRunType}
        </Typography>
      </div>
    );
  }

  // If DryRun fails, show a message
  if (failedToPerformDryRun) {
    return (
      <Typography variant="caption" style={{ display: 'block', marginBottom: 8 }}>
        Failed to perform dry run
      </Typography>
    );
  }

  return (
    <FormatDryRunResponse
      dryRunResponse={dryRunResponse}
      numberOfComponentsInDesign={numberOfElements}
      onErrorClick={handleClose}
    />
  );
};

export const DryRunDesign = withStyles(styles)(DryRunComponent);
