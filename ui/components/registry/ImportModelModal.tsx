/**
 * Import-model modal.
 *
 * Two-step flow:
 *
 *   1. RJSF form driven by the canonical `ModelImportRjsfSchemaV1Beta2`
 *      schema. The upload-type selector (file / url / csv) routes the submit
 *      handler to the right server payload. CSV picks a separate stepper
 *      modal instead of progressing.
 *   2. Finish screen that subscribes to the operations-center event bus and
 *      renders the import progress + result.
 *
 * Built atop the shared `FormModal` and `Modal` primitives. Migrated as part
 * of Phase 5.b.6 (#18754); previously consumed the legacy `RJSFModalWrapper`.
 *
 * Schema setup + DOM file-input helpers live in `./importModelSchema` to
 * keep this file under the 400-line budget — see that module for the
 * v1.2.16 canonical branch layout and the fixes required to keep RJSF v6 /
 * @rjsf/validator-ajv8 happy with hidden CSV fields and the empty-URL guard.
 */
import {
  FormControlLabel,
  FormControl,
  RadioGroup,
  Radio,
  Typography,
  Box,
  ModalButtonPrimary,
} from '@sistent/sistent';
import { useTheme } from '@/theme';
import { Modal, FormModal } from '@/components/shared/Modal';
import CsvStepper from './Stepper/CSVStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import { useState, useEffect, useContext, memo } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/designs/lifecycle/common';
import { NotificationCenterContext } from '@/components/layout/NotificationCenter';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '@/components/layout/NotificationCenter/formatters/model_registration';
import { StyledDocsRedirectLink } from './Stepper/style';
import { updateProgress } from '@/store/slices/mesheryUi';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { readFileAsBytes } from '@/utils/fileUpload';
import {
  UPLOAD_TYPE_CSV,
  UPLOAD_TYPE_FILE,
  UPLOAD_TYPE_URL,
  decodeDataUrlToBytes,
  filenameFromDataUrl,
  findSelectedModelFile,
  importModelSchema,
  importModelUiSchema,
} from './importModelSchema';

const FinishDeploymentStep = ({ deploymentType }: { deploymentType: string }) => {
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [isDeploying, setIsDeploying] = useState(true);
  const [deployEvent, setDeployEvent] = useState<any>();
  const theme = useTheme();

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

  return (
    <Box
      sx={{
        padding: '1rem',
        overflow: 'hidden',
        backgroundColor: theme.palette.background.surfaces,
      }}
    >
      {isDeploying ? (
        <Box sx={{ padding: '1rem' }}>
          <Loading message={progressMessage} />
        </Box>
      ) : (
        <>
          <ModelImportMessages message={deployEvent?.metadata?.ModelImportMessage} />
          <ModelImportedSection modelDetails={deployEvent?.metadata?.ModelDetails} />
        </>
      )}
    </Box>
  );
};

interface CustomRadioWidgetProps {
  options: { enumOptions: Array<{ value: string; label: string }> };
  value: string;
  onChange: (_value: string) => void;
  label: string;
  schema: { description?: string; enumDescriptions?: string[] };
}

