import React, { useEffect } from 'react';
import { NoSsr } from '@mui/material';
import { Paper, SistentThemeProvider } from '@layer5/sistent';
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
      <SistentThemeProvider>
        <Paper sx={{ maxWidth: '90%', margin: 'auto', overflow: 'hidden' }}>
          <MesheryResults />
        </Paper>
      </SistentThemeProvider>
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(Results);
