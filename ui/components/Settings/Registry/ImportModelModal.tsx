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
import { ModelDefinitionV1Beta1OpenApiSchema } from '@meshery/schemas';
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

// Build the Import Model RJSF schema from the canonical
// `@meshery/schemas` OpenAPI definition rather than from sistent's
// derived `importModelSchema` export. Sistent's derivation (since
// v0.21.3) renames `file` → `modelFile` and emits the `uploadType` enum
// without `enumNames`, which doesn't round-trip through this modal —
// `getUnit8ArrayDecodedFile(file)` and the `'File Import'`/`'URL Import'`/
// `'CSV Import'` switch keys can't both move with the upstream rename
// without churn here, and the e2e suite searches for the friendly
// "File Import" heading to begin with.
//
// We pull the canonical strings (title, description, enum tokens,
// enumDescriptions) straight from `ImportRequest.properties.uploadType`
// so wording stays in lockstep with `meshery/schemas`. The friendly
// radio labels come from the canonical `ImportBody.oneOf[i].title` —
// those titles are what RJSF would have surfaced if the schema were
// rendered as a discriminated union, so nothing is being invented
// here. Only the RJSF-specific shape (top-level `file`/`url` fields,
// widget mapping, conditional `required`) is added on top.
const canonicalImportRequest = ModelDefinitionV1Beta1OpenApiSchema.components.schemas.ImportRequest;
const canonicalImportBody = canonicalImportRequest.properties.importBody;
const canonicalUploadType = canonicalImportRequest.properties.uploadType;
const [fileImportBranch, urlImportBranch, csvImportBranch] = canonicalImportBody.oneOf;

// Friendly labels — what the user sees in the radio group and what the
// e2e suite searches for via `getByRole('heading', { name: ... })`.
const FILE_IMPORT_LABEL = fileImportBranch.title;
const URL_IMPORT_LABEL = urlImportBranch.title;
const CSV_IMPORT_LABEL = csvImportBranch.title;
const UPLOAD_TYPE_OPTIONS = [FILE_IMPORT_LABEL, URL_IMPORT_LABEL, CSV_IMPORT_LABEL];

const importModelSchema = {
  $schema: 'https://json-schema.org/draft-07/schema#',
  title: 'Import Model',
  type: 'object',
  properties: {
    uploadType: {
      type: 'string',
      title: canonicalUploadType.title,
      // The canonical enum is the lowercase tokens (`file`/`urlImport`/...).
      // The form keys off the friendly `oneOf` titles instead so the
      // radio labels and the consumer switch share a single vocabulary
      // and the e2e selectors keep working against unchanged copy.
      enum: UPLOAD_TYPE_OPTIONS,
      enumDescriptions: canonicalUploadType.enumDescriptions,
      description: canonicalUploadType.description,
      'x-rjsf-grid-area': '12',
    },
    file: {
      type: 'string',
      title: 'Model file',
      description: fileImportBranch.properties.modelFile.description,
      'x-rjsf-grid-area': '12',
    },
    url: {
      type: 'string',
      format: 'uri',
      title: 'URL',
      description: urlImportBranch.properties.url.description,
      'x-rjsf-grid-area': '12',
    },
  },
  allOf: [
    {
      // The `required: ['uploadType']` inside each `if` ensures the branch
      // only triggers when uploadType is present — without it, JSON Schema's
      // `properties` keyword silently passes when the property is absent, so
      // both `file` and `url` would be marked required on first render.
      if: {
        properties: { uploadType: { const: FILE_IMPORT_LABEL } },
        required: ['uploadType'],
      },
      then: { required: ['file'] },
    },
    {
      if: {
        properties: { uploadType: { const: URL_IMPORT_LABEL } },
        required: ['uploadType'],
      },
      then: { required: ['url'] },
    },
  ],
  required: ['uploadType'],
};

const importModelUiSchema = {
  uploadType: { 'ui:widget': 'radio' },
  file: { 'ui:widget': 'file' },
  'ui:order': ['uploadType', 'file', 'url'],
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
      const { uploadType, url, file } = data;
      let requestBody = null;

      const fileElement = document.getElementById('root_file') as HTMLInputElement | null;

      switch (uploadType) {
        case FILE_IMPORT_LABEL: {
          const fileName = fileElement?.files?.[0]?.name;
          const fileData = file ? getUnit8ArrayDecodedFile(file) : null;
          if (fileData && fileName) {
            // Server expects ImportBody.FileName (json:"file_name"); the legacy
            // `filename` key was silently dropped, leaving the temp file with
            // no name and breaking model registration.
            requestBody = {
              importBody: {
                model_file: fileData,
                url: '',
                file_name: fileName,
              },
              uploadType: 'file',
              register: true,
            };
          } else {
            console.error('Error: File data or file name is empty or invalid');
            return;
          }
          break;
        }
        case URL_IMPORT_LABEL: {
          if (url) {
            requestBody = {
              importBody: {
                url: url,
              },
              uploadType: 'urlImport',
              register: true,
            };
          } else {
            console.error('Error: URL is empty');
            return;
          }
          break;
        }
        case CSV_IMPORT_LABEL: {
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
