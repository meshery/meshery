import React from 'react';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '@/components/general/error-404/index';
import MesheryPatterns from '@/components/designs/patterns/MesheryPatterns';
import { MesheryPage } from '@/components/MesheryPage';

function CatalogPage() {
  return (
    <MesheryPage title="Catalog">
      {CAN(Keys.CatalogManagementViewCatalog.id, Keys.CatalogManagementViewCatalog.function) ||
      false ? (
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
