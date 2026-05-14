/**
 * EnvironmentFormModal
 *
 * Overlay dialog that hosts the full `EnvironmentComponent` management view.
 * Migrated onto the shared `Modal` primitive from `@/components/shared/Modal`
 * so the prop shape, header, and sizing tokens are consistent with the rest
 * of the modal surface. The consumer-facing prop names (`isOpenModal`,
 * `setIsOpenModal`) are preserved to avoid touching callers.
 */
import { FC } from 'react';
import { EnvironmentIcon } from '@sistent/sistent';
import { useTheme, styled } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { EnvironmentComponent } from '../lifecycle';

interface EnvironmentFormModalProps {
  isOpenModal: boolean;
  setIsOpenModal: (open: boolean) => void;
}

const BodyShell = styled('div')({
  maxHeight: '65vh',
});

const EnvironmentFormModal: FC<EnvironmentFormModalProps> = ({ isOpenModal, setIsOpenModal }) => {
  const theme = useTheme();

  if (!isOpenModal) {
    return null;
  }

  return (
    <Modal
      isOpen={isOpenModal}
      onClose={() => setIsOpenModal(false)}
      title="Environment"
      size="xl"
      headerIcon={
        <EnvironmentIcon height={24} width={24} fill={theme.palette.background.constant.white} />
      }
    >
      <BodyShell>
        <EnvironmentComponent />
      </BodyShell>
    </Modal>
  );
};

export default EnvironmentFormModal;
