import { NoSsr } from "@material-ui/core";
import Head from 'next/head';
import React from "react";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import DashboardComponent from "../components/DashboardComponent";
import { getPath } from "../lib/path";
import { updatepagepath } from "../lib/store";

class Index extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Dashboard | Meshery</title>
        </Head>
        <DashboardComponent />
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default connect(
  null,
  mapDispatchToProps
)(Index);