import React from "react";
import { NoSsr } from "@material-ui/core";
import PerformanceProfiles from "../../components/MesheryPerformance/PerformanceProfiles";
import { updatepagepath } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import Head from "next/head";
import { getPath } from "../../lib/path";

class Results extends React.Component {
  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
  }

  render() {
    return (
      <NoSsr>
        <Head>
          <title>Performance Profiles | Meshery</title>
        </Head>
        <PerformanceProfiles />
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch), });

export default connect(null, mapDispatchToProps)(Results);
