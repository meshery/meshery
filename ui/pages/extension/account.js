/* eslint-disable no-unused-vars */

import ExtensionSandbox, { getCapabilities, getComponentTitleFromPathForAccount } from "../../components/ExtensionSandbox";
import { NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import RemoteAccount from "../../components/RemoteAccount";


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

class Account extends React.Component {
  state = { componentTitle : "" }

  componentDidMount() {
    getCapabilities("account", extensions => {
      this.setState({ componentTitle : getComponentTitleFromPathForAccount(extensions, getPath()) });
      this.props.updatepagetitle({ title : getComponentTitleFromPathForAccount(extensions, getPath()) });
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
          <ExtensionSandbox type="account" Extension={(url) => RemoteAccount({ url })} />
        </NoSsr>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch), });

export default connect(null, mapDispatchToProps)(Account);