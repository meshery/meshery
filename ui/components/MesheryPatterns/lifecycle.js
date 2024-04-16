import React, { useEffect, useState } from 'react';

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
import { SelectDeploymentTarget } from '../ConfirmationModal';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import { DeploymentSelectorIcon } from '@/assets/icons/DeploymentSelectorIcon';
import { DryRunDesign } from './DryRunDesign';
import { updateProgress } from 'lib/store';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useDeployPatternMutation, useUndeployPatternMutation } from '@/rtk-query/design';

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
  setDryRunErrors,
  selectedK8sContexts,
  pattern_file,
  pattern_id,
  dryRunType,
}) => {
  const handleErrors = (errors) => {
    console.log('handle errors', errors);
    setDryRunErrors(errors);
  };
  const bypassValidation = bypassDryRun || false;
  const toggleBypassValidation = () => setBypassDryRun(!bypassDryRun);

  return (
    <Box>
      <DryRunDesign
        handleClose={handleClose}
        handleErrors={handleErrors}
        selectedContexts={selectedK8sContexts}
        pattern_file={pattern_file}
        pattern_id={pattern_id}
        dryRunType={dryRunType}
      />
      {dryRunErrors?.length > 0 && (
        <Stack spacing={2} direction="row" alignItems="center">
          <Checkbox value={bypassValidation} onChange={toggleBypassValidation} />
          <Typography variant="body2">Bypass errors and initiate deployment</Typography>
        </Stack>
      )}
    </Box>
  );
};

const DEPLOYMENT_TYPE = {
  DEPLOY: 'deploy',
  UNDEPLOY: 'undeploy',
};

// Component to handle the deployment process ie. deploy, undeploy and update ( as all have similar steps)
export const UpdateDeploymentStepper = ({
  handleClose,
  handleComplete,
  action,
  pattern_file,
  pattern_id,
  name,
  selectedK8sContexts,
}) => {
  const [dryRunErrors, setDryRunErrors] = useState([]);
  const [bypassDryRun, setBypassDryRun] = useState(false);
  const [deployPatternMutation] = useDeployPatternMutation();
  const [undeployPatternMutation] = useUndeployPatternMutation();
  const { notify } = useNotification();

  useEffect(() => {
    console.log('Mounting UpdateDeploymentStepper');
  }, []);
  const steps = [
    {
      component: (props) => <SelectTargetStep handleClose={handleClose} {...props} />,
      icon: DeploymentSelectorIcon,
      label: 'Select Target',
    },
    {
      component: (props) => (
        <DryRunStep
          handleClose={handleClose}
          selectedK8sContexts={selectedK8sContexts}
          pattern_file={pattern_file}
          pattern_id={pattern_id}
          setBypassDryRun={setBypassDryRun}
          dryRunErrors={dryRunErrors}
          setDryRunErrors={setDryRunErrors}
          dryRunType={action}
          {...props}
        />
      ),
      label: 'Dry Run',
      icon: DryRunIcon,
    },
  ];

  const deployStepper = useStepper({
    steps,
  });

  const handleDeploy = async (pattern_file, pattern_id, name) => {
    updateProgress({ showProgress: true });
    const result = await deployPatternMutation({
      pattern_file,
      pattern_id,
      selectedK8sContexts,
    });

    updateProgress({ showProgress: false });
    result.data &&
      notify({
        message: `"${name}" Design Deployed`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    result.error &&
      notify({
        message: `"${name}" Design Failed To Deploy`,
        event_type: EVENT_TYPES.ERROR,
      });
  };

  const handleUndeploy = async (pattern_file, pattern_id, name) => {
    updateProgress({ showProgress: true });
    const result = await undeployPatternMutation({
      pattern_file,
      pattern_id,
      selectedK8sContexts,
    });

    updateProgress({ showProgress: false });
    result.data &&
      notify({
        message: `"${name}" Design Undeployed`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    result.error &&
      notify({
        message: `"${name}" Design Failed To Undeploy`,
        event_type: EVENT_TYPES.ERROR,
      });
  };

  const actionFunction = (sharedData) => {
    const command = {
      [DEPLOYMENT_TYPE.DEPLOY]: handleDeploy,
      [DEPLOYMENT_TYPE.UNDEPLOY]: handleUndeploy,
    };

    command[action]?.(pattern_file, pattern_id, name);
    handleClose?.();
    handleComplete?.(sharedData);
    // designMachineRef.current.send(command[action]);
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
      1: () => selectedK8sContexts?.length > 0,
      2: () => dryRunErrors?.length === 0 || bypassDryRun,
    };
    return CanTransition[activeStep + 1](); // can transition to next step
  };

  const canGoNext = checkCanGoNext(deployStepper);

  const nextButtonText = deployStepper.canGoForward ? 'Next' : action;

  return (
    <>
      <ModalBody>
        <Box style={{ width: '30rem' }}>
          <CustomizedStepper {...deployStepper}>
            <deployStepper.activeStepComponent {...deployStepper} />
          </CustomizedStepper>
        </Box>
      </ModalBody>
      <ModalFooter variant="filled" helpText={`${action} the current design`}>
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
  handleComplete,
  pattern_file,
  pattern_id,
  name,
  selectedK8sContexts,
}) => (
  <UpdateDeploymentStepper
    handleClose={handleClose}
    handleComplete={handleComplete}
    pattern_file={pattern_file}
    pattern_id={pattern_id}
    name={name}
    selectedK8sContexts={selectedK8sContexts}
    action={DEPLOYMENT_TYPE.DEPLOY}
  />
);

export const UnDeployStepper = ({
  handleClose,
  handleComplete,
  pattern_file,
  pattern_id,
  name,
  selectedK8sContexts,
}) => (
  <UpdateDeploymentStepper
    handleClose={handleClose}
    handleComplete={handleComplete}
    pattern_file={pattern_file}
    pattern_id={pattern_id}
    name={name}
    selectedK8sContexts={selectedK8sContexts}
    action={DEPLOYMENT_TYPE.UNDEPLOY}
  />
);
