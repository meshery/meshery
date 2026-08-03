import type { ImportDesignApiArg } from '@meshery/schemas/mesheryApi';
import { resolveImportedDesignFile } from './import-design-file';

export type ImportDesignFormData = {
  uploadType: string;
  name: string;
  url?: string;
  file?: string;
};

export type ImportDesignRequestResult = { requestBody: string } | { errorMessage: string };

/**
 * Wire contract for `POST /api/pattern/import`, sourced from the schemas
 * generated client (`ImportDesignApiArg`, produced from
 * `MesheryPatternImportRequestBody`) rather than re-declared here. Anchoring the
 * request shape to the generated type keeps the wire field names - notably
 * `fileName` - locked to the API contract so they cannot silently drift back to
 * snake_case (the bug fixed in meshery/meshery#21105).
 */
type ImportDesignRequestBody = ImportDesignApiArg['body'];
type ImportDesignFileVariant = Extract<ImportDesignRequestBody, { fileName: string }>;
type ImportDesignUrlVariant = Extract<ImportDesignRequestBody, { url: string }>;

/**
 * File-Upload variant as this client emits it. The fields are pinned to the
 * contract via `ImportDesignFileVariant`; only `file` is widened: the openapi
 * contract types it as a base64 string (`format: byte`), while the browser
 * sends the decoded byte array produced after reading the upload. The server's
 * Go `[]byte` decoder accepts either representation, so the byte array is a
 * conformant alternative and is preserved here unchanged.
 */
type ImportDesignFileWireBody = Omit<ImportDesignFileVariant, 'file'> & { file: number[] };

export const buildImportDesignRequestBody = async (
  data: ImportDesignFormData,
): Promise<ImportDesignRequestResult> => {
  const { uploadType, name, url, file } = data;

  switch (uploadType) {
    case 'File Upload': {
      try {
        const importedFile = await resolveImportedDesignFile(file);
        if (!importedFile) {
          return { errorMessage: 'Please choose a design file before continuing.' };
        }

        const requestBody: ImportDesignFileWireBody = {
          name,
          fileName: importedFile.fileName,
          file: importedFile.fileData,
        };
        return { requestBody: JSON.stringify(requestBody) };
      } catch (error) {
        console.error('Error resolving design import file:', error);
        return { errorMessage: 'Unable to read the selected design file. Please try again.' };
      }
    }
    case 'URL Import': {
      // The "URL Import" branch is only reached with a URL-bearing form; the
      // form model types `url` as optional, so narrow it to the contract's
      // required `url` without altering the emitted wire value.
      const requestBody: ImportDesignUrlVariant = { url: url as string, name };
      return { requestBody: JSON.stringify(requestBody) };
    }
    default:
      return { errorMessage: 'Please choose a valid design import source before continuing.' };
  }
};
