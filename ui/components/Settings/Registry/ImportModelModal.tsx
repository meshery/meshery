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
import React, { useState, useEffect, useContext, useRef } from 'react';
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

const canonicalProps = ModelImportRjsfSchemaV1Beta2?.properties || {};
const canonicalUploadType = canonicalProps.uploadType || {};
const UPLOAD_TYPE_FILE = 'file';
const UPLOAD_TYPE_URL = 'urlImport';
const UPLOAD_TYPE_CSV = 'csv';

const importModelSchema = {
  $schema: 'https://json-schema.org/draft-07/schema#',
  title: ModelImportRjsfSchemaV1Beta2?.title,
  type: 'object',
  properties: {
    uploadType: canonicalUploadType,
    modelFile: canonicalProps.modelFile || {},
    url: canonicalProps.url || {},
  },
  allOf: [
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
    'ui:enumNames': Array.isArray(canonicalUploadType?.enumNames)
      ? [...(canonicalUploadType.enumNames as string[])]
      : undefined,
  },
  modelFile: { 'ui:widget': 'file' },
  'ui:order': ['uploadType', 'modelFile', 'url'],
};

const filenameFromDataUrl = (dataUrl: string | undefined): string | undefined => {
  if (!dataUrl) return undefined;
  const match = dataUrl.match(/;name=([^;]+);/);
  return match ? decodeURIComponent(match[1]) : undefined;
};

const decodeDataUrlToBytes = (dataUrl: string | undefined): number[] | null => {
  if (!dataUrl) return null;
  return getUnit8ArrayDecodedFile(dataUrl);
};

const findSelectedModelFile = (container: HTMLElement | null): File | undefined => {
  if (!container) return undefined;

  const inputs = container.querySelectorAll<HTMLInputElement>('input[type="file"][id$="_modelFile"]');
  for (const input of Array.from(inputs)) {
    const file = input.files?.[0];
    if (file) return file;
  }

  return undefined;
};

const readFileAsBytes = (file: File): Promise<number[]> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const buffer = reader.result as ArrayBuffer;
      resolve(Array.from(new Uint8Array(buffer)));
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsArrayBuffer(file);
  });

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
    const importFormContainerRef = useRef<HTMLDivElement | null>(null);

    const handleClose = () => {
      setIsImportModalOpen(false);
      setIsCsvModalOpen(false);
      setActiveStep(0);
    };

    const handleImportModelSubmit = async (data) => {
      const { uploadType, url, modelFile, fileName: formFileName } = data;
      let requestBody = null;

      switch (uploadType) {
        case UPLOAD_TYPE_FILE: {
          let fileName = formFileName || filenameFromDataUrl(modelFile);
          let fileData: number[] | null = decodeDataUrlToBytes(modelFile);

          if (!fileData) {
            const inputFile = findSelectedModelFile(importFormContainerRef.current);
            if (inputFile) {
              try {
                fileData = await readFileAsBytes(inputFile);
                fileName = fileName || inputFile.name;
              } catch (error) {
                console.error('Error reading model file:', error);
                return;
              }
            }
          }

          if (fileData && fileName) {
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

      setActiveStep(1);
      updateProgress({ showProgress: true });
      try {
        await importModelReq({ importBody: requestBody });
      } finally {
        updateProgress({ showProgress: false });
      }
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
            <Box ref={importFormContainerRef}>
              <RJSFModalWrapper
                schema={importModelSchema}
                uiSchema={importModelUiSchema}
                handleSubmit={handleImportModelSubmit}
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
            </Box>
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