// Hoisted to module scope so the widget identity is stable across renders of
// `ImportModelModal` — defining it inside the parent would unmount/remount
// the radio group on every parent render.
const CustomRadioWidget = (props: CustomRadioWidgetProps) => {
  const { options, value, onChange, label, schema } = props;
  const { enumOptions } = options;

  return (
    <FormControl component="fieldset">
      <RadioGroup
        value={value}
        onChange={(e) => onChange(e.target.value)}
        sx={{ marginTop: '-1.7rem', marginLeft: '-1rem' }}
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

const widgets = { RadioWidget: CustomRadioWidget };

// Pull the upload-type description statically from the canonical schema so
// the footer help text doesn't need state derived from a render-time side
// effect inside the widget.
const importModalDescription =
  ((importModelSchema.properties as Record<string, { description?: string }> | undefined)?.[
    'uploadType'
  ]?.description as string | undefined) ?? '';

interface ImportModelModalProps {
  isImportModalOpen: boolean;
  setIsImportModalOpen: (_open: boolean) => void;
}

const ImportModelModal = memo<ImportModelModalProps>(
  ({ isImportModalOpen, setIsImportModalOpen }) => {
    const [isCsvModalOpen, setIsCsvModalOpen] = useState(false);
    const [importModelReq] = useImportMeshModelMutation();
    const [activeStep, setActiveStep] = useState(0);
    const [pendingRequest, setPendingRequest] = useState<any>(null);
    const { notify } = useNotification();

    const handleClose = () => {
      setIsImportModalOpen(false);
      setIsCsvModalOpen(false);
      setActiveStep(0);
    };

    const handleImportModelSubmit = async (data: any) => {
      // Canonical field names from ModelImportRjsfSchemaV1Beta2:
      // `uploadType` (enum: file/urlImport/csv), `modelFile`, `fileName`,
      // `url`, plus the CSV trio handled by CsvStepper.
      const { uploadType, url, modelFile, fileName: formFileName } = data;
      let requestBody = null;

      switch (uploadType) {
        case UPLOAD_TYPE_FILE: {
          // Three source-of-truth paths for the file:
          //   1. `data.modelFile` (RJSF form data) — populated by
          //      `CustomFileWidget.processFile`'s async FileReader chain.
          //   2. `data.fileName` (RJSF form data) — NOT populated by the
          //      widget; the canonical schema declares it as a separate
          //      field but the widget only emits `modelFile`. Always
          //      `undefined` in practice.
          //   3. DOM input (synchronous) — the browser sets
          //      `<input type=file>.files[0]` the instant the user picks
          //      a file, so it's always readable on submit. The browser-
          //      provided `File.name` is also our only synchronous source
          //      for the filename (the data URL produced by
          //      `readAsDataURL` does NOT embed `;name=` like the
          //      reference RJSF FileWidget would).
          //
          // So: try the form-state bytes (path 1) first, fall back to DOM
          // bytes if the FileReader race lost (paths 1 -> 3 for bytes), and
          // ALWAYS prefer the DOM filename for `fileName` (path 3 for
          // name) — paths 1 and 2 just don't carry it.
          let fileData: number[] | null = decodeDataUrlToBytes(modelFile);
          let fileName: string | undefined = formFileName || filenameFromDataUrl(modelFile);
          if (!fileData || !fileName) {
            const inputFile = findSelectedModelFile();
            if (inputFile) {
              try {
                if (!fileData) {
                  fileData = await readFileAsBytes(inputFile);
                }
                if (!fileName) {
                  fileName = inputFile.name;
                }
              } catch (err) {
                console.error('Error reading file from DOM:', err);
                notify({
                  message: `Unable to read the selected file. Please try again.`,
                  event_type: EVENT_TYPES.ERROR,
                });
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
            notify({
              message: 'Please choose a model file before continuing.',
              event_type: EVENT_TYPES.ERROR,
            });
            return;
          }
          break;
        }
        case UPLOAD_TYPE_URL: {
          if (url) {
            requestBody = {
              importBody: { url },
              uploadType: UPLOAD_TYPE_URL,
              register: true,
            };
          } else {
            notify({
              message: 'Please provide a model URL before continuing.',
              event_type: EVENT_TYPES.ERROR,
            });
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
          notify({
            message: 'Please choose an import method to continue.',
            event_type: EVENT_TYPES.ERROR,
          });
          return;
        }
      }

      // Fire the request asynchronously but transition to the Finish step
      // IMMEDIATELY so the operations-center listener mounts before the
      // fast backend fires the WebSocket success event.
      // If it fails, transition back to the form step so the user is not stranded.
      updateProgress({ showProgress: true });
      setPendingRequest(requestBody);
      setActiveStep(1);
    };

    useEffect(() => {
      if (activeStep === 1 && pendingRequest) {
        (async () => {
          try {
            await importModelReq({ importBody: pendingRequest }).unwrap();
          } catch (err) {
            console.error('Failed to import model:', err);
            notify({
              message: 'Model import failed. Please verify the file or URL and try again.',
              event_type: EVENT_TYPES.ERROR,
            });
            setActiveStep(0); // Move back on failure
          } finally {
            updateProgress({ showProgress: false });
            setPendingRequest(null);
          }
        })();
      }
    }, [activeStep, pendingRequest, importModelReq, notify]);

    const helpText = (
      <>
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
      </>
    );

    return (
      <>
        {activeStep === 0 ? (
          <FormModal
            isOpen={isImportModalOpen}
            onClose={handleClose}
            title="Import Model"
            size="sm"
            schema={importModelSchema}
            uiSchema={importModelUiSchema}
            widgets={widgets}
            onSubmit={handleImportModelSubmit}
            submitText="Next"
            helpText={helpText}
            sx={{ zIndex: 1600 }}
          />
        ) : (
          <Modal
            isOpen={isImportModalOpen}
            onClose={handleClose}
            title="Import Model"
            size="sm"
            helpText="Click Finish to complete the model import process. The imported model will be available in your Model Registry."
            actions={<ModalButtonPrimary onClick={handleClose}>Finish</ModalButtonPrimary>}
            sx={{ zIndex: 1600 }}
          >
            <FinishDeploymentStep deploymentType="register" />
          </Modal>
        )}

        <Modal
          isOpen={isCsvModalOpen}
          onClose={() => setIsCsvModalOpen(false)}
          title="Import CSV"
          size="sm"
          disableBodyWrap
          sx={{ zIndex: 1600 }}
        >
          <CsvStepper handleClose={handleClose} />
        </Modal>
      </>
    );
  },
);

ImportModelModal.displayName = 'ImportModelModal';

export default ImportModelModal;
