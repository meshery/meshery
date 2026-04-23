import React from 'react';
// @ts-ignore
import { useStepper, CustomizedStepper, Box, styled } from '@sistent/sistent';

/**
 * Step interface for the Stepper component.
 */
export interface Step {
  label: string;
  component: React.ReactNode;
  icon?: any;
  helpText?: string;
}

interface StepperProps {
  steps: Step[];
  initialStep?: number;
  contentWrapper?: React.FC<{ children: React.ReactNode }>;
}

const DefaultStepWrapper = styled(Box)({
  padding: 0,
  overflowY: 'auto',
});

/**
 * Reusable Stepper component built on top of Sistent primitives.
 */
const Stepper: React.FC<StepperProps> = ({
  steps,
  initialStep = 0,
  contentWrapper = DefaultStepWrapper,
}) => {
  const stepper = useStepper({
    steps,
    initialStep,
  });

  return (
    <Box>
      <CustomizedStepper {...stepper} ContentWrapper={contentWrapper}>
        {stepper.activeStepComponent}
      </CustomizedStepper>
    </Box>
  );
};

export default Stepper;
