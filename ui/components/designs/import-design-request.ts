import type { ImportDesignApiArg } from '@meshery/schemas/mesheryApi';
import { resolveImportedDesignFile } from './import-design-file';

export type ImportDesignFormData = {
  uploadType: string;
  name: string;
  url?: string;
  file?: string;
};

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
 * Result union of the builder. `requestBody` is the generated wire body itself
 * - not a pre-serialized string - so it can be handed straight to the
 * schemas-generated `importDesign` mutation, which serializes it. `errorMessage`
 * is UI copy for the notification shown when the form data cannot produce a
 * valid body.
 */
export type ImportDesignRequestResult =
  | { requestBody: ImportDesignRequestBody }
  | { errorMessage: string };

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

        const requestBody: ImportDesignFileVariant = {
          name,
          fileName: importedFile.fileName,
          // The contract types `file` as a base64 string (`format: byte`);
          // this client sends the decoded byte array produced after reading
          // the upload. The server's Go `[]byte` decoder accepts either
          // representation and yields identical bytes, so the byte array is a
          // conformant alternative and is preserved here unchanged. The cast is
          // confined to this single property so every other field stays pinned
          // to the generated contract.
          file: importedFile.fileData as unknown as ImportDesignFileVariant['file'],
        };
        return { requestBody };
      } catch (error) {
        console.error('Error resolving design import file:', error);
        return { errorMessage: 'Unable to read the selected design file. Please try again.' };
      }
    }
    case 'URL Import': {
      // Reject a missing or blank URL up front, mirroring the File-Upload guard
      // above. This narrows `url` to the contract's required `string`, so
      // `ImportDesignUrlVariant` is satisfied without a cast, and it avoids
      // emitting a body whose `url` key `JSON.stringify` would silently drop.
      if (typeof url !== 'string' || url.trim().length === 0) {
        return { errorMessage: 'Please enter a design URL before continuing.' };
      }

      const requestBody: ImportDesignUrlVariant = { url, name };
      return { requestBody };
    }
    default:
      return { errorMessage: 'Please choose a valid design import source before continuing.' };
  }
};
