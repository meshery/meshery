import React from "react";
import { NoSsr, withStyles } from "@material-ui/core";
import MesheryPatterns from "../../components/MesheryPatterns";
import { updatepagepath } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from "../../lib/path";
import Wrapper from "../../components/ConfiguratorWrapper"

const styles = {
  paper : {
    maxWidth : '90%',
    margin : 'auto',
    overflow : 'hidden',
  }
};

class Patterns extends React.Component {
  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
  }

  render() {
    return (
      <NoSsr>
        <Head>
          <title>Designs | Meshery</title>
        </Head>
        <Wrapper>
          <MesheryPatterns />
        </Wrapper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(Patterns));