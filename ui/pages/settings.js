import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { NoSsr, withStyles } from '@material-ui/core';
import MesherySettings from '../components/MesherySettings';
import { getPath } from '../lib/path';
import { updatepagepath, updatepagetitle } from '../lib/store';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

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

export default withStyles(styles)(connect(null, mapDispatchToProps)(Settings));
