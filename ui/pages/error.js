import React from 'react';
import CustomErrorMessage from "../components/ErrorPage";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";


class Error extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError : false,
      errorMessage : ""
    };
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI.
    return {
      hasError : true,
      errorMessage : error.toString()
    };
  }

  componentDidCatch(error, errorInfo) {
    // We can log the error to an error reporting service but right now it can be console logged
    console.log({ error, errorInfo });
  }

  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() })
  }

  render() {
    const { hasError, errorMessage } = this.state
    const { children } = this.props;

    return hasError ? (
      <NoSsr>
        <CustomErrorMessage {...{ hasError, errorMessage }}/>
      </NoSsr>
    ) : children;
  }
}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default connect(
  null,
  mapDispatchToProps
)(Error);


