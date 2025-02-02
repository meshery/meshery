import React from 'react';
import { NoSsr } from '@mui/material';
import { styled } from '@mui/material/styles';
import { connect } from 'react-redux';
import Head from 'next/head';
import MesheryConnections from '../../components/connections';

const StyledPageWrapperDiv = styled('div')({
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
});

const Connections = () => {
  return (
    <NoSsr>
      <Head>
        <title>Connections | Meshery</title>
      </Head>
      <StyledPageWrapperDiv>
        <MesheryConnections />
      </StyledPageWrapperDiv>
    </NoSsr>
  );
};

export default connect(null)(Connections);
