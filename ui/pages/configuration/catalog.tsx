import React from 'react';
import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '@/components/General/error-404';
import MesheryPatterns from '@/components/MesheryPatterns/MesheryPatterns';
import { usePageTitle } from '@/utils/hooks';

function CatalogPage() {
  usePageTitle('Catalog');

  return (
    <NoSsr>
      <Head>
        <title>Catalog | Meshery</title>
      </Head>
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
    </NoSsr>
  );
}

export default CatalogPage;
