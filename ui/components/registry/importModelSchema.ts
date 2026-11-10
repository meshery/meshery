/**
 * Schema + helpers for the registry's Import Model modal.
 *
 * Split out of `ImportModelModal.tsx` to keep the rendering file under the
 * 400-line budget. See the inline comments below for the canonical schema
 * shape and the rationale behind each consumer-specific override.
 */
import { ModelImportRjsfSchemaV1Beta2, ModelImportRjsfUiSchemaV1Beta2 } from '@meshery/schemas';
import { findSelectedFileInDialog } from '@/utils/fileUpload';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';

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
  then?: { required?: string[] } & Record<string, unknown>;
};

const canonicalUploadType =
  ((ModelImportRjsfSchemaV1Beta2 as unknown as RJSFNode).properties?.uploadType as
    | { enumNames?: string[] }
    | undefined) ?? {};

export const UPLOAD_TYPE_FILE = 'file';
export const UPLOAD_TYPE_URL = 'urlImport';
export const UPLOAD_TYPE_CSV = 'csv';

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
  const { required: _unused, ...thenRest } = branch.then;
  return { ...branch, then: thenRest };
};

export const importModelSchema = {
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
export const importModelUiSchema = {
  ...ModelImportRjsfUiSchemaV1Beta2,
  uploadType: {
    'ui:widget': 'radio',
    'ui:enumNames': Array.isArray(canonicalUploadType?.enumNames)
      ? [...(canonicalUploadType.enumNames as string[])]
      : undefined,
  },
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
export const filenameFromDataUrl = (dataUrl: string | undefined): string | undefined => {
  if (!dataUrl) return undefined;
  const match = dataUrl.match(/;name=([^;]+);/);
  return match ? decodeURIComponent(match[1]) : undefined;
};

// Decode an RJSF data URL into the byte array the server expects.
export const decodeDataUrlToBytes = (dataUrl: string | undefined): number[] | null => {
  if (!dataUrl) return null;
  return getUnit8ArrayDecodedFile(dataUrl);
};

// Locate the `<input type="file">` rendered for the modelFile field. RJSF
// uses the schema title as the id prefix (`<title>_modelFile`), so we
// can't hard-code `root_modelFile`. Scope the query to the active modal
// dialog to avoid false matches from other file inputs on the page.
export const findSelectedModelFile = (): File | undefined => {
  return findSelectedFileInDialog('input[type="file"][id$="_modelFile"]');
};
