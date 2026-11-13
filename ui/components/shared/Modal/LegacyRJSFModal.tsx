import React, { useState, useEffect, useRef, FC } from 'react';
import {
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
  Typography,
  CircularProgress,
  ModalBody,
  ModalFooter,
  PrimaryActionButtons,
  Modal as SistentModal,
  useTheme,
} from '@sistent/sistent';
import RJSFWrapper from '../../meshery-mesh-interface/PatternService/RJSF_wrapper';
import { ArrowDropDown } from '@/assets/icons';
import { getSchema } from '../../meshery-mesh-interface/PatternService/helper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

interface SchemaVersionProps {
  schema_array: string[];
  type: string;
  schemaChangeHandler: (version: string) => void;
}

const SchemaVersion: FC<SchemaVersionProps> = ({ schema_array, type, schemaChangeHandler }) => {
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <div>
      <Tooltip title="Schema_Changer">
        <IconButton
          aria-label="Change schema version"
          component="span"
          onClick={(e) => setAnchorEl(e.currentTarget)}
        >
          <ArrowDropDown style={{ color: theme.palette.text.primary }} />
        </IconButton>
      </Tooltip>
      <Menu id="schema-menu" anchorEl={anchorEl} open={open} handleClose={handleClose}>
        {schema_array?.map((version, index) => (
          <MenuItem
            id="schema-menu-item"
            key={version || index}
            selected={version === type}
            onClick={() => {
              schemaChangeHandler?.(version);
              handleClose();
            }}
          >
            {version}
          </MenuItem>
        ))}
      </Menu>
    </div>
  );
};

/**
 * Renders common dialog component.
 *
 * @param {Object} props - Component props.
 * @param {boolean} props.open - Determines whether the modal is open or not.
 * @param {string} props.title - The title of the modal.
 * @param {Function} props.handleClose - Function to handle the close event of the modal.
 * @param {Function} props.onChange - Function to handle the change event of the form fields in the modal.
 * @param {Object} props.schema - The JSON schema for the form fields.
 * @param {Object} props.formData - The form data object.
 * @param {Array} props.schema_array - An array of schema versions.
 * @param {string} props.type - The selected schema version.
 * @param {Function} props.schemaChangeHandler - Function to handle the change of the schema version.
 * @param {Function} props.handleSubmit - Function to handle the submit event of the modal.
 * @param {Object} props.payload - The payload for the submit event.
 * @param {Object} props.showInfoIcon - Determines whether to show the info icon adjacent to the modal button.
 * @param {string} props.submitBtnText - The text for the submit button.
 * @param {Object} props.uiSchema - The UI schema for the form fields.
 * @param {boolean} [props.submitDisabled=false] - Optional consumer-controlled flag to disable the submit button (OR-ed with forbidden title and schema-loading pending states).
 */

