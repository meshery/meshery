import { describe, expect, it } from 'vitest';
import { ProviderUiAccessControl } from '../disabledComponents';

describe('ProviderUiAccessControl', () => {
  it('exposes provider capabilities as an array for UI guards', () => {
    const registry = new ProviderUiAccessControl({
      capabilities: [{ feature: 'persist-meshery-patterns' }],
    });

    expect(registry.providerCapabilities).toEqual([{ feature: 'persist-meshery-patterns' }]);
  });

  it('returns an empty capability list when the provider payload is malformed', () => {
    const registry = new ProviderUiAccessControl({
      capabilities: { feature: 'persist-meshery-patterns' },
    });

    expect(registry.providerCapabilities).toEqual([]);
  });

  describe('isNavigatorComponentEnabled', () => {
    it('returns true when not in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({});
      expect(accessControl.isNavigatorComponentEnabled(['designs'])).toBe(true);
      expect(accessControl.isNavigatorComponentEnabled(undefined)).toBe(true);
    });

    it('resolves component enabled state in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({
        restrictedAccess: {
          isMesheryUIRestricted: true,
          allowedComponents: {
            navigator: {
              designs: true,
              filters: false,
            },
          },
        },
      });

      expect(accessControl.isNavigatorComponentEnabled(['designs'])).toBe(true);
      expect(accessControl.isNavigatorComponentEnabled(['filters'])).toBe(false);
      expect(accessControl.isNavigatorComponentEnabled(['nonexistent'])).toBe(false);
    });

    it('safely handles non-array or undefined walker arguments in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({
        restrictedAccess: {
          isMesheryUIRestricted: true,
        },
      });

      expect(() => accessControl.isNavigatorComponentEnabled(undefined)).not.toThrow();
      expect(() => accessControl.isNavigatorComponentEnabled(null)).not.toThrow();
      expect(accessControl.isNavigatorComponentEnabled(undefined)).toBe(false);
      expect(accessControl.isNavigatorComponentEnabled(null)).toBe(false);
    });
  });

  describe('isHeaderComponentEnabled', () => {
    it('returns true when not in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({});
      expect(accessControl.isHeaderComponentEnabled(['userProfile'])).toBe(true);
      expect(accessControl.isHeaderComponentEnabled(undefined)).toBe(true);
    });

    it('resolves header component enabled state in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({
        restrictedAccess: {
          isMesheryUIRestricted: true,
          allowedComponents: {
            header: {
              notificationCenter: true,
              settings: false,
            },
          },
        },
      });

      expect(accessControl.isHeaderComponentEnabled(['notificationCenter'])).toBe(true);
      expect(accessControl.isHeaderComponentEnabled(['settings'])).toBe(false);
      expect(accessControl.isHeaderComponentEnabled(['missing'])).toBe(false);
    });

    it('safely handles non-array or undefined walker arguments in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({
        restrictedAccess: {
          isMesheryUIRestricted: true,
        },
      });

      expect(() => accessControl.isHeaderComponentEnabled(undefined)).not.toThrow();
      expect(() => accessControl.isHeaderComponentEnabled(null)).not.toThrow();
      expect(accessControl.isHeaderComponentEnabled(undefined)).toBe(false);
      expect(accessControl.isHeaderComponentEnabled(null)).toBe(false);
    });
  });

  describe('isExtensionComponentEnabled', () => {
    it('returns true when not in playground environment', () => {
      const accessControl = new ProviderUiAccessControl({});
      expect(accessControl.isExtensionComponentEnabled(['meshmap'])).toBe(true);
      expect(accessControl.isExtensionComponentEnabled(undefined)).toBe(true);
    });

    it('returns false when extensions or navigator is missing or empty', () => {
      const accessControl1 = new ProviderUiAccessControl({
        restrictedAccess: { isMesheryUIRestricted: true },
      });
      expect(accessControl1.isExtensionComponentEnabled(['meshmap'])).toBe(false);

      const accessControl2 = new ProviderUiAccessControl({
        restrictedAccess: { isMesheryUIRestricted: true },
        extensions: { navigator: [] },
      });
      expect(accessControl2.isExtensionComponentEnabled(['meshmap'])).toBe(false);

      const accessControl3 = new ProviderUiAccessControl({
        restrictedAccess: { isMesheryUIRestricted: true },
        extensions: { navigator: [{ allowedTo: null }] },
      });
      expect(accessControl3.isExtensionComponentEnabled(['meshmap'])).toBe(false);
    });

    it('resolves extension component enabled state and handles non-array walker safely', () => {
      const accessControl = new ProviderUiAccessControl({
        restrictedAccess: { isMesheryUIRestricted: true },
        extensions: {
          navigator: [
            {
              allowedTo: {
                meshmap: {
                  designer: true,
                  visualizer: false,
                },
              },
            },
          ],
        },
      });

      expect(accessControl.isExtensionComponentEnabled(['meshmap', 'designer'])).toBe(true);
      expect(accessControl.isExtensionComponentEnabled(['meshmap', 'visualizer'])).toBe(false);
      expect(() => accessControl.isExtensionComponentEnabled(undefined)).not.toThrow();
      expect(() => accessControl.isExtensionComponentEnabled(null)).not.toThrow();
      expect(accessControl.isExtensionComponentEnabled(undefined)).toBe(false);
    });
  });
});
