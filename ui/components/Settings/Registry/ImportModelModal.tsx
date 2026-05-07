import {
  FormControlLabel,
  Button,
  FormControl,
  RadioGroup,
  Radio,
  Typography,
  ModalFooter,
  Box,
  Modal,
  useTheme,
} from '@sistent/sistent';
import { ModelImportRjsfSchemaV1Beta2 } from '@meshery/schemas';
import { RJSFModalWrapper } from '@/components/General/Modals/Modal';
import CsvStepper from './Stepper/CSVStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import React, { useState, useEffect, useContext } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/DesignLifeCycle/common';
import { NotificationCenterContext } from '@/components/NotificationCenter';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '@/components/NotificationCenter/formatters/model_registration';
import { StyledDocsRedirectLink } from './Stepper/style';
import { updateProgress } from '@/store/slices/mesheryUi';

const FinishDeploymentStep = ({
  deploymentType,
  handleClose,
}: {
  deploymentType: string;
  handleClose: () => void;
}) => {
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [isDeploying, setIsDeploying] = useState(true);
  const [deployEvent, setDeployEvent] = useState<any>();

  useEffect(() => {
    const subscription = operationsCenterActorRef.on(
      OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      (event) => {
        const serverEvent = event.data.event;
        if (serverEvent.action === deploymentType) {
          setIsDeploying(false);
          setDeployEvent(serverEvent);
        }
      },
    );

    return () => subscription.unsubscribe();
  }, []);

  const progressMessage = `${capitalize(deploymentType)}ing model`;

  const theme = useTheme();
  return (
    <>
      <Box
        sx={{
          padding: '1rem',
          overflow: 'hidden',
          backgroundColor: theme.palette.background.surfaces,
        }}
      >
        {isDeploying ? (
          <Box style={{ padding: '1rem' }}>
            <Loading message={progressMessage} />
          </Box>
        ) : (
          <>
            <ModelImportMessages message={deployEvent?.metadata?.ModelImportMessage} />
            <ModelImportedSection modelDetails={deployEvent?.metadata?.ModelDetails} />
          </>
        )}
      </Box>
      <ModalFooter
        variant="filled"
        helpText={
          'Click Finish to complete the model import process. The imported model will be available in your Model Registry.'
        }
      >
        <Button variant="contained" color="primary" onClick={handleClose}>
          Finish
        </Button>
      </ModalFooter>
    </>
  );
};

// Canonical RJSF form schemas authored in `meshery/schemas` and validated
// by its `validation/forms_test.go` to be a strict subset of the OpenAPI
// `MesheryModelImportFormPayload` construct. The canonical is the source
// of truth for property shapes, types, descriptions, enum tokens, and
// upload-type labels — pull them through directly.
//
// The canonical describes the full form-data shape (file + url + csv
// branches), but this modal only renders the file-and-url subset; the
// CSV branch is handled by a separate `CsvStepper` modal opened on the
// `csv` upload-type. Keeping the unused fields in the rendered form
// caused two real problems:
//   (a) RJSF v6's @rjsf/validator-ajv8 rejected the form on `Next` —
//       even though `modelCsv` / `componentCsv` / `relationshipCsv` were
//       hidden via `ui:widget: 'hidden'`, their `format: "binary"` and
//       the canonical's csv-branch `allOf` required-list still blocked
//       `validateForm()`, so the import POST never fired and the test
//       hung waiting for the registration event.
//   (b) the canonical's schema-level `enumNames` is silently ignored
//       by @rjsf/utils' `optionsList()` — RJSF v6 reads the radio
//       labels from `ui:enumNames` on the uiSchema instead.
//
// Hold the consumer-shaped subset locally with content pulled from the
// canonical (titles, descriptions, enum tokens, friendly labels). When
// the upstream schema gains a UI-friendly variant or RJSF starts
// honoring schema-level `enumNames`, this override can shrink.
const canonicalProps = ModelImportRjsfSchemaV1Beta2.properties;
const canonicalUploadType = canonicalProps.uploadType;
const UPLOAD_TYPE_FILE = 'file';
const UPLOAD_TYPE_URL = 'urlImport';
const UPLOAD_TYPE_CSV = 'csv';

const importModelSchema = {
  $schema: 'https://json-schema.org/draft-07/schema#',
  title: ModelImportRjsfSchemaV1Beta2.title,
  type: 'object',
  properties: {
    uploadType: canonicalUploadType,
    modelFile: canonicalProps.modelFile,
    url: canonicalProps.url,
  },
  allOf: [
    {
      if: {
        properties: { uploadType: { const: UPLOAD_TYPE_FILE } },
        required: ['uploadType'],
      },
      then: { required: ['modelFile'] },
    },
    {
      if: {
        properties: { uploadType: { const: UPLOAD_TYPE_URL } },
        required: ['uploadType'],
      },
      then: { required: ['url'] },
    },
  ],
  required: ['uploadType'],
};

const importModelUiSchema = {
  uploadType: {
    'ui:widget': 'radio',
    'ui:enumNames': [...(canonicalUploadType.enumNames as string[])],
  },
  modelFile: { 'ui:widget': 'file' },
  'ui:order': ['uploadType', 'modelFile', 'url'],
};

// RJSF's file widget encodes the original filename inside the data URL
// (`data:<mime>;name=<file>;base64,<payload>`). Pull it back out so the
// server-side ImportBody.FileName is populated even though the form
// doesn't surface a separate `fileName` input.
const filenameFromDataUrl = (dataUrl: string | undefined): string | undefined => {
  if (!dataUrl) return undefined;
  const match = dataUrl.match(/;name=([^;]+);/);
  return match ? decodeURIComponent(match[1]) : undefined;
};

