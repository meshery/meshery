import { getErrorMessage } from '../ConnectionTable.constants';

type MeshkitErrorBody = {
  error?: string;
  longDescription?: string[];
  probableCause?: string[];
};

type RtkErrorLike = { data?: MeshkitErrorBody | string };

/**
 * Extracts a human-readable message from a connection API failure. RTK Query
 * nests the parsed meshkit error body under `error.data`; the shared
 * getErrorMessage only understands string payloads, so dig into the meshkit
 * shape (longDescription -> error) first and fall back to it otherwise.
 */
export const formatWizardError = (error: unknown): string => {
  const data = (error as RtkErrorLike)?.data;
  if (data && typeof data === 'object') {
    if (Array.isArray(data.longDescription) && typeof data.longDescription[0] === 'string') {
      return data.longDescription[0];
    }
    if (typeof data.error === 'string') {
      return data.error;
    }
  }
  return getErrorMessage(error);
};
