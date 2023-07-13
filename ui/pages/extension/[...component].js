/* eslint-disable no-unused-vars */

import NavigatorExtension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getComponentTitleFromPath, getComponentIsBetaFromPath } from "../../components/ExtensionSandbox";
import { Box, CircularProgress, NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle, updateExtensionType, updateCapabilities,updatebetabadge } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React from "react";
import RemoteComponent from "../../components/RemoteComponent";
import _ from "lodash"
import { MeshMapEarlyAccessCard } from "../../components/Popup";
import dataFetch from "../../lib/data-fetch";
import ExtensionPointSchemaValidator from "../../utils/ExtensionPointSchemaValidator";
import { withRouter } from "next/router";


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


class RemoteExtension extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      componentTitle : '',
      isLoading : true,
      capabilitiesRegistryObj : null,
    }
  }

  componentWillUnmount() {
    this.setState({
      componentTitle : '',
      isLoading : true,
      capabilitiesRegistryObj : null,
    })
  }

  componentDidMount() {
    dataFetch(
      "/api/provider/capabilities",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        this.props.updatepagepath({ path : getPath() });
        if (result) {
          this.setState({
            capabilitiesRegistryObj : result,
          });
          this.props.updateCapabilities({ capabilitiesRegistry : result })
          this.renderExtension();
        }
      },
      (err) => console.error(err)
    );
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.extensionType !== prevProps.extensionType || this.props.router.query.component != prevProps.router.query.component) {
      this.renderExtension();
    }
  }

  renderExtension = () => {
    let cap = this.props.capabilitiesRegistry;
    // load extension if capabilities are available
    if (cap !== null) {
      let extNames = [];
      for (var key of Object.keys(cap?.extensions)) {
        if (Array.isArray(cap?.extensions[key])) {
          cap?.extensions[key].forEach((comp) => {
            if (comp?.type === "full_page") {
              let ext = {
                name : key,
                uri : comp?.href?.uri
              }
              extNames.push(ext)
            }
          })
        }
      }

      extNames.forEach((ext) => {
        if (matchComponentURI(ext?.uri, getPath())) {
          this.props.updateExtensionType({ extensionType : ext.name });
          let extensions = ExtensionPointSchemaValidator(ext.name)(cap?.extensions[ext.name]);
          this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()), isLoading : false });
          this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
          this.props.updatebetabadge({ isBeta : getComponentIsBetaFromPath(extensions, getPath()) })
        }
      })
    }
    // else, show signup card
    this.setState({ isLoading : false })
  }

  render() {
    const { extensionType, capabilitiesRegistry } = this.props;
    const { componentTitle, isLoading } = this.state;

    return (
      <NoSsr>
        <Head>
          <title>{`${componentTitle} | Meshery` || ""}</title>
        </Head>
        {
          ((this.props.capabilitiesRegistry !== null) && extensionType)?
            (<NoSsr>
              {
                (extensionType === 'navigator') ?
                  <ExtensionSandbox type={extensionType} Extension={NavigatorExtension} />
                  :
                  <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
              }
            </NoSsr>) : (
              !isLoading? (
                <Box display="flex" justifyContent="center">
                  <MeshMapEarlyAccessCard rootStyle={{ position : "relative" }} capabilitiesRegistry={capabilitiesRegistry} />
                </Box>
              ): (
                <CircularProgress />
              )
            )
        }
      </NoSsr>
    )
  }

}

const mapStateToProps = (state) => ({
  extensionType : state.get('extensionType'),
  capabilitiesRegistry : state.get("capabilitiesRegistry")
});

const mapDispatchToProps = (dispatch) => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch),
  updatebetabadge : bindActionCreators(updatebetabadge, dispatch),
  updateExtensionType : bindActionCreators(updateExtensionType, dispatch),
  updateCapabilities : bindActionCreators(updateCapabilities, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(withRouter((RemoteExtension)));


