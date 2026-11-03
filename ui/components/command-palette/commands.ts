import type { ReactNode } from 'react';
import type { PermissionKeySpec } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import { DESIGN } from '@/constants/navigator';
import CAN from '@/utils/can';
import type { ProviderUiAccessControl } from '@/utils/disabledComponents';

export type CommandSection = 'Actions' | 'Navigate' | 'Settings';

export interface Command {
  id: string;
  title: string;
  section: CommandSection;
  subtitle?: string;
  keywords: string[];
  href: string;
  icon?: ReactNode;
  permission?: PermissionKeySpec;
  disabled: boolean;
  requires?: string;
}

export type CanFn = (action: string, subject: string) => boolean;

export interface NavigatorChild {
  id: string;
  title: string;
  href?: string;
  show?: boolean;
  icon?: ReactNode;
  link?: boolean;
  permissionKey?: PermissionKeySpec;
  children?: NavigatorChild[];
}

export interface NavigatorCategory {
  id: string;
  title: string;
  href?: string;
  show?: boolean;
  icon?: ReactNode;
  link?: boolean;
  submenu?: boolean;
  permissionKey?: PermissionKeySpec;
  children?: NavigatorChild[];
}

/**
 * Checks whether the given permission specification is satisfied by the `can` evaluator.
 * Supports a single Key, { anyOf: Key[] }, or { allOf: Key[] }.
 */
export function isAllowed(spec?: PermissionKeySpec, can: CanFn = CAN): boolean {
  if (!spec) {
    return true;
  }
  if ('anyOf' in spec) {
    return Array.isArray(spec.anyOf) && spec.anyOf.some((key) => can(key.id, key.function));
  }
  if ('allOf' in spec) {
    return Array.isArray(spec.allOf) && spec.allOf.every((key) => can(key.id, key.function));
  }
  return can(spec.id, spec.function);
}

/**
 * Produces a human-readable string of the required permission functions.
 */
export function formatRequires(spec?: PermissionKeySpec): string | undefined {
  if (!spec) {
    return undefined;
  }
  if ('anyOf' in spec) {
    return spec.anyOf.map((k) => k.function).join(' or ');
  }
  if ('allOf' in spec) {
    return spec.allOf.map((k) => k.function).join(' and ');
  }
  return spec.function;
}

function isDesignNavNodeShown(categories: NavigatorCategory[]): boolean {
  for (const cat of categories) {
    if (!cat || !cat.show) {
      continue;
    }
    if (cat.id === DESIGN) {
      return true;
    }
    if (Array.isArray(cat.children)) {
      for (const child of cat.children) {
        if (child.id === DESIGN) {
          return child.show !== false;
        }
      }
    }
  }
  return false;
}

type CommandBase = Omit<Command, 'permission' | 'disabled' | 'requires'>;

function withPermission(
  base: CommandBase,
  permission: PermissionKeySpec | undefined,
  can: CanFn,
): Command {
  const allowed = isAllowed(permission, can);
  return {
    ...base,
    permission,
    disabled: !allowed,
    requires: allowed ? undefined : formatRequires(permission),
  };
}

function toNavCommand(
  node: NavigatorCategory | NavigatorChild,
  href: string,
  can: CanFn,
  parent?: NavigatorCategory,
): Command {
  return withPermission(
    {
      id: `nav:${href}`,
      title: node.title,
      section: 'Navigate',
      subtitle: parent?.title,
      keywords: parent ? [parent.title] : [],
      href,
      icon: node.icon,
    },
    node.permissionKey,
    can,
  );
}

/**
 * Builds the complete list of commands for the command palette.
 * Takes the already-resolved navigator array and a ProviderUiAccessControl instance.
 * Does NOT import Navigator.tsx.
 *
 * Output section order: Actions, Navigate, Settings.
 */
export function buildCommands(
  navComponents: NavigatorCategory[] = [],
  providerUiAccessControl?: ProviderUiAccessControl | null,
  can: CanFn = CAN,
): Command[] {
  const rawList = Array.isArray(navComponents) ? navComponents : [];
  const isPlayground = Boolean(providerUiAccessControl?.isPlaygroundEnv);

  const actionsCommands: Command[] = [];
  const navigateCommands: Command[] = [];
  const settingsCommands: Command[] = [];

  // Actions: omitted in playground mode or when the resolved design node is not shown
  if (!isPlayground && isDesignNavNodeShown(rawList)) {
    actionsCommands.push(
      withPermission(
        {
          id: 'action:create-design',
          title: 'Create new Design',
          section: 'Actions',
          keywords: ['design', 'create', 'pattern', 'new', 'configurator'],
          href: '/configuration/designs/configurator',
        },
        { allOf: [Keys.CatalogManagementViewDesigns, Keys.CatalogManagementCreateNewDesign] },
        can,
      ),
      withPermission(
        {
          id: 'action:import-design',
          title: 'Import Design',
          section: 'Actions',
          keywords: ['design', 'import', 'upload', 'pattern'],
          href: '/configuration/designs?action=import',
        },
        { allOf: [Keys.CatalogManagementViewDesigns, Keys.CatalogManagementImportDesign] },
        can,
      ),
    );
  }

  // Navigate: two levels only (top-level nodes and their direct children), deduped by href
  const seenHrefs = new Set<string>();
  const addNavCommand = (
    node: NavigatorCategory | NavigatorChild,
    href: string,
    parent?: NavigatorCategory,
  ) => {
    if (seenHrefs.has(href)) {
      return;
    }
    seenHrefs.add(href);
    navigateCommands.push(toNavCommand(node, href, can, parent));
  };

  for (const category of rawList) {
    if (!category) {
      continue;
    }
    if (
      providerUiAccessControl &&
      !providerUiAccessControl.isNavigatorComponentEnabled([category.id])
    ) {
      continue;
    }
    if (!category.show) {
      continue;
    }

    if (Array.isArray(category.children) && category.children.length > 0) {
      const shownChildren = category.children.filter((child) => {
        if (
          providerUiAccessControl &&
          isPlayground &&
          !providerUiAccessControl.isNavigatorComponentEnabled([category.id, child.id])
        ) {
          return false;
        }
        // A child is skipped only when show is explicitly falsy; undefined show is included
        if (typeof child.show !== 'undefined' && !child.show) {
          return false;
        }
        return Boolean(child.href);
      });
      const shownChildHrefs = new Set(shownChildren.map((child) => child.href));

      // Parent command only when it links somewhere its shown children do not already go
      if (category.link && category.href && !shownChildHrefs.has(category.href)) {
        addNavCommand(category, category.href);
      }
      for (const child of shownChildren) {
        addNavCommand(child, child.href!, category);
      }
    } else if (category.link && category.href) {
      addNavCommand(category, category.href);
    }
  }

  // Settings: omitted in playground mode
  if (!isPlayground) {
    settingsCommands.push(
      withPermission(
        {
          id: 'settings',
          title: 'Settings',
          section: 'Settings',
          keywords: ['settings', 'configure', 'system'],
          href: '/settings',
        },
        Keys.MesherySystemViewSettings,
        can,
      ),
      withPermission(
        {
          id: 'preferences',
          title: 'Preferences',
          section: 'Settings',
          keywords: ['preferences', 'user', 'profile'],
          href: '/user/preferences',
        },
        undefined,
        can,
      ),
    );
  }

  return [...actionsCommands, ...navigateCommands, ...settingsCommands];
}
