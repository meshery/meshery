import { describe, expect, it, vi } from 'vitest';

const jsonParse = vi.fn((input: string) => {
  try {
    return JSON.parse(input);
  } catch {
    return {};
  }
});

vi.mock('../../../../utils/utils', () => ({
  JsonParse: (input: string) => jsonParse(input),
}));

import { getComponentMetadata, getStyleOverrides } from './utils';

describe('getComponentMetadata', () => {
  it('applies the default shape, colors and svg fields when metadata is empty', () => {
    const result = getComponentMetadata({});
    expect(result.shape).toBe('circle');
    expect(result.primaryColor).toBe('#00B39F');
    expect(result.secondaryColor).toBe('#b2e8e2');
    expect(result.svgColor).toBe('static/img/component-svg/meshery.svg');
    expect(result.svgWhite).toBe('static/img/component-svg/meshery.svg');
  });

  it('lets caller-provided fields override the defaults', () => {
    const result = getComponentMetadata({ shape: 'square', primaryColor: '#FFF' });
    expect(result.shape).toBe('square');
    expect(result.primaryColor).toBe('#FFF');
    // Defaults remain for other fields:
    expect(result.secondaryColor).toBe('#b2e8e2');
  });

  it('copies background-image into svgWhite for legacy styled metadata', () => {
    const result = getComponentMetadata({ 'background-image': 'url(custom.svg)' });
    expect(result.svgWhite).toBe('url(custom.svg)');
  });
});

describe('getStyleOverrides', () => {
  it('returns an empty object when no overrides string is supplied', () => {
    expect(getStyleOverrides(undefined)).toEqual({});
    expect(getStyleOverrides('')).toEqual({});
  });

  it('delegates to JsonParse when an overrides string is supplied', () => {
    jsonParse.mockClear();
    jsonParse.mockReturnValueOnce({ 'background-opacity': 0.5 });

    const result = getStyleOverrides('{"background-opacity":0.5}');
    expect(jsonParse).toHaveBeenCalledWith('{"background-opacity":0.5}');
    expect(result).toEqual({ 'background-opacity': 0.5 });
  });
});
