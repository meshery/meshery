import React from 'react';
import { NoSsr, withStyles } from '@material-ui/core';
import { connect } from 'react-redux';
import Head from 'next/head';
import { WorkspacesComponent } from '../../components/Lifecycle';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

const Workspaces = () => {
  return (
    <NoSsr>
      <Head>
        <title>Workspaces | Meshery</title>
      </Head>
      <WorkspacesComponent />
    </NoSsr>
  );
};

export default withStyles(styles)(connect(null)(Workspaces));
