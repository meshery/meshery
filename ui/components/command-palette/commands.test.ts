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

vi.mock('@sistent/sistent', () => ({
  createCanShow: () => () => null,
  EventBus: class {
    on() {
      return { subscribe: () => {} };
    }
  },
  useTheme: () => ({
    palette: {
      background: { constant: { white: 'white' } },
      icon: { default: 'black' },
    },
  }),
  CatalogIcon: () => null,
  TachographDigitalIcon: () => null,
  InsertChartIcon: () => null,
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

    it('denied commands remain listed with disabled: true and requires set', () => {
      ability.update([]); // All denied
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const commands = buildCommands(nav, accessControl, CAN);

      const settingsCmd = commands.find((c) => c.id === 'settings');
      expect(settingsCmd).toBeDefined();
      expect(settingsCmd?.disabled).toBe(true);
      expect(settingsCmd?.requires).toBe(Keys.MesherySystemViewSettings.function);

      const connectionsCmd = commands.find((c) => c.href === '/management/connections');
      expect(connectionsCmd).toBeDefined();
      expect(connectionsCmd?.disabled).toBe(true);
      expect(connectionsCmd?.requires).toBe(Keys.WorkspaceManagementViewConnections.function);
    });

    it('allowed commands have disabled: false and requires: undefined', () => {
      ability.update([{ action: 'manage', subject: 'all' }]);
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const commands = buildCommands(nav, accessControl, CAN);

      const settingsCmd = commands.find((c) => c.id === 'settings');
      expect(settingsCmd?.disabled).toBe(false);
      expect(settingsCmd?.requires).toBeUndefined();

      const connectionsCmd = commands.find((c) => c.href === '/management/connections');
      expect(connectionsCmd?.disabled).toBe(false);
      expect(connectionsCmd?.requires).toBeUndefined();
    });

    it('evaluates Create/Import permissions with partial and full key grants', () => {
      const accessControl = new ProviderUiAccessControl(defaultPayload);
      const nav = getNavigatorComponents(accessControl, testTheme);

      const viewKey = Keys.CatalogManagementViewDesigns;
      const createKey = Keys.CatalogManagementCreateNewDesign;
      const importKey = Keys.CatalogManagementImportDesign;

      // 1. All denied
      ability.update([]);
      let commands = buildCommands(nav, accessControl, CAN);
      let createCmd = commands.find((c) => c.id === 'action:create-design');
      let importCmd = commands.find((c) => c.id === 'action:import-design');
      expect(createCmd?.disabled).toBe(true);
      expect(importCmd?.disabled).toBe(true);

      // 2. Granted only ViewDesigns: Create disabled, requires lists CreateNewDesign; Import disabled, requires lists ImportDesign
      ability.update([{ action: viewKey.id, subject: _.lowerCase(viewKey.function) }]);
      commands = buildCommands(nav, accessControl, CAN);
      createCmd = commands.find((c) => c.id === 'action:create-design');
      importCmd = commands.find((c) => c.id === 'action:import-design');
      expect(createCmd?.disabled).toBe(true);
      expect(createCmd?.requires).toContain(createKey.function);
      expect(importCmd?.disabled).toBe(true);
      expect(importCmd?.requires).toContain(importKey.function);

      // 3. Granted only CreateNewDesign: disabled
      ability.update([{ action: createKey.id, subject: _.lowerCase(createKey.function) }]);
      commands = buildCommands(nav, accessControl, CAN);
      createCmd = commands.find((c) => c.id === 'action:create-design');
      expect(createCmd?.disabled).toBe(true);

      // 4. Granted only ImportDesign: disabled
      ability.update([{ action: importKey.id, subject: _.lowerCase(importKey.function) }]);
      commands = buildCommands(nav, accessControl, CAN);
      importCmd = commands.find((c) => c.id === 'action:import-design');
      expect(importCmd?.disabled).toBe(true);

      // 5. Granted both ViewDesigns and CreateNewDesign: Create allowed
      ability.update([
        { action: viewKey.id, subject: _.lowerCase(viewKey.function) },
        { action: createKey.id, subject: _.lowerCase(createKey.function) },
      ]);
      commands = buildCommands(nav, accessControl, CAN);
      createCmd = commands.find((c) => c.id === 'action:create-design');
      expect(createCmd?.disabled).toBe(false);
      expect(createCmd?.requires).toBeUndefined();

      // 6. Granted both ViewDesigns and ImportDesign: Import allowed
      ability.update([
        { action: viewKey.id, subject: _.lowerCase(viewKey.function) },
        { action: importKey.id, subject: _.lowerCase(importKey.function) },
      ]);
      commands = buildCommands(nav, accessControl, CAN);
      importCmd = commands.find((c) => c.id === 'action:import-design');
      expect(importCmd?.disabled).toBe(false);
      expect(importCmd?.requires).toBeUndefined();
    });

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
        {
          id: 'parent-no-link',
          title: 'Parent No Link',
          href: '/no-link-parent',
          link: false,
          show: true,
          children: [
            {
              id: 'child-1',
              title: 'Child 1',
              href: '/child-1',
              link: true,
              show: true,
            },
          ],
        },
        {
          id: 'parent-dup',
          title: 'Parent Dup',
          href: '/duplicate-href',
          link: true,
          show: true,
          children: [
            {
              id: 'child-dup',
              title: 'Child Dup',
              href: '/duplicate-href',
              link: true,
              show: true,
            },
          ],
        },
        {
          id: 'parent-with-distinct-href',
          title: 'Parent Distinct',
          href: '/distinct-parent',
          link: true,
          show: true,
          children: [
            {
              id: 'child-distinct',
              title: 'Child Distinct',
              href: '/distinct-child',
              link: true,
              show: true,
            },
          ],
        },
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
        {
          id: 'no-link-leaf',
          title: 'No Link Leaf',
          href: '/no-link-leaf',
          link: false,
          show: true,
        },
        {
          id: 'link-leaf',
          title: 'Link Leaf',
          href: '/link-leaf',
          link: true,
          show: true,
        },
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/no-link-leaf')).toBe(false);
      expect(commands.some((c) => c.href === '/link-leaf')).toBe(true);
    });

    it('skips top-level nodes if show is falsy', () => {
      const mockNav: NavigatorCategory[] = [
        {
          id: 'shown-category',
          title: 'Shown',
          show: true,
          link: true,
          href: '/shown',
        },
        {
          id: 'hidden-category',
          title: 'Hidden',
          show: false,
          link: true,
          children: [
            {
              id: 'child',
              title: 'Hidden Child',
              href: '/hidden-child',
              link: true,
              show: true,
            },
          ],
        },
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/shown')).toBe(true);
      expect(commands.some((c) => c.href === '/hidden-child')).toBe(false);
    });

    it('skips child only when typeof show !== "undefined" && !show (undefined show is included)', () => {
      const mockNav: NavigatorCategory[] = [
        {
          id: 'cat',
          title: 'Category',
          show: true,
          link: true,
          children: [
            {
              id: 'explicitly-hidden',
              title: 'Hidden Child',
              href: '/explicitly-hidden',
              link: true,
              show: false,
            },
            {
              id: 'explicitly-shown',
              title: 'Shown Child',
              href: '/explicitly-shown',
              link: true,
              show: true,
            },
            {
              id: 'undefined-show',
              title: 'Undefined Show Child',
              href: '/undefined-show',
              link: true,
              // show is undefined
            },
          ],
        },
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/explicitly-hidden')).toBe(false);
      expect(commands.some((c) => c.href === '/explicitly-shown')).toBe(true);
      expect(commands.some((c) => c.href === '/undefined-show')).toBe(true);
    });

    it('does not recurse into grandchildren', () => {
      const mockNav: NavigatorCategory[] = [
        {
          id: 'lifecycle',
          title: 'Lifecycle',
          show: true,
          link: true,
          href: '/management/adapter',
          children: [
            {
              id: 'service_mesh',
              title: 'Adapters',
              href: '/management/adapter',
              link: true,
              show: true,
              children: [
                {
                  id: 'istio',
                  title: 'Istio Adapter Child',
                  href: '/management/adapter/istio',
                  link: true,
                  show: true,
                },
              ],
            },
          ],
        },
      ];

      const commands = buildCommands(mockNav, null, CAN);
      expect(commands.some((c) => c.href === '/management/adapter')).toBe(true);
      expect(commands.some((c) => c.href === '/management/adapter/istio')).toBe(false);
    });

    it('dedupes commands by href', () => {
      const mockNav: NavigatorCategory[] = [
        {
          id: 'cat-1',
          title: 'Cat 1',
          show: true,
          link: true,
          children: [
            {
              id: 'child-1',
              title: 'Child 1',
              href: '/same-href',
              link: true,
              show: true,
            },
            {
              id: 'child-2',
              title: 'Child 2',
              href: '/same-href',
              link: true,
              show: true,
            },
          ],
        },
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
