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

  const buildHandler = (additionalOptions = '') =>
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

  it('preserves urlError when editing an unrelated field (e.g., qps or profileName)', () => {
    const handleChange = buildHandler();

    handleChange('qps')({ target: { value: '10' } });
    handleChange('profileName')({ target: { value: 'load-test-1' } });

    expect(setUrlError).not.toHaveBeenCalled();
  });

  it('clears urlError and re-enables test when URL is cleared with valid additionalOptions JSON', () => {
    const handleChange = buildHandler('{"headers":{"Authorization":"Bearer token"}}');
    handleChange('url')({ target: { value: '' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(false);
  });

  it('clears urlError and keeps test disabled when URL is cleared with invalid or empty additionalOptions', () => {
    const handleChange = buildHandler('{invalid-json');
    handleChange('url')({ target: { value: '' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(true);
  });

  it('clears urlError and enables test when a valid URL is entered', () => {
    const handleChange = buildHandler();
    handleChange('url')({ target: { value: 'https://example.com' } });

    expect(setUrlError).toHaveBeenCalledWith(false);
    expect(setDisableTest).toHaveBeenCalledWith(false);
  });
});
