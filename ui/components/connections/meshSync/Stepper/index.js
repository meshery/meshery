import * as React from 'react';
import { Step } from '@mui/material';
import { Box, styled } from '@layer5/sistent';
import TipsCarousel from '../../../General/TipsCarousel';
import {
  ConnectionStepperTips,
  registerConnectionContent,
  registerConnectionIcons,
  registerConnectionSteps,
} from './constants';
import { ColorlibConnector, CustomLabelStyle, StepperContainer } from '../../styles';

const StepIconWrapper = styled('div')(({ theme, active, completed }) => ({
  backgroundColor: theme.palette.background.card,
  zIndex: 1,
  color: '#fff',
  width: 50,
  height: 50,
  display: 'flex',
  border: `.2rem solid ${theme.palette.divider}`,
  borderRadius: '50%',
  justifyContent: 'center',
  alignItems: 'center',
  cursor: 'default',
  '@media (max-width:780px)': {
    width: 40,
    height: 40,
  },
  ...(active && {
    background: theme.palette.background.tabs,
    color: '#3C494E',
    border: '.2rem solid #00B39F',
    transition: 'all 0.5s ease-in',
  }),
  ...(completed && {
    border: '.2rem solid #00B39F',
    background: '#00B39F',
    transition: 'all 0.5s ease-in',
  }),
}));

const StepperLayout = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
});

const StepperHeader = styled(Box)({
  display: 'flex',
  justifyContent: 'center',
});

const StepperContent = styled(Box)({
  marginTop: '2rem',
  display: 'flex',
  minHeight: '58vh',
  minWidth: '100%',
  width: '100%',
});

function StepperIcon({ active, completed, stepIcons, icon }) {
  const iconComponent = stepIcons[String(icon)];
  const additionalProps = {
    fill: completed ? 'white' : 'currentColor',
  };

  return (
    <StepIconWrapper active={active} completed={completed}>
      {React.cloneElement(iconComponent, additionalProps)}
    </StepIconWrapper>
  );
}

export default function CustomizedSteppers({ sharedData, setSharedData, connectionData, onClose }) {
  const [activeStep, setActiveStep] = React.useState(0);
  const stepData = {
    stepContent: registerConnectionContent,
    stepIcons: registerConnectionIcons,
    steps: registerConnectionSteps,
  };
  const { stepContent, stepIcons, steps } = stepData;

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
    <StepperLayout>
      <StepperHeader>
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
      </StepperHeader>
      <StepperContent>
        <TipsCarousel tips={ConnectionStepperTips} />
        {React.cloneElement(ActiveStepContent, stepProps)}
      </StepperContent>
    </StepperLayout>
  );
}
