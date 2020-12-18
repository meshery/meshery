import React from "react";
import { connect } from "react-redux";
import { createUseRemoteComponent, getDependencies, createRequires } from "@paciolan/remote-component";
import CircularProgress from "@material-ui/core/CircularProgress";
import { bindActionCreators } from "redux";
import dataFetch from '../lib/data-fetch';
import { updateUser, updateProgress } from '../lib/store';
import { withRouter } from 'next/router';


// const fetcher = (url) => fetch(url, { headers: { "x-auth-token": "welcome123" } }).then((res) => res.text());

const requires = createRequires(getDependencies);

const useRemoteComponent = createUseRemoteComponent({
  requires
});

function RemoteUserPref({ anonymousStats, perfResultStats, startOnZoom, updateProgress, updateUser }) {
  const [loading, err, RemoteComponent] = useRemoteComponent("https://gist.githubusercontent.com/vineethvanga18/7472ef59f3e5afc3c1fad806efe57c33/raw/5593c218a2b6a627ed984a1d341038d0f0baef05/main.js");

  if (loading) {
    return <CircularProgress />;
  }

  if (err != null) {
    return <div>Unknown Error: {err.toString()}</div>;
  }

  return (
    <RemoteComponent
      anonymousStats={anonymousStats}
      perfResultStats={perfResultStats}
      startOnZoom={startOnZoom}
      updateProgress={updateProgress}
      updateUser={updateUser}
      dataFetch={dataFetch}
    />
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateUser: bindActionCreators(updateUser, dispatch)
});

export default connect(
  null,
  mapDispatchToProps,
)(withRouter(RemoteUserPref));