import React, { useEffect } from 'react';
import { Paper } from '@layer5/sistent';
import { NoSsr } from '@mui/material';
import MesheryResults from '../../components/MesheryResults';
import { updatepagepath } from '../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../lib/path';

function Results({ updatepagepath }) {
  useEffect(() => {
    updatepagepath({ path: getPath() });
  }, [updatepagepath]);

  return (
    <NoSsr>
      <Head>
        <title>Performance Test Results | Meshery</title>
      </Head>
      <Paper>
        <MesheryResults />
      </Paper>
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(Results);
