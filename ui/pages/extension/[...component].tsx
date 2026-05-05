import NavigatorExtension from '../../components/NavigatorExtension';
import ExtensionSandbox, {
  getComponentTitleFromPath,
  getComponentIsBetaFromPath,
} from '../../components/ExtensionSandbox';
import { CircularProgress, NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React, { useEffect, useMemo } from 'react';
import RemoteComponent from '../../components/RemoteComponent';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';
import { useRouter } from 'next/router';
import { DynamicFullScreenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import {
  updateBetaBadge,
  updateCapabilities,
  updateExtensionType,
  updatePagePath,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { useDispatch, useSelector } from 'react-redux';
import { GetStaticPaths, GetStaticProps } from 'next';

/**
 * Define static paths for the extension routes.
 * Required for static export with dynamic catch-all routes.
 */
export const getStaticPaths: GetStaticPaths = async () => {
  return {
    paths: [{ params: { component: ['meshmap'] } }],
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  return { props: { component: params?.component || [] } };
};

function getPath(): string {
  return window.location.pathname;
}

function matchComponentURI(extensionURI: string, currentURI: string): boolean {
  return currentURI.includes(extensionURI);
}

function RemoteExtension() {
  const router = useRouter();
  const dispatch = useDispatch();
  const { extensionType } = useSelector((state) => state.ui);
  const { data: capabilitiesRegistry, isLoading } = useGetProviderCapabilitiesQuery();

  // Resolve the active extension that matches the current path. Derived from
  // capabilitiesRegistry rather than mirrored into local state to avoid
  // duplicating Redux state and stale-closure bugs. The `typeof window`
  // guard prevents getPath() from crashing during static prerender
  // (this page uses getStaticPaths/getStaticProps).
  const matchedExtension = useMemo(() => {
    if (typeof window === 'undefined') return null;
    if (!capabilitiesRegistry?.extensions) return null;
    const path = getPath();
    for (const key of Object.keys(capabilitiesRegistry.extensions)) {
      const value = capabilitiesRegistry.extensions[key];
      if (!Array.isArray(value)) continue;
      for (const comp of value) {
        if (
          comp?.type === 'full_page' &&
          comp?.href?.uri &&
          matchComponentURI(comp.href.uri, path)
        ) {
          const extensions = ExtensionPointSchemaValidator(key)(value);
          return {
            name: key,
            title: getComponentTitleFromPath(extensions, path),
            isBeta: getComponentIsBetaFromPath(extensions, path),
          };
        }
      }
    }
    return null;
  }, [capabilitiesRegistry, router.query.component]);

  useEffect(() => {
    if (!capabilitiesRegistry?.extensions) return;
    dispatch(updateCapabilities({ capabilitiesRegistry }));
  }, [dispatch, capabilitiesRegistry]);

  useEffect(() => {
    if (!matchedExtension) return;
    dispatch(updateExtensionType({ extensionType: matchedExtension.name }));
    dispatch(updateTitle({ title: matchedExtension.title }));
    dispatch(updateBetaBadge({ isBeta: matchedExtension.isBeta }));
    dispatch(updatePagePath({ path: getPath() }));

    return () => {
      dispatch(updateExtensionType({ extensionType: null }));
    };
  }, [dispatch, matchedExtension]);

  return (
    <NoSsr>
      <Head>
        <title>{`${matchedExtension?.title ?? ''} | Meshery`}</title>
      </Head>
      <DynamicFullScreenLoader isLoading={isLoading}>
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
        ) : !isLoading ? null : (
          <CircularProgress />
        )}
      </DynamicFullScreenLoader>
    </NoSsr>
  );
}

export default RemoteExtension;
