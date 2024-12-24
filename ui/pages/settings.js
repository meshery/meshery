import MesherySettings from '../components/MesherySettings';
import { NoSsr } from '@mui/material';
import { updatepagepath, updatepagetitle } from '../lib/store';
import { connect } from 'react-redux';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from '../lib/path';
import React, { useEffect } from 'react';
import { Paper, styled } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

const PaperStyle = styled(Paper)(() => ({
  maxWidth: '90%',
  margin: 'auto',
  overflow: 'hidden',
}));

function Settings(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
    props.updatepagetitle({ title: 'Settings' });
  }, []);

  return (
    <UsesSistent>
      <NoSsr>
        <Head>
          <title>Settings | Meshery</title>
        </Head>
        <PaperStyle>
          <MesherySettings />
        </PaperStyle>
      </NoSsr>
    </UsesSistent>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
});

export default connect(null, mapDispatchToProps)(Settings);
