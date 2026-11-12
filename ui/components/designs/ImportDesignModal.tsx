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
import { FC, memo, useCallback, useRef, useState } from 'react';
import { importDesignSchema, importDesignUiSchema } from '@sistent/sistent';
import { FormModal } from '@/components/shared/Modal';
import { DesignModalHeaderIcon } from './design-modal-header';

export interface ImportDesignModalProps {
  /** Called when the user dismisses the modal (cancel / close icon). */
  handleClose: () => void;
  /**
   * Called with the validated import payload (URL, uploaded file, or pasted
   * YAML/JSON) — same shape RJSF emits for `importDesignSchema`.
   * May return a Promise; the modal disables the submit button until it settles.
   */
  handleImportDesign: (formData: unknown) => void | Promise<void>;
}

const ImportDesignModalComponent: FC<ImportDesignModalProps> = ({
  handleClose,
  handleImportDesign,
}) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const isSubmittingRef = useRef(false);

  const handleSubmit = useCallback(
    async (formData: unknown) => {
      if (isSubmittingRef.current) return;
      isSubmittingRef.current = true;
      setIsSubmitting(true);
      try {
        await handleImportDesign(formData);
      } finally {
        isSubmittingRef.current = false;
        setIsSubmitting(false);
      }
    },
    [handleImportDesign],
  );

  return (
    <FormModal
      isOpen
      onClose={isSubmitting ? () => {} : handleClose}
      title="Import Design"
      headerIcon={<DesignModalHeaderIcon />}
      size="sm"
      schema={importDesignSchema}
      uiSchema={importDesignUiSchema}
      submitText={isSubmitting ? 'Importing…' : 'Import'}
      isSubmitDisabled={isSubmitting}
      onSubmit={handleSubmit}
    />
  );
};

ImportDesignModalComponent.displayName = 'ImportDesignModal';

export const ImportDesignModal = memo(ImportDesignModalComponent);

export default ImportDesignModal;
