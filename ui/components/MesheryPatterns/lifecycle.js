import React, { useState } from 'react';

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
import { selectSelectedK8sClusters, updateProgress } from 'lib/store';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useDeployPatternMutation, useUndeployPatternMutation } from '@/rtk-query/design';
import { useSelector } from 'react-redux';

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
    setDryRunErrors(errors);
  };
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
          <Checkbox value={bypassDryRun} onChange={toggleBypassValidation} />
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
}) => {
  const [dryRunErrors, setDryRunErrors] = useState([]);
  const [bypassDryRun, setBypassDryRun] = useState(false);
  const [deployPatternMutation] = useDeployPatternMutation();
  const [undeployPatternMutation] = useUndeployPatternMutation();
  const selectedK8sContexts = useSelector(selectSelectedK8sClusters);

  const { notify } = useNotification();

  const steps = [
    {
      component: <SelectTargetStep handleClose={handleClose} setDryRunErrors={setDryRunErrors} />,
      icon: DeploymentSelectorIcon,
      label: 'Select Environment',
    },
    {
      component: (
        <DryRunStep
          handleClose={handleClose}
          selectedK8sContexts={selectedK8sContexts}
          pattern_file={pattern_file}
          pattern_id={pattern_id}
          setBypassDryRun={setBypassDryRun}
          dryRunErrors={dryRunErrors}
          setDryRunErrors={setDryRunErrors}
          dryRunType={action}
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

  const actionFunction = () => {
    const command = {
      [DEPLOYMENT_TYPE.DEPLOY]: handleDeploy,
      [DEPLOYMENT_TYPE.UNDEPLOY]: handleUndeploy,
    };

    command[action]?.(pattern_file, pattern_id, name);
    handleClose?.();
    handleComplete?.();
    // designMachineRef.current.send(command[action]);
  };

  const handleNext = () => {
    console.log('deployStepper', deployStepper, deployStepper.canGoForward);
    if (deployStepper.canGoForward) {
      deployStepper.handleNext();
    } else {
      actionFunction();
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
    <div>
      <ModalBody>
        <Box style={{ width: '30rem' }}>
          <CustomizedStepper {...deployStepper}>
            {deployStepper.activeStepComponent}
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
    </div>
  );
};

export const DeployStepper = ({ handleClose, handleComplete, pattern_file, pattern_id, name }) => (
  <UpdateDeploymentStepper
    handleClose={handleClose}
    handleComplete={handleComplete}
    pattern_file={pattern_file}
    pattern_id={pattern_id}
    name={name}
    action={DEPLOYMENT_TYPE.DEPLOY}
  />
);

export const UnDeployStepper = ({
  handleClose,
  handleComplete,
  pattern_file,
  pattern_id,
  name,
}) => (
  <UpdateDeploymentStepper
    handleClose={handleClose}
    handleComplete={handleComplete}
    pattern_file={pattern_file}
    pattern_id={pattern_id}
    name={name}
    action={DEPLOYMENT_TYPE.UNDEPLOY}
  />
);
