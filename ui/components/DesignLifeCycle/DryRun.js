import React, { useRef, useEffect } from 'react';
import { List, ListItemText, ListItemIcon, Typography, Collapse, useTheme } from '@layer5/sistent';
import { ExpandLess, ExpandMore } from '@mui/icons-material';
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
import { FormatStructuredData } from '../DataFormatter';
import {
  DryRunComponentLabel,
  DryRunComponentStyled,
  DryRunErrorContainer,
  DryRunSignleError,
  DryRunRootListStyled,
  ValidationSubHeader,
} from './styles';

function breakCapitalizedWords(input) {
  // Use regular expression to split capitalized words
  // into separate words with space in between
  return input.replace(/([a-z])([A-Z])/g, '$1 $2');
}

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
const ExpandableComponentErrors = ({
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
  const theme = useTheme();

  if (!componentName) return null;

  if (!errors.length) return null;
  return (
    <DryRunComponentStyled
      aria-labelledby="nested-list-subheader"
      ref={isCurrentComponent(componentName) ? currentComponentErrorRef : null}
    >
      <DryRunComponentLabel button onClick={() => setIsComponentAccordionOpen((p) => !p)}>
        {componentIcon && <ComponentIcon iconSrc={componentIcon} />}
        <ListItemText primary={componentName} />({errors.length})
        {isComponentAccordionOpen ? <ExpandLess /> : <ExpandMore />}
      </DryRunComponentLabel>

      <Collapse
        in={isComponentAccordionOpen}
        timeout="auto"
        unmountOnExit
        style={{ height: '100%', border: `solid 2px ${theme.palette.background.error.default}` }}
      >
        <List style={{ height: '100%' }}>
          {errors?.length > 0 &&
            errors?.map((err, index) => (
              <DryRunErrorContainer
                disablePadding
                key={index}
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
                  <DryRunSignleError
                    disableTypography
                    primary={
                      <Typography variant="subtitle2" disablePadding>
                        {err}
                      </Typography>
                    }
                    disablePadding
                  ></DryRunSignleError>
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
              </DryRunErrorContainer>
            ))}
        </List>
      </Collapse>
    </DryRunComponentStyled>
  );
};

export const FormatDryRunResponse = ({
  dryRunErrors,
  configurableComponentsCount,
  annotationComponentsCount,
  validationMachine,
  classes,
  currentComponentName,
}) => {
  const totalDryRunErrors = getTotalCountOfDeploymentErrors(dryRunErrors);
  const theme = useTheme();
  const canShowComponentCount =
    annotationComponentsCount !== undefined && annotationComponentsCount !== undefined;

  const dryRunRequestErrors = dryRunErrors
    .filter((error) => error.type === 'RequestError')
    .flatMap((error) => error?.errors || []);
  const componentErrors = dryRunErrors.filter((error) => error.type === 'ComponentError');

  console.log('dryRunRequestErrors', dryRunRequestErrors);

  return (
    <DryRunRootListStyled
      aria-labelledby="nested-list-subheader"
      subheader={
        <ValidationSubHeader disableSticky="true" component="div" id="nested-list-subheader">
          {canShowComponentCount && (
            <Typography varaint="h6" disablePadding style={{ color: theme.palette.text.disabled }}>
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
        </ValidationSubHeader>
      }
    >
      <div style={{ padding: '1rem' }}> </div>
      {dryRunRequestErrors?.map((error, index) => (
        <FormatStructuredData key={index} data={error?.data || error?.error} />
      ))}

      {totalDryRunErrors > 0 ? (
        componentErrors?.map((err) => (
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
    </DryRunRootListStyled>
  );
};

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
