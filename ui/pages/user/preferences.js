import UserPreferences from '../../components/UserPreferences';
import { NoSsr, Paper, withStyles } from '@material-ui/core';
import { updatepagepath, updatepagetitle } from '../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPath } from '../../lib/path';
import Head from 'next/head';
import { promisifiedDataFetch } from '../../lib/data-fetch';
import { ctxUrl } from '../../utils/multi-ctx';
import React, { useEffect, useState } from 'react';

const styles = { paper: { maxWidth: '90%', margin: 'auto', overflow: 'hidden' } };

const UserPref = (props) => {
  const [anonymousStats, setAnonymousStats] = useState(undefined);
  const [perfResultStats, setPerfResultStats] = useState(undefined);

  useEffect(() => {
    handleFetchData(props.selectedK8sContexts);
  }, [props.selectedK8sContext]);

  useEffect(() => {
    props.updatepagepath({ path: getPath() });
    props.updatepagetitle({ title: 'User Preferences' });
  }, []);

  const handleFetchData = async (selectedK8sContexts) => {
    // console.log(`path: ${getPath()}`);

    try {
      const result = await promisifiedDataFetch(ctxUrl('/api/user/prefs', selectedK8sContexts), {
        method: 'GET',
        credentials: 'include',
      });
      if (result) {
        setAnonymousStats(result.anonymousUsageStats);
        setPerfResultStats(result.anonymousPerfResults);
      }
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      {anonymousStats === undefined || perfResultStats === undefined ? (
        <div></div>
      ) : (
        <NoSsr>
          <Head>
            <title>Preferences | Meshery</title>
          </Head>
          <Paper className={props.classes.paper}>
            {/* {should meshmap specific user preferences be placed along with general preferences or from the remote provider} */}
            <UserPreferences
              anonymousStats={anonymousStats}
              perfResultStats={perfResultStats}
              theme={props.theme}
              themeSetter={props.themeSetter}
            />
          </Paper>
        </NoSsr>
      )}
    </>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
});
const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    selectedK8sContexts,
  };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(UserPref));
