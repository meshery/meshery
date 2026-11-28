import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createPerformanceFormChangeHandler } from './performance-handlers';

describe('performance-handlers: URL validation state', () => {
  const setProfileName = vi.fn();
  const setMeshName = vi.fn();
  const setC = vi.fn();
  const setQps = vi.fn();
  const setHeaders = vi.fn();
  const setCookies = vi.fn();
  const setContentType = vi.fn();
  const setReqBody = vi.fn();
  const setLoadGenerator = vi.fn();
  const setUrl = vi.fn();
  const setCaCertificate = vi.fn();
  const setAdditionalOptions = vi.fn();
  const setJsonError = vi.fn();
  const setDisableTest = vi.fn();
  const setUrlError = vi.fn();

  const buildHandler = (additionalOptions = '', url = '') =>
    createPerformanceFormChangeHandler({
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
      url,
      additionalOptions,
    });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('sets urlError and disableTest to true when an invalid URL is entered', () => {
    const handleChange = buildHandler();
    handleChange('url')({ target: { value: 'not-a-url' } });

    expect(setUrlError).toHaveBeenCalledWith(true);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('preserves urlError when editing an unrelated field after entering an invalid URL', () => {
    const handleChange = buildHandler();

    handleChange('url')({ target: { value: 'not-a-url' } });
    handleChange('qps')({ target: { value: '10' } });
    handleChange('profileName')({ target: { value: 'load-test-1' } });

    expect(setUrlError).toHaveBeenCalledTimes(1);
    expect(setUrlError).toHaveBeenCalledWith(true);
    expect(setUrlError).not.toHaveBeenCalledWith(false);
  });

  it('clears urlError and re-enables test when URL is cleared with valid additionalOptions JSON', () => {
    const handleChange = buildHandler('{"headers":{"Authorization":"Bearer token"}}');
    handleChange('url')({ target: { value: '' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(false);
  });

  it('clears urlError and keeps test disabled when URL is cleared with empty additionalOptions', () => {
    const handleChange = buildHandler('');
    handleChange('url')({ target: { value: '' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('clears urlError and keeps test disabled when URL is cleared with invalid additionalOptions', () => {
    const handleChange = buildHandler('{invalid-json');
    handleChange('url')({ target: { value: '' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('clears urlError and enables test when a valid URL is entered with empty or valid additionalOptions', () => {
    const handleChange = buildHandler('');
    handleChange('url')({ target: { value: 'https://example.com' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(false);
  });

  it('clears urlError but keeps test disabled when a valid URL is entered with invalid additionalOptions', () => {
    const handleChange = buildHandler('{invalid-json');
    handleChange('url')({ target: { value: 'https://example.com' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('disables the test when typed additional options become invalid after a valid URL', () => {
    const handleChange = buildHandler('', 'https://example.com');
    handleChange('additional_options')({
      target: { value: '{invalid-json', getAttribute: () => 'text' },
    });

    expect(setJsonError).toHaveBeenCalledWith(true);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('enables the test when typed additional options become valid after an invalid URL', () => {
    const handleChange = buildHandler('', 'not-a-url');
    handleChange('additional_options')({
      target: { value: '{"profile":"fast"}', getAttribute: () => 'text' },
    });

    expect(setJsonError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(false);
  });

  it('reconciles test enablement after an additional-options file upload', () => {
    const originalFileReader = globalThis.FileReader;
    class FileReaderMock {
      onload: ((event: any) => void) | null = null;

      readAsText() {
        this.onload?.({ target: { result: '{invalid-json' } });
      }
    }
    globalThis.FileReader = FileReaderMock as any;

    try {
      const handleChange = buildHandler('', 'https://example.com');
      handleChange('additional_options')({
        target: {
          value: '',
          files: [{ name: 'options.json' }],
          getAttribute: () => 'file',
        },
      });

      expect(setJsonError).toHaveBeenCalledWith(true);
      expect(setDisableTest).toHaveBeenCalledWith(true);
    } finally {
      globalThis.FileReader = originalFileReader;
    }
  });
});
