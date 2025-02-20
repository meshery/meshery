/* eslint-disable no-unused-vars */

import NavigatorExtension from '../../components/NavigatorExtension';
import ExtensionSandbox, {
  getComponentTitleFromPath,
  getComponentIsBetaFromPath,
} from '../../components/ExtensionSandbox';
import { Box, CircularProgress } from '@layer5/sistent';
import { NoSsr } from '@layer5/sistent';
import {
  updatepagepath,
  updatepagetitle,
  updateExtensionType,
  updateCapabilities,
  updatebetabadge,
} from '../../lib/store';
import { connect } from 'react-redux';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import React, { useState, useEffect } from 'react';
import RemoteComponent from '../../components/RemoteComponent';
import _ from 'lodash';
import { MesheryExtensionEarlyAccessCardPopup } from '../../components/Popup';
import dataFetch from '../../lib/data-fetch';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';
import { withRouter } from 'next/router';
import { DynamicFullScrrenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';
import { useGetProviderCapabilitiesQuery } from '../../rtk-query/user';

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

function RemoteExtension(props) {
  const [componentTitle, setComponentTitle] = useState('');

  const {
    extensionType,
    capabilitiesRegistry,
    updatepagepath,
    updatepagetitle,
    updatebetabadge,
    updateExtensionType,
    updateCapabilities,
    router,
  } = props;

  const { data: capabilitiesData, isLoading } = useGetProviderCapabilitiesQuery();

  useEffect(() => {
    if (
      extensionType !== props.extensionType ||
      router.query.component !== props.router.query.component
    ) {
      renderExtension(capabilitiesRegistry);
    }
  }, [extensionType, router.query.component, capabilitiesRegistry]);

  const renderExtension = (cap) => {
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
          setComponentTitle(getComponentTitleFromPath(extensions, getPath()));
          updatepagetitle({ title: getComponentTitleFromPath(extensions, getPath()) });
          updatebetabadge({ isBeta: getComponentIsBetaFromPath(extensions, getPath()) });
        }
      });
    }
  };

  useEffect(() => {
    return () => {
      updateExtensionType({ extensionType: null });
      setComponentTitle('');
    };
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>{`${componentTitle} | Meshery` || ''}</title>
      </Head>
      <DynamicFullScrrenLoader isLoading={isLoading}>
        {capabilitiesRegistry !== null && extensionType ? (
          <NoSsr>
            {extensionType === 'navigator' ? (
              <ExtensionSandbox type={extensionType} Extension={NavigatorExtension} />
            ) : (
              <ExtensionSandbox
                type={extensionType}
                Extension={(url) => RemoteComponent({ url })}
              />
            )}
          </NoSsr>
        ) : !isLoading ? (
          <Box display="flex" justifyContent="center">
            <MesheryExtensionEarlyAccessCardPopup
              rootStyle={{ position: 'relative' }}
              capabilitiesRegistry={capabilitiesRegistry}
            />
          </Box>
        ) : (
          <CircularProgress />
        )}
      </DynamicFullScrrenLoader>
    </NoSsr>
  );
}

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
