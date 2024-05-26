import React, { useEffect } from 'react';
import { NoSsr, withStyles } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { updatepagepath } from '../../lib/store';
import MesheryPatterns from '@/components/MesheryPatterns';
import { VISIBILITY } from '../../utils/Enum';

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
      <MesheryPatterns
        disableCreateImportDesignButton={true}
        disableUniversalFilter={true}
        initialFilters={{ visibility: VISIBILITY.PUBLISHED }}
        pageTitle="Catalog"
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(CatalogPage));
