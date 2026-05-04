import React from 'react';
import { NoSsr, styled } from '@sistent/sistent';
import Head from 'next/head';
import MesheryConnections from '../../components/connections';
import { usePageTitle } from '@/utils/hooks';

const StyledPageWrapperDiv = styled('div')({
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
});

const Connections = () => {
  usePageTitle('Connections');

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

export default Connections;
