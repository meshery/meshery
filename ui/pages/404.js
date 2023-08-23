import React, { useEffect } from 'react';
import CustomErrorMessage from "../components/ErrorPage";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";
import Head from 'next/head';

const Error = () => {
  useEffect(() => {
    console.log(`path: ${getPath()}`);
    updatepagepath({ path : getPath() });
  }, [updatepagepath]);

  return (
    <NoSsr>
      <Head>
        <title>404 - Page Not Found</title>
      </Head>
      <CustomErrorMessage />
    </NoSsr>
  )
}

const mapDispatchToProps = dispatch => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch)
});

export default connect(
  null,
  mapDispatchToProps
)(Error);