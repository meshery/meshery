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
import { RJSFModalWrapper } from '@/components/shared/Modal/Modal';
import CsvStepper from './Stepper/CSVStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import React, { useState, useEffect, useContext, useRef } from 'react';
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
// Since @meshery/schemas v1.2.16, the canonical
// ModelImportRjsfSchemaV1Beta2 routes its `modelFile`, `fileName`,
// `url`, and CSV-trio fields through conditional `allOf[].then`
// branches keyed by the `uploadType` discriminator — they are NO
// LONGER at the schema root.
//
// We reuse the canonical's conditional shape directly so each upload
// branch only renders ITS OWN fields:
//   - File Import: `modelFile` (file widget) + hidden `fileName`
//   - URL Import:  `url` (text widget)
//   - CSV Import:  no fields rendered (CSV opens a separate stepper modal)
//
// Flattening the canonical branches back into a single root-level
// property map — the pre-v1.2.16 shape — was a tempting shortcut, but
// it makes the `url` field always render, and its `format: "uri"`
// constraint then rejects the empty default on every Next-click for
// the File Import branch, leaving validateForm() returning false and
// the submit button silently dead.
type RJSFNode = {
  properties?: Record<string, unknown>;
  required?: string[];

  allOf?: any[];
};
const canonicalUploadType =
  ((ModelImportRjsfSchemaV1Beta2 as unknown as RJSFNode).properties?.uploadType as
    | { enumNames?: string[] }
    | undefined) ?? {};
const UPLOAD_TYPE_FILE = 'file';
const UPLOAD_TYPE_URL = 'urlImport';
const UPLOAD_TYPE_CSV = 'csv';

// Reuse the canonical's discriminator-conditional shape, but drop the
// `required:` clauses on each branch. Two reasons:
//   - File branch: `modelFile` is filled asynchronously by RJSF's
//     FileWidget (FileReader → Promise.then(onChange)). Synchronous
//     Next-click validation loses the race; we trust the user-flow
//     guard in `handleImportModelSubmit` instead, which reads the
//     selected file off the DOM input as a synchronous fallback.
//   - CSV branch: this modal doesn't render the CSV fields (those
//     live in `CsvStepper`), so requiring them would block the submit
//     handler from short-circuiting into the CSV-modal-open flow.
const stripRequired = (branch: RJSFNode): RJSFNode => {
  if (!branch?.then) return branch;

  const { required: _unused, ...thenRest } = branch.then as RJSFNode;
  return { ...branch, then: thenRest };
};
const importModelSchema = {
  title: ModelImportRjsfSchemaV1Beta2.title,
  type: 'object',
  properties: (ModelImportRjsfSchemaV1Beta2 as unknown as RJSFNode).properties,
  allOf: ((ModelImportRjsfSchemaV1Beta2 as unknown as RJSFNode).allOf ?? []).map(stripRequired),
  required: ['uploadType'],
};

// Extend the canonical uiSchema (added in @meshery/schemas v1.2.16) with
// consumer-specific overrides:
//   - uploadType: add ui:enumNames so RJSF v6 renders friendly radio labels
//     (optionsList() reads labels from uiSchema, not schema-level enumNames).
//   - fileName: hidden — derived from the uploaded file at submit time.
//   - CSV fields: not rendered in this modal (handled by CsvStepper).
const importModelUiSchema = {
  ...ModelImportRjsfUiSchemaV1Beta2,
  uploadType: {
    'ui:widget': 'radio',
    'ui:enumNames': Array.isArray(canonicalUploadType?.enumNames)
      ? [...(canonicalUploadType.enumNames as string[])]
      : undefined,
  },
  // The canonical schema's `modelFile` already declares
  // `format: "data-url"` (since v1.2.16), so RJSF auto-routes to the
  // FileWidget. The explicit override here is belt-and-braces in case
  // the canonical or a downstream override regresses.
  modelFile: { 'ui:widget': 'file' },
  fileName: { 'ui:widget': 'hidden' },
  modelCsv: { 'ui:widget': 'hidden' },
  componentCsv: { 'ui:widget': 'hidden' },
  relationshipCsv: { 'ui:widget': 'hidden' },
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

// Decode an RJSF data URL into the byte array the server expects.
// `getUnit8ArrayDecodedFile` is synchronous, so this wrapper is too;
// returns `null` when no data URL is present so callers can fall back
// to the DOM input.
const decodeDataUrlToBytes = (dataUrl: string | undefined): number[] | null => {
  if (!dataUrl) return null;
  return getUnit8ArrayDecodedFile(dataUrl);
};

// Locate the `<input type="file">` rendered for the modelFile field. RJSF
// uses the schema title as the id prefix (`<title>_modelFile`), so we
// can't hard-code `root_modelFile`. Scope the query to the active modal
// dialog to avoid false matches from other file inputs on the page.
const findSelectedModelFile = (): File | undefined => {
  if (typeof document === 'undefined') return undefined;
  const dialog = document.querySelector<HTMLElement>('[role="dialog"]');
  const root = dialog ?? document;
  const inputs = root.querySelectorAll<HTMLInputElement>('input[type="file"][id$="_modelFile"]');
  for (const input of Array.from(inputs)) {
    const file = input.files?.[0];
    if (file) return file;
  }
  return undefined;
};

// Read a `File` synchronously-from-the-DOM-but-asynchronously-into-bytes,
// matching the wire shape `getUnit8ArrayDecodedFile` produces from a
// data URL. Used as the synchronous fallback when RJSF's form data
// hasn't received the data URL yet (the FileReader chain inside
// CustomFileWidget can lose the race against a quick Next-click).
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
    // Prevents RJSFModalWrapper from advancing the step when the submit
    // handler short-circuits (CSV branch, no-file guard, empty-URL guard).
    const skipNextRef = useRef(false);

    const handleNext = () => {
      if (skipNextRef.current) {
        skipNextRef.current = false;
        return;
      }
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
          // bytes if the FileReader race lost (paths 1 ↘ 3 for bytes), and
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
                skipNextRef.current = true;
                return;
              }
            }
          }
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
            skipNextRef.current = true;
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
            skipNextRef.current = true;
            return;
          }
          break;
        }
        case UPLOAD_TYPE_CSV: {
          skipNextRef.current = true;
          handleClose();
          setIsCsvModalOpen(true);
          return;
        }
        default: {
          console.error('Error: Invalid upload type');
          skipNextRef.current = true;
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
