import { NoSsr } from '@material-ui/core';
import Head from 'next/head';
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import DashboardComponent from '../components/DashboardComponent';
import { getPath } from '../lib/path';
import { updatepagepath } from '../lib/store';

function Index(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Dashboard | Meshery</title>
      </Head>
      <DashboardComponent />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(Index);
