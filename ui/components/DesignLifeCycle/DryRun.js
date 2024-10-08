import React, { useRef, useEffect } from 'react';
import {
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  ListSubheader,
  Typography,
  Collapse,
  alpha,
  withStyles,
} from '@material-ui/core';
import { ExpandLess, ExpandMore } from '@material-ui/icons';
import { useState } from 'react';
import { ComponentIcon, DEPLOYMENT_TYPE, Loading, processDesign } from './common';
import {
  designValidatorCommands,
  designValidatorEvents,
  useDryRunValidationResults,
  useIsValidatingDryRun,
} from '../../machines/validator/designValidator';
import { ErrorIcon } from '@layer5/sistent';
import { NOTIFICATIONCOLORS } from '@/themes/index';

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
      gap: '0.5rem',
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
 * @returns
 */
export function getTotalCountOfDeploymentErrors(errors) {
  if (!errors) {
    return 0;
  }

  return errors.reduce((preCount, currEle) => preCount + currEle.errors.length, 0);
}

function getFieldPathString(fieldPath) {
  if (!fieldPath) {
    return '';
  }

  return fieldPath.split('.').splice(2).join(' > ');
}

// errors - [{type, fieldPath, message}]
// 'component' refers to MeshModel Component
// 'componentName' is assumed to be unique
const ExpandableComponentErrors = withStyles(styles)(({
  errors,
  component,
  componentName,
  classes,
  validationMachine,
  currentComponentName, // if dry run is initiated by clicking on node's error badge
}) => {
  const componentIcon = component ? `/${component?.styles?.svgWhite}` : null;

  const isCurrentComponent = (name) => name == currentComponentName;
  const [isComponentAccordionOpen, setIsComponentAccordionOpen] = useState(
    isCurrentComponent(componentName),
  );
  const currentComponentErrorRef = useRef(null);

  const onErrorTap = (error) => {
    if (!validationMachine) {
      return;
    }

    validationMachine.send(
      designValidatorEvents.tapOnError({
        error,
        type: 'dryRun',
        component: component,
      }),
    );
  };

  useEffect(() => {
    // 👇 Will scroll smoothly to the top of component
    if (currentComponentErrorRef.current) {
      currentComponentErrorRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [currentComponentErrorRef]);

  if (!componentName) return null;

  if (!errors.length) return null;

  return (
    <div
      aria-labelledby="nested-list-subheader"
      className={classes.component}
      ref={isCurrentComponent(componentName) ? currentComponentErrorRef : null}
    >
      <ListItem
        button
        className={classes.componentLabel}
        onClick={() => setIsComponentAccordionOpen((p) => !p)}
      >
        {componentIcon && <ComponentIcon iconSrc={componentIcon} />}
        <ListItemText primary={componentName} />({errors.length})
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
            errors?.map((err, index) => (
              <ListItem
                disablePadding
                key={index}
                className={classes.singleErrorRoot}
                style={{ cursor: 'pointer' }}
                onClick={() => onErrorTap(err)}
              >
                <ListItemIcon>
                  {' '}
                  <ErrorIcon
                    height="24px"
                    width="24px"
                    bangFill="#fff"
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

export const FormatDryRunResponse = withStyles(styles)(({
  dryRunErrors,
  configurableComponentsCount,
  annotationComponentsCount,
  validationMachine,
  classes,
  currentComponentName,
}) => {
  const totalDryRunErrors = getTotalCountOfDeploymentErrors(dryRunErrors);

  const canShowComponentCount =
    annotationComponentsCount !== undefined && annotationComponentsCount !== undefined;
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
          {canShowComponentCount && (
            <Typography
              varaint="h6"
              disablePadding
              // style={{ position: "relative", left: "35px" }}
            >
              {configurableComponentsCount} component{configurableComponentsCount > 1 ? 's' : ''}{' '}
              and {annotationComponentsCount} annotations
            </Typography>
          )}
          <Typography
            varaint="h6"
            disablePadding
            className={classes.error}
            style={{
              color: `${
                totalDryRunErrors > 0
                  ? NOTIFICATIONCOLORS.ERROR_DARK
                  : NOTIFICATIONCOLORS.SUCCESS_V2
              }`,
            }}
          >
            {`${totalDryRunErrors} error(s)`}
          </Typography>
        </ListSubheader>
      }
      className={classes.root}
    >
      {totalDryRunErrors > 0 ? (
        dryRunErrors?.map((err) => (
          <ExpandableComponentErrors
            key={`${err.compName}`}
            component={err.component}
            componentName={err.compName}
            validationMachine={validationMachine}
            currentComponentName={currentComponentName}
            errors={err.errors}
          />
        ))
      ) : (
        <Typography varaint="h6" align="center" disablePadding>
          No deployment errors.
        </Typography>
      )}
    </List>
  );
});

const DryRunComponent = (props) => {
  const {
    design,
    validationMachine,
    currentComponentName,
    selectedK8sContexts,
    deployment_type,
    includeDependencies,
  } = props;

  const dryRunErrors = useDryRunValidationResults(validationMachine);
  const isLoading = useIsValidatingDryRun(validationMachine);
  const { configurableComponents, annotationComponents } = processDesign(design);

  useEffect(() => {
    const dryRunCommand =
      deployment_type == DEPLOYMENT_TYPE.DEPLOY
        ? designValidatorCommands.dryRunDesignDeployment
        : designValidatorCommands.dryRunDesignUnDeployment;
    validationMachine.send(
      dryRunCommand({
        design,
        k8sContexts: selectedK8sContexts,
        includeDependencies,
      }),
    );
  }, [includeDependencies]);

  if (isLoading) {
    return <Loading message="Performing a dry run" />;
  }

  if (!dryRunErrors) {
    return null;
  }

  if (typeof dryRunErrors === 'string') {
    return (
      <div>
        <Typography variant="caption" style={{ display: 'block', margin: 8 }}>
          {dryRunErrors}
        </Typography>
      </div>
    );
  }

  return (
    <FormatDryRunResponse
      dryRunErrors={dryRunErrors}
      annotationComponentsCount={annotationComponents.length}
      configurableComponentsCount={configurableComponents.length}
      validationMachine={validationMachine}
      currentComponentName={currentComponentName}
    />
  );
};

export const DryRunDesign = ({
  handleClose,
  currentComponentName,
  design,
  validationMachine,
  selectedK8sContexts,
  deployment_type,
  includeDependencies,
}) => {
  if (!design) {
    return null;
  }

  return (
    <DryRunComponent
      design={design}
      includeDependencies={includeDependencies}
      validationMachine={validationMachine}
      deployment_type={deployment_type}
      handleClose={handleClose}
      selectedK8sContexts={selectedK8sContexts}
      currentComponentName={currentComponentName}
    />
  );
};
