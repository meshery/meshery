import React, { useEffect, useMemo, useRef, useState, type ComponentType } from 'react';
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
import RJSFWrapper from '../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import { ArrowDropDown } from '@/components/icons';
import { getSchema } from '../../MesheryMeshInterface/PatternService/helper';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

type ModalFormData = Record<string, unknown>;

type ModalFormRef = {
  validateForm?: () => boolean;
  state?: {
    formData?: ModalFormData;
  };
};

type WrapperComponentProps = Record<string, unknown> & {
  children?: React.ReactNode;
};

type SchemaVersionProps = {
  schemaArray?: string[];
  type?: string;
  schemaChangeHandler?: (version: string) => void;
};

type ModalProps = {
  open: boolean;
  title?: string;
  handleClose: () => void;
  schema?: Record<string, unknown>;
  schema_array?: string[];
  type?: string;
  schemaChangeHandler?: (version: string) => void;
  handleSubmit: (formData: ModalFormData) => void;
  submitBtnText?: string;
  leftHeaderIcon?: React.ReactNode;
  uiSchema?: Record<string, unknown>;
  helpText?: string;
  RJSFWrapperComponent?: ComponentType<WrapperComponentProps> | null;
  initialData?: ModalFormData;
};

type RJSFModalWrapperProps = {
  handleClose: () => void;
  schema?: Record<string, unknown>;
  uiSchema?: Record<string, unknown>;
  initialData?: ModalFormData;
  handleSubmit: (formData: ModalFormData) => void;
  handleNext?: () => void;
  title?: string;
  submitBtnText?: string;
  helpText?: string;
  widgets?: Record<string, unknown>;
};

const FORBIDDEN_DESIGN_NAME_TOKENS = ['untitled design', 'untitled', 'lfx'] as const;

const containsForbiddenDesignName = (title?: string) => {
  const normalizedTitle = title?.toLowerCase();

  return FORBIDDEN_DESIGN_NAME_TOKENS.some((token) => normalizedTitle?.includes(token));
};

const SchemaVersion = ({ schemaArray = [], type, schemaChangeHandler }: SchemaVersionProps) => {
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const open = Boolean(anchorEl);

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <div>
      <Tooltip title="Schema_Changer">
        <IconButton component="span" onClick={(event) => setAnchorEl(event.currentTarget)}>
          <ArrowDropDown style={{ color: theme.palette.text.primary }} />
        </IconButton>
      </Tooltip>
      <Menu id="schema-menu" anchorEl={anchorEl} open={open} handleClose={handleClose}>
        {schemaArray.map((version) => (
          <MenuItem
            id="schema-menu-item"
            key={version}
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

const loadingFallbackSx = {
  textAlign: 'center',
  padding: '8rem 17rem',
} as const;

// Meshery extensions also uses this modal
function Modal({
  open,
  title,
  handleClose,
  schema,
  schema_array = [],
  type,
  schemaChangeHandler,
  handleSubmit,
  submitBtnText,
  leftHeaderIcon,
  uiSchema = {},
  helpText,
  RJSFWrapperComponent = null,
  initialData = {},
}: ModalProps) {
  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const formRef = useRef<ModalFormRef | null>(null);
  const { notify } = useNotification();
  const isLoadingSchema = useMemo(() => !schema, [schema]);

  useEffect(() => {
    const hasForbiddenDesignName = containsForbiddenDesignName(title);
    setCanNotSubmit(hasForbiddenDesignName);

    if (hasForbiddenDesignName) {
      notify({
        event_type: EVENT_TYPES.WARNING,
        message: 'Design name should not contain Untitled Design, Untitled, LFX',
      });
    }
  }, [notify, title]);

  const handleFormSubmit = () => {
    if (formRef.current?.validateForm?.()) {
      handleClose();
      handleSubmit(formRef.current.state?.formData ?? {});
    }
  };

  return (
    <SistentModal open={open} closeModal={handleClose} title={title} headerIcon={leftHeaderIcon}>
      <Typography variant="h5">
        {schema_array.length > 1 ? (
          <SchemaVersion
            schemaArray={schema_array}
            type={type}
            schemaChangeHandler={schemaChangeHandler}
          />
        ) : null}
      </Typography>
      <ModalBody>
        {isLoadingSchema ? (
          <div style={loadingFallbackSx}>
            <CircularProgress />
          </div>
        ) : (
          <RJSFWrapper
            key={type}
            formData={initialData}
            jsonSchema={schema || getSchema(type)}
            uiSchema={uiSchema}
            liveValidate={false}
            formRef={formRef}
            hideTitle={true}
            {...(RJSFWrapperComponent ? { RJSFWrapperComponent } : {})}
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
  );
}

export default React.memo(Modal);

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
}: RJSFModalWrapperProps) {
  const formRef = useRef<ModalFormRef | null>(null);
  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const { notify } = useNotification();
  const isLoadingSchema = useMemo(() => !schema, [schema]);

  useEffect(() => {
    const hasForbiddenDesignName = containsForbiddenDesignName(title);
    setCanNotSubmit(hasForbiddenDesignName);

    if (hasForbiddenDesignName) {
      notify({
        event_type: EVENT_TYPES.WARNING,
        message: 'Design name should not contain Untitled Design, Untitled, LFX',
      });
    }
  }, [notify, title]);

  const handleFormSubmit = () => {
    if (formRef.current?.validateForm?.()) {
      handleSubmit(formRef.current.state?.formData ?? {});

      if (handleNext) {
        handleNext();
      }
    }
  };

  return (
    <>
      <ModalBody>
        {isLoadingSchema ? (
          <div style={loadingFallbackSx}>
            <CircularProgress />
          </div>
        ) : (
          <RJSFWrapper
            formData={initialData}
            jsonSchema={schema}
            uiSchema={uiSchema}
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
