/* eslint-disable no-unused-vars */

import NavigatorExtension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getCapabilities, getFullPageExtensions, getComponentTitleFromPath } from "../../components/ExtensionSandbox";
import { Box, CircularProgress, NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle, updateExtensionType, updateCapabilities } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React from "react";
import RemoteComponent from "../../components/RemoteComponent";
import { WrappedMeshMapSignupCard } from "../extensions";
import _ from "lodash"
import { MeshMapEarlyAccessCard } from "../../components/Popup";
import { CapabilitiesRegistry } from "../../utils/disabledComponents";
import dataFetch from "../../lib/data-fetch";
import ExtensionPointSchemaValidator from "../../utils/ExtensionPointSchemaValidator";


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
      isLoading : true,
      capabilitiesRegistryObj : null,
      path : ''
    }
  }

  componentDidMount() {
    console.log("componentDidMount")
    console.log("this.props.page:  ", this.props.page)
    dataFetch(
      "/api/provider/capabilities",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        console.log(result)
        if (result) {
          this.setState({
            capabilitiesRegistryObj : result,
          });
          //global state
          this.props.updateCapabilities({ capabilitiesRegistry : result })
          this.renderExtension();
        }
      },
      (err) => console.error(err)
    );
    this.props.updatepagepath({ path : getPath() });
  }

  componentDidUpdate(prevProps, prevState) {
    console.log("componentDidUpdate")
    console.log("prevProps", prevProps)
    console.log("this.props", this.props)
    // re-renders the extension if the extension type (redux store variable) changes
    if (this.props.extensionType !== prevProps.extensionType) {
      this.renderExtension();
    }
  }

  renderExtension = () => {
    console.log("render Ext")
    let result = this.state.capabilitiesRegistryObj;
    let extNames = [];
    for (var key of Object.keys(result?.extensions)) {
      if (Array.isArray(result?.extensions[key])) {
        result?.extensions[key].forEach((comp) => {
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
    console.log("extNames", extNames)
    extNames.forEach((ext) => {
      if (matchComponentURI(ext?.uri, getPath())) {
        console.log("match found");
        this.props.updateExtensionType({ extensionType : ext.name })
        let type = ext.name;
        let extensions = ExtensionPointSchemaValidator(type)(result?.extensions[type])
        this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()), isLoading : false });
        this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
        // getCapabilities(ext.name, extensions => {
        //   this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()), isLoading : false });
        //   this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
        // });
      }
    })
    // getFullPageExtensions(extNames => {
    //   extNames.forEach((ext) => {
    //     if (matchComponentURI(ext?.uri, getPath())) {
    //       this.props.updateExtensionType({ extensionType : ext.name })
    //       getCapabilities(ext.name, extensions => {
    //         this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()), isLoading : false });
    //         this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
    //       });
    //     }
    //   })
    // });

    // loading state may set to false, either if the extension
    // is there or no extension after waiting for
    // setTimeout(() => {
    //   this.setState({ isLoading : false })
    // }, 1500)
  }

  render() {
    const { extensionType } = this.props;
    const { componentTitle, isLoading } = this.state;
    console.log("extensionType: ", extensionType)
    console.log("isLoading: ", isLoading);

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
            </NoSsr>) : (
              !isLoading? (
                <Box display="flex" justifyContent="center">
                  <MeshMapEarlyAccessCard rootStyle={{ position : "relative" }} />
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
});

const mapDispatchToProps = (dispatch) => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch),
  updateExtensionType : bindActionCreators(updateExtensionType, dispatch),
  updateCapabilities : bindActionCreators(updateCapabilities, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(RemoteExtension);


