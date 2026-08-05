import React from 'react';
import { VISIBILITY } from '../../utils/Enum';
import { Keys } from '@meshery/schemas/permissions';
import { useHasPermission } from '@sistent/sistent';
import DefaultError from '@/components/general/error-404/index';
import MesheryPatterns from '@/components/designs/patterns/MesheryPatterns';
import { MesheryPage } from '@/components/general/MesheryPage';

function CatalogPage() {
  const hasPermission = useHasPermission(Keys.CatalogManagementViewCatalog);
  return (
    <MesheryPage title="Catalog">
      {hasPermission ? (
        <MesheryPatterns
          disableCreateImportDesignButton={true}
          disableUniversalFilter={true}
          initialFilters={{ visibility: VISIBILITY.PUBLISHED }}
          hideVisibility={true}
          pageTitle="Catalog"
          arePatternsReadOnly={true}
        />
      ) : (
        <DefaultError />
      )}
    </MesheryPage>
  );
}

export default CatalogPage;
