import type { ConnectionWizardKindConfig } from '../ConnectionWizard.helpers';
import {
  genericCredentialStep,
  genericDetailsStep,
  genericReceiptStep,
  genericRegisterStep,
  selectStep,
} from './genericSteps';
import { kubernetesExtension } from './kubernetesExtension';
import type { ConnectionExtension, WizardMode, WizardStep } from './types';

/**
 * Registered connection extensions. A connection matches the most specific
 * entry: a `kind` + `type` + `subType` match wins over `kind` alone.
 */
const EXTENSIONS: ConnectionExtension[] = [kubernetesExtension];

const matchScore = (extension: ConnectionExtension, config: ConnectionWizardKindConfig): number => {
  const { match } = extension;
  if (match.kind !== config.kind) {
    return -1;
  }
  if (match.type !== undefined && match.type !== config.type) {
    return -1;
  }
  if (match.subType !== undefined && match.subType !== config.subType) {
    return -1;
  }
  // More specified fields => more specific match.
  return (match.type !== undefined ? 1 : 0) + (match.subType !== undefined ? 1 : 0);
};

export const resolveExtension = (
  config?: ConnectionWizardKindConfig | null,
): ConnectionExtension | null => {
  if (!config) {
    return null;
  }
  let best: ConnectionExtension | null = null;
  let bestScore = -1;
  for (const extension of EXTENSIONS) {
    const score = matchScore(extension, config);
    if (score > bestScore) {
      best = extension;
      bestScore = score;
    }
  }
  return bestScore >= 0 ? best : null;
};

/**
 * Composes the ordered step list for a connection. In `configure` mode only the
 * post-config + receipt steps run (for an already-registered connection); in
 * `create` mode the full select -> details -> credential -> register ->
 * post-config -> receipt sequence is built, falling back to generic defaults.
 */
export const buildSteps = (
  config: ConnectionWizardKindConfig | null,
  mode: WizardMode,
): WizardStep[] => {
  const extension = resolveExtension(config);
  const postConfigSteps = extension?.postConfigSteps ?? [];
  const receiptStep = extension?.receiptStep ?? genericReceiptStep;

  if (mode === 'configure') {
    return [...postConfigSteps, receiptStep];
  }

  const steps: WizardStep[] = [selectStep];
  if (!config) {
    return steps;
  }

  steps.push(extension?.detailsStep ?? genericDetailsStep);

  // `credentialStep: null` removes the step; `undefined` falls back to generic.
  const credentialStep =
    extension && 'credentialStep' in extension ? extension.credentialStep : genericCredentialStep;
  if (credentialStep) {
    steps.push(credentialStep);
  }

  steps.push(extension?.registerStep ?? genericRegisterStep);
  steps.push(...postConfigSteps);
  steps.push(receiptStep);

  return steps;
};
