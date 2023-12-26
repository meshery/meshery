import * as React from 'react';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import { ColorlibConnector, useStyles, useColorlibStepIconStyles } from '../../../Connect/styles';
import TipsCarousel from '../../../General/TipsCarousel';
import {
  ConnectionStepperTips,
  registerConnectionContent,
  registerConnectionIcons,
  registerConnectionSteps,
} from './constants';
import Stepper from '@material-ui/core/Stepper';
import clsx from 'clsx';

function StepperIcon(props) {
  const classes = useColorlibStepIconStyles();
  const { active, completed, stepIcons } = props;

  const iconComponent = stepIcons[String(props.icon)];

  const additionalProps = {
    fill: completed ? 'white' : 'currentColor',
  };

  return (
    <div
      className={clsx(classes.icnlist, {
        [classes.active]: active,
        [classes.completed]: completed,
      })}
      style={{ position: 'relative', cursor: 'default' }}
    >
      {React.cloneElement(iconComponent, additionalProps)}
    </div>
  );
}

/* eslint-disable */

export default function CustomizedSteppers({ connectionData, onClose }) {
  const stepData = {
    stepContent: registerConnectionContent,
    stepIcons: registerConnectionIcons,
    steps: registerConnectionSteps,
  };
  const { stepContent, stepIcons, steps } = stepData;
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);

  const [sharedData, setSharedData] = React.useState(null);

  React.useEffect(() => {
    setSharedData({
      metadata: connectionData.metadata,
      capabilities: connectionData.capabilities,
      kind: connectionData.kind,
      onClose: onClose,
    });
  }, [connectionData]);

  const ActiveStepContent = stepContent[String(activeStep + 1)].component;

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const stepProps = stepContent[String(activeStep + 1)]?.props?.reduce((props, propName) => {
    props[propName] = propName === 'handleNext' ? handleNext : eval(propName);
    return props;
  }, {});

  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <Stepper
          alternativeLabel
          activeStep={activeStep}
          connector={<ColorlibConnector />}
          classes={{ root: classes.stepperContainer }}
        >
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel
                StepIconComponent={(props) => <StepperIcon {...props} stepIcons={stepIcons} />}
                classes={{ label: classes.customLabelStyle }}
              >
                {label}
              </StepLabel>
            </Step>
          ))}
        </Stepper>
      </div>
      <div
        style={{
          marginTop: '2rem',
          display: 'flex',
          minHeight: '58vh',
          minWidth: '100%',
          width: '100%',
        }}
      >
        <TipsCarousel tips={ConnectionStepperTips} />
        {React.cloneElement(ActiveStepContent, stepProps)}
      </div>
    </div>
  );
}
