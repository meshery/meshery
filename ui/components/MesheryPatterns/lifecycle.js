// import React, { useContext, useEffect, useState } from 'react';
// import CheckIcon from '../../../assets/CheckIcon';
// import DryRunIcon from '../../../assets/DryRunIcon';
// import { DeploymentSelectorIcon } from '../../../assets/DeploymentSelectorIcon';
//
// import {
//   ModalBody,
//   ModalFooter,
//   useStepper,
//   CustomizedStepper,
//   ModalButtonPrimary,
//   ModalButtonSecondary,
//   Box,
//   Stack,
//   Checkbox,
//   Typography,
// } from '@layer5/sistent';
//
// export const ValidateContent = {
//   btnText: 'Next',
//   cancel: true,
// };
//
// const ValidateDesignStep = ({ handleClose }) => {
//   return <ValidateDesign handleClose={handleClose} />;
// };
//
// const SelectTargetStep = ({ setSharedData, sharedData, ...props }) => {
//   const SelectDepylomentTarget = injectedReactComponents.get(
//     GLOBAL_COMPONENTS.SelectDeploymentTarget,
//   );
//
//   return (
//     <Box>
//       <SelectDepylomentTarget />
//     </Box>
//   );
// };
//
// const DryRunStep = ({ handleClose, bypassDryRun, setBypassDryRun, setDryRunErrors }) => {
//   console.log('Rendering DryRunStep', bypassDryRun, setBypassDryRun, setDryRunErrors);
//
//   useEffect(() => {
//     console.log('Mounting DryRunStep');
//   }, []);
//   const handleErrors = (errors) => {
//     setDryRunErrors(errors);
//   };
//   const bypassValidation = bypassDryRun || false;
//   const toggleBypassValidation = () => setBypassDryRun(!bypassDryRun);
//
//   return (
//     <Box>
//       <DryRunDesign handleClose={handleClose} handleErrors={handleErrors} />
//       <Stack spacing={2} direction="row" alignItems="center">
//         <Checkbox value={bypassValidation} onChange={toggleBypassValidation} />
//         <Typography variant="body2">Bypass errors and initiate deployment</Typography>
//       </Stack>
//     </Box>
//   );
// };
//
// const DEPLOYMENT_TYPE = {
//   DEPLOY: 'deploy',
//   UNDEPLOY: 'undeploy',
// };
//
// // Component to handle the deployment process ie. deploy, undeploy and update ( as all have similar steps)
// export const UpdateDeploymentStepper = ({ handleClose, handleComplete, action }) => {
//   const [dryRunErrors, setDryRunErrors] = useState([]);
//   const [bypassDryRun, setBypassDryRun] = useState(false);
//
//   useEffect(() => {
//     console.log('Mounting UpdateDeploymentStepper');
//   }, []);
//
//   const deployStepper = useStepper({
//     steps: [
//       {
//         component: (props) => <ValidateDesignStep handleClose={handleClose} {...props} />,
//         icon: CheckIcon,
//         label: 'Validate',
//       },
//       {
//         component: (props) => <SelectTargetStep handleClose={handleClose} {...props} />,
//         icon: DeploymentSelectorIcon,
//         label: 'Deployment Selector',
//       },
//       {
//         component: (props) => (
//           <DryRunStep
//             handleClose={handleClose}
//             setBypassDryRun={setBypassDryRun}
//             setDryRunErrors={setDryRunErrors}
//             {...props}
//           />
//         ),
//         label: 'Dry Run',
//         icon: DryRunIcon,
//       },
//     ],
//   });
//
//   const { designMachineRef, selectedK8sContexts } = useContext(DesignerContext);
//
//   const actionFunction = (sharedData) => {
//     const command = {
//       [DEPLOYMENT_TYPE.DEPLOY]: designMachineCommands.deployDesign({
//         context: selectedK8sContexts,
//       }),
//       [DEPLOYMENT_TYPE.UNDEPLOY]: designMachineCommands.undeployDesign({
//         context: selectedK8sContexts,
//       }),
//     };
//
//     handleClose?.();
//     handleComplete?.(sharedData);
//     designMachineRef.current.send(command[action]);
//   };
//
//   const handleNext = () => {
//     console.log('deployStepper', deployStepper, deployStepper.canGoForward);
//     if (deployStepper.canGoForward) {
//       deployStepper.handleNext();
//     } else {
//       actionFunction(deployStepper.sharedData);
//     }
//   };
//
//   const checkCanGoNext = ({ activeStep }) => {
//     // transition map indicating if the next step (key) can be transitioned to
//     // start from 1 because 0 is the first step
//     const CanTransition = {
//       1: () => true,
//       2: () => selectedK8sContexts?.length > 0,
//       3: () => dryRunErrors?.length === 0 || bypassDryRun,
//     };
//     return CanTransition[activeStep + 1](); // can transition to next step
//   };
//
//   const canGoNext = checkCanGoNext(deployStepper);
//
//   const nextButtonText = deployStepper.canGoForward ? 'Next' : action;
//
//   return (
//     <>
//       <ModalBody>
//         <Box style={{ width: '30rem' }}>
//           <CustomizedStepper {...deployStepper} />
//         </Box>
//       </ModalBody>
//       <ModalFooter variant="filled" helpText={`${action} the current design`}>
//         <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'end' }}>
//           <ModalButtonSecondary onClick={deployStepper.goBack} disabled={!deployStepper.canGoBack}>
//             Back
//           </ModalButtonSecondary>
//           <ModalButtonPrimary disabled={!canGoNext} onClick={handleNext}>
//             {nextButtonText}
//           </ModalButtonPrimary>
//         </Box>
//       </ModalFooter>
//     </>
//   );
// };
//
// export const DeployStepper = ({ handleClose, handleComplete }) => (
//   <UpdateDeploymentStepper
//     handleClose={handleClose}
//     handleComplete={handleComplete}
//     action={DEPLOYMENT_TYPE.DEPLOY}
//   />
// );
//
// export const UnDeployStepper = ({ handleClose, handleComplete }) => (
//   <UpdateDeploymentStepper
//     handleClose={handleClose}
//     handleComplete={handleComplete}
//     action={DEPLOYMENT_TYPE.UNDEPLOY}
//   />
// );
