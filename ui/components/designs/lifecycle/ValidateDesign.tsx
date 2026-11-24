import React, { useRef, useEffect } from 'react';
import { List, ListItemText, ListItemIcon, Typography, Collapse } from '@sistent/sistent';
import { ExpandLess, ExpandMore } from '@/assets/icons';
import { ComponentIcon, Loading, getSvgWhiteForComponent, processDesign } from './common';
import {
  designValidatorCommands,
  designValidatorEvents,
  useDesignSchemaValidationResults,
  useIsValidatingDesignSchema,
} from '../../../machines/validator/designValidator';
import AlertIcon from '@/assets/icons/AlertIcon';
import { useTheme } from '@/theme';
import {
  ValidatedComponent,
  ValidationErrorListItem,
  ComponentValidationListItem,
  ValidationResultsListWrapper,
  ValidationSubHeader,
} from './styles';

// Each error deep-links into the designer: tapping it sends TAP_ON_ERROR to the
// validation machine, which emits it to the embedding designer. Kanvas Designer
// listens for that event and opens the component's configurator focused on the
// failing property. Without a designer to navigate into, an error is plain text.
const ComponentErrorList = ({ component, errors, validatorActor, canNavigateToErrors }) => {
  const theme = useTheme();

  const navigateToError = (error) => {
    validatorActor.send(
      designValidatorEvents.tapOnError({ error, type: 'schemaValidation', component }),
    );
  };

  const formatReadableFieldPath = (instancePath) => instancePath.split('/').join('.');

  const message = (error) => {
    const fieldPath = formatReadableFieldPath(error.instancePath);
    return fieldPath ? `Field "${fieldPath}" ${error.message}` : error.message;
  };

  const navigationProps = (error) =>
    canNavigateToErrors
      ? {
          role: 'button',
          tabIndex: 0,
          onClick: () => navigateToError(error),
          onKeyDown: (event) => {
            if (event.key === 'Enter' || event.key === ' ') {
              event.preventDefault();
              navigateToError(error);
            }
          },
        }
      : {};

  return (
    <List>
      {errors.map((error, index) => (
        <ValidationErrorListItem
          disablePadding
          clickable={canNavigateToErrors}
          // instancePath alone is not unique: every `required` error on one
          // object shares it.
          key={`${error.instancePath}-${error.keyword}-${index}`}
          data-testid="validation-error"
          {...navigationProps(error)}
        >
          <ListItemIcon>
            {' '}
            <AlertIcon
              height="27px"
              width="30px"
              outline="#fff"
              bangFill="#fff"
              fill={theme.palette.warning.main}
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
    canNavigateToErrors = false,
  } = props;

  const componentsWithErrors = Object.values(validationResults).filter(
    (result) => result.errors?.length && !result?.component?.metadata?.isAnnotation,
  );

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
            <Typography variant="h6" disablePadding style={{ color: theme.palette.text.default }}>
              {compCount} component{compCount > 1 ? 's' : ''} and {annotationCount} annotations
            </Typography>

            <Typography
              variant="h6"
              disablePadding
              style={{
                color: `${
                  errorCount > 0 ? theme.palette.warning.main : theme.palette.success.main
                }`,
              }}
            >
              {errorCount} error{errorCount === 1 ? '' : 's'}
            </Typography>
          </ValidationSubHeader>
        }
      >
        {errorCount == 0 && (
          <Typography variant="h6" align="center" style={{ marginBlock: '1rem' }} disablePadding>
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
                canNavigateToErrors={canNavigateToErrors}
              />
            </Collapse>
          </ValidatedComponent>
        ))}
      </ValidationResultsListWrapper>
    </div>
  );
};

export const ValidationResults = ValidationResults_;

/**
 * ValidateDesign validates a design's component configurations against their
 * schemas and lists the failures per component.
 *
 * @param design the design to validate
 * @param currentNodeId the component whose errors start expanded and scrolled into view
 * @param validationMachine the design validation machine; Meshery UI injects the same
 *   machine into extensions as `designValidationMachine`
 * @param canNavigateToErrors whether each error deep-links to its component and
 *   property. True by default because the component's purpose when injected into an
 *   extension is Kanvas Designer, which handles TAP_ON_ERROR. A surface with no designer
 *   to navigate into - every Meshery UI page - passes false so errors do not look
 *   clickable. Never derive this from the validator's returnAddress: that names who
 *   receives validation replies, and this component's own validation request clears it.
 */
export const ValidateDesign = ({
  design,
  currentNodeId,
  validationMachine,
  canNavigateToErrors = true,
}) => {
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

  if (isValidating || !validationResults) {
    return <Loading message={isValidating ? 'Validating Design' : 'Loading validation form'} />;
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
      canNavigateToErrors={canNavigateToErrors}
    />
  );
};
