import React, { useEffect } from 'react';
import { NoSsr } from '@mui/material';
import { updatebetabadge, updatepagepath, updatepagetitle } from '../../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';

function DesignConfiguratorPage({ updatepagepath, updatepagetitle, updatebetabadge }) {
  useEffect(() => {
    updatepagepath({ path: getPath(), isBeta: true, title: 'Configure Design' });
    updatepagetitle({ title: 'Configure Design' });
    updatebetabadge({ isBeta: true });
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Designs Configurator</title>
      </Head>
      <DesignConfigurator />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
  updatebetabadge: bindActionCreators(updatebetabadge, dispatch),
});

export default connect(null, mapDispatchToProps)(DesignConfiguratorPage);
