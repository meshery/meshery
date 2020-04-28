import React from 'react';
import MesheryPlayComponent from '../components/MesheryPlayComponent';
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import { getPath } from "../lib/path";

class Manage extends React.Component {
  static getInitialProps = ({ query }) => {
    return { query };
  }
  constructor(props) {
    super(props);
    if (!props.query.adapter) {
      const urlParams = new URLSearchParams(window.location.search);
      this.props.query.adapter = urlParams.get("adapter");
    }
  }
  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  componentDidUpdate() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render() {
    return (
      <NoSsr>
        <React.Fragment>
          <MesheryPlayComponent adapter={this.props.query.adapter} />
        </React.Fragment>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => {
  return {
    updatepagepath: bindActionCreators(updatepagepath, dispatch)
  }
}

export default connect(
  null,
  mapDispatchToProps
)(Manage);