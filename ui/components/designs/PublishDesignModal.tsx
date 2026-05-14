/**
 * PublishDesignModal — catalog-publish flow for designs.
 *
 * Replaces both the unused `components/shared/Modal/PublishModal.tsx` and the
 * inline `PublishModal` previously defined in
 * `components/designs/patterns/MesheryPatternsModals.tsx`.
 *
 * Composed on top of the shared `FormModal` primitive so the catalog publish
 * form (compatibility, license, screenshots, …) renders inside the same
 * header/footer chrome as every other design modal. The caller supplies the
 * already-populated `publishFormSchema` (typically derived from
 * `publishCatalogItemSchema` overlaid with discovered mesh models); when it
 * is absent the static Sistent schema is used as a fallback so the modal
 * still opens with a usable form.
 */
import { FC, memo } from 'react';
import { publishCatalogItemSchema, publishCatalogItemUiSchema } from '@sistent/sistent';
import { FormModal } from '@/components/shared/Modal';
import { DesignModalHeaderIcon } from './design-modal-header';

export interface PublishDesignFormSchema {
  rjsfSchema?: object;
  uiSchema?: object;
}

export interface PublishDesignModalProps {
  /** Dialog heading; typically the design's name. */
  title: string;
  /** Called when the user dismisses the modal (cancel / close icon). */
  handleClose: () => void;
  /** Called with the validated catalog-publish payload. */
  handleSubmit: (formData: unknown) => void;
  /**
   * Pre-built publish form schema (mesh-model-aware). Falls back to the
   * static `publishCatalogItemSchema` from Sistent when omitted.
   */
  publishFormSchema?: PublishDesignFormSchema;
}

const PublishDesignModalComponent: FC<PublishDesignModalProps> = ({
  title,
  handleClose,
  handleSubmit,
  publishFormSchema,
}) => {
  // Prefer the dynamically-populated schema (compatibility enum overlaid
  // with discovered mesh models); fall back to the static Sistent schema
  // when the caller hasn't supplied one yet.
  const schema = publishFormSchema?.rjsfSchema ?? publishCatalogItemSchema;
  const uiSchema = publishFormSchema?.uiSchema ?? publishCatalogItemUiSchema;

  return (
    <FormModal
      isOpen
      onClose={handleClose}
      title={title}
      headerIcon={<DesignModalHeaderIcon />}
      size="sm"
      schema={schema}
      uiSchema={uiSchema}
      submitText="Submit for Approval"
      onSubmit={handleSubmit}
      helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
    />
  );
};

PublishDesignModalComponent.displayName = 'PublishDesignModal';

export const PublishDesignModal = memo(PublishDesignModalComponent);

export default PublishDesignModal;
