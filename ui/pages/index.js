import React from "react";
import DashboardComponent from "../components/DashboardComponent";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";
import Head from 'next/head';

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