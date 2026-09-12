import { URLValidator } from '../../utils/URLValidator';

/**
 * Setters bag for `createPerformanceFormChangeHandler`. Each field maps 1:1
 * to a `setX` call from the parent component's `useState` hooks.
 */
export interface PerformanceFormSetters {
  setProfileName: (v: string) => void;
  setMeshName: (v: string) => void;
  setC: (v: any) => void;
  setQps: (v: any) => void;
  setHeaders: (v: string) => void;
  setCookies: (v: string) => void;
  setContentType: (v: string) => void;
  setReqBody: (v: string) => void;
  setLoadGenerator: (v: string) => void;
  setUrl: (v: string) => void;
  setCaCertificate: (v: { name: string; file: any }) => void;
  setAdditionalOptions: (v: string) => void;
  setJsonError: (v: boolean) => void;
  setDisableTest: (v: boolean) => void;
  setUrlError: (v: boolean) => void;
}

/**
 * createPerformanceFormChangeHandler builds the `handleChange(name)(event)`
 * curry used by the performance form. The shape mirrors the original inline
 * function in `performance/index.tsx`: dispatch by `name`, with extra branches
 * for `url` validation, `caCertificate` file upload, and `additional_options`
 * JSON parsing.
 *
 * Extracted in Phase 5.a so the entry point stays under the 600-line size
 * budget. No behaviour changes — only the setters are now received as a bag.
 */
export const createPerformanceFormChangeHandler = (setters: PerformanceFormSetters) => {
  const {
    setProfileName,
    setMeshName,
    setC,
    setQps,
    setHeaders,
    setCookies,
    setContentType,
    setReqBody,
    setLoadGenerator,
    setUrl,
    setCaCertificate,
    setAdditionalOptions,
    setJsonError,
    setDisableTest,
    setUrlError,
  } = setters;

  return (name: string) => (event: any) => {
    const { value } = event.target;
    if (name === 'caCertificate') {
      if (!event.target.files?.length) return;

      const file = event.target.files[0];
      const reader = new FileReader();
      reader.addEventListener('load', (evt: any) => {
        setCaCertificate({
          name: file.name,
          file: evt.target.result,
        });
      });
      reader.readAsText(file);
    }

    if (name === 'url' && value !== '') {
      let urlPattern = value;

      let val = URLValidator(urlPattern);
      if (!val) {
        setDisableTest(true);
        setUrlError(true);
      } else {
        setDisableTest(false);
        setUrlError(false);
      }
    } else setUrlError(false);

    if (name === 'additional_options') {
      // Check if the target event is an input element (typing) or a file input (upload)
      const isFileUpload = event.target.getAttribute('type') === 'file';

      if (isFileUpload) {
        // Handle file upload
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (event: any) => {
            try {
              const fileContent = event.target.result;
              // Validate JSON
              JSON.parse(fileContent);
              setAdditionalOptions(fileContent);
              setJsonError(false);
            } catch {
              setAdditionalOptions(event.target.result);
              setJsonError(true);
            }
          };
          reader.readAsText(file);
        }
      } else {
        // Handle text input
        try {
          // empty text input exception
          if (value !== '') JSON.parse(value);
          setAdditionalOptions(value);
          setJsonError(false);
        } catch {
          setAdditionalOptions(value);
          setJsonError(true);
        }
      }
    }
    switch (name) {
      case 'profileName':
        setProfileName(value);
        break;
      case 'meshName':
        setMeshName(value);
        break;
      case 'c':
        setC(value);
        break;
      case 'qps':
        setQps(value);
        break;
      case 'headers':
        setHeaders(value);
        break;
      case 'cookies':
        setCookies(value);
        break;
      case 'contentType':
        setContentType(value);
        break;
      case 'reqBody':
        setReqBody(value);
        break;
      case 'loadGenerator':
        setLoadGenerator(value);
        break;
      case 'url':
        setUrl(value);
        break;
      default:
        // Handle any other cases or do nothing if not matched
        break;
    }
  };
};
