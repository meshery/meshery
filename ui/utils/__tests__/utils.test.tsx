import React from 'react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';

// utils.tsx pulls in the next router (via pages/_app), react-redux,
// trueRandom (crypto-backed), and the mesheryEventBus singleton. Mock the
// heavy-but-irrelevant deps so unit tests can run in jsdom without standing up
// a real router / store.

const hoisted = vi.hoisted(() => ({
  eventBus: { publish: vi.fn() },
  useSelectorMock: vi.fn(),
}));
const eventBus = hoisted.eventBus;
const useSelectorMock = hoisted.useSelectorMock;

vi.mock('../../lib/trueRandom', () => ({
  trueRandom: vi.fn(() => 0.5),
}));

// Avoid evaluating `pages/_app` (which pulls notistack/material/relay/etc.).
vi.mock('../../pages/_app', () => ({
  mesheryExtensionRoute: '/extension/meshmap',
}));

vi.mock('../eventBus', () => ({
  mesheryEventBus: hoisted.eventBus,
}));

vi.mock('react-redux', () => ({
  useSelector: (sel: (s: unknown) => unknown) => hoisted.useSelectorMock(sel),
}));

import {
  ConditionalTooltip,
  ResizableCell,
  camelcaseToSnakecase,
  createScrollHandler,
  encodeDesignFile,
  formatToTitleCase,
  generateValidatePayload,
  getColumnValue,
  getComponentFromDesign,
  getComponentsinFile,
  getDecodedFile,
  getDesignVersion,
  getSharableCommonHostAndprotocolLink,
  getUnit8ArrayDecodedFile,
  getUnit8ArrayForDesign,
  getVisibilityColums,
  isDesignOpenInExtension,
  isEmptyArr,
  isEmptyObj,
  isEqualArr,
  isExtensionOpen,
  isInDesignMode,
  isInOperatorMode,
  JsonParse,
  mergeDesignWithCurrent,
  modifyRJSFSchema,
  openDesignInExtension,
  openViewInExtension,
  openViewScopedToDesignInOperator,
  parseDesignFile,
  processDesign,
  randomPatternNameGenerator,
  scrollToTop,
  updateURLs,
  urlEncodeArrayParam,
  urlEncodeParams,
  EXTENSION_MODE,
} from '../utils';

const setLocation = (url: string) => {
  const u = new URL(url);
  Object.defineProperty(window, 'location', {
    configurable: true,
    value: {
      ...window.location,
      href: u.href,
      origin: u.origin,
      protocol: u.protocol,
      host: u.host,
      hostname: u.hostname,
      port: u.port,
      pathname: u.pathname,
      search: u.search,
      hash: u.hash,
    },
  });
};

describe('isEmptyObj', () => {
  it('returns true for plain empty objects, null and undefined', () => {
    expect(isEmptyObj({})).toBe(true);
    expect(isEmptyObj(null as never)).toBe(true);
    expect(isEmptyObj(undefined as never)).toBe(true);
  });

  it('returns false for non-empty objects', () => {
    expect(isEmptyObj({ a: 1 })).toBe(false);
  });
});

describe('isEmptyArr', () => {
  it('returns true for empty arrays', () => {
    expect(isEmptyArr([])).toBe(true);
  });

  it('returns false for populated arrays', () => {
    expect(isEmptyArr([1])).toBe(false);
  });
});

describe('isEqualArr', () => {
  it('returns true for arrays with same elements in order', () => {
    expect(isEqualArr([1, 2, 3], [1, 2, 3])).toBe(true);
  });

  it('returns true when the same reference is passed twice', () => {
    const a = [1, 2];
    expect(isEqualArr(a, a)).toBe(true);
  });

  it('returns false when lengths differ', () => {
    expect(isEqualArr([1, 2, 3], [1, 2])).toBe(false);
  });

  it('returns false when an array is nullish', () => {
    expect(isEqualArr(null, [1])).toBe(false);
    expect(isEqualArr([1], null)).toBe(false);
  });

  it('compares regardless of order when orderMatters=false', () => {
    expect(isEqualArr([1, 2, 3], [3, 2, 1], false)).toBe(true);
  });

  it('returns false when elements differ', () => {
    expect(isEqualArr([1, 2, 3], [1, 2, 4])).toBe(false);
  });
});

