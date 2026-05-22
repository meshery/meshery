import { describe, it } from 'vitest';

// DeployStepper builds a sistent useStepper-driven multi-step modal that
// wires Redux selectors (selectAllSelectedK8sConnections), xstate actor
// hooks (useDryRunValidationResults / useIsValidatingDryRun), several
// child containers (SelectDeploymentTarget, DryRunDesign,
// FinalizeDeployment, FinishDeploymentStep, ValidateDesign), the
// NotificationCenter context, and a router. End-to-end deploy/undeploy
// behavior is covered by playwright. The DEPLOYMENT_TYPE constants and
// individual sub-steps (DryRun summary, ValidateDesign, common helpers)
// are unit-tested separately.
describe.skip('DeployStepper / UpdateDeploymentStepper', () => {
  it('skipped - multi-step stepper bound to xstate + redux + router', () => {});
});
