import React from 'react';
import CustomErrorMessage from "../components/ErrorPage";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";

class Error extends React.Component {

  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() })
  }

  render() {
    return (
      <NoSsr>
        <CustomErrorMessage/>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default connect(
  null,
  mapDispatchToProps
)(Error);