describe('scrollToTop', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('schedules a window.scrollTo on the next tick with the supplied behavior', () => {
    const scrollSpy = vi.spyOn(window, 'scrollTo').mockImplementation(() => {});
    scrollToTop('auto');
    expect(scrollSpy).not.toHaveBeenCalled();
    vi.advanceTimersByTime(0);
    expect(scrollSpy).toHaveBeenCalledWith({ top: 0, left: 0, behavior: 'auto' });
  });

  it('defaults to smooth behavior', () => {
    const scrollSpy = vi.spyOn(window, 'scrollTo').mockImplementation(() => {});
    scrollToTop();
    vi.advanceTimersByTime(0);
    expect(scrollSpy).toHaveBeenCalledWith({ top: 0, left: 0, behavior: 'smooth' });
  });
});

describe('randomPatternNameGenerator', () => {
  it('prefixes with meshery_ and a non-negative integer', () => {
    expect(randomPatternNameGenerator()).toMatch(/^meshery_\d+$/);
  });
});

describe('getComponentsinFile', () => {
  it('returns 0 for empty file content', () => {
    expect(getComponentsinFile('')).toBe(0);
    expect(getComponentsinFile(undefined as never)).toBe(0);
  });

  it('counts services in a single-doc YAML pattern', () => {
    const yaml = 'services:\n  one:\n    x: 1\n  two:\n    y: 2\n';
    expect(getComponentsinFile(yaml)).toBe(2);
  });

  it('falls back to splitting by --- when YAML is a multi-doc stream', () => {
    const yaml = 'foo: 1\n---\nbar: 2\n---\nbaz: 3\n';
    // when js-yaml errors with "expected a single document", we split by --- separator
    expect(getComponentsinFile(yaml)).toBe(3);
  });
});

describe('generateValidatePayload', () => {
  it('reports missing services', () => {
    const result = generateValidatePayload('foo: bar\n', []);
    expect(result).toEqual({ err: 'Services not found in the design' });
  });

  it('builds validation payloads keyed by service id when a workload schema is found', () => {
    const yaml = 'services:\n  svc1:\n    type: web\n    settings:\n      port: 80\n';
    const workloads = [
      {
        workload: {
          oam_definition: { metadata: { name: 'web' } },
          oam_ref_schema: '{"type":"object"}',
        },
      },
    ];
    const result = generateValidatePayload(yaml, workloads) as Record<
      string,
      { schema: string; value: string; valueType: string }
    >;
    expect(result).toHaveProperty('svc1');
    expect(result.svc1).toEqual({
      schema: '{"type":"object"}',
      value: '{"port":80}',
      valueType: 'JSON',
    });
  });

  it('skips services whose workload definition lacks an oam_ref_schema', () => {
    const yaml = 'services:\n  svc1:\n    type: web\n    settings:\n      port: 80\n';
    const workloads = [{ workload: { oam_definition: { metadata: { name: 'web' } } } }];
    expect(generateValidatePayload(yaml, workloads)).toEqual({});
  });

  it('skips services whose workload is missing settings', () => {
    const yaml = 'services:\n  svc1:\n    type: web\n';
    const workloads = [
      {
        workload: {
          oam_definition: { metadata: { name: 'web' } },
          oam_ref_schema: '{}',
        },
      },
    ];
    expect(generateValidatePayload(yaml, workloads)).toEqual({});
  });
});

describe('updateURLs', () => {
  it('adds urls on ADDED/MODIFIED events', () => {
    const set = new Set<string>();
    updateURLs(set, ['/a', '/b'], 'ADDED');
    expect(Array.from(set)).toEqual(['/a', '/b']);
    updateURLs(set, ['/c'], 'MODIFIED');
    expect(set.has('/c')).toBe(true);
  });

  it('removes urls on DELETED events', () => {
    const set = new Set<string>(['/a', '/b']);
    updateURLs(set, ['/a'], 'DELETED');
    expect(set.has('/a')).toBe(false);
    expect(set.has('/b')).toBe(true);
  });
});

describe('base64 decoders', () => {
  it('decodes a data URL into the raw bytes', () => {
    const text = 'hello';
    const dataUrl = `data:text/plain;base64,${btoa(text)}`;
    expect(getDecodedFile(dataUrl)).toBe(text);
  });

  it('decodes a data URL into an array of byte values', () => {
    const text = 'AB';
    const dataUrl = `data:text/plain;base64,${btoa(text)}`;
    expect(getUnit8ArrayDecodedFile(dataUrl)).toEqual([0x41, 0x42]);
  });

  it('converts a design string into an array of byte values', () => {
    expect(getUnit8ArrayForDesign('AB')).toEqual([0x41, 0x42]);
  });
});

