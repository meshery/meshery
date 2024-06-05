import React, { useState } from 'react';
import { ValidateDesign } from './ValidateDesign';
import { DryRunDesign, getTotalCountOfDeploymentErrors } from './DryRun';
import {
  ModalBody,
  ModalFooter,
  useStepper,
  CustomizedStepper,
  ModalButtonPrimary,
  ModalButtonSecondary,
  Box,
  Stack,
  Checkbox,
  Typography,
} from '@layer5/sistent';
import { DEPLOYMENT_TYPE } from './common';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import { SelectDeploymentTarget } from '../ConfirmationModal';
import { DeploymentSelectorIcon } from '@/assets/icons/DeploymentSelectorIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { useFilterK8sContexts } from '../hooks/useKubernetesHook';
import { selectK8sContexts } from 'lib/store';
import { useSelector } from 'react-redux';

export const ValidateContent = {
  btnText: 'Next',
  cancel: true,
};

const SelectTargetStep = () => {
  return (
    <Box>
      <SelectDeploymentTarget />
    </Box>
  );
};

const DryRunStep = ({
  handleClose,
  bypassDryRun,
  setBypassDryRun,
  dryRunErrors,
  deployment_type,
  selectedK8sContexts,
  design,
  validationMachine,
}) => {
  const bypassValidation = bypassDryRun || false;
  const toggleBypassValidation = () => setBypassDryRun(!bypassDryRun);

  return (
    <Box>
      <DryRunDesign
        handleClose={handleClose}
        selectedK8sContexts={selectedK8sContexts}
        design={design}
        validationMachine={validationMachine}
        deployment_type={deployment_type}
      />

      {getTotalCountOfDeploymentErrors(dryRunErrors) > 0 && (
        <Stack spacing={2} direction="row" alignItems="center">
          <Checkbox value={bypassValidation} onChange={toggleBypassValidation} />
          <Typography variant="body2">Bypass errors and initiate deployment</Typography>
        </Stack>
      )}
    </Box>
  );
};

// Component to handle the deployment process ie. deploy, undeploy and update ( as all have similar steps)
export const UpdateDeploymentStepper = ({
  handleClose,
  deployment_type,
  handlePerformDeployment,
  selectedK8sContexts,
  design,
  validationMachine,
}) => {
  const [dryRunErrors, setDryRunErrors] = useState([]);
  const [bypassDryRun, setBypassDryRun] = useState(false);
  const k8sContext = useSelector(selectK8sContexts);
  const selectedDeployableK8scontexts = useFilterK8sContexts(
    k8sContext,
    ({ context, operatorState }) => {
      const isSelected =
        selectedK8sContexts?.includes(context.id) || selectedK8sContexts?.includes('all');
      return isSelected && operatorState !== 'DISABLED';
    },
  );

  const deployStepper = useStepper({
    steps: [
      {
        component: (
          <ValidateDesign
            handleClose={handleClose}
            validationMachine={validationMachine}
            design={design}
          />
        ),
        icon: CheckIcon,
        helpText: `Validate the design before deploying, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the validation process`,
        label: 'Validate Design',
      },
      {
        component: <SelectTargetStep handleClose={handleClose} />,
        helpText:
          'Select the environment to deploy the design , only the kubernetes clusters with the operator enabled are shown,[learn more](https://docs.meshery.io/guides/infrastructure-management/overview)  about the environment selection',
        icon: DeploymentSelectorIcon,

        label: 'Select Environment',
      },
      {
        component: (
          <DryRunStep
            selectedK8sContexts={selectedK8sContexts}
            design={design}
            validationMachine={validationMachine}
            handleClose={handleClose}
            deployment_type={deployment_type}
            setBypassDryRun={setBypassDryRun}
            setDryRunErrors={setDryRunErrors}
            dryRunErrors={dryRunErrors}
          />
        ),
        helpText:
          'Dry Run is a simulation of the deployment process, it helps to identify potential errors before the actual deployment , [learn more](https://docs.meshery.io/guides/infrastructure-management/overview ) about Dry Run.',
        label: 'Dry Run',
        icon: DryRunIcon,
      },
    ],
  });

  const actionFunction = (sharedData) => {
    console.log('sharedData', sharedData);
    handleClose?.();
    handlePerformDeployment({
      design,
      selectedK8sContexts: selectedDeployableK8scontexts.map((context) => context.id),
    });
  };

  const handleNext = () => {
    console.log('deployStepper', deployStepper, deployStepper.canGoForward);
    if (deployStepper.canGoForward) {
      deployStepper.handleNext();
    } else {
      actionFunction(deployStepper.sharedData);
    }
  };

  const checkCanGoNext = ({ activeStep }) => {
    // transition map indicating if the next step (key) can be transitioned to
    // start from 1 because 0 is the first step
    const CanTransition = {
      1: () => true,
      2: () => selectedK8sContexts?.length > 0,
      3: () => getTotalCountOfDeploymentErrors(dryRunErrors) == 0 || bypassDryRun,
    };
    return CanTransition[activeStep + 1](); // can transition to next step
  };

  const canGoNext = checkCanGoNext(deployStepper);

  const nextButtonText = deployStepper.canGoForward ? 'Next' : deployment_type;

  return (
    <>
      <ModalBody>
        <CustomizedStepper {...deployStepper}>
          <Box style={{ overflowY: 'auto' }}>{deployStepper.activeStepComponent}</Box>
        </CustomizedStepper>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={
          deployStepper.steps[deployStepper.activeStep]?.helpText ||
          `${deployment_type} the current design`
        }
      >
        <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'end' }}>
          <ModalButtonSecondary onClick={deployStepper.goBack} disabled={!deployStepper.canGoBack}>
            Back
          </ModalButtonSecondary>
          <ModalButtonPrimary disabled={!canGoNext} onClick={handleNext}>
            {nextButtonText}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </>
  );
};

export const DeployStepper = ({
  handleClose,
  design,
  validationMachine,
  handleDeploy,
  selectedK8sContexts,
}) => (
  <UpdateDeploymentStepper
    selectedK8sContexts={selectedK8sContexts}
    handleClose={handleClose}
    design={design}
    handlePerformDeployment={handleDeploy}
    validationMachine={validationMachine}
    deployment_type={DEPLOYMENT_TYPE.DEPLOY}
  />
);

export const UnDeployStepper = ({
  handleClose,
  design,
  handleUndeploy,
  validationMachine,
  selectedK8sContexts,
}) => (
  <UpdateDeploymentStepper
    selectedK8sContexts={selectedK8sContexts}
    handleClose={handleClose}
    design={design}
    handlePerformDeployment={handleUndeploy}
    validationMachine={validationMachine}
    deployment_type={DEPLOYMENT_TYPE.UNDEPLOY}
  />
);
