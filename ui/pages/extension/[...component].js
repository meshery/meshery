import NavigatorExtension from '../../components/NavigatorExtension';
import ExtensionSandbox, {
  getComponentTitleFromPath,
  getComponentIsBetaFromPath,
} from '../../components/ExtensionSandbox';
import { Box, CircularProgress, NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React, { useEffect, useCallback, useState } from 'react';
import RemoteComponent from '../../components/RemoteComponent';
import { MesheryExtensionEarlyAccessCardPopup } from '../../components/General/Popup';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';
import { useRouter } from 'next/router';
import { DynamicFullScrrenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import {
  updateBetaBadge,
  updateCapabilities,
  updateExtensionType,
  updatePagePath,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { useDispatch, useSelector } from 'react-redux';

/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
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

function RemoteExtension() {
  const [componentTitle, setComponentTitle] = useState('');
  const router = useRouter();
  const dispatch = useDispatch();
  const { extensionType } = useSelector((state) => state.ui);
  const { data: capabilitiesRegistry, isLoading } = useGetProviderCapabilitiesQuery();
  const renderExtension = useCallback(() => {
    if (!capabilitiesRegistry?.extensions) return;
    dispatch(updateCapabilities({ capabilitiesRegistry: capabilitiesRegistry }));

    let extNames = [];
    for (var key of Object.keys(capabilitiesRegistry.extensions)) {
      if (Array.isArray(capabilitiesRegistry.extensions[key])) {
        capabilitiesRegistry.extensions[key].forEach((comp) => {
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
        dispatch(updateExtensionType({ extensionType: ext.name }));
        let extensions = ExtensionPointSchemaValidator(ext.name)(
          capabilitiesRegistry.extensions[ext.name],
        );
        setComponentTitle(getComponentTitleFromPath(extensions, getPath()));
        dispatch(updateTitle({ title: getComponentTitleFromPath(extensions, getPath()) }));
        dispatch(updateBetaBadge({ isBeta: getComponentIsBetaFromPath(extensions, getPath()) }));
        dispatch(updatePagePath({ path: getPath() }));
      }
    });
  }, [capabilitiesRegistry, dispatch]);

  useEffect(() => {
    renderExtension();

    return () => {
      dispatch(updateExtensionType({ extensionType: null }));
      setComponentTitle('');
    };
  }, [capabilitiesRegistry, extensionType, router.query.component, renderExtension, dispatch]);

  return (
    <NoSsr>
      <Head>
        <title>{`${componentTitle} | Meshery`}</title>
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

export default RemoteExtension;
