import React from 'react';
import { useHasPermission } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import DesignConfigurator from '@/components/designs/configurator/MeshModel';
import DefaultError from '@/components/general/error-404';
import { MesheryPage } from '@/components/general/MesheryPage';

/**
 * Access gate for the design configurator.
 *
 * Four call sites navigate here:
 *
 * 1. `components/designs/patterns/MesheryPatterns.tsx:365` - the designs list
 *    row action.
 * 2. `components/designs/patterns/MesheryPatternsToolbar.tsx:61` - the create
 *    new design toolbar button.
 * 3. `components/designs/patterns/MesheryPatternCard.tsx:103` - "Edit In
 *    Configurator" on the design card, reachable from `/configuration/catalog`,
 *    which gates on `CatalogManagementViewCatalog` rather than ViewDesigns.
 *    `arePatternsReadOnly` does not suppress that button: the `isReadOnly`
 *    guard at `MesheryPatternCard.tsx:441` covers a different block, so the
 *    button renders whenever `userCanEdit` (holds `CatalogManagementEditDesign`,
 *    or owns the design).
 * 4. `components/workspaces/SpacesSwitcher/MainDesignsContent.tsx:176` - the
 *    workspace design open, which falls back to the configurator when the
 *    Kanvas designer is unavailable.
 *
 * `CatalogManagementViewDesigns` is the same key that gates the first two: the
 * designs list (`MesheryPatterns.tsx:72`) and the Designs navigator entry
 * (`components/layout/Navigator/navigatorComponents.tsx:159`). It remains the
 * single correct key despite the last two, because no built-in keychain in
 * `server/permissions/keys.csv` holds `Edit Design`, `View Catalog` or
 * `View Workspace` without also holding `View Designs` - verified per role
 * column across User, Team Admin, Academy Admin, Leaner, Workspace Admin, Org
 * Billing Manager, Org Admin and Provider Admin. No default role can reach a
 * deny page from those buttons; only a custom keychain can.
 *
 * Where a custom keychain does make that dead end reachable, the fix belongs at
 * the entry-point buttons - a button that leads to a permission-denied page
 * should not render - and not in widening this gate. The configurator reads and
 * writes full design JSON, whereas `CatalogManagementViewCatalog` grants only
 * read-only published browsing, so accepting it here would grant more than the
 * designs list itself does. That button-side fix is deliberately out of scope
 * here and raised as a separate product question.
 *
 * Gating here rather than inside `DesignConfigurator` keeps the component from
 * mounting at all, so the unconditional `useMeshModelComponents` category fetch
 * never fires for a user who may not read designs.
 *
 * This is a presentation control. Authorization is decided by the configured
 * provider and enforced per handler on Meshery Server; the gate exists so the
 * UI is consistent about what an unpermitted session is shown.
 */
function DesignConfiguratorPage() {
  const canViewDesigns = useHasPermission(Keys.CatalogManagementViewDesigns);

  return (
    <MesheryPage title="Configure Design" headTitle="Designs Configurator">
      {canViewDesigns ? (
        <DesignConfigurator />
      ) : (
        <DefaultError permissionKey={Keys.CatalogManagementViewDesigns} />
      )}
    </MesheryPage>
  );
}

export default DesignConfiguratorPage;
