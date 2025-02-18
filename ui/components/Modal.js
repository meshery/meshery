import React, { useState, useEffect, useRef } from 'react';
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
} from '@layer5/sistent';
import RJSFWrapper from './MesheryMeshInterface/PatternService/RJSF_wrapper';
import { ArrowDropDown } from '@mui/icons-material';
import { getSchema } from './MesheryMeshInterface/PatternService/helper';
import { UsesSistent } from './SistentWrapper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

const SchemaVersion = ({ schema_array, type, schemaChangeHandler }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <div>
      <Tooltip title="Schema_Changer">
        <IconButton component="span" onClick={(e) => setAnchorEl(e.currentTarget)}>
          <ArrowDropDown style={{ color: '#000' }} />
        </IconButton>
      </Tooltip>
      <Menu id="schema-menu" anchorEl={anchorEl} open={open} handleClose={handleClose}>
        {schema_array.map((version, index) => (
          <MenuItem
            id="schema-menu-item"
            key={index}
            selected={version === type}
            onClick={() => {
              schemaChangeHandler(version);
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
  } = props;

  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const formStateRef = useRef({});
  const formRef = React.createRef();
  const [loadingSchema, setLoadingSchema] = useState(true);
  const { notify } = useNotification();
  const handleFormSubmit = () => {
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
    if (schema) {
      setLoadingSchema(false);
    }
  }, [schema]);

  return (
    <UsesSistent>
      <SistentModal open={open} closeModal={handleClose} title={title} headerIcon={leftHeaderIcon}>
        <Typography variant="h5">
          {schema_array?.length < 1 && (
            <SchemaVersion
              schema_array={schema_array}
              type={type}
              schemaChangeHandler={schemaChangeHandler}
            />
          )}
        </Typography>
        <ModalBody>
          {loadingSchema ? (
            <div style={{ textAlign: 'center', padding: '8rem 17rem' }}>
              <CircularProgress />
            </div>
          ) : (
            <RJSFWrapper
              key={type}
              formData={initialData || formStateRef}
              jsonSchema={schema || getSchema(type)}
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
              disabled: canNotSubmit,
            }}
            secondaryButtonProps={{
              onClick: handleClose,
            }}
          />
        </ModalFooter>
      </SistentModal>
    </UsesSistent>
  );
}

export default React.memo(Modal);

function RJSFModalWrapper({
  handleClose,
  schema,
  uiSchema = {},
  initialData = {},
  handleSubmit,
  title,
  submitBtnText,
  helpText,
}) {
  const formRef = useRef();
  const formStateRef = useRef();
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

  const handleFormChange = (data) => {
    formStateRef.current = data;
  };

  useEffect(() => {
    setLoadingSchema(!schema);
  }, [schema]);

  const handleFormSubmit = () => {
    if (formRef.current && formRef.current.validateForm()) {
      handleSubmit(formRef.current.state.formData);
      handleClose();
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
          />
        )}
      </ModalBody>
      <ModalFooter variant="filled" helpText={helpText}>
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
}

export { RJSFModalWrapper };
