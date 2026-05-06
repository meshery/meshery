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

// Local schema/uiSchema for the Import Model modal.
//
// The upstream `importModelSchema`/`importModelUiSchema` exported by
// @sistent/sistent (>=0.21.3) renamed properties (`file` → `modelFile`) and
// switched the `uploadType` enum to internal tokens (`file`/`urlImport`/`csv`)
// without `enumNames`. That broke the consumer code below — which keys off the
// human-readable labels and the `file` form field — and also stripped the
// friendly "File Import" / "URL Import" / "CSV Import" radio labels relied on
// by the e2e tests. Until sistent ships a schema that round-trips through this
// modal we hold the form shape locally so the registry import flow keeps
// working end-to-end.
const importModelSchema = {
  $schema: 'https://json-schema.org/draft-07/schema#',
  title: 'Import Model',
  type: 'object',
  properties: {
    uploadType: {
      type: 'string',
      title: 'Upload method',
      enum: ['File Import', 'URL Import', 'CSV Import'],
      // CustomRadioWidget below reads `schema.enumDescriptions[index]` per
      // option; missing this would crash the modal at render.
      enumDescriptions: [
        'Upload a model file from your local system.',
        'Import a model from a remote URL.',
        'Upload a CSV file containing model definitions.',
      ],
      description:
        "Choose the method you prefer to upload your model file. Select 'File Import' or 'CSV Import' if you have the file on your local system, or 'URL Import' if you have the file hosted online.",
      'x-rjsf-grid-area': '12',
    },
    file: {
      type: 'string',
      title: 'Model file',
      description: 'Model file content. Supported formats: .tar, .tar.gz, .tgz.',
      'x-rjsf-grid-area': '12',
    },
    url: {
      type: 'string',
      format: 'uri',
      title: 'URL',
      description: 'A direct URL to a single model file. Supported formats: .tar, .tar.gz, .tgz.',
      'x-rjsf-grid-area': '12',
    },
  },
  allOf: [
    {
      if: { properties: { uploadType: { const: 'File Import' } } },
      then: { required: ['file'] },
    },
    {
      if: { properties: { uploadType: { const: 'URL Import' } } },
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
        case 'File Import': {
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
        case 'URL Import': {
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
        case 'CSV Import': {
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
