/* eslint-disable no-unused-vars */

import Extension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getCapabilities, getComponentTitleFromPathForNavigator } from "../../components/ExtensionSandbox";
import { NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";


/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * extractComponentName extracts the last part of the
 * given path
 * @param {string} path
 * @returns {string}
 */
function extractComponentName(path) {
  return path.substring(path.lastIndexOf("/") + 1);
}

/**
 * capitalize capitalizes the given string and returns the modified
 * string
 *
 * If the given parameter is not sting then it will return an empty
 * string
 * @param {string} string
 *
 * @returns {string}
 */
function capitalize(string) {
  if (typeof string === "string") return string.charAt(0).toUpperCase() + string.slice(1);
  return "";
}

class Settings extends React.Component {
  state = { componentTitle : "" }

  componentDidMount() {
    getCapabilities("navigator", extensions => {
      this.setState({ componentTitle : getComponentTitleFromPathForNavigator(extensions, getPath()) });
      this.props.updatepagetitle({ title : getComponentTitleFromPathForNavigator(extensions, getPath()) });
    });
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
  }

  render() {
    return (
      <NoSsr>
        <Head>
          <title>{this.state.componentTitle || ""}</title>
        </Head>
        <NoSsr>
          <ExtensionSandbox type="navigator" Extension={Extension} />
        </NoSsr>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch), });

export default connect(null, mapDispatchToProps)(Settings);
