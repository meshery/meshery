import { useEffect, useMemo } from 'react';
import {
  Box,
  CheckIcon,
  DescriptionIcon,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
} from '@sistent/sistent';
import { Modal } from '@/components/shared/Modal';
import { useListConnectionDefinitionsQuery } from '@meshery/schemas/mesheryApi';
import {
  buildConnectionWizardKindConfigs,
  DEFAULT_CONNECTION_DOCS_URL,
  type ConnectionWizardKindConfig,
} from './ConnectionWizard.helpers';
import { useConnectionWizard } from './wizard/useConnectionWizard';
import WizardStepper from './wizard/WizardStepper';

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
  const { data: connectionDefinitionsResponse } = useListConnectionDefinitionsQuery(
    {},
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

  const steps = wizard.stepLabels.map((label, index) => ({
    label,
    icon: index === wizard.stepLabels.length - 1 ? CheckIcon : DescriptionIcon,
    component: <></>,
  }));

  const ActiveBody = wizard.activeStep?.Component;

  const helpText = useMemo(() => {
    const stepHelp = wizard.activeStep?.helpText;
    if (typeof stepHelp === 'function') {
      return stepHelp(wizard.ctx);
    }
    if (typeof stepHelp === 'string' && stepHelp.length > 0) {
      return stepHelp;
    }
    const label = kindConfig?.label || connection?.kind || 'connection';
    const docsUrl = kindConfig?.docsUrl || DEFAULT_CONNECTION_DOCS_URL;
    return `Configure this ${label} connection. [Learn more about connections](${docsUrl}).`;
  }, [wizard.activeStep, wizard.ctx, kindConfig, connection?.kind]);

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
