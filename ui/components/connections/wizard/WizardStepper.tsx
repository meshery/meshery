import { ReactNode, useEffect } from 'react';
import { CustomizedStepper, useStepper } from '@sistent/sistent';

/**
 * A single step descriptor for Sistent's stepper header. `component` is unused
 * by our modals (we render the active body ourselves), but Sistent's
 * `useStepper` reads `steps[activeStep].component`, so it must be present.
 */
export type WizardStepperStep = {
  label: string;
  icon: ReactNode;
  component: ReactNode;
};

type WizardStepperProps = {
  steps: WizardStepperStep[];
  /** The engine-owned active step. Always clamped to `steps` by the caller. */
  activeIndex: number;
  children?: ReactNode;
};

/**
 * Sistent's `useStepper` keeps its own `activeStep` state and reads
 * `steps[activeStep].component` on every render. We drive that state from the
 * wizard engine via a post-render effect, but the `steps` array is rebuilt
 * synchronously during render — so any render where the step list shrinks (a
 * step's `hidden` predicate flips, the selected kind changes, or the wizard
 * resets) lands before the effect lowers the internal index, leaving it past
 * the end of the array: `steps[activeStep]` is `undefined` and the read throws
 * ("can't access property 'component', ... is undefined").
 *
 * Remounting whenever the step count changes resets the internal index to a
 * always-valid `0`; the effect then re-applies the real index. This keeps the
 * internal index and the `steps` array length in lockstep across every render.
 */
const WizardStepperInner = ({ steps, activeIndex, children }: WizardStepperProps) => {
  const stepper = useStepper({ steps });

  // `stepper.setActiveStep` is stable; re-run only when the target changes.
  useEffect(() => {
    stepper.setActiveStep(activeIndex);
  }, [activeIndex]);

  return <CustomizedStepper {...stepper}>{children}</CustomizedStepper>;
};

const WizardStepper = (props: WizardStepperProps) => (
  <WizardStepperInner key={props.steps.length} {...props} />
);

export default WizardStepper;
