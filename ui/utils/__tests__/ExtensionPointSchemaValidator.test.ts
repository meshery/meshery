import { describe, expect, it } from 'vitest';
import ExtensionPointSchemaValidator from '../ExtensionPointSchemaValidator';

describe('ExtensionPointSchemaValidator', () => {
  it('returns a noop decoder for unknown types', () => {
    const decode = ExtensionPointSchemaValidator('does-not-exist');
    expect(decode([{ title: 'X' }])).toBeUndefined();
  });

  describe('navigator decoder', () => {
    const decode = ExtensionPointSchemaValidator('navigator');

    it('returns an empty array for non-array input', () => {
      expect(decode(undefined)).toEqual([]);
      expect(decode(null)).toEqual([]);
      expect(decode({})).toEqual([]);
    });

    it('decodes a single nav entry with defaults for missing fields', () => {
      const out = decode([{ title: 'Home', href: { uri: '/home' }, component: 'HomeView' }]);
      expect(out).toEqual([
        {
          title: 'Home',
          href: '/extension/home',
          component: 'HomeView',
          onClickCallback: 0,
          icon: '',
          show: false,
          children: [],
          full_page: undefined,
          isBeta: false,
        },
      ]);
    });

    it('builds external hrefs as-is and prepends /extension to internal hrefs', () => {
      const out = decode([
        { title: 'Docs', href: { external: true, uri: 'https://docs.example.com' } },
        { title: 'Local', href: { uri: '/inside' } },
      ]);
      expect(out[0].href).toBe('https://docs.example.com');
      expect(out[1].href).toBe('/extension/inside');
    });

    it('uses an empty string for the href when external uri is missing', () => {
      const out = decode([{ title: 'External', href: { external: true } }]);
      expect(out[0].href).toBe('');
    });

    it('builds icon URLs through the provider extension prefix and normalizes leading slash', () => {
      const out = decode([
        { title: 'A', href: { uri: '/a' }, icon: '/icon.svg' },
        { title: 'B', href: { uri: '/b' }, icon: 'icon.svg' },
        { title: 'C', href: { uri: '/c' } },
      ]);
      expect(out[0].icon).toBe('/api/provider/extension/icon.svg');
      expect(out[1].icon).toBe('/api/provider/extension/icon.svg');
      expect(out[2].icon).toBe('');
    });

    it('preserves the show flag as a coerced boolean', () => {
      const out = decode([
        { title: 'A', href: { uri: '/a' }, show: 1 },
        { title: 'B', href: { uri: '/b' }, show: 0 },
      ]);
      expect(out[0].show).toBe(true);
      expect(out[1].show).toBe(false);
    });

    it('preserves isBeta and full_page passthrough', () => {
      const out = decode([
        { title: 'A', href: { uri: '/a' }, isBeta: true, full_page: true },
        { title: 'B', href: { uri: '/b' } },
      ]);
      expect(out[0].isBeta).toBe(true);
      expect(out[0].full_page).toBe(true);
      expect(out[1].isBeta).toBe(false);
      expect(out[1].full_page).toBeUndefined();
    });

    it('recursively decodes children', () => {
      const out = decode([
        {
          title: 'Parent',
          href: { uri: '/p' },
          children: [{ title: 'Child', href: { uri: '/c' } }],
        },
      ]);
      expect(out[0].children).toHaveLength(1);
      expect(out[0].children[0].title).toBe('Child');
      expect(out[0].children[0].href).toBe('/extension/c');
    });
  });

  describe('userPrefs decoder', () => {
    const decode = ExtensionPointSchemaValidator('userPrefs');

    it('returns an empty array for non-array input', () => {
      expect(decode(undefined)).toEqual([]);
    });

    it('extracts only the component field, defaulting to empty string', () => {
      expect(decode([{ component: 'PrefForm', type: 'rjsf' }, {}])).toEqual([
        { component: 'PrefForm' },
        { component: '' },
      ]);
    });
  });

  describe('collaborator decoder', () => {
    const decode = ExtensionPointSchemaValidator('collaborator');

    it('returns an empty array for non-array input', () => {
      expect(decode(null)).toEqual([]);
    });

    it('extracts only the component field, defaulting to empty string', () => {
      expect(decode([{ component: 'A' }, {}])).toEqual([{ component: 'A' }, { component: '' }]);
    });
  });

  describe('account decoder', () => {
    const decode = ExtensionPointSchemaValidator('account');

    it('returns an empty array for non-array input', () => {
      expect(decode(null)).toEqual([]);
    });

    it('decodes account entries with defaults and prepends /extension to internal hrefs', () => {
      const out = decode([
        { title: 'Profile', href: { uri: '/profile' }, component: 'ProfileView', link: true },
      ]);
      expect(out).toEqual([
        {
          title: 'Profile',
          href: '/extension/profile',
          component: 'ProfileView',
          onClickCallback: 0,
          show: false,
          children: [],
          full_page: undefined,
          isHtmlLink: true,
        },
      ]);
    });

    it('recursively decodes children', () => {
      const out = decode([
        {
          title: 'A',
          href: { uri: '/a' },
          children: [{ title: 'B', href: { uri: '/b' } }],
        },
      ]);
      expect(out[0].children).toHaveLength(1);
      expect(out[0].children[0].title).toBe('B');
    });

    it('respects external hrefs', () => {
      const out = decode([
        { title: 'External', href: { external: true, uri: 'https://docs.example.com' } },
      ]);
      expect(out[0].href).toBe('https://docs.example.com');
    });
  });
});
