/* eslint-disable no-unused-vars */

import NavigatorExtension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getCapabilities, getFullPageExtensions, getComponentTitleFromPath } from "../../components/ExtensionSandbox";
import { Box, CircularProgress, NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle, updateExtensionType } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React from "react";
import RemoteComponent from "../../components/RemoteComponent";
import { WrappedMeshMapSignupCard } from "../extensions";
import _ from "lodash"
import { MeshMapEarlyAccessCard } from "../../components/Popup";


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
  const pathSplit = path.split("/")
  return pathSplit[pathSplit.length - 1];
}


/**
 * matchComponent matches the extension URI with current
 * given path
 * @param {string} extensionURI
 * @param {string} currentURI
 * @returns {boolean}
 */
function matchComponentURI(extensionURI, currentURI) {
  return currentURI.includes(extensionURI);
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
    super(props);
    this.state = {
      componentTitle : '',
      isLoading : true
    }
  }

  componentDidMount() {
    this.props.updatepagepath({ path : getPath() });
    this.renderExtension();
  }

  componentDidUpdate(prevProps) {
    // re-renders the extension if the extension type (redux store variable) changes
    if (this.props.extensionType !== prevProps.extensionType) {
      this.renderExtension();
    }
  }

  renderExtension = () => {
    getFullPageExtensions(extNames => {
      extNames.forEach((ext) => {
        if (matchComponentURI(ext?.uri, getPath())) {
          this.props.updateExtensionType({ extensionType : ext.name })
          getCapabilities(ext.name, extensions => {
            this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()), isLoading : false });
            this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
          });
        }
      })
    });

    // loading state may set to false, either if the extension
    // is there or no extension after waiting for
    setTimeout(() => {
      this.setState({ isLoading : false })
    }, 1500)
  }

  render() {
    const { extensionType } = this.props;
    const { componentTitle, isLoading } = this.state;

    return (
      <NoSsr>
        <Head>
          <title>{`${componentTitle} | Meshery` || ""}</title>
        </Head>
        {
          extensionType ?
            (<NoSsr>
              {
                (extensionType === 'navigator') ?
                  <ExtensionSandbox type={extensionType} Extension={NavigatorExtension} />
                  :
                  <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
              }
            </NoSsr>) : (!isLoading && (
              <Box display="flex" justifyContent="center">
                <MeshMapEarlyAccessCard rootStyle={{ position : "relative" }} />
              </Box>)
            )
        }
        {
          isLoading && <CircularProgress />
        }
      </NoSsr>
    )
  }

}

const mapStateToProps = (state) => ({
  extensionType : state.get('extensionType'),
});

const mapDispatchToProps = (dispatch) => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch),
  updateExtensionType : bindActionCreators(updateExtensionType, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(RemoteExtension);


