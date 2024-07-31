/* eslint-disable no-unused-vars */

import React, { useState, useEffect } from 'react';
import { Box, CircularProgress, NoSsr } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import Head from 'next/head';
import _ from 'lodash';

import NavigatorExtension from '../../components/NavigatorExtension';
import ExtensionSandbox, {
  getComponentTitleFromPath,
  getComponentIsBetaFromPath,
} from '../../components/ExtensionSandbox';
import RemoteComponent from '../../components/RemoteComponent';
import { MeshMapEarlyAccessCard } from '../../components/Popup';

import {
  updatepagepath,
  updatepagetitle,
  updateExtensionType,
  updateCapabilities,
  updatebetabadge,
} from '../../lib/store';
import dataFetch from '../../lib/data-fetch';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';

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
  const pathSplit = path.split('/');
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

const RemoteExtension = ({
  updateExtensionType,
  updatepagepath,
  updateCapabilities,
  extensionType,
  router,
  capabilitiesRegistry,
  updatepagetitle,
  updatebetabadge,
}) => {
  const [isLoading, setIsLoading] = useState(true);
  const [componentTitle, setComponentTitle] = useState('');
  const [capabilitiesRegistryObj, setCapabilitiesRegistryObj] = useState(null);

  useEffect(() => {
    const fetchCapabilities = () => {
      dataFetch(
        '/api/provider/capabilities',
        {
          method: 'GET',
          credentials: 'include',
        },
        (result) => {
          updatepagepath({ path: getPath() });
          if (result) {
            setCapabilitiesRegistryObj(result);
            updateCapabilities({ capabilitiesRegistry: result });
            renderExtension(result);
          }
        },
        (err) => console.error(err),
      );
    };

    fetchCapabilities();

    return () => {
      updateExtensionType({ extensionType: null });
      setComponentTitle('');
      setIsLoading(true);
      setCapabilitiesRegistryObj(null);
    };
  }, [updateExtensionType, updatepagepath, updateCapabilities]);

  useEffect(() => {
    renderExtension(capabilitiesRegistryObj);
  }, [extensionType, router.query.component, capabilitiesRegistryObj]);

  const renderExtension = (capabilitiesRegistry) => {
    let cap = capabilitiesRegistry;
    // load extension if capabilities are available
    if (cap !== null) {
      let extNames = [];
      for (var key of Object.keys(cap?.extensions)) {
        if (Array.isArray(cap?.extensions[key])) {
          cap?.extensions[key].forEach((comp) => {
            if (comp?.type === 'full_page') {
              let ext = {
                name: key,
                uri: comp?.href?.uri,
              };
              extNames.push(ext);
            }
          });
        }
      }

      extNames.forEach((ext) => {
        if (matchComponentURI(ext?.uri, getPath())) {
          updateExtensionType({ extensionType: ext.name });
          let extensions = ExtensionPointSchemaValidator(ext.name)(cap?.extensions[ext.name]);
          const componentTitle = getComponentTitleFromPath(extensions, getPath());
          setComponentTitle(componentTitle);
          setIsLoading(false);
          updatepagetitle({ title: getComponentTitleFromPath(extensions, getPath()) });
          updatebetabadge({ isBeta: getComponentIsBetaFromPath(extensions, getPath()) });
        }
      });
    } else {
      setIsLoading(false);
    }
  };
  return (
    <NoSsr>
      <Head>
        <title>{`${componentTitle} | Meshery` || ''}</title>
      </Head>
      {capabilitiesRegistry !== null && extensionType ? (
        <NoSsr>
          {extensionType === 'navigator' ? (
            <ExtensionSandbox type={extensionType} Extension={NavigatorExtension} />
          ) : (
            <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
          )}
        </NoSsr>
      ) : !isLoading ? (
        <Box display="flex" justifyContent="center">
          <MeshMapEarlyAccessCard
            rootStyle={{ position: 'relative' }}
            capabilitiesRegistry={capabilitiesRegistry}
          />
        </Box>
      ) : (
        <CircularProgress />
      )}
    </NoSsr>
  );
};

const mapStateToProps = (state) => ({
  extensionType: state.get('extensionType'),
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
  updatebetabadge: bindActionCreators(updatebetabadge, dispatch),
  updateExtensionType: bindActionCreators(updateExtensionType, dispatch),
  updateCapabilities: bindActionCreators(updateCapabilities, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(withRouter(RemoteExtension));
