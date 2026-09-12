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

  beforeEach(() => {
    setAttributeSpy = vi.fn();
    createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue({
      setAttribute: setAttributeSpy,
      click: vi.fn(),
      remove: vi.fn(),
    } as unknown as HTMLElement);
  });

  afterEach(() => {
    createElementSpy.mockRestore();
  });

  it('builds the pattern download URL with id only when no source_type/params', () => {
    downloadContent({ id: 'p1', type: PATTERN, name: 'pattern.yaml' });
    expect(setAttributeSpy).toHaveBeenCalledWith('href', '/api/pattern/download/p1');
    expect(setAttributeSpy).toHaveBeenCalledWith('download', 'pattern.yaml');
  });

  it('uses source_type in the pattern download URL when supplied', () => {
    downloadContent({ id: 'p1', type: PATTERN, name: 'p.yaml', source_type: 'helm' });
    expect(setAttributeSpy).toHaveBeenCalledWith('href', '/api/pattern/download/p1/helm');
  });

  it('falls back to params in the pattern download URL when source_type is not given', () => {
    downloadContent({ id: 'p1', type: PATTERN, name: 'p.yaml', params: 'export=true' });
    expect(setAttributeSpy).toHaveBeenCalledWith('href', '/api/pattern/download/p1?export=true');
  });

  it('prefers source_type over params when both are provided', () => {
    downloadContent({
      id: 'p1',
      type: PATTERN,
      name: 'p.yaml',
      source_type: 'helm',
      params: 'foo=bar',
    });
    expect(setAttributeSpy).toHaveBeenCalledWith('href', '/api/pattern/download/p1/helm');
  });

  it('builds the filter download URL', () => {
    downloadContent({ id: 'f1', type: FILTER, name: 'filter.wasm' });
    expect(setAttributeSpy).toHaveBeenCalledWith('href', '/api/filter/download/f1');
  });

  it('throws for an unknown content type', () => {
    expect(() => downloadContent({ id: 'x', type: 'unknown', name: 'x' })).toThrow();
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
