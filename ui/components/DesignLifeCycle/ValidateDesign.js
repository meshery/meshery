import React, { useRef, useEffect } from 'react';
import { List, ListItemText, ListItemIcon, Typography, Collapse, useTheme } from '@layer5/sistent';
import { ExpandLess, ExpandMore } from '@mui/icons-material';
import { ComponentIcon, Loading, getSvgWhiteForComponent, processDesign } from './common';
import {
  designValidatorCommands,
  designValidatorEvents,
  useDesignSchemaValidationResults,
  useIsValidatingDesignSchema,
} from '../../machines/validator/designValidator';
import AlertIcon from '@/assets/icons/AlertIcon';
import { NOTIFICATIONCOLORS } from '@/themes/index';
import {
  ValidatedComponent,
  ValidationErrorListItem,
  ComponentValidationListItem,
  ValidationResultsListWrapper,
  ValidationSubHeader,
} from './styles';

const ComponentErrorList = ({ component, errors, validatorActor }) => {
  const onErrorTap = (error) => {
    validatorActor.send(
      designValidatorEvents.tapOnError({ error, type: 'schemaValidation', component }),
    );
  };

  const formatReadableFieldPath = (instancePath) => instancePath.split('/').join('.');

  const message = (error) => {
    const fieldPath = formatReadableFieldPath(error.instancePath);
    return fieldPath ? `Field "${fieldPath}" ${error.message}` : error.message;
  };

  return (
    <List>
      {errors.map((error) => (
        <ValidationErrorListItem
          disablePadding
          key={error.instancePath}
          onClick={() => onErrorTap(error)}
        >
          <ListItemIcon>
            {' '}
            <AlertIcon
              height="27px"
              width="30px"
              outline="#fff"
              bangFill="#fff"
              fill={NOTIFICATIONCOLORS.WARNING}
            />{' '}
          </ListItemIcon>
          <ListItemText primary={message(error)} disableTypography></ListItemText>
        </ValidationErrorListItem>
      ))}
    </List>
  );
};

const ValidationResults_ = (props) => {
  const {
    errorCount,
    compCount,
    annotationCount,
    validationResults,
    currentNodeId,
    validationMachine,
  } = props;

  const componentsWithErrors = Object.values(validationResults).filter(
    (result) => result.errors?.length && !result?.component?.metatadata?.isAnnotation,
  );

  console.log('componentsWithErrors', componentsWithErrors);

  const isCurrentComponent = (err) => err.component.id == currentNodeId;
  const [open, setOpen] = React.useState(componentsWithErrors.map(isCurrentComponent));

  const currentComponentErrorRef = useRef();
  const handleClick = (index) => {
    let updatedState = [...open];
    updatedState[index] = !updatedState[index];
    setOpen(updatedState);
  };

  useEffect(() => {
    // 👇 Will scroll smoothly to the top of component
    if (currentComponentErrorRef.current) {
      currentComponentErrorRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [currentComponentErrorRef]);
  const theme = useTheme();
  return (
    <div title="DesignValidationResults">
      <ValidationResultsListWrapper
        aria-labelledby="nested-list-subheader"
        subheader={
          <ValidationSubHeader disableSticky="true" component="div" disablePadding>
            <Typography varaint="h6" disablePadding style={{ color: theme.palette.text.disabled }}>
              {compCount} component{compCount > 1 ? 's' : ''} and {annotationCount} annotations
            </Typography>

            <Typography
              varaint="h6"
              disablePadding
              style={{
                color: `${
                  errorCount > 0 ? NOTIFICATIONCOLORS.WARNING : NOTIFICATIONCOLORS.SUCCESS_V2
                }`,
              }}
            >
              {errorCount} error{errorCount === 1 ? '' : 's'}
            </Typography>
          </ValidationSubHeader>
        }
      >
        {errorCount == 0 && (
          <Typography varaint="h6" align="center" style={{ marginBlock: '1rem' }} disablePadding>
            No validation errors.
          </Typography>
        )}

        {componentsWithErrors?.map((componentResult, index) => (
          <ValidatedComponent
            key={index}
            ref={isCurrentComponent(componentResult) ? currentComponentErrorRef : null}
          >
            {/*  Errors For A Component */}
            <ComponentValidationListItem button onClick={() => handleClick(index)}>
              {/* key can be error id?? */}
              <ComponentIcon iconSrc={getSvgWhiteForComponent(componentResult.component)} />
              <ListItemText primary={componentResult.component.displayName} disableTypography />(
              {componentResult?.errors?.length}){open[index] ? <ExpandLess /> : <ExpandMore />}
            </ComponentValidationListItem>
            <Collapse
              in={open[index]}
              timeout="auto"
              unmountOnExit
              style={{ border: `solid 2px ${theme.palette.background.cta.default}` }}
            >
              <ComponentErrorList
                component={componentResult.component}
                validatorActor={validationMachine}
                errors={componentResult.errors}
              />
            </Collapse>
          </ValidatedComponent>
        ))}
      </ValidationResultsListWrapper>
    </div>
  );
};

const ValidationResults = ValidationResults_;

/**
 *
 * @param {Design} design is the design file to be validated
 * @returns
 */
export const ValidateDesign = ({ design, currentNodeId, validationMachine }) => {
  const validationResults = useDesignSchemaValidationResults(validationMachine);
  const isValidating = useIsValidatingDesignSchema(validationMachine);

  const designName = design.name;
  const { configurableComponents, annotationComponents } = processDesign(design);

  const totalErrors = Object.values(validationResults || {}).reduce(
    (acc, serviceResult) => acc + (serviceResult?.errors?.length || 0),
    0,
  );

  useEffect(() => {
    validationMachine.send(designValidatorCommands.validateDesignSchema({ design }));
  }, []);

  if (isValidating) {
    return <Loading message="Validating Design" />;
  }

  if (!validationResults) {
    return null;
  }

  if (typeof validationResults === 'string') {
    return <Typography variant="h6">{validationResults}</Typography>;
  }

  return (
    <ValidationResults
      validationResults={validationResults}
      currentNodeId={currentNodeId}
      errorCount={totalErrors}
      compCount={configurableComponents.length}
      annotationCount={annotationComponents.length}
      design={designName}
      validationMachine={validationMachine}
    />
  );
};
