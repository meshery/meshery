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
} from '@layer5/sistent';
import { CheckBoxField, DEPLOYMENT_TYPE, Loading } from './common';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import { DeploymentSelectorIcon } from '@/assets/icons/DeploymentSelectorIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { useFilterK8sContexts } from '../hooks/useKubernetesHook';
import { selectK8sContexts, useLegacySelector } from 'lib/store';
import { DeploymentTargetContext, SelectTargetEnvironments } from './SelectDeploymentTarget';
import { FinalizeDeployment } from './finalizeDeployment';
import { selectSelectedEnvs } from '@/store/slices/globalEnvironmentContext';
import {
  useDryRunValidationResults,
  useIsValidatingDryRun,
} from 'machines/validator/designValidator';
import { useSelectorRtk } from '@/store/hooks';
import { styled } from '@layer5/sistent';
import { useTheme } from '@layer5/sistent';
export const ValidateContent = {
  btnText: 'Next',
  cancel: true,
};

const StepWrapper = styled(Box)({
  padding: 0,
  overflowY: 'auto',
});

const StepContent = styled('div', {
  // shouldForwardProp: (prop) => prop !== 'backgroundColor',
})(({ theme, backgroundColor }) => ({
  paddingInline: theme.spacing(4),
  paddingBlock: theme.spacing(2),
  backgroundColor: backgroundColor || theme.palette.background.constant.white,
}));

export const FinishDeploymentStep = () => {
  return (
    <Box>
      <Loading message="Deploying Design" />
    </Box>
  );
};

const SelectTargetStep = () => {
  const organization = useLegacySelector((state) => state.get('organization'));
  const connectionMetadataState = useLegacySelector((state) =>
    state.get('connectionMetadataState'),
  );
  const meshsyncControllerState = useLegacySelector((state) => state.get('controllerState'));
  return (
    <DeploymentTargetContext.Provider
      value={{ connectionMetadataState, meshsyncControllerState, organization }}
    >
      <SelectTargetEnvironments organization={organization} />
    </DeploymentTargetContext.Provider>
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
  includeDependencies,
  setIncludeDependencies,
}) => {
  const bypassValidation = bypassDryRun || false;
  const toggleBypassValidation = () => setBypassDryRun(!bypassDryRun);
  const toggleIncludeDependencies = () => setIncludeDependencies(!includeDependencies);

  return (
    <Box>
      <DryRunDesign
        handleClose={handleClose}
        selectedK8sContexts={selectedK8sContexts}
        includeDependencies={includeDependencies}
        design={design}
        validationMachine={validationMachine}
        deployment_type={deployment_type}
      />

      {getTotalCountOfDeploymentErrors(dryRunErrors) > 0 && (
        <CheckBoxField
          label="Bypass errors and initiate deployment"
          checked={bypassValidation}
          onChange={toggleBypassValidation}
        />
      )}

      <CheckBoxField
        label="Include Dependencies"
        checked={includeDependencies}
        onChange={toggleIncludeDependencies}
      />
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
  const [includeDependencies, setIncludeDependencies] = useState(false);
  const [bypassDryRun, setBypassDryRun] = useState(false);

  const selectedEnvironments = useSelectorRtk(selectSelectedEnvs);
  const selectedEnvCount = Object.keys(selectedEnvironments).length;
  const dryRunErrors = useDryRunValidationResults(validationMachine);
  const totalDryRunErrors = getTotalCountOfDeploymentErrors(dryRunErrors);
  const isDryRunning = useIsValidatingDryRun(validationMachine);
  const theme = useTheme();

  const k8sContext = useLegacySelector(selectK8sContexts);
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
          <StepContent backgroundColor={theme.palette.background.constant.black}>
            {' '}
            <ValidateDesign
              handleClose={handleClose}
              validationMachine={validationMachine}
              design={design}
            />{' '}
          </StepContent>
        ),
        icon: CheckIcon,
        helpText: `Validate the design before deploying, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the validation process`,
        label: 'Validate Design',
      },
      {
        component: (
          <StepContent>
            {' '}
            <SelectTargetStep handleClose={handleClose} />{' '}
          </StepContent>
        ),
        helpText:
          'Select the environment to deploy the design , only the kubernetes clusters with the operator enabled are shown,[learn more](https://docs.meshery.io/guides/infrastructure-management/overview)  about the environment selection',
        icon: DeploymentSelectorIcon,

        label: 'Select Environment',
      },
      {
        component: (
          <StepContent>
            <DryRunStep
              selectedK8sContexts={selectedK8sContexts}
              design={design}
              validationMachine={validationMachine}
              handleClose={handleClose}
              deployment_type={deployment_type}
              includeDependencies={includeDependencies}
              setIncludeDependencies={setIncludeDependencies}
              setBypassDryRun={setBypassDryRun}
            />
          </StepContent>
        ),
        helpText:
          'Dry Run is a simulation of the deployment process, it helps to identify potential errors before the actual deployment , [learn more](https://docs.meshery.io/guides/infrastructure-management/overview ) about Dry Run.',
        label: 'Dry Run',
        icon: DryRunIcon,
      },
      {
        component: (
          <StepContent>
            <FinalizeDeployment design={design} deployment_type={deployment_type} />
          </StepContent>
        ),
        helpText:
          'Finalize the deployment process, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the finalization process',
        label: 'Finalize Deployment',
        icon: DryRunIcon,
      },
      {
        component: (
          <StepContent backgroundColor="red">
            <FinishDeploymentStep design={design} deployment_type={deployment_type} />{' '}
          </StepContent>
        ),
        helpText:
          'Finalize the deployment process, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the finalization process',
        label: 'Finsh',
        icon: DryRunIcon,
      },
    ],
  });

  const actionFunction = () => {
    handlePerformDeployment({
      design,
      selectedK8sContexts: selectedDeployableK8scontexts.map((context) => context.id),
    });
  };

  const transitionConfig = {
    0: {
      canGoNext: () => true,
      nextButtonText: 'Next',
      nextAction: () => deployStepper.handleNext(),
    },
    1: {
      canGoNext: () => selectedEnvCount > 0,
      nextButtonText: 'Next',
      nextAction: () => deployStepper.handleNext(),
    },
    2: {
      canGoNext: () => !isDryRunning & (totalDryRunErrors == 0 || bypassDryRun),
      nextButtonText: 'Next',
      nextAction: () => deployStepper.handleNext(),
    },
    3: {
      canGoNext: () => true,
      nextButtonText: 'Deploy',
      nextAction: () => {
        actionFunction();
        deployStepper.handleNext();
      },
    },
    4: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleClose,
    },
  };

  const canGoNext = transitionConfig[deployStepper.activeStep].canGoNext();
  const nextButtonText = transitionConfig[deployStepper.activeStep].nextButtonText;

  return (
    <>
      <ModalBody style={{ padding: 0 }}>
        <Box style={{ width: '80vw', maxWidth: '40rem' }}>
          <CustomizedStepper {...deployStepper} ContentWrapper={StepWrapper}>
            {deployStepper.activeStepComponent}
          </CustomizedStepper>
        </Box>
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
          <ModalButtonPrimary
            disabled={!canGoNext}
            onClick={transitionConfig[deployStepper.activeStep].nextAction}
          >
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
