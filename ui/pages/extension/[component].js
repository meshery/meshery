/* eslint-disable no-unused-vars */

import Extension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getCapabilities, getFullPageExtensions, getComponentTitleFromPath } from "../../components/ExtensionSandbox";
import { NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React from "react";
import RemoteComponent from "../../components/RemoteComponent";


/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * extractComponentURI extracts the last part of the
 * given path
 * @param {string} path
 * @returns {string}
 */
function extractComponentURI(path) {
  return path.substring(path.lastIndexOf("/"));
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

class RemoteExtension extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      componentTitle : "",
      extensionType : ""
    }
  }

  componentDidMount() {
    getFullPageExtensions(extNames => {
      extNames.forEach((ext) => {
        const currentURI = extractComponentURI(getPath());
        if (currentURI === ext.uri) {
          this.setState({ extensionType : ext.name })
          getCapabilities(ext.name, extensions => {
            this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()) });
            this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
          });
          console.log(`path: ${getPath()}`);
          this.props.updatepagepath({ path : getPath() });
        }
      })
    });
  }

  render() {
    const { componentTitle, extensionType } = this.state;
    return (
      <NoSsr>
        <Head>
          <title>{componentTitle || ""}</title>
        </Head>
        {
          extensionType &&
          <NoSsr>
            {console.log("ext", extensionType)}
            {
              (extensionType === 'navigator') ?
                <ExtensionSandbox type={extensionType} Extension={Extension} />
                :
                <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
            }
          </NoSsr>
        }
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch), });

export default connect(null, mapDispatchToProps)(RemoteExtension);
