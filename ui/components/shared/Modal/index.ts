/**
 * Shared Modal barrel.
 *
 * Canonical entry point for the modal primitives. Prefer this barrel over
 * deep imports so call sites stay agnostic to internal file layout:
 *
 *   import { Modal, ConfirmModal, InfoModal, FormModal, useModal } from
 *     '@/components/shared/Modal';
 */
export { Modal } from './Modal';
export type { ModalProps, ModalSize, ModalFooterVariant } from './Modal';

export { ConfirmModal } from './ConfirmModal';
export type { ConfirmModalProps, ConfirmModalVariant } from './ConfirmModal';

export { InfoModal } from './InfoModal';
export type { InfoModalProps } from './InfoModal';

export { FormModal } from './FormModal';
export type { FormModalProps, RjsfFormModalProps, FreeformFormModalProps } from './FormModal';

export { useModal } from './useModal';
export type { UseModalReturn } from './useModal';
