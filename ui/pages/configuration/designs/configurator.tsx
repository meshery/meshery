import React from 'react';
import { useHasPermission } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import DesignConfigurator from '@/components/designs/configurator/MeshModel';
import DefaultError from '@/components/general/error-404';
import { MesheryPage } from '@/components/general/MesheryPage';

/**
 * Access gate for the design configurator.
 *
 * Gated on `CatalogManagementViewDesigns` - the same key that gates the designs
 * list and the Designs navigator entry. The configurator reads and writes full
 * design JSON, so a read-only browsing key is deliberately not accepted here;
 * that would grant more than the designs list itself does.
 *
 * Several controls outside the designs page also navigate here, some gated on
 * other keys, so a session can reach this page and be denied. No built-in
 * keychain can: where a custom one does, the fix belongs at the entry-point
 * controls, not in widening this gate. The call-site analysis behind both
 * decisions is in https://github.com/meshery/meshery/issues/21492, which is the
 * place to update when those controls move.
 *
 * Gating here rather than inside `DesignConfigurator` keeps the component from
 * mounting at all, so its unconditional model-category fetch never fires for a
 * user who may not read designs.
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
