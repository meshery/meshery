import React from 'react';
import { NoSsr, withStyles } from '@material-ui/core';
import { connect } from 'react-redux';
import Head from 'next/head';
import { EnvironmentComponent } from '../../components/Lifecycle';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

const Environments = () => {
  return (
    <NoSsr>
      <Head>
        <title>Environments | Meshery</title>
      </Head>
      <EnvironmentComponent />
    </NoSsr>
  );
};

export default withStyles(styles)(connect(null)(Environments));
