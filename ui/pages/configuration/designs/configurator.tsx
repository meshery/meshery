import React from 'react';
import { useHasPermission } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import DesignConfigurator from '@/components/designs/configurator/MeshModel';
import DefaultError from '@/components/general/error-404';
import { MesheryPage } from '@/components/general/MesheryPage';

/**
 * Access gate for the design configurator.
 *
 * `CatalogManagementViewDesigns` is the same key that gates the designs list
 * (`components/designs/patterns/MesheryPatterns.tsx`) and the Designs navigator
 * entry, which are the only ways into this route. Gating here rather than
 * inside `DesignConfigurator` keeps the component from mounting at all, so the
 * unconditional `useMeshModelComponents` category fetch never fires for a user
 * who may not read designs.
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
