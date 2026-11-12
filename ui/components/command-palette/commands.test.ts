import { beforeEach, describe, expect, it, vi } from 'vitest';
import _ from 'lodash';
import CAN, { ability } from '@/utils/can';
import { ProviderUiAccessControl } from '@/utils/disabledComponents';
import { Keys } from '@meshery/schemas/permissions';
import { buildCommands, formatRequires, isAllowed, type NavigatorCategory } from './commands';
import { getNavigatorComponents } from '../layout/Navigator/navigatorComponents';

// Transitive import of can.ts -> store and Sistent
vi.mock('@/store', () => ({
  store: { getState: () => ({}) },
}));

vi.mock('../../css/icons.styles', () => ({
  iconSmall: {},
}));

// Stub the 7 static .js drawer-icon SVG assets that Vite cannot parse as JSX
vi.mock('../../public/static/img/drawer-icons/lifecycle_mgmt_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/extensions_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/pattern_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/lifecycle_hover_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/performance_hover_svg', () => ({
  default: () => null,
}));
vi.mock('../../public/static/img/drawer-icons/configuration_hover_svg', () => ({
  default: () => null,
}));

const testTheme = {
  palette: {
    background: { constant: { white: 'white' } },
    icon: { default: 'black' },
  },
};

const makeNavNode = (
  title: string,
  href?: string,
  id: string = title.toLowerCase().replace(/\s+/g, '-'),
  opts?:
    | NavigatorCategory['children']
    | {
        icon?: NavigatorCategory['icon'];
        permission?: NavigatorCategory['permissionKey'];
        children?: NavigatorCategory['children'];
        show?: boolean;
        link?: boolean;
      },
): NavigatorCategory => {
  const isArr = Array.isArray(opts);
  const conf = isArr ? { children: opts } : opts || {};
  return {
    id,
    title,
    href,
    link: conf.link ?? true,
    ...('show' in conf ? (conf.show !== undefined ? { show: conf.show } : {}) : { show: true }),
    ...(conf.children ? { children: conf.children } : {}),
    ...(conf.icon ? { icon: conf.icon } : {}),
    ...(conf.permission ? { permissionKey: conf.permission } : {}),
  };
};

