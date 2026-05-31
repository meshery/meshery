import React from 'react';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '@/components/general/error-404/index';
import MesheryPatterns from '@/components/designs/patterns/MesheryPatterns';
import { MesheryPage } from '@/components/MesheryPage';

function CatalogPage() {
  return (
    <MesheryPage title="Catalog">
      {CAN(keys.VIEW_CATALOG.action, keys.VIEW_CATALOG.subject) || false ? (
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
