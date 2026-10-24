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

// utils imports componentUniqueKey from rtk-query/meshModel, whose module graph
// reaches the store; mock the store so the real helper can be used.
vi.mock('../../../../store', () => ({ store: { dispatch: vi.fn() } }));

import { getComponentFormKey, getComponentMetadata, getStyleOverrides } from './utils';

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

describe('getComponentFormKey', () => {
  const service = {
    component: { kind: 'Service', version: 'v1' },
    version: 'v1.0.0',
    model: { name: 'kubernetes', version: 'v1.0.0' },
  };

  it('combines the model name, component identity and position', () => {
    expect(getComponentFormKey('kubernetes', service, 3)).toBe(
      'kubernetes-Service-v1-v1.0.0-kubernetes-v1.0.0-3',
    );
  });

  it('gives the same kind at the same position a different key in another model', () => {
    const istioService = { ...service, model: { name: 'istio-base', version: 'v1.0.0' } };
    expect(getComponentFormKey('kubernetes', service, 3)).not.toBe(
      getComponentFormKey('istio-base', istioService, 3),
    );
  });

  it('gives a different key when the definition changes at the same position', () => {
    const serviceV2 = { ...service, component: { kind: 'Service', version: 'v2' } };
    expect(getComponentFormKey('kubernetes', service, 3)).not.toBe(
      getComponentFormKey('kubernetes', serviceV2, 3),
    );
  });

  it('keeps identical siblings distinct by position', () => {
    expect(getComponentFormKey('kubernetes', service, 3)).not.toBe(
      getComponentFormKey('kubernetes', service, 4),
    );
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