describe('commands.ts', () => {
  beforeEach(() => {
    ability.update([]);
  });

  describe('isAllowed and formatRequires', () => {
    it('returns true when spec is undefined or empty', () => {
      expect(isAllowed(undefined, CAN)).toBe(true);
      expect(formatRequires(undefined)).toBeUndefined();
    });

    it('handles a single permission Key (allowed vs denied)', () => {
      const key = Keys.WorkspaceManagementViewConnections;
      expect(isAllowed(key, CAN)).toBe(false);
      expect(formatRequires(key)).toBe(key.function);

      ability.update([{ action: key.id, subject: _.lowerCase(key.function) }]);
      expect(isAllowed(key, CAN)).toBe(true);
    });

    it('handles anyOf (permitted when at least one key is held)', () => {
      const keyA = Keys.WorkspaceManagementViewConnections;
      const keyB = Keys.SecurityManagementViewCredentials;
      const anyOfSpec = { anyOf: [keyA, keyB] };

      expect(isAllowed(anyOfSpec, CAN)).toBe(false);
      expect(formatRequires(anyOfSpec)).toBe(`${keyA.function} or ${keyB.function}`);

      // Grant only keyB
      ability.update([{ action: keyB.id, subject: _.lowerCase(keyB.function) }]);
      expect(isAllowed(anyOfSpec, CAN)).toBe(true);
    });

    it('handles allOf (denied when even one key is missing)', () => {
      const keyA = Keys.CatalogManagementViewDesigns;
      const keyB = Keys.CatalogManagementCreateNewDesign;
      const allOfSpec = { allOf: [keyA, keyB] };

      expect(isAllowed(allOfSpec, CAN)).toBe(false);
      expect(formatRequires(allOfSpec)).toBe(`${keyA.function} and ${keyB.function}`);

      // Grant only keyA -> still denied
      ability.update([{ action: keyA.id, subject: _.lowerCase(keyA.function) }]);
      expect(isAllowed(allOfSpec, CAN)).toBe(false);

      // Grant both keyA and keyB -> allowed
      ability.update([
        { action: keyA.id, subject: _.lowerCase(keyA.function) },
        { action: keyB.id, subject: _.lowerCase(keyB.function) },
      ]);
      expect(isAllowed(allOfSpec, CAN)).toBe(true);
    });
  });

  describe('buildCommands logic', () => {
    const defaultPayload = {
      capabilities: [],
      restrictedAccess: {
        isMesheryUIRestricted: false,
      },
    };

    it.each([
      {
        name: 'denied commands remain listed with disabled: true and requires set',
        abilities: [],
        disabled: true,
        settingsRequires: Keys.MesherySystemViewSettings.function,
        connectionsRequires: Keys.WorkspaceManagementViewConnections.function,
      },
      {
        name: 'allowed commands have disabled: false and requires: undefined',
        abilities: [{ action: 'manage', subject: 'all' }],
        disabled: false,
        settingsRequires: undefined,
        connectionsRequires: undefined,
      },
    ])('$name', ({ abilities, disabled, settingsRequires, connectionsRequires }) => {
      ability.update(abilities);
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const commands = buildCommands(nav, accessControl, CAN);

      const settingsCmd = commands.find((c) => c.id === 'settings');
      expect(settingsCmd).toBeDefined();
      expect(settingsCmd?.disabled).toBe(disabled);
      expect(settingsCmd?.requires).toBe(settingsRequires);

      const connectionsCmd = commands.find((c) => c.href === '/management/connections');
      expect(connectionsCmd).toBeDefined();
      expect(connectionsCmd?.disabled).toBe(disabled);
      expect(connectionsCmd?.requires).toBe(connectionsRequires);
    });

    const viewKey = Keys.CatalogManagementViewDesigns;
    const createKey = Keys.CatalogManagementCreateNewDesign;
    const importKey = Keys.CatalogManagementImportDesign;
    const toGrant = (key: (typeof Keys)[keyof typeof Keys]) => ({
      action: key.id,
      subject: _.lowerCase(key.function),
    });

    it.each([
      {
        name: '1. All denied',
        grants: [],
        createDisabled: true,
        importDisabled: true,
      },
      {
        name: '2. Granted only ViewDesigns: Create disabled, requires lists CreateNewDesign; Import disabled, requires lists ImportDesign',
        grants: [toGrant(viewKey)],
        createDisabled: true,
        createRequires: createKey.function,
        importDisabled: true,
        importRequires: importKey.function,
      },
      {
        name: '3. Granted only CreateNewDesign: disabled',
        grants: [toGrant(createKey)],
        createDisabled: true,
      },
      {
        name: '4. Granted only ImportDesign: disabled',
        grants: [toGrant(importKey)],
        importDisabled: true,
      },
      {
        name: '5. Granted both ViewDesigns and CreateNewDesign: Create allowed',
        grants: [toGrant(viewKey), toGrant(createKey)],
        createDisabled: false,
        createRequiresUndefined: true,
      },
      {
        name: '6. Granted both ViewDesigns and ImportDesign: Import allowed',
        grants: [toGrant(viewKey), toGrant(importKey)],
        importDisabled: false,
        importRequiresUndefined: true,
      },
    ])(
      'evaluates Create/Import permissions with partial and full key grants: $name',
      ({
        grants,
        createDisabled,
        createRequires,
        createRequiresUndefined,
        importDisabled,
        importRequires,
        importRequiresUndefined,
      }) => {
        ability.update(grants);
        const accessControl = new ProviderUiAccessControl(defaultPayload);
        const nav = getNavigatorComponents(accessControl, testTheme);
        const commands = buildCommands(nav, accessControl, CAN);

        if (createDisabled !== undefined) {
          const createCmd = commands.find((c) => c.id === 'action:create-design');
          expect(createCmd?.disabled).toBe(createDisabled);
          if (createRequires) expect(createCmd?.requires).toContain(createRequires);
          if (createRequiresUndefined) expect(createCmd?.requires).toBeUndefined();
        }

        if (importDisabled !== undefined) {
          const importCmd = commands.find((c) => c.id === 'action:import-design');
          expect(importCmd?.disabled).toBe(importDisabled);
          if (importRequires) expect(importCmd?.requires).toContain(importRequires);
          if (importRequiresUndefined) expect(importCmd?.requires).toBeUndefined();
        }
      },
    );

    it('real getNavigatorComponents output includes commands for /performance and /telemetry, deduping /management/connections and /configuration/designs', () => {
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);
      const commands = buildCommands(nav, accessControl, CAN);

      // Parent commands for /performance and /telemetry appear
      expect(commands.some((c) => c.href === '/performance')).toBe(true);
      expect(commands.some((c) => c.href === '/telemetry')).toBe(true);

      // Still lists /management/connections exactly once and /configuration/designs exactly once
      const connectionsCommands = commands.filter((c) => c.href === '/management/connections');
      expect(connectionsCommands).toHaveLength(1);

      const designsCommands = commands.filter((c) => c.href === '/configuration/designs');
      expect(designsCommands).toHaveLength(1);
    });

    it('parent without link is not a command; parent with link whose href equals a shown child href is deduped', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('Parent No Link', '/no-link-parent', 'parent-no-link', {
          link: false,
          children: [makeNavNode('Child 1', '/child-1', 'child-1')],
        }),
        makeNavNode('Parent Dup', '/duplicate-href', 'parent-dup', [
          makeNavNode('Child Dup', '/duplicate-href', 'child-dup'),
        ]),
        makeNavNode('Parent Distinct', '/distinct-parent', 'parent-with-distinct-href', [
          makeNavNode('Child Distinct', '/distinct-child', 'child-distinct'),
        ]),
      ];

      const commands = buildCommands(mockNav, null, CAN);

      // Parent without link is not a command
      expect(commands.some((c) => c.href === '/no-link-parent')).toBe(false);
      expect(commands.some((c) => c.href === '/child-1')).toBe(true);

      // Parent with link whose href equals a shown child href is deduped (only child appears)
      const dupMatches = commands.filter((c) => c.href === '/duplicate-href');
      expect(dupMatches).toHaveLength(1);
      expect(dupMatches[0].subtitle).toBe('Parent Dup');

      // Parent with link whose href is distinct becomes a command before its child
      const distinctParentIdx = commands.findIndex((c) => c.href === '/distinct-parent');
      const distinctChildIdx = commands.findIndex((c) => c.href === '/distinct-child');
      expect(distinctParentIdx).toBeGreaterThanOrEqual(0);
      expect(distinctChildIdx).toBeGreaterThan(distinctParentIdx);
    });

    it('childless top-level node with link: false is not a command', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('No Link Leaf', '/no-link-leaf', 'no-link-leaf', { link: false }),
        makeNavNode('Link Leaf', '/link-leaf', 'link-leaf'),
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/no-link-leaf')).toBe(false);
      expect(commands.some((c) => c.href === '/link-leaf')).toBe(true);
    });

    it('skips top-level nodes if show is falsy', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('Shown', '/shown', 'shown-category'),
        makeNavNode('Hidden', undefined, 'hidden-category', {
          show: false,
          children: [makeNavNode('Hidden Child', '/hidden-child', 'child')],
        }),
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/shown')).toBe(true);
      expect(commands.some((c) => c.href === '/hidden-child')).toBe(false);
    });

    it('skips child only when typeof show !== "undefined" && !show (undefined show is included)', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('Category', undefined, 'cat', [
          makeNavNode('Hidden Child', '/explicitly-hidden', 'explicitly-hidden', { show: false }),
          makeNavNode('Shown Child', '/explicitly-shown', 'explicitly-shown', { show: true }),
          makeNavNode('Undefined Show Child', '/undefined-show', 'undefined-show', {
            show: undefined,
          }),
        ]),
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/explicitly-hidden')).toBe(false);
      expect(commands.some((c) => c.href === '/explicitly-shown')).toBe(true);
      expect(commands.some((c) => c.href === '/undefined-show')).toBe(true);
    });

    it('does not recurse into grandchildren', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('Lifecycle', '/management/adapter', 'lifecycle', [
          makeNavNode('Adapters', '/management/adapter', 'service_mesh', [
            makeNavNode('Istio Adapter Child', '/management/adapter/istio', 'istio'),
          ]),
        ]),
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/management/adapter')).toBe(true);
      expect(commands.some((c) => c.href === '/management/adapter/istio')).toBe(false);
    });

    it('dedupes commands by href', () => {
      const mockNav: NavigatorCategory[] = [
        makeNavNode('Cat 1', undefined, 'cat-1', [
          makeNavNode('Child 1', '/same-href', 'child-1'),
          makeNavNode('Child 2', '/same-href', 'child-2'),
        ]),
      ];

      const commands = buildCommands(mockNav, null, CAN);
      const matches = commands.filter((c) => c.href === '/same-href');
      expect(matches).toHaveLength(1);
    });

    it('restricted playground mode adds no extras and restricts navigation', () => {
      const playgroundPayload = {
        restrictedAccess: {
          isMesheryUIRestricted: true,
          allowedComponents: {
            navigator: {
              dashboard: true,
            },
          },
        },
      };

      const accessControl = new ProviderUiAccessControl(playgroundPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const commands = buildCommands(nav, accessControl, CAN);

      // No extras
      expect(commands.some((c) => c.id === 'action:create-design')).toBe(false);
      expect(commands.some((c) => c.id === 'action:import-design')).toBe(false);
      expect(commands.some((c) => c.id === 'settings')).toBe(false);
      expect(commands.some((c) => c.id === 'preferences')).toBe(false);

      // Only allowed nav
      expect(commands.some((c) => c.href === '/')).toBe(true);
      expect(commands.some((c) => c.href === '/management/connections')).toBe(false);
    });

    it('handles null/undefined nav input or payload safely', () => {
      expect(() => buildCommands(null as any, null)).not.toThrow();
      expect(() => buildCommands(undefined, null)).not.toThrow();
      const safe = buildCommands([], new ProviderUiAccessControl(null));
      expect(Array.isArray(safe)).toBe(true);
    });

    it('omits Create/Import design actions when design nav node is hidden or not shown', () => {
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      // Hide the configuration category or the design child node
      const navWithoutDesign = nav.map((cat) => {
        if (cat.id === 'configuration') {
          return {
            ...cat,
            children: cat.children?.filter((child) => child.id !== 'design'),
          };
        }
        return cat;
      });

      const commands = buildCommands(navWithoutDesign, accessControl, CAN);
      expect(commands.some((c) => c.id === 'action:create-design')).toBe(false);
      expect(commands.some((c) => c.id === 'action:import-design')).toBe(false);
    });

    it('enforces exact section order: Actions, Navigate, Settings', () => {
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const commands = buildCommands(nav, accessControl, CAN);

      const sections = commands.map((c) => c.section);
      const actionsIdx = sections.indexOf('Actions');
      const navigateIdx = sections.indexOf('Navigate');
      const settingsIdx = sections.indexOf('Settings');

      expect(actionsIdx).toBeGreaterThanOrEqual(0);
      expect(navigateIdx).toBeGreaterThan(actionsIdx);
      expect(settingsIdx).toBeGreaterThan(navigateIdx);
    });
  });

  describe('anti-drift test against real getNavigatorComponents', () => {
    it('every non-hidden top-level or direct child href appears in built commands', () => {
      const payload = {
        capabilities: [],
        restrictedAccess: {
          isMesheryUIRestricted: false,
        },
      };

      const accessControl = new ProviderUiAccessControl(payload);
      const nav = getNavigatorComponents(accessControl, testTheme);
      const commands = buildCommands(nav, accessControl, CAN);
      const builtHrefs = new Set(commands.map((c) => c.href));

      expect(builtHrefs.size).toBeGreaterThan(0);

      for (const category of nav) {
        if (!category || !category.show) {
          continue;
        }

        if (Array.isArray(category.children) && category.children.length > 0) {
          const shownChildren = category.children.filter(
            (child) => typeof child.show === 'undefined' || child.show,
          );
          const shownChildHrefs = new Set(shownChildren.map((c) => c.href).filter(Boolean));

          // Top-level parent with children: must be present unless its href equals a shown child's href
          if (category.link && category.href && !shownChildHrefs.has(category.href)) {
            expect(builtHrefs.has(category.href)).toBe(true);
          }

          for (const child of shownChildren) {
            if (child.href) {
              expect(builtHrefs.has(child.href)).toBe(true);
            }
          }
        } else if (category.link && category.href) {
          expect(builtHrefs.has(category.href)).toBe(true);
        }
      }
    });
  });
});
