import React from 'react';
import MesherySMIResults from "../components/MesherySMIResults";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';
import { getPath } from "../lib/path";

class SMIResults extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <React.Fragment>
        <Head>
          <title>SMI Results | Meshery</title>
        </Head>
        <MesherySMIResults />
      </React.Fragment>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default connect(
  null,
  mapDispatchToProps
)(SMIResults);