import React, { useEffect } from 'react';
import { NoSsr } from '@layer5/sistent';
import { updatepagepath, } from '../../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';
import { useDispatchRtk } from '@/store/hooks';
import { updateTitle } from '@/store/slices/mesheryUi';

function DesignConfiguratorPage({ updatepagepath }) {
  const dispatch = useDispatchRtk();
  useEffect(() => {
    updatepagepath({ path: getPath(), isBeta: false, title: 'Configure Design' });
    dispatch(updateTitle({ title: 'Configure Design' }));
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
});

export default connect(null, mapDispatchToProps)(DesignConfiguratorPage);
