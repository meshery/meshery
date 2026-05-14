/**
 * Shared Modal primitive.
 *
 * Canonical base modal for Meshery UI. Wraps Sistent's `Modal`, `ModalBody`,
 * `ModalFooter` so every modal in the application:
 *
 *   - has the same prop shape (`isOpen`/`onClose`/`title`/...).
 *   - composes a consistent header, body, and footer.
 *   - exposes the same `size` token vocabulary (`sm`/`md`/`lg`/`xl`/`fullscreen`).
 *   - threads `actions` through `ModalFooter` so callers don't reimplement
 *     button alignment.
 *
 * Use this primitive directly for ad-hoc dialogs; prefer `ConfirmModal`,
 * `InfoModal`, or `FormModal` for the corresponding higher-level patterns.
 *
 * NOTE: Legacy callers still import the previous RJSF-backed `Modal` default
 * export and the `RJSFModalWrapper` named export from this file path. Those
 * are preserved as a temporary shim re-exported from `./LegacyRJSFModal` and
 * will be removed after the Phase 5.b migration sub-issues (#18752–#18756).
 */
import { FC, ReactNode } from 'react';
import {
  Modal as SistentModal,
  ModalBody as SistentModalBody,
  ModalFooter as SistentModalFooter,
} from '@sistent/sistent';

export type ModalSize = 'sm' | 'md' | 'lg' | 'xl' | 'fullscreen';

export type ModalFooterVariant = 'filled' | 'transparent';

export interface ModalProps {
  /** Whether the modal is currently visible. */
  isOpen: boolean;
  /** Called when the user dismisses the modal (close icon, backdrop, ESC). */
  onClose: () => void;
  /** Heading rendered in the modal header. */
  title: string;
  /** Optional icon rendered to the left of the title. */
  headerIcon?: ReactNode;
  /** Body content. Prefer placing structured content inside a child layout. */
  children?: ReactNode;
  /** Footer content; typically a row of primary/secondary buttons. */
  actions?: ReactNode;
  /** Footer variant. `filled` matches the Sistent default action-bar styling. */
  footerVariant?: ModalFooterVariant;
  /** Optional help-text rendered in the footer alongside `actions`. */
  helpText?: ReactNode;
  /** Pre-set width token. Defaults to `md`. */
  size?: ModalSize;
  /** Toggle Sistent's built-in fullscreen mode button. */
  isFullScreenModeAllowed?: boolean;
  /**
   * When `true`, the modal does NOT wrap children in `ModalBody` and does NOT
   * render a `ModalFooter` for `actions`. Use this for callers (e.g. stepper
   * flows) that emit their own `ModalBody` and `ModalFooter` siblings.
   */
  disableBodyWrap?: boolean;
  /**
   * Forwarded to the underlying Sistent `Modal` (MUI `Dialog`) so consumers
   * can target the dialog root with `styled(Modal)(...)` for project-specific
   * size or stacking overrides.
   */
  className?: string;
  /** Forwarded to the underlying Dialog root for cypress/test selectors. */
  'aria-labelledby'?: string;
  'aria-describedby'?: string;
}

const sizeToMaxWidth = {
  sm: 'sm',
  md: 'md',
  lg: 'lg',
  xl: 'xl',
  fullscreen: false,
} as const;

export const Modal: FC<ModalProps> = ({
  isOpen,
  onClose,
  title,
  headerIcon,
  children,
  actions,
  footerVariant = 'filled',
  helpText,
  size = 'md',
  isFullScreenModeAllowed,
  disableBodyWrap = false,
  className,
  ...ariaProps
}) => {
  return (
    <SistentModal
      open={isOpen}
      closeModal={onClose}
      title={title}
      headerIcon={headerIcon}
      maxWidth={sizeToMaxWidth[size]}
      fullScreen={size === 'fullscreen'}
      fullWidth
      isFullScreenModeAllowed={isFullScreenModeAllowed}
      className={className}
      {...ariaProps}
    >
      {disableBodyWrap ? children : <SistentModalBody>{children}</SistentModalBody>}
      {!disableBodyWrap && actions ? (
        <SistentModalFooter
          variant={footerVariant}
          // Sistent types `helpText` as string but the runtime element renders
          // any node — widen here so callers can pass JSX (e.g. a Docs link).
          helpText={helpText as string}
          hasHelpText={!!helpText}
        >
          {actions}
        </SistentModalFooter>
      ) : null}
    </SistentModal>
  );
};

Modal.displayName = 'Modal';

// -----------------------------------------------------------------------------
// Back-compat re-exports for legacy callers (`import Modal from
// 'shared/Modal/Modal'` and `{ RJSFModalWrapper }`). These will be removed
// once the Phase 5.b migration sub-issues complete.
// -----------------------------------------------------------------------------
import LegacyRJSFModal from './LegacyRJSFModal';

export { RJSFModalWrapper } from './LegacyRJSFModal';

export default LegacyRJSFModal;
