/**
 * ConnectionFormModal
 *
 * Overlay dialog that surfaces the full `ConnectionTable` filtered to
 * Kubernetes connections. Migrated onto the shared `Modal` primitive from
 * `@/components/shared/Modal` so the prop shape, header, and sizing tokens
 * stay consistent across the modal surface. Consumer-facing props
 * (`isOpenModal`, `setIsOpenModal`, ...) are preserved unchanged.
 */
import { FC } from 'react';
import { styled } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import ConnectionIcon from '@/assets/icons/Connection';
import ConnectionTable from './ConnectionTable';

interface ConnectionFormModalProps {
  isOpenModal: boolean;
  setIsOpenModal: (open: boolean) => void;
  /**
   * Preserved for prop compatibility with legacy callers. `ConnectionTable`
   * now sources controller state from the Redux store directly, so this prop
   * is accepted but unused.
   */
  meshsyncControllerState?: unknown;
  /**
   * Preserved for prop compatibility with legacy callers. `ConnectionTable`
   * now sources connection-metadata state from the Redux store directly.
   */
  connectionMetadataState?: unknown;
}

const TableShell = styled('div')(({ theme }) => ({
  marginBlock: theme.spacing(4),
  maxHeight: '65vh',
}));

const ConnectionFormModal: FC<ConnectionFormModalProps> = ({ isOpenModal, setIsOpenModal }) => (
  <Modal
    isOpen={isOpenModal}
    onClose={() => setIsOpenModal(false)}
    title="Connections"
    size="xl"
    headerIcon={<ConnectionIcon height={24} width={24} />}
  >
    <TableShell>
      <ConnectionTable selectedFilter="kubernetes" />
    </TableShell>
  </Modal>
);

export default ConnectionFormModal;
