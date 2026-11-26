import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import downloadContent, { downloadFileFromUrl, downloadFileFromContent } from '../fileDownloader';
import { FILTER, PATTERN } from '@/constants/navigator';

describe('downloadFileFromUrl', () => {
  let createElementSpy: ReturnType<typeof vi.spyOn>;
  let setAttributeSpy: ReturnType<typeof vi.fn>;
  let clickSpy: ReturnType<typeof vi.fn>;
  let removeSpy: ReturnType<typeof vi.fn>;
  let fakeAnchor: HTMLAnchorElement;

  beforeEach(() => {
    setAttributeSpy = vi.fn();
    clickSpy = vi.fn();
    removeSpy = vi.fn();
    fakeAnchor = {
      setAttribute: setAttributeSpy,
      click: clickSpy,
      remove: removeSpy,
    } as unknown as HTMLAnchorElement;

    createElementSpy = vi
      .spyOn(document, 'createElement')
      .mockReturnValue(fakeAnchor as unknown as HTMLElement);
  });

  afterEach(() => {
    createElementSpy.mockRestore();
  });

  it('creates an anchor, sets href/download, clicks, and removes it', () => {
    downloadFileFromUrl('https://example.com/file.json', 'file.json');

    expect(createElementSpy).toHaveBeenCalledWith('a');
    expect(setAttributeSpy).toHaveBeenCalledWith('href', 'https://example.com/file.json');
    expect(setAttributeSpy).toHaveBeenCalledWith('download', 'file.json');
    expect(clickSpy).toHaveBeenCalledTimes(1);
    expect(removeSpy).toHaveBeenCalledTimes(1);
  });
});

describe('downloadContent', () => {
  let createElementSpy: ReturnType<typeof vi.spyOn>;
  let setAttributeSpy: ReturnType<typeof vi.fn>;
  let fetchMock: ReturnType<typeof vi.fn>;
  let createObjectURL: ReturnType<typeof vi.fn>;
  let revokeObjectURL: ReturnType<typeof vi.fn>;
  let blob: Blob;

  beforeEach(() => {
    blob = new Blob(['download bytes'], { type: 'application/wasm' });
    fetchMock = vi.fn().mockResolvedValue({ ok: true, blob: async () => blob });
    vi.stubGlobal('fetch', fetchMock);
    createObjectURL = vi.fn().mockReturnValue('blob:download');
    revokeObjectURL = vi.fn();
    vi.stubGlobal('URL', { createObjectURL, revokeObjectURL });
    setAttributeSpy = vi.fn();
    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue({
      setAttribute: setAttributeSpy,
      click: vi.fn(),
      remove: vi.fn(),
    } as unknown as HTMLElement);
  });

  afterEach(() => {
    createElementSpy.mockRestore();
    vi.unstubAllGlobals();
  });

  it.each([
    [{ id: 'p1', type: PATTERN }, '/api/pattern/download/p1'],
    [{ id: 'p1', type: PATTERN, source_type: 'helm' }, '/api/pattern/download/p1/helm'],
    [{ id: 'p1', type: PATTERN, params: 'export=true' }, '/api/pattern/download/p1?export=true'],
    [
      { id: 'p1', type: PATTERN, source_type: 'helm', params: 'foo=bar' },
      '/api/pattern/download/p1/helm',
    ],
    [{ id: 'f1', type: FILTER }, '/api/filter/download/f1'],
  ])('downloads the response bytes from %j', async (args, uri) => {
    await downloadContent({ ...args, name: 'download' });
    expect(fetchMock).toHaveBeenCalledWith(uri);
    expect(createObjectURL).toHaveBeenCalledWith(blob);
    expect(setAttributeSpy).toHaveBeenCalledWith('href', 'blob:download');
    expect(setAttributeSpy).toHaveBeenCalledWith('download', 'download');
    expect(revokeObjectURL).toHaveBeenCalledWith('blob:download');
  });

  it('waits for the response body before starting the download', async () => {
    let finishBody!: (_blob: Blob) => void;
    fetchMock.mockResolvedValueOnce({
      ok: true,
      blob: () =>
        new Promise<Blob>((resolve) => {
          finishBody = resolve;
        }),
    });
    const download = downloadContent({ id: 'f1', type: FILTER, name: 'filter.wasm' });
    await Promise.resolve();
    expect(createElementSpy).not.toHaveBeenCalled();
    finishBody(blob);
    await download;
    expect(setAttributeSpy).toHaveBeenCalledWith('download', 'filter.wasm');
  });

  it.each([FILTER, PATTERN])('rejects HTTP errors for %s without saving an error body', async (type) => {
    fetchMock.mockResolvedValueOnce({ ok: false, status: 404 });
    await expect(downloadContent({ id: 'missing', type, name: 'missing' })).rejects.toThrow(
      'Download failed (HTTP 404)',
    );
    expect(createElementSpy).not.toHaveBeenCalled();
    expect(createObjectURL).not.toHaveBeenCalled();
  });

  it('propagates network and body-read failures without starting a download', async () => {
    fetchMock.mockRejectedValueOnce(new Error('Network failed'));
    await expect(downloadContent({ id: 'f1', type: FILTER, name: 'f' })).rejects.toThrow(
      'Network failed',
    );
    fetchMock.mockResolvedValueOnce({
      ok: true,
      blob: () => Promise.reject(new Error('Body failed')),
    });
    await expect(downloadContent({ id: 'f1', type: FILTER, name: 'f' })).rejects.toThrow('Body failed');
    expect(createElementSpy).not.toHaveBeenCalled();
  });

  it('rejects an unknown content type', async () => {
    await expect(downloadContent({ id: 'x', type: 'unknown', name: 'x' })).rejects.toThrow();
  });
});

describe('downloadFileFromContent', () => {
  let createObjectURLSpy: ReturnType<typeof vi.fn>;
  let createElementSpy: ReturnType<typeof vi.spyOn>;
  let setAttributeSpy: ReturnType<typeof vi.fn>;
  let originalCreateObjectURL: typeof window.URL.createObjectURL | undefined;

  beforeEach(() => {
    setAttributeSpy = vi.fn();
    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue({
      setAttribute: setAttributeSpy,
      click: vi.fn(),
      remove: vi.fn(),
    } as unknown as HTMLElement);

    createObjectURLSpy = vi.fn().mockReturnValue('blob:fake-url');
    originalCreateObjectURL = window.URL.createObjectURL;
    Object.defineProperty(window.URL, 'createObjectURL', {
      configurable: true,
      writable: true,
      value: createObjectURLSpy,
    });
  });

  afterEach(() => {
    createElementSpy.mockRestore();
    if (originalCreateObjectURL) {
      Object.defineProperty(window.URL, 'createObjectURL', {
        configurable: true,
        writable: true,
        value: originalCreateObjectURL,
      });
    }
  });

  it('creates a Blob with the supplied MIME type and downloads it', () => {
    downloadFileFromContent('{"hello":"world"}', 'data.json', 'application/json');
    expect(createObjectURLSpy).toHaveBeenCalledTimes(1);
    const blob = createObjectURLSpy.mock.calls[0][0] as Blob;
    expect(blob).toBeInstanceOf(Blob);
    expect(blob.type).toBe('application/json');
    expect(setAttributeSpy).toHaveBeenCalledWith('href', 'blob:fake-url');
    expect(setAttributeSpy).toHaveBeenCalledWith('download', 'data.json');
  });
});
