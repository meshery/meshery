/**
 * ImportDesignModal — design-domain import flow.
 *
 * Replaces both the unused `components/shared/Modal/ImportModal.tsx` and the
 * inline `ImportDesignModal` previously defined in
 * `components/designs/patterns/MesheryPatternsModals.tsx`.
 *
 * Renders the Sistent `importDesignSchema` inside the shared `FormModal`
 * primitive so the URL / upload / text-paste flow inherits the standard
 * header, footer, and validation wiring. The shared `FormModal` already
 * handles the RJSF ref + submit plumbing, so callers only need to provide
 * a close handler and a submit callback receiving the form data.
 */
import { FC, memo } from 'react';
import { importDesignSchema, importDesignUiSchema } from '@sistent/sistent';
import { FormModal } from '@/components/shared/Modal';
import { DesignModalHeaderIcon } from './design-modal-header';

export interface ImportDesignModalProps {
  /** Called when the user dismisses the modal (cancel / close icon). */
  handleClose: () => void;
  /**
   * Called with the validated import payload (URL, uploaded file, or pasted
   * YAML/JSON) — same shape RJSF emits for `importDesignSchema`.
   */
  handleImportDesign: (formData: unknown) => void;
}

const ImportDesignModalComponent: FC<ImportDesignModalProps> = ({
  handleClose,
  handleImportDesign,
}) => (
  <FormModal
    isOpen
    onClose={handleClose}
    title="Import Design"
    headerIcon={<DesignModalHeaderIcon />}
    size="sm"
    schema={importDesignSchema}
    uiSchema={importDesignUiSchema}
    submitText="Import"
    onSubmit={handleImportDesign}
  />
);

ImportDesignModalComponent.displayName = 'ImportDesignModal';

export const ImportDesignModal = memo(ImportDesignModalComponent);

export default ImportDesignModal;
