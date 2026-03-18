import React, { useEffect } from 'react';
import { NoSsr, styled } from '@sistent/sistent';
import Head from 'next/head';
import MesheryConnections from '../../components/connections';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

const StyledPageWrapperDiv = styled('div')({
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
});

const Connections = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Connections' }));
  }, []);

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
