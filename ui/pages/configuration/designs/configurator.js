import React, { useEffect } from 'react';
import { NoSsr } from '@layer5/sistent';
import { updatepagepath, updatepagetitle } from '../../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';

function DesignConfiguratorPage({ updatepagepath, updatepagetitle }) {
  useEffect(() => {
    updatepagepath({ path: getPath(), isBeta: false, title: 'Configure Design' });
    updatepagetitle({ title: 'Configure Design' });
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
});

export default connect(null, mapDispatchToProps)(DesignConfiguratorPage);
