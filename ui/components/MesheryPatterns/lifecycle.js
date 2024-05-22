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
import { selectK8sContexts, selectSelectedK8sClusters, updateProgress } from 'lib/store';
import { useDeployPatternMutation, useUndeployPatternMutation } from '@/rtk-query/design';
import { useSelector } from 'react-redux';
import { useFilterK8sContexts } from '../hooks/useKubernetesHook';
import { RenderTooltipContent } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';

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
  const k8sContext = useSelector(selectK8sContexts);
  const selectedK8sContexts = useSelector(selectSelectedK8sClusters);

  const selectedDeployableK8scontexts = useFilterK8sContexts(
    k8sContext,
    ({ context, operatorState }) => {
      const isSelected =
        selectedK8sContexts?.includes(context.id) || selectedK8sContexts?.includes('all');
      return isSelected && operatorState !== 'DISABLED';
    },
  );

  const steps = [
    {
      component: <SelectTargetStep handleClose={handleClose} setDryRunErrors={setDryRunErrors} />,
      icon: DeploymentSelectorIcon,
      label: 'Select Environment',
      helpText: RenderTooltipContent({
        showPriortext:
          'Select the environment to deploy the design , only the kubernetes clusters with the operator enabled are shown,',
        link: 'https://docs.meshery.io/guides/infrastructure-management/overview',
        showAftertext: ' about the environment selection.',
      }),
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
      helpText: RenderTooltipContent({
        showPriortext:
          'Dry Run is a simulation of the deployment process, it helps to identify potential errors before the actual deployment,',
        showAftertext: ' to learn more about Dry Run.',
        link: 'https://docs.meshery.io/guides/infrastructure-management/overview',
      }),
      icon: DryRunIcon,
    },
  ];

  const deployStepper = useStepper({
    steps,
  });

  const handleDeploy = async (pattern_file, pattern_id) => {
    updateProgress({ showProgress: true });
    await deployPatternMutation({
      pattern_file,
      pattern_id,
      selectedK8sContexts: selectedDeployableK8scontexts.map((ctx) => ctx.id),
    });
    updateProgress({ showProgress: false });
  };

  const handleUndeploy = async (pattern_file, pattern_id) => {
    updateProgress({ showProgress: true });
    await undeployPatternMutation({
      pattern_file,
      pattern_id,
      selectedK8sContexts: selectedDeployableK8scontexts.map((ctx) => ctx.id),
    });

    updateProgress({ showProgress: false });
  };

  const actionFunction = () => {
    const command = {
      [DEPLOYMENT_TYPE.DEPLOY]: handleDeploy,
      [DEPLOYMENT_TYPE.UNDEPLOY]: handleUndeploy,
    };

    command[action]?.(pattern_file, pattern_id, name);
    handleClose?.();
    handleComplete?.();
  };

  const handleNext = () => {
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
      1: () => selectedDeployableK8scontexts?.length > 0,
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
            <Box style={{ overflowY: 'auto' }}>{deployStepper.activeStepComponent}</Box>
          </CustomizedStepper>
        </Box>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={
          deployStepper.steps[deployStepper.activeStep]?.helpText || `${action} the current design`
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
