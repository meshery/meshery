import { useEffect, useMemo } from 'react';
import {
  Box,
  DescriptionIcon,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  CheckIcon,
} from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { Modal } from '@/components/shared/Modal';
import ConnectionIcon from '@/assets/icons/Connection';
import type { RootState } from '../../store';
import { useListConnectionDefinitionsQuery } from '@meshery/schemas/mesheryApi';
import {
  buildConnectionWizardKindConfigs,
  DEFAULT_CONNECTION_DOCS_URL,
} from './ConnectionWizard.helpers';
import { useConnectionWizard } from './wizard/useConnectionWizard';
import WizardStepper from './wizard/WizardStepper';

type ConnectionWizardModalProps = {
  isOpen: boolean;
  onClose: () => void;
  /** Pre-select a connection kind once definitions load (e.g. "kubernetes"). */
  presetKind?: string | null;
  /** When true with presetKind, land on the step after "Choose Connection". */
  skipKindSelection?: boolean;
};

// Stepper passes `fill` from theme; force monochrome so ConnectionIcon is not
// brand-green before the step completes (issue #20767).
const StepConnectionIcon = (props: { width?: number; height?: number; fill?: string }) => {
  const { fill, width = 24, height = 24, ...rest } = props;
  return (
    <ConnectionIcon
      width={width}
      height={height}
      {...(fill ? { fill, primaryFill: fill, secondaryFill: fill } : {})}
      {...rest}
    />
  );
};

const ConnectionWizardModal = ({
  isOpen,
  onClose,
  presetKind = null,
  skipKindSelection = false,
}: ConnectionWizardModalProps) => {
  const { connectionMetadataState } = useSelector((state: RootState) => state.ui);
  const { data: connectionDefinitionsResponse, isFetching: isLoadingKinds } =
    useListConnectionDefinitionsQuery({}, { skip: !isOpen });

  // The endpoint returns the page under `connectionDefinitions`; the wizard
  // builds its selectable kinds from that list instead of a hardcoded set.
  const kindConfigs = useMemo(
    () => buildConnectionWizardKindConfigs(connectionDefinitionsResponse?.connectionDefinitions),
    [connectionDefinitionsResponse?.connectionDefinitions],
  );

  const wizard = useConnectionWizard({
    mode: 'create',
    isOpen,
    availableKinds: kindConfigs,
    isLoadingKinds,
    connectionIconMap: connectionMetadataState || undefined,
    presetKind,
    skipKindSelection,
    onComplete: onClose,
  });

  // Reset wizard state whenever the modal closes.
  useEffect(() => {
    if (!isOpen) {
      wizard.reset();
    }
    // wizard.reset is stable; avoid re-running on every render.
  }, [isOpen]);

  const steps = useMemo(
    () =>
      wizard.steps.map((step, index) => ({
        label: step.label,
        icon:
          step.icon ||
          (index === 0
            ? StepConnectionIcon
            : index === wizard.steps.length - 1
              ? CheckIcon
              : DescriptionIcon),
        component: <></>,
      })),
    [wizard.steps],
  );

  const ActiveBody = wizard.activeStep?.Component;

  // Step-specific footer help with docs links (#20767); fall back to kind copy.
  const helpText = useMemo(() => {
    const stepHelp = wizard.activeStep?.helpText;
    if (typeof stepHelp === 'function') {
      return stepHelp(wizard.ctx);
    }
    if (typeof stepHelp === 'string' && stepHelp.length > 0) {
      return stepHelp;
    }
    const kindConfig = wizard.ctx.data.kindConfig;
    if (kindConfig) {
      const docsUrl = kindConfig.docsUrl || DEFAULT_CONNECTION_DOCS_URL;
      return `Meshery connections are first-class constructs. This wizard registers a ${kindConfig.label} connection. [Learn more](${docsUrl}).`;
    }
    return `Choose a supported connection kind to continue. [Learn more about connections](${DEFAULT_CONNECTION_DOCS_URL}).`;
  }, [wizard.activeStep, wizard.ctx]);

  const handleClose = () => {
    if (wizard.isBusy) {
      return;
    }
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Create Connection" size="lg">
      <ModalBody>
        <WizardStepper steps={steps} activeIndex={wizard.activeIndex}>
          {ActiveBody ? <ActiveBody ctx={wizard.ctx} /> : <></>}
        </WizardStepper>
      </ModalBody>
      <ModalFooter variant="filled" helpText={helpText}>
        <Box sx={{ width: '100%', display: 'flex', justifyContent: 'end', gap: 2 }}>
          <ModalButtonSecondary
            onClick={wizard.canGoBack ? wizard.back : handleClose}
            disabled={wizard.isBusy}
          >
            {wizard.canGoBack ? 'Back' : 'Cancel'}
          </ModalButtonSecondary>
          <ModalButtonPrimary onClick={wizard.next} disabled={!wizard.canProceed || wizard.isBusy}>
            {wizard.isBusy ? 'Working...' : wizard.nextLabel}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </Modal>
  );
};

export default ConnectionWizardModal;
