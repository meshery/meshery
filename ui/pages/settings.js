import MesherySettings from '../components/MesherySettings';
import { NoSsr } from '@layer5/sistent';
import { updatepagepath, updatepagetitle } from '../lib/store';
import { connect } from 'react-redux';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from '../lib/path';
import React, { useEffect } from 'react';

function Settings(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
    props.updatepagetitle({ title: 'Settings' });
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Settings | Meshery</title>
      </Head>

      <MesherySettings />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
});

export default connect(null, mapDispatchToProps)(Settings);