type ImportModelModalProps = {
  isImportModalOpen: boolean;
  setIsImportModalOpen: (_open: boolean) => void;
};

const ImportModelModal = React.memo(
  ({ isImportModalOpen, setIsImportModalOpen }: ImportModelModalProps) => {
    const [importModalDescription, setImportModalDescription] = useState('');
    const [isCsvModalOpen, setIsCsvModalOpen] = useState(false);
    const [importModelReq] = useImportMeshModelMutation();
    const [activeStep, setActiveStep] = useState(0);

    const handleNext = () => {
      if (activeStep === 0) {
        setActiveStep(1);
      }
    };

    const handleClose = () => {
      setIsImportModalOpen(false);
      setIsCsvModalOpen(false);
      setActiveStep(0);
    };

    const handleImportModelSubmit = async (data) => {
      // Canonical field names from ModelImportRjsfSchemaV1Beta2:
      // `uploadType` (enum: file/urlImport/csv), `modelFile`, `fileName`,
      // `url`, plus the CSV trio handled by CsvStepper.
      const { uploadType, url, modelFile, fileName: formFileName } = data;
      let requestBody = null;

      // RJSF's file widget encodes the original filename in the data URL,
      // but we also fall back to the input element's `files[0].name` in
      // case the widget shape changes upstream.
      const fileElement = document.getElementById('root_modelFile') as HTMLInputElement | null;

      switch (uploadType) {
        case UPLOAD_TYPE_FILE: {
          const fileName =
            formFileName || filenameFromDataUrl(modelFile) || fileElement?.files?.[0]?.name;
          const fileData = modelFile ? getUnit8ArrayDecodedFile(modelFile) : null;
          if (fileData && fileName) {
            // Server's ImportBody.FileName / .ModelFile are tagged
            // `file_name` / `model_file` (snake_case wire format).
            requestBody = {
              importBody: {
                model_file: fileData,
                url: '',
                file_name: fileName,
              },
              uploadType: UPLOAD_TYPE_FILE,
              register: true,
            };
          } else {
            console.error('Error: File data or file name is empty or invalid');
            return;
          }
          break;
        }
        case UPLOAD_TYPE_URL: {
          if (url) {
            requestBody = {
              importBody: {
                url: url,
              },
              uploadType: UPLOAD_TYPE_URL,
              register: true,
            };
          } else {
            console.error('Error: URL is empty');
            return;
          }
          break;
        }
        case UPLOAD_TYPE_CSV: {
          handleClose();
          setIsCsvModalOpen(true);
          return;
        }
        default: {
          console.error('Error: Invalid upload type');
          return;
        }
      }
      updateProgress({ showProgress: true });
      await importModelReq({ importBody: requestBody });
      updateProgress({ showProgress: false });
    };

    const CustomRadioWidget = (props: {
      options: any;
      value: string;
      onChange: (_value: string) => void;
      label: string;
      schema: any;
    }) => {
      const { options, value, onChange, label, schema } = props;
      const { enumOptions } = options;

      setImportModalDescription(schema.description);

      return (
        <FormControl component="fieldset">
          <RadioGroup
            value={value}
            onChange={(e) => onChange(e.target.value)}
            style={{ marginTop: '-1.7rem', marginLeft: '-1rem' }}
          >
            <Typography fontWeight={'bold'} fontSize={'1rem'}>
              {label}
            </Typography>

            {enumOptions.map((option, index) => (
              <FormControlLabel
                key={option.value}
                value={option.value}
                control={<Radio />}
                label={
                  <div>
                    <Typography variant="subtitle1">{option.label}</Typography>
                    <Typography variant="body2" color="textSecondary" textTransform={'none'}>
                      {schema.enumDescriptions?.[index]}
                    </Typography>
                  </div>
                }
              />
            ))}
          </RadioGroup>
        </FormControl>
      );
    };

    const widgets = {
      RadioWidget: CustomRadioWidget,
    };

    return (
      <>
        <Modal
          open={isImportModalOpen}
          closeModal={handleClose}
          maxWidth="sm"
          title="Import Model"
          style={{
            zIndex: 1500,
          }}
        >
          {activeStep === 0 ? (
            <RJSFModalWrapper
              schema={importModelSchema}
              uiSchema={importModelUiSchema}
              handleSubmit={handleImportModelSubmit}
              handleNext={handleNext}
              submitBtnText="Next"
              handleClose={handleClose}
              widgets={widgets}
              helpText={
                <p>
                  {importModalDescription} <br />
                  Learn more about importing Models in our{' '}
                  <StyledDocsRedirectLink
                    href={`${MESHERY_DOCS_URL}/guides/configuration-management/importing-models`}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    documentation
                  </StyledDocsRedirectLink>
                  .
                </p>
              }
            />
          ) : (
            <FinishDeploymentStep deploymentType="register" handleClose={handleClose} />
          )}
        </Modal>
        <Modal
          open={isCsvModalOpen}
          closeModal={() => setIsCsvModalOpen(false)}
          maxWidth="sm"
          title="Import CSV"
          style={{
            zIndex: 1500,
          }}
        >
          <CsvStepper handleClose={handleClose} />
        </Modal>
      </>
    );
  },
);

ImportModelModal.displayName = 'ImportModelModal';

export default ImportModelModal;
