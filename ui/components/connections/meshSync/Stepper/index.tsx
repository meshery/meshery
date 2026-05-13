import * as React from 'react';
import { Box, styled, Step } from '@sistent/sistent';
import TipsCarousel from '../../../General/TipsCarousel';
import {
  ConnectionStepperTips,
  registerConnectionContent,
  registerConnectionIcons,
  registerConnectionSteps,
} from './constants';
import { ColorlibConnector, CustomLabelStyle, StepperContainer } from '../../styles';

interface StepIconWrapperProps {
  active?: boolean;
  completed?: boolean;
}

const StepIconWrapper = styled('div')<StepIconWrapperProps>(({ theme, active, completed }) => ({
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

function StepperIcon({
  active,
  completed,
  stepIcons,
  icon,
}: StepIconWrapperProps & {
  stepIcons: Record<string, React.ReactElement>;
  icon: React.ReactNode;
}) {
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

type SharedData = {
  metadata?: Record<string, unknown>;
  capabilities?: {
    urls?: string[];
    connection?: boolean;
  };
  kind?: string;
  resourceID?: string;
  onClose?: () => void;
  [key: string]: unknown;
};

type CustomizedSteppersProps = {
  sharedData: SharedData;
  setSharedData: React.Dispatch<React.SetStateAction<SharedData>>;
  connectionData: SharedData;
  onClose: () => void;
  handleRegistrationComplete: (resourceId?: string) => void;
};

export default function CustomizedSteppers({
  sharedData,
  setSharedData,
  connectionData,
  onClose,
  handleRegistrationComplete,
}: CustomizedSteppersProps) {
  const [activeStep, setActiveStep] = React.useState(0);
  const { stepContent, stepIcons, steps } = React.useMemo(
    () => ({
      stepContent: registerConnectionContent,
      stepIcons: registerConnectionIcons,
      steps: registerConnectionSteps,
    }),
    [],
  );

  React.useEffect(() => {
    setSharedData({
      metadata: connectionData.metadata,
      capabilities: connectionData.capabilities,
      kind: connectionData.kind,
      resourceID: connectionData.resourceID,
    });
  }, [connectionData, setSharedData]);

  React.useEffect(() => {
    setSharedData((prevState) =>
      prevState.onClose === onClose
        ? prevState
        : {
            ...prevState,
            onClose,
          },
    );
  }, [onClose, setSharedData]);

  const ActiveStepContent = stepContent[String(activeStep + 1)].component;

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const propHandlerMap = {
    handleNext,
    sharedData,
    setSharedData,
    handleRegistrationComplete,
  };

  const stepProps = React.useMemo(
    () =>
      stepContent[String(activeStep + 1)]?.props?.reduce<Record<string, unknown>>(
        (props, propName) => {
          props[propName] = propHandlerMap[propName];
          return props;
        },
        {},
      ) ?? {},
    [activeStep, handleRegistrationComplete, sharedData, stepContent],
  );

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
