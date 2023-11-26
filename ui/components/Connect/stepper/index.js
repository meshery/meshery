import React, { useEffect, useState } from 'react';
import CustomizedSteppers from './Stepper';
import { useStyles } from '../styles';
import { helmStepContent, helmStepIcons, helmSteps } from '../connections/helm-connect/constants';
import {
  registerConnectionContent,
  registerConnectionIcons,
  registerConnectionSteps,
} from '../connections/register-connection/constants';
import { Paper, Box } from '@material-ui/core';

function ConnectionStepper({ connectionType, wizardType }) {
  const classes = useStyles();
  const [stepData, setStepData] = useState({
    stepContent: {},
    stepIcons: {},
    steps: [],
  });

  useEffect(() => {
    switch (connectionType) {
      case 'helm':
        setStepData({
          stepContent: helmStepContent,
          stepIcons: helmStepIcons,
          steps: helmSteps,
        });
        break;
      case 'register-connection':
        setStepData({
          stepContent: registerConnectionContent,
          stepIcons: registerConnectionIcons,
          steps: registerConnectionSteps,
        });
        break;
      default:
    }
  }, [connectionType]);

  return wizardType === 'modal' ? (
    <Box>{stepData.steps.length > 0 && <CustomizedSteppers stepData={stepData} />}</Box>
  ) : (
    <Box className={classes.welcomeContainer}>
      <Paper className={classes.connectionStepperContent} elevation={3}>
        {stepData.steps.length > 0 && <CustomizedSteppers stepData={stepData} />}
      </Paper>
    </Box>
  );
}

export default React.memo(ConnectionStepper);
