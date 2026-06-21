import { useEffect, useMemo } from 'react';
import {
  Box,
  DescriptionIcon,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
} from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { Modal } from '@/components/shared/Modal';
import ConnectionIcon from '@/assets/icons/Connection';
import CheckIcon from '@/assets/icons/CheckIcon';
import type { RootState } from '@/store/store';
import { useListConnectionDefinitionsQuery } from '@meshery/schemas/mesheryApi';
import { buildConnectionWizardKindConfigs } from './ConnectionWizard.helpers';
import { useConnectionWizard } from './wizard/useConnectionWizard';
import WizardStepper from './wizard/WizardStepper';

type ConnectionWizardModalProps = {
  isOpen: boolean;
  onClose: () => void;
};

// The stepper renders each step's `icon` component bare (no size props), so it
// relies on the icon's own defaults — `CheckIcon`/`DescriptionIcon` default to
// 24px, but `ConnectionIcon` has no size default and would collapse to zero.
// Wrap it with an explicit size so the first step's icon actually shows.
const StepConnectionIcon = (props: { width?: number; height?: number }) => (
  <ConnectionIcon width={24} height={24} {...props} />
);

const ConnectionWizardModal = ({ isOpen, onClose }: ConnectionWizardModalProps) => {
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
    onComplete: onClose,
  });

  // Reset wizard state whenever the modal closes.
  useEffect(() => {
    if (!isOpen) {
      wizard.reset();
    }
    // wizard.reset is stable; avoid re-running on every render.
  }, [isOpen]);

  const steps = wizard.stepLabels.map((label, index) => ({
    label,
    icon:
      index === 0
        ? StepConnectionIcon
        : index === wizard.stepLabels.length - 1
          ? CheckIcon
          : DescriptionIcon,
    component: <></>,
  }));

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
