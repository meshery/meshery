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
import { buildConnectionWizardKindConfigs } from './ConnectionWizard.helpers';
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

// The stepper renders each step's `icon` component bare (no size props), so it
// relies on the icon's own defaults — `CheckIcon`/`DescriptionIcon` default to
// 24px, but `ConnectionIcon` has no size default and would collapse to zero.
// Wrap it with an explicit size so the first step's icon actually shows.
const StepConnectionIcon = (props: { width?: number; height?: number }) => (
  <ConnectionIcon width={24} height={24} {...props} />
);

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
      <ModalFooter
        variant="filled"
        helpText={
          wizard.ctx.data.kindConfig
            ? `Meshery connections are first-class constructs. This wizard registers a ${wizard.ctx.data.kindConfig.label} connection.`
            : 'Choose a supported connection kind to continue.'
        }
      >
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
