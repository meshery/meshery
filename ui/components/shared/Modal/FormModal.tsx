/**
 * Shared FormModal primitive.
 *
 * Composition of the base `Modal` with the project's shared RJSF wrapper so
 * caller code can render a JSON-schema-driven form inside a dialog without
 * re-wiring the form ref, validation, or submit handling each time.
 *
 * The data flow mirrors the legacy `RJSFModalWrapper` so migration of existing
 * callers stays straightforward:
 *
 *   <FormModal
 *     isOpen={isOpen}
 *     onClose={onClose}
 *     title="Create environment"
 *     schema={schema}
 *     uiSchema={uiSchema}
 *     initialData={data}
 *     onSubmit={(formData) => save(formData)}
 *     submitText="Create"
 *   />
 *
 * For freeform forms (e.g. react-hook-form) pass `children` instead of
 * `schema`; the modal will render the children verbatim inside the body and
 * the submit handler is responsible for whatever form library the caller
 * uses. This keeps FormModal a thin layout primitive rather than a forced
 * coupling to RJSF.
 */
import { FC, ReactNode, useRef } from 'react';
import { CircularProgress, ModalButtonPrimary, ModalButtonSecondary } from '@sistent/sistent';
import { styled } from '@/theme';
import RJSFWrapper from '../../meshery-mesh-interface/PatternService/RJSF_wrapper';
import { Modal, ModalSize } from './Modal';

const LoadingShell = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(8, 4),
}));

// `RJSFWrapper` exposes the underlying RJSF form via the ref's `current`. We
// only need `validateForm()` and the working `formData` at submit time; this
// shape mirrors what the RJSF Form component itself assigns to the ref.
type RjsfFormRef = {
  current: null | {
    validateForm: () => boolean;
    state: { formData: unknown };
  };
};

interface FormModalBaseProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  headerIcon?: ReactNode;
  /** Footer help text, e.g. linking to docs. Renders inside the action bar. */
  helpText?: ReactNode;
  /** Primary action label. Defaults to `Submit`. */
  submitText?: string;
  /** Secondary action label. Defaults to `Cancel`. */
  cancelText?: string;
  /** Disable the submit button (e.g. while the action is in flight). */
  isSubmitDisabled?: boolean;
  /** Size token; defaults to `md`. */
  size?: ModalSize;
  sx?: object;
}

export interface RjsfFormModalProps<TFormData = unknown> extends FormModalBaseProps {
  /** JSON schema describing the form fields. */
  schema: object;
  /** Optional RJSF UI schema for widget/layout overrides. */
  uiSchema?: object;
  /** Initial values for the form. */
  initialData?: TFormData;
  /** Custom widget overrides forwarded to RJSF. */
  widgets?: Record<string, unknown>;
  /** Called with validated form data when the user submits. */
  onSubmit: (formData: TFormData) => void;
}

export interface FreeformFormModalProps extends FormModalBaseProps {
  /** Custom form rendering (e.g. react-hook-form). */
  children: ReactNode;
  /**
   * Called when the user clicks the submit button. Caller is responsible for
   * validating its form and reading the latest data.
   */
  onSubmit: () => void;
}

export type FormModalProps<TFormData = unknown> =
  | RjsfFormModalProps<TFormData>
  | FreeformFormModalProps;

const isRjsfMode = <T,>(props: FormModalProps<T>): props is RjsfFormModalProps<T> =>
  'schema' in props;

export const FormModal = <TFormData,>(props: FormModalProps<TFormData>) => {
  const {
    isOpen,
    onClose,
    title,
    headerIcon,
    helpText,
    submitText = 'Submit',
    cancelText = 'Cancel',
    isSubmitDisabled = false,
    size = 'md',
    sx,
  } = props;

  const rjsfRef = useRef<RjsfFormRef['current']>(null);

  let body: ReactNode;
  let submitHandler: () => void;

  if (isRjsfMode(props)) {
    const { schema, uiSchema = {}, initialData, widgets = {} } = props;
    body = schema ? (
      <RJSFWrapper
        formData={initialData}
        jsonSchema={schema}
        uiSchema={uiSchema}
        onChange={undefined}
        liveValidate={false}
        formRef={rjsfRef}
        hideTitle
        widgets={widgets}
      />
    ) : (
      <LoadingShell>
        <CircularProgress />
      </LoadingShell>
    );
    submitHandler = () => {
      const form = rjsfRef.current;
      if (form && form.validateForm()) {
        props.onSubmit(form.state.formData as TFormData);
      }
    };
  } else {
    body = props.children;
    submitHandler = props.onSubmit;
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      headerIcon={headerIcon}
      helpText={helpText}
      size={size}
      sx={sx}
      actions={
        <>
          <ModalButtonSecondary onClick={onClose}>{cancelText}</ModalButtonSecondary>
          <ModalButtonPrimary onClick={submitHandler} disabled={isSubmitDisabled}>
            {submitText}
          </ModalButtonPrimary>
        </>
      }
    >
      {body}
    </Modal>
  );
};

(FormModal as FC).displayName = 'FormModal';

export default FormModal;
