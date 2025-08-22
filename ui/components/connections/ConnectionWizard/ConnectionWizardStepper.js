import React, { useState } from 'react';
import { Box, Step, styled } from '@sistent/sistent';
import TipsCarousel from '../../General/TipsCarousel';
import { getWizardSteps, getWizardContent, getWizardIcons } from './wizardConfig';
import { ColorlibConnector, CustomLabelStyle, StepperContainer } from '../styles';

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

export const ConnectionWizardStepper = ({
  connectionType,
  wizardData,
  setWizardData,
  onComplete,
  onCancel,
}) => {
  const [activeStep, setActiveStep] = useState(0);

  // Get configuration for the specific connection type
  const steps = getWizardSteps(connectionType);
  const stepContent = getWizardContent(connectionType);
  const stepIcons = getWizardIcons(connectionType);

  const handleNext = () => {
    if (activeStep < steps.length - 1) {
      setActiveStep((prevStep) => prevStep + 1);
    }
  };

  const handleBack = () => {
    if (activeStep > 0) {
      setActiveStep((prevStep) => prevStep - 1);
    }
  };

  const handleStepComplete = (stepData) => {
    setWizardData((prevData) => ({
      ...prevData,
      ...stepData,
    }));

    if (activeStep === steps.length - 1) {
      // Final step completed
      onComplete(wizardData);
    } else {
      handleNext();
    }
  };

  const ActiveStepContent = stepContent[activeStep];

  const stepProps = {
    activeStep,
    wizardData,
    setWizardData,
    onNext: handleNext,
    onBack: handleBack,
    onComplete: handleStepComplete,
    onCancel,
    isFirstStep: activeStep === 0,
    isLastStep: activeStep === steps.length - 1,
  };

  return (
    <StepperLayout>
      <StepperHeader>
        <StepperContainer
          alternativeLabel
          activeStep={activeStep}
          connector={<ColorlibConnector />}
        >
          {steps.map((label, index) => (
            <Step key={label}>
              <CustomLabelStyle
                StepIconComponent={(props) => (
                  <StepperIcon {...props} stepIcons={stepIcons} icon={index + 1} />
                )}
              >
                {label}
              </CustomLabelStyle>
            </Step>
          ))}
        </StepperContainer>
      </StepperHeader>
      <StepperContent>
        <TipsCarousel
          tips={[
            'Review all discovered connections before proceeding.',
            'Pre-flight checks help ensure successful connections.',
            'You have full control over which connections to manage.',
          ]}
        />
        {React.cloneElement(ActiveStepContent, stepProps)}
      </StepperContent>
    </StepperLayout>
  );
};
