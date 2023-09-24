import MesherySettings from '../components/MesherySettings';
import { NoSsr, withStyles } from '@material-ui/core';
import { updatepagepath, updatepagetitle } from '../lib/store';
import { connect } from 'react-redux';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from '../lib/path';
import React, { useEffect } from 'react';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

function Settings(props) {
  useEffect(() => {
    console.log(`path: ${getPath()}`);
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
