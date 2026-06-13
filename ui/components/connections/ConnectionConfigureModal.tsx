import { useEffect, useMemo } from 'react';
import {
  Box,
  CheckIcon,
  CustomizedStepper,
  DescriptionIcon,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  useStepper,
} from '@sistent/sistent';
import { Modal } from '@/components/shared/Modal';
import { useGetConnectionDefinitionsQuery } from '@/rtk-query/meshModel';
import {
  buildConnectionWizardKindConfigs,
  type ConnectionWizardKindConfig,
} from './ConnectionWizard.helpers';
import { useConnectionWizard } from './wizard/useConnectionWizard';

export type ConfigurableConnection = {
  id?: string;
  kind?: string;
  type?: string;
  subType?: string;
  name?: string;
  [key: string]: unknown;
};

type ConnectionConfigureModalProps = {
  isOpen: boolean;
  onClose: () => void;
  connection: ConfigurableConnection | null;
};

const matchKindConfig = (
  configs: ConnectionWizardKindConfig[],
  connection: ConfigurableConnection | null,
): ConnectionWizardKindConfig | null => {
  if (!connection?.kind) {
    return null;
  }
  const byKind = configs.filter((config) => config.kind === connection.kind);
  // Prefer the most specific type/subType match when available.
  return (
    byKind.find(
      (config) =>
        (!connection.type || config.type === connection.type) &&
        (!connection.subType || config.subType === connection.subType),
    ) ||
    byKind[0] ||
    null
  );
};

/**
 * Runs a connection's post-registration ("configure") steps + receipt for an
 * already-created connection, reusing the same per-kind extension steps the
 * creation wizard uses.
 */
const ConnectionConfigureModal = ({
  isOpen,
  onClose,
  connection,
}: ConnectionConfigureModalProps) => {
  const { data: connectionDefinitionsResponse } = useGetConnectionDefinitionsQuery(
    { params: { pagesize: 'all' } },
    { skip: !isOpen },
  );

  const kindConfigs = useMemo(
    () => buildConnectionWizardKindConfigs(connectionDefinitionsResponse?.connectionDefinitions),
    [connectionDefinitionsResponse?.connectionDefinitions],
  );

  const kindConfig = useMemo(
    () => matchKindConfig(kindConfigs, connection),
    [kindConfigs, connection],
  );

  const wizard = useConnectionWizard({
    mode: 'configure',
    isOpen,
    initialKindConfig: kindConfig,
    initialRegistrationResult: (connection as Record<string, unknown>) || null,
    onComplete: onClose,
  });

  useEffect(() => {
    if (!isOpen) {
      wizard.reset();
    }
  }, [isOpen]);

  const stepper = useStepper({
    steps: wizard.stepLabels.map((label, index) => ({
      label,
      icon: index === wizard.stepLabels.length - 1 ? CheckIcon : DescriptionIcon,
      component: <></>,
    })),
  });

  useEffect(() => {
    stepper.setActiveStep(wizard.activeIndex);
  }, [wizard.activeIndex, wizard.stepLabels.length]);

  const ActiveBody = wizard.activeStep?.Component;

  const handleClose = () => {
    if (wizard.isBusy) {
      return;
    }
    onClose();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={`Configure ${connection?.name || kindConfig?.label || 'Connection'}`}
      size="lg"
    >
      <ModalBody>
        <CustomizedStepper {...stepper}>
          {ActiveBody ? <ActiveBody ctx={wizard.ctx} /> : <></>}
        </CustomizedStepper>
      </ModalBody>
      <ModalFooter variant="filled">
        <Box sx={{ width: '100%', display: 'flex', justifyContent: 'end', gap: 2 }}>
          <ModalButtonSecondary
            onClick={wizard.canGoBack ? wizard.back : handleClose}
            disabled={wizard.isBusy}
          >
            {wizard.canGoBack ? 'Back' : 'Close'}
          </ModalButtonSecondary>
          <ModalButtonPrimary onClick={wizard.next} disabled={!wizard.canProceed || wizard.isBusy}>
            {wizard.isBusy ? 'Working...' : wizard.nextLabel}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </Modal>
  );
};

export default ConnectionConfigureModal;