describe('modifyRJSFSchema', () => {
  it('clones the schema and mutates a deep property', () => {
    const schema = { properties: { x: { type: 'string' } } };
    const out = modifyRJSFSchema(schema, 'properties.x.type', 'number');
    expect(out.properties.x.type).toBe('number');
    expect(schema.properties.x.type).toBe('string'); // original untouched
    expect(out).not.toBe(schema);
  });
});

describe('getSharableCommonHostAndprotocolLink', () => {
  beforeEach(() => {
    setLocation('http://localhost:9081/');
  });

  it('returns an extension/meshmap link with the application id for application_file resources', () => {
    expect(getSharableCommonHostAndprotocolLink({ application_file: 'foo', id: 'app-1' })).toBe(
      'http://localhost:9081/extension/meshmap?application=app-1',
    );
  });

  it('returns an extension/meshmap link with mode=design for patternFile resources', () => {
    expect(getSharableCommonHostAndprotocolLink({ patternFile: 'foo', id: 'd-1' })).toBe(
      'http://localhost:9081/extension/meshmap?mode=design&design=d-1',
    );
  });

  it('returns an extension/meshmap link with filter id for filter_resource', () => {
    expect(getSharableCommonHostAndprotocolLink({ filter_resource: 'foo', id: 'f-1' })).toBe(
      'http://localhost:9081/extension/meshmap?filter=f-1',
    );
  });

  it('returns an empty string when nothing matches', () => {
    expect(getSharableCommonHostAndprotocolLink({})).toBe('');
  });
});

describe('getColumnValue / getVisibilityColums', () => {
  const columns = [{ name: 'id' }, { name: 'status' }];

  it('returns the value at the index of the named column', () => {
    expect(getColumnValue(['x', 'ok'], 'status', columns)).toBe('ok');
  });

  it('returns undefined when the column is missing', () => {
    expect(getColumnValue(['x', 'ok'], 'missing', columns)).toBeUndefined();
  });

  it('filters out columns with options.display === false', () => {
    expect(
      getVisibilityColums([
        { name: 'a', options: { display: false } },
        { name: 'b' },
        { name: 'c', options: { display: true } },
      ]),
    ).toHaveLength(2);
  });
});

describe('JsonParse', () => {
  it('returns the input when not a string', () => {
    const obj = { a: 1 };
    expect(JsonParse(obj as unknown as string)).toBe(obj);
  });

  it('parses valid JSON strings', () => {
    expect(JsonParse('{"a":1}')).toEqual({ a: 1 });
  });

  it('returns {} for invalid strings in safe mode', () => {
    expect(JsonParse('garbage')).toEqual({});
  });

  it('rethrows in unsafe mode', () => {
    expect(() => JsonParse('garbage', false)).toThrow();
  });

  it('treats empty strings as {} (per implementation)', () => {
    expect(JsonParse('')).toEqual({});
  });
});

describe('ConditionalTooltip & ResizableCell render', () => {
  it('renders the value untruncated when it is shorter than maxLength', () => {
    render(<ConditionalTooltip value="abc" maxLength={10} />);
    expect(screen.getByText('abc')).toBeInTheDocument();
  });

  it('truncates and wraps in tooltip when longer than maxLength', () => {
    render(<ConditionalTooltip value="abcdefghij" maxLength={3} />);
    expect(screen.getByText('abc...')).toBeInTheDocument();
  });

  it('renders the resizable cell with its value', () => {
    render(<ResizableCell value="cell" />);
    expect(screen.getByText('cell')).toBeInTheDocument();
  });
});

describe('createScrollHandler', () => {
  it('advances the page when the scroll target reaches the buffered bottom', () => {
    const setPage = vi.fn();
    const scrollRef: { current: number | null } = { current: null };
    const handler = createScrollHandler('view-a', setPage, scrollRef, 10);

    // event div: scrollTop=90, scrollHeight=100, clientHeight=0, buffer=10 -> at limit
    handler({ target: { scrollTop: 90, scrollHeight: 100, clientHeight: 0 } });

    expect(setPage).toHaveBeenCalledTimes(1);
    const updater = setPage.mock.calls[0][0] as (
      prev: Record<string, number>,
    ) => Record<string, number>;
    expect(updater({ 'view-a': 2 })).toEqual({ 'view-a': 3 });
    expect(scrollRef.current).toBe(90);
  });

  it('does not advance the page when not yet at the bottom', () => {
    const setPage = vi.fn();
    const scrollRef: { current: number | null } = { current: null };
    const handler = createScrollHandler('view-b', setPage, scrollRef, 10);
    handler({ target: { scrollTop: 10, scrollHeight: 1000, clientHeight: 0 } });
    expect(setPage).not.toHaveBeenCalled();
    expect(scrollRef.current).toBe(10);
  });
});

