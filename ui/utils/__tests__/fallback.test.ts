import { describe, expect, it } from 'vitest';
import { getFallbackImageBasedOnKind, normalizeStaticImagePath } from '../fallback';

describe('getFallbackImageBasedOnKind', () => {
  it('returns the meshery logo for the meshery kind', () => {
    expect(getFallbackImageBasedOnKind('meshery')).toBe('static/img/meshery-logo/meshery-logo.png');
  });

  it('returns the kubernetes svg for the kubernetes kind', () => {
    expect(getFallbackImageBasedOnKind('kubernetes')).toBe(
      'static/img/integrations/kubernetes.svg',
    );
  });

  it('returns undefined for unknown kinds', () => {
    expect(getFallbackImageBasedOnKind('istio')).toBeUndefined();
    expect(getFallbackImageBasedOnKind('')).toBeUndefined();
    expect(getFallbackImageBasedOnKind(undefined as unknown as string)).toBeUndefined();
  });
});

describe('normalizeStaticImagePath', () => {
  it('returns an empty string for falsy / blank input', () => {
    expect(normalizeStaticImagePath(undefined)).toBe('');
    expect(normalizeStaticImagePath(null)).toBe('');
    expect(normalizeStaticImagePath('')).toBe('');
    expect(normalizeStaticImagePath('   ')).toBe('');
  });

  it('treats placeholder tokens (empty/none/null/undefined) as empty', () => {
    expect(normalizeStaticImagePath('empty')).toBe('');
    expect(normalizeStaticImagePath('NONE')).toBe('');
    expect(normalizeStaticImagePath(' null ')).toBe('');
    expect(normalizeStaticImagePath('Undefined')).toBe('');
  });

  it('passes through fully-qualified http(s) URLs unchanged (after trim)', () => {
    expect(normalizeStaticImagePath('http://example.com/img.png')).toBe(
      'http://example.com/img.png',
    );
    expect(normalizeStaticImagePath('  https://example.com/img.png  ')).toBe(
      'https://example.com/img.png',
    );
  });

  it('passes through data: and blob: URIs unchanged', () => {
    expect(normalizeStaticImagePath('data:image/png;base64,abc')).toBe('data:image/png;base64,abc');
    expect(normalizeStaticImagePath('blob:something')).toBe('blob:something');
  });

  it('encodes inline SVG markup as a data URI', () => {
    const svg = '<svg viewBox="0 0 1 1"><path d="M0 0h1v1z"/></svg>';
    expect(normalizeStaticImagePath(svg)).toBe(
      `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`,
    );
    // Also handles SVGs that open with an XML prolog.
    const xmlSvg = '<?xml version="1.0"?><svg></svg>';
    expect(normalizeStaticImagePath(xmlSvg)).toBe(
      `data:image/svg+xml;charset=utf-8,${encodeURIComponent(xmlSvg)}`,
    );
  });

  it('prepends a slash to bare static paths', () => {
    expect(normalizeStaticImagePath('static/img/foo.png')).toBe('/static/img/foo.png');
  });

  it('strips leading slashes before re-prepending exactly one', () => {
    expect(normalizeStaticImagePath('///static/img/foo.png')).toBe('/static/img/foo.png');
  });

  it('keeps ui/public/static/img/meshmodels/ paths as-is under a leading slash', () => {
    expect(normalizeStaticImagePath('ui/public/static/img/meshmodels/foo.svg')).toBe(
      '/ui/public/static/img/meshmodels/foo.svg',
    );
  });

  it('strips ui/public/ prefix for other paths', () => {
    expect(normalizeStaticImagePath('ui/public/static/img/foo.png')).toBe('/static/img/foo.png');
  });

  it('trims whitespace around the input', () => {
    expect(normalizeStaticImagePath('  static/img/foo.png  ')).toBe('/static/img/foo.png');
  });
});