// Meshery extensions also uses this modal
function Modal(props) {
  const {
    open,
    title,
    handleClose,
    schema,
    schema_array,
    type,
    schemaChangeHandler,
    handleSubmit,
    submitBtnText,
    leftHeaderIcon,
    uiSchema = {},
    helpText,
    RJSFWrapperComponent = null,
    initialData = {},
    submitDisabled = false,
  } = props;

  const resolvedSchema = schema !== undefined ? schema : getSchema(type);
  const hasResolvedSchema = Boolean(resolvedSchema);
  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const formStateRef = useRef({});
  const formRef = React.createRef();
  const [loadingSchema, setLoadingSchema] = useState(!hasResolvedSchema);
  const { notify } = useNotification();
  const isSubmitDisabled = Boolean(canNotSubmit || loadingSchema || submitDisabled);

  const handleFormSubmit = () => {
    if (isSubmitDisabled) {
      return;
    }
    if (formRef.current && formRef.current.validateForm()) {
      handleClose();
      handleSubmit(formRef.current.state.formData);
    }
  };

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

  const handleFormChange = (data) => {
    formStateRef.current = data;
  };

  useEffect(() => {
    setLoadingSchema(!hasResolvedSchema);
  }, [hasResolvedSchema]);

  return (
    <>
      <SistentModal open={open} closeModal={handleClose} title={title} headerIcon={leftHeaderIcon}>
        {schema_array && schema_array.length > 0 && (
          <Typography component="div" variant="h5">
            <SchemaVersion
              schema_array={schema_array}
              type={type}
              schemaChangeHandler={schemaChangeHandler}
            />
          </Typography>
        )}
        <ModalBody>
          {loadingSchema ? (
            <div style={{ textAlign: 'center', padding: '8rem 17rem' }}>
              <CircularProgress />
            </div>
          ) : (
            <RJSFWrapper
              key={type}
              formData={initialData || formStateRef}
              jsonSchema={resolvedSchema}
              uiSchema={uiSchema}
              onChange={handleFormChange}
              liveValidate={false}
              formRef={formRef}
              hideTitle={true}
              {...(RJSFWrapperComponent && { RJSFWrapperComponent })}
            />
          )}
        </ModalBody>
        <ModalFooter variant="filled" helpText={helpText} hasHelpText={!!helpText}>
          <PrimaryActionButtons
            primaryText={submitBtnText || 'Submit'}
            secondaryText="Cancel"
            primaryButtonProps={{
              onClick: handleFormSubmit,
              disabled: isSubmitDisabled,
            }}
            secondaryButtonProps={{
              onClick: handleClose,
            }}
          />
        </ModalFooter>
      </SistentModal>
    </>
  );
}

export default React.memo(Modal);

/**
 * Renders an RJSF form body and action footer inside an existing modal container.
 *
 * @param {Object} props - Component props.
 * @param {Function} props.handleClose - Function to handle closing the modal.
 * @param {Object} props.schema - The JSON schema for the form fields.
 * @param {Object} [props.uiSchema={}] - The UI schema for the form fields.
 * @param {Object} [props.initialData={}] - Initial form data.
 * @param {Function} props.handleSubmit - Function invoked with validated form data on submit.
 * @param {Function} [props.handleNext] - Optional step callback invoked after submit.
 * @param {string} props.title - Modal title used for forbidden-name checks.
 * @param {string} [props.submitBtnText] - Label for the primary submit button.
 * @param {string} [props.helpText] - Optional footer help text.
 * @param {Object} [props.widgets={}] - Custom RJSF widgets.
 * @param {boolean} [props.isSubmitting=false] - Shows a spinner and disables submit while submitting.
 * @param {boolean} [props.submitDisabled=false] - Optional consumer-controlled flag to disable submit.
 */
function RJSFModalWrapper({
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
  isSubmitting = false,
  submitDisabled = false,
}) {
  const formRef = useRef();
  const formStateRef = useRef();
  const hasSchema = Boolean(schema);
  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const [loadingSchema, setLoadingSchema] = useState(!hasSchema);
  const { notify } = useNotification();
  const isSubmitDisabled = Boolean(canNotSubmit || loadingSchema || isSubmitting || submitDisabled);
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

  const handleFormChange = (data) => {
    formStateRef.current = data;
  };

  useEffect(() => {
    setLoadingSchema(!hasSchema);
  }, [hasSchema]);

  const handleFormSubmit = () => {
    if (isSubmitDisabled) {
      return;
    }
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
      <ModalFooter variant="filled" helpText={helpText}>
        <PrimaryActionButtons
          primaryText={submitBtnText || 'Submit'}
          secondaryText="Cancel"
          primaryButtonProps={{
            onClick: handleFormSubmit,
            disabled: isSubmitDisabled,
            startIcon: isSubmitting ? <CircularProgress size={16} /> : undefined,
          }}
          secondaryButtonProps={{
            onClick: handleClose,
          }}
        />
      </ModalFooter>
    </>
  );
}

export { RJSFModalWrapper };
