import React from 'react';
import {
  Modal as SistentModal,
  ModalBody,
  ModalFooter,
  PrimaryActionButtons,
} from '@sistent/sistent';

/**
 * Base Modal primitive for Meshery UI.
 * Wraps Sistent Modal and providing consistent layout for body and footer.
 */
interface ModalProps {
  open: boolean;
  title?: string;
  headerIcon?: React.ReactNode;
  headerRight?: React.ReactNode;
  closeModal: () => void;
  children: React.ReactNode;
  footer?: React.ReactNode;
  helpText?: string;
  primaryText?: string;
  secondaryText?: string;
  onPrimaryClick?: () => void;
  onSecondaryClick?: () => void;
  primaryDisabled?: boolean;
  maxWidth?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | false;
  fullWidth?: boolean;
  fullScreen?: boolean;
  'aria-labelledby'?: string;
  'aria-describedby'?: string;
}

const Modal: React.FC<ModalProps> = ({
  open,
  title,
  headerIcon,
  headerRight,
  closeModal,
  children,
  footer,
  helpText,
  primaryText,
  secondaryText = 'Cancel',
  onPrimaryClick,
  onSecondaryClick,
  primaryDisabled = false,
  maxWidth = 'sm',
  fullWidth = true,
  fullScreen = false,
  'aria-labelledby': ariaLabelledby,
  'aria-describedby': ariaDescribedby,
}) => {
  return (
    <SistentModal
      open={open}
      closeModal={closeModal}
      title={title}
      headerIcon={headerIcon}
      maxWidth={maxWidth}
      fullWidth={fullWidth}
      fullScreen={fullScreen}
      aria-labelledby={ariaLabelledby}
      aria-describedby={ariaDescribedby}
    >
      <ModalBody>{children}</ModalBody>
      {footer ? (
        <ModalFooter variant="filled" helpText={helpText} hasHelpText={!!helpText}>
          {footer}
        </ModalFooter>
      ) : onPrimaryClick ? (
        <ModalFooter variant="filled" helpText={helpText} hasHelpText={!!helpText}>
          <PrimaryActionButtons
            primaryText={primaryText || 'Submit'}
            secondaryText={secondaryText}
            primaryButtonProps={{
              onClick: onPrimaryClick,
              disabled: primaryDisabled,
            }}
            secondaryButtonProps={{
              onClick: onSecondaryClick || closeModal,
            }}
          />
        </ModalFooter>
      ) : null}
    </SistentModal>
  );
};

export default Modal;