describe('camelcaseToSnakecase / formatToTitleCase', () => {
  it('converts camelCase to lower_snake_case with leading underscore behaviour', () => {
    expect(camelcaseToSnakecase('camelCase')).toBe('camel_case');
    expect(camelcaseToSnakecase('FooBar')).toBe('_foo_bar');
  });

  it('passes through undefined input', () => {
    expect(camelcaseToSnakecase(undefined as never)).toBeUndefined();
  });

  it('formats first letter uppercase, the rest lowercase', () => {
    expect(formatToTitleCase('helloWorld')).toBe('Helloworld');
    expect(formatToTitleCase('HELLO')).toBe('Hello');
  });

  it('returns empty string when value is not a string', () => {
    expect(formatToTitleCase(42 as never)).toBe('');
    expect(formatToTitleCase(undefined as never)).toBe('');
  });
});

describe('parseDesignFile / encodeDesignFile', () => {
  it('round-trips a JSON-compatible design', () => {
    const designJson = { schemaVersion: 'designs.meshery.io/v1beta1', components: [] };
    const yaml = encodeDesignFile(designJson);
    expect(typeof yaml).toBe('string');
    expect(parseDesignFile(yaml as string)).toEqual(designJson);
  });

  it('returns null for invalid YAML', () => {
    const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    expect(parseDesignFile('::not valid\n:::')).toBeNull();
    errSpy.mockRestore();
  });
});

describe('processDesign', () => {
  it('returns empty arrays when the schema version is invalid', () => {
    const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const result = processDesign({ schemaVersion: 'designs.meshery.io/v1alpha1', components: [] });
    expect(result.configurableComponents).toEqual([]);
    expect(result.annotationComponents).toEqual([]);
    expect((result.designJson as { name?: string }).name).toBe('');
    errSpy.mockRestore();
  });

  it('accepts v1beta3 designs and preserves their components', () => {
    const design = {
      schemaVersion: 'designs.meshery.io/v1beta3',
      components: [{ id: 'a' }, { id: 'b', metadata: { isAnnotation: true } }],
    };

    const result = processDesign(design);

    expect(result.components).toEqual(design.components);
    expect(result.configurableComponents.map((c: { id: string }) => c.id)).toEqual(['a']);
    expect(result.annotationComponents.map((c: { id: string }) => c.id)).toEqual(['b']);
    expect(result.designJson).toBe(design);
  });

  it('separates annotation components from configurable ones', () => {
    const design = {
      schemaVersion: 'designs.meshery.io/v1beta1',
      components: [{ id: 'a' }, { id: 'b', metadata: { isAnnotation: true } }, { id: 'c' }],
    };
    const result = processDesign(design);
    expect(result.components).toHaveLength(3);
    expect(result.configurableComponents.map((c: { id: string }) => c.id)).toEqual(['a', 'c']);
    expect(result.annotationComponents.map((c: { id: string }) => c.id)).toEqual(['b']);
    expect(result.designJson).toBe(design);
  });
});

describe('getComponentFromDesign', () => {
  it('returns the component with the matching id', () => {
    const design = { components: [{ id: '1' }, { id: '2' }] };
    expect(getComponentFromDesign(design, '2')).toEqual({ id: '2' });
  });

  it('returns undefined when no component matches', () => {
    const design = { components: [{ id: '1' }] };
    expect(getComponentFromDesign(design, 'x')).toBeUndefined();
  });
});

describe('getDesignVersion', () => {
  it('returns the published version when visibility is published', () => {
    expect(
      getDesignVersion({ visibility: 'published', catalog_data: { published_version: '1.2.3' } }),
    ).toBe('1.2.3');
  });

  it('falls back to the parsed pattern_file version for non-published designs', () => {
    expect(getDesignVersion({ patternFile: 'version: 0.1.0' })).toBe('0.1.0');
  });

  it('logs and swallows YAML errors silently', () => {
    const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    expect(getDesignVersion({ patternFile: '::invalid' })).toBeUndefined();
    errSpy.mockRestore();
  });
});

