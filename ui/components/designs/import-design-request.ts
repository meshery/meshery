import { resolveImportedDesignFile } from './import-design-file';

export type ImportDesignFormData = {
  uploadType: string;
  name: string;
  url?: string;
  file?: string;
};

export type ImportDesignRequestResult = { requestBody: string } | { errorMessage: string };

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

        return {
          requestBody: JSON.stringify({
            name,
            file_name: importedFile.fileName,
            file: importedFile.fileData,
          }),
        };
      } catch (error) {
        console.error('Error resolving design import file:', error);
        return { errorMessage: 'Unable to read the selected design file. Please try again.' };
      }
    }
    case 'URL Import':
      return {
        requestBody: JSON.stringify({
          url,
          name,
        }),
      };
    default:
      return { errorMessage: 'Please choose a valid design import source before continuing.' };
  }
};
