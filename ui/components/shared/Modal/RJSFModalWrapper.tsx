import React, { useRef, useState, useEffect } from 'react';
import { CircularProgress, ModalBody, ModalFooter, PrimaryActionButtons } from '@sistent/sistent';
import RJSFWrapper from '../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

interface RJSFModalWrapperProps {
  handleClose: () => void;
  schema: any;
  uiSchema?: any;
  initialData?: any;
  handleSubmit: (data: any) => void;
  handleNext?: () => void;
  title?: string;
  submitBtnText?: string;
  helpText?: React.ReactNode;
  widgets?: any;
}

/**
 * RJSFModalWrapper provides a standard way to render RJSF forms with a footer
 * inside a ModalBody. It includes validation and specific Meshery logic
 * like the "Untitled Design" check.
 */
export const RJSFModalWrapper: React.FC<RJSFModalWrapperProps> = ({
  handleClose,
  schema,
  uiSchema = {},
  initialData = {},
  handleSubmit,
  handleNext,
  title,
  submitBtnText,
  helpText,
  widgets = {},
}) => {
  const formRef = useRef<any>();
  const formStateRef = useRef<any>();
  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const [loadingSchema, setLoadingSchema] = useState(true);
  const { notify } = useNotification();

  useEffect(() => {
    setCanNotSubmit(false);
    const handleDesignNameCheck = () => {
      const designName = title?.toLowerCase();
      const forbiddenWords = ['untitled design', 'Untitled', 'lfx'];

      for (const word of forbiddenWords) {
        if (designName?.includes(word)) {
          notify({
            event_type: EVENT_TYPES.WARNING,
            message: `Design name should not contain Untitled Design, Untitled, LFX`,
          });
          setCanNotSubmit(true);
          break;
        }
      }
    };
    handleDesignNameCheck();
  }, [title]);

  const handleFormChange = (data: any) => {
    formStateRef.current = data;
  };

  useEffect(() => {
    setLoadingSchema(!schema);
  }, [schema]);

  const handleFormSubmit = () => {
    if (formRef.current && formRef.current.validateForm()) {
      handleSubmit(formRef.current.state.formData);
      if (handleNext) {
        handleNext();
      }
    }
  };

  return (
    <>
      <ModalBody>
        {loadingSchema ? (
          <div style={{ textAlign: 'center', padding: '8rem 17rem' }}>
            <CircularProgress />
          </div>
        ) : (
          <RJSFWrapper
            formData={initialData}
            jsonSchema={schema}
            uiSchema={uiSchema}
            onChange={handleFormChange}
            liveValidate={false}
            formRef={formRef}
            hideTitle={true}
            widgets={widgets}
          />
        )}
      </ModalBody>
      <ModalFooter variant="filled" helpText={helpText} hasHelpText={!!helpText}>
        <PrimaryActionButtons
          primaryText={submitBtnText || 'Submit'}
          secondaryText="Cancel"
          primaryButtonProps={{
            onClick: handleFormSubmit,
            disabled: canNotSubmit,
          }}
          secondaryButtonProps={{
            onClick: handleClose,
          }}
        />
      </ModalFooter>
    </>
  );
};
