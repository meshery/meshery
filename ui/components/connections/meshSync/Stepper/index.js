import * as React from 'react';
import { useColorlibStepIconStyles } from '../../../Connect/styles';
import TipsCarousel from '../../../General/TipsCarousel';
import {
  ConnectionStepperTips,
  registerConnectionContent,
  registerConnectionIcons,
  registerConnectionSteps,
} from './constants';
import clsx from 'clsx';
import { Step } from '@mui/material';
import { ColorlibConnector, CustomLabelStyle, StepperContainer } from '../../styles';

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

export default function CustomizedSteppers({
  sharedData,
  setSharedData,
  connectionData,
  onClose,
  handleRegistrationComplete,
}) {
  const stepData = {
    stepContent: registerConnectionContent,
    stepIcons: registerConnectionIcons,
    steps: registerConnectionSteps,
  };
  const { stepContent, stepIcons, steps } = stepData;
  const [activeStep, setActiveStep] = React.useState(0);

  React.useEffect(() => {
    setSharedData({
      metadata: connectionData.metadata,
      capabilities: connectionData.capabilities,
      kind: connectionData.kind,
    });
  }, [connectionData]);

  React.useEffect(() => {
    setSharedData((prevState) => ({
      ...prevState,
      onClose: onClose,
    }));
  }, [sharedData]);
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
        <StepperContainer
          alternativeLabel
          activeStep={activeStep}
          connector={<ColorlibConnector />}
        >
          {steps.map((label) => (
            <Step key={label}>
              <CustomLabelStyle
                StepIconComponent={(props) => <StepperIcon {...props} stepIcons={stepIcons} />}
              >
                {label}
              </CustomLabelStyle>
            </Step>
          ))}
        </StepperContainer>
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
