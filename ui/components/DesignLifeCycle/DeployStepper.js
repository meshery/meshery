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
  Typography,
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
import { EnvironmentIcon } from '@layer5/sistent';
import { useContext } from 'react';
import { NotificationCenterContext } from '../NotificationCenter';
import { useEffect } from 'react';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import { FormatStructuredData, TextWithLinks } from '../DataFormatter';
import { SEVERITY_STYLE } from '../NotificationCenter/constants';
import { Stack } from '@layer5/sistent';
import { capitalize } from 'lodash';
import { ErrorMetadataFormatter } from '../NotificationCenter/metadata';
import FinishFlagIcon from '@/assets/icons/FinishFlagIcon';
import { alpha } from '@mui/material';

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

// const sampleErrorEvent = {
//   id: 'f2c33cbd-8213-4d12-93e9-a267189f7551',
//   userID: 'ad6f28ae-dd45-4ebb-bf1e-97351fb00fa6',
//   actedUpon: '92ccd93c-747a-49e9-b8c0-bbfd59b58059',
//   operationID: 'efc0deac-802f-4779-8448-c969d77da4fc',
//   systemID: '3a2848fa-8a4f-494c-9391-48e368eebe44',
//   status: 'unread',
//   severity: 'error',
//   action: 'deploy',
//   category: 'pattern',
//   description: "deploy error for design 'nginx'",
//   metadata: {
//     error: {
//       Code: 'meshery-server-1073',
//       Severity: 2,
//       ShortDescription: [
//         'unable to Create Comp Config.',
//         'an empty namespace may not be set during creation',
//       ],
//       LongDescription: ['an empty namespace may not be set during creation'],
//       ProbableCause: [],
//       SuggestedRemediation: [],
//     },
//   },
//   createdAt: '2024-06-05T17:50:45.897572829+05:30',
//   updatedAt: '2024-06-05T17:50:45.89759321+05:30',
//   deletedAt: null,
//   user_id: 'ad6f28ae-dd45-4ebb-bf1e-97351fb00fa6',
//   system_id: '3a2848fa-8a4f-494c-9391-48e368eebe44',
//   updated_at: '2024-06-05T17:50:45.89759321+05:30',
//   created_at: '2024-06-05T17:50:45.897572829+05:30',
//   deleted_at: null,
//   operation_id: 'efc0deac-802f-4779-8448-c969d77da4fc',
// };
const StyledDetailBox = styled(Box)(({ theme, severityColor, bgOpacity }) => ({
  padding: theme.spacing(2),
  backgroundColor: alpha(severityColor, bgOpacity),
  borderLeft: `4px solid ${severityColor}`,
  borderBottom: `1px solid ${theme.palette.divider}`,
  display: 'flex',
}));

export const FinishDeploymentStep = ({ perform_deployment, deployment_type }) => {
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const theme = useTheme();

  const [isDeploying, setIsDeploying] = useState(false);
  const [deployEvent, setDeployEvent] = useState();
  const [deployError, setDeployError] = useState(null);

  useEffect(() => {
    try {
      setIsDeploying(true);
      perform_deployment();
    } catch (error) {
      setDeployError(error);
      setIsDeploying(false);
    }
  }, []);

  useEffect(() => {
    const subscription = operationsCenterActorRef.on(
      OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      (event) => {
        const serverEvent = event.data.event;
        // console.log('serverEvent', serverEvent);
        if (serverEvent.action === deployment_type) {
          setIsDeploying(false);
          setDeployEvent(serverEvent);
        }
      },
    );

    return () => subscription.unsubscribe();
  }, []);

  const progressMessage = `${capitalize(deployment_type)}ing  design`;

  if (isDeploying) {
    return <Loading message={progressMessage} />;
  }

  if (deployError) {
    return (
      <Typography variant="h5" color="error">
        Error deploying design: {JSON.stringify(deployError)}
      </Typography>
    );
  }

  const eventStyle = SEVERITY_STYLE[deployEvent.severity] || {};

  const errors = deployEvent.metadata?.error;
  // const message = deployEvent.metadata?.summary?.message;
  // const hasDetails = !_.isEmpty(details);
  // console.log('deployEvent', eventStyle, deployEvent);

  return (
    <Stack spacing={2}>
      <Typography variant="textB2SemiBold" color={theme.palette.text.secondary}>
        {`${deployment_type}ment Summary`}
      </Typography>

      <Box borderRadius={'0.5rem'} border={`2px solid ${eventStyle.color}`}>
        <StyledDetailBox severityColor={eventStyle.color} bgOpacity={0.1}>
          <TextWithLinks
            text={deployEvent?.description || ''}
            style={{ whiteSpace: 'pre-wrap', color: theme.palette.text.default }}
          />
        </StyledDetailBox>
        {errors && (
          <StyledDetailBox severityColor={eventStyle.color} bgOpacity={0}>
            <ErrorMetadataFormatter metadata={errors} event={deployEvent} />
          </StyledDetailBox>
        )}
        <StyledDetailBox severityColor={eventStyle.color} bgOpacity={0}>
          <FormatStructuredData data={deployEvent.metadata} />
        </StyledDetailBox>
        {/* {hasDetails && (

        {/* )} */}
      </Box>
    </Stack>
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

  const FinalizeBackgroundColor = theme.palette?.background?.blur.light;
  const actionFunction = () => {
    handlePerformDeployment({
      design,
      selectedK8sContexts: selectedDeployableK8scontexts.map((context) => context.id),
    });
  };
  const deployStepper = useStepper({
    steps: [
      {
        component: (
          <StepContent>
            <ValidateDesign
              handleClose={handleClose}
              validationMachine={validationMachine}
              design={design}
            />
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
        icon: EnvironmentIcon,
        label: 'Identify Environments',
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
              dryRunErrors={dryRunErrors}
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
          <StepContent backgroundColor={FinalizeBackgroundColor}>
            <FinalizeDeployment design={design} deployment_type={deployment_type} />
          </StepContent>
        ),
        helpText:
          'Finalize the deployment process, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the finalization process',
        label: 'Finalize Deployment',
        icon: DeploymentSelectorIcon,
      },
      {
        component: (
          <StepContent>
            <FinishDeploymentStep
              design={design}
              deployment_type={deployment_type}
              perform_deployment={actionFunction}
            />{' '}
          </StepContent>
        ),
        helpText:
          'Finalize the deployment process, [learn more](https://docs.meshery.io/guides/infrastructure-management/overview) about the finalization process',
        label: 'Finsh',
        icon: FinishFlagIcon,
      },
    ],
  });

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
      nextButtonText: capitalize(deployment_type),
      nextAction: () => {
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
