import React, { useEffect } from 'react';
import { NoSsr, withStyles } from '@material-ui/core';
import MesheryFilters from '../../components/Filters';
import { updatepagepath } from '../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../lib/path';

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
};

function NewFilters(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Filters | Meshery</title>
      </Head>
      <MesheryFilters />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(NewFilters));