describe('urlEncodeArrayParam / urlEncodeParams', () => {
  it('passes through a string parameter value unchanged', () => {
    expect(urlEncodeArrayParam('id', 'abc')).toBe('abc');
  });

  it('builds repeated key=value pairs from an array', () => {
    expect(urlEncodeArrayParam('id', ['a', 'b'])).toBe('id=a&id=b');
  });

  it('produces a query string from a record', () => {
    expect(urlEncodeParams({ name: 'demo', mode: 'design' })).toBe('name=demo&mode=design');
  });

  it('repeats keys when a value is an array', () => {
    expect(urlEncodeParams({ id: ['a', 'b'] })).toBe('id=a&id=b');
  });

  it('skips nil values', () => {
    expect(urlEncodeParams({ id: undefined, name: 'demo' })).toBe('name=demo');
    expect(urlEncodeParams({ id: null, name: 'demo' })).toBe('name=demo');
  });
});

describe('extension / helpers', () => {
  it('isExtensionOpen returns true under /extension/meshmap', () => {
    setLocation('http://localhost:9081/extension/meshmap');
    expect(isExtensionOpen()).toBe(true);
  });

  it('isExtensionOpen returns false elsewhere', () => {
    setLocation('http://localhost:9081/dashboard');
    expect(isExtensionOpen()).toBe(false);
  });

  it('isDesignOpenInExtension reads `design` + mode=design from query', () => {
    setLocation('http://localhost:9081/x?design=abc&mode=design');
    expect(isDesignOpenInExtension()).toBe(true);
  });

  it('isDesignOpenInExtension returns false when mode is not design', () => {
    setLocation('http://localhost:9081/x?design=abc&mode=operator');
    expect(isDesignOpenInExtension()).toBe(false);
  });

  it('isDesignOpenInExtension returns false when no design param', () => {
    setLocation('http://localhost:9081/x?mode=design');
    expect(isDesignOpenInExtension()).toBe(false);
  });

  it('isInDesignMode / isInOperatorMode read mode from search', () => {
    setLocation('http://localhost:9081/x?mode=design');
    expect(isInDesignMode()).toBe(true);
    expect(isInOperatorMode()).toBe(false);

    setLocation('http://localhost:9081/x?mode=operator');
    expect(isInOperatorMode()).toBe(true);
    expect(isInDesignMode()).toBe(false);
  });

  it('exposes the EXTENSION_MODE constants', () => {
    expect(EXTENSION_MODE).toEqual({ DESIGN: 'design', OPERATOR: 'operator' });
  });
});

describe('event-bus router helpers', () => {
  beforeEach(() => {
    eventBus.publish.mockClear();
  });

  it('publishes OPEN_VIEW_SCOPED_TO_DESIGN inside extension', () => {
    setLocation('http://localhost:9081/extension/meshmap');
    openViewScopedToDesignInOperator('My Design', 'd-1', { push: vi.fn() });
    expect(eventBus.publish).toHaveBeenCalledWith({
      type: 'OPEN_VIEW_SCOPED_TO_DESIGN',
      data: { designId: 'd-1', designName: 'My Design' },
    });
  });

  it('pushes to the operator route outside extension', () => {
    setLocation('http://localhost:9081/dashboard');
    const router = { push: vi.fn() };
    openViewScopedToDesignInOperator('My Design', 'd-1', router);
    expect(eventBus.publish).not.toHaveBeenCalled();
    expect(router.push).toHaveBeenCalledWith(
      '/extension/meshmap?mode=operator&type=view&design_id=d-1',
    );
  });

  it('publishes MERGE_DESIGN unconditionally', () => {
    mergeDesignWithCurrent('d-1', 'D1');
    expect(eventBus.publish).toHaveBeenCalledWith({
      type: 'MERGE_DESIGN',
      data: { id: 'd-1', name: 'D1' },
    });
  });

  it('openDesignInExtension publishes inside extension and routes outside', () => {
    setLocation('http://localhost:9081/extension/meshmap');
    openDesignInExtension('d-1', 'D1', { push: vi.fn() });
    expect(eventBus.publish).toHaveBeenCalled();
    eventBus.publish.mockClear();
    setLocation('http://localhost:9081/dashboard');
    const router = { push: vi.fn() };
    openDesignInExtension('d-1', 'D1', router);
    expect(router.push).toHaveBeenCalledWith('/extension/meshmap?mode=design&type=design&id=d-1');
  });

  it('openViewInExtension publishes inside extension and routes outside', () => {
    setLocation('http://localhost:9081/extension/meshmap');
    openViewInExtension('v-1', 'V1', { push: vi.fn() });
    expect(eventBus.publish).toHaveBeenCalled();
    eventBus.publish.mockClear();
    setLocation('http://localhost:9081/dashboard');
    const router = { push: vi.fn() };
    openViewInExtension('v-1', 'V1', router);
    expect(router.push).toHaveBeenCalledWith('/extension/meshmap?mode=operator&type=view&id=v-1');
  });
});
