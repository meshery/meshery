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

  // Actions section: Create and Import design commands
  // Omitted if in playground mode or if the resolved design node is not shown
  if (!isPlayground && isDesignNavNodeShown(rawList)) {
    const createDesignPermission: PermissionKeySpec = {
      allOf: [Keys.CatalogManagementViewDesigns, Keys.CatalogManagementCreateNewDesign],
    };
    const createAllowed = isAllowed(createDesignPermission, can);
    actionsCommands.push({
      id: 'action:create-design',
      title: 'Create new Design',
      section: 'Actions',
      keywords: ['design', 'create', 'pattern', 'new', 'configurator'],
      href: '/configuration/designs/configurator',
      permission: createDesignPermission,
      disabled: !createAllowed,
      requires: !createAllowed ? formatRequires(createDesignPermission) : undefined,
    });

    const importDesignPermission: PermissionKeySpec = {
      allOf: [Keys.CatalogManagementViewDesigns, Keys.CatalogManagementImportDesign],
    };
    const importAllowed = isAllowed(importDesignPermission, can);
    actionsCommands.push({
      id: 'action:import-design',
      title: 'Import Design',
      section: 'Actions',
      keywords: ['design', 'import', 'upload', 'pattern'],
      href: '/configuration/designs?action=import',
      permission: importDesignPermission,
      disabled: !importAllowed,
      requires: !importAllowed ? formatRequires(importDesignPermission) : undefined,
    });
  }

  // Navigate section: two levels only (top-level nodes and their direct children)
  const seenHrefs = new Set<string>();

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

    // Top-level node: skipped if show is falsy
    if (!category.show) {
      continue;
    }

    const hasChildren = Array.isArray(category.children) && category.children.length > 0;

    if (hasChildren) {
      // 1. Compute shown direct children first
      const shownChildren: NavigatorChild[] = [];
      for (const child of category.children!) {
        if (
          providerUiAccessControl &&
          isPlayground &&
          !providerUiAccessControl.isNavigatorComponentEnabled([category.id, child.id])
        ) {
          continue;
        }

        // Child node: skipped only if typeof show !== 'undefined' && !show
        if (typeof child.show !== 'undefined' && !child.show) {
          continue;
        }

        if (!child.href) {
          continue;
        }

        shownChildren.push(child);
      }

      const shownChildHrefs = new Set<string>(shownChildren.map((c) => c.href!));

      // 2. Parent command: when link && href, unless its href equals the href of one of its shown children or already seen
      if (
        category.link &&
        category.href &&
        !shownChildHrefs.has(category.href) &&
        !seenHrefs.has(category.href)
      ) {
        seenHrefs.add(category.href);
        const allowed = isAllowed(category.permissionKey, can);
        navigateCommands.push({
          id: `nav:${category.href}`,
          title: category.title,
          section: 'Navigate',
          keywords: [],
          href: category.href,
          icon: category.icon,
          permission: category.permissionKey,
          disabled: !allowed,
          requires: !allowed ? formatRequires(category.permissionKey) : undefined,
        });
      }

      // 3. Direct children commands (order: parent first, then its children)
      for (const child of shownChildren) {
        if (seenHrefs.has(child.href!)) {
          continue;
        }
        seenHrefs.add(child.href!);

        const allowed = isAllowed(child.permissionKey, can);
        navigateCommands.push({
          id: `nav:${child.href}`,
          title: child.title,
          section: 'Navigate',
          subtitle: category.title,
          keywords: [category.title],
          href: child.href!,
          icon: child.icon,
          permission: child.permissionKey,
          disabled: !allowed,
          requires: !allowed ? formatRequires(child.permissionKey) : undefined,
        });
      }
    } else if (category.link && category.href) {
      // Node without children becomes one command
      if (!seenHrefs.has(category.href)) {
        seenHrefs.add(category.href);
        const allowed = isAllowed(category.permissionKey, can);
        navigateCommands.push({
          id: `nav:${category.href}`,
          title: category.title,
          section: 'Navigate',
          keywords: [],
          href: category.href,
          icon: category.icon,
          permission: category.permissionKey,
          disabled: !allowed,
          requires: !allowed ? formatRequires(category.permissionKey) : undefined,
        });
      }
    }
  }

  // Settings section: Settings and Preferences
  // Omitted if in playground mode
  if (!isPlayground) {
    const settingsAllowed = isAllowed(Keys.MesherySystemViewSettings, can);
    settingsCommands.push({
      id: 'settings',
      title: 'Settings',
      section: 'Settings',
      keywords: ['settings', 'configure', 'system'],
      href: '/settings',
      permission: Keys.MesherySystemViewSettings,
      disabled: !settingsAllowed,
      requires: !settingsAllowed ? formatRequires(Keys.MesherySystemViewSettings) : undefined,
    });

    settingsCommands.push({
      id: 'preferences',
      title: 'Preferences',
      section: 'Settings',
      keywords: ['preferences', 'user', 'profile'],
      href: '/user/preferences',
      permission: undefined,
      disabled: false,
      requires: undefined,
    });
  }

  return [...actionsCommands, ...navigateCommands, ...settingsCommands];
}
