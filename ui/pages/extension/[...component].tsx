import NavigatorExtension from '../../components/layout/Navigator/NavigatorExtension';
import ExtensionSandbox, {
  getComponentTitleFromPath,
  getComponentIsBetaFromPath,
} from '../../components/ExtensionSandbox';
import { NoSsr, useHasPermission } from '@sistent/sistent';
import Head from 'next/head';
import React, { useEffect, useMemo } from 'react';
import RemoteComponent from '../../components/general/RemoteComponent';
import ExtensionPointSchemaValidator from '../../utils/ExtensionPointSchemaValidator';
import { useRouter } from 'next/router';
import { DynamicFullScreenLoader } from '@/components/shared/LoadingState/DynamicFullscreenLoader';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import {
  updateBetaBadge,
  updateExtensionType,
  updatePagePath,
  updateTitle,
  updateProviderCapabilities,
} from '@/store/slices/mesheryUi';
import { useDispatch } from 'react-redux';
import { GetStaticPaths, GetStaticProps } from 'next';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '@/components/general/error-404';

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
  const { data: providerCapabilities, isLoading } = useGetProviderCapabilitiesQuery();
  const canViewExtensions = useHasPermission(Keys.ExtensibilityViewExtensions);

  // Resolve the active extension that matches the current path. Derived from
  // providerCapabilities rather than mirrored into local state to avoid
  // duplicating Redux state and stale-closure bugs. The `typeof window`
  // guard prevents getPath() from crashing during static prerender
  // (this page uses getStaticPaths/getStaticProps).
  const matchedExtension = useMemo(() => {
    if (typeof window === 'undefined') return null;
    if (!providerCapabilities?.extensions) return null;

    const path = getPath();

    for (const key of Object.keys(providerCapabilities.extensions)) {
      const value = providerCapabilities.extensions[key];

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
  }, [providerCapabilities, router.query.component]);

  useEffect(() => {
    if (!providerCapabilities?.extensions) return;

    dispatch(updateProviderCapabilities({ providerCapabilities }));
  }, [dispatch, providerCapabilities]);

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

  if (!canViewExtensions) {
    return <DefaultError permissionKey={Keys.ExtensibilityViewExtensions} />;
  }

  return (
    <NoSsr>
      <Head>
        <title>{`${matchedExtension?.title ?? ''} | Meshery`}</title>
      </Head>

      <DynamicFullScreenLoader isLoading={isLoading}>
        {matchedExtension ? (
          <NoSsr>
            {matchedExtension.name === 'navigator' ? (
              <ExtensionSandbox type={matchedExtension.name} Extension={NavigatorExtension} />
            ) : (
              <ExtensionSandbox
                type={matchedExtension.name}
                Extension={(url) => RemoteComponent({ url })}
              />
            )}
          </NoSsr>
        ) : null}
      </DynamicFullScreenLoader>
    </NoSsr>
  );
}

export default RemoteExtension;
