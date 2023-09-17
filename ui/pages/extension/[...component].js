/* eslint-disable no-unused-vars */

import NavigatorExtension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getComponentTitleFromPath, getComponentIsBetaFromPath } from "../../components/ExtensionSandbox";
import { Box, CircularProgress, NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle, updateExtensionType, updateCapabilities,updatebetabadge } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React,{ useEffect, useState, useRef } from "react";
import RemoteComponent from "../../components/RemoteComponent";
import _ from "lodash"
import { MeshMapEarlyAccessCard } from "../../components/Popup";
import dataFetch from "../../lib/data-fetch";
import ExtensionPointSchemaValidator from "../../utils/ExtensionPointSchemaValidator";
import { withRouter,useRouter } from "next/router";



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


const RemoteExtension=(props) => {

  const [componentTitle, setComponentTitle] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [capabilitiesRegistryObj, setCapabilitiesRegistryObj] = useState(null);
  const router = useRouter();
  const prevProps = useRef();

  useEffect(() => {
    prevProps.current = props;
  },[props])

  useEffect(() => {
    return () => {
      setComponentTitle('');
      setIsLoading(true);
      setCapabilitiesRegistryObj(null);
    }
  },[])

  useEffect(() => {
    dataFetch(
      "/api/provider/capabilities",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        props.updatepagepath({ path : getPath() });
        if (result) {
          setCapabilitiesRegistryObj(result);
          props.updateCapabilities({ capabilitiesRegistry : result })
          renderExtension();
        }
      },
      (err) => console.error(err)
    );
  },[])

  useEffect(() => {
    if (props.extensionType !== prevProps.current.extensionType || router.query.component != prevProps.current.router.query.component) {
      renderExtension();
    }
  }, [props.extensionType, router.query.component]);



  const renderExtension = () => {
    const router = useRouter();
    let cap = props.capabilitiesRegistry;

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
          props.updateExtensionType({ extensionType : ext.name });
          let extensions = ExtensionPointSchemaValidator(ext.name)(cap?.extensions[ext.name]);
          setComponentTitle(getComponentTitleFromPath(extensions, getPath()));
          setIsLoading(false);
          props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
          props.updatebetabadge({ isBeta : getComponentIsBetaFromPath(extensions, getPath()) });
        }
      })
    }
    // else, show signup card
    setIsLoading(false);
  }




  return (
    <NoSsr>
      <Head>
        <title>{`${componentTitle} | Meshery` || ""}</title>
      </Head>
      {
        ((props.capabilitiesRegistry !== null) && props.extensionType)?
          (<NoSsr>
            {
              (props.extensionType === 'navigator') ?
                <ExtensionSandbox type={props.extensionType} Extension={props.NavigatorExtension} />
                :
                <ExtensionSandbox type={props.extensionType} Extension={(url) => RemoteComponent({ url })} />
            }
          </NoSsr>) : (
            !isLoading? (
              <Box display="flex" justifyContent="center">
                <MeshMapEarlyAccessCard rootStyle={{ position : "relative" }} capabilitiesRegistry={props.capabilitiesRegistry} />
              </Box>
            ): (
              <CircularProgress />
            )
          )
      }
    </NoSsr>
  );
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


