import React, { useEffect } from 'react';
import { NoSsr, withStyles } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { updatepagepath } from '../../lib/store';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from '@/components/General/error-404';
import MesheryPatterns from '@/components/MesheryPatterns';

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
};

function CatalogPage(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, []);

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

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(CatalogPage));
