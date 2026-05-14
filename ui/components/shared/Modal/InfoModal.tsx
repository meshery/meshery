/**
 * Shared InfoModal primitive.
 *
 * Read-only informational dialog: a title plus a body that can be either a
 * string (rendered as markdown via Sistent's `RenderMarkdown`) or arbitrary
 * React content. Use this for help text, release notes, learn-more popups
 * and other non-actionable surfaces.
 *
 * For the more involved catalog/publish flow that the legacy
 * `Information/InfoModal` covered, see `./Information/LegacyInfoModal`. That
 * surface will be modernised under its own migration sub-issue.
 */
import { FC, ReactNode } from 'react';
import { ModalButtonPrimary, RenderMarkdown } from '@sistent/sistent';
import { Modal, ModalSize } from './Modal';

export interface InfoModalProps {
  /** Whether the info modal is visible. */
  isOpen: boolean;
  /** Called when the user dismisses the modal. */
  onClose: () => void;
  /** Modal heading. */
  title: string;
  /**
   * Body content. When a string is provided it's rendered through
   * `RenderMarkdown` so callers can use markdown for emphasis or links;
   * pass React content for custom layouts.
   */
  body: ReactNode;
  /** Optional icon rendered in the header. */
  headerIcon?: ReactNode;
  /** Label for the dismiss button. Defaults to `Close`. */
  closeText?: string;
  /** Optional help-text rendered alongside the action in the footer. */
  helpText?: string;
  /** Size token; defaults to `md`. */
  size?: ModalSize;
  /** Hide the footer action entirely (e.g. when caller renders its own). */
  hideCloseButton?: boolean;
}

export const InfoModal: FC<InfoModalProps> = ({
  isOpen,
  onClose,
  title,
  body,
  headerIcon,
  closeText = 'Close',
  helpText,
  size = 'md',
  hideCloseButton = false,
}) => {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      headerIcon={headerIcon}
      size={size}
      helpText={helpText}
      actions={
        hideCloseButton ? null : (
          <ModalButtonPrimary onClick={onClose}>{closeText}</ModalButtonPrimary>
        )
      }
    >
      {typeof body === 'string' ? <RenderMarkdown content={body} /> : body}
    </Modal>
  );
};

InfoModal.displayName = 'InfoModal';

export default InfoModal;
