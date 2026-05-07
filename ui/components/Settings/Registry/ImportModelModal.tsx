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
import { ModelImportRjsfSchemaV1Beta2, ModelImportRjsfUiSchemaV1Beta2 } from '@meshery/schemas';
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
// `MesheryModelImportFormPayload` construct. Importing them directly
// keeps the modal locked to the canonical field shape (uploadType enum
// tokens, modelFile/fileName/url naming, allOf required branches) — any
// drift surfaces as a build/test failure here rather than as silent
// runtime divergence.
const importModelSchema = ModelImportRjsfSchemaV1Beta2;
const importModelUiSchema = ModelImportRjsfUiSchemaV1Beta2;

// Token vocabulary the canonical schema uses. Pull these from the schema
// (rather than hard-code) so adding/renaming an upload method upstream
// flows through this consumer with one constant change.
const [UPLOAD_TYPE_FILE, UPLOAD_TYPE_URL, UPLOAD_TYPE_CSV] = importModelSchema.properties.uploadType
  .enum as string[];

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

      // RJSF doesn't surface the original filename through form data, so
      // read it off the input element. The id matches the schema field name.
      const fileElement = document.getElementById('root_modelFile') as HTMLInputElement | null;

      switch (uploadType) {
        case UPLOAD_TYPE_FILE: {
          const fileName = formFileName || fileElement?.files?.[0]?.name;
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
