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
import { ComponentIcon, Loading, getSvgWhiteForComponent, processDesign } from './common';
import {
  designValidatorCommands,
  designValidatorEvents,
  useDesignSchemaValidationResults,
  useIsValidatingDesignSchema,
} from '../../machines/validator/designValidator';
import AlertIcon from '@/assets/icons/AlertIcon';
import { NOTIFICATIONCOLORS } from '@/themes/index';

const styles = (theme) => {
  const saffron = NOTIFICATIONCOLORS.WARNING;
  return {
    singleErrorRoot: {
      gap: '0.5rem',
      backgroundColor: theme.palette.secondary.mainBackground2,
      cursor: 'pointer',
      '&:hover': {
        backgroundColor: alpha(NOTIFICATIONCOLORS.WARNING, 0.25),
      },
    },
    singleError: {
      paddingInline: theme.spacing(1),
      paddingBlock: theme.spacing(1),
      marginInline: theme.spacing(0.5),
      overflow: 'hidden',
      whiteSpace: 'wrap',
    },

    componentLabel: {
      gap: '0.5rem',
      backgroundColor: saffron,
      '&:hover': {
        backgroundColor: saffron,
      },
    },
    component: {
      backgroundColor: theme.palette.secondary.mainBackground2,
      color: theme.palette.secondary.text3,
      fontFamily: 'Qanelas Soft, sans-serif',
    },
    errorList: {
      border: `solid 2px ${saffron}`,
    },

    root: {
      width: '100%',
      maxHeight: '18rem',
      marginBottom: '0.5rem',
      display: 'flex',
      flexDirection: 'column',
      gap: '1rem',
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

const ComponentErrorList = ({ component, classes, errors, validatorActor }) => {
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
        <ListItem
          disablePadding
          key={error.instancePath}
          className={classes.singleErrorRoot}
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
          <ListItemText
            primary={message(error)}
            disableTypography
            className={classes.nested}
          ></ListItemText>
        </ListItem>
      ))}
    </List>
  );
};

const ValidationResults_ = (props) => {
  const {
    errorCount,
    compCount,
    annotationCount,
    classes,
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
  return (
    <div title="DesignValidationResults">
      <List
        aria-labelledby="nested-list-subheader"
        subheader={
          <ListSubheader
            disableSticky="true"
            component="div"
            disablePadding
            className={classes.subHeader}
          >
            <Typography
              varaint="h6"
              disablePadding
              // style={{ position: "relative", left: "35px" }}
            >
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
          </ListSubheader>
        }
        className={classes.root}
      >
        {errorCount == 0 && (
          <Typography varaint="h6" align="center" disablePadding>
            No Validation errors.
          </Typography>
        )}

        {componentsWithErrors?.map((componentResult, index) => (
          <div
            style={{ margin: '0.6rem 0rem' }}
            key={index}
            ref={isCurrentComponent(componentResult) ? currentComponentErrorRef : null}
            className={classes.component}
          >
            {/*  Errors For A Component */}
            <ListItem button onClick={() => handleClick(index)} className={classes.componentLabel}>
              {/* key can be error id?? */}
              <ComponentIcon iconSrc={getSvgWhiteForComponent(componentResult.component)} />
              <ListItemText primary={componentResult.component.displayName} disableTypography />(
              {componentResult?.errors?.length}){open[index] ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Collapse in={open[index]} timeout="auto" unmountOnExit className={classes.errorList}>
              <ComponentErrorList
                component={componentResult.component}
                validatorActor={validationMachine}
                errors={componentResult.errors}
                classes={classes}
              />
            </Collapse>
          </div>
        ))}
      </List>
    </div>
  );
};

const ValidationResults = withStyles(styles)(ValidationResults_);

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
